package decision

import "fmt"

type JEVExternalApplyOriginBindingLSPDiagnostic struct {
    Severity               string
    Code                   string
    Message                string
    Status                 string
    MissingStage           string
    OriginDigest           string
    CandidateDigest        string
    RequestEvidenceDigest  string
    BindingEvidenceDigest  string
    EvidenceDigest         string
    Publishable            bool
    NonExecuting           bool
    NonAuthorizing         bool
}

func (d JEVExternalApplyOriginBindingLSPDiagnostic) Validate() error {
    if d.Severity == "" || d.Code == "" || d.Message == "" || d.Status == "" {
        return fmt.Errorf("incomplete external apply origin binding LSP diagnostic")
    }
    if !d.NonExecuting {
        return fmt.Errorf("external apply origin binding LSP diagnostic must be non-executing")
    }
    if !d.NonAuthorizing {
        return fmt.Errorf("external apply origin binding LSP diagnostic must be non-authorizing")
    }
    if d.Publishable && (d.OriginDigest == "" || d.CandidateDigest == "" || d.RequestEvidenceDigest == "" || d.BindingEvidenceDigest == "" || d.EvidenceDigest == "") {
        return fmt.Errorf("publishable external apply origin binding LSP diagnostic requires evidence")
    }
    return nil
}

func ProjectJEVExternalApplyOriginBindingLSP(input JEVExternalApplyOriginBinding) JEVExternalApplyOriginBindingLSPDiagnostic {
    output := JEVExternalApplyOriginBindingLSPDiagnostic{
        Status:                input.Status,
        MissingStage:          input.MissingStage,
        OriginDigest:          input.OriginDigest,
        CandidateDigest:       input.CandidateDigest,
        RequestEvidenceDigest: input.RequestEvidenceDigest,
        BindingEvidenceDigest: input.BindingEvidenceDigest,
        EvidenceDigest:        input.EvidenceDigest,
        NonExecuting:          input.NonExecuting,
        NonAuthorizing:        input.NonAuthorizing,
    }
    if err := input.Validate(); err != nil {
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "JEV external apply origin binding is not publishable: evidence is incomplete or invalid"
        output.Publishable = false
        return output
    }
    output.Publishable = true
    switch input.Status {
    case jevExternalApplyOriginBindingReady:
        output.Severity = "info"
        output.Code = "jev.external-apply.origin-bound"
        output.Message = "External apply request is bound to a resolved origin; no automatic execution occurred"
    default:
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "JEV external apply origin binding is UNKNOWN; missing evidence must be resolved before review"
        output.Publishable = false
    }
    return output
}
