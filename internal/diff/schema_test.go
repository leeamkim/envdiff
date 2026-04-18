package diff

import (
	"testing"
)

func TestValidateSchema_NoIssues(t *testing.T) {
	env := map[string]string{"HOST": "localhost", "PORT": "8080"}
	rules := []SchemaRule{SchemaRuleNoEmptyValues}
	issues := ValidateSchema(env, rules)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d", len(issues))
	}
}

func TestValidateSchema_EmptyValue(t *testing.T) {
	env := map[string]string{"HOST": "", "PORT": "8080"}
	rules := []SchemaRule{SchemaRuleNoEmptyValues}
	issues := ValidateSchema(env, rules)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "HOST" {
		t.Errorf("expected key HOST, got %s", issues[0].Key)
	}
}

func TestSchemaRuleRequiredKeys_AllPresent(t *testing.T) {
	env := map[string]string{"HOST": "localhost", "PORT": "8080"}
	rule := SchemaRuleRequiredKeys([]string{"HOST", "PORT"})
	issues := rule(env)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d", len(issues))
	}
}

func TestSchemaRuleRequiredKeys_Missing(t *testing.T) {
	env := map[string]string{"HOST": "localhost"}
	rule := SchemaRuleRequiredKeys([]string{"HOST", "PORT", "DB_URL"})
	issues := rule(env)
	if len(issues) != 2 {
		t.Fatalf("expected 2 issues, got %d", len(issues))
	}
}

func TestSchemaIssue_String(t *testing.T) {
	issue := SchemaIssue{Key: "PORT", Message: "value must not be empty"}
	expected := "[PORT] value must not be empty"
	if issue.String() != expected {
		t.Errorf("expected %q, got %q", expected, issue.String())
	}
}
