# Week 5 Day 5 - Fail safely

## What OpenTela does now (checked at 4b6e90f)
OpenTela reads GPUs in `src/internal/platform/gpu.go`.
- Line 19 runs nvidia-smi with `exec.Command`, which has no time limit. I did not check if the code that calls `GetGPUInfo()` adds one.
- Lines 21-23: any nvidia-smi error (not installed, broken driver, no permission) gives a Debug log and an empty list.
- Lines 11-15: if the NVIDIA list is empty, it tries rocm-smi (line 48). Lines 51-53 and 57-58 return an empty list again.

So the node table only ever sees `[]`. OpenTela cannot tell "this node has no GPU" from "this node's GPU tool is broken". The reason is only in a Debug log, which is hidden by default.

## What I did
I wrote a profiler that always prints one schema-v1 JSON (week5/day04/schema.json), even when something fails.
- `profiler/errors.go`: the 4 states (ok, none_found, tool_missing, query_failed). A failure always has `gpus: []` and a plain `gpu_error`. Any full GPU uuid inside an error message is cut to `GPU-` + 8 hex.
- `profiler/gpu.go`: live mode (nvidia-smi, 10 s time limit + 2 s WaitDelay) or fixture mode (a saved file).
- `profiler/parse.go`: all or nothing. One bad line gives query_failed, never a half list. `[N/A]` becomes null.
- `profiler/main.go`: builds the JSON. Exit code is 0 even with no GPU, because no GPU is a normal answer.
- `test-failures.sh`: fake nvidia-smi scripts that fail on purpose.

## What we saw (Genie, Ubuntu WSL, 2026-11-20)
| Case | Real or simulated | gpu_status | Schema |
|---|---|---|---|
| live run | real | ok (used 852 MiB, Day 4 was 962) | PASS |
| `env PATH=` (tool hidden) | real program, same PC | tool_missing | PASS |
| no devices | simulated | none_found | PASS |
| no driver / no permission / garbage / leaky error | simulated | query_failed | PASS |
| [N/A] fields | simulated | ok, with nulls | PASS |
| frozen tool | simulated | query_failed after 12 s | PASS |

A broken copy (status ok, no GPUs) FAILS. The fake uuid tail never shows up in the outputs.

Two bugs were found in testing and fixed before committing: uuid_short was always null, and the time limit waited 30 s instead of 10 s.

## Limits
- none_found and query_failed were only simulated. The exact "No devices were found" text was not checked on a real machine without a GPU.
- tool_missing was my PC with nvidia-smi hidden, not a separate computer.
- AMD reading is not written yet. If rocm-smi is found, it honestly reports query_failed.
- Old drivers without the `compute_cap` field would make the whole query fail. The Kaggle T4 capture will test this.
- One GPU only.

## Why it matters
A scheduler that knows "no GPU" vs "GPU tool broken" can send CPU jobs to the first and treat the second as unknown (RQ1). A frozen nvidia-smi without a time limit could delay a node (RQ2). This is a measured difference, not a claim that OpenTela is wrong.

## Extra: real 2-GPU fixture (Kaggle)
I captured 2x Tesla T4 on Kaggle (compute capability 7.5, driver 580.178.04) into fixtures/kaggle-t4x2.csv. The profiler read it in fixture mode: ok with 2 GPUs, schema PASS. This is the first real multi-GPU test. See fixtures/README.md for the fixture time/CPU limit.
