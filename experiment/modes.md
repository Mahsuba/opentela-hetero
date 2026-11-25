# Week 6 Day 6: two modes, OpenTela's routing vs mine

## What I wanted
A fair test needs a reference. My reference is OpenTela's own routing,
unchanged. I ran it next to my change, at the same time, with the same
workers and the same jobs. Only one setting was different.

## The switch
A head node reads `OF_HETERO_PLACEMENT`:
- `off` (also the default): OpenTela's routing. My code gives the worker
  list back untouched and never reads the job's needs.
- `strict`: before picking, the head drops every worker whose hardware
  cannot run the job (Day 5 `Compatible`). If nobody is left, it answers
  503 with a reason for each worker.
- anything else: error 500, so a typo cannot quietly run the wrong mode.

A job states its needs with three request headers: `X-Otela-Need-Arch`,
`X-Otela-Need-Min-CC`, `X-Otela-Need-Min-Free-MiB`. No header = no need.
This copies how OpenTela's `X-Otela-Allowed-Peers` already removes
workers (proxy_handler.go:750-764).

Code: ETH fork, branch hetero/capability-fields, commit e64bd1e.
New files src/internal/server/hetero_placement.go and hetero_router.go
(+ tests). In proxy_handler.go I added 6 lines right after the "No
provider found" check. No OpenTela line was removed.

Tests on my PC: my 7 new tests pass, 4 of them through OpenTela's real
request handler and node table. Full OpenTela suite: 18 ok, 4 with no
test files, 0 FAIL (same as Day 5). gofmt and go vet clean.

## Setup (my PC, 2026-11-25)
| node | seed | HTTP port | what it is |
|---|---|---|---|
| head-off | 1 | 8092 | `OF_HETERO_PLACEMENT=off` |
| head-strict | 4 | 8392 | `OF_HETERO_PLACEMENT=strict` |
| worker rtx | 2 | 8192 | live reading: RTX 5060, cc 12.0, 8151 MiB, 1101 used |
| worker t4 | 3 | 8292 | saved Kaggle capture: 2x Tesla T4, cc 7.5, 15360 MiB, 0 used |

The t4 worker is emulated: the profile is a real capture (2026-11-20),
but the node ran on my PC. Both workers connected to both heads, and
before the run both heads' node tables showed the same two workers.
Each worker passes jobs to a stand-in runner (experiment/backend.py)
that replies with its own name and does no GPU work. Every node had
`OF_SECURITY_REQUIRE_SIGNED_BINARY=false`, like OpenTela's own
integration tests. The picker was OpenTela's default, random.

Jobs: experiment/jobs.tsv (5 jobs). experiment/run-jobs.sh sent each
job 10 times to each head. Raw output, not edited:
experiment/results/off.tsv and strict.tsv (50 requests each).

## Result
| job | needs | can really run it | off | strict |
|---|---|---|---|---|
| any | none | both | 6 rtx, 4 t4 | 6 rtx, 4 t4 |
| cc8 | cc 8.0 | rtx | 8 rtx, 2 t4 | 10 rtx |
| free10g | 10000 MiB free | t4 | 5 rtx, 5 t4 | 10 t4 |
| cc8_free10g | both | nobody | 5 rtx, 5 t4 | 10 x 503 |
| arm64 | arm64 CPU | nobody | 6 rtx, 4 t4 | 10 x 503 |

A wrong placement is a job sent to a worker whose reported hardware
cannot run it. Off: 27 of the 40 jobs with needs. Strict: 0 of 40. The
job with no needs behaved the same in both modes. I wrote these
predictions down before each run, and every one matched.

Strict's 503 replies name the reason, for example
`"cpu_arch amd64, job wants arm64"` for both workers on the arm64 job,
and `"no GPU meets min_cc and min_free_mib"` on cc8_free10g (rtx has
the cc but not the memory, t4 the memory but not the cc).

Cross-check: for every delivered job, the worker named by the head
(`X-Computing-Node`) and the name in the runner's reply agreed: off
50 of 50, strict 30 of 30.

Workbook check, both modes are runnable: met.

## Limits
- 200 means delivered, not run. The stand-in runner does no GPU work,
  so the 27 wrong placements did not fail here. On real hardware they
  would fail or run badly later. That is not measured.
- One PC. The t4 worker is emulated.
- Strict trusts what each worker reports, and a reading is from when
  the worker started. The same GPU showed 915, 991 and 1101 MiB used
  in three nodes' readings.
- 10 tries per job with a random picker, so the off counts change from
  run to run. Week 7 needs more runs and the shortest-queue picker too.
- Strict treats unknown values as "no". A lenient mode is not built.
- The head's log showed about 0.1 ms per refused request. That is a
  hint, not a measurement.
