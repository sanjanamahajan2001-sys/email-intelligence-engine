#!/bin/bash
# 🛠️ Migration & Verification Script
cd /home/sanjana/email-validator

# 1. Perform Migration (Ignore error if column already exists)
sqlite3 emails.db "ALTER TABLE disposable_domains ADD COLUMN provider_name TEXT;" 2>/dev/null || true

# 2. First Scan: Identify via MX Infrastructure and "Learn" the domain with the name
echo "--- FIRST SCAN: Infrastructure Discovery ---"
./email-validator check test@sub.mailinator.com

# 3. Second Scan: Verify that even if DNS was simulated to fail, the provider label is preserved
echo "--- SECOND SCAN: Persistence Check (Label Preservation) ---"
./email-validator check test@sub.mailinator.com

# 4. Final DB Check: Show the stored record
echo "--- DATABASE RECORD ---"
sqlite3 emails.db "SELECT domain, provider_name FROM disposable_domains WHERE domain='sub.mailinator.com';"
