package doctor

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/bkuri/ppc/internal/model"
)

func TestRunDoctorValid(t *testing.T) {
	exitCode := RunDoctor("testdata/valid", false, false, false, false, "")
	if exitCode != 0 {
		t.Errorf("expected exit code 0 for valid modules, got %d", exitCode)
	}
}

func TestRunDoctorInvalidTags(t *testing.T) {
	exitCode := captureOutput(func() int {
		return RunDoctor("testdata/invalid_tags", false, false, false, false, "")
	})

	if exitCode != 2 {
		t.Errorf("expected exit code 2 for invalid tags, got %d", exitCode)
	}
}

func TestRunDoctorMissingRequires(t *testing.T) {
	exitCode := captureOutput(func() int {
		return RunDoctor("testdata/missing_requires", false, false, false, false, "")
	})

	if exitCode != 2 {
		t.Errorf("expected exit code 2 for missing requires, got %d", exitCode)
	}
}

func TestRunDoctorCircular(t *testing.T) {
	exitCode := captureOutput(func() int {
		return RunDoctor("testdata/circular", false, false, false, false, "")
	})

	if exitCode != 2 {
		t.Errorf("expected exit code 2 for circular dependencies, got %d", exitCode)
	}
}

func TestRunDoctorUnreachable(t *testing.T) {
	var output bytes.Buffer
	exitCode := captureOutputTo(&output, func() int {
		return RunDoctor("testdata/unreachable", false, false, false, false, "")
	})

	if exitCode != 0 {
		t.Errorf("expected exit code 0 for unreachable (warning only), got %d", exitCode)
	}

	if !strings.Contains(output.String(), "unreachable modules") {
		t.Error("expected warning about unreachable modules")
	}
}

func TestRunDoctorStrictMode(t *testing.T) {
	exitCode := captureOutput(func() int {
		return RunDoctor("testdata/unreachable", true, false, false, false, "")
	})

	if exitCode != 2 {
		t.Errorf("expected exit code 2 in strict mode with warnings, got %d", exitCode)
	}
}

func TestRunDoctorJSON(t *testing.T) {
	var output bytes.Buffer
	exitCode := captureOutputTo(&output, func() int {
		return RunDoctor("testdata/valid", false, true, false, false, "")
	})

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}

	var report DoctorReport
	if err := json.Unmarshal(output.Bytes(), &report); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}

	if report.Status != "ok" {
		t.Errorf("expected status 'ok', got %q", report.Status)
	}

	if report.Modules != 3 {
		t.Errorf("expected 3 modules, got %d", report.Modules)
	}
}

func TestRunDoctorJSONWithErrors(t *testing.T) {
	var output bytes.Buffer
	exitCode := captureOutputTo(&output, func() int {
		return RunDoctor("testdata/invalid_tags", false, true, false, false, "")
	})

	if exitCode != 2 {
		t.Errorf("expected exit code 2, got %d", exitCode)
	}

	var report DoctorReport
	if err := json.Unmarshal(output.Bytes(), &report); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}

	if report.Status != "failed" {
		t.Errorf("expected status 'failed', got %q", report.Status)
	}

	if len(report.Errors) == 0 {
		t.Error("expected errors in report")
	}
}

func TestRunDoctorStats(t *testing.T) {
	var output bytes.Buffer
	exitCode := captureOutputTo(&output, func() int {
		return RunDoctor("testdata/valid", false, true, true, false, "")
	})

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}

	var report DoctorReport
	if err := json.Unmarshal(output.Bytes(), &report); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}

	if report.Stats == nil {
		t.Fatal("expected stats in report")
	}

	if report.Stats.Modules != 3 {
		t.Errorf("expected 3 modules in stats, got %d", report.Stats.Modules)
	}

	if report.Stats.ByLayer["base"] != 1 {
		t.Errorf("expected 1 base module, got %d", report.Stats.ByLayer["base"])
	}

	if report.Stats.ByLayer["modes"] != 1 {
		t.Errorf("expected 1 modes module, got %d", report.Stats.ByLayer["modes"])
	}
}

func TestCalculateStats(t *testing.T) {
	modByID := map[string]*model.Module{
		"base":          {Layer: 0, Front: model.Frontmatter{ID: "base", Tags: []string{"risk:low"}}},
		"modes/explore": {Layer: 1, Front: model.Frontmatter{ID: "modes/explore", Requires: []string{"base"}}},
		"traits/terse":  {Layer: 2, Front: model.Frontmatter{ID: "traits/terse"}},
		"traits/orphan": {Layer: 2, Front: model.Frontmatter{ID: "traits/orphan"}},
	}

	reachable := map[string]bool{
		"base":          true,
		"modes/explore": true,
		"traits/terse":  true,
	}

	rules := &model.Rules{ExclusiveGroups: []string{"risk"}}

	stats := calculateStats(modByID, rules, reachable)

	if stats.Modules != 4 {
		t.Errorf("expected 4 modules, got %d", stats.Modules)
	}

	if stats.Unreachable != 1 {
		t.Errorf("expected 1 unreachable module, got %d", stats.Unreachable)
	}

	if stats.Groups != 1 {
		t.Errorf("expected 1 group, got %d", stats.Groups)
	}
}

func TestLayerNameFromIndex(t *testing.T) {
	tests := []struct {
		idx      int
		expected string
	}{
		{0, "base"},
		{1, "modes"},
		{2, "traits"},
		{3, "policies"},
		{4, "contracts"},
		{99, "unknown"},
	}

	for _, tc := range tests {
		got := layerNameFromIndex(tc.idx)
		if got != tc.expected {
			t.Errorf("layerNameFromIndex(%d) = %q, want %q", tc.idx, got, tc.expected)
		}
	}
}

func captureOutput(fn func() int) int {
	var buf bytes.Buffer
	return captureOutputTo(&buf, fn)
}

func captureOutputTo(buf *bytes.Buffer, fn func() int) int {
	oldStdout := os.Stdout
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stdout = w
	os.Stderr = w

	exitCode := fn()

	w.Close()
	io.Copy(buf, r)
	os.Stdout = oldStdout
	os.Stderr = oldStderr
	r.Close()

	return exitCode
}
