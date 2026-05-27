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
