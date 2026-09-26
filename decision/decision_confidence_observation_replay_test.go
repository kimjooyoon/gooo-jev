package decision

import (
	"testing"
	"time"
)

func TestReplayDecisionConfidenceObservationPreservesEvidence(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	observation := DecisionConfidenceObservation{
		Schema:                   DecisionConfidenceObservationSchemaV1,
		AssessmentDigest:         "assessment-digest",
		ReverseObservationDigest: "reverse-observation-digest",
		AssessmentStatus:         DecisionConfidenceAssessed,
		ReverseStatus:            ProvenanceObserved,
		Status:                   DecisionConfidenceObserved,
		ObservedAt:               now,
	}
	if err := observation.assignDigest(); err != nil {
		t.Fatalf("observation.assignDigest() error = %v", err)
	}
	replay, err := ReplayDecisionConfidenceObservation(observation, now.Add(time.Minute))
	if err != nil {
		t.Fatalf("ReplayDecisionConfidenceObservation() error = %v", err)
	}
	if replay.Status != DecisionConfidenceObservationReplayed || replay.ObservationDigest != observation.ObservationDigest {
		t.Fatalf("replay = %#v", replay)
	}
	if err := replay.Validate(); err != nil {
		t.Fatalf("replay Validate() error = %v", err)
	}
}

func TestReplayDecisionConfidenceObservationPreservesUnknown(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	observation := DecisionConfidenceObservation{
		Schema:         DecisionConfidenceObservationSchemaV1,
		Status:         DecisionConfidenceObservationUnknown,
		MissingStage:   "assessment:confidence-threshold",
		ObservedAt:     now,
	}
	if err := observation.assignDigest(); err != nil {
		t.Fatalf("observation.assignDigest() error = %v", err)
	}
	replay, err := ReplayDecisionConfidenceObservation(observation, now.Add(time.Minute))
	if err != nil {
		t.Fatalf("ReplayDecisionConfidenceObservation() error = %v", err)
	}
	if replay.Status != DecisionConfidenceObservationReplayUnknown || replay.MissingStage != "observation:assessment:confidence-threshold" {
		t.Fatalf("unknown replay = %#v", replay)
	}
	if err := replay.Validate(); err != nil {
		t.Fatalf("unknown replay Validate() error = %v", err)
	}
}

func TestDecisionConfidenceObservationReplayRejectsTampering(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	replay := DecisionConfidenceObservationReplay{
		Schema:                   DecisionConfidenceObservationReplaySchemaV1,
		ObservationDigest:        "observation-digest",
		AssessmentDigest:         "assessment-digest",
		ReverseObservationDigest: "reverse-observation-digest",
		AssessmentStatus:         DecisionConfidenceAssessed,
		ReverseStatus:            ProvenanceObserved,
		ObservationStatus:        DecisionConfidenceObserved,
		Status:                   DecisionConfidenceObservationReplayed,
		ObservedAt:               now,
	}
	if err := replay.assignDigest(); err != nil {
		t.Fatalf("replay.assignDigest() error = %v", err)
	}
	replay.AssessmentDigest = "tampered-assessment"
	if err := replay.Validate(); err == nil {
		t.Fatal("tampered decision confidence observation replay unexpectedly validated")
	}
}
