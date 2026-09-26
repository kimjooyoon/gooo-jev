package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

const jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackCandidateGateBridgeBound = "generation-feedback-candidate-gate-bound"
const jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackCandidateGateBridgeUnknown = "UNKNOWN"

type JEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackCandidateGateBridgeInput struct {
    GenerationFeedback JEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedback
    CandidateGate      JEVExternalApplyCapabilityReviewRevisionCandidateGate
    NonAuthorizing     bool
}

type JEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackCandidateGateBridge struct {
    Status                   string
    MissingStage             string
    GenerationFeedbackStatus string
    GeneratedCandidateStatus string
    CandidateDigest          string
    GenerationSource         string
    GenerationEvidenceDigest string
    CandidateGateStatus      string
    CandidateDecision        string
    CandidateGateDigest      string
    BridgeDigest             string
    NonExecuting             bool
    NonAuthorizing           bool
}

func (b JEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackCandidateGateBridge) Validate() error {
    if b.Status == "" || b.GenerationFeedbackStatus == "" ||
        b.GeneratedCandidateStatus == "" || b.CandidateGateStatus == "" ||
        b.CandidateDecision == "" || b.CandidateGateDigest == "" ||
        b.BridgeDigest == "" {
        return fmt.Errorf("incomplete JEV generation feedback candidate gate bridge")
    }
    if b.Status != jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackCandidateGateBridgeBound {
        return fmt.Errorf("invalid JEV generation feedback candidate gate bridge status")
    }
    if b.GenerationFeedbackStatus != jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackBound {
        return fmt.Errorf("invalid JEV generation feedback status")
    }
    if b.CandidateGateStatus != jevExternalApplyCapabilityReviewRevisionCandidateGated {
        return fmt.Errorf("invalid JEV candidate gate status")
    }
    switch b.GeneratedCandidateStatus {
    case jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackGenerated:
        if b.CandidateDecision != jevExternalApplyCapabilityReviewRevisionCandidateReady ||
            b.CandidateDigest == "" || b.GenerationSource == "" ||
            b.GenerationEvidenceDigest == "" {
            return fmt.Errorf("generated candidate gate bridge is incomplete")
        }
    case jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackHeld:
        if b.CandidateDecision != jevExternalApplyCapabilityReviewRevisionCandidateHold ||
            b.CandidateDigest != "" || b.GenerationSource != "" ||
            b.GenerationEvidenceDigest != "" {
            return fmt.Errorf("held candidate gate bridge has inconsistent evidence")
        }
    case jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackRejected:
        if b.CandidateDecision != jevExternalApplyCapabilityReviewRevisionCandidateRejected ||
            b.CandidateDigest != "" || b.GenerationSource != "" ||
            b.GenerationEvidenceDigest != "" {
            return fmt.Errorf("rejected candidate gate bridge has inconsistent evidence")
        }
    default:
        return fmt.Errorf("invalid generated candidate status")
    }
    if !b.NonExecuting {
        return fmt.Errorf("JEV generation feedback candidate gate bridge must be non-executing")
    }
    if !b.NonAuthorizing {
        return fmt.Errorf("JEV generation feedback candidate gate bridge must be non-authorizing")
    }
    expected := digestJEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackCandidateGateBridge(
        b.Status,
        b.GenerationFeedbackStatus,
        b.GeneratedCandidateStatus,
        b.CandidateDigest,
        b.GenerationSource,
        b.GenerationEvidenceDigest,
        b.CandidateGateStatus,
        b.CandidateDecision,
        b.CandidateGateDigest,
    )
    if b.BridgeDigest != expected {
        return fmt.Errorf("JEV generation feedback candidate gate bridge digest mismatch")
    }
    return nil
}

func BindJEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackCandidateGateBridge(input JEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackCandidateGateBridgeInput) JEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackCandidateGateBridge {
    output := JEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackCandidateGateBridge{
        Status:         jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackCandidateGateBridgeUnknown,
        NonExecuting:   true,
        NonAuthorizing: true,
    }
    if !input.NonAuthorizing || !input.GenerationFeedback.NonAuthorizing || !input.CandidateGate.NonAuthorizing {
        output.MissingStage = "authorization-boundary"
        return output
    }
    if !input.GenerationFeedback.NonExecuting || !input.CandidateGate.NonExecuting {
        output.MissingStage = "execution-boundary"
        return output
    }
    if err := input.GenerationFeedback.Validate(); err != nil {
        output.MissingStage = "candidate-generation-feedback"
        return output
    }
    if err := input.CandidateGate.Validate(); err != nil {
        output.MissingStage = "revision-candidate-gate"
        return output
    }
    output.Status = jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackCandidateGateBridgeBound
    output.GenerationFeedbackStatus = input.GenerationFeedback.Status
    output.GeneratedCandidateStatus = input.GenerationFeedback.GeneratedCandidateStatus
    output.CandidateGateStatus = input.CandidateGate.Status
    output.CandidateDecision = input.CandidateGate.Decision
    output.CandidateGateDigest = input.CandidateGate.GateDigest
    switch input.GenerationFeedback.GeneratedCandidateStatus {
    case jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackGenerated:
        if input.CandidateGate.Decision != jevExternalApplyCapabilityReviewRevisionCandidateReady {
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackCandidateGateBridgeUnknown
            output.MissingStage = "generation-candidate-consistency"
            return output
        }
        if input.GenerationFeedback.CandidateDigest != input.CandidateGate.CandidateDigest {
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackCandidateGateBridgeUnknown
            output.MissingStage = "candidate-digest-consistency"
            return output
        }
        if input.GenerationFeedback.GenerationSource != input.CandidateGate.CandidateSource {
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackCandidateGateBridgeUnknown
            output.MissingStage = "generation-source-consistency"
            return output
        }
        output.CandidateDigest = input.GenerationFeedback.CandidateDigest
        output.GenerationSource = input.GenerationFeedback.GenerationSource
        output.GenerationEvidenceDigest = input.GenerationFeedback.GenerationEvidenceDigest
    case jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackHeld:
        if input.CandidateGate.Decision != jevExternalApplyCapabilityReviewRevisionCandidateHold {
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackCandidateGateBridgeUnknown
            output.MissingStage = "generation-candidate-consistency"
            return output
        }
    case jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackRejected:
        if input.CandidateGate.Decision != jevExternalApplyCapabilityReviewRevisionCandidateRejected {
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackCandidateGateBridgeUnknown
            output.MissingStage = "generation-candidate-consistency"
            return output
        }
    default:
        output.Status = jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackCandidateGateBridgeUnknown
        output.MissingStage = "generated-candidate-status"
        return output
    }
    output.BridgeDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackCandidateGateBridge(
        output.Status,
        output.GenerationFeedbackStatus,
        output.GeneratedCandidateStatus,
        output.CandidateDigest,
        output.GenerationSource,
        output.GenerationEvidenceDigest,
        output.CandidateGateStatus,
        output.CandidateDecision,
        output.CandidateGateDigest,
    )
    if err := output.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackCandidateGateBridgeUnknown
        output.MissingStage = "generation-feedback-candidate-gate-evidence"
        output.BridgeDigest = ""
    }
    return output
}

func digestJEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackCandidateGateBridge(status, generationFeedbackStatus, generatedCandidateStatus, candidateDigest, generationSource, generationEvidenceDigest, candidateGateStatus, candidateDecision, candidateGateDigest string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s", status, generationFeedbackStatus, generatedCandidateStatus, candidateDigest, generationSource, generationEvidenceDigest, candidateGateStatus, candidateDecision, candidateGateDigest)))
    return hex.EncodeToString(sum[:])
}
