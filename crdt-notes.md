# CRDT Notes

1. A CRDT is a data structure that helps multiple computers combine their updates safely without needing
 a central coordinator.
2. Each computer keeps its own local copy of the shared information, called a replica.
3. When a computer updates shared information, it announces the change through a messaging system called PubSub.
4. Other computers retrieve the required data blocks using a system called Bitswap.
5. Because computers can update information independently, their local copies may temporarily become different.
6. CRDTs provide rules for combining updates so that replicas receiving the same updates eventually reach 
the same state.
7. If a computer cannot retrieve a data block before the configured timeout, 
OpenTela can try synchronization again later.
8. This helps OpenTela maintain shared service discovery, health, and routing information 
even when network connections are unreliable.
