package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

const jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationBridgeBound = "revision-candidate-application-observation-bound"
const jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationBridgeUnknown = "UNKNOWN"
const jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationObserved = "observed"
const jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMismatch = "mismatch"
const jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationHold = "hold"
const jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationRejected = "rejected"

type JEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationBridgeInput struct {
    CandidateApplication       JEVExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateBridge
    ObservationStatus          string
    ObservationSource          string
    ObservationEvidenceDigest  string
    ReverseObservationDigest   string
    NonAuthorizing             bool
}

type JEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationBridge struct {
    Status                       string
    MissingStage                 string
    CandidateApplicationStatus   string
    CandidateApplicationDigest   string
    ObservationStatus            string
    ObservationSource            string
    ObservationEvidenceDigest    string
    ReverseObservationDigest     string
    ObservationDigest            string
    NonExecuting                 bool
    NonAuthorizing               bool
}

func (b JEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationBridge) Validate() error {
    if b.Status == "" || b.CandidateApplicationStatus == "" ||
        b.CandidateApplicationDigest == "" || b.ObservationStatus == "" ||
        b.ObservationDigest == "" {
        return fmt.Errorf("incomplete JEV candidate application observation bridge")
    }
    if b.Status != jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationBridgeBound {
        return fmt.Errorf("invalid JEV candidate application observation bridge status")
    }
    switch b.CandidateApplicationStatus {
    case jevExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateReady:
        if b.ObservationStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationObserved &&
            b.ObservationStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMismatch {
            return fmt.Errorf("ready candidate application observation has invalid status")
        }
        if b.ObservationSource == "" || b.ObservationEvidenceDigest == "" ||
            b.ReverseObservationDigest == "" {
            return fmt.Errorf("candidate application observation evidence is incomplete")
        }
    case jevExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateHold:
        if b.ObservationStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationHold ||
            b.ObservationSource != "" || b.ObservationEvidenceDigest != "" ||
            b.ReverseObservationDigest != "" {
            return fmt.Errorf("held candidate application observation has inconsistent evidence")
        }
    case jevExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateRejected:
        if b.ObservationStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationRejected ||
            b.ObservationSource != "" || b.ObservationEvidenceDigest != "" ||
            b.ReverseObservationDigest != "" {
            return fmt.Errorf("rejected candidate application observation has inconsistent evidence")
        }
    default:
        return fmt.Errorf("invalid candidate application status")
    }
    if !b.NonExecuting {
        return fmt.Errorf("JEV candidate application observation bridge must be non-executing")
    }
    if !b.NonAuthorizing {
        return fmt.Errorf("JEV candidate application observation bridge must be non-authorizing")
    }
    expected := digestJEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationBridge(
        b.Status,
        b.CandidateApplicationStatus,
        b.CandidateApplicationDigest,
        b.ObservationStatus,
        b.ObservationSource,
        b.ObservationEvidenceDigest,
        b.ReverseObservationDigest,
    )
    if b.ObservationDigest != expected {
        return fmt.Errorf("JEV candidate application observation bridge digest mismatch")
    }
    return nil
}

func ObserveJEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationBridge(input JEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationBridgeInput) JEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationBridge {
    output := JEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationBridge{
        Status:                     jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationBridgeUnknown,
        CandidateApplicationStatus: input.CandidateApplication.ApplicationCandidateStatus,
        CandidateApplicationDigest: input.CandidateApplication.BridgeDigest,
        NonExecuting:               true,
        NonAuthorizing:             true,
    }
    if !input.NonAuthorizing || !input.CandidateApplication.NonAuthorizing {
        output.MissingStage = "authorization-boundary"
        return output
    }
    if !input.CandidateApplication.NonExecuting {
        output.MissingStage = "execution-boundary"
        return output
    }
    if err := input.CandidateApplication.Validate(); err != nil {
        output.MissingStage = "candidate-application-candidate-bridge"
        return output
    }
    switch input.CandidateApplication.ApplicationCandidateStatus {
    case jevExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateReady:
        if input.ObservationStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationObserved &&
            input.ObservationStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMismatch {
            output.MissingStage = "application-observation-status"
            return output
        }
        if input.ObservationSource == "" {
            output.MissingStage = "observation-source"
            return output
        }
        if input.ObservationEvidenceDigest == "" {
            output.MissingStage = "observation-evidence"
            return output
        }
        if input.ReverseObservationDigest == "" {
            output.MissingStage = "reverse-observation"
            return output
        }
        output.ObservationStatus = input.ObservationStatus
        output.ObservationSource = input.ObservationSource
        output.ObservationEvidenceDigest = input.ObservationEvidenceDigest
        output.ReverseObservationDigest = input.ReverseObservationDigest
    case jevExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateHold:
        output.ObservationStatus = jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationHold
    case jevExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateRejected:
        output.ObservationStatus = jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationRejected
    default:
        output.MissingStage = "candidate-application-status"
        return output
    }
    output.Status = jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationBridgeBound
    output.ObservationDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationBridge(
        output.Status,
        output.CandidateApplicationStatus,
        output.CandidateApplicationDigest,
        output.ObservationStatus,
        output.ObservationSource,
        output.ObservationEvidenceDigest,
        output.ReverseObservationDigest,
    )
    if err := output.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationBridgeUnknown
        output.MissingStage = "candidate-application-observation-evidence"
        output.ObservationDigest = ""
    }
    return output
}

func digestJEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationBridge(status, candidateApplicationStatus, candidateApplicationDigest, observationStatus, observationSource, observationEvidenceDigest, reverseObservationDigest string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s", status, candidateApplicationStatus, candidateApplicationDigest, observationStatus, observationSource, observationEvidenceDigest, reverseObservationDigest)))
    return hex.EncodeToString(sum[:])
}
