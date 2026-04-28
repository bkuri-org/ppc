package profile

import (
	"fmt"
	"strings"
)

type Profile struct {
	Mode      string         `yaml:"mode"`
	Contract  string         `yaml:"contract"`
	Revisions *int           `yaml:"revisions,omitempty"`
	Traits    []string       `yaml:"traits,omitempty"`
	Vars      map[string]any `yaml:"vars,omitempty"`
}

func (p *Profile) Validate() error {
	if p.Mode == "" {
		return fmt.Errorf("profile: mode is required")
	}
	if p.Contract == "" {
		return fmt.Errorf("profile: contract is required")
	}
	for _, t := range p.Traits {
		if t != "" && !strings.HasPrefix(t, "traits/") {
			return fmt.Errorf("profile: trait %q must be module ID (e.g., traits/conservative)", t)
		}
	}
	return nil
}
