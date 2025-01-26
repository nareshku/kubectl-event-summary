package types

import (
	corev1 "k8s.io/api/core/v1"
)

// Severity represents the event severity level
type Severity string

const (
	SeverityAll     Severity = "all"
	SeverityNormal  Severity = "normal"
	SeverityWarning Severity = "warning"
	SeverityError   Severity = "error"
)

// GroupSummary holds statistics for a group of events
type GroupSummary struct {
	Total    int
	Warnings int
	Errors   int
	Types    map[string]int
	Reasons  map[string]int
	Events   []corev1.Event

	// Initial totals before filtering
	InitialTotal    int
	InitialWarnings int
	InitialErrors   int
} 