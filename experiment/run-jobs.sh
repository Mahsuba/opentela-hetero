#!/bin/sh
# OpenTela-Hetero: send the same jobs (experiment/jobs.tsv) to one head.
# Usage: sh experiment/run-jobs.sh HEAD_HTTP_PORT REPEATS > out.tsv
# One line per request: job, try, HTTP status, node that ran it, reply.
port=$1
n=$2
tab=$(printf '\t')
tmp=$(mktemp -d)
printf 'job\ttry\tstatus\tnode\treply\n'
grep -v '^#' experiment/jobs.tsv | while IFS="$tab" read -r job arch cc mem; do
	i=1
	while [ "$i" -le "$n" ]; do
		set -- -H "X-Otela-Fallback: 2"
		[ "$arch" != "-" ] && set -- "$@" -H "X-Otela-Need-Arch: $arch"
		[ "$cc" != "-" ] && set -- "$@" -H "X-Otela-Need-Min-CC: $cc"
		[ "$mem" != "-" ] && set -- "$@" -H "X-Otela-Need-Min-Free-MiB: $mem"
		rm -f "$tmp/h" "$tmp/b"
		status=$(curl -s -m 10 -D "$tmp/h" -o "$tmp/b" -w '%{http_code}' -X POST "$@" \
			--data '{}' "http://localhost:$port/v1/service/echo/test")
		node=$(grep -i '^x-computing-node:' "$tmp/h" 2>/dev/null | cut -d' ' -f2 | tr -d '\r' | cut -c1-12)
		reply=$(cat "$tmp/b" 2>/dev/null | tr -d '\n')
		printf '%s\t%s\t%s\t%s\t%s\n' "$job" "$i" "$status" "${node:--}" "$reply"
		i=$((i + 1))
	done
done
rm -rf "$tmp"
