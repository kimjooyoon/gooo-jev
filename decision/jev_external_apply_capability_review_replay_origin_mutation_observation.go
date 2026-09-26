package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

const jevExternalApplyCapabilityReviewReplayOriginMutationObservationNotNeeded = "mutation-observation-not-needed"
const jevExternalApplyCapabilityReviewReplayOriginMutationObservationRecorded = "mutation-observation-recorded"
const jevExternalApplyCapabilityReviewReplayOriginMutationObservationHold = "mutation-observation-hold"
const jevExternalApplyCapabilityReviewReplayOriginMutationObservationRejected = "mutation-observation-rejected"
const jevExternalApplyCapabilityReviewReplayOriginMutationObservationUnknown = "UNKNOWN"

type JEVExternalApplyCapabilityReviewReplayOriginMutationObservationInput struct {
    CapabilityBoundary      JEVExternalApplyCapabilityReviewReplayOriginMutationCapabilityBoundary
    ObservedMutationDigest  string
    ReverseMutationDigest   string
    NonAuthorizing          bool
}

type JEVExternalApplyCapabilityReviewReplayOriginMutationObservation struct {
    Status                  string
    MissingStage            string
    CapabilityStatus        string
    CapabilityDigest        string
    ObservedMutationDigest  string
    ReverseMutationDigest   string
    MutationObservationDigest string
    NonExecuting            bool
    NonAuthorizing          bool
}

func (o JEVExternalApplyCapabilityReviewReplayOriginMutationObservation) Validate() error {
    if o.Status == "" || o.CapabilityStatus == "" || o.CapabilityDigest == "" || o.MutationObservationDigest == "" {
        return fmt.Errorf("incomplete JEV external apply capability review replay origin mutation observation")
    }
    if o.Status != jevExternalApplyCapabilityReviewReplayOriginMutationObservationNotNeeded &&
        o.Status != jevExternalApplyCapabilityReviewReplayOriginMutationObservationRecorded &&
        o.Status != jevExternalApplyCapabilityReviewReplayOriginMutationObservationHold &&
        o.Status != jevExternalApplyCapabilityReviewReplayOriginMutationObservationRejected {
        return fmt.Errorf("invalid JEV external apply capability review replay origin mutation observation status")
    }
    if o.CapabilityStatus != jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityNotNeeded &&
        o.CapabilityStatus != jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityRecorded &&
        o.CapabilityStatus != jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityHold &&
        o.CapabilityStatus != jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityRejected {
        return fmt.Errorf("invalid JEV external apply capability review replay origin mutation capability status")
    }
    if o.CapabilityStatus == jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityNotNeeded &&
        o.Status != jevExternalApplyCapabilityReviewReplayOriginMutationObservationNotNeeded {
        return fmt.Errorf("capability-not-needed requires mutation-observation-not-needed status")
    }
    if o.CapabilityStatus == jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityRecorded {
        if o.Status != jevExternalApplyCapabilityReviewReplayOriginMutationObservationRecorded {
            return fmt.Errorf("recorded capability requires mutation-observation-recorded status")
        }
        if o.ObservedMutationDigest == "" || o.ReverseMutationDigest == "" {
            return fmt.Errorf("recorded mutation observation requires forward and reverse evidence")
        }
    }
    if o.CapabilityStatus == jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityHold &&
        o.Status != jevExternalApplyCapabilityReviewReplayOriginMutationObservationHold {
        return fmt.Errorf("capability-hold requires mutation-observation-hold status")
    }
    if o.CapabilityStatus == jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityRejected &&
        o.Status != jevExternalApplyCapabilityReviewReplayOriginMutationObservationRejected {
        return fmt.Errorf("capability-rejected requires mutation-observation-rejected status")
    }
    if !o.NonExecuting {
        return fmt.Errorf("JEV external apply capability review replay origin mutation observation must be non-executing")
    }
    if !o.NonAuthorizing {
        return fmt.Errorf("JEV external apply capability review replay origin mutation observation must be non-authorizing")
    }
    expected := digestJEVExternalApplyCapabilityReviewReplayOriginMutationObservation(
        o.Status,
        o.CapabilityStatus,
        o.CapabilityDigest,
        o.ObservedMutationDigest,
        o.ReverseMutationDigest,
    )
    if o.MutationObservationDigest != expected {
        return fmt.Errorf("JEV external apply capability review replay origin mutation observation digest mismatch")
    }
    return nil
}

func ObserveJEVExternalApplyCapabilityReviewReplayOriginMutation(input JEVExternalApplyCapabilityReviewReplayOriginMutationObservationInput) JEVExternalApplyCapabilityReviewReplayOriginMutationObservation {
    output := JEVExternalApplyCapabilityReviewReplayOriginMutationObservation{
        Status:         jevExternalApplyCapabilityReviewReplayOriginMutationObservationUnknown,
        NonExecuting:   true,
        NonAuthorizing: true,
    }
    if !input.NonAuthorizing || !input.CapabilityBoundary.NonAuthorizing {
        output.MissingStage = "authorization-boundary"
        return output
    }
    if !input.CapabilityBoundary.NonExecuting {
        output.MissingStage = "execution-boundary"
        return output
    }
    if err := input.CapabilityBoundary.Validate(); err != nil {
        output.MissingStage = "mutation-capability-boundary"
        return output
    }
    output.CapabilityStatus = input.CapabilityBoundary.Status
    output.CapabilityDigest = input.CapabilityBoundary.CapabilityDigest
    switch input.CapabilityBoundary.Status {
    case jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityNotNeeded:
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationObservationNotNeeded
    case jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityRecorded:
        if input.ObservedMutationDigest == "" {
            output.MissingStage = "observed-mutation"
            return output
        }
        if input.ReverseMutationDigest == "" {
            output.MissingStage = "reverse-mutation-observation"
            return output
        }
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationObservationRecorded
        output.ObservedMutationDigest = input.ObservedMutationDigest
        output.ReverseMutationDigest = input.ReverseMutationDigest
    case jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityHold:
        if input.ObservedMutationDigest != "" || input.ReverseMutationDigest != "" {
            output.MissingStage = "mutation-observation-lifecycle"
            return output
        }
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationObservationHold
    case jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityRejected:
        if input.ObservedMutationDigest != "" || input.ReverseMutationDigest != "" {
            output.MissingStage = "mutation-observation-lifecycle"
            return output
        }
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationObservationRejected
    default:
        output.MissingStage = "mutation-capability-status"
        return output
    }
    output.MutationObservationDigest = digestJEVExternalApplyCapabilityReviewReplayOriginMutationObservation(
        output.Status,
        output.CapabilityStatus,
        output.CapabilityDigest,
        output.ObservedMutationDigest,
        output.ReverseMutationDigest,
    )
    if err := output.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationObservationUnknown
        output.MissingStage = "mutation-observation-evidence"
        output.MutationObservationDigest = ""
    }
    return output
}

func digestJEVExternalApplyCapabilityReviewReplayOriginMutationObservation(status, capabilityStatus, capabilityDigest, observedMutationDigest, reverseMutationDigest string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s|%s", status, capabilityStatus, capabilityDigest, observedMutationDigest, reverseMutationDigest)))
    return hex.EncodeToString(sum[:])
}
