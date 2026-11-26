#!/bin/sh
# OpenTela-Hetero demo: score one results file made by experiment/run-jobs.sh.
# Usage: sh demo/check.sh results.tsv
# Answer key = which worker's REPORTED hardware fits each job:
#   rtx (seed 2, live):  amd64, cc 12.0, 8151 MiB total -> never 10000 MiB free
#   t4  (seed 3, saved): amd64, cc 7.5, 2x 15360 MiB, 0 used
# Jobs with no fitting worker (cc8_free10g, arm64) should get 503.
[ -f "${1:-}" ] || { echo "usage: sh demo/check.sh results.tsv"; exit 2; }
awk -F'\t' '
BEGIN {
	name["12D3KooWKRyz"] = "rtx"; name["12D3KooWB1b3"] = "t4"
	fits["any"] = " rtx t4 "; fits["cc8"] = " rtx "; fits["free10g"] = " t4 "
	fits["cc8_free10g"] = " "; fits["arm64"] = " "
}
NR == 1 { next }
{
	job = $1; who = name[$4]
	if (!(job in fits)) { print "unknown job: " job; other++; next }
	if (!(job in seen)) { seen[job] = 1; order[++n] = job }
	if (job != "any") needs++
	if ($3 == 200 && who != "") {
		c[job, who]++
		if (index($5, "\"" who "\"") == 0) mismatch++
		if (index(fits[job], " " who " ") == 0) wrong++
	} else if ($3 == 503) {
		c[job, "503"]++
		if (fits[job] != " ") missed++
	} else { c[job, "other"]++; other++ }
}
END {
	printf "%-12s %4s %4s %4s %6s\n", "job", "rtx", "t4", "503", "other"
	for (i = 1; i <= n; i++) {
		j = order[i]
		printf "%-12s %4d %4d %4d %6d\n", j, c[j, "rtx"], c[j, "t4"], c[j, "503"], c[j, "other"]
	}
	printf "wrong placements (sent to a worker whose hardware does not fit): %d of %d jobs with needs\n", wrong, needs
	printf "fitting jobs refused (503 although a worker fits): %d\n", missed
	printf "other replies (not 200 or 503, or unknown node): %d\n", other
	printf "two-witness mismatches (node header vs runner reply): %d\n", mismatch
}' "$1"
