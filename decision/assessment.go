package decision

import (
	"errors"
	"strings"
	"time"
)

type AssessmentStatus string

const (
	AssessmentObserved AssessmentStatus = "observed"
	AssessmentReview   AssessmentStatus = "review"
	AssessmentUnknown  AssessmentStatus = "unknown"
)

type DecisionAssessment struct {
	SpecDigest       string
	ResultDigest     string
	EvidenceDigest   string
	Freshness        FreshnessState
	Confidence       *float64
	Status           AssessmentStatus
	CausalReason     string
	NonAuthorizing   bool
	AssessmentDigest string
}

func AssessDecision(spec Spec, result Result, evaluatedAt time.Time, maxAge time.Duration) (DecisionAssessment, error) {
	if err := result.ValidateFor(spec); err != nil {
		return DecisionAssessment{}, err
	}
	if evaluatedAt.IsZero() {
		return DecisionAssessment{}, errors.New("decision assessment evaluation time is required")
	}
	specDigest, err := Digest(spec)
	if err != nil {
		return DecisionAssessment{}, err
	}
	resultDigest, err := Digest(result)
	if err != nil {
		return DecisionAssessment{}, err
	}
	freshness := ObserveFreshness(result, evaluatedAt, maxAge)
	assessment := DecisionAssessment{
		SpecDigest:     specDigest,
		ResultDigest:   resultDigest,
		EvidenceDigest: result.EvidenceDigest,
		Freshness:      freshness.State,
		Confidence:     result.Confidence,
		Status:         AssessmentObserved,
		CausalReason:   "DECISION_OBSERVED",
		NonAuthorizing: true,
	}
	if result.Status != StatusObserved {
		assessment.Status = AssessmentUnknown
		assessment.CausalReason = "RESULT_STATUS_" + strings.ToUpper(string(result.Status))
	} else if freshness.State == FreshnessUnknown {
		assessment.Status = AssessmentUnknown
		assessment.CausalReason = "DECISION_FRESHNESS_UNKNOWN"
	} else if freshness.State != FreshnessFresh {
		assessment.Status = AssessmentReview
		assessment.CausalReason = "DECISION_FRESHNESS_" + strings.ToUpper(string(freshness.State))
	} else if result.Confidence == nil {
		assessment.Status = AssessmentReview
		assessment.CausalReason = "DECISION_CONFIDENCE_NOT_OBSERVED"
	} else if spec.Threshold != nil && *result.Confidence < *spec.Threshold {
		assessment.Status = AssessmentReview
		assessment.CausalReason = "DECISION_CONFIDENCE_BELOW_THRESHOLD"
	}
	assessment.AssessmentDigest, err = Digest(assessment)
	if err != nil {
		return DecisionAssessment{}, err
	}
	return assessment, nil
}

func (assessment DecisionAssessment) Validate() error {
	if strings.TrimSpace(assessment.SpecDigest) == "" ||
		strings.TrimSpace(assessment.ResultDigest) == "" ||
		strings.TrimSpace(assessment.EvidenceDigest) == "" ||
		strings.TrimSpace(assessment.CausalReason) == "" ||
		strings.TrimSpace(assessment.AssessmentDigest) == "" {
		return errors.New("decision assessment is incomplete")
	}
	if !assessment.NonAuthorizing {
		return errors.New("decision assessment must remain non-authorizing")
	}
	switch assessment.Freshness {
	case FreshnessFresh, FreshnessStale, FreshnessFuture, FreshnessUnknown:
	default:
		return errors.New("unsupported decision assessment freshness")
	}
	switch assessment.Status {
	case AssessmentObserved, AssessmentReview, AssessmentUnknown:
	default:
		return errors.New("unsupported decision assessment status")
	}
	copy := assessment
	copy.AssessmentDigest = ""
	digest, err := Digest(copy)
	if err != nil {
		return err
	}
	if digest != assessment.AssessmentDigest {
		return errors.New("decision assessment digest does not match its evidence")
	}
	return nil
}
