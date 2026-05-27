package intelligence

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"


	"github.com/sanjana/email-validator/internal/db"
)

var (
	// High-confidence keywords for active probing
	disposableKeywords = []string{
		"temporary email",
		"disposable email",
		"10 minute mail",
		"burner email",
		"throwaway email",
		"fake email",
		"anonymous email inbox",
		"forget-me-not email",
		"emailondeck",
	}

	// Shared HTTP Client with optimized transport for parallel probing
	sharedClient = &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     30 * time.Second,
		},
	}
)


var (
	// Safety Whitelist: Known legitimate providers to prevent Hub false positives
	SafetyWhitelist = map[string]bool{
		"google.com":         true,
		"googlemail.com":    true,
		"outlook.com":       true,
		"hotmail.com":       true,
		"microsoft.com":     true,
		"icloud.com":        true,
		"apple.com":         true,
		"yahoo.com":         true,
		"protonmail.com":    true,
		"proton.me":         true,
		"mail.com":          true,
		"gmx.com":           true,
		"fastmail.com":      true,
		"zoho.com":          true,
	}
)

// BootstrapEnterpriseSeeds populates the database with initial high-confidence MX patterns.
func BootstrapEnterpriseSeeds(database *db.DB) error {
	seeds := map[string]string{
		"mail.mailinator.com":      "Mailinator Hub",
		"mail2.mailinator.com":     "Mailinator Hub",
		"mx.emailondeck.com":       "EmailOnDeck",
		"mail.guerrillamail.com":   "GuerrillaMail Hub",
		"mx.10minutemail.com":      "10MinuteMail",
		"mail.mintemail.com":       "MintEmail",
		"mx.mailapi.org":           "Mail7",
		"mx.mail7.io":              "Mail7",
		"mx.yopmail.com":           "Yopmail",
		"gate.poczta.onet.pl":      "Onet (Commonly used by PL-based disposables)",
		"mx.disposable.com":        "Generic Disposable Hub",
		"mail.grr.la":              "GuerrillaMail Hub",
		"mx.maildrop.cc":           "MailDrop Hub",
		"mx.dispostable.com":       "Dispostable Hub",
	}

	return database.SaveMXSignatures(seeds)
}



func IsDisposableByMX(database *db.DB, mxHosts []string) (bool, string) {
	signatures, err := database.GetAllMXSignatures()
	if err != nil {
		return false, ""
	}

	for _, host := range mxHosts {
		host = strings.TrimSuffix(strings.ToLower(host), ".")
		
		// 1. Safety Check: If MX is a known legitimate provider, skip fingerprinting
		for trusted := range SafetyWhitelist {
			if host == trusted || strings.HasSuffix(host, "."+trusted) {
				return false, ""
			}
		}

		// 2. Hub Check: Match against known disposable signatures
		for sig, provider := range signatures {
			if host == sig || strings.HasSuffix(host, "."+sig) {
				return true, provider
			}
		}
	}

	return false, ""
}


// ActiveProbe performs a safe HTTP/S check on a domain to identify a disposable provider.
func ActiveProbe(domain string) (bool, string) {


	urls := []string{
		"https://" + domain,
		"http://" + domain,
	}

	for _, url := range urls {
		resp, err := sharedClient.Get(url)
		if err != nil {
			continue
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			continue
		}

		content := strings.ToLower(string(body))
		for _, kw := range disposableKeywords {
			if strings.Contains(content, kw) {
				return true, fmt.Sprintf("Confirmed via Active Probe: Found keyword '%s'", kw)
			}
		}
	}

	return false, ""
}

// RunDiscoveryPump scans for new domains in the history and verifies their infrastructure.
// It also proactively checks a few known high-rotation seed sites to 'prime' the system.
// RunDiscoveryPump optimized with a parallel worker pool to prevent hanging.
func RunDiscoveryPump(database *db.DB) (int, error) {
	newFindings := 0
	type probeResult struct {
		domain string
		isDisp bool
		reason string
	}

	domainsToProbe := make(chan string, 100)
	results := make(chan probeResult, 100)
	var wg sync.WaitGroup

	// 1. Start Worker Pool (Concurrency: 10)
	for w := 1; w <= 10; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for domain := range domainsToProbe {
				isDisp, reason := ActiveProbe(domain)
				results <- probeResult{domain, isDisp, reason}
			}
		}()
	}

	// 2. Collector: Handle findings in the background
	done := make(chan bool)
	go func() {
		for res := range results {
			if res.isDisp {
				_ = database.SaveDisposableDomains(map[string]string{res.domain: res.reason})
				newFindings++
			}
		}
		done <- true
	}()

	// 3. Feeder: Load domains from seeds and telemetry
	processedCount := 0
	maxProbes := 50 // Safe cap for performance

	// Source A: Proactive Scrape
	seedURLs := []string{
		"https://raw.githubusercontent.com/disposable-email-domains/disposable-email-domains/master/allowlist.conf",
	}

	client := &http.Client{Timeout: 5 * time.Second}
	for _, url := range seedURLs {
		if processedCount >= maxProbes { break }
		resp, err := client.Get(url)
		if err != nil { continue }
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		lines := strings.Split(string(body), "\n")
		for _, line := range lines {
			domain := strings.TrimSpace(line)
			if domain == "" || strings.HasPrefix(domain, "#") { continue }
			domainsToProbe <- domain
			processedCount++
			if processedCount >= maxProbes { break }
		}
	}

	// Source B: Telemetry Scrape
	if processedCount < maxProbes {
		telemetry, err := database.GetRecentDomains(maxProbes - processedCount)
		if err == nil {
			for _, domain := range telemetry {
				domainsToProbe <- domain
				processedCount++
			}
		}
	}

	close(domainsToProbe)
	wg.Wait()
	close(results)
	<-done

	return newFindings, nil
}


