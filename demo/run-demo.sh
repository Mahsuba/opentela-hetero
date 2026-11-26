#!/bin/sh
# OpenTela-Hetero demo: 2 workers (rtx live, t4 saved profile) and 2 heads
# (off = ETH routing, strict = hardware filter) on one PC; the same 50 jobs
# go to both heads. Usage, from the repo root: sh demo/run-demo.sh RUN
# Needs /tmp/otela from branch hetero/capability-fields. Logs: /tmp/d7/
set -u
run=${1:-}
out=demo/runs/$run
w=/tmp/d7
bin=/tmp/otela
profile=$PWD/resourceprofiler-v1/examples/fixture-kaggle-t4x2.json
[ -n "$run" ] || { echo "usage: sh demo/run-demo.sh RUN"; exit 2; }
[ -f experiment/jobs.tsv ] || { echo "STOP: run me from the repo root"; exit 1; }
[ -x "$bin" ] || { echo "STOP: $bin not found, build it first"; exit 1; }
[ -e "$out" ] && { echo "STOP: $out already exists, pick a new name"; exit 1; }
for p in 8092 8192 8292 8392 9001 9002; do
	curl -s -m 2 -o /dev/null "localhost:$p" && { echo "STOP: port $p busy, stop old nodes first"; exit 1; }
done

. demo/helpers.sh
trap stop EXIT
trap 'exit 1' INT TERM

rm -rf "$w"
mkdir -p "$w" "$out"
echo "run $run started $(date -u +%Y-%m-%dT%H:%M:%SZ)" | tee "$out/run-info.txt"
export OF_SECURITY_REQUIRE_SIGNED_BINARY=false
h1=/ip4/127.0.0.1/tcp/43905/p2p/12D3KooWHHzSeKaY8xuZVzkLbKFfvNgPPeKhFBGrMbNzbm5akpqu
h4=/ip4/127.0.0.1/tcp/43935/p2p/12D3KooWE4qDcRrueTuRYWUdQZgcy7APZqBngVeXRt4Y6ytHizKV

python3 experiment/backend.py 9001 rtx > "$w/runner-rtx.log" 2>&1 &
pids="$pids $!"
python3 experiment/backend.py 9002 t4 > "$w/runner-t4.log" 2>&1 &
pids="$pids $!"

OF_HETERO_PLACEMENT=off OF_PORT=8092 "$bin" start --mode standalone --seed 1 \
	--config-dir "$w/h1-cfg" --datadir "$w/h1-data" --tcpport 43905 --udpport 59820 \
	> "$w/head-off.log" 2>&1 &
pids="$pids $!"
wait_for "head-off (port 8092)" 60 up 8092

OF_HETERO_PLACEMENT=strict OF_PORT=8392 "$bin" start --seed 4 \
	--config-dir "$w/h4-cfg" --datadir "$w/h4-data" --tcpport 43935 --udpport 59850 \
	--bootstrap.static "$h1" > "$w/head-strict.log" 2>&1 &
pids="$pids $!"
wait_for "head-strict (port 8392)" 60 up 8392

OF_SERVICE_IDENTITY_GROUP=all OF_PORT=8192 "$bin" start --seed 2 \
	--config-dir "$w/w2-cfg" --datadir "$w/w2-data" --tcpport 43915 --udpport 59830 \
	--bootstrap.static "$h1,$h4" --service.name echo --service.port 9001 \
	> "$w/worker-rtx.log" 2>&1 &
pids="$pids $!"
OF_HETERO_PROFILE="$profile" OF_SERVICE_IDENTITY_GROUP=all OF_PORT=8292 "$bin" start --seed 3 \
	--config-dir "$w/w3-cfg" --datadir "$w/w3-data" --tcpport 43925 --udpport 59840 \
	--bootstrap.static "$h1,$h4" --service.name echo --service.port 9002 \
	> "$w/worker-t4.log" 2>&1 &
pids="$pids $!"
wait_for "head-off sees 2 echo workers" 120 workers 8092
wait_for "head-strict sees 2 echo workers" 120 workers 8392

table 8092 > "$out/table-off.txt"
table 8392 > "$out/table-strict.txt"
sh experiment/run-jobs.sh 8092 10 > "$out/off.tsv"
sh experiment/run-jobs.sh 8392 10 > "$out/strict.tsv"
sh demo/check.sh "$out/off.tsv" > "$out/off-check.txt"
sh demo/check.sh "$out/strict.tsv" > "$out/strict-check.txt"
{
	echo "otela sha256: $(sha256sum "$bin" | cut -c1-16)"
	echo "ETH code: $(git -C "${ETH:-/mnt/c/Users/user/opentela}" describe --always --dirty 2>/dev/null || echo unknown)"
	echo "workspace: $(git describe --always --dirty)"
	echo "filter log lines: strict head $(grep -c 'hetero strict' "$w/head-strict.log"), off head $(grep -c 'hetero strict' "$w/head-off.log")"
} >> "$out/run-info.txt"
echo "== off (ETH routing)"; cat "$out/off-check.txt"
echo "== strict (hardware filter)"; cat "$out/strict-check.txt"
cat "$out/run-info.txt"
