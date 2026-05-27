# 🛡️ Enterprise Email Intelligence Platform

A production-grade, 100% independent email verification and identity intelligence engine written in Go. This system operates with **zero third-party dependencies**, using in-house infrastructure fingerprinting and autonomous learning to detect disposable email services and fraud risks.

## 🚀 Key Features

### 📡 1. Deep Verification Sequence
-   **Syntax & RFC Compliance**: Standard validation against email formatting rules.
    -   **Real-time DNS/MX Resolution**: Verifies existence and health of mail exchange records.
    -   **SMTP Handshake**: Performs deep mailbox checking with **Greylisting Detection** (`4xx` soft-fails).
    -   **Catch-All Discovery**: Uses intelligent "Junk Probing" with a built-in safety whitelist for major providers (Google, Outlook, etc.).

### 🧠 2. Autonomous Infrastructure Intelligence
-   **MX Hub Fingerprinting**: Identifies disposable email providers (e.g., Mailinator, EmailOnDeck) by their underlying mail server patterns. Catch thousands of rotating domains instantly.
    -   **Self-Learning Discovery Pump**: A background worker that monitors validation telemetry and "proactively learns" new disposable domains from your live traffic.
    -   **Recursive Subdomain Protection**: Automatically matches subdomains against parent domain reputation (e.g., `test.sub.mailinator.com` matches `mailinator.com`).
    -   **Zero-Day Name Heuristics**: Advanced pattern matching to flag "throwaway" keyword domains even before they have active MX records.

### 📊 3. Professional Operations & Reporting
-   **Premium CLI Interface**: Columnar reports with pixel-perfect alignment and vivid status indicators.
    -   **Interactive TUI**: A real-time terminal dashboard with history and animated validation states.
    -   **REST API (V1)**: High-performance JSON endpoint for enterprise application integration.
    -   **Persistence**: SQLite-backed history and "Learned Intelligence" database.

## 🛠️ Installation & Build

```bash
# Clone and enter the directory (WSL recommended)
cd /home/sanjana/email-validator

# Build the CLI & API
go build -o email-validator ./cmd/cli/main.go
go build -o email-api ./cmd/api/main.go
```

## 📖 Usage Guide

### 1. Single Domain Check (CLI)
```bash
./email-validator check test@sub.mailinator.com
```

### 2. Dashboard Mode (TUI)
```bash
./email-validator interactive
```

### 3. Start the API Service
```bash
./email-api
```

### 4. Manual Intelligence Sync
```bash
# Triggers the Discovery Pump to scan for new infrastructure patterns
./email-validator sync
```

## 🧪 Testing & Verification Guide

| Test Scenario | Input Email | Expected Result | Signal to Look For |
| :--- | :--- | :--- | :--- |
| **Known Disposable** | `test@mailinator.com` | `Disposable: YES` | `Identified via Static Match` |
| **Infrastructure Hub** | `anything@sub.mailinator.com` | `Disposable: YES` | `Flagged via Infrastructure: Mailinator Hub` |
| **Zero-Day Pattern** | `test@burner-test-123.com` | `Disposable: YES` | `Zero-Day Heuristic: Suspicious Domain Pattern` |
| **Cross-Process Learning** | `test@new-one.com` | `Disposable: YES` | First check via CLI, then check via API to see `Static Match (DB Learned)` |
| **Identity Fraud** | `admin.github@gmail.com` | `Low Score` | `Fraud Risk: Suspicious Local-Part Keyword` |

## 📐 Technical Architecture
-   **Core**: `internal/core/validator.go` (SMTP/DNS logic)
-   **Intelligence**: `internal/intelligence/` (MX Hubs, Discovery Pump, Trust scoring)
-   **Service**: `internal/service/validator_service.go` (Unified logic for all interfaces)
-   **Persistence**: `emails.db` (SQLite)

---
Developed with 🛡️ for high-stakes identity verification.
