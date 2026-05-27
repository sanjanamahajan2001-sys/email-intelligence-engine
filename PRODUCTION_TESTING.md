# 🧪 Production Testing Guide: Load & Worker Pool Validation

This guide provides exact commands for verifying the system's resilience, concurrency, and graceful shutdown capabilities.

---

## 1. Setup: Load Test Environment
To test the worker pool with 30+ concurrent requests, we need a significant amount of "stale" data.

### 📥 Step 1: Seed 50 Stale Records
Run this command in your terminal to create 50 entries that require re-verification (Dated 10 days ago):

```bash
sqlite3 emails.db <<EOF
WITH RECURSIVE cnt(x) AS (
    SELECT 1
    UNION ALL
    SELECT x+1 FROM cnt LIMIT 50
)
INSERT INTO scans (email, lifecycle_state, created_at, is_valid)
SELECT 'test-' || x || '@example.com', 'ACTIVE', datetime('now', '-10 days'), 1
FROM cnt;
EOF
```

---

## 2. Concurrency & Load Testing

### 🚀 Command 1: High-Concurrency Re-verification (30 Workers)
Run the maintenance command with a concurrency level of 30. This replicates a high-load production re-validation event.

```bash
./email-validator update --concurrency=30
```

### 🎯 Expected Validation (Output):
1.  **Worker Pool Activation**: "🚀 Starting maintenance update for 50 emails (Worker Pool Size: 30)"
2.  **Transition Report**: A table showing `Old -> New` transitions for 50 records.
3.  **No Locking Errors**: Thanks to the **WAL mode** update, you should see zero "database is locked" errors during this heavy parallel write.

---

## 3. Resilience & Hardening Tests

### 🛑 Test 2: Graceful Shutdown (SIGINT)
During the update run, manually trigger a cancellation to test the signal handling.

1.  Start the update:
    ```bash
    ./email-validator update --concurrency=5
    ```
2.  Immediately press **`Ctrl + C`**.

### 🎯 Expected Validation (Output):
- The terminal says `(Press Ctrl+C to stop gracefully)`.
- On interruption, the process should finish existing workers and exit cleanly within seconds, rather than hanging or crashing.
- Inspect the database (`sqlite3 emails.db "SELECT count(*) FROM scans;"`) to ensure no corrupt partial records exist.

### 🧹 Test 3: Rate Limit Pruning (Memory Safety)
To verify that the memory pruning service is working:

1.  Make 10 unique requests via Curl (Wait 1 minute between them if necessary to hit limits, but here we test **Memory Cleanup**).
2.  Wait **11 minutes** (The pruning interval is 10 minutes).
3.  The `rateLimitMap` in the service layer will automatically purge entries for those IPs, ensuring memory usage stays flat.

---

## 4. Verification Checklists

- [ ] **WAL Mode**: Check `sqlite3 emails.db "PRAGMA journal_mode;"` (Should return `wal`).
- [ ] **Audit Integrity**: Run `./email-validator audit test-1@example.com` after the load test. You must see exactly two entries (Initial + Re-verified).
- [ ] **Timeout Control**: Start an update and observe the speed. If a network hang occurs, workers should cancel within 10s via the new `context` propagation.

---

> [!TIP]
> If you are testing in an environment that blocks Port 25, the transition report will show `ACTIVE -> ABANDONED` with a `DEMOTED` flag for all test emails, as SMTP connectivity will fail. This is normal behavior for the validator when restricted.
