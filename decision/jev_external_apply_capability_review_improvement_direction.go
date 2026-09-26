package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

const jevExternalApplyCapabilityReviewImprovementDirectionBound = "capability-review-improvement-direction-bound"
const jevExternalApplyCapabilityReviewImprovementDirectionGenerate = "candidate-generation"
const jevExternalApplyCapabilityReviewImprovementDirectionHold = "candidate-hold"
const jevExternalApplyCapabilityReviewImprovementDirectionReject = "candidate-rejection"
const jevExternalApplyCapabilityReviewImprovementDirectionUnknown = "UNKNOWN"

type JEVExternalApplyCapabilityReviewImprovementDirectionInput struct {
    Feedback       JEVExternalApplyCapabilityReviewFeedbackBridge
    Target         string
    CandidateSource string
    NonAuthorizing bool
}

type JEVExternalApplyCapabilityReviewImprovementDirection struct {
    Status           string
    MissingStage     string
    Direction        string
    Target           string
    CandidateSource  string
    FeedbackDigest   string
    DirectionDigest  string
    NonExecuting     bool
    NonAuthorizing   bool
}

func (d JEVExternalApplyCapabilityReviewImprovementDirection) Validate() error {
    if d.Status == "" || d.Direction == "" || d.Target == "" || d.CandidateSource == "" || d.FeedbackDigest == "" || d.DirectionDigest == "" {
        return fmt.Errorf("incomplete JEV external apply capability review improvement direction")
    }
    if d.Status != jevExternalApplyCapabilityReviewImprovementDirectionBound {
        return fmt.Errorf("invalid JEV external apply capability review improvement direction status")
    }
    if d.Direction != jevExternalApplyCapabilityReviewImprovementDirectionGenerate && d.Direction != jevExternalApplyCapabilityReviewImprovementDirectionHold && d.Direction != jevExternalApplyCapabilityReviewImprovementDirectionReject {
        return fmt.Errorf("invalid JEV external apply capability review improvement direction")
    }
    if !d.NonExecuting {
        return fmt.Errorf("JEV external apply capability review improvement direction must be non-executing")
    }
    if !d.NonAuthorizing {
        return fmt.Errorf("JEV external apply capability review improvement direction must be non-authorizing")
    }
    expected := digestJEVExternalApplyCapabilityReviewImprovementDirection(d.Status, d.Direction, d.Target, d.CandidateSource, d.FeedbackDigest)
    if d.DirectionDigest != expected {
        return fmt.Errorf("JEV external apply capability review improvement direction digest mismatch")
    }
    return nil
}

func DeriveJEVExternalApplyCapabilityReviewImprovementDirection(input JEVExternalApplyCapabilityReviewImprovementDirectionInput) JEVExternalApplyCapabilityReviewImprovementDirection {
    output := JEVExternalApplyCapabilityReviewImprovementDirection{
        Status:         jevExternalApplyCapabilityReviewImprovementDirectionUnknown,
        NonExecuting:   true,
        NonAuthorizing: true,
    }
    if !input.NonAuthorizing || !input.Feedback.NonAuthorizing {
        output.MissingStage = "authorization-boundary"
        return output
    }
    if !input.Feedback.NonExecuting {
        output.MissingStage = "execution-boundary"
        return output
    }
    if err := input.Feedback.Validate(); err != nil {
        output.MissingStage = "feedback-bridge"
        return output
    }
    if input.Target == "" {
        output.MissingStage = "improvement-target"
        return output
    }
    if input.CandidateSource == "" {
        output.MissingStage = "candidate-source"
        return output
    }
    direction := jevExternalApplyCapabilityReviewImprovementDirectionHold
    switch input.Feedback.Signal {
    case jevExternalApplyCapabilityReviewTrendSignalImprove:
        direction = jevExternalApplyCapabilityReviewImprovementDirectionGenerate
    case jevExternalApplyCapabilityReviewTrendSignalHold:
        direction = jevExternalApplyCapabilityReviewImprovementDirectionHold
    case jevExternalApplyCapabilityReviewTrendSignalRollback:
        direction = jevExternalApplyCapabilityReviewImprovementDirectionReject
    default:
        output.MissingStage = "feedback-signal"
        return output
    }
    output.Status = jevExternalApplyCapabilityReviewImprovementDirectionBound
    output.Direction = direction
    output.Target = input.Target
    output.CandidateSource = input.CandidateSource
    output.FeedbackDigest = input.Feedback.FeedbackDigest
    output.DirectionDigest = digestJEVExternalApplyCapabilityReviewImprovementDirection(output.Status, output.Direction, output.Target, output.CandidateSource, output.FeedbackDigest)
    if err := output.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewImprovementDirectionUnknown
        output.MissingStage = "improvement-direction-evidence"
        output.FeedbackDigest = ""
        output.DirectionDigest = ""
    }
    return output
}

func digestJEVExternalApplyCapabilityReviewImprovementDirection(status, direction, target, candidateSource, feedbackDigest string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s|%s", status, direction, target, candidateSource, feedbackDigest)))
    return hex.EncodeToString(sum[:])
}
