FROM golang:1.20 AS builder

WORKDIR /app

# Copy entire context so Docker builds work when the build context is repository root
# (some platforms set build context to repo root while using a Dockerfile in a subdir).
COPY . .
RUN go mod download

# Build a fully static linux binary suitable for small runtime images
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o carlock ./cmd/carlock

FROM alpine:latest

# Add CA certs for HTTPS clients
RUN apk add --no-cache ca-certificates

WORKDIR /root/

COPY --from=builder /app/carlock .

EXPOSE 8080

CMD ["./carlock"]