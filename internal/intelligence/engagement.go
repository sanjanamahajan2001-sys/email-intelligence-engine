package intelligence

import (
	"strings"
	"github.com/sanjana/email-validator/internal/core"
	"github.com/sanjana/email-validator/internal/db"
)

// CalculateEngagementProbability returns a score from 0-100 indicating the likelihood of a reply.
// It also populates the EngagementFactors and EngagementInsight fields in the result.
func CalculateEngagementProbability(res *core.EmailResult, database *db.DB) int {
	res.EngagementFactors = []string{}

	// 0. CRITICAL CHECK: Invalid Domain (DNS Fail)
	if !res.DNS {
		res.EngagementFactors = append(res.EngagementFactors, "Critical: Invalid Domain (DNS Failure)")
		res.EngagementInsight = "This domain does not exist or has no mail servers configured."
		res.EngagementProbability = 0
		return 0
	}

	// 1. SMTP Handshake & Bounce Analysis
	isPermanentBounce := false
	isSMTPBlocked := false
	lowerResp := strings.ToLower(res.LastSMTPResponse)
	
	if !res.SMTP && res.LastSMTPResponse != "" {
		// Check for actual "User Not Found" codes vs "IP Blocked"
		if strings.HasPrefix(res.LastSMTPResponse, "550") || strings.HasPrefix(res.LastSMTPResponse, "554") {
			if strings.Contains(lowerResp, "access denied") || 
			   strings.Contains(lowerResp, "banned") || 
			   strings.Contains(lowerResp, "blocked") ||
			   strings.Contains(lowerResp, "spamhaus") {
				isSMTPBlocked = true
			} else {
				isPermanentBounce = true
			}
		}
	}

	if isPermanentBounce {
		res.EngagementFactors = append(res.EngagementFactors, "Critical: Permanent Bounce (Engagement Zero)")
		res.EngagementInsight = "This address is unreachable (Hard Bounce). No engagement is possible."
		res.EngagementProbability = 0
		return 0
	}

	if res.Disposable {
		res.EngagementFactors = append(res.EngagementFactors, "Critical: Disposable/Temporary Domain")
		res.EngagementInsight = "Disposable addresses are short-lived and never yield engagement."
		res.EngagementProbability = 0
		return 0
	}

	// Calculate base score (starting point 35-50 based on provider)
	score := 35.0
	if strings.Contains(res.Provider, "Google") || strings.Contains(res.Provider, "Microsoft") {
		score = 50.0
		res.EngagementFactors = append(res.EngagementFactors, "+25: Tier-1 Infrastructure (Google/Microsoft)")
	} else if strings.Contains(res.Provider, "Enterprise") {
		score = 45.0
		res.EngagementFactors = append(res.EngagementFactors, "+15: Verified Enterprise Provider")
	}

	// 2. Age Factors
	if res.DomainAgeYears > 5 {
		score += 20
		res.EngagementFactors = append(res.EngagementFactors, "+20: Established Legacy Identity (>5 yrs)")
	} else if res.DomainAgeYears < 1 && res.DomainAgeYears >= 0 {
		score -= 20
		res.EngagementFactors = append(res.EngagementFactors, "-20: Brand New Identity (Low Trust)")
	}

	// 3. Connection Quality
	if res.SMTP {
		score += 10
		res.EngagementFactors = append(res.EngagementFactors, "+10: Active Handshake (250 OK)")
	} else {
		// Non-terminal failures
		if strings.Contains(lowerResp, "421") || strings.Contains(lowerResp, "450") || strings.Contains(lowerResp, "451") {
			score -= 15
			res.EngagementFactors = append(res.EngagementFactors, "-15: Temporary Bounce/Greylisting detected")
		}
		if isSMTPBlocked {
			score -= 30
			res.EngagementFactors = append(res.EngagementFactors, "-30: SMTP Delivery Blocked (Security Filtering)")
		}
	}

	// 4. Historical Success Rate
	if database != nil {
		hsr := database.GetHistoricalSuccessRate(res.Email)
		if hsr > 0.8 {
			score += 10
			res.EngagementFactors = append(res.EngagementFactors, "+10: Consistent Historical Activity")
		} else if hsr < 0.3 {
			score -= 20
			res.EngagementFactors = append(res.EngagementFactors, "-20: Poor Historical Success Rate")
		}
	}

	// 5. Risky Attributes (Penalties)
	if res.CatchAll {
		score -= 40
		res.EngagementFactors = append(res.EngagementFactors, "-40: Catch-all domain (Reduced Delivery Confidence)")
	}
	if res.Role {
		score -= 15
		res.EngagementFactors = append(res.EngagementFactors, "-15: Role-based address (info@/admin@)")
	}

	// Clamp the score
	finalScore := int(score)
	if finalScore > 100 { finalScore = 100 }
	if finalScore < 0 { finalScore = 0 }

	// Generate Insight
	if finalScore >= 90 {
		if res.Role {
			res.EngagementInsight = "Exceptional likelihood of reply. This is an established organizational touchpoint."
		} else {
			res.EngagementInsight = "Exceptional likelihood of reply. This is a primary, active identity."
		}
	} else if finalScore >= 70 {
		if res.Role {
			res.EngagementInsight = "High likelihood of reply. Active organizational channel."
		} else {
			res.EngagementInsight = "High likelihood of reply from an established infrastructure."
		}
	} else if finalScore >= 40 {
		res.EngagementInsight = "Moderate likelihood. Common for enterprise or catch-all addresses."
	} else {
		res.EngagementInsight = "Low likelihood of reply due to infrastructure or identity age signals."
	}

	res.EngagementProbability = finalScore
	return finalScore
}
