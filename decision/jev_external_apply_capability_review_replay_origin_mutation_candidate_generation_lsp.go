package decision

import "fmt"

type JEVExternalApplyCapabilityReviewReplayOriginMutationCandidateGenerationLSPDiagnostic struct {
    Severity                  string
    Code                      string
    Message                   string
    Status                    string
    MissingStage              string
    ObservationStatus         string
    CandidateDigest           string
    CandidateSource           string
    ObservationDigest         string
    CandidateGenerationDigest string
    Publishable               bool
    NonExecuting              bool
    NonAuthorizing            bool
}

func (d JEVExternalApplyCapabilityReviewReplayOriginMutationCandidateGenerationLSPDiagnostic) Validate() error {
    if d.Severity == "" || d.Code == "" || d.Message == "" || d.Status == "" {
        return fmt.Errorf("incomplete external apply capability review replay origin mutation candidate generation LSP diagnostic")
    }
    if d.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateNotNeeded &&
        d.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateGenerated &&
        d.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateHold &&
        d.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRejected {
        return fmt.Errorf("invalid mutation candidate generation LSP status")
    }
    if d.ObservationStatus != jevExternalApplyCapabilityReviewReplayOriginMutationObservationNotNeeded &&
        d.ObservationStatus != jevExternalApplyCapabilityReviewReplayOriginMutationObservationRecorded &&
        d.ObservationStatus != jevExternalApplyCapabilityReviewReplayOriginMutationObservationHold &&
        d.ObservationStatus != jevExternalApplyCapabilityReviewReplayOriginMutationObservationRejected {
        return fmt.Errorf("invalid mutation candidate generation LSP observation status")
    }
    if d.ObservationStatus == jevExternalApplyCapabilityReviewReplayOriginMutationObservationRecorded {
        if d.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateGenerated ||
            d.CandidateDigest == "" ||
            d.CandidateSource == "" {
            return fmt.Errorf("generated candidate LSP diagnostic requires candidate digest and source")
        }
    }
    if d.ObservationStatus == jevExternalApplyCapabilityReviewReplayOriginMutationObservationNotNeeded &&
        d.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateNotNeeded {
        return fmt.Errorf("not-needed observation LSP diagnostic requires candidate-not-needed status")
    }
    if d.ObservationStatus == jevExternalApplyCapabilityReviewReplayOriginMutationObservationHold &&
        d.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateHold {
        return fmt.Errorf("held observation LSP diagnostic requires candidate-hold status")
    }
    if d.ObservationStatus == jevExternalApplyCapabilityReviewReplayOriginMutationObservationRejected &&
        d.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRejected {
        return fmt.Errorf("rejected observation LSP diagnostic requires candidate-rejected status")
    }
    if d.ObservationDigest == "" || d.CandidateGenerationDigest == "" {
        return fmt.Errorf("candidate generation LSP diagnostic requires digest evidence")
    }
    if !d.NonExecuting {
        return fmt.Errorf("external apply capability review replay origin mutation candidate generation LSP diagnostic must be non-executing")
    }
    if !d.NonAuthorizing {
        return fmt.Errorf("external apply capability review replay origin mutation candidate generation LSP diagnostic must be non-authorizing")
    }
    return nil
}

func ProjectJEVExternalApplyCapabilityReviewReplayOriginMutationCandidateGenerationLSP(input JEVExternalApplyCapabilityReviewReplayOriginMutationCandidateGeneration) JEVExternalApplyCapabilityReviewReplayOriginMutationCandidateGenerationLSPDiagnostic {
    output := JEVExternalApplyCapabilityReviewReplayOriginMutationCandidateGenerationLSPDiagnostic{
        Status:                    input.Status,
        MissingStage:              input.MissingStage,
        ObservationStatus:         input.ObservationStatus,
        CandidateDigest:            input.CandidateDigest,
        CandidateSource:            input.CandidateSource,
        ObservationDigest:          input.ObservationDigest,
        CandidateGenerationDigest: input.CandidateGenerationDigest,
        NonExecuting:              input.NonExecuting,
        NonAuthorizing:            input.NonAuthorizing,
    }
    if err := input.Validate(); err != nil {
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "JEV candidate generation is not publishable: candidate or observation evidence is incomplete"
        output.Publishable = false
        return output
    }
    output.Publishable = true
    switch input.Status {
    case jevExternalApplyCapabilityReviewReplayOriginMutationCandidateNotNeeded:
        output.Severity = "info"
        output.Code = "jev.external-apply.mutation-candidate-not-needed"
        output.Message = "No revision candidate was needed; no execution or authorization occurred"
    case jevExternalApplyCapabilityReviewReplayOriginMutationCandidateGenerated:
        output.Severity = "info"
        output.Code = "jev.external-apply.mutation-candidate-generated"
        output.Message = "Revision candidate evidence was generated for review; no execution or authorization occurred"
    case jevExternalApplyCapabilityReviewReplayOriginMutationCandidateHold:
        output.Severity = "info"
        output.Code = "jev.external-apply.mutation-candidate-hold"
        output.Message = "Revision candidate generation is on hold; no execution or authorization occurred"
    case jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRejected:
        output.Severity = "warning"
        output.Code = "jev.external-apply.mutation-candidate-rejected"
        output.Message = "Revision candidate generation was rejected; no execution or authorization occurred"
    default:
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "JEV candidate generation is UNKNOWN; evidence must be resolved"
        output.Publishable = false
    }
    return output
}
