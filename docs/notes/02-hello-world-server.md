# 02: Building your first Go HTTP Server

To build a reverse proxy, we first need to understand how to handle HTTP requests in Go.

## The `net/http` Package
Go's `net/http` package is the backbone of web development in Golang. 

### Basic Server Structure
```go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    // 1. Define a handler
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, "Hello from API Sentinel!")
    })

    // 2. Start the server
    fmt.Println("Server starting on :8080...")
    http.ListenAndServe(":8080", nil)
}
```

### What is a Handler?
A handler is just a function that takes two things:
1. `ResponseWriter`: Where you write your response to the client.
2. `Request`: The incoming data from the client (Headers, Body, URL).

## The `httputil.ReverseProxy`
Instead of writing a custom response, a reverse proxy *forwards* the `Request` to another server and *copies* that server's response back to the original client.

Go makes this easy:
```go
target, _ := url.Parse("http://localhost:9000")
proxy := httputil.NewSingleHostReverseProxy(target)
http.ListenAndServe(":8080", proxy)
```

In the next step, we will implement this!
