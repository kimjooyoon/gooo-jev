package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

const jevExternalApplyCapabilityReviewReplayOriginMutationCandidateNotNeeded = "mutation-candidate-not-needed"
const jevExternalApplyCapabilityReviewReplayOriginMutationCandidateGenerated = "mutation-candidate-generated"
const jevExternalApplyCapabilityReviewReplayOriginMutationCandidateHold = "mutation-candidate-hold"
const jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRejected = "mutation-candidate-rejected"
const jevExternalApplyCapabilityReviewReplayOriginMutationCandidateUnknown = "UNKNOWN"

type JEVExternalApplyCapabilityReviewReplayOriginMutationCandidateGenerationInput struct {
    Observation       JEVExternalApplyCapabilityReviewReplayOriginMutationObservation
    CandidateDigest   string
    CandidateSource   string
    NonAuthorizing    bool
}

type JEVExternalApplyCapabilityReviewReplayOriginMutationCandidateGeneration struct {
    Status                    string
    MissingStage              string
    ObservationStatus         string
    CandidateDigest           string
    CandidateSource           string
    ObservationDigest         string
    CandidateGenerationDigest string
    NonExecuting              bool
    NonAuthorizing            bool
}

func (g JEVExternalApplyCapabilityReviewReplayOriginMutationCandidateGeneration) Validate() error {
    if g.Status == "" || g.ObservationStatus == "" || g.ObservationDigest == "" || g.CandidateGenerationDigest == "" {
        return fmt.Errorf("incomplete JEV external apply capability review replay origin mutation candidate generation")
    }
    if g.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateNotNeeded &&
        g.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateGenerated &&
        g.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateHold &&
        g.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRejected {
        return fmt.Errorf("invalid JEV mutation candidate generation status")
    }
    if g.ObservationStatus != jevExternalApplyCapabilityReviewReplayOriginMutationObservationNotNeeded &&
        g.ObservationStatus != jevExternalApplyCapabilityReviewReplayOriginMutationObservationRecorded &&
        g.ObservationStatus != jevExternalApplyCapabilityReviewReplayOriginMutationObservationHold &&
        g.ObservationStatus != jevExternalApplyCapabilityReviewReplayOriginMutationObservationRejected {
        return fmt.Errorf("invalid JEV mutation observation status")
    }
    if g.ObservationStatus == jevExternalApplyCapabilityReviewReplayOriginMutationObservationNotNeeded &&
        g.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateNotNeeded {
        return fmt.Errorf("observation-not-needed requires candidate-not-needed status")
    }
    if g.ObservationStatus == jevExternalApplyCapabilityReviewReplayOriginMutationObservationRecorded {
        if g.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateGenerated ||
            g.CandidateDigest == "" ||
            g.CandidateSource == "" {
            return fmt.Errorf("recorded observation requires candidate digest and source")
        }
    }
    if g.ObservationStatus == jevExternalApplyCapabilityReviewReplayOriginMutationObservationHold &&
        g.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateHold {
        return fmt.Errorf("observation-hold requires candidate-hold status")
    }
    if g.ObservationStatus == jevExternalApplyCapabilityReviewReplayOriginMutationObservationRejected &&
        g.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRejected {
        return fmt.Errorf("observation-rejected requires candidate-rejected status")
    }
    if !g.NonExecuting {
        return fmt.Errorf("JEV mutation candidate generation must be non-executing")
    }
    if !g.NonAuthorizing {
        return fmt.Errorf("JEV mutation candidate generation must be non-authorizing")
    }
    expected := digestJEVExternalApplyCapabilityReviewReplayOriginMutationCandidateGeneration(
        g.Status,
        g.ObservationStatus,
        g.CandidateDigest,
        g.CandidateSource,
        g.ObservationDigest,
    )
    if g.CandidateGenerationDigest != expected {
        return fmt.Errorf("JEV mutation candidate generation digest mismatch")
    }
    return nil
}

func GenerateJEVExternalApplyCapabilityReviewReplayOriginMutationCandidate(input JEVExternalApplyCapabilityReviewReplayOriginMutationCandidateGenerationInput) JEVExternalApplyCapabilityReviewReplayOriginMutationCandidateGeneration {
    output := JEVExternalApplyCapabilityReviewReplayOriginMutationCandidateGeneration{
        Status:         jevExternalApplyCapabilityReviewReplayOriginMutationCandidateUnknown,
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
        output.MissingStage = "mutation-observation"
        return output
    }
    output.ObservationStatus = input.Observation.Status
    output.ObservationDigest = input.Observation.MutationObservationDigest
    switch input.Observation.Status {
    case jevExternalApplyCapabilityReviewReplayOriginMutationObservationNotNeeded:
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationCandidateNotNeeded
    case jevExternalApplyCapabilityReviewReplayOriginMutationObservationRecorded:
        if input.CandidateDigest == "" {
            output.MissingStage = "candidate-digest"
            return output
        }
        if input.CandidateSource == "" {
            output.MissingStage = "candidate-source"
            return output
        }
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationCandidateGenerated
        output.CandidateDigest = input.CandidateDigest
        output.CandidateSource = input.CandidateSource
    case jevExternalApplyCapabilityReviewReplayOriginMutationObservationHold:
        if input.CandidateDigest != "" || input.CandidateSource != "" {
            output.MissingStage = "candidate-generation-lifecycle"
            return output
        }
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationCandidateHold
    case jevExternalApplyCapabilityReviewReplayOriginMutationObservationRejected:
        if input.CandidateDigest != "" || input.CandidateSource != "" {
            output.MissingStage = "candidate-generation-lifecycle"
            return output
        }
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRejected
    default:
        output.MissingStage = "mutation-observation-status"
        return output
    }
    output.CandidateGenerationDigest = digestJEVExternalApplyCapabilityReviewReplayOriginMutationCandidateGeneration(
        output.Status,
        output.ObservationStatus,
        output.CandidateDigest,
        output.CandidateSource,
        output.ObservationDigest,
    )
    if err := output.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationCandidateUnknown
        output.MissingStage = "candidate-generation-evidence"
        output.CandidateGenerationDigest = ""
    }
    return output
}

func digestJEVExternalApplyCapabilityReviewReplayOriginMutationCandidateGeneration(status, observationStatus, candidateDigest, candidateSource, observationDigest string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s|%s", status, observationStatus, candidateDigest, candidateSource, observationDigest)))
    return hex.EncodeToString(sum[:])
}
