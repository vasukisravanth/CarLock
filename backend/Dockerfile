FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

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