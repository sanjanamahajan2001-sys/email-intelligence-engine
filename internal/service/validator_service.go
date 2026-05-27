package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/sanjana/email-validator/internal/core"
	"github.com/sanjana/email-validator/internal/db"
	"github.com/sanjana/email-validator/internal/intelligence"
	"github.com/sanjana/email-validator/internal/tracer"
)

// ProcessEmail handles the full validation lifecycle and returns a unified result.
func ProcessEmail(database *db.DB, fromEmail, email, source, clientIP string, userID int, userAgent string, force bool) (*core.EmailResult, error) {
	start := time.Now()
	// 1. Initial Syntax Validation (Always)
	res := core.NewResult(email)
	res.ClientIP = clientIP
	res.UserID = userID
	res.UserAgent = userAgent
	res.Syntax = core.ValidateSyntax(email)
	if !res.Syntax {
		res.IsValid = false
		res.Message = "Invalid Syntax"
		res.LifecycleState = "INVALID"
		res.ProcessingTimeMs = time.Since(start).Milliseconds()
		return res, nil
	}

	// 2. Intelligent Lifecycle Strategy (Lookup-Before-Scan)
	var previousResult *core.EmailResult
	if database != nil {
		var err error
		previousResult, err = database.GetLatestResult(email)
		if err != nil {
			tracer.Error("Service", "Failed to lookup cache", err)
		}

		if previousResult != nil && !force {
			// TTL Strategy: Standard (30 days), Catch-All (14 days)
			ttlDays := 30
			if previousResult.CatchAll {
				ttlDays = 14
			}
			
			createdAt, err := time.Parse("2006-01-02 15:04:05", previousResult.CreatedAt)
			if err != nil {
				createdAt, _ = time.Parse(time.RFC3339, previousResult.CreatedAt)
			}
			
			if time.Since(createdAt) < time.Duration(ttlDays)*24*time.Hour {
				previousResult.IsCached = true
				previousResult.LifecycleState = "INVALID"
				if previousResult.IsValid {
					previousResult.LifecycleState = "ACTIVE"
				}
				if previousResult.CatchAll {
					previousResult.LifecycleState = "CATCH-ALL"
				}
				return previousResult, nil
			}
			// Result is Stale
			res.STALE = true
		}
	}

	// 3. DNS/MX/Disposable Validation
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return res, fmt.Errorf("invalid email format")
	}
	domain := parts[1]

	isDisp, dispReason, provider := intelligence.IsDisposable(database, domain)
	res.Disposable = isDisp
	res.Provider = provider
	res.Role = core.CheckRoleAccount(email)

	if res.Disposable {
		res.Message = "Identified via " + dispReason
		intelligence.RegisterNewDisposable(database, domain, provider)
	}

	ok, hosts, err := core.ValidateMX(email)
	res.DNS = ok
	if !ok || err != nil {
		res.DNS = false
	} else {
		isDispMX, providerMX := intelligence.IsDisposableByMX(database, hosts)
		if isDispMX {
			res.Disposable = true
			res.Provider = providerMX
			res.Message = "Flagged via Infrastructure: " + providerMX
			intelligence.RegisterNewDisposable(database, domain, providerMX)
		}

		// SMTP Validation & Catch-All
		smtpValid, greylisted, smtpResp, smtpErr := core.ValidateSMTP(fromEmail, email, hosts)
		res.SMTP = smtpValid
		res.Greylisted = greylisted
		res.LastSMTPResponse = smtpResp
		
		if smtpErr != nil {
			errStr := strings.ToLower(smtpErr.Error())
			if strings.Contains(errStr, "timeout") || 
			   strings.Contains(errStr, "refused") || 
			   strings.Contains(errStr, "network is unreachable") || 
			   strings.Contains(errStr, "i/o timeout") || 
			   strings.Contains(errStr, "permission denied") {
				res.SMTPBlocked = true
			}
		}

		res.CatchAll = core.DetectCatchAll(fromEmail, domain, hosts, res.SMTPBlocked, res.LastSMTPResponse)
	}

	// 5. Multi-Signal Intelligence (Domain Age & Trust)
	domAge, _ := intelligence.GetDomainAge(domain)
	res.DomainAgeYears = domAge
	
	trust := intelligence.GetDomainTrust(domain, domAge)
	if trust.IsLowTrust {
		res.TldTrust = "Low"
		if !res.Disposable {
			res.Message = "Suspicious Domain: " + trust.Reason
		}
	} else {
		res.TldTrust = "High"
	}

	// Telemetry Age (System First-Seen)
	teleAge := 0.0
	if database != nil {
		teleAge, _ = database.GetFirstSeenAge(email)
	}

	idAge, idConf, ageStatus := intelligence.GetCombinedAge(email, domAge, teleAge)
	res.IdentityAgeYears = idAge
	res.ConfidenceScore = idConf

	if res.Message == "" {
		if strings.Contains(ageStatus, "Conflict") {
			res.Message = ageStatus
		} else if ageStatus == "Legacy Trust Boost" {
			res.Message = "Trusted: Established Legacy Identity"
		}
	}

	// 6. Overall Scoring & Lifecycle Mapping
	core.CalculateScore(res, ageStatus)
	res.EngagementProbability = intelligence.CalculateEngagementProbability(res, database)
	res.IsValid = res.Syntax && res.DNS && (res.SMTP || res.SMTPBlocked || res.Greylisted) && !res.Disposable

	// Determine Transition State (Advanced State Machine)
	res.LifecycleState = "INVALID"
	if res.IsValid {
		res.LifecycleState = "ACTIVE"
	}
	if res.CatchAll {
		res.LifecycleState = "CATCH-ALL"
	}
	
	if previousResult != nil {
		if previousResult.IsValid && !res.IsValid {
			// Was valid before, now it's not (and it's not just a transient block)
			if !res.SMTP && !res.Greylisted && !res.SMTPBlocked {
				res.LifecycleState = "ABANDONED"
				res.Message = "Identity Abandoned: Previously Valid address now persistent failure"
			}
		} else if !previousResult.IsValid && res.IsValid {
			// Identity was previously invalid/bounced, but now it's reachable again
			res.LifecycleState = "ACTIVE"
			res.Message = "Identity Recovered: Previously unreachable address is now ACTIVE"
		}

		// STALE Logic: If this scan was triggered by a stale cache
		if res.STALE {
			if res.Message == "" {
				res.Message = "Re-verified: Outdated record (STALE) successfully updated"
			}
		}
	}

	res.ProcessingTimeMs = time.Since(start).Milliseconds()

	// 7. Persist Result
	if database != nil {
		res.Source = source
		res.CreatedAt = time.Now().UTC().Format("2006-01-02 15:04:05")
		res.LastVerifiedAt = res.CreatedAt
		err := database.SaveScan(res, source)
		if err != nil {
			tracer.Error("Service", "Failed to save scan result", err)
		}

		if ok && !res.Disposable {
			err = database.AddToDiscoveryQueue(domain, hosts)
			if err != nil {
				tracer.Error("Service", "Failed to update discovery queue", err)
			}
		}
	}

	return res, nil
}

