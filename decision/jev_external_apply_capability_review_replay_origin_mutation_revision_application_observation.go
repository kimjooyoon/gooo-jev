package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

const jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationNotNeeded = "application-observation-not-needed"
const jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationRecorded = "application-observation-recorded"
const jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationMismatch = "application-observation-mismatch"
const jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationHold = "application-observation-hold"
const jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationRejected = "application-observation-rejected"
const jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationUnknown = "UNKNOWN"

const jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationObserved = "observed"
const jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationMismatched = "mismatch"
const jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationNotNeededStatus = "not-needed"
const jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationHoldStatus = "hold"
const jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationRejectedStatus = "rejected"

type JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationInput struct {
    Candidate                 JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidate
    ObservationStatus         string
    ObservationSource         string
    ObservationEvidenceDigest string
    ReverseObservationDigest  string
    NonAuthorizing            bool
}

type JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservation struct {
    Status                    string
    MissingStage              string
    CandidateStatus            string
    CandidateDigest            string
    ObservationStatus          string
    ObservationSource          string
    ObservationEvidenceDigest  string
    ReverseObservationDigest   string
    ObservationDigest          string
    NonExecuting               bool
    NonAuthorizing             bool
}

func (o JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservation) Validate() error {
    if o.Status == "" || o.CandidateStatus == "" || o.CandidateDigest == "" || o.ObservationStatus == "" || o.ObservationDigest == "" {
        return fmt.Errorf("incomplete JEV revision application observation")
    }
    if o.Status != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationNotNeeded &&
        o.Status != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationRecorded &&
        o.Status != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationMismatch &&
        o.Status != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationHold &&
        o.Status != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationRejected {
        return fmt.Errorf("invalid JEV revision application observation status")
    }
    switch o.Status {
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationNotNeeded:
        if o.CandidateStatus != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateNotNeeded ||
            o.ObservationStatus != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationNotNeededStatus ||
            o.ObservationSource != "" || o.ObservationEvidenceDigest != "" ||
            o.ReverseObservationDigest != "" {
            return fmt.Errorf("not-needed application observation has inconsistent evidence")
        }
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationRecorded:
        if o.CandidateStatus != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateReady ||
            o.ObservationStatus != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationObserved ||
            o.ObservationSource == "" || o.ObservationEvidenceDigest == "" ||
            o.ReverseObservationDigest == "" {
            return fmt.Errorf("recorded application observation is incomplete")
        }
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationMismatch:
        if o.CandidateStatus != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateReady ||
            o.ObservationStatus != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationMismatched ||
            o.ObservationSource == "" || o.ObservationEvidenceDigest == "" ||
            o.ReverseObservationDigest == "" {
            return fmt.Errorf("mismatched application observation is incomplete")
        }
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationHold:
        if o.CandidateStatus != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateHold ||
            o.ObservationStatus != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationHoldStatus {
            return fmt.Errorf("held application observation has inconsistent evidence")
        }
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationRejected:
        if o.CandidateStatus != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateRejected ||
            o.ObservationStatus != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationRejectedStatus {
            return fmt.Errorf("rejected application observation has inconsistent evidence")
        }
    }
    if !o.NonExecuting {
        return fmt.Errorf("JEV revision application observation must be non-executing")
    }
    if !o.NonAuthorizing {
        return fmt.Errorf("JEV revision application observation must be non-authorizing")
    }
    expected := digestJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservation(
        o.Status,
        o.CandidateStatus,
        o.CandidateDigest,
        o.ObservationStatus,
        o.ObservationSource,
        o.ObservationEvidenceDigest,
        o.ReverseObservationDigest,
    )
    if o.ObservationDigest != expected {
        return fmt.Errorf("JEV revision application observation digest mismatch")
    }
    return nil
}

func ObserveJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplication(input JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationInput) JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservation {
    output := JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservation{
        Status:           jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationUnknown,
        CandidateStatus:  input.Candidate.Status,
        CandidateDigest:  input.Candidate.CandidateDigest,
        NonExecuting:     true,
        NonAuthorizing:   true,
    }
    if !input.NonAuthorizing || !input.Candidate.NonAuthorizing {
        output.MissingStage = "authorization-boundary"
        return output
    }
    if !input.Candidate.NonExecuting {
        output.MissingStage = "execution-boundary"
        return output
    }
    if err := input.Candidate.Validate(); err != nil {
        output.MissingStage = "revision-application-candidate"
        return output
    }
    switch input.Candidate.Status {
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateNotNeeded:
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationNotNeeded
        output.ObservationStatus = jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationNotNeededStatus
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateReady:
        if input.ObservationStatus == "" {
            output.MissingStage = "application-observation-status"
            return output
        }
        if input.ObservationStatus != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationObserved &&
            input.ObservationStatus != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationMismatched {
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
        if input.ObservationStatus == jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationObserved {
            output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationRecorded
        } else {
            output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationMismatch
        }
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateHold:
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationHold
        output.ObservationStatus = jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationHoldStatus
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateRejected:
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationRejected
        output.ObservationStatus = jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationRejectedStatus
    default:
        output.MissingStage = "application-candidate-status"
        return output
    }
    output.ObservationDigest = digestJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservation(
        output.Status,
        output.CandidateStatus,
        output.CandidateDigest,
        output.ObservationStatus,
        output.ObservationSource,
        output.ObservationEvidenceDigest,
        output.ReverseObservationDigest,
    )
    if err := output.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationUnknown
        output.MissingStage = "application-observation-evidence"
        output.ObservationDigest = ""
    }
    return output
}

func digestJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservation(status, candidateStatus, candidateDigest, observationStatus, observationSource, observationEvidenceDigest, reverseObservationDigest string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s", status, candidateStatus, candidateDigest, observationStatus, observationSource, observationEvidenceDigest, reverseObservationDigest)))
    return hex.EncodeToString(sum[:])
}