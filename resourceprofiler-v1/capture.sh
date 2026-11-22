#!/bin/sh
# capture.sh - save THIS machine's GPU reading as a fixture file.
# Read-only: it only asks nvidia-smi, date and uname.
# Full GPU uuids are cut to GPU- + 8 hex (privacy rule, Day 3).
# Usage: sh capture.sh > my-machine.csv
# If nvidia-smi is missing or fails, it prints NOTHING to the file and
# exits 1, so a broken tool can never be saved as "no GPU".
cut='s/(GPU-[0-9a-fA-F]{8})[-0-9a-fA-F]+/\1/g'
if ! command -v nvidia-smi >/dev/null 2>&1; then
  echo "capture.sh: nvidia-smi not found, nothing saved" >&2; exit 1
fi
if ! out=$(nvidia-smi --query-gpu=name,memory.total,memory.used,compute_cap,driver_version,uuid,pci.bus_id \
  --format=csv,noheader,nounits 2>&1); then
  echo "capture.sh: nvidia-smi failed, nothing saved: $(echo "$out" | head -1 | sed -E "$cut")" >&2; exit 1
fi
echo "# collected_at=$(date -u +%Y-%m-%dT%H:%M:%SZ)"
echo "# arch=$(uname -m)"
echo "$out" | sed -E "$cut"
