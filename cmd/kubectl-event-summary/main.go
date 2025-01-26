// cmd/kubectl-event-summary/main.go
package main

import (
    "os"

    "kubectl-event-summary/pkg/cmd"
)

func main() {
    command := cmd.NewEventSummaryCommand()
    if err := command.Execute(); err != nil {
        os.Exit(1)
    }
}