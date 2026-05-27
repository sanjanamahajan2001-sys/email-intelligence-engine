package main

import (
	"fmt"
	"github.com/sanjana/email-validator/internal/core"
)

func main() {
	testEmails := []struct {
		input    string
		expected string
		hasAlias bool
	}{
		{"john.doe@gmail.com", "johndoe@gmail.com", true},
		{"j.o.h.n.d.o.e@gmail.com", "johndoe@gmail.com", true},
		{"john.doe+extra@gmail.com", "johndoe@gmail.com", true},
		{"normal.email@example.com", "normal.email@example.com", false},
		{"alias+test@outlook.com", "alias@outlook.com", true},
	}

	fmt.Println("Running Gmail Normalization Tests...")
	fmt.Println("---------------------------------------")
	
	passed := 0
	for _, tc := range testEmails {
		result, hasAlias := core.NormalizeEmail(tc.input)
		if result == tc.expected && hasAlias == tc.hasAlias {
			fmt.Printf("✅ PASS: %s -> %s (Alias: %v)\n", tc.input, result, hasAlias)
			passed++
		} else {
			fmt.Printf("❌ FAIL: %s | Expected: %s (%v) | Got: %s (%v)\n", 
				tc.input, tc.expected, tc.hasAlias, result, hasAlias)
		}
	}
	
	fmt.Printf("\nSummary: %d/%d tests passed.\n", passed, len(testEmails))
}
