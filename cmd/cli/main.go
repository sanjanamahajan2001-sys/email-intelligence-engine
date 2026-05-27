package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/sanjana/email-validator/internal/config"
	"github.com/sanjana/email-validator/internal/core"
	"github.com/sanjana/email-validator/internal/db"
	"github.com/sanjana/email-validator/internal/intelligence"
	"github.com/sanjana/email-validator/internal/service"
	"github.com/sanjana/email-validator/internal/tracer"
	"github.com/sanjana/email-validator/internal/tui"
)

var (
	database *db.DB
	appConfig *config.AppConfig
)

func main() {
	if len(os.Args) < 2 {
		showUsage()
		os.Exit(1)
	}

	appConfig = config.LoadConfig()

	var err error
	database, err = db.InitDB(appConfig.DBPath)
	if err != nil {
		tracer.Error("CLI", "Error initializing database", err)
		os.Exit(1)
	}

	// Initialize Intelligence Modules & Check for Auto-Sync
	lastSync, err := intelligence.InitDisposable(database)
	if err != nil {
		tracer.Error("CLI", "Warning: Failed to load disposable domains", err)
	}

	// Automated Background Sync (if older than 24 hours)
	if lastSync.IsZero() || time.Since(lastSync) > 24*time.Hour {
		go func() {
			_, _ = intelligence.SyncDisposable(database)
		}()
	}

	command := os.Args[1]
	switch command {
	case "check":
		if len(os.Args) < 3 {
			fmt.Println("Usage: email-validator check <email> [--force]")
			os.Exit(1)
		}
		force := false
		if len(os.Args) > 3 && (os.Args[3] == "--force" || os.Args[3] == "--refresh") {
			force = true
		}
		handleCheck(os.Args[2], force)
	case "audit":
		if len(os.Args) < 3 {
			fmt.Println("Usage: email-validator audit <email>")
			os.Exit(1)
		}
		handleAudit(os.Args[2])
	case "history":
		limit := 10
		if len(os.Args) > 2 {
			if l, err := strconv.Atoi(os.Args[2]); err == nil {
				limit = l
			}
		}
		results, err := database.GetHistory(limit)
		if err != nil {
			fmt.Printf("Error retrieving history: %v\n", err)
			os.Exit(1)
		}
		renderHistory(results)
	case "export":
		if len(os.Args) < 3 {
			fmt.Println("Usage: email-validator export <filename.csv|json>")
			os.Exit(1)
		}
		handleExport(os.Args[2])
	case "sync":
		handleSync()
	case "interactive":
		if err := tui.Start(database, appConfig); err != nil {
			fmt.Printf("Error starting TUI: %v\n", err)
		}
	case "user":
		if len(os.Args) < 3 {
			fmt.Println("Usage: email-validator user <create|list> ...")
			os.Exit(1)
		}
		subcommand := os.Args[2]
		switch subcommand {
		case "create":
			if len(os.Args) < 5 {
				fmt.Println("Usage: email-validator user create <username> <password>")
				os.Exit(1)
			}
			handleUserCreate(os.Args[3], os.Args[4])
		case "delete":
			if len(os.Args) < 4 {
				fmt.Println("Usage: email-validator user delete <username>")
				os.Exit(1)
			}
			handleUserDelete(os.Args[3])
		case "password-reset":
			if len(os.Args) < 5 {
				fmt.Println("Usage: email-validator user password-reset <username> <new_password>")
				os.Exit(1)
			}
			handleUserPasswordReset(os.Args[3], os.Args[4])
		case "list":
			handleUserList()
		default:
			fmt.Println("Unknown user subcommand:", subcommand)
		}
	default:
		showUsage()
	}
}

func showUsage() {
	fmt.Println("Usage:")
	fmt.Println("  email-validator check <email> [--force] Single email check (use --force to bypass cache)")
	fmt.Println("  email-validator audit <email>          View detailed transition/lifecycle log")
	fmt.Println("  email-validator interactive           Start TUI mode")
	fmt.Println("  email-validator history [count]       View recent scans")
	fmt.Println("  email-validator export <filename>     Export history (csv/json)")
	fmt.Println("  email-validator sync                  Sync disposable domain lists")
	fmt.Println("  email-validator user create <u/n> <p> Create a new API user")
	fmt.Println("  email-validator user delete <u/n>    Permanently remove an API user")
	fmt.Println("  email-validator user password-reset <u/n> <p> Reset an API user's password")
	fmt.Println("  email-validator user list             List all API users")
}

func handleCheck(email string, force bool) {
	res, err := service.ProcessEmail(database, appConfig.SMTPSender, email, "CLI", "127.0.0.1", 0, "Validator-CLI", force)
	if err != nil {
		fmt.Printf("Error during validation: %v\n", err)
		return
	}
	renderReport(res, email)
}

func handleAudit(email string) {
	history, err := database.GetScanHistory(email)
	if err != nil {
		fmt.Printf("Error retrieving transitions: %v\n", err)
		return
	}
	if len(history) == 0 {
		fmt.Println("No historical transitions found for this address.")
		return
	}
	renderAuditTrail(history, email)
}

func handleExport(filename string) {
	results, err := database.GetHistory(1000)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	f, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer f.Close()

	if strings.HasSuffix(filename, ".json") {
		json.NewEncoder(f).Encode(results)
	} else {
		w := csv.NewWriter(f)
		// Master Forensic Header
		w.Write([]string{
			"Email", "Is_Valid", "Syntax", "DNS", "SMTP", "Reputation", "Risk", "Provider", 
			"Disposable", "Role", "Catch_All", "Domain_Age", "Source", "Requester_IP", "User_ID", 
			"User_Agent", "Latency_MS", "Date_UTC", "Engagement_Factors",
		})
		
		for _, r := range results {
			// Convert bools to 1/0 for easier analysis
			w.Write([]string{
				r.Email,
				fmt.Sprintf("%v", r.IsValid),
				fmt.Sprintf("%v", r.Syntax),
				fmt.Sprintf("%v", r.DNS),
				fmt.Sprintf("%v", r.SMTP),
				strconv.Itoa(r.ReputationScore),
				r.RiskLevel,
				r.Provider,
				fmt.Sprintf("%v", r.Disposable),
				fmt.Sprintf("%v", r.Role),
				fmt.Sprintf("%v", r.CatchAll),
				fmt.Sprintf("%.2f", r.DomainAgeYears),
				r.Source,
				r.ClientIP,
				strconv.Itoa(r.UserID),
				r.UserAgent,
				fmt.Sprintf("%d", r.ProcessingTimeMs),
				r.CreatedAt,
				strings.Join(r.EngagementFactors, "; "),
			})
		}
		w.Flush()
	}
	fmt.Printf("✨ Exported %d records to %s\n", len(results), filename)
}

func handleSync() {
	fmt.Println("🔄 Starting In-House Infrastructure Discovery...")
	fmt.Println("   - Bootstrapping MX Signatures")
	fmt.Println("   - Running Discovery Pump on validation telemetry")
	
	count, err := intelligence.SyncDisposable(database)
	if err != nil {
		fmt.Printf("❌ Discovery failed: %v\n", err)
		return
	}
	fmt.Printf("✅ Discovery complete! Automatically identified %d new domains.\n", count)
}


func getStatusIcon(ok bool) string {
	if ok {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Render("✅")
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Render("❌")
}

func renderHistory(results []*core.EmailResult) {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7D56F4")).Padding(1, 0)
	fmt.Println(titleStyle.Render("🕒 RECENT VERIFICATION HISTORY"))
	
	header := fmt.Sprintf("%-4s | %-12s | %-15s | %-30s | %-10s | %-5s", 
		"ID", "TIMESTAMP", "REQUESTER_IP", "EMAIL", "STATUS", "SCORE")
	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#777777")).Render(header))
	fmt.Println(strings.Repeat("-", 95))

	for i, res := range results {
		status := "✅ VALID"
		if !res.IsValid { status = "❌ INVALID" }
		ts := formatRelativeTime(res.CreatedAt)
		fmt.Printf("%-4d | %-12s | %-15s | %-30s | %-10s | %d/100\n", 
			i+1, ts, res.ClientIP, res.Email, status, res.ReputationScore)
	}
}

func renderReport(res *core.EmailResult, email string) {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7D56F4")).Padding(0, 1)
	headerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#555555")).Italic(true)
	
	report := titleStyle.Render("🛡️  VERIFICATION REPORT: " + res.Email) + "\n"
	if res.IsCached {
		cacheStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4")).Bold(true).Padding(0, 1).MarginLeft(1)
		report += cacheStyle.Render("⚡ CACHED RESULT") + "\n"
	}
	report += headerStyle.Render("   Analysis timestamp: " + time.Now().Format("2006-01-02 15:04:05")) + "\n\n"

	col1Label := lipgloss.NewStyle().Width(20).Foreground(lipgloss.Color("#777777"))
	col1Value := lipgloss.NewStyle().Width(32).Bold(true)
	col2Label := lipgloss.NewStyle().Width(20).Foreground(lipgloss.Color("#777777"))
	col2Value := lipgloss.NewStyle().Width(45).Bold(true)

	renderRow := func(l1, v1, l2, v2 string) string {
		return "  " + lipgloss.JoinHorizontal(lipgloss.Top,
			col1Label.Render(l1),
			col1Value.Render(v1),
			col2Label.Render(l2),
			col2Value.Render(v2),
		) + "\n"
	}

	smtpValue := getStatusIcon(res.SMTP)
	if res.Greylisted {
		smtpValue = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00")).Render("⏳ Greylisted")
	} else if res.SMTPBlocked && !res.SMTP {
		smtpValue = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFA500")).Render("🛡️  Blocked (P25)")
	}

	aliasStr := "No"
	if res.HasAlias { aliasStr = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00")).Render("Yes (+alias)") }

	roleStr := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Render("No")
	if res.Role { roleStr = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Render("Yes (Role)") }

	dispStr := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Render("No")
	if res.Disposable { dispStr = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Bold(true).Render("⚠️  YES") }

	ageStr := fmt.Sprintf("%.1f yrs", res.DomainAgeYears)
	if res.DomainAgeYears < 0 { ageStr = "New/Unknown" }

	report += renderRow("📂 [Identity/Alias]", aliasStr, "📡 [DNS/MX Records]", getStatusIcon(res.DNS))
	report += renderRow("🏢 [Base Address]", res.BaseEmail, "📧 [SMTP Delivery]", smtpValue)
	report += renderRow("✨ [RFC Syntax]", getStatusIcon(res.Syntax), "⏳ [Domain Age]", ageStr)
	report += renderRow("👤 [User Role]", roleStr, "🏛️ [TLD Reputation]", res.TldTrust)
	
	catchAllIcon := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Render("No")
	if res.CatchAll { catchAllIcon = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00")).Render("⚠️  YES") }
	report += renderRow("🗑️ [Disposable]", dispStr, "🕳️ [Catch-All Hub]", catchAllIcon)

	engagementStr := fmt.Sprintf("%d%%", res.EngagementProbability)
	engagementColor := "#00FF00"
	if res.EngagementProbability < 40 { engagementColor = "#FF0000" } else if res.EngagementProbability < 75 { engagementColor = "#FFFF00" }
	engagementValue := lipgloss.NewStyle().Foreground(lipgloss.Color(engagementColor)).Bold(true).Render(engagementStr)

	shortResp := core.GetShortSMTPResponse(res.LastSMTPResponse)
	report += renderRow("🤝 [Engagement]", engagementValue, "📡 [Last Response]", shortResp)

	// Lifecycle Reporting
	stateColor := "#00FF00"
	if res.LifecycleState == "STALE" { stateColor = "#FFFF00" } else if res.LifecycleState == "ABANDONED" || res.LifecycleState == "INVALID" { stateColor = "#FF0000" }
	stateValue := lipgloss.NewStyle().Foreground(lipgloss.Color(stateColor)).Bold(true).Render(res.LifecycleState)
	
	validatedStr := res.LastVerifiedAt
	if validatedStr == "" { validatedStr = "Unknown" }
	report += renderRow("🔄 [Lifecycle State]", stateValue, "📅 [Last Validated]", validatedStr)

	if res.Message != "" {
		report += "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Bold(true).PaddingLeft(2).Render("⚠️  "+res.Message) + "\n"
	}

	// 🔍 Detailed Telemetry Insight for Failures
	if !res.SMTP && res.LastSMTPResponse != "" && !strings.Contains(res.LastSMTPResponse, "250") {
		telemetryStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#777777")).
			Italic(true).
			PaddingLeft(2).
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(lipgloss.Color("#444444")).
			Padding(0, 1)
		
		report += "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("#555555")).Bold(true).PaddingLeft(2).Render("📡 FULL TELEMETRY") + "\n"
		// Wrap telemetry to maintain box integrity
		faintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#777777")).Italic(true)
		wrappedTelemetry := faintStyle.Width(80).Render(res.LastSMTPResponse)
		report += telemetryStyle.Render(wrappedTelemetry) + "\n"
	}

	scoreColor := "#00FF00"
	if res.ReputationScore < 50 { scoreColor = "#FF0000" } else if res.ReputationScore < 85 { scoreColor = "#FFFF00" }
	
	scoreBox := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color(scoreColor)).
		Padding(1, 4).MarginTop(1).
		Render(fmt.Sprintf("REPUTATION SCORE: %d/100\nRISK LEVEL:       %s\nPROVIDER:         %s",
			res.ReputationScore, res.RiskLevel, res.Provider))
	
	report += "\n" + scoreBox + "\n"
	
	// 🧠 Engagement Analysis & Factors
	if len(res.EngagementFactors) > 0 {
		analysisStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true).
			PaddingTop(1).
			PaddingLeft(2)
		
		factorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#AAAAAA")).
			PaddingLeft(4)

		report += analysisStyle.Render("🧠 ENGAGEMENT ANALYSIS FACTORS:") + "\n"
		for _, factor := range res.EngagementFactors {
			report += factorStyle.Render("• " + factor) + "\n"
		}
	}

	fmt.Println("\n" + lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1, 2).Render(report))
	fmt.Println("\n✨ Validation Complete!")
	
	insightStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF00")).
		Bold(true).
		Padding(0, 1).
		MarginLeft(1)

	fmt.Printf("\n💡 Engagement Insight: %s\n", insightStyle.Render(res.EngagementInsight))
}

func formatRelativeTime(createdAt string) string {
	t, err := time.Parse("2006-01-02 15:04:05", createdAt)
	if err != nil {
		t, _ = time.Parse(time.RFC3339, createdAt)
	}
	diff := time.Now().UTC().Sub(t)
	if diff < time.Minute { return "Just now" }
	if diff < time.Hour { return fmt.Sprintf("%dm ago", int(diff.Minutes())) }
	return t.Format("Jan 02")
}

func renderAuditTrail(history []*core.EmailResult, email string) {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7D56F4")).Padding(1, 0)
	fmt.Println(titleStyle.Render("📜 LIFECYCLE TRANSITIONS LOG: " + email))
	
	header := fmt.Sprintf("%-20s | %-12s | %-15s | %-5s | %-15s", 
		"TIMESTAMP", "STATUS", "CLIENT_IP", "UID", "TRANSITION")
	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#777777")).Render(header))
	fmt.Println(strings.Repeat("-", 85))

	for i, res := range history {
		state := res.LifecycleState
		if state == "" {
			state = "ACTIVE"
			if !res.IsValid { state = "INVALID" }
		}
		
		// Determine transition label
		transition := "Initial Scan"
		if i > 0 {
			prev := history[i-1]
			prevState := prev.LifecycleState
			if prevState == "" {
				prevState = "ACTIVE"
				if !prev.IsValid { prevState = "INVALID" }
			}

			if prevState != state {
				transition = fmt.Sprintf("➡️  %s", state)
			} else {
				transition = "No State Change"
			}
		}

		fmt.Printf("%-20s | %-12s | %-15s | %-5d | %-15s\n", 
			res.CreatedAt, state, res.ClientIP, res.UserID, transition)
	}
	fmt.Println("\n✨ End of Audit Trail")
}

func handleUserCreate(username, password string) {
	hash, err := service.HashPassword(password)
	if err != nil {
		fmt.Printf("Error hashing password: %v\n", err)
		return
	}
	err = database.CreateUser(username, hash)
	if err != nil {
		fmt.Printf("Error creating user: %v\n", err)
		return
	}
	fmt.Printf("✅ User '%s' created successfully!\n", username)
}

func handleUserList() {
	users, err := database.GetUsers()
	if err != nil {
		fmt.Printf("❌ Error retrieving user list: %v\n", err)
		return
	}

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7D56F4")).Padding(1, 0)
	fmt.Println(titleStyle.Render("👥 REGISTERED API USERS"))

	header := fmt.Sprintf("%-4s | %-20s | %-20s", "ID", "USERNAME", "CREATED_AT (UTC)")
	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#777777")).Render(header))
	fmt.Println(strings.Repeat("-", 50))

	for _, u := range users {
		fmt.Printf("%-4v | %-20s | %-20s\n", u["id"], u["username"], u["created_at"])
	}
	fmt.Printf("\nTotal Users: %d\n", len(users))
}

func handleUserDelete(username string) {
	err := database.DeleteUser(username)
	if err != nil {
		fmt.Printf("❌ Error deleting user '%s': %v\n", username, err)
		return
	}
	fmt.Printf("✅ User '%s' has been permanently removed.\n", username)
}

func handleUserPasswordReset(username, password string) {
	hash, err := service.HashPassword(password)
	if err != nil {
		fmt.Printf("❌ Error hashing new password: %v\n", err)
		return
	}
	err = database.UpdateUserPassword(username, hash)
	if err != nil {
		fmt.Printf("❌ Error resetting password for '%s': %v\n", username, err)
		return
	}
	fmt.Printf("✅ Password for user '%s' has been reset successfully.\n", username)
}
