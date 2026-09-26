package decision

import "testing"
import "time"

func TestDeriveDecisionConfidenceImprovementSignalCreatesCandidate(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	assessment := DecisionConfidenceAssessment{
		Schema:            DecisionConfidenceAssessmentSchemaV1,
		DecisionDigest:    "decision-digest",
		ResultDigest:      "result-digest",
		Confidence:        0.4,
		MinimumConfidence: 0.8,
		Status:            DecisionConfidenceAssessmentReview,
		MissingStage:      "confidence-threshold",
		NonAuthorizing:    true,
		ObservedAt:        now,
	}
	if err := assessment.assignDigest(); err != nil {
		t.Fatalf("assessment.assignDigest() error = %v", err)
	}
	observation := DecisionConfidenceObservation{Schema: DecisionConfidenceObservationSchemaV1, Status: DecisionConfidenceObservationUnknown, MissingStage: "assessment:confidence-threshold", ObservedAt: now}
	if err := observation.assignDigest(); err != nil {
		t.Fatalf("observation.assignDigest() error = %v", err)
	}
	replay := DecisionConfidenceObservationReplay{Schema: DecisionConfidenceObservationReplaySchemaV1, Status: DecisionConfidenceObservationReplayUnknown, MissingStage: "observation:assessment:confidence-threshold", ObservedAt: now}
	if err := replay.assignDigest(); err != nil {
		t.Fatalf("replay.assignDigest() error = %v", err)
	}
	signal, err := DeriveDecisionConfidenceImprovementSignal(assessment, observation, replay, now.Add(time.Minute))
	if err != nil {
		t.Fatalf("DeriveDecisionConfidenceImprovementSignal() error = %v", err)
	}
	if signal.Status != DecisionConfidenceCandidate || signal.Action != "request-additional-confidence-evidence" {
		t.Fatalf("signal = %#v", signal)
	}
	if err := signal.Validate(); err != nil {
		t.Fatalf("signal Validate() error = %v", err)
	}
}

func TestDeriveDecisionConfidenceImprovementSignalPreservesNoChange(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	assessment := DecisionConfidenceAssessment{Schema: DecisionConfidenceAssessmentSchemaV1, DecisionDigest: "decision", ResultDigest: "result", Confidence: 0.9, MinimumConfidence: 0.8, Status: DecisionConfidenceAssessed, NonAuthorizing: true, ObservedAt: now}
	if err := assessment.assignDigest(); err != nil {
		t.Fatalf("assessment.assignDigest() error = %v", err)
	}
	observation := DecisionConfidenceObservation{Schema: DecisionConfidenceObservationSchemaV1, AssessmentDigest: assessment.AssessmentDigest, ReverseObservationDigest: "reverse", AssessmentStatus: DecisionConfidenceAssessed, ReverseStatus: ProvenanceObserved, Status: DecisionConfidenceObserved, ObservedAt: now}
	if err := observation.assignDigest(); err != nil {
		t.Fatalf("observation.assignDigest() error = %v", err)
	}
	replay := DecisionConfidenceObservationReplay{Schema: DecisionConfidenceObservationReplaySchemaV1, ObservationDigest: observation.ObservationDigest, AssessmentDigest: assessment.AssessmentDigest, ReverseObservationDigest: observation.ReverseObservationDigest, AssessmentStatus: DecisionConfidenceAssessed, ReverseStatus: ProvenanceObserved, ObservationStatus: DecisionConfidenceObserved, Status: DecisionConfidenceObservationReplayed, ObservedAt: now}
	if err := replay.assignDigest(); err != nil {
		t.Fatalf("replay.assignDigest() error = %v", err)
	}
	signal, err := DeriveDecisionConfidenceImprovementSignal(assessment, observation, replay, now.Add(time.Minute))
	if err != nil {
		t.Fatalf("DeriveDecisionConfidenceImprovementSignal() error = %v", err)
	}
	if signal.Status != DecisionConfidenceNoChange || signal.Action != "retain-observed-provenance" {
		t.Fatalf("no-change signal = %#v", signal)
	}
	if err := signal.Validate(); err != nil {
		t.Fatalf("signal Validate() error = %v", err)
	}
}

func TestDecisionConfidenceImprovementSignalRejectsTampering(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	signal := DecisionConfidenceImprovementSignal{Schema: DecisionConfidenceImprovementSignalSchemaV1, AssessmentDigest: "assessment", ObservationDigest: "observation", ReplayDigest: "replay", Status: DecisionConfidenceNoChange, Action: "retain-observed-provenance", NonAuthorizing: true, ObservedAt: now}
	if err := signal.assignDigest(); err != nil {
		t.Fatalf("signal.assignDigest() error = %v", err)
	}
	signal.Action = "mutate-without-review"
	if err := signal.Validate(); err == nil {
		t.Fatal("tampered decision confidence improvement signal unexpectedly validated")
	}
}
