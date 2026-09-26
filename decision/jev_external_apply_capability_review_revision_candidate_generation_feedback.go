package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

const jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackBound = "candidate-generation-feedback-bound"
const jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackUnknown = "UNKNOWN"
const jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackGenerated = "candidate-generated"
const jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackHeld = "candidate-generation-held"
const jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackRejected = "candidate-generation-rejected"

type JEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackInput struct {
    Direction                 JEVExternalApplyCapabilityReviewRevisionCandidateObservationMetricDirectionBridge
    CandidateDigest           string
    GenerationSource          string
    GenerationEvidenceDigest string
    NonAuthorizing            bool
}

type JEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedback struct {
    Status                   string
    MissingStage             string
    DirectionStatus          string
    Direction                string
    Target                   string
    FeedbackDigest           string
    DirectionDigest          string
    GeneratedCandidateStatus string
    CandidateDigest          string
    GenerationSource         string
    GenerationEvidenceDigest string
    BridgeDigest             string
    NonExecuting             bool
    NonAuthorizing           bool
}

func (f JEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedback) Validate() error {
    if f.Status == "" || f.DirectionStatus == "" || f.Direction == "" ||
        f.Target == "" || f.FeedbackDigest == "" || f.DirectionDigest == "" ||
        f.GeneratedCandidateStatus == "" || f.BridgeDigest == "" {
        return fmt.Errorf("incomplete JEV candidate generation feedback")
    }
    if f.Status != jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackBound {
        return fmt.Errorf("invalid JEV candidate generation feedback status")
    }
    if f.DirectionStatus != jevExternalApplyCapabilityReviewRevisionCandidateObservationMetricDirectionBridgeBound {
        return fmt.Errorf("invalid JEV candidate direction status")
    }
    switch f.Direction {
    case jevExternalApplyCapabilityReviewImprovementDirectionGenerate:
        if f.GeneratedCandidateStatus != jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackGenerated ||
            f.CandidateDigest == "" || f.GenerationSource == "" ||
            f.GenerationEvidenceDigest == "" {
            return fmt.Errorf("generated candidate feedback is incomplete")
        }
    case jevExternalApplyCapabilityReviewImprovementDirectionHold:
        if f.GeneratedCandidateStatus != jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackHeld ||
            f.CandidateDigest != "" || f.GenerationSource != "" ||
            f.GenerationEvidenceDigest != "" {
            return fmt.Errorf("held candidate feedback has inconsistent evidence")
        }
    case jevExternalApplyCapabilityReviewImprovementDirectionReject:
        if f.GeneratedCandidateStatus != jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackRejected ||
            f.CandidateDigest != "" || f.GenerationSource != "" ||
            f.GenerationEvidenceDigest != "" {
            return fmt.Errorf("rejected candidate feedback has inconsistent evidence")
        }
    default:
        return fmt.Errorf("invalid JEV candidate generation direction")
    }
    if !f.NonExecuting {
        return fmt.Errorf("JEV candidate generation feedback must be non-executing")
    }
    if !f.NonAuthorizing {
        return fmt.Errorf("JEV candidate generation feedback must be non-authorizing")
    }
    expected := digestJEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedback(
        f.Status,
        f.DirectionStatus,
        f.Direction,
        f.Target,
        f.FeedbackDigest,
        f.DirectionDigest,
        f.GeneratedCandidateStatus,
        f.CandidateDigest,
        f.GenerationSource,
        f.GenerationEvidenceDigest,
    )
    if f.BridgeDigest != expected {
        return fmt.Errorf("JEV candidate generation feedback digest mismatch")
    }
    return nil
}

func BindJEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedback(input JEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackInput) JEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedback {
    output := JEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedback{
        Status:         jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackUnknown,
        NonExecuting:   true,
        NonAuthorizing: true,
    }
    if !input.NonAuthorizing || !input.Direction.NonAuthorizing {
        output.MissingStage = "authorization-boundary"
        return output
    }
    if !input.Direction.NonExecuting {
        output.MissingStage = "execution-boundary"
        return output
    }
    if err := input.Direction.Validate(); err != nil {
        output.MissingStage = "candidate-observation-metric-direction"
        return output
    }
    output.Status = jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackBound
    output.DirectionStatus = input.Direction.Status
    output.Direction = input.Direction.Direction
    output.Target = input.Direction.Target
    output.FeedbackDigest = input.Direction.FeedbackDigest
    output.DirectionDigest = input.Direction.DirectionDigest
    switch input.Direction.Direction {
    case jevExternalApplyCapabilityReviewImprovementDirectionGenerate:
        if input.CandidateDigest == "" {
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackUnknown
            output.MissingStage = "candidate-digest"
            return output
        }
        if input.GenerationSource == "" {
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackUnknown
            output.MissingStage = "generation-source"
            return output
        }
        if input.GenerationEvidenceDigest == "" {
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackUnknown
            output.MissingStage = "generation-evidence"
            return output
        }
        output.GeneratedCandidateStatus = jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackGenerated
        output.CandidateDigest = input.CandidateDigest
        output.GenerationSource = input.GenerationSource
        output.GenerationEvidenceDigest = input.GenerationEvidenceDigest
    case jevExternalApplyCapabilityReviewImprovementDirectionHold:
        output.GeneratedCandidateStatus = jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackHeld
    case jevExternalApplyCapabilityReviewImprovementDirectionReject:
        output.GeneratedCandidateStatus = jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackRejected
    default:
        output.Status = jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackUnknown
        output.MissingStage = "improvement-direction"
        return output
    }
    output.BridgeDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedback(
        output.Status,
        output.DirectionStatus,
        output.Direction,
        output.Target,
        output.FeedbackDigest,
        output.DirectionDigest,
        output.GeneratedCandidateStatus,
        output.CandidateDigest,
        output.GenerationSource,
        output.GenerationEvidenceDigest,
    )
    if err := output.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackUnknown
        output.MissingStage = "candidate-generation-feedback-evidence"
        output.BridgeDigest = ""
    }
    return output
}

func digestJEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedback(status, directionStatus, direction, target, feedbackDigest, directionDigest, generatedCandidateStatus, candidateDigest, generationSource, generationEvidenceDigest string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%s", status, directionStatus, direction, target, feedbackDigest, directionDigest, generatedCandidateStatus, candidateDigest, generationSource, generationEvidenceDigest)))
    return hex.EncodeToString(sum[:])
}
