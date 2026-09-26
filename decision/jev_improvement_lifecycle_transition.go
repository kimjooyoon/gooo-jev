package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

const (
    jevLifecycleTransitionAllowed = "transition-allowed"
    jevLifecycleTransitionUnknown = "UNKNOWN"
)

type JEVImprovementLifecycleTransitionInput struct {
    FromStatus          string
    ToStatus            string
    InputEvidenceDigest string
    NonAuthorizing      bool
}

type JEVImprovementLifecycleTransition struct {
    Status              string
    MissingStage        string
    FromStatus          string
    ToStatus            string
    InputEvidenceDigest string
    EvidenceDigest      string
    NonExecuting        bool
    NonAuthorizing      bool
}

func (t JEVImprovementLifecycleTransition) Validate() error {
    if t.Status == "" || t.FromStatus == "" || t.ToStatus == "" || t.InputEvidenceDigest == "" || t.EvidenceDigest == "" {
        return fmt.Errorf("incomplete JEV improvement lifecycle transition")
    }
    if t.Status != jevLifecycleTransitionAllowed {
        return fmt.Errorf("invalid JEV improvement lifecycle transition status")
    }
    if !t.NonExecuting {
        return fmt.Errorf("JEV improvement lifecycle transition must be non-executing")
    }
    if !t.NonAuthorizing {
        return fmt.Errorf("JEV improvement lifecycle transition must be non-authorizing")
    }
    expected := digestJEVImprovementLifecycleTransition(t.Status, t.FromStatus, t.ToStatus, t.InputEvidenceDigest)
    if t.EvidenceDigest != expected {
        return fmt.Errorf("JEV improvement lifecycle transition evidence mismatch")
    }
    return nil
}

func ValidateJEVImprovementLifecycleTransition(input JEVImprovementLifecycleTransitionInput) JEVImprovementLifecycleTransition {
    output := JEVImprovementLifecycleTransition{
        Status: jevLifecycleTransitionUnknown,
        NonExecuting: true,
        NonAuthorizing: true,
    }
    if !input.NonAuthorizing {
        output.MissingStage = "authorization-boundary"
        return output
    }
    if input.FromStatus == "" {
        output.MissingStage = "from-status"
        return output
    }
    if input.ToStatus == "" {
        output.MissingStage = "to-status"
        return output
    }
    if input.InputEvidenceDigest == "" {
        output.MissingStage = "transition-evidence"
        return output
    }
    if !allowedJEVImprovementLifecycleTransition(input.FromStatus, input.ToStatus) {
        output.MissingStage = "lifecycle-transition"
        return output
    }
    output.Status = jevLifecycleTransitionAllowed
    output.FromStatus = input.FromStatus
    output.ToStatus = input.ToStatus
    output.InputEvidenceDigest = input.InputEvidenceDigest
    output.EvidenceDigest = digestJEVImprovementLifecycleTransition(output.Status, output.FromStatus, output.ToStatus, output.InputEvidenceDigest)
    if err := output.Validate(); err != nil {
        output.Status = jevLifecycleTransitionUnknown
        output.MissingStage = "lifecycle-transition-evidence"
        output.EvidenceDigest = ""
    }
    return output
}

func allowedJEVImprovementLifecycleTransition(fromStatus, toStatus string) bool {
    switch fromStatus {
    case jevImprovementCandidateReady:
        return toStatus == jevImprovementCandidateReviewConfirmed || toStatus == jevImprovementCandidateReviewAbort || toStatus == jevImprovementCandidateReviewHold
    case jevImprovementCandidateReviewConfirmed:
        return toStatus == jevExternalApplyRequestReady
    case jevExternalApplyRequestReady:
        return toStatus == jevExternalApplyObservedApplied || toStatus == jevExternalApplyObservedRejected || toStatus == jevExternalApplyObservedUnknown
    case jevExternalApplyObservedApplied:
        return toStatus == jevPostApplyCycleReplayRequired
    case jevExternalApplyObservedRejected:
        return toStatus == jevPostApplyCycleRevisionRequired
    case jevExternalApplyObservedUnknown:
        return toStatus == jevPostApplyCycleEvidenceRequired
    case jevImprovementDirectiveRevision:
        return toStatus == jevImprovementRevisionCandidateReady
    case jevImprovementRevisionCandidateReady:
        return toStatus == jevRevisionCandidateReviewBridgeReady
    case jevRevisionCandidateReviewBridgeReady:
        return toStatus == jevImprovementCandidateReady
    case jevImprovementReplayReady:
        return toStatus == jevImprovementReplayReproduced || toStatus == jevImprovementReplayCounterexample || toStatus == jevImprovementReplayObservationUnknown
    case jevImprovementReplayReproduced:
        return toStatus == jevReplayFeedbackConfirmed
    case jevImprovementReplayCounterexample:
        return toStatus == jevReplayFeedbackRefuted
    case jevReplayFeedbackConfirmed, jevReplayFeedbackRefuted, jevReplayFeedbackUnknown:
        return toStatus == jevImprovementFeedbackStableForReview || toStatus == jevImprovementFeedbackNeedsRevision || toStatus == jevImprovementFeedbackHold
    default:
        return false
    }
}

func digestJEVImprovementLifecycleTransition(status, fromStatus, toStatus, inputEvidenceDigest string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s", status, fromStatus, toStatus, inputEvidenceDigest)))
    return hex.EncodeToString(sum[:])
}
