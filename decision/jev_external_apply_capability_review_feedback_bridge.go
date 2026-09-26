package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

const jevExternalApplyCapabilityReviewFeedbackBound = "capability-review-feedback-bound"
const jevExternalApplyCapabilityReviewFeedbackUnknown = "UNKNOWN"

type JEVExternalApplyCapabilityReviewFeedbackBridgeInput struct {
    Signal                 JEVExternalApplyCapabilityReviewTrendSignal
    FeedbackEvidenceDigest string
    NonAuthorizing         bool
}

type JEVExternalApplyCapabilityReviewFeedbackBridge struct {
    Status                 string
    MissingStage           string
    Signal                 string
    SignalDigest            string
    TrendDigest             string
    FeedbackEvidenceDigest  string
    FeedbackDigest          string
    NonExecuting            bool
    NonAuthorizing          bool
}

func (b JEVExternalApplyCapabilityReviewFeedbackBridge) Validate() error {
    if b.Status == "" || b.Signal == "" || b.SignalDigest == "" || b.TrendDigest == "" || b.FeedbackEvidenceDigest == "" || b.FeedbackDigest == "" {
        return fmt.Errorf("incomplete JEV external apply capability review feedback bridge")
    }
    if b.Status != jevExternalApplyCapabilityReviewFeedbackBound {
        return fmt.Errorf("invalid JEV external apply capability review feedback bridge status")
    }
    if !b.NonExecuting {
        return fmt.Errorf("JEV external apply capability review feedback bridge must be non-executing")
    }
    if !b.NonAuthorizing {
        return fmt.Errorf("JEV external apply capability review feedback bridge must be non-authorizing")
    }
    expected := digestJEVExternalApplyCapabilityReviewFeedback(b.Signal, b.SignalDigest, b.TrendDigest, b.FeedbackEvidenceDigest)
    if b.FeedbackDigest != expected {
        return fmt.Errorf("JEV external apply capability review feedback digest mismatch")
    }
    return nil
}

func BindJEVExternalApplyCapabilityReviewFeedback(input JEVExternalApplyCapabilityReviewFeedbackBridgeInput) JEVExternalApplyCapabilityReviewFeedbackBridge {
    output := JEVExternalApplyCapabilityReviewFeedbackBridge{
        Status:         jevExternalApplyCapabilityReviewFeedbackUnknown,
        NonExecuting:   true,
        NonAuthorizing: true,
    }
    if !input.NonAuthorizing || !input.Signal.NonAuthorizing {
        output.MissingStage = "authorization-boundary"
        return output
    }
    if !input.Signal.NonExecuting {
        output.MissingStage = "execution-boundary"
        return output
    }
    if err := input.Signal.Validate(); err != nil {
        output.MissingStage = "trend-signal"
        return output
    }
    if input.FeedbackEvidenceDigest == "" {
        output.MissingStage = "feedback-evidence"
        return output
    }
    output.Status = jevExternalApplyCapabilityReviewFeedbackBound
    output.Signal = input.Signal.Signal
    output.SignalDigest = input.Signal.SignalDigest
    output.TrendDigest = input.Signal.TrendDigest
    output.FeedbackEvidenceDigest = input.FeedbackEvidenceDigest
    output.FeedbackDigest = digestJEVExternalApplyCapabilityReviewFeedback(output.Signal, output.SignalDigest, output.TrendDigest, output.FeedbackEvidenceDigest)
    if err := output.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewFeedbackUnknown
        output.MissingStage = "feedback-binding-evidence"
        output.SignalDigest = ""
        output.TrendDigest = ""
        output.FeedbackDigest = ""
    }
    return output
}

func digestJEVExternalApplyCapabilityReviewFeedback(signal, signalDigest, trendDigest, feedbackEvidenceDigest string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s", signal, signalDigest, trendDigest, feedbackEvidenceDigest)))
    return hex.EncodeToString(sum[:])
}
