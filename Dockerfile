# ── Build stage ───────────────────────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /isk-server ./cmd/server

# ── Runtime stage ─────────────────────────────────────────────────────────────
FROM scratch

COPY --from=builder /isk-server /isk-server

# Railway sets $PORT dynamically; default 8080
EXPOSE 8080
ENTRYPOINT ["/isk-server"]
