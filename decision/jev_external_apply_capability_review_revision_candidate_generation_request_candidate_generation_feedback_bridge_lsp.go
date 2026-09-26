package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "strings"
)

const jevExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGenerationFeedbackLSPProjectionBound = "generation-request-candidate-generation-feedback-lsp-bound"

type JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGenerationFeedbackLSPProjection struct {
    Severity                         string
    Code                             string
    Message                          string
    Status                           string
    MissingStage                     string
    GenerationRequestStatus          string
    GenerationRequestDigest          string
    GenerationRequestSource          string
    GenerationRequestEvidenceDigest  string
    GeneratorIdentity                string
    ProposalDecision                 string
    Feedback                         JEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackLSPDiagnostic
    GenerationRequestBridgeDigest    string
    FeedbackBridgeDigest             string
    BridgeDigest                     string
    ProjectionDigest                 string
    Publishable                      bool
    NonExecuting                     bool
    NonAuthorizing                   bool
}

func (d JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGenerationFeedbackLSPProjection) Validate() error {
    if d.Severity == "" || d.Code == "" || d.Message == "" || d.Status == "" {
        return fmt.Errorf("incomplete JEV generation request candidate feedback LSP projection")
    }
    if !d.NonExecuting || !d.NonAuthorizing {
        return fmt.Errorf("JEV generation request candidate feedback LSP projection must be non-executing and non-authorizing")
    }
    if d.Status == jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown {
        if d.Publishable || d.MissingStage == "" || d.BridgeDigest != "" || d.ProjectionDigest != "" {
            return fmt.Errorf("UNKNOWN JEV generation request candidate feedback projection must remain non-publishable")
        }
        return nil
    }
    if d.Status != jevExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGenerationFeedbackBridgeBound ||
        !d.Publishable || d.MissingStage != "" || d.GenerationRequestStatus == "" ||
        d.GenerationRequestDigest == "" || d.GenerationRequestSource == "" ||
        d.GenerationRequestEvidenceDigest == "" || d.GeneratorIdentity == "" ||
        d.GenerationRequestBridgeDigest == "" || d.FeedbackBridgeDigest == "" ||
        d.BridgeDigest == "" || d.ProjectionDigest == "" {
        return fmt.Errorf("bound JEV generation request candidate feedback projection is incomplete")
    }
    if err := d.Feedback.Validate(); err != nil {
        return fmt.Errorf("invalid candidate generation feedback projection: %w", err)
    }
    if d.Feedback.BridgeDigest != d.FeedbackBridgeDigest || !d.Feedback.Publishable {
        return fmt.Errorf("candidate generation feedback projection digest is inconsistent")
    }
    expected := digestJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGenerationFeedbackLSPProjection(
        d.Status, d.GenerationRequestBridgeDigest, d.FeedbackBridgeDigest,
        d.BridgeDigest, d.ProposalDecision,
    )
    if d.ProjectionDigest != expected {
        return fmt.Errorf("JEV generation request candidate feedback projection digest mismatch")
    }
    switch d.ProposalDecision {
    case "revise":
        if d.Feedback.Code != "jev.external-apply.candidate-generated" {
            return fmt.Errorf("revise generation request projection has inconsistent feedback code")
        }
    case "retain":
        if d.Feedback.Code != "jev.external-apply.candidate-generation-held" {
            return fmt.Errorf("retain generation request projection has inconsistent feedback code")
        }
    case "rollback":
        if d.Feedback.Code != "jev.external-apply.candidate-generation-rejected" {
            return fmt.Errorf("rollback generation request projection has inconsistent feedback code")
        }
    default:
        return fmt.Errorf("invalid generation request proposal decision in LSP projection")
    }
    return nil
}

func ProjectJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGenerationFeedbackLSPProjection(
    input JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGenerationFeedbackBridge,
) JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGenerationFeedbackLSPProjection {
    unknown := func(stage string) JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGenerationFeedbackLSPProjection {
        if stage == "" {
            stage = "generation-request-candidate-generation-feedback-lsp"
        }
        return JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGenerationFeedbackLSPProjection{
            Severity:       "error",
            Code:           "jev.provenance.unknown",
            Message:        "Generation request candidate feedback is UNKNOWN; evidence must be resolved",
            Status:         jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown,
            MissingStage:   stage,
            NonExecuting:   true,
            NonAuthorizing: true,
        }
    }
    if input.Status == jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown {
        return unknown(input.MissingStage)
    }
    if err := input.Validate(); err != nil {
        return unknown("generation-request-candidate-generation-feedback")
    }
    feedback := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackLSP(input.Feedback)
    output := JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGenerationFeedbackLSPProjection{
        Status:                        input.Status,
        GenerationRequestStatus:       input.GenerationRequestStatus,
        GenerationRequestDigest:       input.GenerationRequestDigest,
        GenerationRequestSource:       input.GenerationRequestSource,
        GenerationRequestEvidenceDigest: input.GenerationRequestEvidenceDigest,
        GeneratorIdentity:             input.GeneratorIdentity,
        ProposalDecision:              input.ProposalDecision,
        Feedback:                      feedback,
        GenerationRequestBridgeDigest: input.GenerationRequestBridgeDigest,
        FeedbackBridgeDigest:          input.FeedbackBridgeDigest,
        BridgeDigest:                  input.BridgeDigest,
        Publishable:                   false,
        NonExecuting:                  true,
        NonAuthorizing:                true,
    }
    if err := feedback.Validate(); err != nil {
        return unknown("candidate-generation-feedback-lsp")
    }
    output.Publishable = true
    output.Severity = "info"
    output.Code = "jev.external-apply.generation-request-candidate-generated"
    output.Message = "Generation request is projected to candidate generation feedback without execution or authorization"
    if input.ProposalDecision == "retain" {
        output.Code = "jev.external-apply.generation-request-candidate-held"
        output.Message = "Generation request is held and projected without execution or authorization"
    }
    if input.ProposalDecision == "rollback" {
        output.Severity = "warning"
        output.Code = "jev.external-apply.generation-request-candidate-rejected"
        output.Message = "Generation request is rejected and projected without execution or authorization"
    }
    output.ProjectionDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGenerationFeedbackLSPProjection(
        output.Status, output.GenerationRequestBridgeDigest, output.FeedbackBridgeDigest,
        output.BridgeDigest, output.ProposalDecision,
    )
    if err := output.Validate(); err != nil {
        return unknown("generation-request-candidate-generation-feedback-lsp")
    }
    return output
}

func digestJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGenerationFeedbackLSPProjection(
    status, generationRequestBridgeDigest, feedbackBridgeDigest, bridgeDigest, proposalDecision string,
) string {
    sum := sha256.Sum256([]byte(strings.Join([]string{
        status, generationRequestBridgeDigest, feedbackBridgeDigest, bridgeDigest, proposalDecision,
    }, "|")))
    return hex.EncodeToString(sum[:])
}
