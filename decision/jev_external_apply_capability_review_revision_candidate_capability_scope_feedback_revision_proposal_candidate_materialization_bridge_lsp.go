package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "strconv"
    "strings"
)

const jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridgeLSPBoundCode = "jev.provenance.revision-proposal-candidate-materialized"
const jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridgeLSPUnknownCode = "jev.provenance.unknown"

type JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridgeLSPDiagnostic struct {
    Severity                    string
    Code                        string
    Message                     string
    Status                      string
    MissingStage                string
    CandidateDigest             string
    CandidateSource             string
    CandidateGateDigest         string
    FeedbackDirection             string
    FeedbackDigest                string
    FeedbackSource                string
    FeedbackEvidenceDigest        string
    FeedbackBridgeDigest          string
    ProposalDecision              string
    ProposalTarget                string
    ProposalDigest                string
    ProposalSource                string
    ProposalEvidenceDigest        string
    RevisionProposalBridgeDigest string
    MaterializationStatus         string
    GeneratedCandidateDigest      string
    GeneratedCandidateSource      string
    GenerationInputDigest         string
    GenerationEvidenceDigest      string
    GeneratorIdentity             string
    BridgeDigest                  string
    ProjectionDigest              string
    Publishable                   bool
    NonExecuting                  bool
    NonAuthorizing                bool
}

func (d JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridgeLSPDiagnostic) Validate() error {
    if d.Severity == "" || d.Code == "" || d.Message == "" || d.Status == "" ||
        d.ProjectionDigest == "" {
        return fmt.Errorf("incomplete JEV candidate materialization LSP diagnostic")
    }
    if !d.NonExecuting {
        return fmt.Errorf("JEV candidate materialization LSP diagnostic must be non-executing")
    }
    if !d.NonAuthorizing {
        return fmt.Errorf("JEV candidate materialization LSP diagnostic must be non-authorizing")
    }
    switch d.Status {
    case jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridgeUnknown:
        if d.Severity != "warning" || d.Code != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridgeLSPUnknownCode ||
            d.MissingStage == "" || d.Publishable {
            return fmt.Errorf("unknown candidate materialization LSP diagnostic is inconsistent")
        }
    case jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridgeBound:
        if d.Severity != "info" || d.Code != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridgeLSPBoundCode ||
            !d.Publishable || d.BridgeDigest == "" || d.FeedbackBridgeDigest == "" ||
            d.RevisionProposalBridgeDigest == "" || d.MaterializationStatus == "" {
            return fmt.Errorf("bound candidate materialization LSP diagnostic is incomplete")
        }
        if d.ProposalDecision == jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalRevise &&
            (d.CandidateDigest == "" || d.CandidateSource == "" || d.CandidateGateDigest == "" ||
                d.FeedbackDirection == "" || d.FeedbackDigest == "" || d.FeedbackSource == "" ||
                d.FeedbackEvidenceDigest == "" || d.ProposalTarget == "" || d.ProposalDigest == "" ||
                d.ProposalSource == "" || d.ProposalEvidenceDigest == "" ||
                d.GeneratedCandidateDigest == "" || d.GeneratedCandidateSource == "" ||
                d.GenerationInputDigest == "" || d.GenerationEvidenceDigest == "" ||
                d.GeneratorIdentity == "") {
            return fmt.Errorf("revise candidate materialization LSP diagnostic lost evidence")
        }
    default:
        return fmt.Errorf("invalid candidate materialization LSP status")
    }
    expected := digestJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridgeLSPDiagnostic(
        d.Severity, d.Code, d.Message, d.Status, d.MissingStage,
        d.CandidateDigest, d.CandidateSource, d.CandidateGateDigest,
        d.FeedbackDirection, d.FeedbackDigest, d.FeedbackSource, d.FeedbackEvidenceDigest,
        d.FeedbackBridgeDigest, d.ProposalDecision, d.ProposalTarget, d.ProposalDigest,
        d.ProposalSource, d.ProposalEvidenceDigest, d.RevisionProposalBridgeDigest,
        d.MaterializationStatus, d.GeneratedCandidateDigest, d.GeneratedCandidateSource,
        d.GenerationInputDigest, d.GenerationEvidenceDigest, d.GeneratorIdentity,
        d.BridgeDigest, d.Publishable,
    )
    if d.ProjectionDigest != expected {
        return fmt.Errorf("JEV candidate materialization LSP projection digest mismatch")
    }
    return nil
}

func ProjectJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridgeLSP(input JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridge) JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridgeLSPDiagnostic {
    output := JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridgeLSPDiagnostic{
        Status:                      input.Status,
        MissingStage:                input.MissingStage,
        CandidateDigest:              input.CandidateDigest,
        CandidateSource:              input.CandidateSource,
        CandidateGateDigest:          input.CandidateGateDigest,
        FeedbackDirection:            input.FeedbackDirection,
        FeedbackDigest:               input.FeedbackDigest,
        FeedbackSource:               input.FeedbackSource,
        FeedbackEvidenceDigest:       input.FeedbackEvidenceDigest,
        FeedbackBridgeDigest:         input.FeedbackBridgeDigest,
        ProposalDecision:             input.ProposalDecision,
        ProposalTarget:               input.ProposalTarget,
        ProposalDigest:               input.ProposalDigest,
        ProposalSource:               input.ProposalSource,
        ProposalEvidenceDigest:       input.ProposalEvidenceDigest,
        RevisionProposalBridgeDigest: input.RevisionProposalBridgeDigest,
        MaterializationStatus:         input.MaterializationStatus,
        GeneratedCandidateDigest:      input.GeneratedCandidateDigest,
        GeneratedCandidateSource:      input.GeneratedCandidateSource,
        GenerationInputDigest:         input.GenerationInputDigest,
        GenerationEvidenceDigest:      input.GenerationEvidenceDigest,
        GeneratorIdentity:             input.GeneratorIdentity,
        BridgeDigest:                  input.BridgeDigest,
        NonExecuting:                  true,
        NonAuthorizing:                true,
    }
    if input.Status == jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridgeUnknown || input.MissingStage != "" {
        output.Severity = "warning"
        output.Code = jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridgeLSPUnknownCode
        output.Message = "candidate materialization provenance is incomplete"
        output.Publishable = false
    } else {
        output.Severity = "info"
        output.Code = jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridgeLSPBoundCode
        output.Message = "revision proposal candidate materialization provenance is bound"
        output.Publishable = true
    }
    output.ProjectionDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridgeLSPDiagnostic(
        output.Severity, output.Code, output.Message, output.Status, output.MissingStage,
        output.CandidateDigest, output.CandidateSource, output.CandidateGateDigest,
        output.FeedbackDirection, output.FeedbackDigest, output.FeedbackSource, output.FeedbackEvidenceDigest,
        output.FeedbackBridgeDigest, output.ProposalDecision, output.ProposalTarget, output.ProposalDigest,
        output.ProposalSource, output.ProposalEvidenceDigest, output.RevisionProposalBridgeDigest,
        output.MaterializationStatus, output.GeneratedCandidateDigest, output.GeneratedCandidateSource,
        output.GenerationInputDigest, output.GenerationEvidenceDigest, output.GeneratorIdentity,
        output.BridgeDigest, output.Publishable,
    )
    if err := output.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridgeUnknown
        output.MissingStage = "lsp-projection"
        output.Publishable = false
        output.ProjectionDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridgeLSPDiagnostic(
            output.Severity, output.Code, output.Message, output.Status, output.MissingStage,
            output.CandidateDigest, output.CandidateSource, output.CandidateGateDigest,
            output.FeedbackDirection, output.FeedbackDigest, output.FeedbackSource, output.FeedbackEvidenceDigest,
            output.FeedbackBridgeDigest, output.ProposalDecision, output.ProposalTarget, output.ProposalDigest,
            output.ProposalSource, output.ProposalEvidenceDigest, output.RevisionProposalBridgeDigest,
            output.MaterializationStatus, output.GeneratedCandidateDigest, output.GeneratedCandidateSource,
            output.GenerationInputDigest, output.GenerationEvidenceDigest, output.GeneratorIdentity,
            output.BridgeDigest, output.Publishable,
        )
    }
    return output
}

func digestJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridgeLSPDiagnostic(severity, code, message, status, missingStage, candidateDigest, candidateSource, candidateGateDigest, feedbackDirection, feedbackDigest, feedbackSource, feedbackEvidenceDigest, feedbackBridgeDigest, proposalDecision, proposalTarget, proposalDigest, proposalSource, proposalEvidenceDigest, revisionProposalBridgeDigest, materializationStatus, generatedCandidateDigest, generatedCandidateSource, generationInputDigest, generationEvidenceDigest, generatorIdentity, bridgeDigest string, publishable bool) string {
    values := []string{
        severity, code, message, status, missingStage, candidateDigest, candidateSource,
        candidateGateDigest, feedbackDirection, feedbackDigest, feedbackSource, feedbackEvidenceDigest,
        feedbackBridgeDigest, proposalDecision, proposalTarget, proposalDigest, proposalSource,
        proposalEvidenceDigest, revisionProposalBridgeDigest, materializationStatus,
        generatedCandidateDigest, generatedCandidateSource, generationInputDigest,
        generationEvidenceDigest, generatorIdentity, bridgeDigest, strconv.FormatBool(publishable),
    }
    sum := sha256.Sum256([]byte(strings.Join(values, "|")))
    return hex.EncodeToString(sum[:])
}
