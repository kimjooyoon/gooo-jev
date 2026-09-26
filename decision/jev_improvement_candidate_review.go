package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

const (
    jevImprovementCandidateReviewAccepted = "accepted"
    jevImprovementCandidateReviewRejected = "rejected"
    jevImprovementCandidateReviewUnknown = "unknown"
    jevImprovementCandidateReviewConfirmed = "review-confirmed-for-external-apply"
    jevImprovementCandidateReviewAbort = "abort"
    jevImprovementCandidateReviewHold = "hold"
    jevImprovementCandidateReviewStatusUnknown = "UNKNOWN"
)

type JEVImprovementCandidateReviewInput struct {
    Selection       JEVImprovementCandidateSelection
    ReviewerReference string
    Decision        string
    ReviewEvidenceDigest string
    NonAuthorizing  bool
}

type JEVImprovementCandidateReview struct {
    Status                   string
    MissingStage             string
    CandidateDigest          string
    ReviewerReference        string
    Decision                 string
    SelectionEvidenceDigest  string
    ReviewEvidenceDigest     string
    EvidenceDigest           string
    NonExecuting             bool
    NonAuthorizing           bool
}

func (r JEVImprovementCandidateReview) Validate() error {
    if r.Status == "" || r.CandidateDigest == "" || r.ReviewerReference == "" || r.Decision == "" || r.SelectionEvidenceDigest == "" || r.ReviewEvidenceDigest == "" || r.EvidenceDigest == "" {
        return fmt.Errorf("incomplete JEV improvement candidate review")
    }
    if !r.NonExecuting {
        return fmt.Errorf("JEV improvement candidate review must be non-executing")
    }
    if !r.NonAuthorizing {
        return fmt.Errorf("JEV improvement candidate review must be non-authorizing")
    }
    expected := digestJEVImprovementCandidateReview(r.Status, r.CandidateDigest, r.ReviewerReference, r.Decision, r.SelectionEvidenceDigest, r.ReviewEvidenceDigest)
    if r.EvidenceDigest != expected {
        return fmt.Errorf("JEV improvement candidate review evidence digest mismatch")
    }
    return nil
}

func ReviewJEVImprovementCandidate(input JEVImprovementCandidateReviewInput) JEVImprovementCandidateReview {
    output := JEVImprovementCandidateReview{
        Status: jevImprovementCandidateReviewStatusUnknown,
        NonExecuting: true,
        NonAuthorizing: true,
    }
    if !input.NonAuthorizing || !input.Selection.NonAuthorizing {
        output.MissingStage = "authorization-boundary"
        return output
    }
    if !input.Selection.NonExecuting {
        output.MissingStage = "execution-boundary"
        return output
    }
    if err := input.Selection.Validate(); err != nil {
        output.MissingStage = "candidate-selection"
        return output
    }
    if input.Selection.Status != jevImprovementCandidateReady {
        output.MissingStage = "external-review-candidate"
        return output
    }
    if input.ReviewerReference == "" {
        output.MissingStage = "reviewer-reference"
        return output
    }
    if input.ReviewEvidenceDigest == "" {
        output.MissingStage = "review-evidence"
        return output
    }
    output.CandidateDigest = input.Selection.CandidateDigest
    output.ReviewerReference = input.ReviewerReference
    output.Decision = input.Decision
    output.SelectionEvidenceDigest = input.Selection.EvidenceDigest
    output.ReviewEvidenceDigest = input.ReviewEvidenceDigest
    switch input.Decision {
    case jevImprovementCandidateReviewAccepted:
        output.Status = jevImprovementCandidateReviewConfirmed
    case jevImprovementCandidateReviewRejected:
        output.Status = jevImprovementCandidateReviewAbort
    case jevImprovementCandidateReviewUnknown:
        output.Status = jevImprovementCandidateReviewHold
    default:
        output.MissingStage = "review-decision"
        output.CandidateDigest = ""
        output.SelectionEvidenceDigest = ""
        output.ReviewEvidenceDigest = ""
        return output
    }
    output.EvidenceDigest = digestJEVImprovementCandidateReview(output.Status, output.CandidateDigest, output.ReviewerReference, output.Decision, output.SelectionEvidenceDigest, output.ReviewEvidenceDigest)
    if err := output.Validate(); err != nil {
        output.Status = jevImprovementCandidateReviewStatusUnknown
        output.MissingStage = "review-evidence-binding"
        output.EvidenceDigest = ""
    }
    return output
}

func digestJEVImprovementCandidateReview(status, candidateDigest, reviewerReference, decision, selectionEvidenceDigest, reviewEvidenceDigest string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s|%s|%s", status, candidateDigest, reviewerReference, decision, selectionEvidenceDigest, reviewEvidenceDigest)))
    return hex.EncodeToString(sum[:])
}
