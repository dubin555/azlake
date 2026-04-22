package stats

import (
	"context"
	"time"

	"github.com/dubin555/azlake/pkg/logging"
)

type Collector interface {
	CollectEvent(ev Event)
	CollectMetadata(metadata *Metadata)
	SetInstallationID(installationID string)
	Close()
}

type Event struct {
	Class string
	Name  string
	Count int
}

type Metadata struct {
	InstallationID string
	Entries        []MetadataEntry
}

type MetadataEntry struct {
	Name  string
	Value string
}

type NullCollector struct{}

func (n *NullCollector) CollectEvent(ev Event)             {}
func (n *NullCollector) CollectMetadata(metadata *Metadata) {}
func (n *NullCollector) SetInstallationID(id string)       {}
func (n *NullCollector) Close()                            {}

type TimeFn func() time.Time

type FlushTicker struct{}

func NewNullCollector() Collector {
	return &NullCollector{}
}

type BufferedCollector struct{}

type BufferedCollectorOpts struct {
	Address          string
	InstallationID   string
	ProcessID        string
	FlushInterval    time.Duration
	FlushSize        int
	SendTimeout      time.Duration
}

func NewBufferedCollector(installationID string, cfg interface{}, opts ...BufferedCollectorOpts) Collector {
	return &NullCollector{}
}

func NewBufferedCollectorFromConfig(ctx context.Context, cfg interface{}) Collector {
	return &NullCollector{}
}

type UsageReporterOperations interface {
	// Stub - will be filled in later
}

// LoggerAdapter adapts a logging.Logger to the retryablehttp.LeveledLogger interface.
type LoggerAdapter struct {
	Logger logging.Logger
}

func (l *LoggerAdapter) Printf(msg string, args ...any) {
	l.Logger.Tracef(msg, args...)
}

