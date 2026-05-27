FROM golang:1.20 AS builder

WORKDIR /app

# Copy entire context so Docker builds work when the build context is repository root
# (some platforms set build context to repo root while using a Dockerfile in a subdir).
COPY . .

# Detect where go.mod is (repo root vs backend subdir), install modules
# and build a single static binary at /app/carlock so this Dockerfile works
# whether the build context is the repo root or the `backend/` folder.
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN set -eux; \
	# Use `go env GOMOD` to reliably locate the go.mod file inside the copied context
	MODFILE=$(go env GOMOD); \
	if [ -n "$MODFILE" ] && [ -f "$MODFILE" ]; then \
		MODDIR=$(dirname "$MODFILE"); \
	else \
		# check common locations relative to the build context
		for d in . backend car-lock-system/backend; do \
			if [ -f "$d/go.mod" ]; then MODDIR=$d; break; fi; \
		done; \
	fi; \
	if [ -z "$MODDIR" ]; then echo "go.mod not found in context"; exit 1; fi; \
	cd "$MODDIR"; \
	echo "Using module dir: $MODDIR"; ls -la; \
	go mod download; \
	# determine build path for cmd/carlock relative to module dir
	if [ -d ./cmd/carlock ]; then \
		BUILD_TARGET=./cmd/carlock; \
	elif [ -d ./car-lock-system/backend/cmd/carlock ]; then \
		BUILD_TARGET=./car-lock-system/backend/cmd/carlock; \
	elif [ -d ../backend/cmd/carlock ]; then \
		BUILD_TARGET=../backend/cmd/carlock; \
	else \
		echo "cmd/carlock not found relative to module dir ($MODDIR)"; exit 1; \
	fi; \
	go build -ldflags='-s -w' -o /app/carlock "$BUILD_TARGET"

FROM alpine:latest

# Add CA certs for HTTPS clients
RUN apk add --no-cache ca-certificates

WORKDIR /root/

COPY --from=builder /app/carlock .

EXPOSE 8080

CMD ["./carlock"]