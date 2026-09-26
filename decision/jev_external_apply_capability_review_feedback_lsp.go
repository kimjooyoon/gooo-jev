package decision

import "fmt"

type JEVExternalApplyCapabilityReviewFeedbackLSPDiagnostic struct {
    Severity                string
    Code                    string
    Message                 string
    Status                  string
    MissingStage            string
    Signal                  string
    SignalDigest            string
    TrendDigest             string
    FeedbackEvidenceDigest  string
    FeedbackDigest          string
    Publishable             bool
    NonExecuting            bool
    NonAuthorizing          bool
}

func (d JEVExternalApplyCapabilityReviewFeedbackLSPDiagnostic) Validate() error {
    if d.Severity == "" || d.Code == "" || d.Message == "" || d.Status == "" {
        return fmt.Errorf("incomplete external apply capability review feedback LSP diagnostic")
    }
    if !d.NonExecuting {
        return fmt.Errorf("external apply capability review feedback LSP diagnostic must be non-executing")
    }
    if !d.NonAuthorizing {
        return fmt.Errorf("external apply capability review feedback LSP diagnostic must be non-authorizing")
    }
    if d.Publishable && (d.Signal == "" || d.SignalDigest == "" || d.TrendDigest == "" || d.FeedbackEvidenceDigest == "" || d.FeedbackDigest == "") {
        return fmt.Errorf("publishable external apply capability review feedback LSP diagnostic requires evidence")
    }
    return nil
}

func ProjectJEVExternalApplyCapabilityReviewFeedbackLSP(input JEVExternalApplyCapabilityReviewFeedbackBridge) JEVExternalApplyCapabilityReviewFeedbackLSPDiagnostic {
    output := JEVExternalApplyCapabilityReviewFeedbackLSPDiagnostic{
        Status:                input.Status,
        MissingStage:          input.MissingStage,
        Signal:                input.Signal,
        SignalDigest:          input.SignalDigest,
        TrendDigest:           input.TrendDigest,
        FeedbackEvidenceDigest: input.FeedbackEvidenceDigest,
        FeedbackDigest:        input.FeedbackDigest,
        NonExecuting:          input.NonExecuting,
        NonAuthorizing:        input.NonAuthorizing,
    }
    if err := input.Validate(); err != nil {
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "JEV capability review feedback is not publishable: evidence is incomplete or invalid"
        output.Publishable = false
        return output
    }
    output.Publishable = true
    switch input.Status {
    case jevExternalApplyCapabilityReviewFeedbackBound:
        output.Severity = "info"
        output.Code = "jev.external-apply.feedback-bound"
        output.Message = "Capability review signal is bound to feedback evidence; no automatic change occurred"
    default:
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "JEV capability review feedback is UNKNOWN; evidence must be resolved before use"
        output.Publishable = false
    }
    return output
}
