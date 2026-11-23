# Week 6 Day 3: profiler hook (filling the capability fields)

OpenTela code: fork Mahsuba/OpenTela, branch hetero/capability-fields,
commit 98bdc25 (on top of Day 2's 8293bf0, base 4b6e90f). Not sent to ETH.

## What was missing

On Day 2 I added 4 empty fields (cpu_arch, gpu_status, collected_at,
compute_capability). Nothing filled them, so a node still reported only
GPU name and memory.

## What I changed

OpenTela writes a node's hardware in two places:
- at start: src/internal/protocol/node_table.go:725
- at re-announce: src/internal/protocol/registrar.go:253 and :266

Before, both did `Hardware.GPUs = platform.GetGPUInfo()`. Now both do
`Hardware = platform.GetHardwareSpec()`. That is 3 changed lines in ETH's
files. All new code is in 3 new files in src/internal/platform/:
capability.go, profile_file.go, capability_test.go.

GetHardwareSpec has two modes:
- Live (normal): GPU name and memory still come from ETH's own
  GetGPUInfo, unchanged. I add cpu_arch (Go's runtime.GOARCH), the time
  (UTC), gpu_status, and compute capability from a second nvidia-smi
  question (name,compute_cap, 10 s limit). Compute capability is filled
  only if the GPU names match ETH's list one to one; otherwise it stays
  empty (unknown).
- Saved profile: if OF_HETERO_PROFILE points to a resourceprofiler-v1
  JSON file, the node reports that file instead. This is for Week 7
  emulated nodes. A broken file gives gpu_status query_failed with no
  GPUs and no arch, so it never pretends to be the real PC.

## What we saw (Genie, RTX 5060, 2026-11-23)

- platform tests: our 4 PASS, ETH's 4 PASS + AMD SKIP (as before).
- All OpenTela tests: 17 ok, 4 no test files, 0 FAIL (same as Day 2).
- Real node, private (`--mode standalone`, bootstrap_connected=false),
  asked with `curl localhost:8092/v1/self`:
  RTX 5060, 8151 / 720 MiB, compute_capability 12.0, cpu_arch amd64,
  gpu_status ok, collected_at 2026-11-23T14:38:29Z (same second as
  "Server started").
- Same node with OF_HETERO_PROFILE = fixture-kaggle-t4x2.json:
  2x Tesla T4, 15360 MiB, 7.5, amd64, ok, collected_at
  2026-11-20T20:37:34Z (Kaggle capture time). RTX 5060 not shown.

Workbook check met: a node can report its capability profile.

## Limits

- Live mode never says none_found. ETH's GetGPUInfo cannot tell "no GPU"
  from "tool broken", so if a GPU tool exists but returns nothing, I say
  query_failed (unknown).
- Live mode runs nvidia-smi twice per reading (ETH's + mine). Extra cost
  to measure in Week 7 (RQ2).
- A saved profile keeps its own collected_at, so emulated nodes look
  old. Day 6 decides if experiments need "now" instead (H3).
- Both runs had the same peer id even with a new --datadir. Day 4 needs
  two different identities (check --seed / --ephemeral-identity).
- Without --mode standalone, a node looks up ETH's public network
  (root.go:50). Always use standalone for experiments.
