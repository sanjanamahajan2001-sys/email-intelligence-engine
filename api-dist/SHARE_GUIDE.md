# 🛡️ Email Validator API - Integration & Sandbox Guide

This guide is for the engineering team responsible for integrating the Email Validator API into the company's registration frontends or security sandbox.

---

## 🏗️ 1. Architecture Overview
The API is a "Black Box" service. It does not require a frontend; it exposes a REST interface that your application calls to verify the authenticity of an email address.

**Primary Endpoint**: `POST /v1/web-validate`

---

## 🚀 2. Getting Started (Sandbox Setup)

You have been provided with a distribution folder (`api-dist/`). Unzip this folder and open your terminal inside it.

### A. The Containerized Way (Recommended for High Availability)
If you have Docker installed, this is the safest and most portable method. It handles dependencies and restarts automatically.

```bash
# 1. Enter the package directory
cd api-dist/

# 2. Build and start the containerized service
docker-compose up -d --build

# 3. Verify it is running
docker ps
```
*Note: Your local port `8080` maps to the API's port `8080`.*

### B. The Standalone Binary Way (Background Mode)
If you want to run the API in the background so it stays alive even after you close your terminal or logout:

```bash
# 1. Enter the package directory
cd api-dist/

# 2. Start the API in the background with logs directed to api.log
nohup ./start.sh > api.log 2>&1 &

# 3. To stop the API later:
lsof -ti:8080 | xargs kill -9
```

---

## 🔐 3. Authentication

The API is secured using JWT (JSON Web Tokens). All requests to the validation endpoints must include a valid Bearer token.

### A. How to get a Token
Use the `/v1/auth/login` endpoint with your credentials. 
*Note: A default user `admin` with password `admin123` is created on the first run if no other users exist.*

```bash
curl -X POST "http://localhost:8080/v1/auth/login" \
     -H 'Content-Type: application/json' \
     -d '{"username": "admin", "password": "admin123"}'
```
**Response**:
```json
{
  "access_token": "eyJhbGci...",
  "refresh_token": "def456...",
  "expires_in": 900
}
```

### B. How to use a Token
Include the `access_token` in the `Authorization` header as a `Bearer` token for all subsequent requests.

### C. How to Refresh a Token
If your `access_token` expires (401 error), use your `refresh_token` to get a new pair without logging in again:
```bash
curl -X POST "http://localhost:8080/v1/auth/refresh" \
     -H 'Content-Type: application/json' \
     -d '{"refresh_token": "your_refresh_token_here"}'
```

---

## 📡 4. Data Integration (Request & Response)

### Registration Flow Example
When a user submits an email on your web form, your backend should perform a server-side request to this API.

#### Request Example
```bash
# Set your token from the login step
TOKEN="your_access_token_here"

curl -X POST "http://localhost:8080/v1/web-validate" \
     -H "Authorization: Bearer $TOKEN" \
     -H 'Content-Type: application/json' \
     -d '{"email": "contact@stripe.com"}'
```

#### Detailed Response Structure
The API returns a rich JSON object. Use the following fields to drive your frontend logic:

| JSON Field | Purpose | Integration Advice |
| :--- | :--- | :--- |
| `is_valid` | Strict reachability | Use for form validation block. |
| `reputation_score` | 0-100 Trust Score | Score > 80 is High Trust. |
| `authenticity_status` | "Verified" / "Suspicious" | Use for UI badges (Green/Yellow). |
| `recommendation` | "Accept" / "Flag" / "Reject" | **Crucial**: Drive your registration logic with this. |
| `engagement.factors` | Evidence list | (Optional) Show to admins in dashboards. |

### 🚨 Common Status Codes
| Code | Meaning | Action |
| :--- | :--- | :--- |
| `200` | Success | Process the validation result. |
| `401` | Unauthorized | Access token expired. Trigger "Refresh" flow. |
| `429` | Too Many Requests | Rate limit hit. Wait 60 seconds before retrying. |
| `500` | Server Error | Internal issue. Check `api.log`. |

---

## 🏥 5. Health & Monitoring
For automated monitoring (uptime checks), use the dedicated health endpoint:
```bash
curl -i "http://localhost:8080/v1/health"
# Expected: 200 OK with {"status": "healthy"}
```

---

## 💾 6. Persistence & Database Audit

The API uses a local `emails.db` (SQLite) to store validation history and learn over time.

> [!IMPORTANT]
> **Data Safety**: All data is stored on the physical host (your machine). Even if the Docker container is deleted or the binary is stopped, your validation history is **SAFE** in the `emails.db` file.

### How to query the Scan History
You can use any SQLite browser or the command line to inspect the results.

**Prerequisite: Install SQLite3**
If the `sqlite3` command is not found on your system, install it using:
```bash
# Ubuntu / Debian / WSL
sudo apt-get update
sudo apt-get install sqlite3 -y

# CentOS / RHEL
sudo yum install sqlite
```

**Query Examples:**
```bash
# View last 10 scans with Audit Trail (IP and User ID)
sqlite3 emails.db "SELECT created_at, email, client_ip, user_id, reputation_score FROM scans ORDER BY id DESC LIMIT 10;"
```

### Manual Intelligence Sync
To manually trigger the "Learning Pump" and update the infrastructure intelligence:
```bash
# Call the internal sync endpoint
curl -X POST "http://localhost:8080/v1/sync-disposable"
```

---

## 🕵️ 7. Forensic Auditing
The API automatically captures environment metadata for every request to provide a complete audit trail.

*   **Client IP**: The `client_ip` field records the remote address that initiated the request.
*   **User ID**: The `user_id` field tracks exactly which internal API user authorized the request.
*   **User Agent**: Captures the software or tool used to call the API (e.g., Python-Requests, Chrome, or internal CLI).
*   **Latency (ms)**: High-resolution timing of exactly how long the validation cycle took for performance SLA monitoring.
*   **Timestamp (UTC)**: All audit records use the **UTC Standard (Universal Time Coordinated)** to ensure consistency across different server timezones and regional offices.

This data is stored permanently in the `scans` table, allowing security teams to audit usage patterns or identify compromised credentials.

---

---

## 🛠️ 8. Admin CLI & Forensic Tools (Terminal)

The distribution includes a powerful CLI tool called `email-validator`. This tool allows administrators to manage users and perform deep forensic analysis directly from the server.

### A. View Scan History
To see a live summary of recent requests across your entire user base.
**Usage:**
```bash
./email-validator history
# Or: ./email-validator history 50 (to see more)
```
**Example Output:**
```text
ID   | TIMESTAMP    | REQUESTER_IP    | EMAIL              | STATUS     | SCORE
------------------------------------------------------------------------------------
1    | Just now     | 203.0.113.1     | contact@stripe.com | ✅ VALID   | 100/100
2    | 15m ago      | 203.0.113.1     | info@example.com   | ❌ INVALID | 0/100
```

### B. Deep Forensic Audit
To see every lifecycle transition of a specific address. If an email goes from "ACTIVE" to "INVALID" or "BOUNCED," this trail shows you exactly when it happened and which user/IP was responsible.
**Usage:**
```bash
./email-validator audit contact@stripe.com
```
**Example Output:**
```text
📜 LIFECYCLE TRANSITIONS LOG: contact@stripe.com

TIMESTAMP            | STATUS       | CLIENT_IP       | UID   | TRANSITION
-------------------------------------------------------------------------------------
2026-04-16 06:13:17  | ACTIVE       | 203.0.113.1     | 2     | Initial Scan
```

### C. Master Data Export (CSV/JSON)
To export the full forensic database for compliance teams or manual review. The export includes all 17+ data points per scan, including internal infrastructure metadata.
**Usage:**
```bash
./email-validator export logs.csv
```
**Example Row structure:**
> `email, is_valid, syntax, dns, smtp, reputation, risk, provider, disposable, role, catch_all, domain_age, source, requester_ip, user_id, date_utc, engagement_factors`

### D. User Management
Manage API access for your team:
```bash
# Create a new user
./email-validator user create manager secure_pass_123

# List all users
./email-validator user list

# Reset a password
./email-validator user password-reset manager new_pass_456

# Delete a user (Critical: Revokes all active tokens)
./email-validator user delete manager
```

### E. Manual Intelligence Sync
To manually trigger the "Learning Pump" and update the infrastructure intelligence:
```bash
./email-validator sync
```

---

## 🏗️ 9. Integration Best Practices (Pro Recommendations)

### A. Secure Architecture (Server-to-Server)
**NEVER** call this API directly from a frontend (JavaScript/Browser). This would expose your credentials and tokens.
*   **Correct Flow**: Frontend Form -> **Your Backend** -> Email Validator API -> **Your Backend** -> Frontend Response.

### B. Automated Token Management
Your backend should handle the token lifecycle automatically to prevent manual intervention:
1.  **On Startup**: Your backend should call `/v1/auth/login` and store the `access_token` and `refresh_token` in memory.
2.  **During Requests**: Include the `access_token` in the header.
3.  **On 401 Unauthorized**: If a request fails with 401, catch the error, call `/v1/auth/refresh` using your `refresh_token`, update your held tokens, and retry the original request.

---

## ⚙️ 10. Configuration (.env)
...
...
You can customize the service behavior in the provided `.env` file:
- `API_PORT`: Customize the listening port (Default: 8080).
- `SMTP_SENDER`: The "From" address used for mailbox checking.
- `DB_PATH`: Location of the SQLite database file.
