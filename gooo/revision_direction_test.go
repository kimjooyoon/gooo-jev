package gooo

import "testing"

func revisionFeedbackForDirection(t *testing.T) RevisionFeedback {
	t.Helper()
	return revisionFeedbackForQuality(t)
}

func revisionFeedbackForQuality(t *testing.T) RevisionFeedback {
	t.Helper()
	metrics, err := MeasureRevision(validContract, metricsRevisionApplication(t))
	if err != nil {
		t.Fatalf("MeasureRevision() error = %v", err)
	}
	quality, err := EvaluateRevisionMetrics(metrics)
	if err != nil {
		t.Fatalf("EvaluateRevisionMetrics() error = %v", err)
	}
	feedback, err := DeriveRevisionFeedback(quality)
	if err != nil {
		t.Fatalf("DeriveRevisionFeedback() error = %v", err)
	}
	return feedback
}

func TestProposeRevisionDirectionBindsFeedback(t *testing.T) {
	feedback := revisionFeedbackForDirection(t)
	proposal, err := ProposeRevisionDirection(feedback)
	if err != nil {
		t.Fatalf("ProposeRevisionDirection() error = %v", err)
	}
	if proposal.Status != "BOUND" || proposal.ChangeClass != "structural" || proposal.Direction != RepairRevision {
		t.Fatalf("unexpected direction proposal: %#v", proposal)
	}
	if proposal.SourceDigest != feedback.SourceDigest || proposal.FeedbackDigest != feedback.FeedbackDigest {
		t.Fatalf("direction provenance mismatch: %#v", proposal)
	}
	if err := proposal.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestProposeRevisionDirectionRejectsTamperedFeedback(t *testing.T) {
	feedback := revisionFeedbackForDirection(t)
	feedback.FeedbackDigest = digestString("tampered")
	proposal, err := ProposeRevisionDirection(feedback)
	if err == nil {
		t.Fatal("ProposeRevisionDirection() error = nil, want feedback validation failure")
	}
	if proposal.Status != "UNKNOWN" || proposal.MissingStage != "revision-direction-feedback" {
		t.Fatalf("unexpected unknown proposal: %#v", proposal)
	}
}

func TestProposeRevisionDirectionIsDeterministic(t *testing.T) {
	feedback := revisionFeedbackForDirection(t)
	first, err := ProposeRevisionDirection(feedback)
	if err != nil {
		t.Fatalf("first ProposeRevisionDirection() error = %v", err)
	}
	second, err := ProposeRevisionDirection(feedback)
	if err != nil {
		t.Fatalf("second ProposeRevisionDirection() error = %v", err)
	}
	if first.DirectionDigest != second.DirectionDigest {
		t.Fatal("same feedback produced different direction digest")
	}
}
