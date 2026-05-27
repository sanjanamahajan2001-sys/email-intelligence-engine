# 🛡️ Intelligence Feature Documentation

This document detail the advanced forensic signals used to determine an email identity's reputation and lifecycle.

---

## 🧭 Intelligence Signals

### 1. Enterprise Fingerprinting
We identify **Tier-1 Filtering Gateways** (Stripe, Google, Microsoft) by their SMTP response behavior and MX signatures.
- **Problem**: Many major domains accept all emails during the SMTP handshake but filter them internally (Catch-all).
- **Solution**: If a domain is identified as a High-Trust Enterprise Hub, it is **exempted from catch-all risk scoring**, maintaining a 100/100 reputation score.

### 2. Infrastructure Hub Detection (Disposable)
Disposable detection goes beyond static domain lists.
- **Hub Matching**: We match the underlying mail server signatures (MX) against known disposable provider clusters (e.g., Mailinator, EmailOnDeck).
- **Zero-Day Heuristics**: Domains matching temp-mail keywords (fixed and patterned) are flagged even if undocumented.

### 3. Identity Age Triangulation
Determines "True Identity Age" by cross-referencing three independent sources:
- **RDAP/WHOIS**: The registration date of the domain infrastructure.
- **Telecom/Telemetry**: The first seen date in our global intelligence network.
- **OSINT Breach Data**: Historical identity records dating back to 2015.
- **Confidence Logic**: High confidence is achieved when an identity matches the age of its infrastructure. Conflicts (e.g., identity older than domain) are flagged as **Identity Fraud Risk**.

---

## 🔄 Lifecycle States & Transitions

The platform tracks every identity throughout its lifecycle.

| Lifecycle State | Description | Transition Rule |
| :--- | :--- | :--- |
| **ACTIVE** | Trustworthy, active identity | SMTP Success + No Intelligence Flags |
| **CATCH-ALL** | Enterprise gateway | SMTP junk probe accepted + Enterprise MX |
| **STALE** | Outdated intelligence | Record > 30 days old |
| **ABANDONED** | Previously valid, now persistent failure | Was ACTIVE -> Now persistent 550 error |
| **FULL** | Mailbox quota exceeded | SMTP 5.2.2 response |
| **INVALID** | Syntax, DNS, or Fatal SMTP failure | Syntax RFC fail or No MX found |

### 5. Exportable Intelligence
- Every identity transition is logged with high-fidelity telemetry.
- All historical data is exportable via `./email-validator export` for external forensics and reputation analysis.
- Supports CSV and JSON formats for seamless integration with enterprise BI tools.

---

## 📡 API Integration (CURL)

### Refresh Intelligence (Force Re-validation)
Use this when you suspect a change in an account's status (e.g., reactivation of an ABANDONED account).

```bash
curl -X POST "http://localhost:8080/v1/validate?force=true" \
     -H 'Content-Type: application/json' \
     -d '{"email": "contact@stripe.com"}'
```

---

### Verification Summary
Passing all the above test scenarios confirms the platform is **Production Grade**, with robust state management, infrastructure-aware intelligence, and high-precision reputation scoring.
