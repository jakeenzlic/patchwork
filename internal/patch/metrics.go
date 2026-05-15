package patch

import (
	"fmt"
	"strings"
	"time"
)

// MetricsEntry records statistics for a single patch application run.
type MetricsEntry struct {
	PatchID   string        `json:"patch_id"`
	AppliedAt time.Time     `json:"applied_at"`
	Duration  time.Duration `json:"duration_ns"`
	OpCount   int           `json:"op_count"`
	Status    string        `json:"status"` // "success" | "failure"
	Error     string        `json:"error,omitempty"`
}

// MetricsSummary aggregates entries for reporting.
type MetricsSummary struct {
	Total     int
	Succeeded int
	Failed    int
	AvgNs     int64
}

// Summarise computes aggregate statistics from a slice of entries.
func Summarise(entries []MetricsEntry) MetricsSummary {
	var s MetricsSummary
	var totalNs int64
	for _, e := range entries {
		s.Total++
		if e.Status == "success" {
			s.Succeeded++
		} else {
			s.Failed++
		}
		totalNs += e.Duration.Nanoseconds()
	}
	if s.Total > 0 {
		s.AvgNs = totalNs / int64(s.Total)
	}
	return s
}

// FormatMetrics returns a human-readable report of the entries.
func FormatMetrics(entries []MetricsEntry) string {
	if len(entries) == 0 {
		return "no metrics recorded"
	}
	var sb strings.Builder
	summary := Summarise(entries)
	sb.WriteString(fmt.Sprintf("total: %d  succeeded: %d  failed: %d  avg_duration: %s\n",
		summary.Total, summary.Succeeded, summary.Failed,
		time.Duration(summary.AvgNs).String()))
	sb.WriteString(strings.Repeat("-", 60) + "\n")
	for _, e := range entries {
		line := fmt.Sprintf("[%s] %-30s ops=%-3d dur=%-12s status=%s",
			e.AppliedAt.Format(time.RFC3339),
			e.PatchID,
			e.OpCount,
			e.Duration.String(),
			e.Status,
		)
		if e.Error != "" {
			line += " err=" + e.Error
		}
		sb.WriteString(line + "\n")
	}
	return sb.String()
}
