# 🏆 Master Testing Guide: Secure Email Validator API

This guide provides a step-by-step sequence to verify all security, rate-limiting, and user management features implemented in the current version.

---

## 🛠️ Phase 1: Environment Preparation

1.  **Ensure no old processes are running**:
    ```bash
    lsof -ti:8080 | xargs kill -9 2>/dev/null || true
    ```

2.  **Enter the distribution folder**:
    ```bash
    cd ~/email-validator/api-dist
    ```

3.  **Wipe existing database to test Auto-Provisioning**:
    ```bash
    rm emails.db && touch emails.db
    ```

---

## 🔐 Phase 2: Start & Initial Authentication

1.  **Start the API**:
    ```bash
    ./email-api
    ```
    > [!NOTE]
    > Verify the logs show: `[INFO] Auth: AUTO-PROVISIONING: Created default admin user`.

2.  **Verify Access Control (Negative Test)**:
    In a new terminal, try a request without a token. It should fail.
    ```bash
    curl -X POST "http://localhost:8080/v1/web-validate" \
         -H 'Content-Type: application/json' \
         -d '{"email": "test@example.com"}'
    # Expected: {"error": "Authorization header required"} (401)
    ```

3.  **Login with Admin Credentials**:
    ```bash
    curl -X POST "http://localhost:8080/v1/auth/login" \
         -H 'Content-Type: application/json' \
         -d '{"username": "admin", "password": "admin123"}'
    # Expected: Successful JSON response with access_token and refresh_token.
    ```

---

## 👥 Phase 3: Advanced User Management

1.  **Create a New User via CLI**:
    ```bash
    ./email-validator user create manager pass456
    # Expected: ✅ User 'manager' created successfully!
    ```

2.  **Verify New User Login**:
    ```bash
    curl -X POST "http://localhost:8080/v1/auth/login" \
         -H 'Content-Type: application/json' \
         -d '{"username": "manager", "password": "pass456"}'
    # Expected: Success. Save the "access_token" to a variable:
    TOKEN="<YOUR_TOKEN_HERE>"
    ```

3.  **Reset Password**:
    ```bash
    ./email-validator user password-reset manager newsecret789
    ```

4.  **Verify Password Reset**:
    Try logging in with the OLD password first (should fail), then the NEW one (should succeed).

5.  **Delete a User**:
    ```bash
    ./email-validator user delete manager
    ```
    Try logging in with `manager` again (should fail).

---

## 🚦 Phase 4: Rate Limiting Verification

1.  **Perform rapid requests**:
    Using a valid token, run this command 10+ times in quick succession (or until blocked):
    ```bash
    for i in {1..15}; do curl -s -o /dev/null -w "%{http_code}\n" -X POST "http://localhost:8080/v1/web-validate" -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' -d '{"email": "test@stripe.com"}'; done
    ```
    > [!IMPORTANT]
    > You should see several `200` responses followed by `429` (Too many requests) once the bucket is empty.

---

## 📊 Phase 5: Database Persistence Audit

1.  **Ensure SQLite is installed**:
    ```bash
    sudo apt-get install sqlite3 -y
    ```

2.  **Query the users table**:
    ```bash
    sqlite3 emails.db "SELECT id, username FROM users;"
    ```

3.  **Query validation history**:
    ```bash
    sqlite3 emails.db "SELECT email, reputation_score FROM scans LIMIT 5;"
    ```

---

## 🏁 Phase 6: Sync & Maintenance

1.  **Manual Intelligence Sync**:
    ```bash
    curl -X POST "http://localhost:8080/v1/sync-disposable" -H "Authorization: Bearer $TOKEN"
    # Expected: Success message with count of domains.
    ```
