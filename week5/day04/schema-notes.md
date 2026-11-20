# Week 5 Day 4 - Schema (schema.json)

## Why
Day 1, 2 and 3 each printed in a different style. The same number (used memory 966)
had two different labels. A scheduler cannot read that reliably.
So I made one fixed shape (a JSON schema) that every computer must use.

## Files
- schema.json: the rules (which fields, what type, which values are allowed)
- examples/genie-live.json: a real profile of my PC that follows the rules
- validate.py: checks a JSON file against schema.json and prints PASS or FAIL

## Fields and why each one is there
| field | reason |
|-------|--------|
| schema_version | if the shape changes later, readers know which rules apply |
| collected_at (UTC) | H3: how old the information is |
| source.mode / source.tool | honesty: live or saved file, and which tool |
| cpu.arch | H1: wrong architecture = Exec format error (Day 1) |
| gpu_status | "no GPU" and "could not ask" are different; an empty list alone cannot show that |
| gpu_error | the error text, so a person can see why |
| vendor | H1: CUDA jobs need NVIDIA |
| name | OpenTela already has it (GPUSpec.Name) |
| memory_total_mib, memory_used_mib | H2: free memory = total - used |
| compute_capability | H1: a job may need a minimum GPU generation |
| driver_version | H1 candidate; RQ4 may show it is not needed |
| uuid_short, pci_bus_id | Day 3: which physical card / which slot |

Left out on purpose: serial (my card said 0, and real serials are private),
full uuid (privacy), index (only an order number), memory_free (can be calculated),
hostname (OpenTela already has a peer ID).

Rules I chose: memory is a whole number in MiB. compute_capability is text ("12.0").
null means "not available" (Day 3 lesson: never fill in a fake value).

## Where the example values came from
I typed examples/genie-live.json by hand from real readings on 2026-11-20:
date -u gave 2026-11-20T14:35:22Z, nvidia-smi gave the GPU values,
step1.go gave amd64, Day 3 gave GPU-9877209e. Later the profiler will write this file itself.

## Tests
- Real example: PASS.
- 5 broken copies in /tmp, each must FAIL. All 5 FAILED for the right reason:
  x86_64 instead of amd64, a long uuid, memory as text "962", an extra hostname field,
  and tool_missing with a GPU still listed.
- A paste into nano got cut at line 105. I found it because json.tool failed.
  After fixing, sha256sum matched Claude's tested copy (schema.json a672033f...).

## What the schema cannot check yet
- That used memory is not bigger than total memory (will be a Day 6 test).
- That fixture_file is there when mode is fixture.
- OpenTela's GPUSpec memory units are not checked yet. Check in Week 6 before mapping
  memory_total_mib -> TotalMemory and memory_used_mib -> UsedMemory.
- jsonschema 4.19.2 on my PC. Error wording can differ in other versions, but PASS/FAIL is the same.
