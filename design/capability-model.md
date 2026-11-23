# Capability model (Week 6 Day 1)

Baseline: OpenTela 4b6e90f. Profiler: resourceprofiler-v1 (schema v1).

## Rule
A field goes into OpenTela only if it changes a placement decision in one of my hypotheses.

## What OpenTela has today
GPUSpec (common/structs.go:19-23) has Name, TotalMemory and UsedMemory. Memory is in MiB for NVIDIA (gpu.go:19, no conversion) and for AMD (gpu.go:86 and :90 turn bytes into MiB). HardwareSpec (:25-30) has the GPU list and three host memory fields that are never filled (always 0, Week 4 Day 4). There is no CPU architecture, no GPU status, no reading time and no compute capability. A request cannot say what hardware it needs: none of the X- headers in server/ is about hardware.

## The three decisions
- D1 (H1): Can this node run the job? Yes or no.
- D2 (H2): Of the nodes that can, which is best? The one with the most free GPU memory.
- D3 (H3, bonus): Is the reading too old to trust?

## Node fields
| Field | Where | Decision | Real values I have |
|---|---|---|---|
| cpu_arch (new) | HardwareSpec | D1 | amd64 (my PC, Kaggle) |
| gpu_status (new) | HardwareSpec | D1 | ok, tool_missing |
| compute_capability (new) | GPUSpec | D1 | 12.0 (RTX 5060), 7.5 (T4) |
| number of GPUs (exists) | len(GPUs) | D1 | 1 (my PC), 2 (Kaggle) |
| TotalMemory, UsedMemory (exist, MiB) | GPUSpec | D1, D2 | 8151 / ~719-892 (RTX 5060), 15360 / 0 (T4) |
| collected_at (new) | HardwareSpec | D3 | from the profiler |

All new fields are text. Empty text means "unknown". Old nodes do not send them, so they read as unknown and nothing breaks: nodes read each other with plain json.Unmarshal (crdt.go:165), which skips fields it does not know.

## Job fields (new, request side)
- arch, for example arm64 (D1)
- min_cc, for example 8.0 (D1)
- min_free_mib, for example 6000 (D1 and D2)

How a request carries them is decided on Day 5. There is a pattern to copy: X-Otela-Allowed-Peers (proxy_handler.go:750-764) already cuts the candidate list and returns 503 if no node is left.

## Rules
- A request with no requirements behaves exactly like the baseline, so the baseline stays reproducible.
- If a job asks for something and the node's value is unknown, the node is not compatible (strict). Lenient (unknown = allowed) is an experiment setting.
- gpu_status none_found means "no GPU" (fine for a CPU job). tool_missing and query_failed mean "unknown", not "no GPU".
- Compute capability is compared as numbers. As text, "12.0" comes before "8.0", which is wrong.
- free = TotalMemory - UsedMemory. A node with several GPUs counts its best GPU.

## Left out, and why
- vendor: I have no AMD machine to test it. A job with min_cc already refuses AMD GPUs, because OpenTela's AMD path gives no compute capability (unknown). Add it back if the RX 6600 capture happens.
- driver_version: no decision in H1-H3 uses it.
- uuid_short, pci_bus_id: one GPU per node in my experiments, so they cannot change a choice. They stay in the profiler JSON for labelling.
- host memory fields: already exist, always 0, not needed by my hypotheses. Left untouched.

## Cost (RQ2)
With example values, the new fields add 103 bytes of JSON per node with one GPU (75 node + 28 per GPU, counted with wc -c). Real network cost is measured in Week 7.

## Open questions
- Real arm64 node for D1: so far only a hand-written test fixture. I need a real capture or clearly labelled emulated nodes.
- D2 plugs into scoreCandidates (proxy_handler.go:502-512), which today gives every node 1.0 and gets only node IDs.
