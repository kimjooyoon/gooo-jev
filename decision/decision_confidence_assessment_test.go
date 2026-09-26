package decision

import "testing"
import "time"

func TestAssessDecisionConfidenceBindsObservedResult(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	confidence := 0.9
	result := Result{
		SpecID:          "spec",
		Kind:            KindChoice,
		Value:           Value{Choice: "allow"},
		Confidence:      &confidence,
		EvidenceDigest:  "evidence",
		Provider:        "provider",
		Status:          StatusObserved,
		ObservedAt:      now,
	}
	resultDigest, err := Digest(result)
	if err != nil {
		t.Fatalf("Digest(result) error = %v", err)
	}
	receipt := Receipt{
		Schema:         SchemaV1,
		SpecDigest:     "spec-digest",
		StateDigest:    "state-digest",
		ResultDigest:   resultDigest,
		PolicyDigest:   "policy-digest",
		Provider:       "provider",
		Status:         StatusObserved,
		ObservedAt:     now,
		NonAuthorizing: true,
		DecisionDigest: "decision-digest",
	}
	assessment, err := AssessDecisionConfidence(receipt, result, 0.8, now.Add(time.Minute))
	if err != nil {
		t.Fatalf("AssessDecisionConfidence() error = %v", err)
	}
	if assessment.Status != DecisionConfidenceAssessed || assessment.Confidence != confidence {
		t.Fatalf("assessment = %#v", assessment)
	}
	if err := assessment.Validate(); err != nil {
		t.Fatalf("assessment Validate() error = %v", err)
	}
}

func TestAssessDecisionConfidenceRequiresReviewBelowThreshold(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	confidence := 0.4
	result := Result{
		SpecID:         "spec",
		Kind:            KindChoice,
		Value:          Value{Choice: "review"},
		Confidence:     &confidence,
		EvidenceDigest: "evidence",
		Provider:       "provider",
		Status:         StatusObserved,
		ObservedAt:     now,
	}
	resultDigest, err := Digest(result)
	if err != nil {
		t.Fatalf("Digest(result) error = %v", err)
	}
	receipt := Receipt{
		Schema:         SchemaV1,
		SpecDigest:     "spec-digest",
		StateDigest:    "state-digest",
		ResultDigest:   resultDigest,
		PolicyDigest:   "policy-digest",
		Provider:       "provider",
		Status:         StatusObserved,
		ObservedAt:     now,
		NonAuthorizing: true,
		DecisionDigest: "decision-digest",
	}
	assessment, err := AssessDecisionConfidence(receipt, result, 0.8, now.Add(time.Minute))
	if err != nil {
		t.Fatalf("AssessDecisionConfidence() error = %v", err)
	}
	if assessment.Status != DecisionConfidenceAssessmentReview || assessment.MissingStage != "confidence-threshold" {
		t.Fatalf("review assessment = %#v", assessment)
	}
	if err := assessment.Validate(); err != nil {
		t.Fatalf("review assessment Validate() error = %v", err)
	}
}

func TestAssessDecisionConfidenceRejectsResultBindingTampering(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	confidence := 0.9
	result := Result{Confidence: &confidence, EvidenceDigest: "evidence", Provider: "provider", Status: StatusObserved, ObservedAt: now}
	assessment, err := AssessDecisionConfidence(Receipt{Schema: SchemaV1, ResultDigest: "tampered-result", DecisionDigest: "decision", NonAuthorizing: true, ObservedAt: now}, result, 0.8, now.Add(time.Minute))
	if err != nil {
		t.Fatalf("AssessDecisionConfidence() error = %v", err)
	}
	if assessment.Status != DecisionConfidenceAssessmentUnknown || assessment.MissingStage != "receipt" {
		t.Fatalf("unknown assessment = %#v", assessment)
	}
	if err := assessment.Validate(); err != nil {
		t.Fatalf("unknown assessment Validate() error = %v", err)
	}
}
