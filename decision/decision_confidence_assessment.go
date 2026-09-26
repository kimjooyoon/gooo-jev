package decision

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const DecisionConfidenceAssessmentSchemaV1 = "gooo/jev-decision-confidence-assessment/v1"

type DecisionConfidenceAssessmentStatus string

const (
	DecisionConfidenceAssessed        DecisionConfidenceAssessmentStatus = "assessed"
	DecisionConfidenceAssessmentReview DecisionConfidenceAssessmentStatus = "review"
	DecisionConfidenceAssessmentUnknown DecisionConfidenceAssessmentStatus = "unknown"
)

type DecisionConfidenceAssessment struct {
	Schema             string
	DecisionDigest     string
	ResultDigest       string
	Confidence         float64
	MinimumConfidence  float64
	Status             DecisionConfidenceAssessmentStatus
	MissingStage       string
	NonAuthorizing     bool
	ObservedAt         time.Time
	AssessmentDigest   string
}

func AssessDecisionConfidence(
	receipt Receipt,
	result Result,
	minimumConfidence float64,
	observedAt time.Time,
) (DecisionConfidenceAssessment, error) {
	if observedAt.IsZero() {
		return DecisionConfidenceAssessment{}, errors.New("decision confidence observation time is required")
	}
	assessment := DecisionConfidenceAssessment{
		Schema:            DecisionConfidenceAssessmentSchemaV1,
		DecisionDigest:    receipt.DecisionDigest,
		ResultDigest:      receipt.ResultDigest,
		MinimumConfidence: minimumConfidence,
		Status:            DecisionConfidenceAssessed,
		NonAuthorizing:    true,
		ObservedAt:        observedAt.UTC(),
	}
	resultDigest, digestErr := Digest(result)
	switch {
	case receipt.Validate() != nil:
		assessment.Status = DecisionConfidenceAssessmentUnknown
		assessment.MissingStage = "receipt"
	case digestErr != nil:
		return DecisionConfidenceAssessment{}, fmt.Errorf("digest decision result: %w", digestErr)
	case receipt.ResultDigest != resultDigest:
		assessment.Status = DecisionConfidenceAssessmentUnknown
		assessment.MissingStage = "result-binding"
	case result.Confidence == nil:
		assessment.Status = DecisionConfidenceAssessmentUnknown
		assessment.MissingStage = "confidence"
	case *result.Confidence < 0 || *result.Confidence > 1:
		assessment.Status = DecisionConfidenceAssessmentUnknown
		assessment.MissingStage = "confidence-range"
	case minimumConfidence < 0 || minimumConfidence > 1:
		assessment.Status = DecisionConfidenceAssessmentUnknown
		assessment.MissingStage = "threshold-range"
	default:
		assessment.Confidence = *result.Confidence
		if result.Status != StatusObserved {
			assessment.Status = DecisionConfidenceAssessmentReview
			assessment.MissingStage = "decision-status"
		} else if assessment.Confidence < minimumConfidence {
			assessment.Status = DecisionConfidenceAssessmentReview
			assessment.MissingStage = "confidence-threshold"
		}
	}
	if err := assessment.assignDigest(); err != nil {
		return DecisionConfidenceAssessment{}, fmt.Errorf("digest decision confidence assessment: %w", err)
	}
	if err := assessment.Validate(); err != nil {
		return DecisionConfidenceAssessment{}, err
	}
	return assessment, nil
}

func (assessment *DecisionConfidenceAssessment) assignDigest() error {
	digest, err := assessment.computeDigest()
	if err != nil {
		return err
	}
	assessment.AssessmentDigest = digest
	return nil
}

func (assessment DecisionConfidenceAssessment) Validate() error {
	if err := assessment.validateShape(); err != nil {
		return err
	}
	expected, err := assessment.computeDigest()
	if err != nil {
		return fmt.Errorf("digest decision confidence assessment: %w", err)
	}
	if assessment.AssessmentDigest != expected {
		return errors.New("decision confidence assessment digest mismatch")
	}
	return nil
}

func (assessment DecisionConfidenceAssessment) validateShape() error {
	if assessment.Schema != DecisionConfidenceAssessmentSchemaV1 {
		return errors.New("unsupported decision confidence assessment schema")
	}
	if assessment.ObservedAt.IsZero() {
		return errors.New("decision confidence observation time is required")
	}
	if strings.TrimSpace(assessment.AssessmentDigest) == "" {
		return errors.New("decision confidence assessment digest is required")
	}
	if !assessment.NonAuthorizing {
		return errors.New("decision confidence assessment must remain non-authorizing")
	}
	switch assessment.Status {
	case DecisionConfidenceAssessed:
		if strings.TrimSpace(assessment.DecisionDigest) == "" ||
			strings.TrimSpace(assessment.ResultDigest) == "" ||
			assessment.Confidence < 0 || assessment.Confidence > 1 ||
			assessment.MinimumConfidence < 0 || assessment.MinimumConfidence > 1 ||
			strings.TrimSpace(assessment.MissingStage) != "" {
			return errors.New("assessed decision confidence is incomplete")
		}
	case DecisionConfidenceAssessmentReview:
		if strings.TrimSpace(assessment.DecisionDigest) == "" ||
			strings.TrimSpace(assessment.ResultDigest) == "" ||
			assessment.Confidence < 0 || assessment.Confidence > 1 ||
			assessment.MinimumConfidence < 0 || assessment.MinimumConfidence > 1 ||
			strings.TrimSpace(assessment.MissingStage) == "" {
			return errors.New("review decision confidence is incomplete")
		}
	case DecisionConfidenceAssessmentUnknown:
		if strings.TrimSpace(assessment.MissingStage) == "" {
			return errors.New("unknown decision confidence requires a missing stage")
		}
	default:
		return fmt.Errorf("unsupported decision confidence status %q", assessment.Status)
	}
	return nil
}

func (assessment DecisionConfidenceAssessment) computeDigest() (string, error) {
	return Digest(struct {
		Schema            string
		DecisionDigest    string
		ResultDigest      string
		Confidence        float64
		MinimumConfidence float64
		Status            DecisionConfidenceAssessmentStatus
		MissingStage      string
		NonAuthorizing    bool
		ObservedAt        time.Time
	}{
		Schema:            assessment.Schema,
		DecisionDigest:    assessment.DecisionDigest,
		ResultDigest:      assessment.ResultDigest,
		Confidence:        assessment.Confidence,
		MinimumConfidence: assessment.MinimumConfidence,
		Status:            assessment.Status,
		MissingStage:      assessment.MissingStage,
		NonAuthorizing:    assessment.NonAuthorizing,
		ObservedAt:        assessment.ObservedAt.UTC(),
	})
}
