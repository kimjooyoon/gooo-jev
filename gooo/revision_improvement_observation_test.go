package gooo

import "testing"

func summaryForReplacement(t *testing.T, replacement string) RevisionEvidenceSummary {
	t.Helper()
	chain, metrics := revisionTrendChainForReplacement(t, replacement)
	document, err := Parse(validContract)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	generation, err := Generate(document)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	binding, err := BindRevisionEvidenceGeneration(chain, generation)
	if err != nil {
		t.Fatalf("BindRevisionEvidenceGeneration() error = %v", err)
	}
	summary, err := SummarizeRevisionEvidence(chain, binding, metrics)
	if err != nil {
		t.Fatalf("SummarizeRevisionEvidence() error = %v", err)
	}
	return summary
}

func TestObserveRevisionImprovementReportsObservedDifferences(t *testing.T) {
	baseline := summaryForReplacement(t, "lineage")
	candidate := summaryForReplacement(t, "lineage-expanded")
	observation, err := ObserveRevisionImprovement(baseline, candidate)
	if err != nil {
		t.Fatalf("ObserveRevisionImprovement() error = %v", err)
	}
	if observation.Status != "BOUND" || observation.ObservationClass == "" || observation.ReviewSignal == "" {
		t.Fatalf("unexpected revision improvement observation: %#v", observation)
	}
	if observation.ProposedIRChanged != (baseline.ProposedIRDigest != candidate.ProposedIRDigest) ||
		observation.GeneratedIRChanged != (baseline.GeneratedIRDigest != candidate.GeneratedIRDigest) ||
		observation.StructureStable != (baseline.StructureDigest == candidate.StructureDigest) {
		t.Fatalf("observation did not preserve exact differences: %#v", observation)
	}
	if err := observation.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestObserveRevisionImprovementReportsStableEvidence(t *testing.T) {
	baseline := summaryForReplacement(t, "lineage")
	observation, err := ObserveRevisionImprovement(baseline, baseline)
	if err != nil {
		t.Fatalf("ObserveRevisionImprovement() error = %v", err)
	}
	if observation.ObservationClass != "stable" || observation.ReviewSignal != "observe" || observation.StageCountDelta != 0 {
		t.Fatalf("unexpected stable observation: %#v", observation)
	}
}

func TestObserveRevisionImprovementRetainsBaselineFailure(t *testing.T) {
	baseline := summaryForReplacement(t, "lineage")
	candidate := summaryForReplacement(t, "lineage-expanded")
	baseline.SummaryDigest = digestString("tampered")
	observation, err := ObserveRevisionImprovement(baseline, candidate)
	if err == nil {
		t.Fatal("ObserveRevisionImprovement() error = nil, want baseline failure")
	}
	if observation.Status != "UNKNOWN" || observation.MissingStage != "revision-improvement-observation-baseline" {
		t.Fatalf("unexpected unknown baseline observation: %#v", observation)
	}
}

func TestObserveRevisionImprovementIsDeterministic(t *testing.T) {
	baseline := summaryForReplacement(t, "lineage")
	candidate := summaryForReplacement(t, "lineage-expanded")
	first, err := ObserveRevisionImprovement(baseline, candidate)
	if err != nil {
		t.Fatalf("first ObserveRevisionImprovement() error = %v", err)
	}
	second, err := ObserveRevisionImprovement(baseline, candidate)
	if err != nil {
		t.Fatalf("second ObserveRevisionImprovement() error = %v", err)
	}
	if first.ObservationDigest != second.ObservationDigest {
		t.Fatal("same summary pair produced different observation digest")
	}
}