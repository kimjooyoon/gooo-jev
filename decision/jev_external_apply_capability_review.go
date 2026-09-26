package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "strings"
)

const jevExternalApplyCapabilityReviewConfirmed = "capability-review-confirmed"
const jevExternalApplyCapabilityReviewRefuted = "capability-review-refuted"
const jevExternalApplyCapabilityReviewUnknown = "UNKNOWN"

type JEVExternalApplyCapabilityReviewInput struct {
    Capability          JEVExternalApplyCapabilityBoundary
    Outcome             string
    ReviewerPrincipal   string
    ReviewEvidenceDigest string
    NonAuthorizing      bool
}

type JEVExternalApplyCapabilityReview struct {
    Status               string
    MissingStage         string
    CapabilityDigest     string
    ReviewerPrincipal    string
    ReviewEvidenceDigest string
    ReviewDigest         string
    NonExecuting         bool
    NonAuthorizing       bool
}

func (r JEVExternalApplyCapabilityReview) Validate() error {
    if r.Status == "" || r.CapabilityDigest == "" || r.ReviewerPrincipal == "" || r.ReviewEvidenceDigest == "" || r.ReviewDigest == "" {
        return fmt.Errorf("incomplete JEV external apply capability review")
    }
    if r.Status != jevExternalApplyCapabilityReviewConfirmed && r.Status != jevExternalApplyCapabilityReviewRefuted {
        return fmt.Errorf("invalid JEV external apply capability review status")
    }
    if !strings.HasPrefix(r.ReviewerPrincipal, "spiffe://") {
        return fmt.Errorf("JEV external apply capability reviewer must use SPIFFE URI form")
    }
    if !r.NonExecuting {
        return fmt.Errorf("JEV external apply capability review must be non-executing")
    }
    if !r.NonAuthorizing {
        return fmt.Errorf("JEV external apply capability review must be non-authorizing")
    }
    expected := digestJEVExternalApplyCapabilityReview(r.Status, r.CapabilityDigest, r.ReviewerPrincipal, r.ReviewEvidenceDigest)
    if r.ReviewDigest != expected {
        return fmt.Errorf("JEV external apply capability review digest mismatch")
    }
    return nil
}

func ObserveJEVExternalApplyCapabilityReview(input JEVExternalApplyCapabilityReviewInput) JEVExternalApplyCapabilityReview {
    output := JEVExternalApplyCapabilityReview{
        Status:         jevExternalApplyCapabilityReviewUnknown,
        NonExecuting:   true,
        NonAuthorizing: true,
    }
    if !input.NonAuthorizing || !input.Capability.NonAuthorizing {
        output.MissingStage = "authorization-boundary"
        return output
    }
    if !input.Capability.NonExecuting {
        output.MissingStage = "execution-boundary"
        return output
    }
    if err := input.Capability.Validate(); err != nil {
        output.MissingStage = "capability-boundary"
        return output
    }
    if input.Outcome != jevExternalApplyCapabilityReviewConfirmed && input.Outcome != jevExternalApplyCapabilityReviewRefuted {
        output.MissingStage = "review-outcome"
        return output
    }
    if !strings.HasPrefix(input.ReviewerPrincipal, "spiffe://") {
        output.MissingStage = "reviewer-principal"
        return output
    }
    if input.ReviewEvidenceDigest == "" {
        output.MissingStage = "review-evidence"
        return output
    }
    output.Status = input.Outcome
    output.CapabilityDigest = input.Capability.CapabilityDigest
    output.ReviewerPrincipal = input.ReviewerPrincipal
    output.ReviewEvidenceDigest = input.ReviewEvidenceDigest
    output.ReviewDigest = digestJEVExternalApplyCapabilityReview(output.Status, output.CapabilityDigest, output.ReviewerPrincipal, output.ReviewEvidenceDigest)
    if err := output.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewUnknown
        output.MissingStage = "capability-review-evidence"
        output.CapabilityDigest = ""
        output.ReviewDigest = ""
    }
    return output
}

func digestJEVExternalApplyCapabilityReview(status, capabilityDigest, reviewerPrincipal, reviewEvidenceDigest string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s", status, capabilityDigest, reviewerPrincipal, reviewEvidenceDigest)))
    return hex.EncodeToString(sum[:])
}
