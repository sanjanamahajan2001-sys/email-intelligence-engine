INSERT INTO scans (email, base_email, is_valid, syntax, dns, smtp, domain_age_years, reputation_score, risk_level, provider, lifecycle_state, created_at) 
VALUES ('stale-test@gmail.com', 'stale-test@gmail.com', 1, 1, 1, 1, 30.0, 100, 'Low', 'Google Workspace', 'ACTIVE', datetime('now', '-31 days'));
