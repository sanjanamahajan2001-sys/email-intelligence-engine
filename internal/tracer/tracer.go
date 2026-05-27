package tracer

import (
	"fmt"
	"net"
	"strings"
	"time"
)

// Log level constants
const (
	LevelInfo  = "INFO"
	LevelError = "ERROR"
	LevelTrace = "TRACE"
)

// Log formats and prints a message with a timestamp and level.
func Log(level, component, message string) {
	timestamp := time.Now().UTC().Format("2006-01-02 15:04:05")
	fmt.Printf("[%s] [%s] %s: %s\n", timestamp, level, component, message)
}

// Info logs an informational message.
func Info(component, message string) {
	Log(LevelInfo, component, message)
}

// Error logs an error message.
func Error(component, message string, err error) {
	msg := message
	if err != nil {
		msg = fmt.Sprintf("%s: %v", message, err)
	}
	Log(LevelError, component, msg)
}

// TraceResult holds information about the email's origin and reputation.
type TraceResult struct {
	Domain     string   `json:"domain"`
	IPs        []string `json:"ips"`
	Provider   string   `json:"provider"`
	Reputation string   `json:"reputation"`
	Risky      bool     `json:"risky"`
}

// TraceEmail analyzes the domain of the email to identify its source.
func TraceEmail(email string) (*TraceResult, error) {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid email")
	}
	domain := parts[1]

	ips, _ := net.LookupHost(domain)
	
	result := &TraceResult{
		Domain:     domain,
		IPs:        ips,
		Reputation: "Neutral", // Default
		Risky:      false,
	}

	// Basic provider identification
	if strings.Contains(domain, "google.com") || strings.Contains(domain, "gmail.com") {
		result.Provider = "Google"
	} else if strings.Contains(domain, "outlook.com") || strings.Contains(domain, "hotmail.com") {
		result.Provider = "Microsoft"
	} else if strings.Contains(domain, "yahoo.com") {
		result.Provider = "Yahoo"
	} else {
		result.Provider = "Custom/Other"
	}

	return result, nil
}
