# 07: Data Loss Prevention (DLP) 🛡️💨

A common security failure is when a backend server accidentally leaks sensitive data in a response. This could be:
- **PII (Personally Identifiable Information):** Social Security Numbers, phone numbers.
- **Financial Data:** Credit Card numbers.
- **Secrets:** Internal API keys, database connection strings left in debug mode.

## How DLP Works in a Proxy
API Sentinel sits between the backend and the user. While our `SecurityInspector` looks at *incoming* requests, our `DLPMiddleware` looks at *outgoing* responses.

### The Challenge: Capturing the Response
In Go, once you call `next.ServeHTTP(w, r)`, the response is sent immediately to the client. To inspect it, we must use a **Response Wrapper**. We provide a custom `ResponseWriter` to the next handler that captures the data in a buffer instead of sending it directly to the network.

### Pattern Matching
Just like the inspector, we use RegEx to look for sensitive patterns.
- **Credit Cards (Luhn-like):** `\b(?:\d[ -]*?){13,16}\b`
- **SSN:** `\b\d{3}-\d{2}-\d{4}\b`

## Blocking vs. Masking
- **Blocking:** Return a `500 Internal Server Error` and log the leak. (Safest)
- **Masking:** Replace the sensitive parts with `****` and let the response pass. (Better UX)

We will start with **Blocking** to ensure maximum security!
