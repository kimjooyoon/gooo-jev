package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

const (
    changePlanFeedbackDispositionUnknown = "UNKNOWN"
    changePlanFeedbackDispositionReady = "ready-for-external-apply-review"
    changePlanFeedbackDispositionAbort = "abort"
    changePlanFeedbackDispositionHold = "hold"
    changePlanFeedbackStatusConfirmed = "confirmed"
    changePlanFeedbackStatusRefuted = "refuted"
)

type DecisionConfidenceChangePlanFeedbackDispositionInput struct {
    Plan            DecisionConfidenceChangePlan
    Verification    DecisionConfidenceChangePlanVerification
    FeedbackBinding DecisionConfidenceChangePlanReplayFeedbackBinding
    NonAuthorizing  bool
}

type DecisionConfidenceChangePlanFeedbackDisposition struct {
    Status                 string
    FeedbackStatus         string
    MissingStage           string
    ChangePlanDigest       string
    VerificationDigest     string
    DispositionDigest      string
    FeedbackEvidenceDigest string
    EvidenceDigest         string
    NonExecuting           bool
    NonAuthorizing         bool
}

func (d DecisionConfidenceChangePlanFeedbackDisposition) Validate() error {
    if d.Status == "" || d.FeedbackStatus == "" || d.ChangePlanDigest == "" || d.VerificationDigest == "" || d.DispositionDigest == "" || d.FeedbackEvidenceDigest == "" || d.EvidenceDigest == "" {
        return fmt.Errorf("incomplete change plan feedback disposition")
    }
    if !d.NonExecuting {
        return fmt.Errorf("change plan feedback disposition must be non-executing")
    }
    if !d.NonAuthorizing {
        return fmt.Errorf("change plan feedback disposition must be non-authorizing")
    }
    expected := digestChangePlanFeedbackDisposition(d.ChangePlanDigest, d.VerificationDigest, d.DispositionDigest, d.FeedbackEvidenceDigest, d.FeedbackStatus)
    if d.EvidenceDigest != expected {
        return fmt.Errorf("change plan feedback disposition evidence digest mismatch")
    }
    return nil
}

func DispositionDecisionConfidenceChangePlanWithFeedback(input DecisionConfidenceChangePlanFeedbackDispositionInput) DecisionConfidenceChangePlanFeedbackDisposition {
    output := DecisionConfidenceChangePlanFeedbackDisposition{
        Status: changePlanFeedbackDispositionUnknown,
        FeedbackStatus: input.FeedbackBinding.FeedbackStatus,
        MissingStage: "authorization-boundary",
        NonExecuting: true,
        NonAuthorizing: true,
    }
    if !input.NonAuthorizing || !input.FeedbackBinding.NonAuthorizing {
        return output
    }
    if !input.Plan.NonExecuting || !input.Verification.NonExecuting || !input.FeedbackBinding.NonExecuting {
        output.MissingStage = "execution-boundary"
        return output
    }
    if err := input.Plan.Validate(); err != nil {
        output.MissingStage = "change-plan"
        return output
    }
    if err := input.Verification.Validate(); err != nil {
        output.MissingStage = "verification"
        return output
    }
    if err := input.FeedbackBinding.Validate(); err != nil {
        output.MissingStage = "feedback-binding"
        return output
    }
    output.ChangePlanDigest = input.Plan.ChangePlanDigest
    output.VerificationDigest = input.Verification.VerificationDigest
    output.FeedbackEvidenceDigest = input.FeedbackBinding.FeedbackEvidenceDigest
    if input.FeedbackBinding.ChangePlanDigest != input.Plan.ChangePlanDigest {
        output.MissingStage = "feedback-change-plan-binding"
        return output
    }
    disposition, err := DispositionDecisionConfidenceChangePlan(input.Plan, input.Verification)
    if err != nil {
        output.MissingStage = "change-plan-disposition"
        return output
    }
    if disposition.Status == "" || disposition.DispositionDigest == "" {
        output.MissingStage = "change-plan-disposition"
        return output
    }
    output.DispositionDigest = disposition.DispositionDigest
    switch input.FeedbackBinding.FeedbackStatus {
    case changePlanFeedbackStatusRefuted:
        output.Status = changePlanFeedbackDispositionAbort
    case changePlanFeedbackStatusConfirmed:
        output.Status = disposition.Status
    default:
        output.MissingStage = "feedback-disposition"
        return output
    }
    output.EvidenceDigest = digestChangePlanFeedbackDisposition(output.ChangePlanDigest, output.VerificationDigest, output.DispositionDigest, output.FeedbackEvidenceDigest, output.FeedbackStatus)
    if err := output.Validate(); err != nil {
        output.Status = changePlanFeedbackDispositionUnknown
        output.MissingStage = "feedback-disposition-evidence"
        output.EvidenceDigest = ""
    }
    return output
}

func digestChangePlanFeedbackDisposition(changePlanDigest, verificationDigest, dispositionDigest, feedbackEvidenceDigest, feedbackStatus string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s|%s", changePlanDigest, verificationDigest, dispositionDigest, feedbackEvidenceDigest, feedbackStatus)))
    return hex.EncodeToString(sum[:])
}
