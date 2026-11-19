# Fixtures (real captures, not made up)

## rtx5060-wsl.csv
- Source: real nvidia-smi output on Mahi's PC (hostname Genie), Ubuntu WSL on Windows
- GPU: NVIDIA GeForce RTX 5060
- Date: 2026-11-19
- Command: nvidia-smi --query-gpu=name,memory.total,memory.used,compute_cap,driver_version --format=csv,noheader,nounits
- Columns: name, memory.total (MiB), memory.used (MiB), compute_cap, driver_version
- Note: driver_version is the Windows driver (WSL borrows it from Windows).
- Note: memory.used changes all the time (859 then 836 MiB, one minute apart).
