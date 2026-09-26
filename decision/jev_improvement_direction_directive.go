package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

const (
    jevImprovementDirectiveExternalReview = "request-external-apply-review"
    jevImprovementDirectiveRevision = "generate-revision-candidate"
    jevImprovementDirectiveEvidence = "await-missing-evidence"
    jevImprovementDirectiveUnknown = "UNKNOWN"
)

type JEVImprovementDirectionDirectiveInput struct {
    Aggregation    JEVImprovementFeedbackAggregation
    CandidateDigest string
    CandidateSource string
    NonAuthorizing bool
}

type JEVImprovementDirectionDirective struct {
    Status              string
    Directive           string
    MissingStage        string
    CandidateDigest     string
    CandidateSource     string
    InputEvidenceDigest string
    EvidenceDigest      string
    NonExecuting        bool
    NonAuthorizing      bool
}

func (d JEVImprovementDirectionDirective) Validate() error {
    if d.Status == "" || d.Directive == "" || d.CandidateDigest == "" || d.CandidateSource == "" || d.InputEvidenceDigest == "" || d.EvidenceDigest == "" {
        return fmt.Errorf("incomplete JEV improvement direction directive")
    }
    switch d.Directive {
    case jevImprovementDirectiveExternalReview, jevImprovementDirectiveRevision, jevImprovementDirectiveEvidence:
    default:
        return fmt.Errorf("invalid JEV improvement direction directive")
    }
    if !d.NonExecuting {
        return fmt.Errorf("JEV improvement direction directive must be non-executing")
    }
    if !d.NonAuthorizing {
        return fmt.Errorf("JEV improvement direction directive must be non-authorizing")
    }
    expected := digestJEVImprovementDirectionDirective(d.Status, d.Directive, d.CandidateDigest, d.CandidateSource, d.InputEvidenceDigest)
    if d.EvidenceDigest != expected {
        return fmt.Errorf("JEV improvement direction directive evidence digest mismatch")
    }
    return nil
}

func DeriveJEVImprovementDirectionDirective(input JEVImprovementDirectionDirectiveInput) JEVImprovementDirectionDirective {
    output := JEVImprovementDirectionDirective{
        Status: jevImprovementDirectiveUnknown,
        NonExecuting: true,
        NonAuthorizing: true,
    }
    if !input.NonAuthorizing || !input.Aggregation.NonAuthorizing {
        output.MissingStage = "authorization-boundary"
        return output
    }
    if !input.Aggregation.NonExecuting {
        output.MissingStage = "execution-boundary"
        return output
    }
    if err := input.Aggregation.Validate(); err != nil {
        output.MissingStage = "feedback-aggregation"
        return output
    }
    if input.CandidateDigest == "" {
        output.MissingStage = "candidate-digest"
        return output
    }
    if input.CandidateSource == "" {
        output.MissingStage = "candidate-source"
        return output
    }
    output.CandidateDigest = input.CandidateDigest
    output.CandidateSource = input.CandidateSource
    output.InputEvidenceDigest = input.Aggregation.EvidenceDigest
    switch input.Aggregation.Status {
    case jevImprovementFeedbackStableForReview:
        output.Status = jevImprovementDirectiveExternalReview
        output.Directive = jevImprovementDirectiveExternalReview
    case jevImprovementFeedbackNeedsRevision:
        output.Status = jevImprovementDirectiveRevision
        output.Directive = jevImprovementDirectiveRevision
    case jevImprovementFeedbackHold:
        output.Status = jevImprovementDirectiveEvidence
        output.Directive = jevImprovementDirectiveEvidence
    default:
        output.MissingStage = "feedback-direction"
        return output
    }
    output.EvidenceDigest = digestJEVImprovementDirectionDirective(output.Status, output.Directive, output.CandidateDigest, output.CandidateSource, output.InputEvidenceDigest)
    if err := output.Validate(); err != nil {
        output.Status = jevImprovementDirectiveUnknown
        output.MissingStage = "direction-directive-evidence"
        output.EvidenceDigest = ""
    }
    return output
}

func digestJEVImprovementDirectionDirective(status, directive, candidateDigest, candidateSource, inputEvidenceDigest string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s|%s", status, directive, candidateDigest, candidateSource, inputEvidenceDigest)))
    return hex.EncodeToString(sum[:])
}
