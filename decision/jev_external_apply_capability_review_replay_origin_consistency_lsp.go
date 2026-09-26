package decision

import "fmt"

type JEVExternalApplyCapabilityReviewReplayOriginConsistencyLSPDiagnostic struct {
    Severity             string
    Code                 string
    Message              string
    Status               string
    MissingStage         string
    ObservationStatus    string
    DeclaredOriginDigest string
    ReverseOriginDigest  string
    ObservationDigest    string
    ConsistencyDigest    string
    Publishable          bool
    NonExecuting         bool
    NonAuthorizing       bool
}

func (d JEVExternalApplyCapabilityReviewReplayOriginConsistencyLSPDiagnostic) Validate() error {
    if d.Severity == "" || d.Code == "" || d.Message == "" || d.Status == "" {
        return fmt.Errorf("incomplete external apply capability review replay origin consistency LSP diagnostic")
    }
    if d.Status != jevExternalApplyCapabilityReviewReplayOriginConsistent &&
        d.Status != jevExternalApplyCapabilityReviewReplayOriginMismatch {
        return fmt.Errorf("invalid external apply capability review replay origin consistency LSP status")
    }
    if d.ObservationStatus != jevExternalApplyCapabilityReviewReplayObserved {
        return fmt.Errorf("origin consistency LSP diagnostic requires replay-observed status")
    }
    if d.DeclaredOriginDigest == "" || d.ReverseOriginDigest == "" || d.ObservationDigest == "" || d.ConsistencyDigest == "" {
        return fmt.Errorf("origin consistency LSP diagnostic requires complete evidence")
    }
    if d.Status == jevExternalApplyCapabilityReviewReplayOriginConsistent &&
        d.DeclaredOriginDigest != d.ReverseOriginDigest {
        return fmt.Errorf("origin-consistent LSP diagnostic requires equal origin digests")
    }
    if d.Status == jevExternalApplyCapabilityReviewReplayOriginMismatch &&
        d.DeclaredOriginDigest == d.ReverseOriginDigest {
        return fmt.Errorf("origin-mismatch LSP diagnostic requires different origin digests")
    }
    if !d.NonExecuting {
        return fmt.Errorf("external apply capability review replay origin consistency LSP diagnostic must be non-executing")
    }
    if !d.NonAuthorizing {
        return fmt.Errorf("external apply capability review replay origin consistency LSP diagnostic must be non-authorizing")
    }
    return nil
}

func ProjectJEVExternalApplyCapabilityReviewReplayOriginConsistencyLSP(input JEVExternalApplyCapabilityReviewReplayOriginConsistency) JEVExternalApplyCapabilityReviewReplayOriginConsistencyLSPDiagnostic {
    output := JEVExternalApplyCapabilityReviewReplayOriginConsistencyLSPDiagnostic{
        Status:               input.Status,
        MissingStage:         input.MissingStage,
        ObservationStatus:    input.ObservationStatus,
        DeclaredOriginDigest: input.DeclaredOriginDigest,
        ReverseOriginDigest:  input.ReverseOriginDigest,
        ObservationDigest:    input.ObservationDigest,
        ConsistencyDigest:    input.ConsistencyDigest,
        NonExecuting:         input.NonExecuting,
        NonAuthorizing:       input.NonAuthorizing,
    }
    if err := input.Validate(); err != nil {
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "JEV replay origin consistency is not publishable: evidence is incomplete or invalid"
        output.Publishable = false
        return output
    }
    output.Publishable = true
    switch input.Status {
    case jevExternalApplyCapabilityReviewReplayOriginConsistent:
        output.Severity = "info"
        output.Code = "jev.external-apply.origin-consistent"
        output.Message = "Declared and reverse-observed origins are consistent; no execution or authorization occurred"
    case jevExternalApplyCapabilityReviewReplayOriginMismatch:
        output.Severity = "warning"
        output.Code = "jev.external-apply.origin-mismatch"
        output.Message = "Declared and reverse-observed origins differ; review is required and no execution or authorization occurred"
    default:
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "JEV replay origin consistency is UNKNOWN; evidence must be resolved before publication"
        output.Publishable = false
    }
    return output
}
