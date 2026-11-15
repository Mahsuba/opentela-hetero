# Week 4 Day 3 — Node state in OpenTela

Commit inspected: 4b6e90f (upstream eth-easl/opentela main)

## 1. Peer — src/internal/protocol/node_table.go, line 162
- What it is: One Peer = one computer (node) in the network, as seen by this computer.
- Fields similar to my Node: ID and Status match exactly. GPU/VRAM live inside the Hardware field.
- Other useful fields: LastSeen (like my heartbeat), Labels map[string]string (free-form tags), Latency, Role, Service.
- Hardware field: `Hardware common.HardwareSpec` (defined in src/internal/common).

## 2. NodeTable — same file, line 278
- What it is: `map[string]Peer` — a phonebook of all known nodes.
- Key: peer ID with a leading "/" (evidence: line 546, `table["/"+peerId]`).
- Value: the Peer.

## 3. HardwareSpec + GPUSpec — src/internal/common/structs.go, lines 19-30
- GPUSpec: Name, TotalMemory, UsedMemory.
- HardwareSpec: GPUs []GPUSpec, Memory, MemoryBandwidth, UsedMemory.
- My Node.GPU ~ GPUSpec.Name; my Node.VRAM ~ GPUSpec.TotalMemory.

## 4. Candidate gap (NOT verified yet)
HardwareSpec has no explicit field for:
- CPU architecture
- GPU vendor
- GPU compute capability
- driver / CUDA version
To check: does discovery fill these elsewhere (e.g. Labels)? Does the scheduler use them? -> Days 4-5.

## 5. Clue for Day 6
UpdateNodeTable() (node_table.go ~line 287) writes the peer to the CRDT store ("broadcast the peer to the network").

## 6. Evidence for the candidate gap (verified by me)
Command (run in opentela/src):
grep -rniE "GOARCH|ComputeCap|compute_cap|CUDAVersion|cuda_version|vendor|Architecture" --include=*.go . | grep -v _test
Result: only 1 hit -> entry/cmd/update.go:15 `arch := runtime.GOARCH`
- Used ONLY by doUpdate() to download the right binary (otela-<arch>).
- NOT stored in Peer or HardwareSpec -> other nodes never learn a node's CPU arch.
- No GPU vendor / compute capability / CUDA version field anywhere in Go code.
Caveat: Labels (map[string]string) is free-form; an operator could add arch manually. Check on Day 4 how Labels is filled.

## 7. What I learned today
- Today I opened the real OpenTela code for the first time to find where it keeps info about computers.
- In my Week 3 sim I made a `Node` struct. OpenTela's version is called `Peer`. Same idea, just bigger.
- All the Peers are kept in a `NodeTable`. It's just a Go map, like a phonebook: node ID -> Peer.
- Peer has a `Hardware` field. Inside it there's a list of GPUs with their name and memory. So OpenTela already knows GPU name and VRAM, like my sim did.
- A struct can sit inside another struct (Hardware inside Peer). And `common.` in front of a type means it lives in another folder (package).
- The `json:"..."` tags I learned in Week 2 are used here too. They let a node turn its Peer into JSON and send it to others.
- I searched the WHOLE codebase for CPU architecture, GPU vendor, compute capability and CUDA version. Only one hit, and it's just for downloading updates. So right now other nodes can't learn these things about a node.
- This might be the gap for my project. But I can't say for sure yet. I still need to check how Labels is filled (Day 4) and what the scheduler actually uses (Day 5).
- Biggest lesson: don't guess. Search the code, find the exact file and line, then write it down.
- Small lesson: `>` overwrites a file, `>>` adds to the end. Always check `pwd` before creating a file.
