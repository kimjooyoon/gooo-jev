package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

const jevExternalApplyRequestReady = "ready-for-external-apply"
const jevExternalApplyRequestUnknown = "UNKNOWN"

type JEVExternalApplyRequestInput struct {
    Review              JEVImprovementCandidateReview
    Aggregation         JEVImprovementFeedbackAggregation
    TargetReference     string
    NonAuthorizing      bool
}

type JEVExternalApplyRequest struct {
    Status                  string
    MissingStage            string
    CandidateDigest         string
    ReviewEvidenceDigest    string
    AggregationEvidenceDigest string
    TargetReference         string
    RequestEvidenceDigest   string
    EvidenceDigest          string
    NonExecuting            bool
    NonAuthorizing          bool
}

func (r JEVExternalApplyRequest) Validate() error {
    if r.Status == "" || r.CandidateDigest == "" || r.ReviewEvidenceDigest == "" || r.AggregationEvidenceDigest == "" || r.TargetReference == "" || r.RequestEvidenceDigest == "" || r.EvidenceDigest == "" {
        return fmt.Errorf("incomplete JEV external apply request")
    }
    if r.Status != jevExternalApplyRequestReady {
        return fmt.Errorf("invalid JEV external apply request status")
    }
    if !r.NonExecuting {
        return fmt.Errorf("JEV external apply request must be non-executing")
    }
    if !r.NonAuthorizing {
        return fmt.Errorf("JEV external apply request must be non-authorizing")
    }
    expectedRequest := digestJEVExternalApplyRequestEvidence(r.CandidateDigest, r.ReviewEvidenceDigest, r.AggregationEvidenceDigest, r.TargetReference)
    if r.RequestEvidenceDigest != expectedRequest {
        return fmt.Errorf("JEV external apply request evidence mismatch")
    }
    expected := digestJEVExternalApplyRequest(r.Status, r.CandidateDigest, r.RequestEvidenceDigest)
    if r.EvidenceDigest != expected {
        return fmt.Errorf("JEV external apply request digest mismatch")
    }
    return nil
}

func CreateJEVExternalApplyRequest(input JEVExternalApplyRequestInput) JEVExternalApplyRequest {
    output := JEVExternalApplyRequest{
        Status: jevExternalApplyRequestUnknown,
        NonExecuting: true,
        NonAuthorizing: true,
    }
    if !input.NonAuthorizing || !input.Review.NonAuthorizing || !input.Aggregation.NonAuthorizing {
        output.MissingStage = "authorization-boundary"
        return output
    }
    if !input.Review.NonExecuting || !input.Aggregation.NonExecuting {
        output.MissingStage = "execution-boundary"
        return output
    }
    if err := input.Review.Validate(); err != nil {
        output.MissingStage = "candidate-review"
        return output
    }
    if input.Review.Status != jevImprovementCandidateReviewConfirmed {
        output.MissingStage = "review-confirmation"
        return output
    }
    if err := input.Aggregation.Validate(); err != nil {
        output.MissingStage = "feedback-aggregation"
        return output
    }
    if input.Aggregation.Status != jevImprovementFeedbackStableForReview {
        output.MissingStage = "stable-feedback-aggregation"
        return output
    }
    if input.TargetReference == "" {
        output.MissingStage = "target-reference"
        return output
    }
    output.Status = jevExternalApplyRequestReady
    output.CandidateDigest = input.Review.CandidateDigest
    output.ReviewEvidenceDigest = input.Review.EvidenceDigest
    output.AggregationEvidenceDigest = input.Aggregation.EvidenceDigest
    output.TargetReference = input.TargetReference
    output.RequestEvidenceDigest = digestJEVExternalApplyRequestEvidence(output.CandidateDigest, output.ReviewEvidenceDigest, output.AggregationEvidenceDigest, output.TargetReference)
    output.EvidenceDigest = digestJEVExternalApplyRequest(output.Status, output.CandidateDigest, output.RequestEvidenceDigest)
    if err := output.Validate(); err != nil {
        output.Status = jevExternalApplyRequestUnknown
        output.MissingStage = "external-apply-request-evidence"
        output.EvidenceDigest = ""
    }
    return output
}

func digestJEVExternalApplyRequestEvidence(candidateDigest, reviewEvidenceDigest, aggregationEvidenceDigest, targetReference string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s", candidateDigest, reviewEvidenceDigest, aggregationEvidenceDigest, targetReference)))
    return hex.EncodeToString(sum[:])
}

func digestJEVExternalApplyRequest(status, candidateDigest, requestEvidenceDigest string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s", status, candidateDigest, requestEvidenceDigest)))
    return hex.EncodeToString(sum[:])
}
