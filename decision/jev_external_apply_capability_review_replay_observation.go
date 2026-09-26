package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

const jevExternalApplyCapabilityReviewReplayObserved = "replay-observed"
const jevExternalApplyCapabilityReviewReplayObservationHold = "replay-observation-hold"
const jevExternalApplyCapabilityReviewReplayObservationRejected = "replay-observation-rejected"
const jevExternalApplyCapabilityReviewReplayObservationUnknown = "UNKNOWN"

type JEVExternalApplyCapabilityReviewReplayObservationInput struct {
    Preparation             JEVExternalApplyCapabilityReviewReplayPreparation
    ObservedArtifactDigest  string
    ReverseObservationDigest string
    NonAuthorizing          bool
}

type JEVExternalApplyCapabilityReviewReplayObservation struct {
    Status                  string
    MissingStage            string
    PreparationDecision     string
    ReplayScope             string
    ReplayEnvironmentDigest string
    PreparationDigest       string
    CandidateDigest         string
    ObservedArtifactDigest  string
    ReverseObservationDigest string
    ObservationDigest       string
    NonExecuting            bool
    NonAuthorizing          bool
}

func (o JEVExternalApplyCapabilityReviewReplayObservation) Validate() error {
    if o.Status == "" || o.PreparationDecision == "" || o.ReplayScope == "" || o.PreparationDigest == "" || o.ObservationDigest == "" {
        return fmt.Errorf("incomplete JEV external apply capability review replay observation")
    }
    if o.Status != jevExternalApplyCapabilityReviewReplayObserved &&
        o.Status != jevExternalApplyCapabilityReviewReplayObservationHold &&
        o.Status != jevExternalApplyCapabilityReviewReplayObservationRejected {
        return fmt.Errorf("invalid JEV external apply capability review replay observation status")
    }
    if o.PreparationDecision != jevExternalApplyCapabilityReviewReplayReady &&
        o.PreparationDecision != jevExternalApplyCapabilityReviewReplayHold &&
        o.PreparationDecision != jevExternalApplyCapabilityReviewReplayRejected {
        return fmt.Errorf("invalid JEV external apply capability review replay preparation decision")
    }
    if o.PreparationDecision == jevExternalApplyCapabilityReviewReplayReady {
        if o.Status != jevExternalApplyCapabilityReviewReplayObserved {
            return fmt.Errorf("replay-ready preparation requires replay-observed status")
        }
        if o.ObservedArtifactDigest == "" || o.ReverseObservationDigest == "" {
            return fmt.Errorf("replay-observed result requires forward and reverse observation evidence")
        }
    }
    if o.PreparationDecision == jevExternalApplyCapabilityReviewReplayHold &&
        o.Status != jevExternalApplyCapabilityReviewReplayObservationHold {
        return fmt.Errorf("replay-hold preparation requires replay-observation-hold status")
    }
    if o.PreparationDecision == jevExternalApplyCapabilityReviewReplayRejected &&
        o.Status != jevExternalApplyCapabilityReviewReplayObservationRejected {
        return fmt.Errorf("replay-rejected preparation requires replay-observation-rejected status")
    }
    if !o.NonExecuting {
        return fmt.Errorf("JEV external apply capability review replay observation must be non-executing")
    }
    if !o.NonAuthorizing {
        return fmt.Errorf("JEV external apply capability review replay observation must be non-authorizing")
    }
    expected := digestJEVExternalApplyCapabilityReviewReplayObservation(
        o.Status,
        o.PreparationDecision,
        o.ReplayScope,
        o.ReplayEnvironmentDigest,
        o.PreparationDigest,
        o.CandidateDigest,
        o.ObservedArtifactDigest,
        o.ReverseObservationDigest,
    )
    if o.ObservationDigest != expected {
        return fmt.Errorf("JEV external apply capability review replay observation digest mismatch")
    }
    return nil
}

func ObserveJEVExternalApplyCapabilityReviewReplay(input JEVExternalApplyCapabilityReviewReplayObservationInput) JEVExternalApplyCapabilityReviewReplayObservation {
    output := JEVExternalApplyCapabilityReviewReplayObservation{
        Status:         jevExternalApplyCapabilityReviewReplayObservationUnknown,
        NonExecuting:   true,
        NonAuthorizing: true,
    }
    if !input.NonAuthorizing || !input.Preparation.NonAuthorizing {
        output.MissingStage = "authorization-boundary"
        return output
    }
    if !input.Preparation.NonExecuting {
        output.MissingStage = "execution-boundary"
        return output
    }
    if err := input.Preparation.Validate(); err != nil {
        output.MissingStage = "replay-preparation"
        return output
    }
    if input.Preparation.ReplayScope == "" {
        output.MissingStage = "replay-scope"
        return output
    }

    output.PreparationDecision = input.Preparation.Decision
    output.ReplayScope = input.Preparation.ReplayScope
    output.ReplayEnvironmentDigest = input.Preparation.ReplayEnvironmentDigest
    output.PreparationDigest = input.Preparation.PreparationDigest
    output.CandidateDigest = input.Preparation.CandidateDigest
    output.ObservedArtifactDigest = input.ObservedArtifactDigest
    output.ReverseObservationDigest = input.ReverseObservationDigest
    switch input.Preparation.Decision {
    case jevExternalApplyCapabilityReviewReplayReady:
        if input.ObservedArtifactDigest == "" {
            output.MissingStage = "observed-artifact"
            return output
        }
        if input.ReverseObservationDigest == "" {
            output.MissingStage = "reverse-observation"
            return output
        }
        output.Status = jevExternalApplyCapabilityReviewReplayObserved
    case jevExternalApplyCapabilityReviewReplayHold:
        if input.ObservedArtifactDigest != "" || input.ReverseObservationDigest != "" {
            output.MissingStage = "observation-lifecycle"
            return output
        }
        output.Status = jevExternalApplyCapabilityReviewReplayObservationHold
    case jevExternalApplyCapabilityReviewReplayRejected:
        if input.ObservedArtifactDigest != "" || input.ReverseObservationDigest != "" {
            output.MissingStage = "observation-lifecycle"
            return output
        }
        output.Status = jevExternalApplyCapabilityReviewReplayObservationRejected
    default:
        output.MissingStage = "replay-decision"
        return output
    }
    output.ObservationDigest = digestJEVExternalApplyCapabilityReviewReplayObservation(
        output.Status,
        output.PreparationDecision,
        output.ReplayScope,
        output.ReplayEnvironmentDigest,
        output.PreparationDigest,
        output.CandidateDigest,
        output.ObservedArtifactDigest,
        output.ReverseObservationDigest,
    )
    if err := output.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewReplayObservationUnknown
        output.MissingStage = "replay-observation-evidence"
        output.ObservationDigest = ""
    }
    return output
}

func digestJEVExternalApplyCapabilityReviewReplayObservation(status, preparationDecision, replayScope, replayEnvironmentDigest, preparationDigest, candidateDigest, observedArtifactDigest, reverseObservationDigest string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s", status, preparationDecision, replayScope, replayEnvironmentDigest, preparationDigest, candidateDigest, observedArtifactDigest, reverseObservationDigest)))
    return hex.EncodeToString(sum[:])
}
