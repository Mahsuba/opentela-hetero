# Code map: resource discovery (OpenTela @ 4b6e90f)

Date: 2026-11-16 (Week 4, Day 4)
Repo: /mnt/c/Users/user/opentela, commit 4b6e90f (same as upstream/main)
Question for today: how does a node find out what hardware it has, and when does it check again?

lemme picture it:
OpenTela measures GPUs by running a tool (nvidia-smi or rocm-smi), gets only
3 facts per GPU, re-checks only when certain events happen, and never fills
host RAM or CPU info into the node record.

---

## 1. GPU discovery (how)

Where: src/internal/platform/gpu.go

- gpu.go:10-16 GetGPUInfo() tries NVIDIA first. If it finds at least one
  NVIDIA GPU it stops there. Only if it finds none does it try AMD.
  (So on a machine with both, only the NVIDIA GPUs would be reported.)
- gpu.go:19 It runs the program nvidia-smi and asks for only 3 things:
  name, memory.total, memory.used.
- gpu.go:21-23 If nvidia-smi is missing or fails, it returns an empty list.
  The error is logged only at Debug level, so normally you see nothing.
- gpu.go:49 For AMD it runs rocm-smi instead.
- gpu.go:84-91 AMD memory comes in bytes and is converted to MiB so it
  matches the NVIDIA units.
- gpu.go:100-112 Special case: some AMD GPUs (MCM, e.g. MI250X) show up as
  2 pieces. The code joins each pair back into 1 GPU.

What we get per GPU (structs.go:19-23, GPUSpec):
  Name (text), TotalMemory (MiB), UsedMemory (MiB)

What we do NOT get (nothing in gpu.go asks for it):
  vendor as its own field, compute capability, driver version,
  CUDA/ROCm version, GPU architecture, interconnect (NVLink etc.).
  The vendor and model are only hidden inside the Name text,
  e.g. "NVIDIA A100-SXM4-40GB".

How it gets into the node record (Peer):
- node_table.go:725  myself.Hardware.GPUs = platform.GetGPUInfo()   (at startup)
- registrar.go:253 + :266  measured again and stored, inside ReannounceLocalServices()

Who reads Hardware.GPUs directly (grep "Hardware\." in src, tests removed):
- analytics_hook.go:77-79 only counts GPUs and takes the first GPU's name,
  for statistics. That is the only direct read the grep found.
  (Does the scheduler read it some other way? -> Day 5.)

## 2. Labels (how)

- entry/cmd/root.go:71  The start command has a --label flag.
  Format key=value, can be repeated. Example: --label gpu=a100
- node_table.go:251-270  parseLabels() splits each "key=value" into a map.
  Bad entries without "=" are skipped (node_table.go:258-262).
- node_table.go:654  Labels: parseLabels()  puts the map into the Peer at startup.

In plain words: labels are DECLARED, typed by a human. Nothing checks they
are true. If I write --label gpu=h100 on a laptop, OpenTela believes it.
The only place I found that sets Peer.Labels is node_table.go:654 (startup).

## 3. Host memory / CPU (what is missing)

Host memory:
- structs.go:25-30 HardwareSpec has Memory, MemoryBandwidth, UsedMemory.
- But grep "Hardware\." shows the only writes are to Hardware.GPUs
  (node_table.go:725, registrar.go:266). Nothing ever writes Memory,
  MemoryBandwidth or UsedMemory. So these fields are always empty (0).

CPU:
- No CPU architecture, CPU model or core count is put into Peer.
- entry/cmd/update.go:15  runtime.GOARCH is used only to download the right
  file when OpenTela updates itself. It is not shared with other nodes.
- ingest/server.go:70-87  collectStats() counts CPU cores (runtime.NumCPU, :76),
  but it fills a different record called SystemStats, not Peer.
  Careful: runtime.MemStats (:71-72) is memory used by the OpenTela program
  itself, NOT the machine's RAM.
  ingest/server.go:85 GPU list here is empty with
  "TODO: Implement NVIDIA/AMD GPU collection logic".

Extra (seen while reading, not part of the main path):
- server/registration.go:208  when a relay registers a remote peer,
  Hardware is cleared to empty on purpose (with Service and Load at :206-207).

## 4. Freshness (when is it measured again)

GPU info is measured:
1. once at startup (node_table.go:725)
2. again every time ReannounceLocalServices() runs (registrar.go:247, :253)

ReannounceLocalServices() is called from 4 places, all of them EVENTS:
- host.go:328  when this node connects to another node,
  but only if this node already has a service (host.go:327 hasLocalServices())
- host.go:807-809  when the node picks a NEW relay
- service_health.go:184  when a service fails its health checks and is withdrawn
- service_health.go:195  when a service becomes healthy again
- service_health.go:196  if health did not change (Noop), nothing is re-announced

In plain words: there is no timer among these callers. If nothing happens,
the GPU info other nodes see stays old. UsedMemory changes all the time,
so it is the field most likely to be stale. (This links to RQ3 and my
Week 3 stale-information simulation, commit 5e40c54.)

Labels are never re-read after startup (only set at node_table.go:654).

## 5. questions in my mind 

For Day 5 (scheduler):
- Does the scheduler use Hardware.GPUs at all? Name? TotalMemory? UsedMemory?
- Does it use Labels for placement?
- Does it try to read GPU type out of the Name text?
- What happens to a job asking for "NVIDIA, 16 GB" today?

For Day 6 (propagation):
- After store.Put of the Peer (node_table.go, registrar.go), how fast do
  other nodes see the new value?

Possible hypotheses (to test later, not claims):
- H1: Because GPU identity is only a free-text Name, matching a job to the
  right GPU type is fragile or not done.
- H2: Because re-measuring only happens on events, UsedMemory can be stale
  for a long time on a quiet node.
- H3: Adding a few fields (vendor, compute capability, CPU arch) would cost
  very little extra metadata, because the Peer record is already sent anyway.

---

## question myself and answers

Q: How does an OpenTela node know it has a GPU?
A: It runs nvidia-smi (gpu.go:19). If no NVIDIA GPU is found, it runs
   rocm-smi (gpu.go:15, :49). It keeps name, total memory and used memory.

Q: Is that information trusted?
A: GPUs are measured by a tool, so yes, but only 3 fields. Labels are typed
   by a human (root.go:71) and never checked.

Q: How fresh is it?
A: Startup plus 4 event types (host.go:328, :809, service_health.go:184, :195).
   I found no timer, so on a quiet node it can be old.

Q: Does OpenTela know host RAM or CPU type?
A: The fields for RAM exist (structs.go:27-29) but nothing fills them.
   CPU type is not in Peer at all.

Q: So is there a gap?
A: There is a gap in what is DISCOVERED. Whether it matters depends on what
   the scheduler uses. I check that on Day 5 before claiming anything.
