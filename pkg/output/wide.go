package output

import (
    "fmt"
    "io"
    
    "github.com/nareshku/kubectl-event-summary/pkg/events"
)

type WideFormatter struct {
    out     io.Writer
    compact bool
}

func (f *WideFormatter) Format(groups map[string]*events.GroupSummary, keys []string) error {
    for _, key := range keys {
        summary := groups[key]
        
        // Print group header with summary
        fmt.Fprintf(f.out, "\n=== %s ===\n", key)
        fmt.Fprintf(f.out, "Total Events: %d (Warnings: %d)\n", summary.Total, summary.Warnings)
        
        // Print type distribution
        fmt.Fprintf(f.out, "Event Types: ")
        for eventType, count := range summary.Types {
            fmt.Fprintf(f.out, "%s=%d ", eventType, count)
        }
        fmt.Fprintln(f.out)

        // Print reason distribution
        fmt.Fprintf(f.out, "Event Reasons: ")
        for reason, count := range summary.Reasons {
            fmt.Fprintf(f.out, "%s=%d ", reason, count)
        }
        fmt.Fprintln(f.out)

        if !f.compact {
            for _, event := range summary.Events {
                fmt.Fprintf(f.out, "[%s] %s/%s: %s (count: %d)\n",
                    event.Type,
                    event.InvolvedObject.Namespace,
                    event.InvolvedObject.Name,
                    event.Message,
                    event.Count)
            }
        }
    }
    return nil
} 