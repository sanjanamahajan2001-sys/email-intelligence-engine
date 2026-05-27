package intelligence

import (
	"strings"
)

// TrustScore represents the calculated trust for a domain.
type TrustScore struct {
	Score       int
	Reason      string
	IsLowTrust  bool
}

// GetDomainTrust calculates a local "Domain Authority" proxy.
func GetDomainTrust(domain string, ageYears float64) TrustScore {
	score := 100
	var reasons []string

	// 1. Age Factor (Critical for Disposable detection)
	if ageYears >= 0 {
		if ageYears < 0.1 { // Less than ~36 days
			score -= 50
			reasons = append(reasons, "Very New Domain")
		} else if ageYears < 0.5 { // Less than 6 months
			score -= 20
			reasons = append(reasons, "Recently Registered")
		} else if ageYears > 2.0 {
			score += 10 // Trust boost for established domains
		}
	}

	// 2. TLD Reputation Factor
	domain = strings.ToLower(domain)
	suspiciousTLDs := map[string]int{
		".tk":  40,
		".ml":  40,
		".ga":  40,
		".cf":  40,
		".gq":  40,
		".xyz": 20,
		".icu": 20,
		".top": 15,
		".pw":  15,
	}

	for tld, penalty := range suspiciousTLDs {
		if strings.HasSuffix(domain, tld) {
			score -= penalty
			reasons = append(reasons, "Low-Trust TLD: "+tld)
			break
		}
	}

	// 3. Length/Entropy Heuristic (Common in auto-generated domains)
	parts := strings.Split(domain, ".")
	if len(parts) > 0 {
		name := parts[0]
		if len(name) > 15 {
			score -= 10
			reasons = append(reasons, "Unusually Long Name")
		}
	}

	if score < 0 {
		score = 0
	}

	return TrustScore{
		Score:      score,
		Reason:     strings.Join(reasons, ", "),
		IsLowTrust: score < 50,
	}
}
