
# Stale Information in a Distributed System

## 1. What is stale information?

Stale information is information that was correct earlier but is no longer up to date.

In a distributed system, computers communicate over a network. Updates may take time to reach other computers.

Because of this, a scheduler may temporarily use an older version of a node's information.

## 2. Our Three Computers

Imagine three computers in a distributed compute system:

| Computer | Location | GPU | VRAM |
|---|---|---|---|
| Computer 1 | Dhaka | RTX 5060 | 8 GB |
| Computer 2 | London | RTX 5070 | 12 GB |
| Computer 3 | Zurich | RTX 5090 | 32 GB |

The scheduler collects hardware information from these computers to decide where workloads should run.

## 3. The Staleness Story

Let's focus on Computer 2 in London.

At 10:00 AM, Computer 2 has an RTX 5070 with 12 GB VRAM.

At 10:05 AM, its GPU is replaced with an RTX 5090 with 32 GB VRAM.

Computer 2 detects the hardware change and sends an updated profile to the scheduler.

However, the network delays the update.

At 10:07 AM, the scheduler still believes Computer 2 has an RTX 5070.

The scheduler is therefore using stale information.

## 4. Timeline Diagram

```text
TIME             COMPUTER 2                 SCHEDULER
----------------------------------------------------------------

10:00 AM         RTX 5070                   RTX 5070
                 12 GB VRAM                 12 GB VRAM
                 Initial profile            Initial profile

10:05 AM         GPU replaced
                 RTX 5090
                 32 GB VRAM

10:06 AM         Detects hardware change
                 Sends updated profile  ----->

10:07 AM         RTX 5090                   RTX 5070
                 32 GB VRAM                 12 GB VRAM
                                             OLD INFORMATION
                                             |
                                             |
                                      Network delay
                                             |
10:09 AM         RTX 5090                   RTX 5090
                 32 GB VRAM                 32 GB VRAM
                                             UPDATED PROFILE
```

## 5. Why Does Information Become Stale?

There are several possible reasons:

1. Network delay: Updates take time to travel between computers.
2. Disconnection: A node may temporarily lose its network connection.
3. Update intervals: Nodes may report their hardware information periodically.
4. Processing delay: The receiving system may take time to process an update.

## 6. Why Does This Matter for Scheduling?

Lets say an AI workload needs 24 GB of VRAM.

At 10:07 AM, the scheduler believes Computer 2 has only 12 GB VRAM.

It concludes that Computer 2 cannot handle the workload.

However, the actual GPU has 32 GB VRAM.

The scheduler might therefore miss a suitable computer because it is using outdated information.

The opposite can also happen.

If the scheduler believes a node has enough resources but the actual available resources have decreased, it might assign a workload that cannot run successfully.

This can affect placement correctness and useful performance.

## 7. Connection to OpenTela-Hetero

OpenTela coordinates computing resources across multiple computers.

For hardware-aware placement, the scheduler needs sufficiently accurate and up-to-date hardware information.

Our research interest includes understanding:

- How frequently nodes should update hardware metadata.
- How long metadata takes to reach the scheduler.
- How network delays affect metadata freshness.
- How stale information affects placement decisions.
- What communication and processing overhead is introduced by frequent updates.

This creates an important trade-off:

More frequent updates may improve information freshness, but they can also increase communication overhead.

## 8. What I Learned Today

1. Information can become outdated even if it was originally correct.
2. Network delays and disconnections can cause stale information.
3. A scheduler may make decisions using an older hardware profile.
4. Stale metadata can cause unsuitable placement or missed placement opportunities.
5. Distributed schedulers must balance information freshness and communication overhead.

## 9. Age of Information

A useful concept in distributed systems is Age of Information (AoI).

It measures how old the latest information available to a receiver is.

In our example:

- The scheduler last received a profile generated at 10:00 AM.
- At 10:07 AM, the scheduler still has that old profile.
- Therefore, the information available to the scheduler is 7 minutes old.

At 10:06 AM, Computer 2 generated a new hardware profile.

However, the update was delayed by the network and arrived at 10:09 AM.

At 10:09 AM:

- The latest profile was generated at 10:06 AM.
- The profile arrived at 10:09 AM.
- Its age upon arrival was 3 minutes.

This distinction matters:

- Time since hardware change at 10:07: 2 minutes.
- Age of the scheduler's latest received profile at 10:07: 7 minutes.
- Age of the newly received profile at 10:09: 3 minutes.

The two measurements are not necessarily the same.

In our research, we can investigate how metadata freshness affects hardware-aware placement.


**Main lesson:**

A scheduler does not always know the current state of every computer. It makes decisions using the information it has received, which may be slightly out of date.
