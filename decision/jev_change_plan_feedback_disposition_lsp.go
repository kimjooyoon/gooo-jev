package decision

import "fmt"

type DecisionConfidenceChangePlanFeedbackDispositionLSPDiagnostic struct {
    Severity        string
    Code            string
    Message         string
    Status          string
    MissingStage    string
    EvidenceDigest  string
    Publishable     bool
    NonAuthorizing  bool
}

func (d DecisionConfidenceChangePlanFeedbackDispositionLSPDiagnostic) Validate() error {
    if d.Severity == "" || d.Code == "" || d.Message == "" || d.Status == "" {
        return fmt.Errorf("incomplete change plan feedback disposition LSP diagnostic")
    }
    if !d.NonAuthorizing {
        return fmt.Errorf("change plan feedback disposition LSP diagnostic must be non-authorizing")
    }
    if d.Publishable && d.EvidenceDigest == "" {
        return fmt.Errorf("publishable change plan feedback disposition LSP diagnostic requires evidence")
    }
    return nil
}

func ProjectDecisionConfidenceChangePlanFeedbackDispositionLSP(input DecisionConfidenceChangePlanFeedbackDisposition) DecisionConfidenceChangePlanFeedbackDispositionLSPDiagnostic {
    output := DecisionConfidenceChangePlanFeedbackDispositionLSPDiagnostic{
        Status: input.Status,
        MissingStage: input.MissingStage,
        EvidenceDigest: input.EvidenceDigest,
        NonAuthorizing: input.NonAuthorizing,
    }
    if err := input.Validate(); err != nil {
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "JEV disposition is not publishable: evidence is incomplete or invalid"
        output.Publishable = false
        return output
    }
    output.Publishable = true
    switch input.Status {
    case changePlanFeedbackDispositionReady:
        output.Severity = "info"
        output.Code = "jev.change-plan.ready-for-external-apply-review"
        output.Message = "Change plan is ready for external apply review; no automatic apply occurred"
    case changePlanFeedbackDispositionAbort:
        output.Severity = "warning"
        output.Code = "jev.change-plan.abort"
        output.Message = "Feedback refuted the change plan; application is blocked"
    case changePlanFeedbackDispositionHold:
        output.Severity = "warning"
        output.Code = "jev.change-plan.hold"
        output.Message = "Change plan remains on hold because feedback is inconclusive"
    default:
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "JEV disposition is UNKNOWN; missing evidence must be resolved before review"
        output.Publishable = false
    }
    return output
}
