package gooo

import "testing"

func candidateObservationInputs(t *testing.T) (RevisionImprovementObservation, RevisionCandidateMaterialization) {
	t.Helper()
	baseline := summaryForReplacement(t, "lineage")
	candidate := summaryForReplacement(t, "lineage-expanded")
	observation, err := ObserveRevisionImprovement(baseline, candidate)
	if err != nil {
		t.Fatalf("ObserveRevisionImprovement() error = %v", err)
	}
	assessment := choiceAssessmentForRevision(t)
	qualityMetrics, err := MeasureRevision(validContract, mustApplyRevisionForCandidateObservation(t, assessment))
	if err != nil {
		t.Fatalf("MeasureRevision() error = %v", err)
	}
	quality, err := EvaluateRevisionMetrics(qualityMetrics)
	if err != nil {
		t.Fatalf("EvaluateRevisionMetrics() error = %v", err)
	}
	feedback, err := DeriveRevisionFeedback(quality)
	if err != nil {
		t.Fatalf("DeriveRevisionFeedback() error = %v", err)
	}
	proposal, err := ProposeRevisionDirection(feedback)
	if err != nil {
		t.Fatalf("ProposeRevisionDirection() error = %v", err)
	}
	materialization, err := MaterializeRevisionCandidate(feedback, proposal, assessment)
	if err != nil {
		t.Fatalf("MaterializeRevisionCandidate() error = %v", err)
	}
	return observation, materialization
}

func mustApplyRevisionForCandidateObservation(t *testing.T, assessment DecisionAssessment) RevisionApplication {
	t.Helper()
	candidate, err := ProposeRevision(assessment, RepairRevision)
	if err != nil {
		t.Fatalf("ProposeRevision() error = %v", err)
	}
	edit := SourceEdit{Start: Position{Line: 4, Column: 3}, End: Position{Line: 4, Column: 11}, Replacement: "lineage"}
	edit.Digest = digestSourceEdit(edit)
	application, err := ApplyRevision(validContract, digestString(validContract), candidate, edit)
	if err != nil {
		t.Fatalf("ApplyRevision() error = %v", err)
	}
	return application
}

func TestObserveRevisionCandidateForImprovementBindsMaterialization(t *testing.T) {
	observation, materialization := candidateObservationInputs(t)
	result, err := ObserveRevisionCandidateForImprovement(observation, materialization)
	if err != nil {
		t.Fatalf("ObserveRevisionCandidateForImprovement() error = %v", err)
	}
	if result.Status != "BOUND" || !result.CandidateAvailable || result.MaterializationSignal == "" {
		t.Fatalf("unexpected revision candidate observation: %#v", result)
	}
	if result.SourceDigest != materialization.SourceDigest || result.CandidateDigest != materialization.CandidateDigest {
		t.Fatalf("candidate observation lost materialization links: %#v", result)
	}
	if err := result.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestObserveRevisionCandidateForImprovementRetainsComparisonFailure(t *testing.T) {
	observation, materialization := candidateObservationInputs(t)
	observation.ObservationDigest = digestString("tampered")
	result, err := ObserveRevisionCandidateForImprovement(observation, materialization)
	if err == nil {
		t.Fatal("ObserveRevisionCandidateForImprovement() error = nil, want comparison failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "revision-candidate-observation-comparison" {
		t.Fatalf("unexpected unknown comparison result: %#v", result)
	}
}

func TestObserveRevisionCandidateForImprovementRetainsMaterializationFailure(t *testing.T) {
	observation, materialization := candidateObservationInputs(t)
	materialization.MaterializationDigest = digestString("tampered")
	result, err := ObserveRevisionCandidateForImprovement(observation, materialization)
	if err == nil {
		t.Fatal("ObserveRevisionCandidateForImprovement() error = nil, want materialization failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "revision-candidate-observation-materialization" {
		t.Fatalf("unexpected unknown materialization result: %#v", result)
	}
}

func TestObserveRevisionCandidateForImprovementRetainsSourceLinkFailure(t *testing.T) {
	observation, materialization := candidateObservationInputs(t)
	materialization.SourceDigest = digestString("other-source")
	materialization.MaterializationDigest = digestRevisionCandidateMaterialization(materialization)
	result, err := ObserveRevisionCandidateForImprovement(observation, materialization)
	if err == nil {
		t.Fatal("ObserveRevisionCandidateForImprovement() error = nil, want source link failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "revision-candidate-observation-link" {
		t.Fatalf("unexpected unknown source link result: %#v", result)
	}
}

func TestObserveRevisionCandidateForImprovementIsDeterministic(t *testing.T) {
	observation, materialization := candidateObservationInputs(t)
	first, err := ObserveRevisionCandidateForImprovement(observation, materialization)
	if err != nil {
		t.Fatalf("first ObserveRevisionCandidateForImprovement() error = %v", err)
	}
	second, err := ObserveRevisionCandidateForImprovement(observation, materialization)
	if err != nil {
		t.Fatalf("second ObserveRevisionCandidateForImprovement() error = %v", err)
	}
	if first.ResultDigest != second.ResultDigest {
		t.Fatal("same observation and materialization produced different digest")
	}
}