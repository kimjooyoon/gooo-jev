package gooo

import "testing"

func metricsRevisionApplication(t *testing.T) RevisionApplication {
	t.Helper()
	edit := SourceEdit{
		Start:       Position{Line: 4, Column: 3},
		End:         Position{Line: 4, Column: 11},
		Replacement: "lineage",
	}
	edit.Digest = digestSourceEdit(edit)
	application, err := ApplyRevision(validContract, digestString(validContract), revisionCandidateForApplication(t), edit)
	if err != nil {
		t.Fatalf("ApplyRevision() error = %v", err)
	}
	return application
}

func TestMeasureRevisionBindsByteLineAndIRMetrics(t *testing.T) {
	application := metricsRevisionApplication(t)
	metrics, err := MeasureRevision(validContract, application)
	if err != nil {
		t.Fatalf("MeasureRevision() error = %v", err)
	}
	if metrics.Status != "BOUND" || metrics.MissingStage != "" {
		t.Fatalf("unexpected metrics status: %#v", metrics)
	}
	if metrics.SourceBytesBefore != len([]byte(validContract)) {
		t.Fatalf("SourceBytesBefore = %d, want %d", metrics.SourceBytesBefore, len([]byte(validContract)))
	}
	if metrics.SourceBytesAfter != len([]byte(application.ProposedSource)) {
		t.Fatalf("SourceBytesAfter = %d, want %d", metrics.SourceBytesAfter, len([]byte(application.ProposedSource)))
	}
	if metrics.SourceBytesDelta != metrics.SourceBytesAfter-metrics.SourceBytesBefore {
		t.Fatalf("SourceBytesDelta = %d, want arithmetic delta", metrics.SourceBytesDelta)
	}
	if metrics.ChangedByteCount <= 0 || metrics.ChangedLineCount != 1 {
		t.Fatalf("unexpected change counts: %#v", metrics)
	}
	if metrics.SourceLinesDelta != 0 {
		t.Fatalf("SourceLinesDelta = %d, want 0", metrics.SourceLinesDelta)
	}
	if metrics.IRChanged {
		t.Fatal("comment-only edit unexpectedly changed IR")
	}
	if !metrics.NonExecuting || !metrics.NonAuthorizing {
		t.Fatal("metrics must remain non-executing and non-authorizing")
	}
	if err := metrics.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestMeasureRevisionRejectsSourceMismatchAsUnknown(t *testing.T) {
	metrics, err := MeasureRevision("stale-source", metricsRevisionApplication(t))
	if err == nil {
		t.Fatal("MeasureRevision() error = nil, want source mismatch")
	}
	if metrics.Status != "UNKNOWN" || metrics.MissingStage != "revision-metrics-source" {
		t.Fatalf("unexpected source mismatch: %#v", metrics)
	}
}

func TestMeasureRevisionRejectsTamperedApplicationAsUnknown(t *testing.T) {
	application := metricsRevisionApplication(t)
	application.ApplicationDigest = digestString("tampered")
	metrics, err := MeasureRevision(validContract, application)
	if err == nil {
		t.Fatal("MeasureRevision() error = nil, want application validation failure")
	}
	if metrics.Status != "UNKNOWN" || metrics.MissingStage != "revision-metrics-application" {
		t.Fatalf("unexpected application failure: %#v", metrics)
	}
}

func TestMeasureRevisionIsDeterministic(t *testing.T) {
	first, err := MeasureRevision(validContract, metricsRevisionApplication(t))
	if err != nil {
		t.Fatalf("first MeasureRevision() error = %v", err)
	}
	second, err := MeasureRevision(validContract, metricsRevisionApplication(t))
	if err != nil {
		t.Fatalf("second MeasureRevision() error = %v", err)
	}
	if first.MetricsDigest != second.MetricsDigest {
		t.Fatal("same revision produced different metrics digest")
	}
}
