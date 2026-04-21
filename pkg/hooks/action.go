package hooks

import (
	"fmt"

	"github.com/dubin555/azlake/pkg/graveler"
	"gopkg.in/yaml.v3"
)

// Action represents a named set of hooks triggered by specific events.
type Action struct {
	Name  string                           `yaml:"name" json:"name"`
	On    map[graveler.EventType]*ActionOn `yaml:"on" json:"on"`
	Hooks []HookConfig                     `yaml:"hooks" json:"hooks"`
}

// ActionOn specifies branch filters for an event trigger.
type ActionOn struct {
	Branches []string `yaml:"branches,omitempty" json:"branches,omitempty"`
}

// Validate checks that the action definition is well-formed.
func (a *Action) Validate() error {
	if a.Name == "" {
		return fmt.Errorf("'name' is required: %w", ErrInvalidAction)
	}
	if len(a.On) == 0 {
		return fmt.Errorf("'on' is required: %w", ErrInvalidAction)
	}
	for event := range a.On {
		if !event.IsValid() {
			return fmt.Errorf("event '%s' not supported: %w", event, ErrInvalidAction)
		}
	}
	ids := make(map[string]struct{})
	for i, h := range a.Hooks {
		if h.ID == "" {
			return fmt.Errorf("hook[%d] missing ID: %w", i, ErrInvalidAction)
		}
		if _, dup := ids[h.ID]; dup {
			return fmt.Errorf("hook[%d] duplicate ID '%s': %w", i, h.ID, ErrInvalidAction)
		}
		ids[h.ID] = struct{}{}
		if h.Type != HookTypeWebhook {
			return fmt.Errorf("hook[%d] unsupported type '%s': %w", i, h.Type, ErrInvalidAction)
		}
		if h.URL == "" {
			return fmt.Errorf("hook[%d] missing url: %w", i, ErrInvalidAction)
		}
	}
	return nil
}

// ParseAction unmarshals YAML bytes into a validated Action.
func ParseAction(data []byte) (*Action, error) {
	var act Action
	if err := yaml.Unmarshal(data, &act); err != nil {
		return nil, err
	}
	if err := act.Validate(); err != nil {
		return nil, err
	}
	return &act, nil
}
