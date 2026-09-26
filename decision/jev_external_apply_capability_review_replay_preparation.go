package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

const jevExternalApplyCapabilityReviewReplayPreparationRecorded = "replay-preparation-recorded"
const jevExternalApplyCapabilityReviewReplayReady = "replay-ready"
const jevExternalApplyCapabilityReviewReplayHold = "replay-hold"
const jevExternalApplyCapabilityReviewReplayRejected = "replay-rejected"
const jevExternalApplyCapabilityReviewReplayUnknown = "UNKNOWN"

type JEVExternalApplyCapabilityReviewReplayPreparationInput struct {
    Transition              JEVExternalApplyCapabilityReviewCandidateLifecycleTransition
    ReplayScope             string
    ReplayEnvironmentDigest string
    NonAuthorizing          bool
}

type JEVExternalApplyCapabilityReviewReplayPreparation struct {
    Status                  string
    MissingStage            string
    LifecycleState          string
    Decision                string
    ReplayScope             string
    ReplayEnvironmentDigest string
    TransitionDigest        string
    CandidateDigest         string
    PreparationDigest       string
    NonExecuting            bool
    NonAuthorizing          bool
}

func (p JEVExternalApplyCapabilityReviewReplayPreparation) Validate() error {
    if p.Status == "" || p.LifecycleState == "" || p.Decision == "" || p.ReplayScope == "" || p.TransitionDigest == "" || p.PreparationDigest == "" {
        return fmt.Errorf("incomplete JEV external apply capability review replay preparation")
    }
    if p.Status != jevExternalApplyCapabilityReviewReplayPreparationRecorded {
        return fmt.Errorf("invalid JEV external apply capability review replay preparation status")
    }
    if p.LifecycleState != jevExternalApplyCapabilityReviewCandidateLifecycleReady &&
        p.LifecycleState != jevExternalApplyCapabilityReviewCandidateLifecycleHold &&
        p.LifecycleState != jevExternalApplyCapabilityReviewCandidateLifecycleRejected {
        return fmt.Errorf("invalid JEV external apply capability review replay lifecycle state")
    }
    if p.Decision != jevExternalApplyCapabilityReviewReplayReady &&
        p.Decision != jevExternalApplyCapabilityReviewReplayHold &&
        p.Decision != jevExternalApplyCapabilityReviewReplayRejected {
        return fmt.Errorf("invalid JEV external apply capability review replay decision")
    }
    if p.LifecycleState == jevExternalApplyCapabilityReviewCandidateLifecycleReady &&
        p.Decision != jevExternalApplyCapabilityReviewReplayReady {
        return fmt.Errorf("ready lifecycle state requires replay-ready decision")
    }
    if p.LifecycleState == jevExternalApplyCapabilityReviewCandidateLifecycleHold &&
        p.Decision != jevExternalApplyCapabilityReviewReplayHold {
        return fmt.Errorf("hold lifecycle state requires replay-hold decision")
    }
    if p.LifecycleState == jevExternalApplyCapabilityReviewCandidateLifecycleRejected &&
        p.Decision != jevExternalApplyCapabilityReviewReplayRejected {
        return fmt.Errorf("rejected lifecycle state requires replay-rejected decision")
    }
    if p.Decision == jevExternalApplyCapabilityReviewReplayReady && p.ReplayEnvironmentDigest == "" {
        return fmt.Errorf("replay-ready preparation requires replay environment digest")
    }
    if !p.NonExecuting {
        return fmt.Errorf("JEV external apply capability review replay preparation must be non-executing")
    }
    if !p.NonAuthorizing {
        return fmt.Errorf("JEV external apply capability review replay preparation must be non-authorizing")
    }
    expected := digestJEVExternalApplyCapabilityReviewReplayPreparation(
        p.Status,
        p.LifecycleState,
        p.Decision,
        p.ReplayScope,
        p.ReplayEnvironmentDigest,
        p.TransitionDigest,
        p.CandidateDigest,
    )
    if p.PreparationDigest != expected {
        return fmt.Errorf("JEV external apply capability review replay preparation digest mismatch")
    }
    return nil
}

func PrepareJEVExternalApplyCapabilityReviewReplay(input JEVExternalApplyCapabilityReviewReplayPreparationInput) JEVExternalApplyCapabilityReviewReplayPreparation {
    output := JEVExternalApplyCapabilityReviewReplayPreparation{
        Status:         jevExternalApplyCapabilityReviewReplayUnknown,
        NonExecuting:   true,
        NonAuthorizing: true,
    }
    if !input.NonAuthorizing || !input.Transition.NonAuthorizing {
        output.MissingStage = "authorization-boundary"
        return output
    }
    if !input.Transition.NonExecuting {
        output.MissingStage = "execution-boundary"
        return output
    }
    if err := input.Transition.Validate(); err != nil {
        output.MissingStage = "candidate-lifecycle-transition"
        return output
    }
    if input.ReplayScope == "" {
        output.MissingStage = "replay-scope"
        return output
    }
    if input.Transition.ToState == jevExternalApplyCapabilityReviewCandidateLifecycleReady &&
        input.ReplayEnvironmentDigest == "" {
        output.MissingStage = "replay-environment"
        return output
    }

    output.Status = jevExternalApplyCapabilityReviewReplayPreparationRecorded
    output.LifecycleState = input.Transition.ToState
    output.ReplayScope = input.ReplayScope
    output.ReplayEnvironmentDigest = input.ReplayEnvironmentDigest
    output.TransitionDigest = input.Transition.TransitionDigest
    output.CandidateDigest = input.Transition.CandidateDigest
    switch input.Transition.ToState {
    case jevExternalApplyCapabilityReviewCandidateLifecycleReady:
        output.Decision = jevExternalApplyCapabilityReviewReplayReady
    case jevExternalApplyCapabilityReviewCandidateLifecycleHold:
        output.Decision = jevExternalApplyCapabilityReviewReplayHold
    case jevExternalApplyCapabilityReviewCandidateLifecycleRejected:
        output.Decision = jevExternalApplyCapabilityReviewReplayRejected
    default:
        output.Status = jevExternalApplyCapabilityReviewReplayUnknown
        output.MissingStage = "replay-lifecycle-state"
        return output
    }
    output.PreparationDigest = digestJEVExternalApplyCapabilityReviewReplayPreparation(
        output.Status,
        output.LifecycleState,
        output.Decision,
        output.ReplayScope,
        output.ReplayEnvironmentDigest,
        output.TransitionDigest,
        output.CandidateDigest,
    )
    if err := output.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewReplayUnknown
        output.MissingStage = "replay-preparation-evidence"
        output.TransitionDigest = ""
        output.PreparationDigest = ""
    }
    return output
}

func digestJEVExternalApplyCapabilityReviewReplayPreparation(status, lifecycleState, decision, replayScope, replayEnvironmentDigest, transitionDigest, candidateDigest string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s", status, lifecycleState, decision, replayScope, replayEnvironmentDigest, transitionDigest, candidateDigest)))
    return hex.EncodeToString(sum[:])
}
