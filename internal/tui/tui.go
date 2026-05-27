package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sanjana/email-validator/internal/config"
	"github.com/sanjana/email-validator/internal/core"
	"github.com/sanjana/email-validator/internal/db"
	"github.com/sanjana/email-validator/internal/intelligence"
	"github.com/sanjana/email-validator/internal/service"
)


var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4")).
			MarginBottom(1)
	
	validStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
	invalidStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
	warnStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00"))
)

type model struct {
	textInput  textinput.Model
	err        error
	validating bool
	result     *core.EmailResult
	history    []*core.EmailResult
	database   *db.DB
	config     *config.AppConfig
}

func (m model) getStyle(valid bool) lipgloss.Style {
	if valid { return validStyle }
	return invalidStyle
}

type validationMsg *core.EmailResult

func InitialModel(database *db.DB, cfg *config.AppConfig) model {
	ti := textinput.New()
	ti.Placeholder = "Enter email to validate..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 30

	var history []*core.EmailResult
	if database != nil {
		history, _ = database.GetHistory(5)
		// Ensure intelligence modules are initialized for the TUI's DB connection
		_, _ = intelligence.InitDisposable(database)
	}

	return model{
		textInput: ti,
		err:       nil,
		history:   history,
		database:  database,
		config:    cfg,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			if m.textInput.Value() != "" {
				val := m.textInput.Value()
				m.textInput.SetValue("") // Clear input
				m.validating = true
				m.result = nil
				return m, m.validateEmail(val)
			}
		}

	case validationMsg:
		m.validating = false
		m.result = msg
		m.history = append([]*core.EmailResult{msg}, m.history...)
		if len(m.history) > 5 {
			m.history = m.history[:5]
		}
		return m, nil

	case error:
		m.err = msg
		m.validating = false
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) validateEmail(email string) tea.Cmd {
	return func() tea.Msg {
		res, err := service.ProcessEmail(m.database, m.config.SMTPSender, email, "TUI", "127.0.0.1", 0, "Validator-TUI", false)
		if err != nil {
			return err
		}
		return validationMsg(res)
	}
}

func (m model) View() string {
	s := titleStyle.Render("🚀 Email-Validator Pro ")
	s += lipgloss.NewStyle().Foreground(lipgloss.Color("#555555")).Render("(Interactive Mode)") + "\n\n"
	
	s += m.textInput.View() + "\n\n"
	
	var mainContent string
	if m.validating {
		mainContent = warnStyle.Render("⏳ Validating Intelligence Signals...") + "\n"
	} else if m.result != nil {
		res := m.result
		statusStyle := m.getStyle(res.IsValid)

		aliasStr := "No"
		if res.HasAlias { aliasStr = warnStyle.Render("Yes (+alias)") }
		
		smtpStr := validStyle.Render("✅ Valid")
		if res.Greylisted {
			smtpStr = warnStyle.Render("⏳ Greylisted")
		} else if res.SMTPBlocked && !res.SMTP {
			smtpStr = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFA500")).Render("🛡️ Blocked (P25)")
		} else if !res.SMTP {
			smtpStr = invalidStyle.Render("❌ NotFound")
		}

		dispStr := validStyle.Render("No")
		if res.Disposable { dispStr = invalidStyle.Bold(true).Render("⚠️  YES") }

		roleStr := validStyle.Render("No")
		if res.Role { roleStr = invalidStyle.Render("Yes") }

		catchAllStr := validStyle.Render("No")
		if res.CatchAll { catchAllStr = warnStyle.Render("⚠️  YES") }

		ageStr := fmt.Sprintf("%.1f yrs", res.DomainAgeYears)
		if res.DomainAgeYears < 0 { ageStr = "New/Unknown" }

		engagementColor := "#00FF00"
		if res.EngagementProbability < 40 { engagementColor = "#FF0000" } else if res.EngagementProbability < 75 { engagementColor = "#FFFF00" }
		engagementStr := lipgloss.NewStyle().Foreground(lipgloss.Color(engagementColor)).Bold(true).Render(fmt.Sprintf("%d%%", res.EngagementProbability))

		// Truncate long emails for TUI display
		displayEmail := res.Email
		if len(displayEmail) > 30 { displayEmail = displayEmail[:27] + "..." }
		displayBase := res.BaseEmail
		if len(displayBase) > 25 { displayBase = displayBase[:22] + "..." }

		mainContent = fmt.Sprintf(
				"Email Address:   %s\n"+
				"Alias/Identity:  %s (Base: %s)\n"+
				"Reputation:     %d/100 (%s)\n"+
				"Engagement:     %s Probability\n"+
				"──────────────────────────────────────────────────\n"+
				"DNS/MX Record: %v\n"+
				"SMTP Delivery: %s\n"+
				"Last Response: %s\n"+
				"Disposable:    %s\n"+
				"Role-based:    %s\n"+
				"Catch-All Hub: %s\n"+
				"Lifecycle:     %s (Last: %s)\n"+
				"Domain Age:    %s (TLD: %s)\n"+
				"Provider:      %s\n"+
				"Requester IP:  %s",
				statusStyle.Render(displayEmail),
				aliasStr, displayBase,
				res.ReputationScore, res.RiskLevel,
				engagementStr,
				res.DNS,
				smtpStr,
				core.GetShortSMTPResponse(res.LastSMTPResponse),
				dispStr,
				roleStr,
				catchAllStr,
				res.LifecycleState, res.LastVerifiedAt,
				ageStr, res.TldTrust,
				res.Provider,
				lipgloss.NewStyle().Foreground(lipgloss.Color("#777777")).Render(res.ClientIP),
			)
	// 📡 FULL TELEMETRY (Consistent with CLI)
		if !res.SMTP && res.LastSMTPResponse != "" && !strings.Contains(res.LastSMTPResponse, "250") {
			telemetryTitle := lipgloss.NewStyle().Foreground(lipgloss.Color("#555555")).Bold(true).PaddingTop(1).Render("📡 FULL TELEMETRY")
			wrappedTelemetry := lipgloss.NewStyle().Foreground(lipgloss.Color("#777777")).Italic(true).Width(50).Render(res.LastSMTPResponse)
			mainContent += "\n\n" + telemetryTitle + "\n" + wrappedTelemetry
		}

		if res.Message != "" {
			mainContent += "\n\n" + invalidStyle.Bold(true).Render("⚠️  " + res.Message)
		}

		if len(res.EngagementFactors) > 0 {
			analysisTitle := lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4")).Bold(true).PaddingTop(1).Render("ENGAGEMENT ANALYSIS:")
			mainContent += "\n\n" + analysisTitle + "\n"
			for i, factor := range res.EngagementFactors {
				if i >= 3 { break } // Show top 3 in TUI to save space
				mainContent += lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAAA")).Render("• " + factor) + "\n"
			}
		}

	} else {
		mainContent = lipgloss.NewStyle().Foreground(lipgloss.Color("#444444")).Render("Waiting for input...")
	}


	historyTitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")).
		Bold(true).
		Underline(true).
		PaddingLeft(2).
		Render("RECENT SCANS")
	
	historyContent := historyTitle + "\n"
	if len(m.history) == 0 {
		historyContent += lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#444444")).Render("No history yet")
	} else {
		for _, h := range m.history {
			statusIcon := "✅"
			if !h.IsValid { statusIcon = "❌" }
			
			emailStr := h.Email
			if len(emailStr) > 30 {
				emailStr = emailStr[:27] + "..."
			}
			
			emailStyle := lipgloss.NewStyle().Width(32).PaddingLeft(1)
			iconStyle := lipgloss.NewStyle().Width(3)
			historyContent += lipgloss.JoinHorizontal(lipgloss.Top,
				iconStyle.Render(statusIcon),
				emailStyle.Render(emailStr),
			) + "\n"
		}
	}

	mainBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#444444")).
		Padding(1, 2).
		Width(54).
		Render(mainContent)

	sidebarBox := lipgloss.NewStyle().
		Padding(1, 2).
		Render(historyContent)

	s += lipgloss.JoinHorizontal(lipgloss.Top, mainBox, sidebarBox) + "\n\n"
	s += lipgloss.NewStyle().Foreground(lipgloss.Color("#555555")).Render(" (esc to quit)") + "\n"
	return s
}

func Start(database *db.DB, cfg *config.AppConfig) error {
	p := tea.NewProgram(InitialModel(database, cfg))
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}
