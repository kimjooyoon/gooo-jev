package decision

import (
	"testing"
	"time"
)

func makeImprovementReplayFeedback(t *testing.T, kind FeedbackKind) FeedbackObservation {
	t.Helper()
	feedback := FeedbackObservation{
		ReceiptDigest:  "receipt-digest",
		LedgerDigest:   "ledger-digest",
		Kind:           kind,
		MetricName:     "metric",
		MetricValue:    0.75,
		EvidenceDigest: "feedback-evidence",
		RecordedAt:     time.Unix(1_800_000_000, 0).UTC(),
		NonAuthorizing: true,
	}
	digest, err := Digest(feedback)
	if err != nil {
		t.Fatalf("digest feedback: %v", err)
	}
	feedback.FeedbackDigest = digest
	if err := feedback.Validate(); err != nil {
		t.Fatalf("validate feedback: %v", err)
	}
	return feedback
}

func replayFeedbackInput(feedback FeedbackObservation) ImprovementReplayFeedbackBindingInput {
	return ImprovementReplayFeedbackBindingInput{
		ReplayBindingStatus: "replayed",
		CycleDigest:         "cycle-digest",
		ReplayDigest:        "replay-digest",
		CandidateDigest:     "candidate-digest",
		Feedback:            feedback,
		NonAuthorizing:      true,
	}
}

func TestBindImprovementReplayFeedbackRecordsConfirmedAndRefuted(t *testing.T) {
	for _, kind := range []FeedbackKind{FeedbackConfirmed, FeedbackRefuted} {
		output := BindImprovementReplayFeedback(replayFeedbackInput(makeImprovementReplayFeedback(t, kind)))
		if output.Status != "feedback-recorded" || output.FeedbackKind != kind {
			t.Fatalf("kind %q: unexpected output: %+v", kind, output)
		}
		if output.MetricName != "metric" || output.FeedbackDigest == "" || output.EvidenceDigest == "" {
			t.Fatalf("kind %q: feedback provenance missing: %+v", kind, output)
		}
		if !output.NonExecuting || !output.NonAuthorizing {
			t.Fatalf("kind %q: safety boundary changed: %+v", kind, output)
		}
	}
}

func TestBindImprovementReplayFeedbackKeepsUnknownInReview(t *testing.T) {
	output := BindImprovementReplayFeedback(replayFeedbackInput(makeImprovementReplayFeedback(t, FeedbackUnknown)))
	if output.Status != "review" || output.MissingStage != "feedback-unknown" {
		t.Fatalf("unexpected unknown feedback output: %+v", output)
	}
}

func TestBindImprovementReplayFeedbackDoesNotPromoteDivergence(t *testing.T) {
	input := replayFeedbackInput(makeImprovementReplayFeedback(t, FeedbackConfirmed))
	input.ReplayBindingStatus = "diverged"
	output := BindImprovementReplayFeedback(input)
	if output.Status != "review" || output.MissingStage != "replay-diverged" {
		t.Fatalf("unexpected diverged output: %+v", output)
	}
}

func TestBindImprovementReplayFeedbackRejectsTamperingAndAuthorization(t *testing.T) {
	feedback := makeImprovementReplayFeedback(t, FeedbackConfirmed)
	feedback.FeedbackDigest = "tampered"
	output := BindImprovementReplayFeedback(replayFeedbackInput(feedback))
	if output.Status != "UNKNOWN" || output.MissingStage != "feedback-observation" {
		t.Fatalf("unexpected tampered output: %+v", output)
	}

	input := replayFeedbackInput(makeImprovementReplayFeedback(t, FeedbackConfirmed))
	input.NonAuthorizing = false
	output = BindImprovementReplayFeedback(input)
	if output.Status != "UNKNOWN" || output.NonAuthorizing || output.MissingStage != "authorization-boundary" {
		t.Fatalf("unexpected authorization output: %+v", output)
	}
}
