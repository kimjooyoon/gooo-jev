package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

const jevExternalApplyCapabilityReviewReplayOriginConsistent = "origin-consistent"
const jevExternalApplyCapabilityReviewReplayOriginMismatch = "origin-mismatch"
const jevExternalApplyCapabilityReviewReplayOriginUnknown = "UNKNOWN"

type JEVExternalApplyCapabilityReviewReplayOriginConsistencyInput struct {
    Observation          JEVExternalApplyCapabilityReviewReplayObservation
    DeclaredOriginDigest string
    ReverseOriginDigest string
    NonAuthorizing       bool
}

type JEVExternalApplyCapabilityReviewReplayOriginConsistency struct {
    Status               string
    MissingStage         string
    ObservationStatus    string
    DeclaredOriginDigest string
    ReverseOriginDigest  string
    ObservationDigest    string
    ConsistencyDigest    string
    NonExecuting         bool
    NonAuthorizing       bool
}

func (c JEVExternalApplyCapabilityReviewReplayOriginConsistency) Validate() error {
    if c.Status == "" || c.ObservationStatus == "" || c.DeclaredOriginDigest == "" || c.ReverseOriginDigest == "" || c.ObservationDigest == "" || c.ConsistencyDigest == "" {
        return fmt.Errorf("incomplete JEV external apply capability review replay origin consistency")
    }
    if c.Status != jevExternalApplyCapabilityReviewReplayOriginConsistent &&
        c.Status != jevExternalApplyCapabilityReviewReplayOriginMismatch {
        return fmt.Errorf("invalid JEV external apply capability review replay origin consistency status")
    }
    if c.ObservationStatus != jevExternalApplyCapabilityReviewReplayObserved {
        return fmt.Errorf("origin consistency requires replay-observed status")
    }
    if c.Status == jevExternalApplyCapabilityReviewReplayOriginConsistent && c.DeclaredOriginDigest != c.ReverseOriginDigest {
        return fmt.Errorf("origin-consistent status requires equal origin digests")
    }
    if c.Status == jevExternalApplyCapabilityReviewReplayOriginMismatch && c.DeclaredOriginDigest == c.ReverseOriginDigest {
        return fmt.Errorf("origin-mismatch status requires different origin digests")
    }
    if !c.NonExecuting {
        return fmt.Errorf("JEV external apply capability review replay origin consistency must be non-executing")
    }
    if !c.NonAuthorizing {
        return fmt.Errorf("JEV external apply capability review replay origin consistency must be non-authorizing")
    }
    expected := digestJEVExternalApplyCapabilityReviewReplayOriginConsistency(
        c.Status,
        c.ObservationStatus,
        c.DeclaredOriginDigest,
        c.ReverseOriginDigest,
        c.ObservationDigest,
    )
    if c.ConsistencyDigest != expected {
        return fmt.Errorf("JEV external apply capability review replay origin consistency digest mismatch")
    }
    return nil
}

func ReconcileJEVExternalApplyCapabilityReviewReplayOrigin(input JEVExternalApplyCapabilityReviewReplayOriginConsistencyInput) JEVExternalApplyCapabilityReviewReplayOriginConsistency {
    output := JEVExternalApplyCapabilityReviewReplayOriginConsistency{
        Status:         jevExternalApplyCapabilityReviewReplayOriginUnknown,
        NonExecuting:   true,
        NonAuthorizing: true,
    }
    if !input.NonAuthorizing || !input.Observation.NonAuthorizing {
        output.MissingStage = "authorization-boundary"
        return output
    }
    if !input.Observation.NonExecuting {
        output.MissingStage = "execution-boundary"
        return output
    }
    if err := input.Observation.Validate(); err != nil {
        output.MissingStage = "replay-observation"
        return output
    }
    if input.Observation.Status != jevExternalApplyCapabilityReviewReplayObserved {
        output.MissingStage = "replay-observed"
        return output
    }
    if input.DeclaredOriginDigest == "" {
        output.MissingStage = "declared-origin"
        return output
    }
    if input.ReverseOriginDigest == "" {
        output.MissingStage = "reverse-origin"
        return output
    }

    output.ObservationStatus = input.Observation.Status
    output.DeclaredOriginDigest = input.DeclaredOriginDigest
    output.ReverseOriginDigest = input.ReverseOriginDigest
    output.ObservationDigest = input.Observation.ObservationDigest
    if input.DeclaredOriginDigest == input.ReverseOriginDigest {
        output.Status = jevExternalApplyCapabilityReviewReplayOriginConsistent
    } else {
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMismatch
    }
    output.ConsistencyDigest = digestJEVExternalApplyCapabilityReviewReplayOriginConsistency(
        output.Status,
        output.ObservationStatus,
        output.DeclaredOriginDigest,
        output.ReverseOriginDigest,
        output.ObservationDigest,
    )
    if err := output.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewReplayOriginUnknown
        output.MissingStage = "replay-origin-consistency-evidence"
        output.ConsistencyDigest = ""
    }
    return output
}

func digestJEVExternalApplyCapabilityReviewReplayOriginConsistency(status, observationStatus, declaredOriginDigest, reverseOriginDigest, observationDigest string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s|%s", status, observationStatus, declaredOriginDigest, reverseOriginDigest, observationDigest)))
    return hex.EncodeToString(sum[:])
}
