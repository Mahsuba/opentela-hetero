#!/bin/bash
# test-failures.sh - Day 5: make nvidia-smi fail ON PURPOSE and check the
# profiler still prints valid schema-v1 JSON.
# The fake nvidia-smi scripts are SIMULATED failures, not real broken machines.
# Usage (from repo root): bash week5/day05/test-failures.sh /tmp/profiler
set -u
P=${1:-/tmp/profiler}
OUT=week5/day05/examples
FAKE=$(mktemp -d)
UUID="GPU-9877209e-1111-2222-3333-444455556666" # fake tail, not the real uuid

fake() { # fake <case> <what the fake nvidia-smi does>
  mkdir -p "$FAKE/$1"
  printf '#!/bin/sh\n%s\n' "$2" > "$FAKE/$1/nvidia-smi"
  chmod +x "$FAKE/$1/nvidia-smi"
}
fake no-devices    'echo "No devices were found"; exit 6'
fake no-driver     'echo "NVIDIA-SMI has failed because it could not communicate with the NVIDIA driver."; exit 9'
fake no-permission 'echo "Failed to initialize NVML: Insufficient Permissions"; exit 4'
fake garbage       'echo "hello world"'
fake not-available "echo \"NVIDIA Test GPU, 15360, 0, [N/A], [N/A], $UUID, 00000000:00:04.0\""
fake leaky-error   "echo \"Unknown Error on $UUID\"; exit 15"
fake frozen        '/bin/sleep 30'

for c in no-devices no-driver no-permission garbage not-available leaky-error frozen; do
  start=$(date +%s)
  env PATH="$FAKE/$c" "$P" > "$OUT/simulated-$c.json"
  code=$?
  status=$(grep -o '"gpu_status": "[a-z_]*"' "$OUT/simulated-$c.json" | cut -d'"' -f4)
  echo "$c: exit=$code status=$status time=$(( $(date +%s) - start ))s"
done
rm -rf "$FAKE"

echo "--- schema check (all must PASS)"
python3 week5/day04/validate.py week5/day04/schema.json $OUT/simulated-*.json

echo "--- negative control (must FAIL): status ok but no GPUs"
sed 's/"none_found"/"ok"/' $OUT/simulated-no-devices.json > /tmp/broken.json
python3 week5/day04/validate.py week5/day04/schema.json /tmp/broken.json

echo "--- leak check: fake uuid tail 1111-2222 must NOT appear"
if grep -q "1111-2222" $OUT/simulated-*.json; then echo "LEAK"; else echo "no full uuid"; fi
