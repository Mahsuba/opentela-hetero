# Week 6 Day 7: the demo

## What this shows
Two heads on one PC send the same 50 jobs to the same 2 workers. Head-off uses ETH's normal routing (OF_HETERO_PLACEMENT=off). Head-strict uses our hardware filter (OF_HETERO_PLACEMENT=strict). That switch is the only difference between the two heads.

Workers:
- rtx: my RTX 5060, read live: amd64, compute capability 12.0, 8151 MiB in total.
- t4: the Kaggle capture (2x Tesla T4, 7.5, 15360 MiB each), loaded from a saved profile. This node is emulated: the hardware was on Kaggle, the jobs ran on my PC.

Jobs (experiment/jobs.tsv, 10 tries each): any (no needs), cc8 (compute capability 8.0 or more), free10g (10000 MiB free), cc8_free10g (both), arm64.

## How to run
From the repo root, with /tmp/otela built from branch hetero/capability-fields (e64bd1e):

    sh demo/run-demo.sh run1

It starts 2 stand-in runners and 4 nodes in the background, waits until both heads see both workers, sends the jobs to each head, marks the results, and stops everything. One run takes about 20 seconds. Results go to demo/runs/RUN/. Logs stay in /tmp/d7 and are not in the repo, because they contain the node owner key and the PC name.

Files:
- run-demo.sh: the recipe.
- helpers.sh: wait, stop, and read a head's table.
- check.sh: marks each job against an answer key (which worker's reported hardware fits it).

## Results (2026-11-26, 3 runs)
| run | off: wrong placements | strict: wrong placements |
|---|---|---|
| run1 | 30 of 40 | 0 of 40 |
| run2 | 28 of 40 | 0 of 40 |
| run3 | 34 of 40 | 0 of 40 |
| total | 92 of 120 | 0 of 120 |

A wrong placement is a job with needs that went to a worker whose reported hardware does not fit it. In every run there were 0 fitting jobs refused, 0 other replies, and 0 mismatches between the node header and the runner's reply. The strict head logged 40 filter lines (one per job with needs), the off head 0.

Strict gave the same rows in all 3 runs: cc8 10 to rtx, free10g 10 to t4, cc8_free10g and arm64 10 refusals (503) each. Only "any" changed, because both workers fit it. Off changed every run because ETH's default picker is random. If chance alone decides, about 30 of 40 should be wrong: all 20 jobs nobody can run, plus about half of the 20 that only one worker can run. 92 of 120 is close to that.

Before the runs, check.sh scored the Day 6 hand-run results as 27 of 40 and 0 of 40, the same as experiment/modes.md.

## Something we saw that we did not plan
In all 3 runs, head-strict's table showed head-off with no hardware ("[] None None []"), while head-off's table showed head-strict with its full hardware. In ETH's code, protocol/host.go:332-343 creates a minimal entry (only the id) when a node connects to a peer it does not know yet; the full record replaces it later. The tables were saved about 8 seconds after head-strict started, which matches head-strict still holding only the minimal entry for the node that was there before it. This did not change the results, because heads offer no service. It matters for Week 7: a node that just joined can, for a short time, see others with no hardware. Strict mode refuses such a node (unknown = no). How long it lasts is not measured yet.

## Limits
- 200 means the job was delivered to the runner, not that it ran on a GPU. The runner is a stand-in that only answers with its name.
- One PC. The t4 worker is emulated from a real Kaggle capture.
- The answer key and the filter both trust what the workers report.
- Memory is read when a node starts (rtx reported 810 MiB used in run1).
- 3 runs of 50 jobs with a random picker in off mode. Week 7 needs more runs and shortest-queue as a second baseline.
- Every node runs with OF_SECURITY_REQUIRE_SIGNED_BINARY=false, as ETH's own integration tests do.
