package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "strconv"
    "strings"
)

const jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridgeLSPBoundCode = "jev.provenance.feedback-revision-proposal-bound"
const jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridgeLSPUnknownCode = "jev.provenance.unknown"

type JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridgeLSPDiagnostic struct {
    Severity                  string
    Code                      string
    Message                   string
    Status                    string
    MissingStage              string
    CandidateDigest           string
    CandidateSource           string
    CandidateGateDigest       string
    FeedbackStatus             string
    FeedbackDirection          string
    FeedbackDigest             string
    FeedbackSource             string
    FeedbackEvidenceDigest    string
    FeedbackBridgeDigest      string
    RevisionProposalStatus    string
    ProposalDecision           string
    ProposalTarget             string
    ProposalDigest             string
    ProposalSource             string
    ProposalEvidenceDigest    string
    BridgeDigest              string
    ProjectionDigest          string
    Publishable               bool
    NonExecuting              bool
    NonAuthorizing            bool
}

func (d JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridgeLSPDiagnostic) Validate() error {
    if d.Severity == "" || d.Code == "" || d.Message == "" || d.Status == "" ||
        d.ProjectionDigest == "" {
        return fmt.Errorf("incomplete JEV feedback revision proposal LSP diagnostic")
    }
    if !d.NonExecuting {
        return fmt.Errorf("JEV feedback revision proposal LSP diagnostic must be non-executing")
    }
    if !d.NonAuthorizing {
        return fmt.Errorf("JEV feedback revision proposal LSP diagnostic must be non-authorizing")
    }
    switch d.Status {
    case jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridgeUnknown:
        if d.Severity != "warning" || d.Code != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridgeLSPUnknownCode ||
            d.MissingStage == "" || d.Publishable {
            return fmt.Errorf("unknown feedback revision proposal LSP diagnostic is inconsistent")
        }
    case jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridgeBound:
        if d.Severity != "info" || d.Code != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridgeLSPBoundCode ||
            !d.Publishable || d.BridgeDigest == "" || d.FeedbackBridgeDigest == "" ||
            d.FeedbackStatus == "" || d.RevisionProposalStatus == "" {
            return fmt.Errorf("bound feedback revision proposal LSP diagnostic is incomplete")
        }
        if d.FeedbackStatus == jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBound &&
            (d.CandidateDigest == "" || d.CandidateSource == "" || d.CandidateGateDigest == "" ||
                d.FeedbackDirection == "" || d.FeedbackDigest == "" || d.FeedbackSource == "" ||
                d.FeedbackEvidenceDigest == "" ||
                d.ProposalDecision == "" || d.ProposalTarget == "" || d.ProposalDigest == "" ||
                d.ProposalSource == "" || d.ProposalEvidenceDigest == "") {
            return fmt.Errorf("bound feedback revision proposal LSP diagnostic lost evidence")
        }
    default:
        return fmt.Errorf("invalid feedback revision proposal LSP status")
    }
    expected := digestJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridgeLSPDiagnostic(
        d.Severity, d.Code, d.Message, d.Status, d.MissingStage,
        d.CandidateDigest, d.CandidateSource, d.CandidateGateDigest,
        d.FeedbackStatus, d.FeedbackDirection, d.FeedbackDigest, d.FeedbackSource,
        d.FeedbackEvidenceDigest, d.FeedbackBridgeDigest, d.RevisionProposalStatus,
        d.ProposalDecision, d.ProposalTarget, d.ProposalDigest, d.ProposalSource,
        d.ProposalEvidenceDigest, d.BridgeDigest, d.Publishable,
    )
    if d.ProjectionDigest != expected {
        return fmt.Errorf("JEV feedback revision proposal LSP projection digest mismatch")
    }
    return nil
}

func ProjectJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridgeLSP(input JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridge) JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridgeLSPDiagnostic {
    output := JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridgeLSPDiagnostic{
        Status:                 input.Status,
        MissingStage:           input.MissingStage,
        CandidateDigest:        input.CandidateDigest,
        CandidateSource:        input.CandidateSource,
        CandidateGateDigest:    input.CandidateGateDigest,
        FeedbackStatus:         input.FeedbackStatus,
        FeedbackDirection:      input.FeedbackDirection,
        FeedbackDigest:         input.FeedbackDigest,
        FeedbackSource:         input.FeedbackSource,
        FeedbackEvidenceDigest: input.FeedbackEvidenceDigest,
        FeedbackBridgeDigest:   input.FeedbackBridgeDigest,
        RevisionProposalStatus: input.RevisionProposalStatus,
        ProposalDecision:       input.ProposalDecision,
        ProposalTarget:         input.ProposalTarget,
        ProposalDigest:         input.ProposalDigest,
        ProposalSource:         input.ProposalSource,
        ProposalEvidenceDigest: input.ProposalEvidenceDigest,
        BridgeDigest:            input.BridgeDigest,
        NonExecuting:            true,
        NonAuthorizing:          true,
    }
    if input.Status == jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridgeUnknown || input.MissingStage != "" {
        output.Severity = "warning"
        output.Code = jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridgeLSPUnknownCode
        output.Message = "feedback revision proposal provenance is incomplete"
        output.Publishable = false
    } else {
        output.Severity = "info"
        output.Code = jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridgeLSPBoundCode
        output.Message = "feedback revision proposal provenance is bound"
        output.Publishable = true
    }
    output.ProjectionDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridgeLSPDiagnostic(
        output.Severity, output.Code, output.Message, output.Status, output.MissingStage,
        output.CandidateDigest, output.CandidateSource, output.CandidateGateDigest,
        output.FeedbackStatus, output.FeedbackDirection, output.FeedbackDigest, output.FeedbackSource,
        output.FeedbackEvidenceDigest, output.FeedbackBridgeDigest, output.RevisionProposalStatus,
        output.ProposalDecision, output.ProposalTarget, output.ProposalDigest, output.ProposalSource,
        output.ProposalEvidenceDigest, output.BridgeDigest, output.Publishable,
    )
    if err := output.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridgeUnknown
        output.MissingStage = "lsp-projection"
        output.Publishable = false
        output.ProjectionDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridgeLSPDiagnostic(
            output.Severity, output.Code, output.Message, output.Status, output.MissingStage,
            output.CandidateDigest, output.CandidateSource, output.CandidateGateDigest,
            output.FeedbackStatus, output.FeedbackDirection, output.FeedbackDigest, output.FeedbackSource,
            output.FeedbackEvidenceDigest, output.FeedbackBridgeDigest, output.RevisionProposalStatus,
            output.ProposalDecision, output.ProposalTarget, output.ProposalDigest, output.ProposalSource,
            output.ProposalEvidenceDigest, output.BridgeDigest, output.Publishable,
        )
    }
    return output
}

func digestJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridgeLSPDiagnostic(severity, code, message, status, missingStage, candidateDigest, candidateSource, candidateGateDigest, feedbackStatus, feedbackDirection, feedbackDigest, feedbackSource, feedbackEvidenceDigest, feedbackBridgeDigest, revisionProposalStatus, proposalDecision, proposalTarget, proposalDigest, proposalSource, proposalEvidenceDigest, bridgeDigest string, publishable bool) string {
    values := []string{
        severity, code, message, status, missingStage, candidateDigest, candidateSource,
        candidateGateDigest, feedbackStatus, feedbackDirection, feedbackDigest, feedbackSource,
        feedbackEvidenceDigest, feedbackBridgeDigest, revisionProposalStatus, proposalDecision,
        proposalTarget, proposalDigest, proposalSource, proposalEvidenceDigest, bridgeDigest,
        strconv.FormatBool(publishable),
    }
    sum := sha256.Sum256([]byte(strings.Join(values, "|")))
    return hex.EncodeToString(sum[:])
}
