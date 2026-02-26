# 13: Request Tracing & Correlation IDs 🆔

When your system grows and you have many backends, debugging becomes hard. A user might report an error, but how do you find that *specific* request in your logs?

## The Solution: Correlation IDs
A Correlation ID (usually the header `X-Request-ID`) is a unique string (like a UUID) assigned to an HTTP request at the very first point of entry—**The Proxy**.

### How it works:
1.  **Generation:** API Sentinel receives a request. It checks if `X-Request-ID` already exists. If not, it generates a new, unique one.
2.  **Injection:** It adds this ID to the request headers before forwarding it to the backend.
3.  **Logging:** The ID is included in all our audit logs.
4.  **Backend Propagation:** The backend should also include this ID in its own logs.

Now, if a user has an issue, they can give you their Request ID, and you can search all your systems for that exact string.

Next, we will implement the Tracing middleware!
