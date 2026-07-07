# ---------- Build Stage ----------
FROM golang:1.25.10 AS builder

WORKDIR /app

# Download dependencies first (better layer caching)
COPY go.mod ./
COPY go.sum* ./
RUN go mod download

# Copy source
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o duke ./cmd/duke

# ---------- Runtime Stage ----------
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/duke .

EXPOSE 8000
EXPOSE 9000

ENTRYPOINT ["./duke"]
