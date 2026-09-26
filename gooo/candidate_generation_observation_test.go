package gooo

import "testing"

func candidateMaterializationForGeneration(t *testing.T) RevisionCandidateMaterialization {
	t.Helper()
	feedback := revisionFeedbackForQuality(t)
	proposal, err := ProposeRevisionDirection(feedback)
	if err != nil {
		t.Fatalf("ProposeRevisionDirection() error = %v", err)
	}
	materialization, err := MaterializeRevisionCandidate(feedback, proposal, choiceAssessmentForRevision(t))
	if err != nil {
		t.Fatalf("MaterializeRevisionCandidate() error = %v", err)
	}
	return materialization
}

func TestObserveCandidateGenerationBindsReverseObservation(t *testing.T) {
	document, err := Parse(validContract)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	observation, err := ObserveCandidateGeneration(document, candidateMaterializationForGeneration(t))
	if err != nil {
		t.Fatalf("ObserveCandidateGeneration() error = %v", err)
	}
	if observation.Status != "BOUND" || !observation.StructureMatch || !observation.ReverseObserved {
		t.Fatalf("unexpected candidate generation observation: %#v", observation)
	}
	if observation.SourceDigest != document.SourceDigest || observation.CandidateDigest == "" {
		t.Fatalf("candidate generation provenance missing: %#v", observation)
	}
	if err := observation.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestObserveCandidateGenerationRetainsUnknownMaterializationStage(t *testing.T) {
	document, err := Parse(validContract)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	materialization := candidateMaterializationForGeneration(t)
	materialization.MaterializationDigest = digestString("tampered")
	observation, err := ObserveCandidateGeneration(document, materialization)
	if err == nil {
		t.Fatal("ObserveCandidateGeneration() error = nil, want materialization failure")
	}
	if observation.Status != "UNKNOWN" || observation.MissingStage != "candidate-generation-materialization" {
		t.Fatalf("unexpected unknown observation: %#v", observation)
	}
}

func TestObserveCandidateGenerationIsDeterministic(t *testing.T) {
	document, err := Parse(validContract)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	materialization := candidateMaterializationForGeneration(t)
	first, err := ObserveCandidateGeneration(document, materialization)
	if err != nil {
		t.Fatalf("first ObserveCandidateGeneration() error = %v", err)
	}
	second, err := ObserveCandidateGeneration(document, materialization)
	if err != nil {
		t.Fatalf("second ObserveCandidateGeneration() error = %v", err)
	}
	if first.ObservationDigest != second.ObservationDigest {
		t.Fatal("same document and candidate produced different observation digest")
	}
}
