package events

import (
    "context"
    "fmt"
    "strings"
    "time"

    "github.com/spf13/cobra"
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/utils/ptr"

    "kubectl-event-summary/pkg/types"
)

// Complete completes all the required options
func (o *EventSummaryOptions) Complete(cmd *cobra.Command, args []string) error {
    return nil
}

// Validate validates the provided options
func (o *EventSummaryOptions) Validate() error {
    switch o.Severity {
    case types.SeverityAll, types.SeverityNormal, types.SeverityWarning, types.SeverityError:
        // valid severity
    default:
        return fmt.Errorf("invalid severity: %s, must be one of: all, normal, warning, error", o.Severity)
    }

    if o.Format != "wide" && o.Format != "json" && o.Format != "yaml" {
        return fmt.Errorf("invalid format: %s, must be one of: wide, json, yaml", o.Format)
    }

    if o.SortBy != "lastTimestamp" && o.SortBy != "count" {
        return fmt.Errorf("invalid sort-by: %s, must be one of: lastTimestamp, count", o.SortBy)
    }

    return nil
}

// Run executes the command
func (o *EventSummaryOptions) Run() error {
    config, err := o.ConfigFlags.ToRESTConfig()
    if err != nil {
        return fmt.Errorf("failed to get client config: %v", err)
    }

    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        return fmt.Errorf("failed to create clientset: %v", err)
    }

    var namespace string
    if !o.AllNs {
        var explicit bool
        namespace, explicit, err = o.ConfigFlags.ToRawKubeConfigLoader().Namespace()
        if err != nil {
            return fmt.Errorf("failed to get namespace: %v", err)
        }
        if !explicit {
            fmt.Fprintf(o.ErrOut, "Using namespace %q\n", namespace)
        }
    }

    eventList, err := clientset.CoreV1().Events(namespace).List(context.TODO(), metav1.ListOptions{
        TimeoutSeconds: ptr.To[int64](10),
    })
    if err != nil {
        return fmt.Errorf("failed to list events: %v", err)
    }

    // First, count total events and warnings/errors before filtering
    totalInitialEvents := len(eventList.Items)
    totalInitialWarnings := 0
    totalInitialErrors := 0
    for _, event := range eventList.Items {
        if event.Type == "Warning" {
            totalInitialWarnings++
            if strings.Contains(strings.ToLower(event.Reason), "error") ||
               strings.Contains(strings.ToLower(event.Reason), "failed") ||
               strings.Contains(strings.ToLower(event.Reason), "backoff") {
                totalInitialErrors++
            }
        }
    }

    // Filter events by time window and search string
    var filteredEvents []corev1.Event
    timeWindow := time.Now().Add(-o.Since)
    for _, event := range eventList.Items {
        // Use EventTime if available, otherwise use LastTimestamp, then FirstTimestamp
        eventTime := event.EventTime.Time
        if eventTime.IsZero() {
            eventTime = event.LastTimestamp.Time
        }
        if eventTime.IsZero() {
            eventTime = event.FirstTimestamp.Time
        }

        // Include events that happened at or after the time window
        if eventTime.Before(timeWindow) {
            continue
        }

        // Apply search filter if specified
        if o.Search != "" {
            searchLower := strings.ToLower(o.Search)
            // Check various fields for the search string
            if !strings.Contains(strings.ToLower(event.InvolvedObject.Name), searchLower) &&
               !strings.Contains(strings.ToLower(event.Message), searchLower) &&
               !strings.Contains(strings.ToLower(event.Reason), searchLower) &&
               !strings.Contains(strings.ToLower(event.InvolvedObject.Namespace), searchLower) &&
               !strings.Contains(strings.ToLower(event.InvolvedObject.Kind), searchLower) {
                continue
            }
        }

        filteredEvents = append(filteredEvents, event)
    }

    // If no events found after filtering, show a message with total events
    if len(filteredEvents) == 0 {
        fmt.Fprintf(o.Out, "\nTotal Events in cluster: %d (Warnings: %d, Errors: %d)\n", 
            totalInitialEvents, 
            totalInitialWarnings,
            totalInitialErrors)
        if o.Search != "" {
            fmt.Fprintf(o.Out, "No events found matching search term: %q\n", o.Search)
        } else {
            fmt.Fprintf(o.Out, "No events found matching the specified criteria\n")
        }
        return nil
    }

    // Group events if grouping is requested
    var groups map[string]*types.GroupSummary
    var keys []string
    if o.GroupBy != "" {
        groups, keys, err = groupEvents(filteredEvents, o.GroupBy, o.Filter, o.Severity)
        if err != nil {
            return err
        }
        // Add initial totals to all groups
        for _, summary := range groups {
            summary.InitialTotal = totalInitialEvents
            summary.InitialWarnings = totalInitialWarnings
            summary.InitialErrors = totalInitialErrors
        }
    } else {
        // Create a single group for all events
        // First filter events by severity
        var severityFilteredEvents []corev1.Event
        for _, event := range filteredEvents {
            if shouldIncludeEvent(event, o.Severity) {
                severityFilteredEvents = append(severityFilteredEvents, event)
            }
        }
        
        // Use severity as the key instead of "all"
        groupKey := string(o.Severity)
        if groupKey == string(types.SeverityAll) {
            groupKey = "all events"
        }
        
        groups = map[string]*types.GroupSummary{
            groupKey: {
                Types:           make(map[string]int),
                Reasons:         make(map[string]int),
                Events:         severityFilteredEvents,
                InitialTotal:   totalInitialEvents,
                InitialWarnings: totalInitialWarnings,
                InitialErrors:  totalInitialErrors,
            },
        }
        keys = []string{groupKey}
        
        // Update the summary counts
        summary := groups[groupKey]
        for _, event := range severityFilteredEvents {
            summary.Total++
            if event.Type == "Warning" {
                summary.Warnings++
                // Count errors based on reason
                if strings.Contains(strings.ToLower(event.Reason), "error") ||
                   strings.Contains(strings.ToLower(event.Reason), "failed") ||
                   strings.Contains(strings.ToLower(event.Reason), "backoff") {
                    summary.Errors++
                }
            }
            summary.Types[event.Type]++
            summary.Reasons[event.Reason]++
        }
    }

    // Format and display events based on output format
    switch o.Format {
    case "wide":
        return o.printWideFormat(groups, keys)
    case "json":
        return o.printJSONFormat(groups, keys)
    case "yaml":
        return o.printYAMLFormat(groups, keys)
    default:
        return fmt.Errorf("unsupported format: %s", o.Format)
    }
}

// Helper functions for different output formats
func (o *EventSummaryOptions) printWideFormat(groups map[string]*types.GroupSummary, keys []string) error {
    // Get initial totals from any group (they're all the same)
    var initialTotals *types.GroupSummary
    for _, summary := range groups {
        initialTotals = summary
        break
    }

    // Calculate totals for filtered events
    totalFilteredEvents := 0
    totalFilteredWarnings := 0
    totalFilteredErrors := 0
    
    for _, summary := range groups {
        totalFilteredEvents += summary.Total
        totalFilteredWarnings += summary.Warnings
        totalFilteredErrors += summary.Errors
    }

    // Print overall cluster events summary first
    fmt.Fprintf(o.Out, "\nTotal Events in cluster: %d (Warnings: %d, Errors: %d)\n", 
        initialTotals.InitialTotal, 
        initialTotals.InitialWarnings,
        initialTotals.InitialErrors)

    // Always show filtered events summary when using severity filter or grouping
    if o.Severity != types.SeverityAll || o.GroupBy != "" {
        fmt.Fprintf(o.Out, "Filtered Events: %d (Warnings: %d, Errors: %d)\n", 
            totalFilteredEvents, 
            totalFilteredWarnings, 
            totalFilteredErrors)
    }
    fmt.Fprintln(o.Out, "---")

    // Print group details
    for _, key := range keys {
        summary := groups[key]
        fmt.Fprintf(o.Out, "\n=== %s ===\n", key)
        fmt.Fprintf(o.Out, "Events in group: %d (Warnings: %d, Errors: %d)\n", 
            summary.Total, 
            summary.Warnings, 
            summary.Errors)
        
        if !o.Compact {
            for _, event := range summary.Events {
                fmt.Fprintf(o.Out, "[%s] %s/%s: %s (count: %d)\n",
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

func (o *EventSummaryOptions) printJSONFormat(groups map[string]*types.GroupSummary, keys []string) error {
    // TODO: Implement JSON output
    return fmt.Errorf("JSON output not implemented yet")
}

func (o *EventSummaryOptions) printYAMLFormat(groups map[string]*types.GroupSummary, keys []string) error {
    // TODO: Implement YAML output
    return fmt.Errorf("YAML output not implemented yet")
}