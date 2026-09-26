package decision

import "testing"

func replayFeedbackLSPInput(status string) JEVImprovementReplayFeedback {
	return JEVImprovementReplayFeedback{
		Status:                  status,
		NonExecuting:            true,
		NonAuthorizing:          true,
		CandidateDigest:         "candidate-digest",
		ReplayObservationDigest: "replay-observation-digest",
		MetricDigest:            "metric-digest",
		FeedbackEvidenceDigest:  "feedback-evidence-digest",
		EvidenceDigest:          "feedback-evidence",
	}
}

func TestProjectExecutionEnvelopeJEVImprovementReplayFeedbackLSP(t *testing.T) {
	for _, status := range []string{jevReplayFeedbackConfirmed, jevReplayFeedbackRefuted} {
		t.Run(status, func(t *testing.T) {
			got := ProjectExecutionEnvelopeJEVImprovementReplayFeedbackLSP(ExecutionEnvelopeJEVImprovementReplayFeedbackLSPInput{
				Feedback:       replayFeedbackLSPInput(status),
				NonAuthorizing: true,
			})
			want := "confirmed"
			if status == jevReplayFeedbackRefuted {
				want = "refuted"
			}
			if got.Status != want || got.EvidenceDigest == "" || !got.NonExecuting || !got.NonAuthorizing {
				t.Fatalf("got %+v", got)
			}
		})
	}
}

func TestProjectExecutionEnvelopeJEVImprovementReplayFeedbackLSPPreservesUnknownStage(t *testing.T) {
	feedback := replayFeedbackLSPInput(jevReplayFeedbackConfirmed)
	feedback.EvidenceDigest = ""
	feedback.MissingStage = "feedback-ledger"
	got := ProjectExecutionEnvelopeJEVImprovementReplayFeedbackLSP(ExecutionEnvelopeJEVImprovementReplayFeedbackLSPInput{
		Feedback:       feedback,
		NonAuthorizing: true,
	})
	if got.Status != "UNKNOWN" || got.MissingStage != "feedback-ledger" || got.EvidenceDigest == "" {
		t.Fatalf("got %+v", got)
	}
}

func TestProjectExecutionEnvelopeJEVImprovementReplayFeedbackLSPRejectsAuthorization(t *testing.T) {
	got := ProjectExecutionEnvelopeJEVImprovementReplayFeedbackLSP(ExecutionEnvelopeJEVImprovementReplayFeedbackLSPInput{})
	if got.Status != "UNKNOWN" || got.MissingStage != "authorization-boundary" || got.NonAuthorizing || got.EvidenceDigest == "" {
		t.Fatalf("got %+v", got)
	}
}
