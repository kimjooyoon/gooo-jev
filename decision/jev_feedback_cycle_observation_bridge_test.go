package decision

import "testing"

func feedbackCycleAggregation(status string) JEVImprovementFeedbackAggregation {
	aggregation := JEVImprovementFeedbackAggregation{Status: status, Total: 1, InputEvidenceDigest: "feedback-input-digest", NonExecuting: true, NonAuthorizing: true}
	switch status {
	case "stable-for-review":
		aggregation.Confirmed = 1
	case "needs-revision":
		aggregation.Refuted = 1
	case "hold":
		aggregation.Unknown = 1
	}
	aggregation.EvidenceDigest = digestJEVImprovementFeedbackAggregation(aggregation.Status, aggregation.Total, aggregation.Confirmed, aggregation.Refuted, aggregation.Unknown, aggregation.InputEvidenceDigest)
	return aggregation
}

func feedbackCycleInput(status string) ExecutionEnvelopeJEVImprovementFeedbackCycleInput {
	return ExecutionEnvelopeJEVImprovementFeedbackCycleInput{
		DeclarationDigest:        "declaration-digest",
		IRDigest:                 "ir-digest",
		GenerationDigest:         "generation-digest",
		ReverseObservationDigest: "reverse-observation-digest",
		MetricDigest:             "metric-digest",
		ChangePlanDigest:         "change-plan-digest",
		Feedback:                 feedbackCycleAggregation(status),
		NonAuthorizing:           true,
	}
}

func TestObserveJEVImprovementFeedbackCycle(t *testing.T) {
	for _, status := range []string{"stable-for-review", "needs-revision", "hold"} {
		t.Run(status, func(t *testing.T) {
			got := ObserveJEVImprovementFeedbackCycle(feedbackCycleInput(status))
			if got.Status != status || got.DispositionStatus != status || got.EvidenceDigest == "" || got.ChangePlanDigest == "change-plan-digest" {
				t.Fatalf("got %+v", got)
			}
			if err := got.Validate(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestObserveJEVImprovementFeedbackCyclePreservesMissingStage(t *testing.T) {
	input := feedbackCycleInput("stable-for-review")
	input.ReverseObservationDigest = ""
	got := ObserveJEVImprovementFeedbackCycle(input)
	if got.Status != "UNKNOWN" || got.MissingStage != "reverse-observation" || got.EvidenceDigest != "" {
		t.Fatalf("got %+v", got)
	}
}

func TestObserveJEVImprovementFeedbackCyclePreservesFeedbackStage(t *testing.T) {
	input := feedbackCycleInput("stable-for-review")
	input.Feedback.Status = "UNKNOWN"
	input.Feedback.MissingStage = "ledger[1]-predecessor"
	input.Feedback.EvidenceDigest = ""
	got := ObserveJEVImprovementFeedbackCycle(input)
	if got.Status != "UNKNOWN" || got.MissingStage != "ledger[1]-predecessor" || got.EvidenceDigest != "" {
		t.Fatalf("got %+v", got)
	}
}

func TestObserveJEVImprovementFeedbackCycleRejectsAuthorization(t *testing.T) {
	input := feedbackCycleInput("stable-for-review")
	input.NonAuthorizing = false
	got := ObserveJEVImprovementFeedbackCycle(input)
	if got.Status != "UNKNOWN" || got.MissingStage != "authorization-boundary" || got.NonAuthorizing || got.EvidenceDigest != "" {
		t.Fatalf("got %+v", got)
	}
}
