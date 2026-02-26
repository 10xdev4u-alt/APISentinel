# 12: Active Health Checks & Self-Healing 🩺

A Load Balancer is only as good as its knowledge of the backend. If a backend server dies, we shouldn't wait for a user to report it!

## Passive vs. Active Health Checks
- **Passive:** Wait for a request to fail, then mark the server as down. (A bit late!)
- **Active:** Proactively send "heartbeat" requests to the server every few seconds. If it doesn't respond, take it out of the rotation.

## Our Implementation
We will add a background **Monitor** to each `LoadBalancer`. 
Every 10 seconds, it will:
1.  Try to connect to each backend URL.
2.  Mark it as `Healthy` or `Unhealthy`.
3.  The Round Robin logic will only pick `Healthy` targets.

## Self-Healing
Once an `Unhealthy` server comes back online and starts responding to heartbeats, the monitor will mark it `Healthy` again, and it will automatically rejoin the rotation.

Next, we will implement the Health Monitor!
