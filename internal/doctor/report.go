package doctor

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/bkuri/ppc/internal/model"
	"github.com/bkuri/ppc/internal/resolver"
)

// DoctorStats represents module statistics
type DoctorStats struct {
	Modules     int            `json:"modules"`
	ByLayer     map[string]int `json:"by_layer"`
	Unreachable int            `json:"unreachable"`
	Tags        int            `json:"tags"`
	Groups      int            `json:"groups"`
	Orphaned    int            `json:"orphaned"`
}

// DoctorReport represents the complete doctor report
type DoctorReport struct {
	Status   string       `json:"status"`
	Modules  int          `json:"modules"`
	Errors   []string     `json:"errors,omitempty"`
	Warnings []string     `json:"warnings,omitempty"`
	Stats    *DoctorStats `json:"stats,omitempty"`
}

// printDoctorJSON outputs doctor results as JSON
// Returns exit code: 0=ok, 2=failed
func printDoctorJSON(moduleCount int, errs, warns []string, strict bool, stats *DoctorStats) int {
	status := "ok"
	exitCode := 0

	if len(errs) > 0 {
		status = "failed"
		exitCode = 2
	} else if strict && len(warns) > 0 {
		status = "failed"
		exitCode = 2
	}

	report := DoctorReport{
		Status:   status,
		Modules:  moduleCount,
		Errors:   errs,
		Warnings: warns,
		Stats:    stats,
	}

	b, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "json marshal error: %v\n", err)
		return 2
	}

	fmt.Println(string(b))
	return exitCode
}

// calculateStats computes module statistics
func calculateStats(modByID map[string]*model.Module, rules *model.Rules, reachable map[string]bool) *DoctorStats {
	byLayer := map[string]int{"base": 0, "modes": 0, "traits": 0, "policies": 0, "contracts": 0}

	for _, m := range modByID {
		layerName := model.LayerName(m.Layer)
		byLayer[layerName]++
	}

	tagValues := map[string]bool{}
	for _, m := range modByID {
		for _, t := range m.Front.Tags {
			if g, v, ok := resolver.ParseKeyedTag(t); ok {
				tagValues[g+v] = true
			}
		}
	}

	unreachable := 0
	for id := range modByID {
		if !reachable[id] {
			unreachable++
		}
	}

	orphaned := 0
	for _, m := range modByID {
		for _, r := range m.Front.Requires {
			if _, ok := modByID[r]; !ok {
				orphaned++
			}
		}
	}

	return &DoctorStats{
		Modules:     len(modByID),
		ByLayer:     byLayer,
		Unreachable: unreachable,
		Tags:        len(tagValues),
		Groups:      len(rules.ExclusiveGroups),
		Orphaned:    orphaned,
	}
}
