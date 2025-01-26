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
go build -o kubectl-event-summary ./cmd/kubectl-event-summary

# Move to PATH

sudo mv kubectl-event-summary /usr/local/bin/
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
# Search eventswith a string
$ ./kubectl-event-summary -n kube-system --since 1h --search coredns

Total Events in cluster: 7 (Warnings: 1, Errors: 0)
---

=== all events ===
Events in group: 7 (Warnings: 1, Errors: 0)
[Normal] kube-system/coredns-668d6bf9bc-jmpqz: Stopping container coredns (count: 1)
[Warning] kube-system/coredns-668d6bf9bc-jmpqz: Readiness probe failed: Get "http://10.244.0.5:8181/ready": dial tcp 10.244.0.5:8181: connect: connection refused (count: 1)
[Normal] kube-system/coredns-668d6bf9bc-wkdgd: Successfully assigned kube-system/coredns-668d6bf9bc-wkdgd to kind-control-plane (count: 0)
[Normal] kube-system/coredns-668d6bf9bc-wkdgd: Container image "registry.k8s.io/coredns/coredns:v1.11.3" already present on machine (count: 1)
[Normal] kube-system/coredns-668d6bf9bc-wkdgd: Created container: coredns (count: 1)
[Normal] kube-system/coredns-668d6bf9bc-wkdgd: Started container coredns (count: 1)
[Normal] kube-system/coredns-668d6bf9bc: Created pod: coredns-668d6bf9bc-wkdgd (count: 1)
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
