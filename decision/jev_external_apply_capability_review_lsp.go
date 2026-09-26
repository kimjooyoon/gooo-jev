package decision

import "fmt"

type JEVExternalApplyCapabilityReviewLSPDiagnostic struct {
    Severity               string
    Code                   string
    Message               string
    Status                 string
    MissingStage           string
    CapabilityDigest       string
    ReviewerPrincipal      string
    ReviewEvidenceDigest   string
    ReviewDigest           string
    Publishable            bool
    NonExecuting           bool
    NonAuthorizing         bool
}

func (d JEVExternalApplyCapabilityReviewLSPDiagnostic) Validate() error {
    if d.Severity == "" || d.Code == "" || d.Message == "" || d.Status == "" {
        return fmt.Errorf("incomplete external apply capability review LSP diagnostic")
    }
    if !d.NonExecuting {
        return fmt.Errorf("external apply capability review LSP diagnostic must be non-executing")
    }
    if !d.NonAuthorizing {
        return fmt.Errorf("external apply capability review LSP diagnostic must be non-authorizing")
    }
    if d.Publishable && (d.CapabilityDigest == "" || d.ReviewerPrincipal == "" || d.ReviewEvidenceDigest == "" || d.ReviewDigest == "") {
        return fmt.Errorf("publishable external apply capability review LSP diagnostic requires evidence")
    }
    return nil
}

func ProjectJEVExternalApplyCapabilityReviewLSP(input JEVExternalApplyCapabilityReview) JEVExternalApplyCapabilityReviewLSPDiagnostic {
    output := JEVExternalApplyCapabilityReviewLSPDiagnostic{
        Status:               input.Status,
        MissingStage:         input.MissingStage,
        CapabilityDigest:     input.CapabilityDigest,
        ReviewerPrincipal:    input.ReviewerPrincipal,
        ReviewEvidenceDigest: input.ReviewEvidenceDigest,
        ReviewDigest:         input.ReviewDigest,
        NonExecuting:         input.NonExecuting,
        NonAuthorizing:       input.NonAuthorizing,
    }
    if err := input.Validate(); err != nil {
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "JEV capability review is not publishable: evidence is incomplete or invalid"
        output.Publishable = false
        return output
    }
    output.Publishable = true
    switch input.Status {
    case jevExternalApplyCapabilityReviewConfirmed:
        output.Severity = "info"
        output.Code = "jev.external-apply.review-confirmed"
        output.Message = "External capability review is confirmed; no authorization or execution occurred"
    case jevExternalApplyCapabilityReviewRefuted:
        output.Severity = "warning"
        output.Code = "jev.external-apply.review-refuted"
        output.Message = "External capability review is refuted; external application remains blocked"
    default:
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "JEV capability review is UNKNOWN; evidence must be resolved before review"
        output.Publishable = false
    }
    return output
}
