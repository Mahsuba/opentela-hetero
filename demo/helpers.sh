# OpenTela-Hetero demo: small helpers for demo/run-demo.sh (it loads this file).

pids="" # every process the demo starts, so stop can end them all

alive() { for p in $pids; do kill -0 "$p" 2>/dev/null && return 0; done; return 1; }

# stop: ask everything to end (like Ctrl+C), wait up to 20 s, then force.
stop() {
	trap - EXIT
	[ -n "$pids" ] || return
	echo "stopping all nodes and runners"
	kill $pids 2>/dev/null
	i=0
	while alive && [ "$i" -lt 20 ]; do sleep 1; i=$((i + 1)); done
	alive && { echo "some did not stop in 20 s, forcing"; kill -9 $pids 2>/dev/null; }
	echo "stopped"
}

# wait_for WHAT SECONDS COMMAND...: try COMMAND once a second.
wait_for() {
	what=$1 t=$2
	shift 2
	i=0
	until "$@"; do
		i=$((i + 1))
		[ "$i" -ge "$t" ] && { echo "STOP: $what not ready after $t s (logs in $w)"; exit 1; }
		sleep 1
	done
	echo "ready after ${i} s: $what"
}

# up PORT: true when a node answers on that HTTP port.
up() { curl -s -m 1 -o /dev/null "localhost:$1/v1/health"; }

# workers PORT: true when that head's table shows 2 workers offering echo.
workers() {
	curl -s -m 2 "localhost:$1/v1/dnt/table" | python3 -c '
import json, sys
try:
    t = json.load(sys.stdin)
except ValueError:
    t = {}
n = sum(1 for p in t.values() if any(s.get("name") == "echo" for s in p.get("service") or []))
sys.exit(0 if n == 2 else 1)'
}

# table PORT: one short line per node in that head's table (no owner/hostname).
table() {
	curl -s -m 2 "localhost:$1/v1/dnt/table" | python3 -c '
import json, sys
for k, p in sorted(json.load(sys.stdin).items()):
    h = p["hardware"]
    print(k[1:13], [s["name"] for s in p.get("service") or []], h.get("cpu_arch"), h.get("gpu_status"),
          [(g["name"], g["total_memory"], g["used_memory"], g.get("compute_capability")) for g in h.get("gpus") or []])'
}
