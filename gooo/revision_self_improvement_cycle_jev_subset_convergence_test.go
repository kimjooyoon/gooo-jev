package gooo

import "testing"

func testJEVSubsetConvergenceInput() RevisionSelfImprovementCycleJEVSubsetConvergenceInput {
	return RevisionSelfImprovementCycleJEVSubsetConvergenceInput{
		SubsetName:               "gooo.provenance.v1",
		DeclarationDigest:        digestString("declaration"),
		IRDigest:                 digestString("ir"),
		GenerationDigest:         digestString("generation"),
		ReverseObservationDigest: digestString("reverse-observation"),
		ConvergenceIteration:     3,
		Stable:                   true,
		ReverseObserved:          true,
		SourceLineCount:          120,
		GeneratedLineCount:       24,
	}
}

func TestObserveRevisionSelfImprovementCycleJEVSubsetConvergence(t *testing.T) {
	observation := ObserveRevisionSelfImprovementCycleJEVSubsetConvergence(
		testJEVSubsetConvergenceInput(),
	)
	if observation.Status != "BOUND" || observation.MissingStage != "" {
		t.Fatalf("expected bounded subset convergence, got %#v", observation)
	}
	if err := observation.Validate(); err != nil {
		t.Fatalf("expected valid subset convergence: %v", err)
	}
}

func TestObserveRevisionSelfImprovementCycleJEVSubsetConvergencePreservesLineageUnknown(
	t *testing.T,
) {
	input := testJEVSubsetConvergenceInput()
	input.IRDigest = ""
	observation := ObserveRevisionSelfImprovementCycleJEVSubsetConvergence(input)
	if observation.Status != "UNKNOWN" ||
		observation.MissingStage != "revision-self-improvement-cycle-jev-subset-convergence-lineage" {
		t.Fatalf("expected missing IR lineage to remain unknown, got %#v", observation)
	}
	if err := observation.Validate(); err != nil {
		t.Fatalf("expected valid unknown subset convergence: %v", err)
	}
}

func TestObserveRevisionSelfImprovementCycleJEVSubsetConvergenceRequiresReverseObservation(
	t *testing.T,
) {
	input := testJEVSubsetConvergenceInput()
	input.ReverseObserved = false
	observation := ObserveRevisionSelfImprovementCycleJEVSubsetConvergence(input)
	if observation.Status != "UNKNOWN" ||
		observation.MissingStage != "revision-self-improvement-cycle-jev-subset-convergence-reverse-observation" {
		t.Fatalf("expected missing reverse observation to remain unknown, got %#v", observation)
	}
}

func TestRevisionSelfImprovementCycleJEVSubsetConvergenceRejectsTampering(t *testing.T) {
	observation := ObserveRevisionSelfImprovementCycleJEVSubsetConvergence(
		testJEVSubsetConvergenceInput(),
	)
	observation.StabilityStatus = "unknown"
	if err := observation.Validate(); err == nil {
		t.Fatal("expected stability tampering to be rejected")
	}
}