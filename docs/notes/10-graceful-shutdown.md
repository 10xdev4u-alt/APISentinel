# 10: Graceful Shutdown - Being a Good System Citizen 🛑

In production, applications don't just "exit". They are shut down by the operating system (SIGTERM) or a developer (SIGINT/Ctrl+C).

If you just let the program die, you might:
- Lose the last few lines of your Audit Log.
- Drop active connections that were halfway through a download.
- Leave internal state in an inconsistent mess.

## What is Graceful Shutdown?
It means telling the server:
1. "Stop accepting *new* requests."
2. "Wait for existing requests to finish (with a timeout)."
3. "Close all files, database connections, and logs."
4. "Finally, exit."

## How it works in Go
Go's `http.Server` has a `Shutdown(ctx context.Context)` method.
We combine this with a **Signal Channel** to listen for OS signals.

```go
quit := make(chan os.Signal, 1)
signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
<-quit // Block here until we get a signal
```

Next, we will implement this in `main.go`!
