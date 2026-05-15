package patch

import (
	"strings"
	"testing"
	"time"
)

func makeEntry(id, status string, dur time.Duration, ops int, errMsg string) MetricsEntry {
	return MetricsEntry{
		PatchID:   id,
		AppliedAt: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
		Duration:  dur,
		OpCount:   ops,
		Status:    status,
		Error:     errMsg,
	}
}

func TestSummarise_Empty(t *testing.T) {
	s := Summarise(nil)
	if s.Total != 0 || s.Succeeded != 0 || s.Failed != 0 || s.AvgNs != 0 {
		t.Errorf("expected zero summary, got %+v", s)
	}
}

func TestSummarise_MixedStatuses(t *testing.T) {
	entries := []MetricsEntry{
		makeEntry("p1", "success", 100*time.Millisecond, 3, ""),
		makeEntry("p2", "failure", 200*time.Millisecond, 1, "bad op"),
		makeEntry("p3", "success", 300*time.Millisecond, 5, ""),
	}
	s := Summarise(entries)
	if s.Total != 3 {
		t.Errorf("expected total 3, got %d", s.Total)
	}
	if s.Succeeded != 2 {
		t.Errorf("expected succeeded 2, got %d", s.Succeeded)
	}
	if s.Failed != 1 {
		t.Errorf("expected failed 1, got %d", s.Failed)
	}
	expectedAvg := (100 + 200 + 300) * int64(time.Millisecond) / 3
	if s.AvgNs != expectedAvg {
		t.Errorf("expected avg %d, got %d", expectedAvg, s.AvgNs)
	}
}

func TestFormatMetrics_Empty(t *testing.T) {
	out := FormatMetrics(nil)
	if out != "no metrics recorded" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatMetrics_ContainsPatchID(t *testing.T) {
	entries := []MetricsEntry{
		makeEntry("my-patch-v1", "success", 50*time.Millisecond, 2, ""),
	}
	out := FormatMetrics(entries)
	if !strings.Contains(out, "my-patch-v1") {
		t.Errorf("expected patch id in output, got:\n%s", out)
	}
	if !strings.Contains(out, "success") {
		t.Errorf("expected status in output, got:\n%s", out)
	}
}

func TestFormatMetrics_ShowsError(t *testing.T) {
	entries := []MetricsEntry{
		makeEntry("bad-patch", "failure", 10*time.Millisecond, 1, "unknown op"),
	}
	out := FormatMetrics(entries)
	if !strings.Contains(out, "unknown op") {
		t.Errorf("expected error in output, got:\n%s", out)
	}
}

func TestFormatMetrics_SummaryLine(t *testing.T) {
	entries := []MetricsEntry{
		makeEntry("p1", "success", 100*time.Millisecond, 2, ""),
		makeEntry("p2", "failure", 100*time.Millisecond, 1, "err"),
	}
	out := FormatMetrics(entries)
	if !strings.Contains(out, "total: 2") {
		t.Errorf("expected total in summary, got:\n%s", out)
	}
	if !strings.Contains(out, "succeeded: 1") {
		t.Errorf("expected succeeded in summary, got:\n%s", out)
	}
}
