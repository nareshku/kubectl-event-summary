package cmd

import (
	"github.com/spf13/cobra"
	"time"
	
	"github.com/nareshku/kubectl-event-summary/pkg/events"
	"github.com/nareshku/kubectl-event-summary/pkg/types"
)

// AddFlags adds flags to the specified command.
func AddFlags(cmd *cobra.Command, o *events.EventSummaryOptions) {
	o.ConfigFlags.AddFlags(cmd.Flags())
	cmd.Flags().BoolVarP(&o.AllNs, "all-namespaces", "A", false, "If present, summarize events across all namespaces")
	cmd.Flags().StringVar(&o.SortBy, "sort-by", "lastTimestamp", "Sort events by (lastTimestamp, count)")
	cmd.Flags().StringVarP(&o.Format, "output", "o", "wide", "Output format. One of: wide|json|yaml")
	cmd.Flags().DurationVar(&o.Since, "since", 15*time.Minute, "Show events from the last duration (e.g., 5m, 1h)")
	cmd.Flags().StringVar(&o.GroupBy, "group-by", "", "Group events by (comma-separated): kind,namespace,reason,type")
	cmd.Flags().BoolVar(&o.Compact, "compact", false, "Show only group summaries")
	cmd.Flags().StringVar(&o.Filter, "filter", "", "Filter groups by prefix (e.g., 'kind=Pod')")
	cmd.Flags().StringVar((*string)(&o.Severity), "severity", string(types.SeverityAll),
		"Filter events by severity (all|normal|warning|error)")
	cmd.Flags().StringVar(&o.Search, "search", "", 
		"Search string to filter events (searches in name, message, reason, and namespace)")
} 