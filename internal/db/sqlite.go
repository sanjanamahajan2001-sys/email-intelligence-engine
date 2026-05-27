package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sanjana/email-validator/internal/core"
)

// DB represents the database connection.
type DB struct {
	conn *sql.DB
}

// InitDB initializes the SQLite database with production hardening (WAL mode & busy timeout).
func InitDB(dbPath string) (*DB, error) {
	dsn := dbPath + "?_busy_timeout=5000&_journal_mode=WAL"
	conn, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	db := &DB{conn: conn}
	if err := db.createTables(); err != nil {
		return nil, err
	}
	return db, nil
}

func (db *DB) createTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS scans (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT NOT NULL,
			base_email TEXT,
			has_alias BOOLEAN,
			is_valid BOOLEAN,
			syntax BOOLEAN,
			dns BOOLEAN,
			smtp BOOLEAN,
			disposable BOOLEAN,
			role BOOLEAN,
			domain_age_years REAL,
			reputation_score INTEGER,
			risk_level TEXT,
			provider TEXT,
			tld_trust TEXT,
			source TEXT,
			created_at DATETIME,
			catch_all BOOLEAN,
			greylisted BOOLEAN,
			identity_age_years REAL,
			confidence_score INTEGER,
			engagement_probability INTEGER,
			last_smtp_response TEXT,
			engagement_insight TEXT,
			lifecycle_state TEXT,
			engagement_factors TEXT,
			message TEXT,
			client_ip TEXT,
			user_id INTEGER,
			user_agent TEXT,
			processing_time_ms INTEGER
		)`,
		`CREATE INDEX IF NOT EXISTS idx_scans_email ON scans(email)`,
		`CREATE TABLE IF NOT EXISTS disposable_domains (
			domain TEXT PRIMARY KEY,
			provider_name TEXT,
			source TEXT,
			created_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS discovery_queue (
			domain TEXT PRIMARY KEY,
			mx_hosts TEXT,
			status TEXT DEFAULT 'PENDING',
			created_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS metadata (
			key TEXT PRIMARY KEY,
			val TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS mx_signatures (
			signature TEXT PRIMARY KEY,
			provider_name TEXT,
			created_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS domain_intelligence (
			domain TEXT PRIMARY KEY,
			age REAL,
			provider TEXT,
			trust TEXT,
			created_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			created_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS refresh_tokens (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			token TEXT UNIQUE NOT NULL,
			expires_at DATETIME NOT NULL,
			FOREIGN KEY(user_id) REFERENCES users(id)
		)`,
	}

	for _, query := range queries {
		if _, err := db.conn.Exec(query); err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) GetLatestResult(email string) (*core.EmailResult, error) {
	query := `SELECT 
		email, COALESCE(base_email, ''), COALESCE(has_alias, 0), COALESCE(is_valid, 0), COALESCE(syntax, 0), COALESCE(dns, 0), COALESCE(smtp, 0), COALESCE(disposable, 0), COALESCE(role, 0), 
		COALESCE(domain_age_years, -1), COALESCE(reputation_score, 0), COALESCE(risk_level, 'Unknown'), COALESCE(provider, 'Unknown'), COALESCE(tld_trust, 'Unknown'), COALESCE(source, 'Internal'), COALESCE(created_at, ''),
		COALESCE(catch_all, 0), COALESCE(greylisted, 0), COALESCE(identity_age_years, 0), COALESCE(confidence_score, 0), COALESCE(engagement_probability, 0), 
		COALESCE(last_smtp_response, ''), COALESCE(engagement_insight, ''), 
		COALESCE(lifecycle_state, 'ACTIVE'), COALESCE(engagement_factors, '[]'),
		COALESCE(message, ''), COALESCE(client_ip, ''), COALESCE(user_id, 0)
		FROM scans 
		WHERE email = ? 
		ORDER BY id DESC 
		LIMIT 1`
	
	res := core.NewResult(email)
	var factorsJSON string
	err := db.conn.QueryRow(query, email).Scan(
		&res.Email, &res.BaseEmail, &res.HasAlias, &res.IsValid, &res.Syntax, &res.DNS, &res.SMTP, &res.Disposable, &res.Role,
		&res.DomainAgeYears, &res.ReputationScore, &res.RiskLevel, &res.Provider, &res.TldTrust, &res.Source, &res.CreatedAt,
		&res.CatchAll, &res.Greylisted, &res.IdentityAgeYears, &res.ConfidenceScore, &res.EngagementProbability, 
		&res.LastSMTPResponse, &res.EngagementInsight, &res.LifecycleState, &factorsJSON,
		&res.Message, &res.ClientIP, &res.UserID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	json.Unmarshal([]byte(factorsJSON), &res.EngagementFactors)
	res.LastVerifiedAt = res.CreatedAt
	return res, nil
}

func (db *DB) GetScanHistory(email string) ([]*core.EmailResult, error) {
	query := `SELECT 
		email, COALESCE(base_email, ''), COALESCE(has_alias, 0), COALESCE(is_valid, 0), COALESCE(syntax, 0), COALESCE(dns, 0), COALESCE(smtp, 0), COALESCE(disposable, 0), COALESCE(role, 0), 
		COALESCE(domain_age_years, -1), COALESCE(reputation_score, 0), COALESCE(risk_level, ''), COALESCE(provider, ''), COALESCE(tld_trust, ''), COALESCE(source, ''), COALESCE(created_at, ''),
		COALESCE(catch_all, 0), COALESCE(greylisted, 0), COALESCE(identity_age_years, 0), COALESCE(confidence_score, 0), COALESCE(engagement_probability, 0), 
		COALESCE(last_smtp_response, ''), COALESCE(engagement_insight, ''), 
		COALESCE(lifecycle_state, ''), COALESCE(engagement_factors, '[]'),
			COALESCE(message, ''), COALESCE(client_ip, ''), COALESCE(user_id, 0),
			COALESCE(user_agent, ''), COALESCE(processing_time_ms, 0)
			FROM scans 
			WHERE email = ? 
			ORDER BY created_at ASC`
		
		rows, err := db.conn.Query(query, email)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
	
		var results []*core.EmailResult
		for rows.Next() {
			res := core.NewResult(email)
			var factorsJSON string
			err := rows.Scan(
				&res.Email, &res.BaseEmail, &res.HasAlias, &res.IsValid, &res.Syntax, &res.DNS, &res.SMTP, &res.Disposable, &res.Role,
				&res.DomainAgeYears, &res.ReputationScore, &res.RiskLevel, &res.Provider, &res.TldTrust, &res.Source, &res.CreatedAt,
				&res.CatchAll, &res.Greylisted, &res.IdentityAgeYears, &res.ConfidenceScore, &res.EngagementProbability, 
				&res.LastSMTPResponse, &res.EngagementInsight, &res.LifecycleState, &factorsJSON,
				&res.Message, &res.ClientIP, &res.UserID, &res.UserAgent, &res.ProcessingTimeMs,
			)
		if err == nil {
			json.Unmarshal([]byte(factorsJSON), &res.EngagementFactors)
			results = append(results, res)
		}
	}
	return results, nil
}

func (db *DB) SaveScan(res *core.EmailResult, source string) error {
	factorsJSON, _ := json.Marshal(res.EngagementFactors)
	query := `INSERT INTO scans (
		email, base_email, has_alias, is_valid, syntax, dns, smtp, disposable, role, 
		domain_age_years, reputation_score, risk_level, provider, tld_trust, source, created_at,
		catch_all, greylisted, identity_age_years, confidence_score, engagement_probability, 
		last_smtp_response, engagement_insight, lifecycle_state, engagement_factors, message, client_ip, user_id, user_agent, processing_time_ms
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := db.conn.Exec(query,
		res.Email, res.BaseEmail, res.HasAlias, res.IsValid, res.Syntax, res.DNS, res.SMTP, res.Disposable, res.Role,
		res.DomainAgeYears, res.ReputationScore, res.RiskLevel, res.Provider, res.TldTrust, source, res.CreatedAt,
		res.CatchAll, res.Greylisted, res.IdentityAgeYears, res.ConfidenceScore, res.EngagementProbability, 
		res.LastSMTPResponse, res.EngagementInsight, res.LifecycleState, string(factorsJSON), res.Message, res.ClientIP, res.UserID,
		res.UserAgent, res.ProcessingTimeMs,
	)
	return err
}

func (db *DB) GetHistory(limit int) ([]*core.EmailResult, error) {
	query := `SELECT 
		email, COALESCE(base_email, ''), COALESCE(has_alias, 0), COALESCE(is_valid, 0), COALESCE(syntax, 0), COALESCE(dns, 0), COALESCE(smtp, 0), COALESCE(disposable, 0), COALESCE(role, 0), 
		COALESCE(domain_age_years, -1), COALESCE(reputation_score, 0), COALESCE(risk_level, ''), COALESCE(provider, ''), COALESCE(tld_trust, ''), COALESCE(source, ''), COALESCE(created_at, ''),
		COALESCE(catch_all, 0), COALESCE(greylisted, 0), COALESCE(identity_age_years, 0), COALESCE(confidence_score, 0), COALESCE(engagement_probability, 0), 
		COALESCE(last_smtp_response, ''), COALESCE(engagement_insight, ''), 
		COALESCE(lifecycle_state, ''), COALESCE(engagement_factors, '[]'),
			COALESCE(message, ''), COALESCE(client_ip, ''), COALESCE(user_id, 0),
			COALESCE(user_agent, ''), COALESCE(processing_time_ms, 0)
			FROM scans 
			ORDER BY created_at DESC 
			LIMIT ?`
		
		rows, err := db.conn.Query(query, limit)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
	
		var results []*core.EmailResult
		for rows.Next() {
			res := core.NewResult("")
			var factorsJSON string
			err := rows.Scan(
				&res.Email, &res.BaseEmail, &res.HasAlias, &res.IsValid, &res.Syntax, &res.DNS, &res.SMTP, &res.Disposable, &res.Role,
				&res.DomainAgeYears, &res.ReputationScore, &res.RiskLevel, &res.Provider, &res.TldTrust, &res.Source, &res.CreatedAt,
				&res.CatchAll, &res.Greylisted, &res.IdentityAgeYears, &res.ConfidenceScore, &res.EngagementProbability, 
				&res.LastSMTPResponse, &res.EngagementInsight, &res.LifecycleState, &factorsJSON,
				&res.Message, &res.ClientIP, &res.UserID, &res.UserAgent, &res.ProcessingTimeMs,
			)
		if err == nil {
			json.Unmarshal([]byte(factorsJSON), &res.EngagementFactors)
			results = append(results, res)
		}
	}
	return results, nil
}

func (db *DB) GetFirstSeenAge(email string) (float64, error) {
	var firstSeen string
	query := `SELECT MIN(created_at) FROM scans WHERE email = ?`
	err := db.conn.QueryRow(query, email).Scan(&firstSeen)
	if err != nil || firstSeen == "" {
		return 0, nil
	}

	t, err := core.ParseFlexibleTime(firstSeen)
	if err != nil {
		return 0, nil
	}

	return time.Since(t).Hours() / 24 / 365, nil
}

func (db *DB) GetHistoricalSuccessRate(email string) float64 {
	var total, valid int
	// We count success based on SMTP being true
	err := db.conn.QueryRow(`SELECT COUNT(*), SUM(CASE WHEN smtp = 1 THEN 1 ELSE 0 END) FROM scans WHERE email = ?`, email).Scan(&total, &valid)
	if err != nil || total == 0 {
		return 0
	}
	return float64(valid) / float64(total)
}

func (db *DB) SaveDisposableDomains(domains map[string]string) error {
	tx, err := db.conn.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`INSERT OR REPLACE INTO disposable_domains (domain, provider_name, source, created_at) VALUES (?, ?, ?, ?)`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	now := time.Now().UTC().Format(time.RFC3339)
	for domain, provider := range domains {
		if _, err = stmt.Exec(domain, provider, "DISCOVERY", now); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (db *DB) GetAllDisposableDomains() (map[string]string, error) {
	rows, err := db.conn.Query(`SELECT domain, provider_name FROM disposable_domains`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make(map[string]string)
	for rows.Next() {
		var domain, provider string
		if err := rows.Scan(&domain, &provider); err == nil {
			results[domain] = provider
		}
	}
	return results, nil
}

func (db *DB) IsDisposable(domain string) (bool, string, string) {
	var provider string
	query := `SELECT provider_name FROM disposable_domains WHERE domain = ?`
	err := db.conn.QueryRow(query, domain).Scan(&provider)
	if err == nil {
		return true, provider, "Static Match"
	}
	return false, "", ""
}

func (db *DB) AddToDiscoveryQueue(domain string, mxHosts []string) error {
	mxJSON, _ := json.Marshal(mxHosts)
	query := `INSERT OR IGNORE INTO discovery_queue (domain, mx_hosts, created_at) VALUES (?, ?, ?)`
	_, err := db.conn.Exec(query, domain, string(mxJSON), time.Now().UTC().Format(time.RFC3339))
	return err
}

func (db *DB) GetMetadata(key string) (string, error) {
	var val string
	query := `SELECT val FROM metadata WHERE key = ?`
	err := db.conn.QueryRow(query, key).Scan(&val)
	if err != nil {
		return "", err
	}
	return val, nil
}

func (db *DB) SetMetadata(key string, val string) error {
	query := `INSERT OR REPLACE INTO metadata (key, val) VALUES (?, ?)`
	_, err := db.conn.Exec(query, key, val)
	return err
}

func (db *DB) SaveMXSignatures(signatures map[string]string) error {
	tx, err := db.conn.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`INSERT OR REPLACE INTO mx_signatures (signature, provider_name, created_at) VALUES (?, ?, ?)`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	now := time.Now().UTC().Format(time.RFC3339)
	for sig, provider := range signatures {
		if _, err = stmt.Exec(sig, provider, now); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (db *DB) GetAllMXSignatures() (map[string]string, error) {
	rows, err := db.conn.Query(`SELECT signature, provider_name FROM mx_signatures`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make(map[string]string)
	for rows.Next() {
		var sig, provider string
		if err := rows.Scan(&sig, &provider); err == nil {
			results[sig] = provider
		}
	}
	return results, nil
}

func (db *DB) GetRecentDomains(limit int) ([]string, error) {
	query := `SELECT DISTINCT SUBSTR(email, INSTR(email, '@') + 1) as domain 
	          FROM scans ORDER BY created_at DESC LIMIT ?`
	rows, err := db.conn.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var domains []string
	for rows.Next() {
		var domain string
		if err := rows.Scan(&domain); err == nil {
			domains = append(domains, domain)
		}
	}
	return domains, nil
}

func (db *DB) GetDomainIntelligence(domain string) (float64, string, string, bool) {
	var age float64
	var provider, trust string
	query := `SELECT age, provider, trust FROM domain_intelligence WHERE domain = ?`
	err := db.conn.QueryRow(query, domain).Scan(&age, &provider, &trust)
	if err != nil {
		return -1, "", "", false
	}
	return age, provider, trust, true
}

func (db *DB) SaveDomainIntelligence(domain string, age float64, provider string, trust string) error {
	query := `INSERT OR REPLACE INTO domain_intelligence (domain, age, provider, trust, created_at) VALUES (?, ?, ?, ?, ?)`
	_, err := db.conn.Exec(query, domain, age, provider, trust, time.Now().UTC())
	return err
}

func (db *DB) GetEmailsForUpdate() ([]string, error) {
	query := `
		SELECT email FROM (
			SELECT email, COALESCE(lifecycle_state, 'ACTIVE') as state, MAX(created_at) as last_seen 
			FROM scans 
			GROUP BY email
		) AS latest
		WHERE (state IN ('ACTIVE', 'FULL') AND last_seen < datetime('now', '-7 days'))
		   OR (state IN ('CATCH-ALL', 'STALE') AND last_seen < datetime('now', '-14 days'))
		   OR (state = 'ABANDONED' AND last_seen < datetime('now', '-30 days'))
		   OR (state IN ('INVALID', 'DISPOSABLE') AND last_seen < datetime('now', '-90 days'))
	`
	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var emails []string
	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err == nil {
			emails = append(emails, email)
		}
	}
	return emails, nil
}
func (db *DB) CountUsers() (int, error) {
	var count int
	err := db.conn.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&count)
	return count, err
}

func (db *DB) GetUsers() ([]map[string]interface{}, error) {
	rows, err := db.conn.Query(`SELECT id, username, created_at FROM users ORDER BY id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var id int
		var username, createdAt string
		if err := rows.Scan(&id, &username, &createdAt); err == nil {
			users = append(users, map[string]interface{}{
				"id":         id,
				"username":   username,
				"created_at": createdAt,
			})
		}
	}
	return users, nil
}

func (db *DB) CreateUser(username, passwordHash string) error {
	query := `INSERT INTO users (username, password_hash, created_at) VALUES (?, ?, ?)`
	_, err := db.conn.Exec(query, username, passwordHash, time.Now().UTC())
	return err
}

func (db *DB) GetUserByUsername(username string) (int, string, string, error) {
	var id int
	var uname, hash string
	query := `SELECT id, username, password_hash FROM users WHERE username = ?`
	err := db.conn.QueryRow(query, username).Scan(&id, &uname, &hash)
	return id, uname, hash, err
}

func (db *DB) SaveRefreshToken(userID int, token string, expiresAt time.Time) error {
	query := `INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES (?, ?, ?)`
	_, err := db.conn.Exec(query, userID, token, expiresAt.UTC())
	return err
}

func (db *DB) ValidateRefreshToken(token string) (int, error) {
	var userID int
	var expiresAt time.Time
	query := `SELECT user_id, expires_at FROM refresh_tokens WHERE token = ?`
	err := db.conn.QueryRow(query, token).Scan(&userID, &expiresAt)
	if err != nil {
		return 0, err
	}
	if time.Now().After(expiresAt) {
		db.DeleteRefreshToken(token)
		return 0, fmt.Errorf("token expired")
	}
	return userID, nil
}

func (db *DB) DeleteRefreshToken(token string) error {
	query := `DELETE FROM refresh_tokens WHERE token = ?`
	_, err := db.conn.Exec(query, token)
	return err
}

func (db *DB) GetUsernameByID(id int) (string, error) {
	var username string
	query := `SELECT username FROM users WHERE id = ?`
	err := db.conn.QueryRow(query, id).Scan(&username)
	return username, err
}

func (db *DB) DeleteUser(username string) error {
	// 1. Get user ID first
	var userID int
	err := db.conn.QueryRow(`SELECT id FROM users WHERE username = ?`, username).Scan(&userID)
	if err != nil {
		return err
	}

	// Start a transaction for safe deletion
	tx, err := db.conn.Begin()
	if err != nil {
		return err
	}

	// 2. Delete refresh tokens
	_, err = tx.Exec(`DELETE FROM refresh_tokens WHERE user_id = ?`, userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 3. Delete user
	_, err = tx.Exec(`DELETE FROM users WHERE id = ?`, userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (db *DB) UpdateUserPassword(username, newHash string) error {
	query := `UPDATE users SET password_hash = ? WHERE username = ?`
	res, err := db.conn.Exec(query, newHash, username)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}
