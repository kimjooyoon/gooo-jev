package decision

import (
	"testing"
	"time"
)

func TestObserveDecisionConfidenceBindsReverseObservation(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	assessment := DecisionConfidenceAssessment{
		Schema:            DecisionConfidenceAssessmentSchemaV1,
		DecisionDigest:    "decision-digest",
		ResultDigest:      "result-digest",
		Confidence:        0.9,
		MinimumConfidence: 0.8,
		Status:            DecisionConfidenceAssessed,
		NonAuthorizing:    true,
		ObservedAt:        now,
	}
	if err := assessment.assignDigest(); err != nil {
		t.Fatalf("assessment.assignDigest() error = %v", err)
	}
	reverse := ReverseObservation{
		Schema:                 ReverseObservationSchemaV1,
		ExecutionReceiptDigest: "execution-receipt-digest",
		ObservedOutputDigest:   "output-digest",
		VerifierDigest:         "verifier-digest",
		Status:                 ProvenanceObserved,
		ObservedAt:             now,
	}
	if err := func() error {
		digest, err := reverse.computeDigest()
		reverse.ObservationDigest = digest
		return err
	}(); err != nil {
		t.Fatalf("reverse digest error = %v", err)
	}
	observation, err := ObserveDecisionConfidence(assessment, reverse, now.Add(time.Minute))
	if err != nil {
		t.Fatalf("ObserveDecisionConfidence() error = %v", err)
	}
	if observation.Status != DecisionConfidenceObserved || observation.MissingStage != "" {
		t.Fatalf("observation = %#v", observation)
	}
	if err := observation.Validate(); err != nil {
		t.Fatalf("observation Validate() error = %v", err)
	}
}

func TestObserveDecisionConfidencePreservesReview(t *testing.T) {
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
	observation, err := ObserveDecisionConfidence(assessment, ReverseObservation{}, now.Add(time.Minute))
	if err != nil {
		t.Fatalf("ObserveDecisionConfidence() error = %v", err)
	}
	if observation.Status != DecisionConfidenceObservationUnknown || observation.MissingStage != "assessment:confidence-threshold" {
		t.Fatalf("unknown observation = %#v", observation)
	}
	if err := observation.Validate(); err != nil {
		t.Fatalf("unknown observation Validate() error = %v", err)
	}
}

func TestDecisionConfidenceObservationRejectsTampering(t *testing.T) {
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
	observation.ReverseObservationDigest = "tampered-reverse"
	if err := observation.Validate(); err == nil {
		t.Fatal("tampered decision confidence observation unexpectedly validated")
	}
}
