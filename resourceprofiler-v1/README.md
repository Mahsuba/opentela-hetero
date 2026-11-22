# resourceprofiler v1

A small command-line tool that prints what GPU hardware a computer has,
as one JSON "capability profile" (schema v1). Built for OpenTela-Hetero:
OpenTela (at commit 4b6e90f) stores only GPU name, total memory and used
memory. This tool also reads CPU architecture, compute capability, driver
version, a short GPU id and the PCI address, and says clearly when it
could not read them.

## Build and test (Ubuntu / WSL, Go 1.24 or newer)

    cd resourceprofiler-v1
    go test ./...
    go build -o /tmp/resourceprofiler .

The binary goes to /tmp so it is never committed.

## Run

    /tmp/resourceprofiler                              # live: this machine, now
    /tmp/resourceprofiler testdata/kaggle-t4x2.csv     # fixture: a saved capture

Live mode runs nvidia-smi (time limit 10 s). Fixture mode reads a file
saved earlier, so tests and experiments can use real readings from
machines we do not have.

## gpu_status (what the answer means)

- ok: at least one GPU was read.
- none_found: the tool ran fine and saw no GPU.
- tool_missing: nvidia-smi and rocm-smi are not installed. Unknown, not "no GPU".
- query_failed: the tool exists but gave no usable answer (error, frozen,
  bad or impossible numbers). gpu_error says why. gpus is then empty.

A field the GPU cannot report (for example "[N/A]") becomes null.
Exit code: 0 for every gpu_status. 2 = wrong arguments or a bad fixture file.

## Fixture files

A fixture is the raw nvidia-smi lines with two header lines on top:

    # collected_at=2026-11-20T20:37:34Z
    # arch=x86_64

The profile uses this time and CPU, not those of the machine replaying it.
A fixture without them is refused. Make one on another machine with:

    sh capture.sh > my-machine.csv

capture.sh is read-only and cuts full GPU uuids to 8 hex characters.
If nvidia-smi is missing or fails, it saves nothing and exits 1.

## Schema and examples

- schema/schema.json: JSON Schema v1 (same file as week5/day04).
- Check outputs: `python3 schema/validate.py schema/schema.json examples/*.json`
- examples/: real outputs (live RTX 5060, Kaggle 2x T4 fixture, and a
  run with no GPU tool on the search path).

## Known limits

- NVIDIA only. AMD (rocm-smi) is detected but not read: no AMD GPU to test on.
- A GPU name containing a comma is rejected (query_failed).
- Drivers too old to report compute_cap make the whole query fail.
- none_found text ("No devices were found") is from memory, not seen on real hardware.
- A fixture can only store a working nvidia-smi reading.
- Hardware is self-reported; nothing checks it.
