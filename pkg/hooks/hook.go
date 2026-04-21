package hooks

import (
	"context"
	"errors"
	"fmt"

	"github.com/dubin555/azlake/pkg/graveler"
)

// HookType enumerates supported hook types.
type HookType string

const (
	HookTypeWebhook HookType = "webhook"
)

var (
	ErrUnknownHookType = errors.New("unknown hook type")
	ErrInvalidAction   = errors.New("invalid action")
)

// Hook is the interface for a runnable hook.
type Hook interface {
	Run(ctx context.Context, record graveler.HookRecord) error
}

// HookConfig describes a single hook from configuration.
type HookConfig struct {
	ID         string            `yaml:"id" json:"id"`
	Type       HookType          `yaml:"type" json:"type"`
	ActionName string            `yaml:"-" json:"-"`
	URL        string            `yaml:"url" json:"url"`
	Timeout    string            `yaml:"timeout,omitempty" json:"timeout,omitempty"`
	Headers    map[string]string `yaml:"headers,omitempty" json:"headers,omitempty"`
}

// NewHook creates a Hook from config. Only webhook type is supported.
func NewHook(cfg HookConfig) (Hook, error) {
	switch cfg.Type {
	case HookTypeWebhook:
		timeout := webhookDefaultTimeout
		if cfg.Timeout != "" {
			d, err := parseTimeout(cfg.Timeout)
			if err != nil {
				return nil, err
			}
			timeout = d
		}
		return &Webhook{
			ID:         cfg.ID,
			ActionName: cfg.ActionName,
			URL:        cfg.URL,
			Timeout:    timeout,
			Headers:    cfg.Headers,
		}, nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnknownHookType, cfg.Type)
	}
}
