# 01: Introduction to Proxies

## What is a Proxy?
At its simplest, a **proxy** is an intermediary.
Imagine you want to buy something from a store, but you don't want the shopkeeper to know who you are. You send your friend. Your friend is the proxy.

### Forward Proxy vs. Reverse Proxy
- **Forward Proxy:** Sits in front of a *client* (e.g., your computer at work). It hides the client from the internet.
- **Reverse Proxy:** Sits in front of a *server* (e.g., your web app). It hides the server from the internet.

## Why use a Reverse Proxy?
1. **Load Balancing:** Distribute traffic to multiple servers.
2. **Security:** Hide the identity and structure of your backend.
3. **SSL Termination:** Handle HTTPS at the proxy so the backend doesn't have to.
4. **Security Middleware (Our Goal):** Inspect traffic for malicious patterns.

## Why Go for Proxies?
Go's `net/http` and `net/http/httputil` libraries are industry-standard. They are:
- **Fast:** Compiled and concurrent (Goroutines).
- **Simple:** The standard library provides almost everything we need to build a robust proxy.

---
*Next: Setting up our first HTTP Server in Go.*
