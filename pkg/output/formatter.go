package output

import (
	"io"
	
	"github.com/nareshku/kubectl-event-summary/pkg/types"
)

// Formatter defines the interface for output formatters
type Formatter interface {
	Format(groups map[string]*types.GroupSummary, keys []string) error
}

// NewFormatter creates a new formatter based on the format string
func NewFormatter(format string, out io.Writer, compact bool) (Formatter, error) {
	switch format {
	case "wide":
		return &WideFormatter{out: out, compact: compact}, nil
	case "json":
		return &JSONFormatter{out: out}, nil
	case "yaml":
		return &YAMLFormatter{out: out}, nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
} 