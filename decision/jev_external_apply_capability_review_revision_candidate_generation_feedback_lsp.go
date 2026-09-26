package decision

import "fmt"

type JEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackLSPDiagnostic struct {
    Severity                  string
    Code                      string
    Message                   string
    Status                    string
    MissingStage              string
    DirectionStatus           string
    Direction                 string
    Target                    string
    FeedbackDigest            string
    DirectionDigest           string
    GeneratedCandidateStatus  string
    CandidateDigest           string
    GenerationSource          string
    GenerationEvidenceDigest  string
    BridgeDigest              string
    Publishable               bool
    NonExecuting              bool
    NonAuthorizing            bool
}

func (d JEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackLSPDiagnostic) Validate() error {
    if d.Severity == "" || d.Code == "" || d.Message == "" || d.Status == "" {
        return fmt.Errorf("incomplete JEV candidate generation feedback LSP diagnostic")
    }
    if !d.NonExecuting {
        return fmt.Errorf("JEV candidate generation feedback LSP diagnostic must be non-executing")
    }
    if !d.NonAuthorizing {
        return fmt.Errorf("JEV candidate generation feedback LSP diagnostic must be non-authorizing")
    }
    if d.Status == jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackUnknown {
        if d.Publishable || d.MissingStage == "" {
            return fmt.Errorf("UNKNOWN JEV candidate generation feedback diagnostic must remain non-publishable")
        }
        return nil
    }
    if d.Status != jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackBound {
        return fmt.Errorf("invalid JEV candidate generation feedback LSP status")
    }
    if !d.Publishable || d.DirectionStatus != jevExternalApplyCapabilityReviewRevisionCandidateObservationMetricDirectionBridgeBound ||
        d.Direction == "" || d.Target == "" || d.FeedbackDigest == "" ||
        d.DirectionDigest == "" || d.GeneratedCandidateStatus == "" ||
        d.BridgeDigest == "" {
        return fmt.Errorf("bound JEV candidate generation feedback diagnostic is incomplete")
    }
    switch d.Direction {
    case jevExternalApplyCapabilityReviewImprovementDirectionGenerate:
        if d.GeneratedCandidateStatus != jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackGenerated ||
            d.CandidateDigest == "" || d.GenerationSource == "" ||
            d.GenerationEvidenceDigest == "" {
            return fmt.Errorf("generated candidate feedback diagnostic is incomplete")
        }
    case jevExternalApplyCapabilityReviewImprovementDirectionHold:
        if d.GeneratedCandidateStatus != jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackHeld ||
            d.CandidateDigest != "" || d.GenerationSource != "" ||
            d.GenerationEvidenceDigest != "" {
            return fmt.Errorf("held candidate feedback diagnostic has inconsistent evidence")
        }
    case jevExternalApplyCapabilityReviewImprovementDirectionReject:
        if d.GeneratedCandidateStatus != jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackRejected ||
            d.CandidateDigest != "" || d.GenerationSource != "" ||
            d.GenerationEvidenceDigest != "" {
            return fmt.Errorf("rejected candidate feedback diagnostic has inconsistent evidence")
        }
    default:
        return fmt.Errorf("invalid candidate generation direction")
    }
    return nil
}

func ProjectJEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackLSP(
    input JEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedback,
) JEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackLSPDiagnostic {
    output := JEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackLSPDiagnostic{
        Status:                   input.Status,
        MissingStage:             input.MissingStage,
        DirectionStatus:           input.DirectionStatus,
        Direction:                input.Direction,
        Target:                   input.Target,
        FeedbackDigest:           input.FeedbackDigest,
        DirectionDigest:          input.DirectionDigest,
        GeneratedCandidateStatus: input.GeneratedCandidateStatus,
        CandidateDigest:          input.CandidateDigest,
        GenerationSource:          input.GenerationSource,
        GenerationEvidenceDigest: input.GenerationEvidenceDigest,
        BridgeDigest:              input.BridgeDigest,
        Publishable:              false,
        NonExecuting:             true,
        NonAuthorizing:           true,
    }
    if err := input.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackUnknown
        if output.MissingStage == "" {
            output.MissingStage = "candidate-generation-feedback-evidence"
        }
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "Candidate generation feedback is UNKNOWN; evidence must be resolved"
        return output
    }
    output.Publishable = true
    output.Severity = "info"
    output.Code = "jev.external-apply.candidate-generated"
    output.Message = "Candidate generation feedback is bound without execution or authorization"
    if input.GeneratedCandidateStatus == jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackHeld {
        output.Code = "jev.external-apply.candidate-generation-held"
        output.Message = "Candidate generation is held without execution or authorization"
    }
    if input.GeneratedCandidateStatus == jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackRejected {
        output.Severity = "warning"
        output.Code = "jev.external-apply.candidate-generation-rejected"
        output.Message = "Candidate generation is rejected without execution or authorization"
    }
    return output
}
