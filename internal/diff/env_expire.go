package diff

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// ExpireIssue represents a key that is approaching or past its expiration date.
type ExpireIssue struct {
	Key       string
	ExpiresAt time.Time
	Expired   bool
}

func (e ExpireIssue) String() string {
	status := "expires"
	if e.Expired {
		status = "expired"
	}
	return fmt.Sprintf("%s: %s on %s", e.Key, status, e.ExpiresAt.Format("2006-01-02"))
}

// ExpiryRule maps a key pattern to an expiration date string (YYYY-MM-DD).
type ExpiryRule struct {
	Pattern   string
	ExpiresAt time.Time
}

// CheckExpiry checks env keys against expiry rules.
// Keys matching a pattern are flagged if they are expired or within warnDays of expiry.
func CheckExpiry(env map[string]string, rules []ExpiryRule, warnDays int, now time.Time) []ExpireIssue {
	var issues []ExpireIssue
	warnCutoff := now.AddDate(0, 0, warnDays)

	for key := range env {
		for _, rule := range rules {
			if matchesExpiryPattern(key, rule.Pattern) {
				if rule.ExpiresAt.Before(now) || rule.ExpiresAt.Equal(now) {
					issues = append(issues, ExpireIssue{Key: key, ExpiresAt: rule.ExpiresAt, Expired: true})
				} else if rule.ExpiresAt.Before(warnCutoff) {
					issues = append(issues, ExpireIssue{Key: key, ExpiresAt: rule.ExpiresAt, Expired: false})
				}
				break
			}
		}
	}

	sort.Slice(issues, func(i, j int) bool {
		return issues[i].Key < issues[j].Key
	})
	return issues
}

func matchesExpiryPattern(key, pattern string) bool {
	if pattern == "*" {
		return true
	}
	if strings.HasSuffix(pattern, "*") {
		return strings.HasPrefix(key, strings.TrimSuffix(pattern, "*"))
	}
	return key == pattern
}

// FormatExpireIssues returns a human-readable summary of expiry issues.
func FormatExpireIssues(issues []ExpireIssue) string {
	if len(issues) == 0 {
		return "no expiry issues found"
	}
	var sb strings.Builder
	for _, issue := range issues {
		sb.WriteString(issue.String())
		sb.WriteString("\n")
	}
	return strings.TrimRight(sb.String(), "\n")
}
