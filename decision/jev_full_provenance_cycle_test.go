package decision

import "testing"

func fullProvenanceCycleBinding() ExecutionEnvelopeFullProvenanceSourceBinding {
	return ExecutionEnvelopeFullProvenanceSourceBinding{
		Status: "complete",
		DeclarationID: "gooo://gooo-jev/declaration/cycle",
		ContractID: "gooo://gooo-jev/contract/cycle",
		DeclarationDigest: "declaration-digest",
		IRDigest: "ir-digest",
		GenerationDigest: "generation-digest",
		BindingDigest: "binding-digest",
		ReverseObservationSourceDigest: "reverse-source-digest",
		ReverseObservationDigest: "reverse-observation-digest",
		MetricDigest: "metric-digest",
		EvidenceDigest: "evidence-digest",
		CompletenessDigest: "completeness-digest",
		NonExecuting: true,
		NonAuthorizing: true,
	}
}

func TestObserveJEVImprovementCycleFromFullProvenance(t *testing.T) {
	for _, status := range []string{"stable-for-review", "needs-revision", "hold"} {
		t.Run(status, func(t *testing.T) {
			got := ObserveJEVImprovementCycleFromFullProvenance(ExecutionEnvelopeJEVFullProvenanceCycleInput{
				Provenance: fullProvenanceCycleBinding(),
				ChangePlanDigest: "change-plan-digest",
				Feedback: feedbackCycleAggregation(status),
				NonAuthorizing: true,
			})
			if got.Status != status || got.DispositionStatus != status || got.EvidenceDigest == "" || !got.NonExecuting || !got.NonAuthorizing {
				t.Fatalf("got %+v", got)
			}
			if err := got.Validate(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestObserveJEVImprovementCycleFromFullProvenancePreservesMissingStage(t *testing.T) {
	provenance := fullProvenanceCycleBinding()
	provenance.MetricDigest = ""
	got := ObserveJEVImprovementCycleFromFullProvenance(ExecutionEnvelopeJEVFullProvenanceCycleInput{
		Provenance: provenance,
		ChangePlanDigest: "change-plan-digest",
		Feedback: feedbackCycleAggregation("stable-for-review"),
		NonAuthorizing: true,
	})
	if got.Status != "UNKNOWN" || got.MissingStage != "metric" || got.EvidenceDigest != "" {
		t.Fatalf("got %+v", got)
	}
}

func TestObserveJEVImprovementCycleFromFullProvenancePreservesIncompleteProvenance(t *testing.T) {
	provenance := fullProvenanceCycleBinding()
	provenance.Status = "UNKNOWN"
	provenance.MissingStage = "reverse-observation"
	got := ObserveJEVImprovementCycleFromFullProvenance(ExecutionEnvelopeJEVFullProvenanceCycleInput{
		Provenance: provenance,
		ChangePlanDigest: "change-plan-digest",
		Feedback: feedbackCycleAggregation("stable-for-review"),
		NonAuthorizing: true,
	})
	if got.Status != "UNKNOWN" || got.MissingStage != "reverse-observation" || got.EvidenceDigest != "" {
		t.Fatalf("got %+v", got)
	}
}
