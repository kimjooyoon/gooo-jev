package gooo

import "testing"

func revisionQualityForFeedback(t *testing.T) RevisionQuality {
	t.Helper()
	metrics, err := MeasureRevision(validContract, metricsRevisionApplication(t))
	if err != nil {
		t.Fatalf("MeasureRevision() error = %v", err)
	}
	quality, err := EvaluateRevisionMetrics(metrics)
	if err != nil {
		t.Fatalf("EvaluateRevisionMetrics() error = %v", err)
	}
	return quality
}

func TestDeriveRevisionFeedbackBindsStructuralHint(t *testing.T) {
	feedback, err := DeriveRevisionFeedback(revisionQualityForFeedback(t))
	if err != nil {
		t.Fatalf("DeriveRevisionFeedback() error = %v", err)
	}
	if feedback.Status != "BOUND" || feedback.ChangeClass != "structural" || feedback.ActionHint != "inspect-ir" {
		t.Fatalf("unexpected feedback: %#v", feedback)
	}
	if !feedback.NonExecuting || !feedback.NonAuthorizing {
		t.Fatal("feedback must remain non-executing and non-authorizing")
	}
	if err := feedback.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestDeriveRevisionFeedbackRetainsUnknownForTamperedQuality(t *testing.T) {
	quality := revisionQualityForFeedback(t)
	quality.QualityDigest = digestString("tampered")
	feedback, err := DeriveRevisionFeedback(quality)
	if err == nil {
		t.Fatal("DeriveRevisionFeedback() error = nil, want quality validation failure")
	}
	if feedback.Status != "UNKNOWN" || feedback.MissingStage != "revision-feedback-quality" {
		t.Fatalf("unexpected unknown feedback: %#v", feedback)
	}
}

func TestDeriveRevisionFeedbackIsDeterministic(t *testing.T) {
	quality := revisionQualityForFeedback(t)
	first, err := DeriveRevisionFeedback(quality)
	if err != nil {
		t.Fatalf("first DeriveRevisionFeedback() error = %v", err)
	}
	second, err := DeriveRevisionFeedback(quality)
	if err != nil {
		t.Fatalf("second DeriveRevisionFeedback() error = %v", err)
	}
	if first.FeedbackDigest != second.FeedbackDigest {
		t.Fatal("same quality signal produced different feedback digest")
	}
}
