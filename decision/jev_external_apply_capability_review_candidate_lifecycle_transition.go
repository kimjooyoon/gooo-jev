package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

const jevExternalApplyCapabilityReviewCandidateLifecycleTransitionRecorded = "candidate-lifecycle-transition-recorded"
const jevExternalApplyCapabilityReviewCandidateLifecycleReady = "candidate-ready"
const jevExternalApplyCapabilityReviewCandidateLifecycleHold = "candidate-on-hold"
const jevExternalApplyCapabilityReviewCandidateLifecycleRejected = "candidate-rejected"
const jevExternalApplyCapabilityReviewCandidateLifecycleUnknown = "UNKNOWN"

type JEVExternalApplyCapabilityReviewCandidateLifecycleTransitionInput struct {
    Gate           JEVExternalApplyCapabilityReviewRevisionCandidateGate
    FromState      string
    NonAuthorizing bool
}

type JEVExternalApplyCapabilityReviewCandidateLifecycleTransition struct {
    Status           string
    MissingStage     string
    FromState        string
    ToState          string
    Decision         string
    GateDigest       string
    CandidateDigest  string
    TransitionDigest string
    NonExecuting     bool
    NonAuthorizing   bool
}

func (t JEVExternalApplyCapabilityReviewCandidateLifecycleTransition) Validate() error {
    if t.Status == "" || t.FromState == "" || t.ToState == "" || t.Decision == "" || t.GateDigest == "" || t.TransitionDigest == "" {
        return fmt.Errorf("incomplete JEV external apply capability review candidate lifecycle transition")
    }
    if t.Status != jevExternalApplyCapabilityReviewCandidateLifecycleTransitionRecorded {
        return fmt.Errorf("invalid JEV external apply capability review candidate lifecycle transition status")
    }
    if t.Decision != jevExternalApplyCapabilityReviewRevisionCandidateReady && t.Decision != jevExternalApplyCapabilityReviewRevisionCandidateHold && t.Decision != jevExternalApplyCapabilityReviewRevisionCandidateRejected {
        return fmt.Errorf("invalid JEV external apply capability review candidate lifecycle transition decision")
    }
    if t.ToState != jevExternalApplyCapabilityReviewCandidateLifecycleReady && t.ToState != jevExternalApplyCapabilityReviewCandidateLifecycleHold && t.ToState != jevExternalApplyCapabilityReviewCandidateLifecycleRejected {
        return fmt.Errorf("invalid JEV external apply capability review candidate lifecycle transition target")
    }
    if !t.NonExecuting {
        return fmt.Errorf("JEV external apply capability review candidate lifecycle transition must be non-executing")
    }
    if !t.NonAuthorizing {
        return fmt.Errorf("JEV external apply capability review candidate lifecycle transition must be non-authorizing")
    }
    expected := digestJEVExternalApplyCapabilityReviewCandidateLifecycleTransition(t.Status, t.FromState, t.ToState, t.Decision, t.GateDigest, t.CandidateDigest)
    if t.TransitionDigest != expected {
        return fmt.Errorf("JEV external apply capability review candidate lifecycle transition digest mismatch")
    }
    return nil
}

func TransitionJEVExternalApplyCapabilityReviewCandidateLifecycle(input JEVExternalApplyCapabilityReviewCandidateLifecycleTransitionInput) JEVExternalApplyCapabilityReviewCandidateLifecycleTransition {
    output := JEVExternalApplyCapabilityReviewCandidateLifecycleTransition{
        Status:         jevExternalApplyCapabilityReviewCandidateLifecycleUnknown,
        NonExecuting:   true,
        NonAuthorizing: true,
    }
    if !input.NonAuthorizing || !input.Gate.NonAuthorizing {
        output.MissingStage = "authorization-boundary"
        return output
    }
    if !input.Gate.NonExecuting {
        output.MissingStage = "execution-boundary"
        return output
    }
    if err := input.Gate.Validate(); err != nil {
        output.MissingStage = "revision-candidate-gate"
        return output
    }
    if input.FromState == "" {
        output.MissingStage = "lifecycle-source-state"
        return output
    }
    toState := ""
    switch input.Gate.Decision {
    case jevExternalApplyCapabilityReviewRevisionCandidateReady:
        toState = jevExternalApplyCapabilityReviewCandidateLifecycleReady
    case jevExternalApplyCapabilityReviewRevisionCandidateHold:
        toState = jevExternalApplyCapabilityReviewCandidateLifecycleHold
    case jevExternalApplyCapabilityReviewRevisionCandidateRejected:
        toState = jevExternalApplyCapabilityReviewCandidateLifecycleRejected
    default:
        output.MissingStage = "lifecycle-decision"
        return output
    }
    output.Status = jevExternalApplyCapabilityReviewCandidateLifecycleTransitionRecorded
    output.FromState = input.FromState
    output.ToState = toState
    output.Decision = input.Gate.Decision
    output.GateDigest = input.Gate.GateDigest
    output.CandidateDigest = input.Gate.CandidateDigest
    output.TransitionDigest = digestJEVExternalApplyCapabilityReviewCandidateLifecycleTransition(output.Status, output.FromState, output.ToState, output.Decision, output.GateDigest, output.CandidateDigest)
    if err := output.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewCandidateLifecycleUnknown
        output.MissingStage = "candidate-lifecycle-transition-evidence"
        output.GateDigest = ""
        output.TransitionDigest = ""
    }
    return output
}

func digestJEVExternalApplyCapabilityReviewCandidateLifecycleTransition(status, fromState, toState, decision, gateDigest, candidateDigest string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s|%s|%s", status, fromState, toState, decision, gateDigest, candidateDigest)))
    return hex.EncodeToString(sum[:])
}
