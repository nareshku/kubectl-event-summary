package events

import (
	"time"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	
	"github.com/nareshku/kubectl-event-summary/pkg/types"
)

// EventSummaryOptions contains the options for the event-summary command
type EventSummaryOptions struct {
	ConfigFlags *genericclioptions.ConfigFlags
	AllNs       bool
	SortBy      string
	Format      string
	Since       time.Duration
	GroupBy     string
	Compact     bool
	Filter      string
	Severity    types.Severity
	Search      string

	genericclioptions.IOStreams
}

// NewEventSummaryOptions returns initialized EventSummaryOptions
func NewEventSummaryOptions(streams genericclioptions.IOStreams) *EventSummaryOptions {
	return &EventSummaryOptions{
		ConfigFlags: genericclioptions.NewConfigFlags(true),
		IOStreams:   streams,
		Severity:    types.SeverityAll,
	}
} 