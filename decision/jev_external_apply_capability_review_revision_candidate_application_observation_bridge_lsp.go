package decision

import "fmt"

type JEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationBridgeLSPDiagnostic struct {
    Severity                    string
    Code                        string
    Message                     string
    Status                     string
    MissingStage                string
    CandidateApplicationStatus  string
    CandidateApplicationDigest  string
    ObservationStatus           string
    ObservationSource           string
    ObservationEvidenceDigest   string
    ReverseObservationDigest    string
    ObservationDigest           string
    Publishable                 bool
    NonExecuting                bool
    NonAuthorizing              bool
}

func (d JEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationBridgeLSPDiagnostic) Validate() error {
    if d.Severity == "" || d.Code == "" || d.Message == "" || d.Status == "" {
        return fmt.Errorf("incomplete JEV candidate application observation LSP diagnostic")
    }
    if !d.NonExecuting {
        return fmt.Errorf("JEV candidate application observation LSP diagnostic must be non-executing")
    }
    if !d.NonAuthorizing {
        return fmt.Errorf("JEV candidate application observation LSP diagnostic must be non-authorizing")
    }
    if d.Status == jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationBridgeUnknown {
        if d.Publishable || d.MissingStage == "" {
            return fmt.Errorf("UNKNOWN JEV candidate application observation diagnostic must remain non-publishable")
        }
        return nil
    }
    if d.Status != jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationBridgeBound {
        return fmt.Errorf("invalid JEV candidate application observation LSP status")
    }
    if !d.Publishable || d.CandidateApplicationDigest == "" || d.ObservationStatus == "" ||
        d.ObservationDigest == "" {
        return fmt.Errorf("bound JEV candidate application observation diagnostic is incomplete")
    }
    switch d.CandidateApplicationStatus {
    case jevExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateReady:
        if d.ObservationStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationObserved &&
            d.ObservationStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMismatch {
            return fmt.Errorf("ready candidate application observation diagnostic has invalid status")
        }
        if d.ObservationSource == "" || d.ObservationEvidenceDigest == "" ||
            d.ReverseObservationDigest == "" {
            return fmt.Errorf("candidate application observation diagnostic evidence is incomplete")
        }
    case jevExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateHold:
        if d.ObservationStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationHold ||
            d.ObservationSource != "" || d.ObservationEvidenceDigest != "" ||
            d.ReverseObservationDigest != "" {
            return fmt.Errorf("held candidate application observation diagnostic has inconsistent evidence")
        }
    case jevExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateRejected:
        if d.ObservationStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationRejected ||
            d.ObservationSource != "" || d.ObservationEvidenceDigest != "" ||
            d.ReverseObservationDigest != "" {
            return fmt.Errorf("rejected candidate application observation diagnostic has inconsistent evidence")
        }
    default:
        return fmt.Errorf("invalid candidate application status")
    }
    return nil
}

func ProjectJEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationBridgeLSP(
    input JEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationBridge,
) JEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationBridgeLSPDiagnostic {
    output := JEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationBridgeLSPDiagnostic{
        Status:                   input.Status,
        MissingStage:             input.MissingStage,
        CandidateApplicationStatus: input.CandidateApplicationStatus,
        CandidateApplicationDigest: input.CandidateApplicationDigest,
        ObservationStatus:         input.ObservationStatus,
        ObservationSource:         input.ObservationSource,
        ObservationEvidenceDigest: input.ObservationEvidenceDigest,
        ReverseObservationDigest:  input.ReverseObservationDigest,
        ObservationDigest:          input.ObservationDigest,
        Publishable:                false,
        NonExecuting:              true,
        NonAuthorizing:            true,
    }
    if err := input.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationBridgeUnknown
        if output.MissingStage == "" {
            output.MissingStage = "candidate-application-observation-evidence"
        }
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "Candidate application observation is UNKNOWN; reverse evidence must be resolved"
        return output
    }
    output.Publishable = true
    output.Severity = "info"
    output.Code = "jev.external-apply.candidate-application-observation-recorded"
    output.Message = "Candidate application observation is bound to reverse evidence without execution or authorization"
    if input.ObservationStatus == jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMismatch {
        output.Severity = "warning"
        output.Code = "jev.external-apply.candidate-application-observation-mismatch"
        output.Message = "Candidate application observation mismatched reverse evidence without execution or authorization"
    }
    return output
}
