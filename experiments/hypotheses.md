# Week 7 hypotheses (frozen before any Week 7 run)

This file was committed before any Week 7 experiment ran. The commit time is
the proof that the predictions came first. If something must change later,
it goes under "Changes after freezing" at the bottom, with a date and a
reason. Nothing above that heading is edited once results exist.

## Words used here

- off: `OF_HETERO_PLACEMENT=off`. ETH's routing, unchanged.
- strict: our filter. A node whose reported hardware does not fit the job is
  removed before ETH picks a node.
- score: new this week. strict, then keep only the node with the most
  reported free GPU memory.
- Wrong placement: a job sent to a node whose reported hardware does not fit.
- False refusal: a 503 from our filter although some node fits.
- Each setting runs 5 times. A difference counts only if every run of one
  side is better than every run of the other side (the ranges do not overlap).

## H1. The filter stops wrong placements (RQ1, RQ4)

I expect strict to send 0 jobs to nodes that cannot run them, while off sends
them about as often as chance. When all nodes are the same, strict should
change nothing for jobs every node can run.

- Baseline: off (ETH's random pick).
- Intervention: strict. Both heads use the same workers at the same time.
- Jobs: CPU-only, GPU-required, capability-specific and unsupported (Day 2).
- Clusters: same-hardware, mixed-architecture, mixed-GPU (Day 3).
- Metrics: wrong placements out of jobs with needs; false refusals; which
  detail (arch, compute capability, memory, no GPU) caused each refusal.
- Prediction: strict 0 wrong and 0 false refusals in every run. off close to
  chance (the share of nodes that do not fit). Week 6 demo: off 92 of 120,
  strict 0 of 120.
- Wrong if: strict has any wrong placement or false refusal in any run, or
  off is not worse than strict in the mixed clusters (then the filter fixes
  nothing).

## H2. Free-memory scoring helps only while the numbers are true (RQ1, RQ3)

ETH's docs suggest scoring nodes by free GPU memory (scheduling.mdx:177).
OpenTela reads a node's GPU memory at start (node_table.go:725) and again
only on a few events: a new connection while it offers a service
(host.go:328), a relay change (host.go:809), a service health change
(service_health.go:184, :195). All of these call ReannounceLocalServices
(registrar.go:247-253). There is no timer, so if a GPU fills up later, other
nodes keep seeing the old number.

I expect score to finish more jobs than random and shortest-queue when the
reported memory is true. When one node's GPU has filled up after it
reported, I expect score to finish fewer jobs than shortest-queue, because
it keeps choosing that node.

- Baselines: strict + random (ETH default) and strict + shortest-queue
  (`--lb-policy shortest-queue`, loadbalancer.go:62-117; it counts requests
  in flight per node, not memory). Both use strict, so only the choice
  among fitting nodes differs.
- Intervention: score.
- Jobs: each job holds 3000 MiB of emulated GPU memory on the stand-in
  runner for 2 s. A runner answers 507 if the job does not fit on one of its
  GPUs. ETH retries only on transport errors (proxy_handler.go:1017-1021), so
  a 507 is final. 8 jobs are sent at once, 10 times per run.
- Workers: rtx (live, about 7000 MiB free: room for 2 jobs) and t4 (saved
  profile, 2 GPUs with 15360 MiB free each: room for 10).
- Settings: "true" = the t4 runner has the room its report says. "stale" =
  the t4 runner acts as if another program took 12000 MiB on each GPU after
  the report (room for 2), while the report still says 15360 free.
- Metric: jobs finished (not 507) out of 8 per round.
- Prediction (per round of 8): true: score 8, shortest-queue about 6,
  random about 6. stale: score 2, shortest-queue about 4, random about 4.
- Wrong if: in "true", score does not finish more jobs than both baselines;
  or in "stale", score finishes as many jobs as shortest-queue or more.

## H3. The extra information is cheap (RQ2)

I expect the cost of our fields and filter to be small next to what
OpenTela already does.

- Baseline: off, and ETH's node record without our 4 fields.
- Intervention: strict (and score), node record with our 4 fields.
- Metrics and predictions:
  - Node record size: +103 bytes with 1 GPU, +129 with 2 GPUs (already
    measured in Week 6, so this is a repeat check, not a prediction).
  - Filter time (Go benchmark of filterByHardware, fake table with 2, 10,
    100 and 1000 nodes): under 1 ms per request at 1000 nodes.
  - End-to-end: median time of the same placed jobs, strict minus off:
    under 1 ms.
  - Live reading: our extra `nvidia-smi` call (name, compute_cap) on my PC:
    median under 500 ms over 20 calls. It runs only at start and on the
    events listed in H2, not per request.
- Wrong if: any of the three new numbers is above its limit.

## Rules for every run

- Fixed job lists and fixed seeds, written before the run (Day 2, Day 3).
- Raw outputs are saved by scripts and never edited by hand.
- Emulated nodes (saved profiles) and emulated memory are labelled as such.
  Jobs do not use the real GPU.
- All nodes run on one PC with `OF_SECURITY_REQUIRE_SIGNED_BINARY=false`,
  like ETH's own integration tests.
- No raw node tables or node logs in the repo (they hold the owner key and
  the hostname).

## Not tested this week

Network bytes of the CRDT traffic; a timer that re-reads GPU memory; how long
a new node sees older nodes without hardware (seen in the Week 6 demo); an
arm64 or AMD machine running as a node (a profile captured on a real arm64
Mac may be used as a saved profile); more than one PC.

## Changes after freezing

(none yet)
