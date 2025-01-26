package events

import (
	"strings"
	"sort"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"github.com/nareshku/kubectl-event-summary/pkg/types"
)

func buildGroupKey(event corev1.Event, groupLevels []string) string {
	var parts []string
	for _, level := range groupLevels {
		var value string
		switch level {
		case "kind":
			value = event.InvolvedObject.Kind
		case "namespace":
			value = event.InvolvedObject.Namespace
		case "reason":
			value = event.Reason
		case "type":
			value = event.Type
		}
		parts = append(parts, level+"="+value)
	}
	return strings.Join(parts, ",")
}

func shouldIncludeEvent(event corev1.Event, severity types.Severity) bool {
	switch severity {
	case types.SeverityAll:
		return true
	case types.SeverityNormal:
		return event.Type == "Normal"
	case types.SeverityWarning:
		return event.Type == "Warning"
	case types.SeverityError:
		// In Kubernetes events, errors are typically marked as Warning type
		// but we can check for specific error-like reasons
		return event.Type == "Warning" && (
			strings.Contains(strings.ToLower(event.Reason), "error") ||
			strings.Contains(strings.ToLower(event.Reason), "failed") ||
			strings.Contains(strings.ToLower(event.Reason), "backoff"))
	default:
		return true
	}
}

func groupEvents(events []corev1.Event, groupBy string, filter string, severity types.Severity) (map[string]*types.GroupSummary, []string, error) {
	groupLevels := strings.Split(groupBy, ",")
	groups := make(map[string]*types.GroupSummary)
	
	// Build groups and collect statistics
	for _, event := range events {
		// Check severity filter
		if !shouldIncludeEvent(event, severity) {
			continue
		}

		groupKey := buildGroupKey(event, groupLevels)
		
		// Apply filter if specified
		if filter != "" {
			filterParts := strings.Split(filter, "=")
			if len(filterParts) != 2 {
				return nil, nil, fmt.Errorf("invalid filter format. Use 'field=value'")
			}
			if !strings.HasPrefix(groupKey, filterParts[0]+"="+filterParts[1]) {
				continue
			}
		}

		if groups[groupKey] == nil {
			groups[groupKey] = &types.GroupSummary{
				Types:    make(map[string]int),
				Reasons:  make(map[string]int),
			}
		}
		
		summary := groups[groupKey]
		summary.Total++
		if event.Type == "Warning" {
			summary.Warnings++
		}
		summary.Types[event.Type]++
		summary.Reasons[event.Reason]++
		summary.Events = append(summary.Events, event)
	}

	// Sort group keys
	var keys []string
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return groups, keys, nil
} 