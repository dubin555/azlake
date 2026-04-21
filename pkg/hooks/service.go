package hooks

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/dubin555/azlake/pkg/graveler"
	"github.com/dubin555/azlake/pkg/logging"
)

func parseTimeout(s string) (time.Duration, error) {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0, fmt.Errorf("invalid timeout %q: %w", s, err)
	}
	return d, nil
}

// Service manages hook actions and fires webhooks on graveler events.
type Service struct {
	actions []*Action
	logger  logging.Logger
}

// NewService creates a Service from a list of actions.
func NewService(actions []*Action, logger logging.Logger) *Service {
	return &Service{actions: actions, logger: logger}
}

// matchActions returns actions that match the given event type and branch.
func (s *Service) matchActions(eventType graveler.EventType, branchID string) []*Action {
	var matched []*Action
	for _, act := range s.actions {
		on, ok := act.On[eventType]
		if !ok {
			continue
		}
		if on == nil || len(on.Branches) == 0 {
			matched = append(matched, act)
			continue
		}
		for _, pattern := range on.Branches {
			if matchBranch(pattern, branchID) {
				matched = append(matched, act)
				break
			}
		}
	}
	return matched
}

func matchBranch(pattern, branch string) bool {
	if pattern == "*" || pattern == branch {
		return true
	}
	if strings.HasPrefix(pattern, "*") && strings.HasSuffix(branch, pattern[1:]) {
		return true
	}
	if strings.HasSuffix(pattern, "*") && strings.HasPrefix(branch, pattern[:len(pattern)-1]) {
		return true
	}
	return false
}

// isPreEvent returns true if the event type is a pre-* (blocking) event.
func isPreEvent(et graveler.EventType) bool {
	return strings.HasPrefix(string(et), "pre-")
}

// RunHooks fires all matching webhooks for the given record.
// For pre-* events, any hook failure returns an error (blocking).
// For post-* events, errors are logged but not returned.
func (s *Service) RunHooks(ctx context.Context, record graveler.HookRecord) error {
	matched := s.matchActions(record.EventType, record.BranchID.String())
	if len(matched) == 0 {
		return nil
	}

	pre := isPreEvent(record.EventType)

	for _, act := range matched {
		for _, hcfg := range act.Hooks {
			hcfg.ActionName = act.Name
			hook, err := NewHook(hcfg)
			if err != nil {
				if pre {
					return fmt.Errorf("action %q hook %q: %w", act.Name, hcfg.ID, err)
				}
				s.logger.WithError(err).
					WithField("action", act.Name).
					WithField("hook", hcfg.ID).
					Error("failed to create hook")
				continue
			}

			if err := hook.Run(ctx, record); err != nil {
				if pre {
					return fmt.Errorf("action %q hook %q: %w", act.Name, hcfg.ID, err)
				}
				s.logger.WithError(err).
					WithField("action", act.Name).
					WithField("hook", hcfg.ID).
					Error("post-event hook failed")
			}
		}
	}
	return nil
}
