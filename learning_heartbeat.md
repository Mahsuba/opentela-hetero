 Week 3 — Day 4
 Topic: Heartbeats and Failure Detection

1. What I learned

Today I learned about heartbeats in distributed systems.

A heartbeat is a recurring signal sent by a computer to indicate that it is still responding.

I implemented a small Go simulation containing three nodes:

 computer-1 (Dhaka)
 computer-2 (London)
 computer-3 (Zurich)

The first two nodes continued sending heartbeats, while the third node was intentionally stopped.

I implemented a timeout mechanism using Go's time package. If a node's last heartbeat becomes older than the configured timeout of three seconds, the monitoring system suspects that the node is unavailable.

2. What initially confused me

Initially, I found the two Boolean fields `Stopped` and `Alive` confusing.

For example:

```go
Stopped: false,
Alive:   true,
```

I initially wondered why both fields were necessary and whether they represented the same thing.

I later understood that they describe two different perspectives:

 `IntentionallyStopped` describes the configured behavior of the node in our simulation.
 `HeartbeatConfirmed` describes what the monitoring system has observed.

This helped me understand an important distinction: the actual condition of a computer and what another computer knows about its condition are not necessarily the same.

3. An important insight

A missing heartbeat does not necessarily prove that a computer has crashed.

The computer might be experiencing network delays, or its heartbeat message might be delayed.

Therefore, a timeout-based monitoring system can only suspect that a node is unavailable.

This distinction helped me understand why failure detection in distributed systems is not always straightforward.

4. What I implemented

- Defined a Node structure.
- Stored each node's last heartbeat timestamp.
- Simulated recurring heartbeats.
- Configured a three-second timeout.
- Detected a node that stopped sending heartbeats.
- Displayed the final monitoring results.

Implementation:** `heartbeat_sim.go`

5. What I want to understand further

I would like to explore how distributed systems handle delayed heartbeats, false suspicions, and nodes that recover after becoming unavailable.

These are possible directions for further experiments.



Reflection: Today's exercise helped me understand that monitoring a distributed system is not simply about knowing whether a computer is running. It is also about reasoning from incomplete and potentially delayed information.
