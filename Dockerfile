# Build Stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies for CGO (required for SQLite)
RUN apk add --no-cache gcc musl-dev

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the API binary
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o email-api ./cmd/api/main.go

# Final Stage
FROM alpine:latest

# Install sqlite for database management
RUN apk add --no-cache sqlite-libs ca-certificates

WORKDIR /root/

# Copy the binary from the builder
COPY --from=builder /app/email-api .
# Copy default env if needed (or use env variables)
COPY --from=builder /app/README.md .

# Expose API port
EXPOSE 8080

# Run the API
CMD ["./email-api"]
