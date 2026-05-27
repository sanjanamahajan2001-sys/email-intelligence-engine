#!/bin/bash

# Exit on error
set -e

echo "🚀 Starting Production API Build..."

# Define paths
SOURCE_DIR="./cmd/api"
DIST_DIR="./api-dist"
BINARY_NAME="email-api"

# Ensure Go is in PATH
if ! command -v go &> /dev/null; then
    if [ -f "/usr/local/go/bin/go" ]; then
        export PATH=$PATH:/usr/local/go/bin
    else
        echo "❌ Error: Go is not installed or not in PATH."
        exit 1
    fi
fi

# Clean previous builds
rm -rf $DIST_DIR
mkdir -p $DIST_DIR

# 1. Compile the API and CLI (Linux/WSL)
echo "📦 Compiling binaries (Linux optimized)..."
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $DIST_DIR/$BINARY_NAME $SOURCE_DIR
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $DIST_DIR/email-validator ./cmd/cli

# 2. Copy necessary assets
echo "📂 Packaging configuration and database templates..."
cp README.md $DIST_DIR/
cp SHARE_GUIDE.md $DIST_DIR/
cp docker-compose.yml $DIST_DIR/

# 3. Create a Standalone Dockerfile for distribution (Compatible with Ubuntu binary)
cat <<EOF > $DIST_DIR/Dockerfile
FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y libsqlite3-0 ca-certificates && rm -rf /var/lib/apt/lists/*
WORKDIR /root/
# Copy already-built binaries
COPY email-api .
COPY email-validator .
COPY README.md .
COPY SHARE_GUIDE.md .
EXPOSE 8080
CMD ["./email-api"]
EOF

# 4. Create a default .env if it doesn't exist (using config defaults)
cat <<EOF > $DIST_DIR/.env
API_PORT=8080
DB_PATH=emails.db
SMTP_SENDER=sanjanamaahi2001@gmail.com
RATE_LIMIT_IP_MIN=10
JWT_SECRET=$(openssl rand -base64 32 2>/dev/null || echo "default-secret-key-$(date +%s)")
ACCESS_TOKEN_DURATION=15m
REFRESH_TOKEN_DURATION=168h
EOF

# 5. Create a docker-compose.yml for easy deployment
cat <<EOF > $DIST_DIR/docker-compose.yml
services:
  email-validator-api:
    build: .
    ports:
      - "8080:8080"
    volumes:
      - .:/root/
    environment:
      - API_PORT=8080
      - DB_PATH=emails.db
      - SMTP_SENDER=sanjanamaahi2001@gmail.com
      - RATE_LIMIT_IP_MIN=10
      - ACCESS_TOKEN_DURATION=15m
      - REFRESH_TOKEN_DURATION=168h
    restart: unless-stopped
EOF

# Ensure emails.db exists as a file (prevents Docker directory-mount bug)
touch $DIST_DIR/emails.db

# 5. Create a simple runner script
cat <<EOF > $DIST_DIR/start.sh
#!/bin/bash
./email-api
EOF
chmod +x $DIST_DIR/start.sh

echo "✅ Build Complete! Your distribution is ready in: $DIST_DIR"
echo "--------------------------------------------------------"
echo "To share: Zip the '$DIST_DIR' folder and send it."
echo "--------------------------------------------------------"
