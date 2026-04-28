package substitute

import (
	"testing"
)

func TestSubstituteSimple(t *testing.T) {
	vars := Vars{
		"goals": map[string]any{
			"target": "production",
		},
	}

	content := "Target: {{goals.target}}"
	result, unresolved := Substitute(content, vars)

	expected := "Target: production"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
	if len(unresolved) != 0 {
		t.Errorf("expected no unresolved vars, got %v", unresolved)
	}
}

func TestSubstituteNested(t *testing.T) {
	vars := Vars{
		"user": map[string]any{
			"profile": map[string]any{
				"name":  "alice",
				"email": "alice@example.com",
			},
		},
	}

	content := "Name: {{user.profile.name}}, Email: {{user.profile.email}}"
	result, unresolved := Substitute(content, vars)

	expected := "Name: alice, Email: alice@example.com"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
	if len(unresolved) != 0 {
		t.Errorf("expected no unresolved vars, got %v", unresolved)
	}
}

func TestSubstituteMissing(t *testing.T) {
	vars := Vars{
		"known": "value",
	}

	content := "Known: {{known}}, Unknown: {{missing}}"
	result, unresolved := Substitute(content, vars)

	if result != "Known: value, Unknown: {{missing}}" {
		t.Errorf("expected unresolved variable to remain, got %q", result)
	}
	if len(unresolved) != 1 || unresolved[0] != "missing" {
		t.Errorf("expected unresolved [missing], got %v", unresolved)
	}
}

func TestSubstituteString(t *testing.T) {
	vars := Vars{"name": "test"}
	result, _ := Substitute("Hello {{name}}", vars)
	if result != "Hello test" {
		t.Errorf("expected %q, got %q", "Hello test", result)
	}
}

func TestSubstituteInt(t *testing.T) {
	vars := Vars{"count": 42}
	result, _ := Substitute("Count: {{count}}", vars)
	if result != "Count: 42" {
		t.Errorf("expected %q, got %q", "Count: 42", result)
	}
}

func TestSubstituteFloat(t *testing.T) {
	vars := Vars{"rate": 3.14159}
	result, _ := Substitute("Rate: {{rate}}", vars)
	if result != "Rate: 3.14159" {
		t.Errorf("expected %q, got %q", "Rate: 3.14159", result)
	}
}

func TestSubstituteBool(t *testing.T) {
	vars := Vars{"enabled": true, "disabled": false}

	result1, _ := Substitute("Enabled: {{enabled}}", vars)
	if result1 != "Enabled: true" {
		t.Errorf("expected %q, got %q", "Enabled: true", result1)
	}

	result2, _ := Substitute("Disabled: {{disabled}}", vars)
	if result2 != "Disabled: false" {
		t.Errorf("expected %q, got %q", "Disabled: false", result2)
	}
}

func TestSubstituteFloatWholeNumber(t *testing.T) {
	vars := Vars{"count": 42.0}
	result, _ := Substitute("Count: {{count}}", vars)
	if result != "Count: 42" {
		t.Errorf("expected whole number format %q, got %q", "Count: 42", result)
	}
}

func TestResolvePathSimple(t *testing.T) {
	vars := Vars{"name": "test"}
	val, ok := ResolvePath(vars, "name")
	if !ok || val != "test" {
		t.Errorf("expected (test, true), got (%v, %v)", val, ok)
	}
}

func TestResolvePathNested(t *testing.T) {
	vars := Vars{
		"level1": map[string]any{
			"level2": map[string]any{
				"value": "deep",
			},
		},
	}
	val, ok := ResolvePath(vars, "level1.level2.value")
	if !ok || val != "deep" {
		t.Errorf("expected (deep, true), got (%v, %v)", val, ok)
	}
}

func TestResolvePathMissing(t *testing.T) {
	vars := Vars{"name": "test"}
	_, ok := ResolvePath(vars, "missing")
	if ok {
		t.Error("expected false for missing path")
	}
}

func TestResolvePathMissingNested(t *testing.T) {
	vars := Vars{
		"level1": map[string]any{
			"level2": "value",
		},
	}
	_, ok := ResolvePath(vars, "level1.level2.level3")
	if ok {
		t.Error("expected false for missing nested path")
	}
}

func TestSubstituteMultipleInLine(t *testing.T) {
	vars := Vars{
		"first": "John",
		"last":  "Doe",
	}
	result, _ := Substitute("{{first}} {{last}}", vars)
	if result != "John Doe" {
		t.Errorf("expected %q, got %q", "John Doe", result)
	}
}

func TestSubstituteNoVariables(t *testing.T) {
	vars := Vars{"name": "test"}
	result, _ := Substitute("No variables here", vars)
	if result != "No variables here" {
		t.Errorf("expected unchanged content, got %q", result)
	}
}

func TestSubstituteMultipleUnresolved(t *testing.T) {
	vars := Vars{}
	content := "{{foo}} and {{bar}} and {{foo}}"
	result, unresolved := Substitute(content, vars)
	if result != "{{foo}} and {{bar}} and {{foo}}" {
		t.Errorf("expected unchanged content, got %q", result)
	}
	if len(unresolved) != 3 {
		t.Errorf("expected 3 unresolved vars, got %d: %v", len(unresolved), unresolved)
	}
}
