# Stage 1: Build
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum (if it exists)
COPY go.mod ./
# RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o apisentinel ./cmd/apisentinel/main.go

# Stage 2: Final Image
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/apisentinel .

# Expose the proxy port
EXPOSE 8080

# Command to run the proxy
CMD ["./apisentinel"]
