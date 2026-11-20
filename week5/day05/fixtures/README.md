# Day 5 fixtures

## kaggle-t4x2.csv
- What: real nvidia-smi output from a Kaggle notebook (accelerator "GPU T4 x2"), not my PC.
- Captured: 2026-11-20T20:37:34Z (date -u in the notebook).
- CPU of that machine: uname -m = x86_64 (= amd64).
- Query: nvidia-smi --query-gpu=name,memory.total,memory.used,compute_cap,driver_version,uuid,pci.bus_id --format=csv,noheader,nounits
- Privacy: each uuid was cut to GPU- + 8 hex with sed BEFORE it was shown. The full uuid was never saved.
- used = 0 MiB on both cards: idle server GPUs with no desktop on them.
- The two cards have the same name, memory, compute capability and driver. Only uuid_short and pci_bus_id tell them apart.

## Known limit
In fixture mode the profiler writes collected_at = the time it was run, and cpu.arch = the CPU of the computer running it (Genie), not the capture machine. Here both are x86_64, so the JSON is still right except collected_at. To be fixed: fixture files must carry their own time and CPU.

## Real CPU-only machine (Kaggle, accelerator None)
- 2026-11-20T20:41:19Z, uname -m = x86_64, and nvidia-smi gave "command not found".
- That is exactly the case where our profiler says tool_missing (nvidia-smi is not on PATH).
- The profiler itself was not run there (no Go on Kaggle). So this backs up the PATH= test on Genie but does not replace it.
