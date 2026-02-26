# 04: The Inspection Engine (The Banger Feature) 🛡️

Now we get to the core of API Sentinel: **Application Security (AppSec)**.

## Why inspect requests?
Security headers and rate limiting are great, but they don't stop a hacker from trying to steal data using:
- **XSS (Cross-Site Scripting):** Injecting `<script>` tags into your pages.
- **SQL Injection:** Injecting `' OR 1=1 --` into your login forms to bypass authentication.

## How do we stop them?
We will build a middleware that:
1.  Reads the **Query Parameters** (e.g., `?id=1' OR 1=1 --`).
2.  Reads the **Request Body** (e.g., JSON or form data).
3.  Runs **Pattern Matching (RegEx)** against a list of known malicious signatures.
4.  If a match is found, we block the request *before* it reaches the backend.

## The Challenge with Bodies
In Go, the `r.Body` is a "Stream". Once you read it, it's gone.
Since the proxy *also* needs to read it to forward it, we must:
1.  Read the body.
2.  Inspect it.
3.  **Restore** the body so the proxy can read it again.

```go
body, _ := io.ReadAll(r.Body)
// ... inspect ...
r.Body = io.NopCloser(bytes.NewBuffer(body)) // Restore!
```

Next, we will implement the `Inspector` middleware!
