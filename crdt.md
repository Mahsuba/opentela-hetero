# my understaning

Ana lives in Zurich, Switzerland.
Mahi lives in Dhaka, Bangladesh.
Yao lives in Beijing, China.

Each person owns a computer connected to the same decentralized computing network.

Ana's computer is Computer A.
Mahi's computer is Computer B.
Yao's computer is Computer C.

All three computers maintain their own local copies of shared network information.

These local copies are called replicas.


## Chapter 1: Ana Makes an Update

One morning, Ana discovers that her Zurich computer has an NVIDIA RTX 5090 GPU available.

She updates the shared network information to say:

"Computer A in Zurich has an RTX 5090 available."

Ana's computer now has this new information.

Yao's computer can receive the announcement about the update.

But Mahi's computer is disconnected from the network.

Therefore, Mahi does not receive the update yet.

His local copy still contains the older information.


## Chapter 2: Mahi Makes His Own Update

While disconnected from the network, Mahi checks his Dhaka computer.

He discovers that his computer has an NVIDIA RTX 5060 GPU available.

Mahi updates his local copy of the shared information:

"Computer B in Dhaka has an RTX 5060 available."

Now Ana and Mahi have made separate updates.

Ana's update describes the RTX 5090 in Zurich.

Mahi's update describes the RTX 5060 in Dhaka.

Neither computer knows about the other's latest update yet.





## Chapter 3: Mahi Reconnects

After some time, Mahi's internet connection comes back.

His computer can now communicate with the other computers.

Mahi's computer discovers that it missed Ana's update.

It retrieves the missing data from the network.

The computers exchange their updates and merge them according to the CRDT's rules.

Ana's Zurich RTX 5090 information is preserved.

Mahi's Dhaka RTX 5060 information is also preserved.

The shared information can now contain both updates.




## Chapter 4: Yao Synchronizes

Yao's computer in Beijing also needs the latest shared information.

It receives announcements about updates and retrieves the required data blocks.

It merges the updates according to the CRDT's rules.

Eventually, Ana, Mahi, and Yao can all have the same shared information.

Their computers now know about:

- Computer A in Zurich: RTX 5090
- Computer B in Dhaka: RTX 5060

This process is called convergence.

Convergence means that replicas eventually reach the same state when they have received the same updates 
and applied the same merging rules.


## Chapter 5: All Three Computers Agree

After synchronization:

- Ana's computer knows about the RTX 5090 in Zurich.
- Mahi's computer knows about the RTX 5060 in Dhaka.
- Yao's computer knows about both GPUs.

All three computers eventually have the same shared information.

This does not mean they have identical hardware.

It means their copies of the shared network information agree.

That is the main idea behind CRDT convergence.

The computers can make updates independently, exchange information later, and merge updates using shared rules.


