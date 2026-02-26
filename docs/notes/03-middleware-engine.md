# 03: The Power of Middleware

In the "Beginner" phase, we built a simple proxy that just passes traffic. In the "Intermediate" phase, we make it smart.

## What is Middleware?
Middleware is like a checkpoint on a road. Every request must pass through it. It can:
- **Inspect** the request (Log it, check for auth).
- **Modify** the request (Add headers).
- **Reject** the request (Rate limiting, blocking bad actors).
- **Pass** it to the next checkpoint.

## The Middleware Pattern in Go
In Go, a middleware is usually a function that takes an `http.Handler` and returns an `http.Handler`.

```go
func SecurityHeaders(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 1. Do something before the request
        w.Header().Set("X-Content-Type-Options", "nosniff")
        
        // 2. Pass to the next handler
        next.ServeHTTP(w, r)
        
        // 3. (Optional) Do something after the response
    })
}
```

## Security Headers
These are simple but powerful. They tell the browser how to behave securely.
- `X-Content-Type-Options: nosniff`: Prevents MIME-type sniffing.
- `X-Frame-Options: DENY`: Prevents Clickjacking.
- `Content-Security-Policy`: Tells the browser which sources are trusted (prevents XSS).

Next, we will build a middleware engine and add these headers!
