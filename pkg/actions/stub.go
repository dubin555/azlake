package actions

import (
	"errors"
	"time"

	"github.com/dubin555/azlake/pkg/graveler"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrParamConflict = errors.New("parameter conflict")
	ErrActionFailed  = errors.New("action failed")
)

type RunResult struct {
	RunID     string
	BranchID  string
	CommitID  string
	SourceRef string
	EventType string
	StartTime time.Time
	EndTime   time.Time
	Passed    bool
}

type TaskResult struct {
	RunID      string
	HookRunID  string
	HookID     string
	ActionName string
	StartTime  time.Time
	EndTime    time.Time
	Passed     bool
}

type RunResultIterator interface {
	Next() bool
	Value() *RunResult
	Err() error
	Close()
}

type TaskResultIterator interface {
	Next() bool
	Value() *TaskResult
	Err() error
	Close()
}

// Source is a stub interface for the actions source.
type Source interface {
	List(ctx interface{}, record graveler.HookRecord) ([]string, error)
	Load(ctx interface{}, record graveler.HookRecord, name string) ([]byte, error)
}
