# Week 6 Day 4: propagation (can other nodes see the new fields?)

OpenTela code: fork Mahsuba/OpenTela, branch hetero/capability-fields,
commit 98bdc25 (same as Day 3). Today I changed no OpenTela code.

## The question

A scheduler can only use hardware facts that reach it. On Day 3 a node
could report its own 4 new fields. Today's question: does another node
receive them? In Week 4 (code-map-propagation.md) I wrote, from reading
the code, that new HardwareSpec fields travel inside the Peer record
with no new sharing code. Today I tested that on two running nodes.

## Setup (one PC, private)

- Node A: live reading of my RTX 5060. `--mode standalone --seed 1`,
  web port 8092, tcp 43905.
- Node B: saved Kaggle profile (OF_HETERO_PROFILE =
  fixture-kaggle-t4x2.json). `--seed 2`, web port 8192, tcp 43915,
  `--bootstrap.static` = only node A's address. With that set, ETH's
  public list is not used (bootstrap.go:80-83).
- Each node got its own `--config-dir` and `--datadir`. A node's id comes
  from a saved key file, and the seed is used only when no key file
  exists yet (host.go:123-146). With one shared key file both nodes
  would have the same id, which is what happened on Day 3.
- Ids: A = 12D3KooWHHzSe... , B = 12D3KooWKRyzV... . Both match the ids
  I worked out from the seeds before running.

## First try: both tables empty

`curl localhost:PORT/v1/dnt/table` gave `{}` on both nodes, but
`/v1/dnt/peers_status` showed A and B Connected to each other. So the
network link worked, but neither node put the other in its table.

Reason, from the code: `security.require_signed_binary` is true by
default (root.go:154). A node then rejects any peer whose program is not
signed by ETH (node_table.go:402-407), and logs this only at debug
level. I built my program myself, so it is not signed. ETH's own
multi-node tests turn this off with OF_SECURITY_REQUIRE_SIGNED_BINARY=false
(peer_discovery_test.go:196, stale_head_test.go:291). I did the same.
Note: ETH's CLAUDE.md:79 says this check is off by default for
development, but the code and docs/configuration.md:90 say true.

## Second try: each node sees the other (2026-11-23)

Same commands, plus OF_SECURITY_REQUIRE_SIGNED_BINARY=false on both.

Node A's table shows B:
12D3KooWKRyz True amd64 ok 2026-11-20T20:37:34Z
[('Tesla T4', 15360, '7.5'), ('Tesla T4', 15360, '7.5')]

Node B's table shows A:
12D3KooWHHzS True amd64 ok 2026-11-23T18:05:58Z
[('NVIDIA GeForce RTX 5060', 8151, '12.0')]

B's row is the saved Kaggle profile. A's row is a live reading, taken
the same second A started (10:05:58 local = 18:05:58 UTC).

Workbook check met: two local nodes can observe each other's profile,
through OpenTela's own shared table, with no new sharing code.

## Size of the new fields (first number for RQ2)

For each row, I counted the bytes of the row as compact JSON, then
again with our 4 fields removed:

| Row seen | Row bytes | Without our fields | Added |
|---|---|---|---|
| B (2 GPUs) seen by A | 765 | 636 | 129 (+20%) |
| A (1 GPU) seen by B | 699 | 596 | 103 (+17%) |

103 is exactly the Day 1 estimate. Each extra GPU adds about 27 bytes.
I predicted 103 and 129 before running.

## Limits

- Bytes are counted on the row as the table shows it, not the bytes
  sent over the network. Real traffic cost is for Week 7.
- I did not measure how long a change takes to reach the other node.
  That is RQ3 (Week 7).
- Signed-program check turned off. Real ETH nodes keep it on, so only
  signed builds would share these fields there.
- Both nodes run on one PC. B's hardware is a saved Kaggle profile
  (emulated node), not a second machine.
- Hardware is self-reported. A node can claim anything; nothing checks it.
- B (normal mode) prints "Connectivity incomplete ... relay_reserved=false;
  retrying" (connectivity.go:324). It tries to reserve a relay node and
  there is none. It did not stop the sharing.
