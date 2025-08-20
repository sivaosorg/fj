package fj

// Shared utility functions for Phase 2 query processing system

// isEmptyString checks if a string is empty
func isEmptyString(s string) bool {
	return len(s) == 0
}

// containsRune checks if a string contains a specific rune
func containsRune(s string, r rune) bool {
	for _, char := range s {
		if char == r {
			return true
		}
	}
	return false
}

// containsAny checks if a string contains any of the given substrings
func containsAny(s string, substrings []string) bool {
	for _, substr := range substrings {
		for i := 0; i <= len(s)-len(substr); i++ {
			if s[i:i+len(substr)] == substr {
				return true
			}
		}
	}
	return false
}

// simpleFloatParse converts string to float with basic parsing
func simpleFloatParse(s string) float64 {
	if isEmptyString(s) {
		return 0
	}
	
	var result float64
	var decimal bool
	var decimalPlace float64 = 0.1
	
	for _, r := range s {
		if r >= '0' && r <= '9' {
			if decimal {
				result += float64(r-'0') * decimalPlace
				decimalPlace *= 0.1
			} else {
				result = result*10 + float64(r-'0')
			}
		} else if r == '.' && !decimal {
			decimal = true
		} else {
			break // Stop at first non-numeric character
		}
	}
	return result
}