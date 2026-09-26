package decision

import "fmt"

type JEVExternalApplyCapabilityReviewReplayOriginMutationObservationLSPDiagnostic struct {
    Severity                  string
    Code                      string
    Message                   string
    Status                    string
    MissingStage              string
    CapabilityStatus          string
    CapabilityDigest           string
    ObservedMutationDigest    string
    ReverseMutationDigest     string
    MutationObservationDigest string
    Publishable               bool
    NonExecuting              bool
    NonAuthorizing            bool
}

func (d JEVExternalApplyCapabilityReviewReplayOriginMutationObservationLSPDiagnostic) Validate() error {
    if d.Severity == "" || d.Code == "" || d.Message == "" || d.Status == "" {
        return fmt.Errorf("incomplete external apply capability review replay origin mutation observation LSP diagnostic")
    }
    if d.Status != jevExternalApplyCapabilityReviewReplayOriginMutationObservationNotNeeded &&
        d.Status != jevExternalApplyCapabilityReviewReplayOriginMutationObservationRecorded &&
        d.Status != jevExternalApplyCapabilityReviewReplayOriginMutationObservationHold &&
        d.Status != jevExternalApplyCapabilityReviewReplayOriginMutationObservationRejected {
        return fmt.Errorf("invalid mutation observation LSP status")
    }
    if d.CapabilityStatus != jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityNotNeeded &&
        d.CapabilityStatus != jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityRecorded &&
        d.CapabilityStatus != jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityHold &&
        d.CapabilityStatus != jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityRejected {
        return fmt.Errorf("invalid mutation observation LSP capability status")
    }
    if d.CapabilityStatus == jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityRecorded {
        if d.Status != jevExternalApplyCapabilityReviewReplayOriginMutationObservationRecorded ||
            d.ObservedMutationDigest == "" ||
            d.ReverseMutationDigest == "" {
            return fmt.Errorf("recorded mutation observation LSP diagnostic requires forward and reverse evidence")
        }
    }
    if d.CapabilityStatus == jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityNotNeeded &&
        d.Status != jevExternalApplyCapabilityReviewReplayOriginMutationObservationNotNeeded {
        return fmt.Errorf("not-needed capability LSP diagnostic requires observation-not-needed status")
    }
    if d.CapabilityStatus == jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityHold &&
        d.Status != jevExternalApplyCapabilityReviewReplayOriginMutationObservationHold {
        return fmt.Errorf("held capability LSP diagnostic requires observation-hold status")
    }
    if d.CapabilityStatus == jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityRejected &&
        d.Status != jevExternalApplyCapabilityReviewReplayOriginMutationObservationRejected {
        return fmt.Errorf("rejected capability LSP diagnostic requires observation-rejected status")
    }
    if d.CapabilityDigest == "" || d.MutationObservationDigest == "" {
        return fmt.Errorf("mutation observation LSP diagnostic requires digest evidence")
    }
    if !d.NonExecuting {
        return fmt.Errorf("external apply capability review replay origin mutation observation LSP diagnostic must be non-executing")
    }
    if !d.NonAuthorizing {
        return fmt.Errorf("external apply capability review replay origin mutation observation LSP diagnostic must be non-authorizing")
    }
    return nil
}

func ProjectJEVExternalApplyCapabilityReviewReplayOriginMutationObservationLSP(input JEVExternalApplyCapabilityReviewReplayOriginMutationObservation) JEVExternalApplyCapabilityReviewReplayOriginMutationObservationLSPDiagnostic {
    output := JEVExternalApplyCapabilityReviewReplayOriginMutationObservationLSPDiagnostic{
        Status:                    input.Status,
        MissingStage:              input.MissingStage,
        CapabilityStatus:           input.CapabilityStatus,
        CapabilityDigest:           input.CapabilityDigest,
        ObservedMutationDigest:    input.ObservedMutationDigest,
        ReverseMutationDigest:     input.ReverseMutationDigest,
        MutationObservationDigest: input.MutationObservationDigest,
        NonExecuting:              input.NonExecuting,
        NonAuthorizing:            input.NonAuthorizing,
    }
    if err := input.Validate(); err != nil {
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "JEV mutation observation is not publishable: forward or reverse evidence is incomplete"
        output.Publishable = false
        return output
    }
    output.Publishable = true
    switch input.Status {
    case jevExternalApplyCapabilityReviewReplayOriginMutationObservationNotNeeded:
        output.Severity = "info"
        output.Code = "jev.external-apply.mutation-observation-not-needed"
        output.Message = "Mutation observation was not needed; no execution or authorization occurred"
    case jevExternalApplyCapabilityReviewReplayOriginMutationObservationRecorded:
        output.Severity = "info"
        output.Code = "jev.external-apply.mutation-observation-recorded"
        output.Message = "Mutation observation includes forward and reverse evidence; no execution or authorization occurred"
    case jevExternalApplyCapabilityReviewReplayOriginMutationObservationHold:
        output.Severity = "info"
        output.Code = "jev.external-apply.mutation-observation-hold"
        output.Message = "Mutation observation is on hold; no execution or authorization occurred"
    case jevExternalApplyCapabilityReviewReplayOriginMutationObservationRejected:
        output.Severity = "warning"
        output.Code = "jev.external-apply.mutation-observation-rejected"
        output.Message = "Mutation observation is rejected; no execution or authorization occurred"
    default:
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "JEV mutation observation is UNKNOWN; evidence must be resolved"
        output.Publishable = false
    }
    return output
}
