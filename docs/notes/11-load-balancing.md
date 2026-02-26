# 11: Multi-Target Load Balancing (The Pro Tier) ⚖️

What happens if your backend server crashes? In our current setup, the user gets a "502 Bad Gateway" (or connection refused). 

In a professional environment, you run **multiple** copies of your backend server. API Sentinel's job is to:
1.  **Distribute** traffic between them so no single server is overwhelmed.
2.  **Failover** to a healthy server if one goes down.

## Load Balancing Strategies
1.  **Round Robin (Simplest):** Go through the list of servers one by one (1, 2, 3, 1, 2, 3...).
2.  **Least Connections:** Send the request to the server with the fewest active requests.
3.  **Weighted Round Robin:** Some servers are more powerful; send more traffic to them.

## Our Goal: Round Robin
We will implement a **Round Robin** load balancer. Each route in our `config.yaml` can now have multiple targets!

```yaml
routes:
  - path: "/"
    targets:
      - "http://localhost:9000"
      - "http://localhost:9002"
```

Next, we will implement the Load Balancing Proxy!
