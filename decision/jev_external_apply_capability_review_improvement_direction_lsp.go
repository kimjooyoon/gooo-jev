package decision

import "fmt"

type JEVExternalApplyCapabilityReviewImprovementDirectionLSPDiagnostic struct {
    Severity        string
    Code            string
    Message         string
    Status          string
    MissingStage    string
    Direction       string
    Target          string
    CandidateSource string
    FeedbackDigest  string
    DirectionDigest string
    Publishable     bool
    NonExecuting    bool
    NonAuthorizing  bool
}

func (d JEVExternalApplyCapabilityReviewImprovementDirectionLSPDiagnostic) Validate() error {
    if d.Severity == "" || d.Code == "" || d.Message == "" || d.Status == "" {
        return fmt.Errorf("incomplete external apply capability review improvement direction LSP diagnostic")
    }
    if d.Direction != jevExternalApplyCapabilityReviewImprovementDirectionGenerate && d.Direction != jevExternalApplyCapabilityReviewImprovementDirectionHold && d.Direction != jevExternalApplyCapabilityReviewImprovementDirectionReject {
        return fmt.Errorf("invalid external apply capability review improvement direction LSP direction")
    }
    if !d.NonExecuting {
        return fmt.Errorf("external apply capability review improvement direction LSP diagnostic must be non-executing")
    }
    if !d.NonAuthorizing {
        return fmt.Errorf("external apply capability review improvement direction LSP diagnostic must be non-authorizing")
    }
    if d.Publishable && (d.Target == "" || d.CandidateSource == "" || d.FeedbackDigest == "" || d.DirectionDigest == "") {
        return fmt.Errorf("publishable external apply capability review improvement direction LSP diagnostic requires evidence")
    }
    return nil
}

func ProjectJEVExternalApplyCapabilityReviewImprovementDirectionLSP(input JEVExternalApplyCapabilityReviewImprovementDirection) JEVExternalApplyCapabilityReviewImprovementDirectionLSPDiagnostic {
    output := JEVExternalApplyCapabilityReviewImprovementDirectionLSPDiagnostic{
        Status:          input.Status,
        MissingStage:    input.MissingStage,
        Direction:       input.Direction,
        Target:          input.Target,
        CandidateSource: input.CandidateSource,
        FeedbackDigest:  input.FeedbackDigest,
        DirectionDigest: input.DirectionDigest,
        NonExecuting:    input.NonExecuting,
        NonAuthorizing:  input.NonAuthorizing,
    }
    if err := input.Validate(); err != nil {
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "JEV improvement direction is not publishable: evidence is incomplete or invalid"
        output.Publishable = false
        return output
    }
    output.Publishable = true
    switch input.Direction {
    case jevExternalApplyCapabilityReviewImprovementDirectionGenerate:
        output.Severity = "info"
        output.Code = "jev.external-apply.candidate-generation"
        output.Message = "Evidence supports candidate generation for external review; no automatic change occurred"
    case jevExternalApplyCapabilityReviewImprovementDirectionHold:
        output.Severity = "info"
        output.Code = "jev.external-apply.candidate-hold"
        output.Message = "Evidence supports holding candidate generation until further review"
    case jevExternalApplyCapabilityReviewImprovementDirectionReject:
        output.Severity = "warning"
        output.Code = "jev.external-apply.candidate-rejection"
        output.Message = "Evidence supports rejecting the candidate direction; no automatic change occurred"
    default:
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "JEV improvement direction is UNKNOWN; evidence must be resolved before review"
        output.Publishable = false
    }
    return output
}
