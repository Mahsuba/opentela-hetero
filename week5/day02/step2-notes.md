# W5D2 - GPU inventory (step2.go)

## Idea
The GPU is a separate device from the CPU. It has its own name, memory,
generation (compute capability) and driver. Asking the CPU (step1) tells
us nothing about the GPU. For NVIDIA we ask the GPU's own tool: nvidia-smi.

## What step2.go does
- No argument: runs nvidia-smi now, labels output "LIVE".
- File argument: reads a saved real capture, labels output "FIXTURE <file>".
- Splits each line at commas, trims spaces, prints 5 values per GPU.
- If nvidia-smi is missing: prints "no NVIDIA GPU found" + the reason. Never makes up numbers.

## Query
nvidia-smi --query-gpu=name,memory.total,memory.used,compute_cap,driver_version --format=csv,noheader,nounits

## Real results (Mahi's PC, Ubuntu WSL, 2026-11-19)
- name: NVIDIA GeForce RTX 5060
- memory total: 8151 MiB (not 8192 = "8 GB" on the box; 41 MiB less)
- memory used: 859 -> 836 -> 892 MiB within a few minutes (Windows desktop uses the GPU)
- compute capability: 12.0 (Blackwell)
- driver: 616.92 (Windows driver; WSL borrows it. Kubuntu uses a different driver, 595 series)

## Tests
- LIVE: go run week5/day02/profiler/step2.go -> real values
- FIXTURE: go run week5/day02/profiler/step2.go week5/day02/fixtures/rtx5060-wsl.csv -> same values, used = 836 (frozen)
- No GPU: go build -o /tmp/step2 ...; env PATH= /tmp/step2 -> "no NVIDIA GPU found (exec: "nvidia-smi": executable file not found in $PATH)"

## What I learned
1. Total memory is fixed; used memory changes every minute -> goes stale fast (H2, H3, RQ3).
2. A job needing 7.5 GB fits by total (8151) but not by free (~7259) -> scheduler needs used/free (H2).
3. What the GPU reports depends on where you ask from (WSL vs Kubuntu driver) -> every fixture needs a label.
4. Fixture = a photo of one moment. Fine for repeatable tests if labelled.

## Link to OpenTela (verified W4D3/W4D4 at 4b6e90f)
- platform/gpu.go:19 asks nvidia-smi only for name, memory.total, memory.used.
- GPUSpec (common/structs.go:19-30) has no field for compute capability or driver version.
- So compute capability is "available from nvidia-smi but not collected", same pattern as GOARCH in W5D1.

## Limits
- One machine, one GPU. No AMD hardware, so no AMD fixture (not made up on purpose).
