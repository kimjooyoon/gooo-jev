package decision

import "fmt"

type JEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackCandidateGateBridgeLSPDiagnostic struct {
    Severity                  string
    Code                      string
    Message                   string
    Status                    string
    MissingStage              string
    GenerationFeedbackStatus  string
    GeneratedCandidateStatus  string
    CandidateDigest           string
    GenerationSource          string
    GenerationEvidenceDigest  string
    CandidateGateStatus       string
    CandidateDecision         string
    CandidateGateDigest       string
    BridgeDigest              string
    Publishable               bool
    NonExecuting              bool
    NonAuthorizing            bool
}

func (d JEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackCandidateGateBridgeLSPDiagnostic) Validate() error {
    if d.Severity == "" || d.Code == "" || d.Message == "" || d.Status == "" {
        return fmt.Errorf("incomplete JEV generation feedback candidate gate LSP diagnostic")
    }
    if !d.NonExecuting {
        return fmt.Errorf("JEV generation feedback candidate gate LSP diagnostic must be non-executing")
    }
    if !d.NonAuthorizing {
        return fmt.Errorf("JEV generation feedback candidate gate LSP diagnostic must be non-authorizing")
    }
    if d.Status == jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackCandidateGateBridgeUnknown {
        if d.Publishable || d.MissingStage == "" {
            return fmt.Errorf("UNKNOWN JEV generation feedback candidate gate diagnostic must remain non-publishable")
        }
        return nil
    }
    if d.Status != jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackCandidateGateBridgeBound {
        return fmt.Errorf("invalid JEV generation feedback candidate gate LSP status")
    }
    if !d.Publishable || d.GenerationFeedbackStatus == "" ||
        d.GeneratedCandidateStatus == "" || d.CandidateGateStatus == "" ||
        d.CandidateDecision == "" || d.CandidateGateDigest == "" ||
        d.BridgeDigest == "" {
        return fmt.Errorf("bound JEV generation feedback candidate gate diagnostic is incomplete")
    }
    switch d.GeneratedCandidateStatus {
    case jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackGenerated:
        if d.CandidateDecision != jevExternalApplyCapabilityReviewRevisionCandidateReady ||
            d.CandidateDigest == "" || d.GenerationSource == "" ||
            d.GenerationEvidenceDigest == "" {
            return fmt.Errorf("generated candidate gate diagnostic is incomplete")
        }
    case jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackHeld:
        if d.CandidateDecision != jevExternalApplyCapabilityReviewRevisionCandidateHold ||
            d.CandidateDigest != "" || d.GenerationSource != "" ||
            d.GenerationEvidenceDigest != "" {
            return fmt.Errorf("held candidate gate diagnostic has inconsistent evidence")
        }
    case jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackRejected:
        if d.CandidateDecision != jevExternalApplyCapabilityReviewRevisionCandidateRejected ||
            d.CandidateDigest != "" || d.GenerationSource != "" ||
            d.GenerationEvidenceDigest != "" {
            return fmt.Errorf("rejected candidate gate diagnostic has inconsistent evidence")
        }
    default:
        return fmt.Errorf("invalid generated candidate status")
    }
    return nil
}

func ProjectJEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackCandidateGateBridgeLSP(
    input JEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackCandidateGateBridge,
) JEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackCandidateGateBridgeLSPDiagnostic {
    output := JEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackCandidateGateBridgeLSPDiagnostic{
        Status:                   input.Status,
        MissingStage:             input.MissingStage,
        GenerationFeedbackStatus: input.GenerationFeedbackStatus,
        GeneratedCandidateStatus: input.GeneratedCandidateStatus,
        CandidateDigest:          input.CandidateDigest,
        GenerationSource:         input.GenerationSource,
        GenerationEvidenceDigest: input.GenerationEvidenceDigest,
        CandidateGateStatus:       input.CandidateGateStatus,
        CandidateDecision:         input.CandidateDecision,
        CandidateGateDigest:       input.CandidateGateDigest,
        BridgeDigest:              input.BridgeDigest,
        Publishable:               false,
        NonExecuting:              true,
        NonAuthorizing:            true,
    }
    if err := input.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackCandidateGateBridgeUnknown
        if output.MissingStage == "" {
            output.MissingStage = "generation-feedback-candidate-gate-evidence"
        }
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "Generation feedback candidate gate evidence is UNKNOWN; evidence must be resolved"
        return output
    }
    output.Publishable = true
    output.Severity = "info"
    output.Code = "jev.external-apply.generation-feedback-candidate-gate-bound"
    output.Message = "Generation feedback is consistent with the candidate gate without execution or authorization"
    if input.GeneratedCandidateStatus == jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackRejected {
        output.Severity = "warning"
        output.Code = "jev.external-apply.generation-feedback-candidate-gate-rejected"
        output.Message = "Generation feedback and candidate gate reject generation without execution or authorization"
    }
    return output
}
