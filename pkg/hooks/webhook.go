package hooks

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/dubin555/azlake/pkg/graveler"
	"github.com/dubin555/azlake/pkg/logging"
)

const webhookDefaultTimeout = 1 * time.Minute

var (
	errWebhookRequestFailed = errors.New("webhook request failed")
)

// Webhook fires an HTTP POST with event information to a configured URL.
type Webhook struct {
	ID         string
	ActionName string
	URL        string
	Timeout    time.Duration
	Headers    map[string]string
}

// Run executes the webhook HTTP call.
func (w *Webhook) Run(ctx context.Context, record graveler.HookRecord) error {
	logging.FromContext(ctx).
		WithField("hook_type", "webhook").
		WithField("event_type", record.EventType).
		Debug("webhook executing")

	eventData, err := marshalEventInformation(w.ActionName, w.ID, record)
	if err != nil {
		return fmt.Errorf("marshaling event: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, w.URL, bytes.NewReader(eventData))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range w.Headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{Timeout: w.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("executing webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("%w (status code: %d)", errWebhookRequestFailed, resp.StatusCode)
	}
	return nil
}

// RunWithLog is like Run but also writes request/response details to buf.
func (w *Webhook) RunWithLog(ctx context.Context, record graveler.HookRecord, buf *bytes.Buffer) error {
	eventData, err := marshalEventInformation(w.ActionName, w.ID, record)
	if err != nil {
		return fmt.Errorf("marshaling event: %w", err)
	}

	_, _ = fmt.Fprintf(buf, "POST %s\n", w.URL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, w.URL, bytes.NewReader(eventData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range w.Headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{Timeout: w.Timeout}
	start := time.Now()
	resp, err := client.Do(req)
	elapsed := time.Since(start)
	_, _ = fmt.Fprintf(buf, "Duration: %s\n", elapsed)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if dumpResp, dErr := httputil.DumpResponse(resp, true); dErr == nil {
		buf.Write(dumpResp)
	}

	// Ensure json import is used (for potential future response parsing)
	_ = json.Decoder{}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("%w (status code: %d)", errWebhookRequestFailed, resp.StatusCode)
	}
	return nil
}
