Week 5 Day 1 — CPU architecture (step1.go)


## What is "architecture"?
A CPU only understands one family of instructions. That family is its architecture.
Like plugs: a UK plug does not fit a US socket. A program made for one family
cannot run on another.

Main families:
- amd64   = normal Windows/Linux PCs (Intel or AMD chips)
- arm64   = phones, Raspberry Pi, new Macs
- riscv64 = open design, used in research boards

## What I built
File: week5/day01/profiler/step1.go
It prints the architecture using Go's runtime.GOARCH.

    go run week5/day01/profiler/step1.go
    architecture: amd64

runtime.GOARCH = the family the program was BUILT for.
Usually that is the same as the computer it runs on.

## Same family, different names
    uname -m                  -> x86_64
    Go (runtime.GOARCH)       -> amd64
    file (arm64 program)      -> aarch64

| Go name | Linux name |
|---------|------------|
| amd64   | x86_64     |
| arm64   | aarch64    |

Why it matters: if one computer says "x86_64" and a job asks for "amd64",
a simple scheduler compares the words, sees they are different, and says
"no match" even though they are the same. So the profiler must use ONE name
per family.

## Experiment: wrong plug
    GOARCH=arm64 go build -o /tmp/step1-arm64 week5/day01/profiler/step1.go
    /tmp/step1-arm64
    -> cannot execute binary file: Exec format error

My PC is amd64, so it cannot run an arm64 program. It does not even start.
This is what happens if a scheduler sends a job to the wrong kind of computer
(an incompatible placement -> RQ1, H1).

## What OpenTela has (checked at baseline 4b6e90f)
    grep -rn "GOARCH" /mnt/c/Users/user/opentela --include=*.go
    -> only one hit: src/entry/cmd/update.go:15

OpenTela HAS:
- It reads its own architecture (update.go:15, arch := runtime.GOARCH).
- It uses it only for "otela update": to download the right version of itself
  from GitHub, e.g. otela-amd64 (update.go:16).

But OpenTela does NOT:
- keep it (arch is a temporary variable inside doUpdate()),
- put it in HardwareSpec or Peer (common/structs.go:19-30, Week 4 Day 3),
- share it with other computers or the scheduler.

So: OpenTela KNOWS its architecture, but does not SHARE it.
The gap is "known but not shared", not "unknown".

## Decision
Use Go names (amd64, arm64, riscv64) as the one name per family.
Reason: ETH's own release files use Go names (otela-amd64), and the
fleet-manager docs also write amd64/arm64.

## Check (workbook)
uname -m says x86_64, Go says amd64 -> same family -> output is correct.
