package decision

import (
	"testing"
	"time"
)

func TestReplayDecisionConfidenceAssessmentPreservesEvidence(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	assessment := DecisionConfidenceAssessment{
		Schema:             DecisionConfidenceAssessmentSchemaV1,
		DecisionDigest:     "decision-digest",
		ResultDigest:       "result-digest",
		Confidence:         0.9,
		MinimumConfidence:  0.8,
		Status:             DecisionConfidenceAssessed,
		NonAuthorizing:     true,
		ObservedAt:         now,
	}
	if err := assessment.assignDigest(); err != nil {
		t.Fatalf("assessment.assignDigest() error = %v", err)
	}
	replay, err := ReplayDecisionConfidenceAssessment(assessment, now.Add(time.Minute))
	if err != nil {
		t.Fatalf("ReplayDecisionConfidenceAssessment() error = %v", err)
	}
	if replay.Status != DecisionConfidenceAssessmentReplayed || replay.Confidence != assessment.Confidence || replay.MinimumConfidence != assessment.MinimumConfidence {
		t.Fatalf("replay = %#v", replay)
	}
	if err := replay.Validate(); err != nil {
		t.Fatalf("replay Validate() error = %v", err)
	}
}

func TestReplayDecisionConfidenceAssessmentPreservesUnknown(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	assessment := DecisionConfidenceAssessment{
		Schema:          DecisionConfidenceAssessmentSchemaV1,
		Status:          DecisionConfidenceAssessmentUnknown,
		MissingStage:    "confidence",
		NonAuthorizing:  true,
		ObservedAt:      now,
	}
	if err := assessment.assignDigest(); err != nil {
		t.Fatalf("assessment.assignDigest() error = %v", err)
	}
	replay, err := ReplayDecisionConfidenceAssessment(assessment, now.Add(time.Minute))
	if err != nil {
		t.Fatalf("ReplayDecisionConfidenceAssessment() error = %v", err)
	}
	if replay.Status != DecisionConfidenceAssessmentReplayUnknown || replay.MissingStage != "assessment:confidence" {
		t.Fatalf("unknown replay = %#v", replay)
	}
	if err := replay.Validate(); err != nil {
		t.Fatalf("unknown replay Validate() error = %v", err)
	}
}

func TestDecisionConfidenceReplayRejectsTampering(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	replay := DecisionConfidenceAssessmentReplay{
		Schema:             DecisionConfidenceAssessmentReplaySchemaV1,
		AssessmentDigest:   "assessment-digest",
		DecisionDigest:     "decision-digest",
		ResultDigest:       "result-digest",
		Confidence:         0.9,
		MinimumConfidence:  0.8,
		AssessmentStatus:   DecisionConfidenceAssessed,
		Status:             DecisionConfidenceAssessmentReplayed,
		ObservedAt:         now,
	}
	if err := replay.assignDigest(); err != nil {
		t.Fatalf("replay.assignDigest() error = %v", err)
	}
	replay.ResultDigest = "tampered-result"
	if err := replay.Validate(); err == nil {
		t.Fatal("tampered decision confidence replay unexpectedly validated")
	}
}
