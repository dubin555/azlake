package hooks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dubin555/azlake/pkg/logging"
)

// EventType represents the type of hook event.
type EventType string

const (
	EventPreCommit  EventType = "pre-commit"
	EventPostCommit EventType = "post-commit"
	EventPreMerge   EventType = "pre-merge"
	EventPostMerge  EventType = "post-merge"
)

// HookConfig represents a single webhook hook configuration.
type HookConfig struct {
	ID         string            `yaml:"id"`
	URL        string            `yaml:"url"`
	Timeout    time.Duration     `yaml:"timeout"`
	Headers    map[string]string `yaml:"headers"`
	QueryParams map[string]string `yaml:"query_params"`
}

// ActionConfig represents a set of hooks for a specific event.
type ActionConfig struct {
	Name  string       `yaml:"name"`
	On    EventType    `yaml:"on"`
	Hooks []HookConfig `yaml:"hooks"`
}

// EventData is the payload sent to webhook endpoints.
type EventData struct {
	EventType    EventType         `json:"event_type"`
	ActionName   string            `json:"action_name"`
	HookID       string            `json:"hook_id"`
	RepositoryID string            `json:"repository_id"`
	BranchID     string            `json:"branch_id,omitempty"`
	SourceRef    string            `json:"source_ref,omitempty"`
	CommitID     string            `json:"commit_id,omitempty"`
	Committer    string            `json:"committer,omitempty"`
	Message      string            `json:"commit_message,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

// Service manages webhook hook execution.
type Service struct {
	actions    []ActionConfig
	httpClient *http.Client
}

// NewService creates a new hooks service.
func NewService(actions []ActionConfig) *Service {
	return &Service{
		actions: actions,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Run fires all hooks configured for the given event type.
// For pre-* events, returns an error if any hook fails (blocking).
// For post-* events, logs errors but does not block.
func (s *Service) Run(ctx context.Context, eventType EventType, data EventData) error {
	log := logging.FromContext(ctx).WithField("event_type", string(eventType))

	for _, action := range s.actions {
		if action.On != eventType {
			continue
		}
		for _, hook := range action.Hooks {
			data.ActionName = action.Name
			data.HookID = hook.ID
			data.EventType = eventType

			err := s.fireWebhook(ctx, hook, data)
			if err != nil {
				log.WithField("hook_id", hook.ID).WithField("url", hook.URL).WithError(err).Error("webhook failed")
				// Pre-events are blocking: if a hook fails, the operation is aborted
				if eventType == EventPreCommit || eventType == EventPreMerge {
					return fmt.Errorf("hook %s failed: %w", hook.ID, err)
				}
				// Post-events are non-blocking: log and continue
			} else {
				log.WithField("hook_id", hook.ID).Debug("webhook succeeded")
			}
		}
	}
	return nil
}

func (s *Service) fireWebhook(ctx context.Context, hook HookConfig, data EventData) error {
	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal event data: %w", err)
	}

	timeout := hook.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, hook.URL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	for k, v := range hook.Headers {
		req.Header.Set(k, v)
	}
	q := req.URL.Query()
	for k, v := range hook.QueryParams {
		q.Set(k, v)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}
	return nil
}
