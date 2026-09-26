package decision

import (
	"errors"
	"strings"
)

type CausalReviewStatus string

const (
	CausalReviewReview  CausalReviewStatus = "review"
	CausalReviewUnknown CausalReviewStatus = "unknown"
)

type CausalReview struct {
	AssessmentDigest       string
	CounterexampleSetDigest string
	FeedbackHistoryDigest  string
	Status                 CausalReviewStatus
	CausalReason           string
	NonAuthorizing         bool
	ReviewDigest           string
}

func BuildCausalReview(
	assessment DecisionAssessment,
	counterexamples CounterexampleSet,
	history FeedbackHistory,
) (CausalReview, error) {
	if err := assessment.Validate(); err != nil {
		return CausalReview{}, err
	}
	if err := counterexamples.Validate(); err != nil {
		return CausalReview{}, err
	}
	if err := history.Validate(); err != nil {
		return CausalReview{}, err
	}
	if counterexamples.AssessmentDigest != assessment.AssessmentDigest {
		return CausalReview{}, errors.New("causal review assessment digests do not match")
	}
	review := CausalReview{
		AssessmentDigest:        assessment.AssessmentDigest,
		CounterexampleSetDigest: counterexamples.SetDigest,
		FeedbackHistoryDigest:   history.HistoryDigest,
		Status:                  CausalReviewReview,
		CausalReason:            "COUNTEREXAMPLES_PRESENT",
		NonAuthorizing:          true,
	}
	if assessment.Status == AssessmentUnknown {
		review.Status = CausalReviewUnknown
		review.CausalReason = "ASSESSMENT_UNKNOWN:" + assessment.CausalReason
	}
	var err error
	review.ReviewDigest, err = Digest(review)
	if err != nil {
		return CausalReview{}, err
	}
	return review, nil
}

func (review CausalReview) Validate() error {
	if strings.TrimSpace(review.AssessmentDigest) == "" ||
		strings.TrimSpace(review.CounterexampleSetDigest) == "" ||
		strings.TrimSpace(review.FeedbackHistoryDigest) == "" ||
		strings.TrimSpace(review.CausalReason) == "" ||
		strings.TrimSpace(review.ReviewDigest) == "" {
		return errors.New("causal review is incomplete")
	}
	if !review.NonAuthorizing {
		return errors.New("causal review must remain non-authorizing")
	}
	switch review.Status {
	case CausalReviewReview, CausalReviewUnknown:
	default:
		return errors.New("unsupported causal review status")
	}
	copy := review
	copy.ReviewDigest = ""
	digest, err := Digest(copy)
	if err != nil {
		return err
	}
	if digest != review.ReviewDigest {
		return errors.New("causal review digest does not match its evidence")
	}
	return nil
}
