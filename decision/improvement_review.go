package decision

import (
    "errors"
    "strings"
)

type ImprovementReviewDecision string

const (
    ReviewRequired ImprovementReviewDecision = "required"
    ReviewPassed   ImprovementReviewDecision = "passed"
    ReviewFailed   ImprovementReviewDecision = "failed"
    ReviewUnknown  ImprovementReviewDecision = "unknown"
)

type ImprovementReviewReceipt struct {
    CandidateDigest      string
    ReviewerReference    string
    ReviewEvidenceDigest string
    Decision             ImprovementReviewDecision
    ExecutionGranted     bool
    ReviewDigest         string
}

func ReviewImprovementCandidate(
    candidate ImprovementCandidate,
    reviewerReference,
    reviewEvidenceDigest string,
    decision ImprovementReviewDecision,
) (ImprovementReviewReceipt, error) {
    if err := candidate.Validate(); err != nil {
        return ImprovementReviewReceipt{}, err
    }
    review := ImprovementReviewReceipt{
        CandidateDigest:      candidate.CandidateDigest,
        ReviewerReference:    reviewerReference,
        ReviewEvidenceDigest: reviewEvidenceDigest,
        Decision:             decision,
        ExecutionGranted:     false,
    }
    if err := review.validateShape(); err != nil {
        return ImprovementReviewReceipt{}, err
    }
    digest, err := Digest(review)
    if err != nil {
        return ImprovementReviewReceipt{}, err
    }
    review.ReviewDigest = digest
    return review, nil
}

func (review ImprovementReviewReceipt) Validate() error {
    if err := review.validateShape(); err != nil {
        return err
    }
    if strings.TrimSpace(review.ReviewDigest) == "" {
        return errors.New("improvement review digest is missing")
    }
    copy := review
    copy.ReviewDigest = ""
    digest, err := Digest(copy)
    if err != nil {
        return err
    }
    if digest != review.ReviewDigest {
        return errors.New("improvement review digest does not match its evidence")
    }
    return nil
}

func (review ImprovementReviewReceipt) validateShape() error {
    if strings.TrimSpace(review.CandidateDigest) == "" ||
        strings.TrimSpace(review.ReviewerReference) == "" ||
        strings.TrimSpace(review.ReviewEvidenceDigest) == "" {
        return errors.New("improvement review is incomplete")
    }
    if review.ExecutionGranted {
        return errors.New("improvement review must not grant execution")
    }
    switch review.Decision {
    case ReviewRequired, ReviewPassed, ReviewFailed, ReviewUnknown:
        return nil
    default:
        return errors.New("unsupported improvement review decision")
    }
}
