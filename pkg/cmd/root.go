// /Users/mocha/code/kubectl-event-summary/pkg/cmd/root.go
package cmd

import (
    "context"
    "fmt"
    "os"
    "sort"
    "time"

    "github.com/spf13/cobra"
    "k8s.io/cli-runtime/pkg/genericclioptions"
    "k8s.io/client-go/kubernetes"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/utils/ptr"
    corev1 "k8s.io/api/core/v1"

    "kubectl-event-summary/pkg/events"
)

// EventSummaryOptions contains the options for the event-summary command
type EventSummaryOptions struct {
    configFlags *genericclioptions.ConfigFlags
    allNs       bool
    sortBy      string
    format      string
    since       time.Duration  // Added for time window configuration

    genericclioptions.IOStreams
}

// NewEventSummaryOptions returns initialized EventSummaryOptions
func NewEventSummaryOptions(streams genericclioptions.IOStreams) *EventSummaryOptions {
    return &EventSummaryOptions{
        configFlags: genericclioptions.NewConfigFlags(true),
        IOStreams:  streams,
    }
}

// NewEventSummaryCommand creates the event-summary command
func NewEventSummaryCommand() *cobra.Command {
    o := events.NewEventSummaryOptions(genericclioptions.IOStreams{
        In:     os.Stdin,
        Out:    os.Stdout,
        ErrOut: os.Stderr,
    })

    cmd := &cobra.Command{
        Use:          "event-summary [flags]",
        Short:        "Summarize Kubernetes events",
        SilenceUsage: true,
        RunE: func(c *cobra.Command, args []string) error {
            if err := o.Complete(c, args); err != nil {
                return err
            }
            if err := o.Validate(); err != nil {
                return err
            }
            if err := o.Run(); err != nil {
                return err
            }
            return nil
        },
    }

    AddFlags(cmd, o)
    return cmd
}

// Complete completes all the required options
func (o *EventSummaryOptions) Complete(cmd *cobra.Command, args []string) error {
    return nil
}

// Validate validates the provided options
func (o *EventSummaryOptions) Validate() error {
    if o.allNs && o.configFlags.Namespace != nil && *o.configFlags.Namespace != "" {
        return fmt.Errorf("--namespace and --all-namespaces cannot be used together")
    }
    return nil
}

// Run runs the command
func (o *EventSummaryOptions) Run() error {
    config, err := o.configFlags.ToRESTConfig()
    if err != nil {
        return fmt.Errorf("failed to get client config: %v", err)
    }

    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        return fmt.Errorf("failed to create clientset: %v", err)
    }

    timeWindow := time.Now().Add(-o.since)

    var namespace string
    if !o.allNs {
        namespace = *o.configFlags.Namespace
    }

    eventList, err := clientset.CoreV1().Events(namespace).List(context.TODO(), metav1.ListOptions{
        TimeoutSeconds: ptr.To[int64](10),
    })
    if err != nil {
        return fmt.Errorf("failed to list events: %v", err)
    }

    // Filter and sort events
    var filteredEvents []corev1.Event
    for _, event := range eventList.Items {
        if event.LastTimestamp.Time.After(timeWindow) {
            filteredEvents = append(filteredEvents, event)
        }
    }

    // Sort events based on sortBy flag
    switch o.sortBy {
    case "count":
        sort.Slice(filteredEvents, func(i, j int) bool {
            return filteredEvents[i].Count > filteredEvents[j].Count
        })
    case "lastTimestamp":
        sort.Slice(filteredEvents, func(i, j int) bool {
            return filteredEvents[i].LastTimestamp.Time.After(filteredEvents[j].LastTimestamp.Time)
        })
    }

    // Format and display events based on output format
    switch o.format {
    case "wide":
        for _, event := range filteredEvents {
            fmt.Fprintf(o.Out, "[%s] %s/%s: %s (count: %d)\n",
                event.Type,
                event.InvolvedObject.Namespace,
                event.InvolvedObject.Name,
                event.Message,
                event.Count)
        }
    case "json":
        // Add JSON output
        fmt.Println("JSON output not implemented")
    case "yaml":
        // Add YAML output
        fmt.Println("YAML output not implemented")
    }

    return nil
}