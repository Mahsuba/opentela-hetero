# Week 5 Day 3 fixtures

## rtx5060-identity.csv
This is a real capture from my PC (Genie, Ubuntu WSL) on 2026-11-20. Nothing is made up.

I made it with this command:
nvidia-smi --query-gpu=index,pci.bus_id,uuid,serial,name,memory.used --format=csv,noheader,nounits

The columns are: index, PCI address, uuid, serial, name, used memory (MiB).

I made one change: I cut the uuid short to "GPU-" plus the first 8 characters,
so my full GPU ID is not on GitHub. Everything else is exactly what nvidia-smi printed.

The serial says 0 because that is what nvidia-smi printed. It is not a real serial number.

Used memory (969 MiB) is the value at the moment I saved it. On the live GPU it keeps changing.
