# Week 4 · Day 7 — Research gap memo

Baseline: eth-easl/opentela @ `4b6e90f`.
Evidence: week4/day03–day06 code maps + fleet-manager doc (read on Day 7).

## 1. How to read this memo

- **Fact** = seen in code or docs, with file:line.
- **Candidate gap** = a guess about something missing. Not proven.
- **Hypothesis** = a guess written so an experiment can prove it right or wrong.

Rule: nothing goes into Week 5 unless it is a fact or clearly labelled as a guess.

## 2. Already exists (facts)

### Node info (Day 3)
- `Peer` struct — `protocol/node_table.go:162` (ID, Status, LastSeen, Labels, Hardware).
- `HardwareSpec` / `GPUSpec` — `common/structs.go:19-30`. GPU name, total memory, used memory.
  No CPU architecture, GPU vendor, compute capability, or driver/CUDA version.

### Discovery (Day 4)
- `platform/gpu.go`: tries `nvidia-smi` first (:19); `rocm-smi` (:49) only if no NVIDIA.
  Reads name, memory.total, memory.used. On failure: empty list, Debug log only.
- Host Memory / MemoryBandwidth / UsedMemory are never written → always 0. CPU info is not in `Peer`.
- Labels come from `--label` (`entry/cmd/root.go:71`), parsed in `node_table.go:251-270`. Never verified.
- Nodes re-share their info only when something happens (`host.go:328`, `:809`; `service_health.go:184/195`). No regular timer.

### Sharing between nodes (Day 6)
- Engine: Merkle-CRDT (vendored IPFS `go-ds-crdt`). Topic `"ocf-crdt"` (`crdt.go:26`).
- Rebroadcast every 5 s by default, 60 s if `scalability.crdt_tuned` (`crdt.go:150-155`).
- `PutHook` (`crdt.go:163-193`) replaces the whole `Peer` from JSON.
  → Any new `HardwareSpec` field would spread to other nodes with no new sharing code.
- Relay registration clears Hardware (`server/registration.go:208`) — deliberate (`:198` "unverified", relay-only).

### Choosing a node for a request (Day 5) — files under `src/internal/server/`
- Pool: `GetAllProviders` (`protocol/node_table.go:611-631`) — only checks Connected + service name.
- `selectCandidates` (`proxy_handler.go:400-462`): exact identity group > `"*"` > `"all"`; `X-Otela-Fallback` 0/1/2, default = exact only.
- Final pick (`proxy_handler.go:919-942`): affinity > price-aware > weighted_routing > `lb.Pick`.
  Default policy = random (`entry/cmd/root.go:75`). weighted_routing off (`root.go:200`).
- `scoreCandidates` (`proxy_handler.go:502-512`) returns a constant 1.0 and only gets node IDs.
  Its comment calls it an extension point (latency, load, trust). Hardware is not listed.
- Hardware/Labels in `server/`: used only in `analytics_hook.go:77-79` (stats) and `registration.go:208` (clear). Zero uses in `proxy_handler.go`.

### Heterogeneity elsewhere in the repo
- ETH's hardware-aware scheduler exists in the **Python simulator** only:
  `contrib/serving-simulator/simulator/core/scheduler.py:402` (H100/A100/3090 tiers, out-of-memory filtering). No Go file uses it.
- Fleet manager doc (`docs/content/docs/extensions/fleet-manager.mdx`):
  - `:253` each preset defines SLURM parameters.
  - `:260` `gpus` = count or type, e.g. `4` or `"rtx_3090:1"`. No GPU-memory field.
  - `:223` cluster `arch` = `amd64` or `arm64`, written by hand.
  - Only file named "fleet" in the repo → the fleet manager code probably lives elsewhere (not checked).

## 3. Not a gap (by design)

- **Picker trusts the model name.** For nodes launched by the fleet manager, a human preset + SLURM already put the model on suitable GPUs. So the picker skipping hardware is most likely deliberate.
- **Relay clears Hardware.** Deliberate (relay-only, unverified).
- We do **not** claim "OpenTela ignores hardware". It collects hardware, shares it, and uses GPU type at deploy time.

## 4. Candidate gaps (guesses — not proven)

- **G1 — Outsider nodes.** Nodes that join without the fleet manager, or requests that fall back to `"*"`/`"all"`, get no hardware check in the Go request path.
  Unknown: what `"*"`/`"all"` nodes actually serve; how common self-joined nodes are.
- **G2 — Live GPU load not used.** Used GPU memory is collected and shared, but the picker is random / constant score.
- **G3 — Old news.** Info is only re-shared on events, so "how full is this GPU" can be any age. This is separate from delivery delay (rebroadcast + fetch).
- **G4 — Missing identity fields.** CPU arch, GPU vendor, compute capability, driver/CUDA version are not in live node info. (CPU arch is known at deploy time, `fleet-manager.mdx:223`, but not shared.)
- **G5 — Self-reported hardware.** Nothing verifies it. `attestation/attestation.go` exists (unread). Out of scope for now.

## 5. Hypotheses (testable)

### H1 — Outsider case (RQ1, RQ4)
- **Guess:** In a cluster mixing fleet-launched and self-joined nodes (or with fallback on), the current picker sends some requests to nodes that can't run them. A hardware-compatibility filter cuts this to about zero with very small extra pick time.
- **Proven wrong if:** baseline already has about zero wrong picks.
- **Measure:** % incompatible placements, request success rate, extra pick latency.

### H2 — Busy GPU case (RQ1, RQ4)
- **Guess:** With uneven load, picking by free GPU memory (instead of random) lowers the slowest 5% of response times (p95) and the reject rate.
- **Proven wrong if:** no improvement, or worse because everyone piles onto the same "free" node (herding on old info).
- **Measure:** p50/p95 latency, rejects/out-of-memory errors.

### H3 — Old news case (RQ2, RQ3)
- **Guess:** Info age grows with time since the last event. A regular re-share every T seconds keeps it around T + delivery delay, at about (number of nodes × Peer size ÷ T) extra bytes per second.
- **Proven wrong if:** traffic doesn't grow like that, or fresher info doesn't help H2.
- **Measure:** info age, bytes/sec per node, effect on the H2 results.

## 6. Checks still open before Week 5

- [ ] What the `"*"` and `"all"` identity groups mean in practice.
- [ ] Whether self-joined (non-fleet) nodes are a supported setup.
- [ ] Where the fleet manager code lives.
- [ ] Skim `attestation/attestation.go` (decide scope).
- [ ] Optional: `stale_head_test.go`, `controlplane_overhead_test.go` (may give a method for measuring H3).

## 7. Rule check

Every fact above has a file:line from `4b6e90f` or the fleet-manager doc. Everything else is labelled as a guess.
