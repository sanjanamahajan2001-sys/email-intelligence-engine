# 🛡️ Comprehensive Document: Enterprise Email Intelligence Architecture

This document provides a technical deep-dive into the platform's internal logic, scoring, and production operations.

---

## 🏗️ 1. Technical Architecture & Data Layers

The platform is designed as a zero-dependency, local-first intelligence engine.

- **Persistence Layer**: Structured via SQLite (`emails.db`). All raw SMTP responses and intelligence signals are stored to generate a high-fidelity audit trail.
- **Service Layer**: A multi-threaded orchestrator that coordinates syntax, DNS, SMTP, and Intelligence modules.
- **Intelligence Modules**: 
    - `Disposable`: Static domain blocking and Infrastructure-aware Hub matching.
    - `Identity`: RDAP/WHOIS age analysis and OSINT signal triangulation.
    - `Engagement`: Behavioral heuristics (success-rate history, mailbox quota).

---

## 🔄 2. Global Lifecycle States

The validator maps every identity into a distinct stage of its lifecycle.

| Lifecycle State | Description | Transition Rule |
| :--- | :--- | :--- |
| **ACTIVE** | Trustworthy, active identity | SMTP Success + No Intelligence Flags |
| **CATCH-ALL** | Enterprise gateway | SMTP junk probe accepted + Enterprise MX |
| **STALE** | Outdated intelligence | Age > Tiered Threshold (7-90 days) |
| **ABANDONED** | Previously valid, now persistent failure | Was ACTIVE -> Now persistent 550 error |
| **FULL** | Mailbox quota exceeded | SMTP 5.2.2 response |
| **INVALID** | Syntax, DNS, or Fatal SMTP failure | Syntax RFC fail or No MX found |

---

## ⚒️ 3. Automated Maintenance & Tiered Aging

The platform enforces a production-grade maintenance lifecycle to prevent data decay. The `update` command automatically identifies stale records using a tiered priority system.

### Tiered Aging Policy (Thresholds)
| Priority Tier | Lifecycle States | Maintenance Interval |
| :--- | :--- | :--- |
| **Tier 1 (Critical)** | ACTIVE, FULL | 7 Days |
| **Tier 2 (Stale)** | CATCH-ALL, STALE | 14 Days |
| **Tier 3 (Inactive)** | ABANDONED | 30 Days |
| **Tier 4 (Archive)** | INVALID, DISPOSABLE | 90 Days |

### High-Concurrency Worker Pool
The bulk maintenance engine scales re-verification across multiple parallel workers:
- **Default Concurrency**: 5 workers (optimized for sender reputation).
- **Production Tuning**: Use `--concurrency=N` to scale up to 50+ workers in high-bandwidth environments.
- **Resilience**: Integrated signal handling (SIGINT/SIGTERM) enables graceful stops without database corruption.

---

## 📡 3. REST API Interface (Production Reference)

### Endpoints
- `POST /v1/validate`: Primary identity verification endpoint.
- `POST /v1/sync-disposable`: Trigger global intelligence discovery.
- `GET /v1/health`: System health and uptime.

### Production CURL Examples
#### Standard Intelligence Check (Audit-Trail Aware)
```bash
curl -X POST "http://localhost:8080/v1/validate" \
     -H 'Content-Type: application/json' \
     -d '{"email": "contact@stripe.com"}'
```

#### Forced Intelligence Refresh (Bypass All Caches)
```bash
curl -X POST "http://localhost:8080/v1/validate?force=true" \
     -H 'Content-Type: application/json' \
     -d '{"email": "contact@stripe.com"}'
```

#### Global Hub Sync
```bash
curl -X POST "http://localhost:8080/v1/sync-disposable"
```

---

## 📊 5. Data Sovereignty & Export

To support external SIEM, Excel, or BI analysis, the platform provides a robust export engine.

### Secure Local Storage
By default, all exports are stored in the user's home directory to remain independent of the application repository:
- **Location**: `~/email-exports/`
- **Permissions**: `0755` (Directory), `0644` (Files)

### Export Schema (CSV)
The export includes up to 50,000 historical records with the following enterprise headers:
- `Email`, `Base_Email`, `Alias`, `IsValid`
- `Lifecycle_State`, `Reputation_Score`, `Engagement_Prob`
- `Risk_Level`, `Provider`, `Tld_Trust`
- `Syntax`, `DNS`, `SMTP`, `Disposable`, `CatchAll`, `Role`
- `Domain_Age_Yrs`, `Identity_Age_Yrs`, `Source`, `Date_Verified`

---

## 🕒 6. Operational Safety & Rate Limiting

The platform is hardened for production via two safety mechanisms:

1. **Per-IP Rate Limiting**: The system enforces a strict limit of **5 requests per minute per IP source**. Requests exceeding this will return `429 Too Many Requests`.
2. **SMTP Jitter & Human Simulation**: Forced re-validations include a randomized `100ms - 300ms` delay to prevent IP blacklisting.
3. **Audit Trail Deduplication**: Standard requests for the same identity within a **5-minute window** are de-duplicated if the state hasn't changed.

---

### Final Verification Summary
The platform is fully audited against all security and intelligence benchmarks. It displays consistent behavior across CLI, TUI, and API interfaces.
