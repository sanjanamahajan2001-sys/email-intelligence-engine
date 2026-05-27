package intelligence

import (
	"strings"
	"sync"
	"time"

	"github.com/sanjana/email-validator/internal/db"
)


var (
	disposableMap map[string]string
	mu            sync.RWMutex
	initialized   bool
)


// InitDisposable initializes the in-memory cache and returns the last sync time.
func InitDisposable(database *db.DB) (time.Time, error) {
	mu.Lock()
	defer mu.Unlock()

	// 1. Ensure MX seeds are bootstrapped
	_ = BootstrapEnterpriseSeeds(database)

	var err error
	disposableMap, err = database.GetAllDisposableDomains()
	if err != nil {
		return time.Time{}, err
	}

	initialized = true


	// Check last discovery time
	lastSyncStr, _ := database.GetMetadata("last_disposable_discovery")
	if lastSyncStr != "" {
		lastSync, err := time.Parse(time.RFC3339, lastSyncStr)
		if err == nil {
			return lastSync, nil
		}
	}

	return time.Time{}, nil
}

// IsDisposable checks if a domain or its parent domains are in the disposable list.
// Returns (isDisposable, matchType, providerName)
func IsDisposable(dbConn *db.DB, domain string) (bool, string, string) {
	mu.RLock()
	defer mu.RUnlock()

	if !initialized {
		return false, "", ""
	}

	domain = strings.ToLower(domain)
	
	// Check the full domain first (in-memory)
	if provider, ok := disposableMap[domain]; ok {
		return true, "Static Match", provider
	}

	// Cross-Process Fallback: Check the database directly if not in our process memory
	if dbConn != nil {
		if exists, provider, _ := dbConn.IsDisposable(domain); exists {
			// Update in-memory map for O(1) next time
			mu.Lock()
			if disposableMap == nil { disposableMap = make(map[string]string) }
			disposableMap[domain] = provider
			mu.Unlock()
			return true, "Static Match (DB Learned)", provider
		}
	}



	// Recursive parent-domain check (e.g., test.sub.mailinator.com -> sub.mailinator.com -> mailinator.com)
	parts := strings.Split(domain, ".")
	for len(parts) > 2 { // We stop at length 2 (e.g., example.com)
		parts = parts[1:]
		parent := strings.Join(parts, ".")
		if provider, ok := disposableMap[parent]; ok {
			return true, "Subdomain Match: " + parent, provider
		}
		
		// DB Fallback for subdomains
		if dbConn != nil {
			if exists, provider, _ := dbConn.IsDisposable(parent); exists {
				return true, "Subdomain Match (DB): " + parent, provider
			}
		}
	}



	// 3. Zero-Day Name Heuristics
	suspiciousKeywords := []string{"burner", "tempmail", "throwaway", "disposable", "10minutemail", "fakeinbox", "guerrillamail"}
	for _, kw := range suspiciousKeywords {
		if strings.Contains(domain, kw) {
			return true, "Zero-Day Heuristic: Suspicious Domain Pattern", ""
		}
	}

	return false, "", ""
}




// RegisterNewDisposable adds a newly discovered domain to the cache and database.
func RegisterNewDisposable(database *db.DB, domain string, providerName string) {
	domain = strings.ToLower(strings.TrimSpace(domain))
	if domain == "" {
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if disposableMap == nil {
		disposableMap = make(map[string]string)
	}

	if p, exists := disposableMap[domain]; !exists || (p == "" && providerName != "") {
		disposableMap[domain] = providerName
		_ = database.SaveDisposableDomains(map[string]string{domain: providerName})
	}
}



// SyncDisposable replaces the old GitHub sync with the new Discovery Pump.
func SyncDisposable(database *db.DB) (int, error) {
	count, err := RunDiscoveryPump(database)
	if err != nil {
		return 0, err
	}

	// Update metadata
	_ = database.SetMetadata("last_disposable_discovery", time.Now().UTC().Format(time.RFC3339))

	// Refresh cache
	mu.Lock()
	disposableMap, _ = database.GetAllDisposableDomains()
	mu.Unlock()



	return count, nil
}