package gooo

import "testing"

func TestEvaluateRevisionMetricsClassifiesLocalizedChange(t *testing.T) {
	metrics, err := MeasureRevision(validContract, metricsRevisionApplication(t))
	if err != nil {
		t.Fatalf("MeasureRevision() error = %v", err)
	}
	quality, err := EvaluateRevisionMetrics(metrics)
	if err != nil {
		t.Fatalf("EvaluateRevisionMetrics() error = %v", err)
	}
	if quality.Status != "BOUND" || quality.ChangeClass != "localized" || quality.ReviewSignal != "observe" {
		t.Fatalf("unexpected quality signal: %#v", quality)
	}
	if quality.SourceDigest != metrics.SourceDigest || quality.MetricsDigest != metrics.MetricsDigest {
		t.Fatalf("quality provenance mismatch: %#v", quality)
	}
	if err := quality.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestEvaluateRevisionMetricsRetainsUnknownForTamperedMetrics(t *testing.T) {
	metrics, err := MeasureRevision(validContract, metricsRevisionApplication(t))
	if err != nil {
		t.Fatalf("MeasureRevision() error = %v", err)
	}
	metrics.MetricsDigest = digestString("tampered")
	quality, err := EvaluateRevisionMetrics(metrics)
	if err == nil {
		t.Fatal("EvaluateRevisionMetrics() error = nil, want metrics validation failure")
	}
	if quality.Status != "UNKNOWN" || quality.MissingStage != "revision-quality-metrics" {
		t.Fatalf("unexpected unknown quality: %#v", quality)
	}
}

func TestEvaluateRevisionMetricsIsDeterministic(t *testing.T) {
	metrics, err := MeasureRevision(validContract, metricsRevisionApplication(t))
	if err != nil {
		t.Fatalf("MeasureRevision() error = %v", err)
	}
	first, err := EvaluateRevisionMetrics(metrics)
	if err != nil {
		t.Fatalf("first EvaluateRevisionMetrics() error = %v", err)
	}
	second, err := EvaluateRevisionMetrics(metrics)
	if err != nil {
		t.Fatalf("second EvaluateRevisionMetrics() error = %v", err)
	}
	if first.QualityDigest != second.QualityDigest {
		t.Fatal("same metrics produced different quality digest")
	}
}
