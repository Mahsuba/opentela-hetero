
# OpenTela Glossary

## 1. What is OpenTela?

OpenTela (also called OpenFabric) is a decentralized computing platform.

It connects computing resources across multiple computers and coordinates them as a distributed network.

It uses peer-to-peer networking, gossip protocols, and CRDT-based state management.

Its purpose includes service discovery, health monitoring, request routing, and fault tolerance.

OpenTela is used for distributed LLM serving, including SwissAI Serving.

## 2. Distributed Systems

| Term | Meaning |
|---|---|
| Distributed System | A system where multiple computers work together. |
| Decentralized | No single central coordinator controls the entire system. |
| Computing Fabric | An interconnected network of computing resources. |
| Orchestration | Coordinating computing resources and services. |
| Peer | A participant in a peer-to-peer network. |
| Peer ID | A unique identifier assigned to a peer. |
| Peer Addressing | The network address used to communicate with a peer. |
| Gossip Protocol | A mechanism for spreading information between peers. |
| Network Partition | A situation where parts of a network cannot communicate. |

## 3. State and Failure Management

| Term | Meaning |
|---|---|
| CRDT | A data structure designed to help distributed replicas synchronize and converge. |
| Registry | A record of information about services and resources. |
| Heartbeat | A periodic signal indicating that a node is alive. |
| Fault Tolerance | The ability to continue operating despite failures. |
| Service Discovery | Finding available services in a network. |
| Replica | Another copy of distributed data or a service. |
| Stale Information | Information that is outdated compared with the current state. |

## 4. Identity and Routing

| Term | Meaning |
|---|---|
| Identity Group | A group of peers that can be treated as interchangeable for a particular service. |
| Permissionless Partition | The public OpenTela routing partition. |
| Trusted Region | A managed group of peers eligible for trusted routing. |
| Service Exposure | Determines the partition assignment of a service. |
| Exact Service Name | The precise, case-sensitive identifier of a service. |
| Request Routing | Directing requests to appropriate services or peers. |

## 5. Security and Authorization

| Term | Meaning |
|---|---|
| Policy Scope | Determines whether authorization applies to a peer or an individual service. |
| Service ACL | Rules defining who can access a service. |
| Authorization | Determining whether someone is permitted to access a resource. |
| User Space | The environment where ordinary applications run without directly modifying the kernel. |

## 6. Hardware and Computing

| Term | Meaning |
|---|---|
| GPU | A processor designed for parallel computing. |
| VRAM | Memory available on a GPU. |
| HPC | High-Performance Computing. |
| Slurm | A workload and job scheduler used in computing clusters. |
| vLLM | An engine for serving large language models. |
| SGLang | An engine for serving language models. |
| LLM | Large Language Model. |
| Preemption | When a running job or resource allocation is interrupted. |

## 7. My Week 3 Connections

- Gossip: I built a small gossip simulation.
- Heartbeat: I simulated a computer disappearing.
- CRDT: I learned how distributed replicas synchronize.
- Stale Information: I studied outdated GPU capability information.
- Distributed Fabric: I built a three-node simulation.

These concepts connect my earlier exercises to OpenTela's architecture.

## 8. Research Questions for OpenTela-Hetero

1. What hardware capabilities does OpenTela currently discover?
2. Where is hardware capability information stored?
3. How is capability information shared between peers?
4. How does OpenTela use hardware information when routing requests?
5. What happens when hardware information becomes stale?
6. What is the cost of discovering and propagating additional capability metadata?

## 9. Sources

- Official repository: https://github.com/eth-easl/OpenTela
- README.md
- Official glossary: docs/content/docs/tutorial/glossary.mdx
