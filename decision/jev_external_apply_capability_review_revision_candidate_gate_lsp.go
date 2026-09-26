package decision

import "fmt"

type JEVExternalApplyCapabilityReviewRevisionCandidateGateLSPDiagnostic struct {
    Severity        string
    Code            string
    Message         string
    Status          string
    MissingStage    string
    Decision        string
    Direction       string
    Target          string
    CandidateSource string
    CandidateDigest string
    DirectionDigest string
    GateDigest      string
    Publishable     bool
    NonExecuting    bool
    NonAuthorizing  bool
}

func (d JEVExternalApplyCapabilityReviewRevisionCandidateGateLSPDiagnostic) Validate() error {
    if d.Severity == "" || d.Code == "" || d.Message == "" || d.Status == "" {
        return fmt.Errorf("incomplete external apply capability review revision candidate gate LSP diagnostic")
    }
    if d.Decision != jevExternalApplyCapabilityReviewRevisionCandidateReady && d.Decision != jevExternalApplyCapabilityReviewRevisionCandidateHold && d.Decision != jevExternalApplyCapabilityReviewRevisionCandidateRejected {
        return fmt.Errorf("invalid external apply capability review revision candidate gate LSP decision")
    }
    if d.Decision == jevExternalApplyCapabilityReviewRevisionCandidateReady && d.CandidateDigest == "" {
        return fmt.Errorf("ready candidate gate LSP diagnostic requires candidate digest")
    }
    if !d.NonExecuting {
        return fmt.Errorf("external apply capability review revision candidate gate LSP diagnostic must be non-executing")
    }
    if !d.NonAuthorizing {
        return fmt.Errorf("external apply capability review revision candidate gate LSP diagnostic must be non-authorizing")
    }
    if d.Publishable && (d.Target == "" || d.CandidateSource == "" || d.DirectionDigest == "" || d.GateDigest == "") {
        return fmt.Errorf("publishable external apply capability review revision candidate gate LSP diagnostic requires evidence")
    }
    return nil
}

func ProjectJEVExternalApplyCapabilityReviewRevisionCandidateGateLSP(input JEVExternalApplyCapabilityReviewRevisionCandidateGate) JEVExternalApplyCapabilityReviewRevisionCandidateGateLSPDiagnostic {
    output := JEVExternalApplyCapabilityReviewRevisionCandidateGateLSPDiagnostic{
        Status:          input.Status,
        MissingStage:    input.MissingStage,
        Decision:        input.Decision,
        Direction:       input.Direction,
        Target:          input.Target,
        CandidateSource: input.CandidateSource,
        CandidateDigest: input.CandidateDigest,
        DirectionDigest: input.DirectionDigest,
        GateDigest:      input.GateDigest,
        NonExecuting:    input.NonExecuting,
        NonAuthorizing:  input.NonAuthorizing,
    }
    if err := input.Validate(); err != nil {
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "JEV revision candidate gate is not publishable: evidence is incomplete or invalid"
        output.Publishable = false
        return output
    }
    output.Publishable = true
    switch input.Decision {
    case jevExternalApplyCapabilityReviewRevisionCandidateReady:
        output.Severity = "info"
        output.Code = "jev.external-apply.candidate-ready"
        output.Message = "Revision candidate is ready for external review; no automatic application occurred"
    case jevExternalApplyCapabilityReviewRevisionCandidateHold:
        output.Severity = "info"
        output.Code = "jev.external-apply.candidate-hold"
        output.Message = "Revision candidate generation is on hold pending further evidence"
    case jevExternalApplyCapabilityReviewRevisionCandidateRejected:
        output.Severity = "warning"
        output.Code = "jev.external-apply.candidate-rejected"
        output.Message = "Revision candidate direction is rejected; no automatic application occurred"
    default:
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "JEV revision candidate gate is UNKNOWN; evidence must be resolved before review"
        output.Publishable = false
    }
    return output
}
