FROM golang:1.22-alpine3.19 AS builder

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY *.go ./

# Build Go
RUN go build -o /cloudflare-ddns

# Create a smaller image for running the binary
FROM alpine:3.20.3

WORKDIR /app

COPY --from=builder /cloudflare-ddns /cloudflare-ddns

CMD ["/cloudflare-ddns"]
