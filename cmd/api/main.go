package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sanjana/email-validator/internal/config"
	"github.com/sanjana/email-validator/internal/db"
	"github.com/sanjana/email-validator/internal/intelligence"
	"github.com/sanjana/email-validator/internal/service"
	"github.com/sanjana/email-validator/internal/tracer"
	"golang.org/x/time/rate"
)

var (
	database *db.DB
	appConfig *config.AppConfig
)

type ValidationRequest struct {
	Email string `json:"email" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"` // Seconds
}

// WebValidationResponse provides a rich, consumer-friendly response for web frontends.
type WebValidationResponse struct {
	Email           string       `json:"email"`
	IsValid         bool         `json:"is_valid"`
	Authenticity    string       `json:"authenticity_status"` // "Verified", "Suspicious", "Invalid"
	IsTemporary     bool         `json:"is_temporary"`
	ReputationScore int          `json:"reputation_score"` // 0-100
	Recommendation  string       `json:"recommendation"`   // "Accept", "Flag", "Reject"
	LifecycleState  string       `json:"lifecycle_state"`  // "ACTIVE", "CATCH-ALL", "STALE", etc.
	Engagement      Engagement   `json:"engagement"`
	DetailedInfo    DetailedInfo `json:"detailed_info"`
}

type Engagement struct {
	Probability int      `json:"probability"` // 0-100
	Insight     string   `json:"insight"`     // Human readable advice
	Factors     []string `json:"factors"`     // List of trust/risk signals
}

type DetailedInfo struct {
	DNSActive       bool    `json:"dns_active"`
	SMTPDeliverable bool    `json:"smtp_deliverable"`
	SMTPResponse    string  `json:"smtp_response"`
	Provider        string  `json:"provider"`
	RiskLevel       string  `json:"risk_level"`
	TrustLevel      string  `json:"trust_level"`
	Message         string  `json:"message"`
	DomainAge       float64 `json:"domain_age_years"`
	IdentityAge     float64 `json:"identity_age_years"`
	Confidence      int     `json:"confidence_score"`
	HasAlias        bool    `json:"has_alias"`
	BaseEmail       string  `json:"base_email"`
}

// RateLimiter manages IP-based request counts using token bucket
type RateLimiter struct {
	sync.RWMutex
	ips map[string]*rate.Limiter
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{ips: make(map[string]*rate.Limiter)}
}

func (rl *RateLimiter) Middleware(limit int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		rl.Lock()
		limiter, exists := rl.ips[ip]
		if !exists {
			// Refill rate: 'limit' tokens per minute
			refillRate := rate.Every(time.Minute / time.Duration(limit))
			limiter = rate.NewLimiter(refillRate, limit)
			rl.ips[ip] = limiter
		}
		rl.Unlock()

		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Please try again later.",
			})
			return
		}
		c.Next()
	}
}

// CORSMiddleware enables requests from web frontends
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	appConfig = config.LoadConfig()
	tracer.Info("API", "Starting Email Intelligence API...")

	var err error
	database, err = db.InitDB(appConfig.DBPath)
	if err != nil {
		tracer.Error("API", "Failed to initialize database", err)
		panic(err)
	}

	// Initialize Intelligence Modules
	lastSync, err := intelligence.InitDisposable(database)
	if err != nil {
		tracer.Error("API", "Warning: Failed to load disposable domains", err)
	}

	// Automated Background Sync
	if lastSync.IsZero() || time.Since(lastSync) > 24*time.Hour {
		go func() {
			tracer.Info("Sync", "Background Sync: Updating disposable domain list...")
			count, err := intelligence.SyncDisposable(database)
			if err != nil {
				tracer.Error("Sync", "Background Sync failed", err)
			} else {
				tracer.Info("Sync", fmt.Sprintf("Background Sync complete! Tracked %d domains.", count))
			}
		}()
	}

	startTime := time.Now()
	r := gin.Default()
	
	// Apply CORS
	r.Use(CORSMiddleware())
	
	limiter := NewRateLimiter()

	// v1 endpoints
	v1 := r.Group("/v1")
	{
		// Auth endpoints (Public)
		auth := v1.Group("/auth")
		{
			auth.POST("/login", handleLogin)
			auth.POST("/refresh", handleRefresh)
		}

		// Protected endpoints
		protected := v1.Group("/")
		protected.Use(service.AuthMiddleware(appConfig))
		protected.Use(limiter.Middleware(appConfig.RateLimitPerMin))
		{
			protected.POST("/validate", handleValidate)
			protected.POST("/web-validate", handleWebValidate)
			protected.POST("/sync-disposable", handleSyncDisposable)
		}

		// Public Validation (Demo route for Vercel Web UI, rate-limited to 5/min per IP)
		v1.POST("/public-validate", limiter.Middleware(5), handlePublicValidate)

		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "healthy",
				"uptime":  time.Since(startTime).String(),
				"db":      "connected",
				"version": "1.2.0 (Secure Auth)",
			})
		})
	}

	// Auto-Provision Default Admin User
	count, _ := database.CountUsers()
	if count == 0 {
		defaultPassword := "admin123"
		hash, _ := service.HashPassword(defaultPassword)
		err := database.CreateUser("admin", hash)
		if err == nil {
			tracer.Info("Auth", "************************************************")
			tracer.Info("Auth", "AUTO-PROVISIONING: Created default admin user")
			tracer.Info("Auth", "Username: admin")
			tracer.Info("Auth", "Password: "+defaultPassword)
			tracer.Info("Auth", "************************************************")
		}
	}

	tracer.Info("API", "Listening on :"+appConfig.APIPort)
	r.Run(":" + appConfig.APIPort)
}

func handleValidate(c *gin.Context) {
	var req ValidationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	uID := 0
	if id, ok := userID.(int); ok {
		uID = id
	}

	res, err := service.ProcessEmail(database, appConfig.SMTPSender, req.Email, "API", c.ClientIP(), uID, c.Request.UserAgent(), false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func handleWebValidate(c *gin.Context) {
	var req ValidationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Full validation cycle
	userID, _ := c.Get("user_id")
	uID := 0
	if id, ok := userID.(int); ok {
		uID = id
	}

	res, err := service.ProcessEmail(database, appConfig.SMTPSender, req.Email, "WEB-INTEGRATION", c.ClientIP(), uID, c.Request.UserAgent(), false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Map internal EmailResult to WebValidationResponse
	authenticity := "Invalid"
	if res.IsValid {
		if res.SMTP {
			if res.ReputationScore > 60 {
				authenticity = "Verified"
			} else {
				authenticity = "Suspicious"
			}
		} else if res.SMTPBlocked {
			if res.ReputationScore > 75 {
				authenticity = "Verified"
			} else {
				authenticity = "Suspicious"
			}
		} else {
			authenticity = "Suspicious"
		}
	}

	recommendation := "Reject"
	if authenticity == "Verified" && !res.Disposable && !res.SMTPBlocked {
		recommendation = "Accept"
	} else if authenticity == "Suspicious" || res.CatchAll || res.Role || res.SMTPBlocked {
		recommendation = "Flag"
	}

	response := WebValidationResponse{
		Email:           res.Email,
		IsValid:         res.IsValid,
		Authenticity:    authenticity,
		IsTemporary:     res.Disposable,
		ReputationScore: res.ReputationScore,
		Recommendation:  recommendation,
		LifecycleState:  res.LifecycleState,
		Engagement: Engagement{
			Probability: res.EngagementProbability,
			Insight:     res.EngagementInsight,
			Factors:     res.EngagementFactors,
		},
		DetailedInfo: DetailedInfo{
			DNSActive:       res.DNS,
			SMTPDeliverable: res.SMTP,
			SMTPResponse:    res.LastSMTPResponse,
			Provider:        res.Provider,
			RiskLevel:       res.RiskLevel,
			TrustLevel:      res.TldTrust,
			Message:         res.Message,
			DomainAge:       res.DomainAgeYears,
			IdentityAge:     res.IdentityAgeYears,
			Confidence:      res.ConfidenceScore,
			HasAlias:        res.HasAlias,
			BaseEmail:       res.BaseEmail,
		},
	}

	c.JSON(http.StatusOK, response)
}

func handlePublicValidate(c *gin.Context) {
	var req ValidationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Process email anonymously (userID = 0) with a public demo source tag
	res, err := service.ProcessEmail(database, appConfig.SMTPSender, req.Email, "WEB-PUBLIC-DEMO", c.ClientIP(), 0, c.Request.UserAgent(), false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Map internal EmailResult to WebValidationResponse
	authenticity := "Invalid"
	if res.IsValid {
		if res.SMTP {
			if res.ReputationScore > 60 {
				authenticity = "Verified"
			} else {
				authenticity = "Suspicious"
			}
		} else if res.SMTPBlocked {
			if res.ReputationScore > 75 {
				authenticity = "Verified"
			} else {
				authenticity = "Suspicious"
			}
		} else {
			authenticity = "Suspicious"
		}
	}

	recommendation := "Reject"
	if authenticity == "Verified" && !res.Disposable && !res.SMTPBlocked {
		recommendation = "Accept"
	} else if authenticity == "Suspicious" || res.CatchAll || res.Role || res.SMTPBlocked {
		recommendation = "Flag"
	}

	response := WebValidationResponse{
		Email:           res.Email,
		IsValid:         res.IsValid,
		Authenticity:    authenticity,
		IsTemporary:     res.Disposable,
		ReputationScore: res.ReputationScore,
		Recommendation:  recommendation,
		LifecycleState:  res.LifecycleState,
		Engagement: Engagement{
			Probability: res.EngagementProbability,
			Insight:     res.EngagementInsight,
			Factors:     res.EngagementFactors,
		},
		DetailedInfo: DetailedInfo{
			DNSActive:       res.DNS,
			SMTPDeliverable: res.SMTP,
			SMTPResponse:    res.LastSMTPResponse,
			Provider:        res.Provider,
			RiskLevel:       res.RiskLevel,
			TrustLevel:      res.TldTrust,
			Message:         res.Message,
			DomainAge:       res.DomainAgeYears,
			IdentityAge:     res.IdentityAgeYears,
			Confidence:      res.ConfidenceScore,
			HasAlias:        res.HasAlias,
			BaseEmail:       res.BaseEmail,
		},
	}

	c.JSON(http.StatusOK, response)
}

func handleSyncDisposable(c *gin.Context) {
	count, err := intelligence.SyncDisposable(database)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully synchronized disposable domains",
		"count":   count,
	})
}

func handleLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, username, hash, err := database.GetUserByUsername(req.Username)
	if err != nil || !service.CheckPasswordHash(req.Password, hash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	accessToken, err := service.GenerateAccessToken(appConfig, id, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken := service.GenerateRefreshToken()
	err = database.SaveRefreshToken(id, refreshToken, time.Now().Add(appConfig.RefreshTokenDuration))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save refresh token"})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(appConfig.AccessTokenDuration.Seconds()),
	})
}

func handleRefresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := database.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	username, err := database.GetUsernameByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user info"})
		return
	}

	accessToken, err := service.GenerateAccessToken(appConfig, userID, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: req.RefreshToken,
		ExpiresIn:    int(appConfig.AccessTokenDuration.Seconds()),
	})
}
