# Week 4 · Day 6 — Code map: distributed propagation

Repo: eth-easl/opentela @ 4b6e90f (read-only). All paths under `src/`.

**Question:** When computer A learns its own hardware, how does that info reach computer B's NodeTable, where B's request picker can see it?

**Short answer:** Through a CRDT shared key-value store. A writes its whole `Peer` (with `Hardware` inside) as JSON under key `"/<peerID>"`. The store announces the change on a pubsub topic. B fetches it and, in a hook, replaces its whole copy of A except `Connected`. By default this is the only road for hardware info. SWIM is off by default and only carries alive/dead.

---

## 1. The engine: copied Merkle-CRDT library

- `internal/protocol/go-ds-crdt/` is a copy of the IPFS **go-ds-crdt** library. It is *not* ETH-written: the header in `go-ds-crdt/crdt.go:1-18` cites the "Merkle-CRDTs" paper by Sanjuán, Pöyhtäri and Teixeira.
- What it is: a replicated key-value store. Every change is a "commit" in a Merkle-DAG (like Git).
- It has two helpers:
  - **Broadcaster**: announces "my newest commit (head) is X."
  - **DAGSyncer**: fetches commits a replica is missing (like `git pull`).
- Claims from the header comment (not measured by me):
  - New replicas need time to fully sync, depending on how fast they fetch the DAG.
  - Newer values may arrive before older ones.
  - Without compaction the store only grows, even when keys are deleted. This is why `tombstone_compactor.go` and `tombstone_manager.go` exist.

## 2. ETH's setup of the engine: `internal/protocol/crdt.go`

| Line(s) | What | Meaning |
|---|---|---|
| 26 | `pubsubTopic = "ocf-crdt"` | All nodes share one pubsub topic (one "group chat"). |
| 120-124 | publish `"ping"` every 20 s | Keep-alive message on the topic. Carries no data. (The condition that starts this loop is above line 120 and not checked.) |
| 131 | `crdt.NewPubSubBroadcaster(..., pubsubTopic)` | The broadcaster is attached to that topic. |
| 135-145 | `DAGSyncerTimeout` = 30 s default; env `OF_CRDT_DAG_SYNCER_TIMEOUT` overrides | Comment: reduced from 5 min "to prevent cascading blockage". A failed fetch is retried at the next rebroadcast. |
| 149 | `MultiHeadProcessing = true` | Incoming heads are processed in parallel, so one slow peer (e.g. NAT'd) doesn't stall the others. |
| 150-155 | `RebroadcastInterval`: **5 s** by default; **60 s** if `scalability.crdt_tuned=true` (also `NumWorkers` 16) | How often a node re-announces its current head. |
| 161-162 | `MaxFetchFailures = 50`, `MinFetchFailureAge = 6h` | Give up on a commit that can't be fetched (orphan). Banned commits stay banned after restart (see `go-ds-crdt/autoban.go`). |
| 163-193 | **PutHook**: runs on B for every incoming put | Turns the JSON back into a `Peer`. Skips B's own key (174-177). Unknown peer → `Connected=false` (180-182). Known peer → **takes the whole new Peer but keeps B's own `Connected`** (183-185). Then `UpdateNodeTableHook` (189) writes it into the NodeTable. |
| 194-197 | **DeleteHook** | A deleted key leads to `DeleteNodeTableHook` (the row is removed). |
| 199 | `crdt.New(store, ..., ipfs, pubsubBC, opts)` | Store, broadcaster and options are joined together here. |
| 210 | `startTombstoneCompactor(crdtStore)` | Background cleanup of deletion markers. |

Defaults verified in `entry/cmd/root.go`:
- `:198` `scalability.swim_enabled = false`
- `:199` `scalability.crdt_tuned = false` → 5 s rebroadcast by default
- `:227` `crdt.tuned_rebroadcast_interval = "60s"`

## 3. Sending side: every place that writes into the shared store

Found with `grep -rn "\.Put(" internal/protocol/*.go` (test files excluded). Each line is listed with the function it sits in.

| Put at | Function | What it writes |
|---|---|---|
| node_table.go:745 | `InitializeMyself` (639-751) | Own Peer at startup. Hardware is filled at :725 (Day 4). |
| registrar.go:279 | `ReannounceLocalServices` (247-end) | Own Peer again. Hardware is re-filled at :253/:266 (Day 4). |
| registrar.go:223 | `provideService` (205-230) | Own Peer after adding a service. |
| node_table.go:321 | `UpdateNodeTable` (290-325) | General update of own entry. |
| node_table.go:347 | `MarkSelfAsBootstrap` (326-352) | Marks self as bootstrap node. |
| node_table.go:374 | `AnnounceLeave` (353-380) | Leaving the network (churn case). |
| node_table.go:786 | `RegisterRemotePeer` (775-788) | A **remote** peer's entry. Only caller: `server/registration.go:211`. |
| liveness.go:335 | `publishEviction` (317-352) | Marks a peer as evicted (dead). |

Assumption (not yet read line by line): these Puts write `json.Marshal(Peer)`. It's consistent with the PutHook decoding JSON into `Peer`.

## 4. Relay registration strips Hardware (on purpose)

`internal/server/registration.go`:
- :193-196 rejects anything whose role is not `relay`.
- :198 comment: "clear fields that are **unverified** or must not come from the caller".
- :204 forces role = `["relay"]`.
- :206-208 sets `Service = nil`, `Load = nil`, `Hardware = {}`. Comment at :205: "fields that only the node itself should set".
- :211 `protocol.RegisterRemotePeer(req.Peer)` → CRDT.

My reading: this is a **deliberate trust decision, not a placement bug**. Relays carry no services, so the picker never routes jobs to them.

Open (low priority): if the relay later Puts its own Peer, whichever write comes last wins (PutHook replaces the whole Peer).

## 5. SWIM: optional second channel, membership only

- `node_table.go:34-35`: `InitScalableNodeTable` (copy-on-write node table plus SWIM) runs only when `scalability.swim_enabled=true`. That is **false by default** (root.go:198).
- `node_table.go:57-67` `StartSWIM`: all config fields are probe and suspect timings (`ProbeInterval`, `ProbeTimeout`, `IndirectProbes`, `SuspectTimeout`, `RetransmitMult`). Events are `swim.MemberEvent`.
- My reading: SWIM decides alive/suspect/dead. It does **not** carry Hardware. Not verified: what fields `MemberEvent` contains.

---

## 6. The full path (default config)

```
A: fill Hardware (node_table.go:725 / registrar.go:253,266)
 → A: store.Put("/<A>", Peer-as-JSON)       (node_table.go:745 / registrar.go:279)
 → engine: new commit; head announced on "ocf-crdt" (crdt.go:131)
 → head re-announced every 5 s              (crdt.go:154)
 → B: fetch missing commits, 30 s timeout   (crdt.go:145)
 → B: PutHook: replace whole Peer, keep own Connected (crdt.go:163-193)
 → B: NodeTable row updated                 (UpdateNodeTableHook, node_table.go:381)
 → B's picker reads the NodeTable           (GetAllProviders, node_table.go:611 — Day 5)
```

## 7. Check: where must new capability data flow?

- New capability fields (CPU arch, GPU vendor, compute capability, CUDA version) belong in `common.HardwareSpec` / `GPUSpec` (common/structs.go:19-30, Day 3).
- They must be filled in where Hardware is already filled (node_table.go:725, registrar.go:253/266).
- After that they **ride the existing road with no new propagation code**: the whole Peer is JSON, and the PutHook copies every field except `Connected`.
- Relay registration will strip them, the same way it strips the current Hardware.

## 8. What this means for the research questions (hypotheses, not results)

- **RQ2 (cost):** Every Put sends the whole Peer. Bigger HardwareSpec = bigger payload on every update. The store also only grows until compaction. Measure bytes per Put and total store growth. ETH already has `server/controlplane_overhead_test.go` (not read yet) as a possible starting point.
- **RQ3 (staleness):** Keep two delays apart:
  1. **Delivery delay**: bounded by rebroadcast (5 s default / 60 s tuned) plus fetch time (timeout 30 s).
  2. **Content age**: new hardware values are only written on startup or re-announce. These are event-driven, with no timer (Day 4: host.go:328, :809; service_health.go:184/195). So `UsedMemory` can be old even if delivery is fast. `integration/stale_head_test.go` exists (not read yet).
- **Experiment knobs that need no code change:** `scalability.crdt_tuned` (5 s vs 60 s), env `OF_CRDT_DAG_SYNCER_TIMEOUT`.
- **Trust:** Hardware is self-reported. ETH's own word is "unverified" (registration.go:198). Richer capability info means more claims nobody checks. `internal/attestation/attestation.go` exists (not read yet).


