# Week 5 Day 6 - Tests

## Goal
Tests that check the profiler without needing a real GPU, so the result is the
same on any computer. Workbook check: `go test ./...` passes.

## What I did
I copied the Day 5 profiler into week5/day06/profiler and made a go.mod
(`go mod init opentela-hetero/week5/day06`) so `go test ./...` can run.
Then I added 3 test files with 9 tests:
- parse_test.go: the real Kaggle 2x T4 file (2 GPUs, Tesla T4, 15360 MiB,
  compute 7.5, two different uuid_short and pci_bus_id), a table of 10
  made-up inputs, and a check that [N/A] becomes null and a full uuid is cut.
- errors_test.go: how a broken nvidia-smi is labelled (time limit, no
  devices, error text, no output), the uuid privacy cut, and that a failure
  gives "gpus": [] and never null.
- gpu_test.go: fake nvidia-smi scripts in a temporary folder (works, no
  devices, no driver, no permission, garbage), no tool at all, and a missing
  fixture file. These are simulated failures. Needs Linux (/bin/sh).

## What the tests found
The first run failed 1 test: "used more than total" (100 MiB total, 200 used)
was accepted as ok. My Day 4 schema notes said "used <= total" was not
enforced yet. I added 3 lines to parse.go: if used > total, the answer is
query_failed. Then all 9 tests passed.

## Proof the tests work
I broke a copy of the code in /tmp on purpose (removed the lowercase step for
uuids). TestNotAvailableAndPrivacy failed with "got GPU-ABCDEF12". So the tests
really catch this bug. My real parse.go was not changed.
The real program still works: live run on Genie and the Kaggle fixture both
PASS the Day 4 schema.

## How to run
```
cd week5/day06
go test ./...
go build -o /tmp/profiler ./profiler
```

## Link to OpenTela (commit 4b6e90f)
- OpenTela has GPU tests: src/internal/platform/gpu_test.go (5 tests).
- Most of them read the real machine. TestGetNvidiaGPUs (:29-33) skips itself
  when nvidia-smi is missing, so on a computer without an NVIDIA GPU the
  NVIDIA reading code is not checked. Only TestCsvField uses fixed input.
- In gpu.go:18-46, running nvidia-smi and reading its text are one function,
  so a test cannot give it saved text. In my profiler, parse() is separate.
- gpu.go:37-38 ignore number errors (`total, _ := strconv.ParseInt(...)`), so
  a broken memory number becomes 0. There is no used vs total check.
- On Genie (CGO_ENABLED=0, because my PC has gcc but no C headers):
  4 PASS, TestGetAMDGPUs SKIP (no rocm-smi). git status stayed clean.
- This is not a bug claim. Today the router never reads Hardware. It matters
  in Week 6, when a scheduler uses free memory = total - used.

## Limits
- One real fixture (Kaggle). The frozen case (12 s) is only in test-failures.sh.
- Fixture mode still writes the run time and the CPU of the computer running
  it, not the capture machine. Fix on Day 7.
