# kubectl-event-summary

A kubectl plugin that provides a summarized view of Kubernetes events with powerful filtering, grouping, and search capabilities.

## Features

- **Time-based Filtering**: View events from specific time windows (e.g., last 1h, 30m)
- **Severity Filtering**: Filter events by severity (all|normal|warning|error)
- **Flexible Grouping**: Group events by:
  - Resource kind
  - Event type
  - Event reason
  - Namespace
- **Search**: Search events across multiple fields:
  - Object name
  - Event message
  - Event reason
  - Namespace
  - Resource kind
- **Comprehensive Statistics**: View:
  - Total cluster events
  - Filtered events count
  - Warning and error counts
  - Per-group statistics

## Installation
```
# Clone the repository

git clone https://github.com/nareshku/kubectl-event-summary.git

# Build the plugin

cd kubectl-event-summary
go build -o kubectl-event_summary ./cmd/kubectl-event-summary

# Move to PATH

sudo mv kubectl-event_summary /usr/local/bin/
```

## Usage Examples
1. View all events in current namespace:
```
kubectl event-summary
```

2. View events from all namespaces:
```
kubectl event-summary -A
```

3. View only warning events from last hour:
```
kubectl event-summary --severity warning --since 1h
```

4. Group events by type:
```
kubectl event-summary --group-by type
```


5. Search for specific events using search string:
```
kubectl event-summary --search coredns
```

6. Combine multiple filters:
```
kubectl event-summary -n kube-system --severity warning --search api --since 1h
```

## Sample Output
```
kubectl event-summary -n kube-system --since 1h --group-by type
Total Events in cluster: 50 (Warnings: 5, Errors: 1)
Filtered Events: 10 (Warnings: 2, Errors: 0)
---

=== type=Normal ===
Events in group: 8 (Warnings: 0, Errors: 0)
[Normal] kube-system/coredns-abc: Started container...

=== type=Warning ===
Events in group: 2 (Warnings: 2, Errors: 0)
[Warning] kube-system/coredns-xyz: Readiness probe failed...
```


## Available Flags

- `--all-namespaces, -A`: Show events from all namespaces
- `--since duration`: Show events from the last duration (default: 15m)
- `--severity string`: Filter by severity (all|normal|warning|error)
- `--group-by string`: Group events by field(s) (kind,namespace,reason,type)
- `--search string`: Search string to filter events
- `--compact`: Show only group summaries
- `--output, -o`: Output format (wide|json|yaml)

## Contributing

Contributions are welcome! Feel free to submit issues and pull requests.
