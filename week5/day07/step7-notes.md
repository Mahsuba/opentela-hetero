# Week 5 Day 7 - resourceprofiler v1

## What I made
I packaged the Day 6 profiler as its own tool in resourceprofiler-v1/,
so someone with only the GitHub link can build, test and run it.
It has its own go.mod (module resourceprofiler, go 1.24), a README,
the schema v1 and its checker (copied from week5/day04), test data,
three example outputs and a small script, capture.sh.

## The bug I fixed
On Day 5 I found that fixture mode used the time and CPU of the machine
that replays the file, not of the machine that captured it. A Kaggle
file from Nov 20 replayed today said "today", and a file from an ARM
machine would have said amd64. That would break H1.
Now a fixture starts with two lines:
  # collected_at=2026-11-20T20:37:34Z
  # arch=x86_64
The profile takes the time and CPU from these lines (x86_64 becomes the
Go name amd64, aarch64 becomes arm64). A fixture without them is refused
with exit code 2 and no JSON. The tool does not guess.

## capture.sh
A read-only script that saves the header lines plus the raw nvidia-smi
lines on another machine. It cuts full GPU uuids to 8 hex characters.
If nvidia-smi is missing or fails, it saves nothing and exits 1, so a
broken tool can never be saved as "no GPU".
I ran it on my PC: RTX 5060, arch x86_64, 0 lines with the cut-off
uuid piece, 1 line with the short id (positive control). Replaying the
capture gave a profile that PASSES the schema.

## Tests
12 tests, all PASS on my Go 1.26.8. Five are new (fixture_test.go):
real Kaggle file keeps Nov 20 + amd64; a made-up aarch64 file says
arm64 on my amd64 PC; 5 broken headers are refused; a missing file is
refused; error line numbers match the file.
Mutation check: in a copy in /tmp I put the Day 5 bug back. Two tests
failed. The Kaggle test only caught the time, because Kaggle and my PC
are both amd64. The aarch64 test caught both time and arch. So both
tests are needed.

## Example outputs (examples/, all PASS the schema)
- live-rtx5060.json: live on my PC, ok
- fixture-kaggle-t4x2.json: Kaggle 2x T4, collected_at 2026-11-20T20:37:34Z
- no-gpu-tool-missing.json: PATH emptied, tool_missing
Wrong use (2 arguments) and the old Day 6 fixture without header both
gave exit code 2 with a clear message.

## Clean checkout (workbook check: MET)
Commit ff12208. I cloned it fresh into /tmp/clean: go test ok, build ok,
live run PASS, fixture run PASS, 0 files changed after building.
Claude did the same on another machine with Go 1.24.7: tests ok,
fixture run PASS, 0 files changed (no GPU there, so no live run).

## Limits
NVIDIA only (AMD detected, not read). A GPU name with a comma is
rejected. Old drivers without compute_cap fail the whole query. A
fixture can only store a working nvidia-smi reading. Hardware is
self-reported.

## Week 5 deliverable
A tested resourceprofiler-v1 with a stable JSON schema (v1, unchanged
since Day 4) and honest handling of unavailable hardware: no GPU, no
tool, a broken tool, and missing values are all reported, never made up.
