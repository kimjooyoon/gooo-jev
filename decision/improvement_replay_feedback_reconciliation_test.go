package decision

import (
	"testing"
	"time"
)

func makeReplayFeedbackHistory(t *testing.T, observations ...FeedbackObservation) FeedbackHistory {
	t.Helper()
	summary, err := SummarizeFeedback(observations, "metric")
	if err != nil {
		t.Fatalf("summarize feedback: %v", err)
	}
	history, err := BuildFeedbackHistory([]FeedbackSummary{summary}, time.Unix(1_800_000_100, 0).UTC())
	if err != nil {
		t.Fatalf("build feedback history: %v", err)
	}
	return history
}

func makeReplayFeedbackBinding(t *testing.T, kind FeedbackKind) ImprovementReplayFeedbackBinding {
	t.Helper()
	binding := BindImprovementReplayFeedback(replayFeedbackInput(makeImprovementReplayFeedback(t, kind)))
	if binding.Status != "feedback-recorded" {
		t.Fatalf("unexpected replay feedback binding: %+v", binding)
	}
	return binding
}

func TestReconcileImprovementReplayFeedbackClassifiesConfirmedAndRefuted(t *testing.T) {
	cases := []struct {
		name   string
		kind   FeedbackKind
		status string
	}{
		{name: "confirmed", kind: FeedbackConfirmed, status: "confirmed"},
		{name: "refuted", kind: FeedbackRefuted, status: "refuted"},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			feedback := makeImprovementReplayFeedback(t, testCase.kind)
			output := ReconcileImprovementReplayFeedback(ImprovementReplayFeedbackReconciliationInput{
				ReplayFeedback:  makeReplayFeedbackBinding(t, testCase.kind),
				FeedbackHistory: makeReplayFeedbackHistory(t, feedback),
				NonAuthorizing:  true,
			})
			if output.Status != testCase.status || output.MissingStage != "" {
				t.Fatalf("unexpected reconciliation: %+v", output)
			}
			if output.EvidenceDigest == "" || !output.NonExecuting || !output.NonAuthorizing {
				t.Fatalf("reconciliation provenance or safety missing: %+v", output)
			}
		})
	}
}

func TestReconcileImprovementReplayFeedbackKeepsUnknownAndTieInReview(t *testing.T) {
	unknown := makeImprovementReplayFeedback(t, FeedbackUnknown)
	output := ReconcileImprovementReplayFeedback(ImprovementReplayFeedbackReconciliationInput{
		ReplayFeedback:  makeReplayFeedbackBinding(t, FeedbackConfirmed),
		FeedbackHistory: makeReplayFeedbackHistory(t, unknown),
		NonAuthorizing:  true,
	})
	if output.Status != "review" || output.MissingStage != "feedback-unknown" {
		t.Fatalf("unexpected unknown reconciliation: %+v", output)
	}

	confirmed := makeImprovementReplayFeedback(t, FeedbackConfirmed)
	refuted := makeImprovementReplayFeedback(t, FeedbackRefuted)
	output = ReconcileImprovementReplayFeedback(ImprovementReplayFeedbackReconciliationInput{
		ReplayFeedback:  makeReplayFeedbackBinding(t, FeedbackConfirmed),
		FeedbackHistory: makeReplayFeedbackHistory(t, confirmed, refuted),
		NonAuthorizing:  true,
	})
	if output.Status != "review" || output.MissingStage != "feedback-tie" {
		t.Fatalf("unexpected tie reconciliation: %+v", output)
	}
}

func TestReconcileImprovementReplayFeedbackFailsClosed(t *testing.T) {
	feedback := makeImprovementReplayFeedback(t, FeedbackConfirmed)
	history := makeReplayFeedbackHistory(t, feedback)
	history.HistoryDigest = "tampered"
	output := ReconcileImprovementReplayFeedback(ImprovementReplayFeedbackReconciliationInput{
		ReplayFeedback:  makeReplayFeedbackBinding(t, FeedbackConfirmed),
		FeedbackHistory: history,
		NonAuthorizing:  true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "feedback-history" {
		t.Fatalf("unexpected tampered history output: %+v", output)
	}

	input := ImprovementReplayFeedbackReconciliationInput{
		ReplayFeedback:  makeReplayFeedbackBinding(t, FeedbackConfirmed),
		FeedbackHistory: makeReplayFeedbackHistory(t, feedback),
		NonAuthorizing:  false,
	}
	output = ReconcileImprovementReplayFeedback(input)
	if output.Status != "UNKNOWN" || output.NonAuthorizing || output.MissingStage != "authorization-boundary" {
		t.Fatalf("unexpected authorization output: %+v", output)
	}
}
