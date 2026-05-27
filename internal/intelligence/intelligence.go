package intelligence

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// RDAPResponse represents the structure of an RDAP domain response.
type RDAPResponse struct {
	Events []struct {
		EventAction string `json:"eventAction"`
		EventDate   string `json:"eventDate"`
	} `json:"events"`
}

// GetDomainAge returns the approximate age of a domain in years using RDAP.
func GetDomainAge(domain string) (float64, error) {
	url := fmt.Sprintf("https://rdap.org/domain/%s", domain)
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return -1, fmt.Errorf("RDAP request failed with status: %d", resp.StatusCode)
	}

	var rdap RDAPResponse
	if err := json.NewDecoder(resp.Body).Decode(&rdap); err != nil {
		return -1, err
	}

	for _, event := range rdap.Events {
		if event.EventAction == "registration" {
			regTime, err := time.Parse(time.RFC3339, event.EventDate)
			if err != nil {
				continue
			}
			age := time.Since(regTime).Hours() / 24 / 365
			return age, nil
		}
	}
	return -1, fmt.Errorf("registration date not found in RDAP response")
}

// GetBreachAge simulates an OSINT check against historical data breaches.
// Returns (Age in years, Confidence 0-100)
func GetBreachAge(email string) (int, int) {
	if strings.Contains(email, "2015") || strings.Contains(email, "legacy") {
		return 9, 90
	}
	if strings.Contains(email, "2020") {
		return 4, 85
	}
	return 0, 0
}

// GetCombinedAge triangulates multiple signals to find the "True Identity Age".
func GetCombinedAge(email string, domainAge float64, telemetryAge float64) (float64, int, string) {
	breachAge, breachConf := GetBreachAge(email)
	
	// Start with the baseline: Infrastructure Age
	finalAge := domainAge
	conf := 60 // Baseline confidence for RDAP
	status := "Infrastructure Only"

	// 1. Conflict Resolution (Identity > Domain: High Risk)
	// If breach data or telemetry indicates the identity is older than the domain by > 1 year
	if (float64(breachAge) > domainAge + 1 && breachConf > 0) || (telemetryAge > domainAge + 1) {
		maxIdentityAge := float64(breachAge)
		if telemetryAge > maxIdentityAge { maxIdentityAge = telemetryAge }
		return maxIdentityAge, 95, "Conflict: Identity older than Domain (High Risk)"
	}

	// 2. Telemetry Integration
	if telemetryAge > 0 && telemetryAge > finalAge {
		finalAge = telemetryAge
		conf = 80 // Internal data is high confidence
		status = "Telemetry Verified"
	}

	// 3. Breach Data Integration
	if breachConf > 0 {
		if float64(breachAge) > finalAge {
			finalAge = float64(breachAge)
			conf = breachConf
			status = "Identity Verified (OSINT)"
		}
	}

	// 4. Legacy Trust Boost
	// "If both breach data and domain age indicate 5+ years, we assign a 'legacy trust' boost."
	if float64(breachAge) >= 5 && domainAge >= 5 {
		status = "Legacy Trust Boost"
		conf = 100
	}

	return finalAge, conf, status
}
