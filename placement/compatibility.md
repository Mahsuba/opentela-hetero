# Week 6 Day 5: Can this node run the job?

## What I built
Before OpenTela picks a node for a job, we now have a way to ask "can this node run it?".
The code is on my OpenTela fork, branch `hetero/capability-fields`, commit `ecea7be`.
It is a new folder `src/internal/placement/` with two files:
- `compatibility.go` (91 lines): the check `Compatible(job, node)`.
- `compatibility_test.go` (74 lines): 23 test cases.

Nothing in OpenTela calls the check yet. Day 6 connects it to the router.

## The rule
A job can ask for up to three things. An empty one means "no need".
- `Arch`: CPU type, for example `amd64`.
- `MinCC`: lowest GPU compute capability, for example `8.0`.
- `MinFreeMiB`: free GPU memory needed, in MiB.

The check answers yes or no. Every "no" comes with a reason.
1. No needs: always yes. This is exactly how OpenTela works today.
2. CPU type: the node's `cpu_arch` must be the same. If the node does not say, the answer is no.
3. GPU needs: the node's `gpu_status` must be `ok`. `none_found` gives "no GPU". Anything else (`tool_missing`, `query_failed`, or empty from a node without our fields) gives "gpu_status unknown".
4. One GPU must meet every GPU need, because the job runs on one GPU.
5. Compute capability is compared as numbers, not as text.
6. Free memory = total - used. A broken reading never passes.

## Why these choices
- Unknown means no. This is the strict rule from design/capability-model.md (Day 1). A lenient mode (unknown means yes) can be an experiment setting on Day 6.
- As text, "12.0" comes before "8.0". A text comparison would call my RTX 5060 (12.0) too old for an 8.0 job.
- A node with one new but full GPU and one old but empty GPU cannot run a job that needs "new and empty". Each need alone is met, but not by the same GPU.
- The reasons let Week 7 count why nodes were refused. That helps answer RQ4 (which details matter).

## Evidence
On my PC (Genie, Go 1.26.8, Ubuntu WSL):
- `CGO_ENABLED=0 go test -count=1 -v ./internal/placement/` -> PASS.
- `gofmt -l` printed nothing. `go vet` -> no problems.
- Broken on purpose (compute capability compared as text) -> FAIL on "cc 12.0 >= 8.0 (as text it would fail)". Then the good file was put back (fingerprint c9686a01, same as before).
- All OpenTela tests: 18 ok, 4 with no tests, 0 FAIL. The baseline was 17 ok; the one extra is the new folder.

In Claude's sandbox, 5 more broken versions were each caught by a test: unknown CPU type allowed, `tool_missing` treated as `ok`, each need met by a different GPU, the "no needs means yes" step removed, negative used memory allowed.

The test nodes: 2 are real readings (my RTX 5060: 12.0, 8151 MiB total, 704 used; Kaggle 2x Tesla T4: 7.5, 15360 MiB). The other 8 are hand-written, not captures.

## Limits
- The router does not use the check yet, and a request cannot state its needs yet (Day 6).
- The check trusts what each node says about itself. Nothing verifies it.
- Used memory changes over time. The check sees the last value the node shared (H3).
- OpenTela's AMD reading gives no compute capability, so a job with `MinCC` refuses AMD nodes.
