FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o carlock ./cmd/carlock

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/carlock .

CMD ["./carlock"] 

EXPOSE 8080