

## What is a Distributed System?

### 1. Definition

A distributed system consists of multiple separate computers
that communicate over a network.

### 2. Our Three Nodes

Node A:
- Location: Dhaka
- GPU: RTX 5060

Node B:
- Location: London
- GPU: RTX 5070

Node C:
- Location: Zurich
- GPU: RTX 5090

### 3. Important Concepts

- Node: An individual computer in the system.
- Network: A way for computers to communicate.
- Message: Information sent from one computer to another.
- State: Information a computer currently stores or knows.

### 4. Single Computer vs Distributed System

Single computer:
One computer executes the work.

Distributed system:
Multiple computers communicate and coordinate their work.

### 5. Connection to OpenTela

OpenTela coordinates computing resources.
Different computers may have different hardware capabilities.

A placement system needs information about those resources
to make suitable decisions.

### 6. What I Learned

A distributed system is not simply one powerful computer.
It involves multiple independent computers communicating
through a network.

### 7. Distributed System Diagram

```text
Node A (RTX 5060) -----> Node B (RTX 5070)
       |                       |
       |                       |
       v                       v
Node C (RTX 5090) <-------------+



### 8. Network Failure

A computer may still be running even when its network connection
is unavailable. Not receiving a reply does not always mean
the computer has crashed.
