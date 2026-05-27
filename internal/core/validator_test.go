package core

import (
	"testing"
)

func TestNormalizeEmail(t *testing.T) {
	tests := []struct {
		email     string
		wantBase  string
		wantAlias bool
	}{
		{"user@gmail.com", "user@gmail.com", false},
		{"user+alias@gmail.com", "user@gmail.com", true},
		{"test+123@example.com", "test@example.com", true},
		{"invalid", "invalid", false},
	}

	for _, tt := range tests {
		base, hasAlias := NormalizeEmail(tt.email)
		if base != tt.wantBase || hasAlias != tt.wantAlias {
			t.Errorf("NormalizeEmail(%s) = %s, %v; want %s, %v", tt.email, base, hasAlias, tt.wantBase, tt.wantAlias)
		}
	}
}

func TestCalculateScore(t *testing.T) {
	res := &EmailResult{
		Syntax:         true,
		DNS:            true,
		SMTP:           true,
		Disposable:     false,
		Role:           false,
		DomainAgeYears: 10,
	}

	CalculateScore(res, "")
	if res.ReputationScore != 100 { // 100 base + 10 for age, capped at 100
		t.Errorf("Expected score 100, got %d", res.ReputationScore)
	}

	res.Disposable = true
	CalculateScore(res)
	if res.ReputationScore >= 100 {
		t.Errorf("Expected score < 100 for disposable, got %d", res.ReputationScore)
	}
}
