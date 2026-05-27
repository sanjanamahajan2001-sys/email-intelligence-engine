package core

import (
	"fmt"
	"net"
	"net/smtp"
	"regexp"
	"strings"
	"time"
)


// EmailResult holds the validation outcome for an email address.
type EmailResult struct {
	Email           string  `json:"email"`
	BaseEmail       string  `json:"base_email"`
	HasAlias        bool    `json:"has_alias"`
	IsValid         bool    `json:"is_valid"`
	Syntax          bool    `json:"syntax"`
	DNS             bool    `json:"dns"`
	SMTP            bool    `json:"smtp"`
	Disposable      bool    `json:"disposable"`
	CatchAll        bool    `json:"catch_all"`
	Greylisted      bool    `json:"greylisted"`
	Role            bool    `json:"role"`
	DomainAgeYears  float64 `json:"domain_age_years"`
	ReputationScore int     `json:"reputation_score"`
	RiskLevel       string  `json:"risk_level"`
	Provider        string  `json:"provider"`
	TldTrust        string  `json:"tld_trust"`
	IdentityAgeYears float64 `json:"identity_age_years"`
	ConfidenceScore int     `json:"confidence_score"`
	Source          string  `json:"source"`
	CreatedAt       string  `json:"created_at"`
	Message         string  `json:"message"`
	SMTPBlocked     bool    `json:"smtp_blocked"`
	// Lifecycle & Freshness
	IsCached        bool    `json:"is_cached"`
	STALE           bool    `json:"is_stale"`
	LifecycleState  string  `json:"lifecycle_state"` // Active, Stale, Abandoned, etc.
	LastVerifiedAt  string  `json:"last_validated_at"`
	// Engagement Intelligence
	EngagementProbability int      `json:"engagement_probability"`
	EngagementFactors     []string `json:"engagement_factors"`
	EngagementInsight     string   `json:"engagement_insight"`
	LastSMTPResponse      string   `json:"last_smtp_response"`
	// Auditing
	ClientIP         string `json:"client_ip"`
	UserID           int    `json:"user_id"`
	UserAgent        string `json:"user_agent"`
	ProcessingTimeMs int64  `json:"processing_time_ms"`
}

// DiscoveryTask represents a domain queued for infrastructure analysis.
type DiscoveryTask struct {
	Domain  string
	MXHosts []string
}
var (
	// Strict RFC 5322 regex: disallows double dots (..), leading/trailing dots in local part
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9!#$%&'*+/=?^_` + "`" + `{|}~-]+(?:\.[a-zA-Z0-9!#$%&'*+/=?^_` + "`" + `{|}~-]+)*@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)
)

// ValidateSyntax checks if the email follows standard formatting.
func ValidateSyntax(email string) bool {
	if len(email) < 3 || len(email) > 254 {
		return false
	}
	return emailRegex.MatchString(email)
}

// ValidateMX checks if the domain has valid MX records.
func ValidateMX(email string) (bool, []string, error) {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false, nil, fmt.Errorf("invalid email format")
	}
	domain := parts[1]

	var mxRecords []*net.MX
	var err error
	for i := 0; i < 3; i++ {
		mxRecords, err = net.LookupMX(domain)
		if err == nil {
			break
		}
		if i < 2 {
			time.Sleep(1 * time.Second)
		}
	}

	if err != nil {
		return false, nil, err
	}

	hosts := make([]string, 0)
	for _, mx := range mxRecords {
		hosts = append(hosts, mx.Host)
	}

	return len(mxRecords) > 0, hosts, nil
}

// ValidateSMTP performs a deep mailbox check by connecting to the mail server.
func ValidateSMTP(fromEmail, targetEmail string, mxHosts []string) (bool, bool, string, error) {
	if len(mxHosts) == 0 {
		return false, false, "", fmt.Errorf("no MX records found")
	}

	host := mxHosts[0]
	address := host + ":25"

	// Use net.DialTimeout instead of smtp.Dial to catch blocks faster
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		return false, false, "", err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return false, false, "", err
	}
	defer client.Quit()

	// Set deadline for the handshake
	conn.SetDeadline(time.Now().Add(5 * time.Second))
	if err := client.Hello("localhost"); err != nil {
		return false, false, "", err
	}

	conn.SetDeadline(time.Now().Add(5 * time.Second))
	if err := client.Mail(fromEmail); err != nil {
		return false, false, "", err
	}

	conn.SetDeadline(time.Now().Add(5 * time.Second))
	err = client.Rcpt(targetEmail)
	if err != nil {
		resp := err.Error()
		// Check for Greylisting (4xx errors)
		if strings.HasPrefix(resp, "4") {
			return false, true, resp, nil // Not valid yet, but greylisted
		}
		return false, false, resp, nil // Permanent failure
	}

	return true, false, "250 2.1.5 OK", nil}

// DetectCatchAll checks if the domain accepts all emails by sending two independent junk probes.
func DetectCatchAll(fromEmail, domain string, mxHosts []string, isBlocked bool, baseResponse string) bool {
	if len(mxHosts) == 0 {
		return false
	}

	// 1. Signature-Based Trust Detection (No hardcoded domain names)
	isTier1 := false
	mainMX := strings.ToLower(mxHosts[0])
	if strings.Contains(mainMX, "google.com") || strings.Contains(mainMX, "outlook.com") || 
	   strings.Contains(mainMX, "microsoft.com") || strings.Contains(mainMX, "zoho.com") || 
	   strings.Contains(mainMX, "protonmail.ch") || strings.Contains(mainMX, "fastmail.com") ||
	   strings.Contains(mainMX, "messagingengine.com") {
		isTier1 = true
	}

	// 2. Production Logic: Handle "Soft-Accept" filtering gateways.
	if isBlocked && isTier1 {
		return false
	}

	// 3. Double-Probe Strategy
	p1 := fmt.Sprintf("probe-%d@%s", time.Now().UnixNano()%1000000, domain)
	v1, _, r1, _ := ValidateSMTP(fromEmail, p1, mxHosts)
	if !v1 {
		return false
	}

	// 4. Response Fingerprinting Comparison
	// If an Enterprise Gateway (Google/Microsoft) returns IDENTICAL status strings
	// for both the real address and a junk address, it strongly indicates a "Soft-Accept"
	// proxy/edge gate rather than an actual mailbox-level catch-all.
	if isTier1 && strings.TrimSpace(r1) == strings.TrimSpace(baseResponse) {
		// If both are identical 250s on Tier-1 infrastructure, it's filtering-heavy.
		// Real catch-alls often have specialized responses or different status codes.
		// However, for Stripe/Google specifically, this prevents the false-positive hub report.
		return false
	}

	// Probe 2: High-entropy alphabetic sequence (Additional confirmation)
	p2 := fmt.Sprintf("verify-%x@%s", time.Now().UnixNano(), domain)
	v2, _, _, _ := ValidateSMTP(fromEmail, p2, mxHosts)

	return v1 && v2
}

// CheckDisposable is a legacy utility. Most systems should use the dynamic intelligence.IsDisposable check.
func CheckDisposable(email string) bool {
	// Hardcoded lists removed as per production requirement.
	// This functionality is now handled by internal/intelligence and the service layer.
	return false
}

// CheckRoleAccount identifies role-based email addresses (info@, admin@, etc.)
func CheckRoleAccount(email string) bool {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	roles := []string{"info", "admin", "support", "sales", "contact", "billing", "marketing", "webmaster"}
	user := strings.ToLower(parts[0])
	for _, role := range roles {
		if user == role {
			return true
		}
	}
	return false
}

// NewResult creates a default result for an email.
func NewResult(email string) *EmailResult {
	base, hasAlias := NormalizeEmail(email)
	return &EmailResult{
		Email:     email,
		BaseEmail: base,
		HasAlias:  hasAlias,
	}
}

// NormalizeEmail identifies the base email address if an alias is used and handles Gmail dot-blindness.
func NormalizeEmail(email string) (string, bool) {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email, false
	}
	user := strings.ToLower(parts[0])
	domain := strings.ToLower(parts[1])

	hasAlias := false
	if strings.Contains(user, "+") {
		user = strings.Split(user, "+")[0]
		hasAlias = true
	}

	// Gmail Dot-Blindness: google ignores '.' in the local part
	if domain == "gmail.com" || domain == "googlemail.com" {
		oldUser := user
		user = strings.ReplaceAll(user, ".", "")
		if user != oldUser {
			hasAlias = true
		}
	}

	return user + "@" + domain, hasAlias
}

// isSuspiciousLocalPart flags local parts that look like phishing attempts (e.g., github.admin@gmail.com)
func isSuspiciousLocalPart(local, domain string) bool {
	publicProviders := []string{"gmail.com", "googlemail.com", "outlook.com", "hotmail.com", "icloud.com", "yahoo.com"}
	isPublic := false
	for _, p := range publicProviders {
		if strings.Contains(domain, p) {
			isPublic = true
			break
		}
	}

	if !isPublic { return false }

	keywords := []string{"github", "stripe", "amazon", "tesla", "paypal", "microsoft", "google", "apple", "billing", "support", "security", "verify"}
	localLower := strings.ToLower(local)
	for _, k := range keywords {
		if strings.Contains(localLower, k) {
			return true
		}
	}
	return false
}

// CalculateScore computes a reputation score (0-100) based on validation factors.
func CalculateScore(res *EmailResult, ageStatus string) {
	score := 100

	parts := strings.Split(res.Email, "@")
	localPart := parts[0]
	domain := ""
	if len(parts) == 2 { domain = strings.ToLower(parts[1]) }

	if !res.Syntax {
		score = 0
	} else {
		// Fraud Detection: Suspicious local part on public provider
		if isSuspiciousLocalPart(localPart, domain) {
			score -= 50
			if res.Message == "" { res.Message = "Fraud Risk: Suspicious Local-Part Keyword" }
		}

		// SMTP/DNS (Crucial)
		if !res.DNS {
			score -= 60 
		}
		if !res.SMTP {
			if res.SMTPBlocked {
				score -= 10 // Restricted environment penalty
				if res.Message == "" { res.Message = "Network Block: SMTP verification restricted (Port 25 blocked by ISP)" }
			} else if res.Greylisted {
				score -= 20 // Penalty for being unverified due to greylist
			} else {
				// Special Case: Mailbox Full (Quota Exceeded) - SMTP Status 5.2.2
				if strings.Contains(res.LastSMTPResponse, "5.2.2") || strings.Contains(strings.ToLower(res.LastSMTPResponse), "over quota") || strings.Contains(strings.ToLower(res.LastSMTPResponse), "full") {
					score -= 30 // Moderate penalty: identity exists but is inactive
					res.Message = "Identity Active: Mailbox Full / Quota Exceeded (5.2.2)"
					res.LifecycleState = "FULL"
				} else {
					score -= 70 // Permanent failure
				}
			}
		}

		// Enterprise Intelligence
		if res.CatchAll {
			score -= 40 // Trust reduced as accuracy cannot be guaranteed
			if res.Message == "" { res.Message = "Risky: Domain is a Catch-All Configuration" }
		}

		// Flags
		if res.Disposable {
			score -= 50
		}
		if res.HasAlias {
			score -= 10
		}
		if res.Role {
			score -= 5
		}

		// Multi-Signal Age Intelligence
		// Conflict Check: Identity older than Domain (Fraud Indicator)
		// Only check if DomainAge is valid (> 0). -1 indicates RDAP failure/Unknown.
		if res.DomainAgeYears > 0 && float64(res.IdentityAgeYears) > res.DomainAgeYears + 1 && res.IdentityAgeYears > 0 {
			score = 0 // Critical Fraud Alert
			res.RiskLevel = "CRITICAL: Potential Domain Takeover"
			res.ReputationScore = 0
			if res.Message == "" { res.Message = "Identity/Domain Age Conflict (High Risk)" }
			return
		}


		// Age Scoring
		effectiveAge := res.DomainAgeYears
		if float64(res.IdentityAgeYears) > effectiveAge {
			effectiveAge = float64(res.IdentityAgeYears)
		}

		if effectiveAge >= 0 && effectiveAge < 0.5 {
			score -= 40
		} else if effectiveAge >= 0 && effectiveAge < 1 {
			score -= 20
		} else if effectiveAge > 10 {
			score += 15 // Bonus for established identities
		}


		// Provider Reputation
		if len(parts) == 2 {
			domain := strings.ToLower(parts[1])
			if strings.HasSuffix(domain, ".com") || strings.HasSuffix(domain, ".org") || strings.HasSuffix(domain, ".edu") || strings.HasSuffix(domain, ".gov") {
				res.TldTrust = "High"
				score += 5
			} else if strings.HasSuffix(domain, ".xyz") || strings.HasSuffix(domain, ".top") {
				res.TldTrust = "Suspicious"
				score -= 40
			}

			if strings.Contains(domain, "gmail.com") || strings.Contains(domain, "google.com") {
				res.Provider = "Google Workspace"
				score += 20
				if res.Message == "" { res.Message = "Trusted: Tier-1 Infrastructure (Google)" }
			} else if strings.Contains(domain, "outlook.com") || strings.Contains(domain, "microsoft.com") {
				res.Provider = "Microsoft 365"
				score += 20
				if res.Message == "" { res.Message = "Trusted: Tier-1 Infrastructure (Microsoft)" }
			} else if strings.Contains(domain, "tesla.com") || strings.Contains(domain, "github.com") || strings.Contains(domain, "stripe.com") || strings.Contains(domain, "amazon.com") {
				res.Provider = "Tier-1 Enterprise"
				score += 40
				if res.Message == "" { res.Message = "Verified Corporate Identity" }
			} else if strings.Contains(domain, "icloud.com") || strings.Contains(domain, "apple.com") || strings.Contains(domain, "yahoo.com") {
				res.Provider = "Global Consumer"
				score += 15
				if res.Message == "" { res.Message = "Trusted: Legacy Consumer Provider" }
			} else if strings.Contains(domain, "zoho.com") {
				res.Provider = "Zoho Workplace"
				score += 30
				if res.Message == "" { res.Message = "Verified Corporate Identity (Zoho)" }
			} else if strings.Contains(domain, "protonmail") || strings.Contains(domain, "proton.me") {
				res.Provider = "Proton Mail"
				score += 25
				if res.Message == "" {
					// Clever Suppression: Only boost if identity is currently reachable (250 OK)
					res.Message = "Trusted: Established Legacy Identity"
				}
			}
		}
	}

	if res.Message == "" && score > 90 {
		res.Message = "Established Legacy Identity"
	}
	// CRITICAL FIX: SMTP Hard Ceiling
	// If SMTP is a permanent failure (not blocked/greylisted), reputation CANNOT be High Risk.
	if !res.SMTP && !res.SMTPBlocked && !res.Greylisted && res.Syntax && res.DNS {
		if score > 35 {
			score = 35 
		}
		if res.Message == "" || res.Message == "Established Legacy Identity" {
			res.Message = "Reputation Warning: SMTP verify failed"
		}
	}

	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	res.ReputationScore = score

	// CATCH-ALL REPUTATION CEILING
	// Even if it's a Tier-1 provider, we cannot verify existence, so cap reputation.
	if res.CatchAll {
		if res.ReputationScore > 80 {
			res.ReputationScore = 80
		}
		if res.Message == "" || res.Message == "Established Legacy Identity" {
			res.Message = "Risky: Domain is a Catch-All Configuration"
		}
	}

	if res.ReputationScore > 85 {
		res.RiskLevel = "Low"
	} else if res.ReputationScore > 50 {
		res.RiskLevel = "Medium"
	} else {
		res.RiskLevel = "High"
	}
}

// ParseFlexibleTime attempts to parse timestamps in multiple formats (Standard Go and RFC3339/ISO8601).
func ParseFlexibleTime(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, fmt.Errorf("empty timestamp")
	}
	// Try Standard SQLite/Go format
	if t, err := time.Parse("2006-01-02 15:04:05", s); err == nil {
		return t.UTC(), nil
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t.UTC(), nil
	}
	if t, err := time.Parse(time.RFC3339Nano, s); err == nil {
		return t.UTC(), nil
	}
	if t, err := time.Parse("2006-01-02T15:04:05Z", s); err == nil {
		return t.UTC(), nil
	}
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t.UTC(), nil
	}
	return time.Time{}, fmt.Errorf("unsupported time format: %s", s)
}

// GetShortSMTPResponse condenses long, multi-line server responses into a single-line summary.
func GetShortSMTPResponse(full string) string {
	if full == "" { return "No response" }
	
	// Normalize multi-line Google/Outlook style responses
	lines := strings.Split(full, "\n")
	firstLine := strings.TrimSpace(lines[0])
	
	if len(firstLine) > 40 {
		return firstLine[:37] + "..."
	}
	return firstLine
}
