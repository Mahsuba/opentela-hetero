The two computers

Ana’s computer is a big, strong machine at ETH’s computer centre. Let’s say it has 4 big NVIDIA graphics cards.

Mahi’s computer is your own PC at home. Let’s pretend it has one small graphics card.

Both want to join the same OpenTela team.

Chapter 1: Ana starts a model (fleet manager + SLURM)

Ana wants to run an AI model, say “model X”, on the team.

She doesn’t pick a machine by hand. She writes a short recipe, called a preset, in the fleet manager:

“Model X needs gpus: 4.”

The fleet manager passes this to SLURM. 
SLURM is the booking system of the computer centre: it hands out machines to whoever asks.
The request goes in a line like #SBATCH --gpus-per-node=4.

SLURM looks around, finds a machine with 4 graphics cards (Ana’s computer), and starts OpenTela + model X on it.

Key point: a person (Ana) said what hardware is needed, and SLURM found a machine that fits. 
OpenTela itself didn’t check anything.




Chapter 2: Ana’s computer looks at itself (nvidia-smi, HardwareSpec)

When OpenTela starts on Ana’s computer, the code in platform/gpu.go runs a small NVIDIA tool called nvidia-smi. It asks: “What graphics cards do I have?”

(If there’s no NVIDIA card, it tries rocm-smi, the same kind of tool for AMD cards.)

It writes down only three things for each card:

name, e.g. “A100”
memory.total, how big the card is
memory.used, how full it is right now

These go into a box called HardwareSpec. That box sits inside Ana’s card about herself, which the code calls a Peer (node_table.go:162).

So Ana’s Peer says: “I’m Ana’s computer. I’m connected. I run model X. My hardware: 4 × A100, this much memory, this much used.”





Chapter 3: Ana’s computer tells everyone (CRDT, ocf-crdt, PutHook)

Now the rest of the team needs to see Ana’s card.

OpenTela uses a shared notebook that every computer keeps its own copy of.
The proper name is a CRDT, a notebook that many computers can write in without fighting over it. ETH didn’t write it themselves; they use a ready-made one called go-ds-crdt.

Ana’s computer writes its Peer card into the notebook.
The computers talk on a channel named ocf-crdt.
Every 5 seconds, Ana’s computer calls out on that channel: “My notebook has new pages!”
Mahi’s computer hears it and fetches the new page.
A piece of code called PutHook runs on Mahi’s computer. It throws away Ana’s old card and puts in the new one, whole.
Mahi’s computer keeps its own note, “can I actually reach Ana right now?” (called Connected). Ana’s card is not allowed to change that note.

the catch: Ana’s computer only writes a new card when something happens, like starting or stopping a model. 
There’s no clock saying “update every minute”. So if Ana’s cards fill up later, everyone else may still be reading an old card.








Chapter 4: Mahi’s computer joins on its own

Mahi’s computer joins the team by itself. No preset, no fleet manager, no SLURM.

It does exactly what Ana’s did: runs nvidia-smi, writes a Peer card with its small graphics card, and shares it through the CRDT notebook.

Now both computers have cards in everyone’s notebook. Let’s pretend both say “I run model X”.













Chapter 5: A job arrives (GetAllProviders, selectCandidates, random pick)

Someone sends a job: “I want model X.”

Here’s what OpenTela does, in the file proxy_handler.go:

GetAllProviders looks in the notebook: “Who is connected and runs model X?” → Ana’s and Mahi’s.
selectCandidates narrows it by group (exact group first, then "*", then "all"), plus checks like price and permission.
Then it chooses one. There’s a function called scoreCandidates meant to give each computer a score, but right now it gives everyone the same score, 1.0.
So the default chooser, lb.Pick, just picks at random.

Key point: in this whole path, OpenTela never reads the HardwareSpec box. We searched proxy_handler.go and it’s used zero times. Ana’s and Mahi’s computers each have a 50/50 chance.

(ETH did write a smarter chooser that looks at A100 vs H100 vs 3090, but only in a practice program in Python, the simulator. The real Go program doesn’t use it.)










Chapter 6: Is that a mistake? Usually no.

If every computer came in through the fleet manager + SLURM, a person already made sure each one is strong enough for its model. Then picking at random is perfectly fine.

So it’s on purpose, not a bug. We must never say “OpenTela ignores hardware”.







Chapter 7: Where it might go wrong (my guesses)
Mahi’s computer is the outsider. Nobody checked it. Maybe model X runs on it, but very slowly, or it runs out of memory on big jobs. (We still need to check exactly how this can happen, for example what the "*" and "all" groups really mean. It’s on our to-check list.)
Ana’s computer is busy. Its cards are 95% full, but random picking still sends it half the jobs. The memory.used info exists, but nobody reads it.
The card is old. Even if someone did read memory.used, it might be from an hour ago, because cards only update when something happens.
Chapter 8: What we’ll test next
Test 1: mix in outsider computers like Mahi’s. Count how often jobs land on a computer that can’t handle them, with and without a hardware check.
Test 2: make Ana’s computer busy. Compare “random” with “pick the less busy one”, and see how slow the slowest jobs get.
Test 3: make computers update their card every few seconds. Is the info fresher, and how many extra messages does it cost?
