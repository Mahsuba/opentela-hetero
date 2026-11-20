# Week 5 Day 3 - Device identity (step3.go)

## What I wanted to find out
On Day 2 I got my GPU's name and memory. But the name only tells the model.
If a computer has two RTX 5060 cards, both have the same name.
Today I wanted a way to tell which physical card is which.

## What nvidia-smi can give
- index: the GPU's number in the list (0, 1, 2...). This can change if the order changes.
- pci.bus_id (PCI address): which slot on the motherboard the card is plugged into.
  It changes only if the card is moved to another slot.
- uuid: a unique ID stored in the card itself. It stays with the card.
So the best ID is the uuid. If there is no uuid, use the PCI address. If neither, use the index.

## What my GPU said (2026-11-20)
I ran: nvidia-smi --query-gpu=index,pci.bus_id,uuid,serial,name --format=csv
- index: 0
- PCI address: 00000000:01:00.0
- uuid: GPU-9877209e... (I cut it short in these notes)
- serial: 0
- name: NVIDIA GeForce RTX 5060

## The serial problem
I expected the serial to say [N/A]. It said 0 instead.
0 is not a real serial number. It is what nvidia-smi prints when it has no serial.
If my program saved 0 as a serial, every gaming GPU would look like the same card.
So step3.go treats 0, N/A, [N/A], Not Supported and empty as "not available".

## Does OpenTela collect GPU identity?
I checked ETH's code on branch main, commit 4b6e90f.
1. src/internal/platform/gpu.go line 19: OpenTela asks nvidia-smi only for
   name, memory.total and memory.used. It does not ask for uuid or PCI address.
2. I searched the whole hardware folder (src/internal/platform) for pci, bus_id and uuid. Nothing was found.
   To make sure my search worked, I searched the same folder for memory.total. It found gpu.go line 19.
   So the search works, and the empty result is real.
3. src/internal/common/structs.go lines 19-23: GPUSpec has only Name, TotalMemory and UsedMemory.
   There is no ID field. So if a computer had two identical cards, other computers
   could not tell them apart.
Result: OpenTela does not collect GPU identity. nvidia-smi can give it, but OpenTela does not ask for it.
This is the same as Day 1 (CPU architecture) and Day 2 (compute capability).

## Running it 3 times
I ran step3.go 3 times, 10 seconds apart:
| time     | uuid         | PCI address      | used memory |
|----------|--------------|------------------|-------------|
| 05:59:57 | GPU-9877209e | 00000000:01:00.0 | 964 MiB     |
| 06:00:07 | GPU-9877209e | 00000000:01:00.0 | 952 MiB     |
| 06:00:17 | GPU-9877209e | 00000000:01:00.0 | 953 MiB     |
The uuid and PCI address did not change. Used memory changed every time.
This means the ID only needs to be sent once, but used memory must be sent again and again,
or other computers will see an old number (this connects to H3).

## What step3.go does
1. Asks nvidia-smi for index, PCI address, uuid, serial, name and used memory.
2. If a value is 0 or N/A, it prints "not available" instead.
3. Cuts the uuid short so the full ID is not printed.
4. Picks the best ID: uuid first, then PCI address, then index.
5. Prints the ID fields (should stay the same) separately from used memory (changes).

## How I tested it
- Live: go run week5/day03/profiler/step3.go
  It showed my real GPU and picked the uuid as the ID.
- From a saved file: go run week5/day03/profiler/step3.go week5/day03/fixtures/rtx5060-identity.csv
  It said FIXTURE and showed used memory 969 from the file.
- No GPU: go build -o /tmp/step3 week5/day03/profiler/step3.go && env PATH= /tmp/step3
  It printed: no NVIDIA GPU identity found (exec: "nvidia-smi": executable file not found in $PATH)
  It did not crash and did not make up numbers.

## Why this matters
In our experiments each computer will have one GPU, so the ID probably will not change which computer gets the job.
But it is still useful:
- In Week 7 many test nodes will use my real GPU data. The uuid shows they all came from the same card.
- To check if data is old (H3), I need to follow the same GPU over time.
Whether OpenTela really needs this field will be decided later (RQ4).

## What I have not tested yet
- I have only one GPU, so I did not test two identical cards.
- I only ran it for 30 seconds. I did not check if the uuid stays the same after a restart.
- I did not compare the WSL PCI address with Kubuntu (lspci), so I don't know yet if WSL shows the real slot.
- Only NVIDIA. I have no AMD graphics card, so the AMD side was not tested.
- The short uuid (8 characters) is enough for a small cluster, but two cards could match in theory.
