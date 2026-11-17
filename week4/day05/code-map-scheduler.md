What I wanted to find out;;;;;;;

On Day 3 and Day 4 I found that every node stores some hardware info about itself (GPU name, GPU memory, labels) and shares it with the others.

Today's question was simple: when a request comes in and OpenTela has to pick a node, does it actually look at that hardware info? Or is it just sitting there?



i tried to sketch :::


User sends: POST /v1/service/llm/...
        |
        v
1. Front door        server.go:281-284
        |
        v
2. Doorman           GlobalServiceForwardHandler   (proxy_handler.go:673)
        |            just reads the service name and passes it on
        v
3. Make the list     GetAllProviders               (node_table.go:611)
        |            keeps nodes that are online + offer that service name
        v
4. Badge check       selectCandidates              (proxy_handler.go:400)
        |            keeps nodes whose badge matches the model asked for
        v
5. More crossing off billing / access / trust / too-few-workers check
        |
        v
6. Pick one          proxy_handler.go:919-942
        |            5 ways, default is random (root.go:75)
        v
   Chosen node       -> hardware was never checked


Front door (src/internal/server/server.go:281-284) Any request to /v1/service/<name>/<rest> goes to GlobalServiceForwardHandler. The address only says which service you want. 
It doesn't say anything like "I need 16 GB of GPU memory".




Doorman (src/internal/server/proxy_handler.go:673-680) Takes the service name out of the address, 
marks it as the public/open network, and hands everything to globalServiceForwardWithScope (line 691). It doesn't pick anything itself.




Make the list (src/internal/protocol/node_table.go:611-631) Goes through the node table one by one. Keeps a node only if:
it's online (peer.Connected), and
it offers a service with the right name (service.Name == serviceName). That's it. No GPU, no memory, no labels.







Badge check (src/internal/server/proxy_handler.go:400-462, called at line 743) Each node's service has a badge (identity group). Three kinds:
exact, like model=Qwen -> best match
wildcard model=* -> "any model"
all -> "send me anything" The request's model name is compared to the badge. By default only exact matches count (the X-Otela-Fallback header can loosen this). It's all name matching, not hardware.





More crossing off (proxy_handler.go, after line 743)
billing: only nodes the user can afford (X-Otela-Allowed-Peers)
access rules (control plane), if turned on
trusted region: only nodes with a direct connection
trust level: drop nodes below the trust number asked for (X-Otela-Trust)
if there are too few workers, it may reject the request (admission control) None of these look at hardware.







Pick one (proxy_handler.go:919-942) Checked top to bottom, first one that fits wins:
same conversation as before -> reuse the same node (session affinity)
price mode, first try -> cheaper nodes get picked more often
price mode, retry -> next cheapest node
weighted mode on -> scoreCandidates + weighted random
otherwise -> load balancing policy: random, round-robin or shortest-queue Default policy is random (src/entry/cmd/root.go:75). If the request fails, it retries on a different node (default up to 3 times).









Does the picker use Hardware or Labels?

No.

Proof:

grep -n "Hardware" internal/server/proxy_handler.go -> nothing found. This is the file where all the choosing happens.
grep -rn "\.Hardware\|\.Labels" internal/server -> only 2 places, and neither is about choosing:
analytics_hook.go:77-79: counts GPUs and saves the first GPU name, for stats/reports only
registration.go:208: wipes the hardware box empty when a relay registers
Labels isn't used anywhere in the server folder.

So the hardware info gets collected (Day 4) and shared (Day 3), but the code that picks a node for a request never reads it.





The ETH simulator-----im not claiming ,,,,,i mean overclaiming you can say ;)


ETH does have a mixed-hardware scheduler: HeterogeneousRiskAwareScheduler in contrib/serving-simulator/simulator/core/scheduler.py:402. It groups GPUs into tiers (H100, A100, RTX 3090) and does "hard constraint filtering (OOM prevention)", which basically means not sending a job to a GPU that's too small for it. That's very close to our idea.

But:

it's written in Python, and the real OpenTela is Go
it's a simulator (a pretend version for testing ideas), not the live system
I searched every Go file for serving-simulator, and nothing mentions it. The live system never calls it.

So what I can honestly say: ETH has hardware-aware placement in their simulator, but the live Go request path does not use the hardware info nodes already share.

I should NOT say "OpenTela ignores heterogeneity". Their paper and simulator clearly think about it.
















2 more thing: scoreCandidates (proxy_handler.go:502-512) gives every node the same score, 1.0. The comment says it's the "extension point for richer scoring (latency, load, trust, etc.) in the future".

Two things I noticed:

it only gets a list of node IDs ([]string), not the node info, so right now it can't even see hardware
their list of future ideas doesn't mention hardware

This looks like the natural place to plug in a hardware-aware check or score later. It's only used when scalability.weighted_routing is turned on (off by default).









one more thing i noticed , The scheduling docs (docs/content/docs/advanced/scheduling.mdx) describe the picking steps, but they don't mention session affinity or price-aware picking, which are in the code. 
So the docs are a bit behind the code.
lets Trust the code.

