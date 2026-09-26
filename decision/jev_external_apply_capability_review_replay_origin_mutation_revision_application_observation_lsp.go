package decision

import "fmt"

type JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationLSPDiagnostic struct {
    Severity                  string
    Code                      string
    Message                   string
    Status                    string
    MissingStage              string
    CandidateStatus           string
    CandidateDigest           string
    ObservationStatus         string
    ObservationSource         string
    ObservationEvidenceDigest string
    ReverseObservationDigest  string
    ObservationDigest         string
    Publishable               bool
    NonExecuting              bool
    NonAuthorizing            bool
}

func (d JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationLSPDiagnostic) Validate() error {
    if d.Severity == "" || d.Code == "" || d.Message == "" || d.Status == "" {
        return fmt.Errorf("incomplete JEV revision application observation LSP diagnostic")
    }
    if !d.NonExecuting {
        return fmt.Errorf("JEV revision application observation LSP diagnostic must be non-executing")
    }
    if !d.NonAuthorizing {
        return fmt.Errorf("JEV revision application observation LSP diagnostic must be non-authorizing")
    }
    switch d.Status {
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationNotNeeded:
        if d.CandidateStatus != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateNotNeeded ||
            d.ObservationStatus != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationNotNeeded ||
            d.CandidateDigest == "" || d.ObservationDigest == "" {
            return fmt.Errorf("not-needed application observation LSP diagnostic is incomplete")
        }
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationRecorded:
        if d.CandidateStatus != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateReady ||
            d.ObservationStatus != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationObserved ||
            d.CandidateDigest == "" || d.ObservationSource == "" ||
            d.ObservationEvidenceDigest == "" || d.ReverseObservationDigest == "" ||
            d.ObservationDigest == "" {
            return fmt.Errorf("recorded application observation LSP diagnostic is incomplete")
        }
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationMismatch:
        if d.CandidateStatus != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateReady ||
            d.ObservationStatus != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationMismatched ||
            d.CandidateDigest == "" || d.ObservationSource == "" ||
            d.ObservationEvidenceDigest == "" || d.ReverseObservationDigest == "" ||
            d.ObservationDigest == "" {
            return fmt.Errorf("mismatched application observation LSP diagnostic is incomplete")
        }
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationHold:
        if d.CandidateStatus != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateHold ||
            d.ObservationStatus != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationHoldStatus ||
            d.CandidateDigest == "" || d.ObservationDigest == "" {
            return fmt.Errorf("held application observation LSP diagnostic is incomplete")
        }
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationRejected:
        if d.CandidateStatus != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateRejected ||
            d.ObservationStatus != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationRejectedStatus ||
            d.CandidateDigest == "" || d.ObservationDigest == "" {
            return fmt.Errorf("rejected application observation LSP diagnostic is incomplete")
        }
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationUnknown:
        if d.Publishable || d.MissingStage == "" {
            return fmt.Errorf("UNKNOWN application observation LSP diagnostic must remain non-publishable")
        }
    default:
        return fmt.Errorf("invalid JEV revision application observation LSP status")
    }
    return nil
}

func ProjectJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationLSP(input JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservation) JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationLSPDiagnostic {
    output := JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationLSPDiagnostic{
        Status:                    input.Status,
        MissingStage:              input.MissingStage,
        CandidateStatus:            input.CandidateStatus,
        CandidateDigest:            input.CandidateDigest,
        ObservationStatus:          input.ObservationStatus,
        ObservationSource:          input.ObservationSource,
        ObservationEvidenceDigest:  input.ObservationEvidenceDigest,
        ReverseObservationDigest:   input.ReverseObservationDigest,
        ObservationDigest:          input.ObservationDigest,
        NonExecuting:               input.NonExecuting,
        NonAuthorizing:             input.NonAuthorizing,
    }
    if err := input.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationUnknown
        if output.MissingStage == "" {
            output.MissingStage = "application-observation-evidence"
        }
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "JEV revision application observation is not publishable; evidence must be resolved"
        output.Publishable = false
        return output
    }
    output.Publishable = true
    switch input.Status {
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationNotNeeded:
        output.Severity = "info"
        output.Code = "jev.external-apply.application-observation-not-needed"
        output.Message = "No application observation was needed; no execution or authorization occurred"
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationRecorded:
        output.Severity = "info"
        output.Code = "jev.external-apply.application-observation-recorded"
        output.Message = "Application observation was recorded with reverse evidence; no execution or authorization occurred"
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationMismatch:
        output.Severity = "warning"
        output.Code = "jev.external-apply.application-observation-mismatch"
        output.Message = "Application observation mismatched reverse evidence; no execution or authorization occurred"
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationHold:
        output.Severity = "info"
        output.Code = "jev.external-apply.application-observation-hold"
        output.Message = "Application observation is held; no execution or authorization occurred"
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationRejected:
        output.Severity = "warning"
        output.Code = "jev.external-apply.application-observation-rejected"
        output.Message = "Application observation was rejected; no execution or authorization occurred"
    default:
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationUnknown
        output.MissingStage = "application-observation-status"
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "JEV revision application observation status is UNKNOWN; evidence must be resolved"
        output.Publishable = false
    }
    return output
}