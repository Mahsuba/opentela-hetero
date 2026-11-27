#!/usr/bin/env python3
# OpenTela-Hetero: turn experiments/workloads.json into the exact list of
# requests, one line each (tab-separated). No randomness: the same file
# always gives the same lines. Usage: python3 experiments/expand.py h1|h2
import json
import sys

part = sys.argv[1] if len(sys.argv) > 1 else ""
if part not in ("h1", "h2"):
    sys.exit("usage: python3 experiments/expand.py h1|h2 [workloads.json]")
path = sys.argv[2] if len(sys.argv) > 2 else "experiments/workloads.json"
w = json.load(open(path))[part]


def cell(v):
    return "-" if v is None else str(v)


if part == "h1":
    plan = [(j, 1, t) for j in w["jobs"] for t in range(1, w["tries_per_job"] + 1)]
else:
    plan = [(w["job"], r, t) for r in range(1, w["rounds"] + 1)
            for t in range(1, w["jobs_per_round"] + 1)]
print("seq\tjob\tround\ttry\tarch\tmin_cc\tmin_free_mib")
for seq, (j, r, t) in enumerate(plan, 1):
    print(seq, j["name"], r, t, cell(j["arch"]), cell(j["min_cc"]), cell(j["min_free_mib"]), sep="\t")
