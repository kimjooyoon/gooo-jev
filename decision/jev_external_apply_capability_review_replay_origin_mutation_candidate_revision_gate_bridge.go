package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

const jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateNotNeeded = "revision-candidate-not-needed"
const jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateReady = "revision-candidate-ready"
const jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateHold = "revision-candidate-hold"
const jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateRejected = "revision-candidate-rejected"
const jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateUnknown = "UNKNOWN"

type JEVExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateBridgeInput struct {
    Direction        JEVExternalApplyCapabilityReviewImprovementDirection
    Candidate        JEVExternalApplyCapabilityReviewReplayOriginMutationCandidateGeneration
    RevisionSource   string
    NonAuthorizing   bool
}

type JEVExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateBridge struct {
    Status                    string
    MissingStage              string
    CandidateStatus           string
    CandidateGenerationDigest string
    Gate                      JEVExternalApplyCapabilityReviewRevisionCandidateGate
    BridgeDigest              string
    NonExecuting              bool
    NonAuthorizing            bool
}

func (b JEVExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateBridge) Validate() error {
    if b.Status == "" || b.CandidateStatus == "" || b.CandidateGenerationDigest == "" || b.BridgeDigest == "" {
        return fmt.Errorf("incomplete JEV mutation candidate revision gate bridge")
    }
    if b.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateNotNeeded &&
        b.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateReady &&
        b.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateHold &&
        b.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateRejected {
        return fmt.Errorf("invalid JEV mutation candidate revision gate bridge status")
    }
    if b.CandidateStatus != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateNotNeeded &&
        b.CandidateStatus != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateGenerated &&
        b.CandidateStatus != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateHold &&
        b.CandidateStatus != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRejected {
        return fmt.Errorf("invalid JEV mutation candidate generation status")
    }
    if b.CandidateStatus == jevExternalApplyCapabilityReviewReplayOriginMutationCandidateNotNeeded {
        if b.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateNotNeeded || b.Gate.Status != "" {
            return fmt.Errorf("candidate-not-needed requires an empty revision gate")
        }
    }
    if b.CandidateStatus == jevExternalApplyCapabilityReviewReplayOriginMutationCandidateGenerated {
        if b.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateReady {
            return fmt.Errorf("generated candidate requires revision-candidate-ready status")
        }
        if err := b.Gate.Validate(); err != nil {
            return fmt.Errorf("invalid bridged revision candidate gate: %w", err)
        }
        if b.Gate.Decision != jevExternalApplyCapabilityReviewRevisionCandidateReady {
            return fmt.Errorf("generated candidate requires ready revision gate decision")
        }
    }
    if b.CandidateStatus == jevExternalApplyCapabilityReviewReplayOriginMutationCandidateHold {
        if b.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateHold {
            return fmt.Errorf("held candidate requires revision-candidate-hold status")
        }
        if err := b.Gate.Validate(); err != nil {
            return fmt.Errorf("invalid held revision candidate gate: %w", err)
        }
        if b.Gate.Decision != jevExternalApplyCapabilityReviewRevisionCandidateHold {
            return fmt.Errorf("held candidate requires hold revision gate decision")
        }
    }
    if b.CandidateStatus == jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRejected {
        if b.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateRejected {
            return fmt.Errorf("rejected candidate requires revision-candidate-rejected status")
        }
        if err := b.Gate.Validate(); err != nil {
            return fmt.Errorf("invalid rejected revision candidate gate: %w", err)
        }
        if b.Gate.Decision != jevExternalApplyCapabilityReviewRevisionCandidateRejected {
            return fmt.Errorf("rejected candidate requires rejected revision gate decision")
        }
    }
    if !b.NonExecuting {
        return fmt.Errorf("JEV mutation candidate revision gate bridge must be non-executing")
    }
    if !b.NonAuthorizing {
        return fmt.Errorf("JEV mutation candidate revision gate bridge must be non-authorizing")
    }
    expected := digestJEVExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateBridge(
        b.Status,
        b.CandidateStatus,
        b.CandidateGenerationDigest,
        b.Gate.GateDigest,
    )
    if b.BridgeDigest != expected {
        return fmt.Errorf("JEV mutation candidate revision gate bridge digest mismatch")
    }
    return nil
}

func GateJEVExternalApplyCapabilityReviewReplayOriginMutationCandidate(input JEVExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateBridgeInput) JEVExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateBridge {
    output := JEVExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateBridge{
        Status:         jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateUnknown,
        NonExecuting:   true,
        NonAuthorizing: true,
    }
    if !input.NonAuthorizing || !input.Candidate.NonAuthorizing || !input.Direction.NonAuthorizing {
        output.MissingStage = "authorization-boundary"
        return output
    }
    if !input.Candidate.NonExecuting || !input.Direction.NonExecuting {
        output.MissingStage = "execution-boundary"
        return output
    }
    if err := input.Candidate.Validate(); err != nil {
        output.MissingStage = "mutation-candidate-generation"
        return output
    }
    if err := input.Direction.Validate(); err != nil {
        output.MissingStage = "improvement-direction"
        return output
    }
    output.CandidateStatus = input.Candidate.Status
    output.CandidateGenerationDigest = input.Candidate.CandidateGenerationDigest
    switch input.Candidate.Status {
    case jevExternalApplyCapabilityReviewReplayOriginMutationCandidateNotNeeded:
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateNotNeeded
    case jevExternalApplyCapabilityReviewReplayOriginMutationCandidateGenerated:
        if input.Direction.Direction != jevExternalApplyCapabilityReviewImprovementDirectionGenerate {
            output.MissingStage = "candidate-direction"
            return output
        }
        if input.RevisionSource == "" {
            output.MissingStage = "revision-source"
            return output
        }
        if input.RevisionSource != input.Candidate.CandidateSource {
            output.MissingStage = "revision-source-binding"
            return output
        }
        output.Gate = GateJEVExternalApplyCapabilityReviewRevisionCandidate(JEVExternalApplyCapabilityReviewRevisionCandidateGateInput{
            Direction:       input.Direction,
            CandidateDigest: input.Candidate.CandidateDigest,
            RevisionSource:  input.RevisionSource,
            NonAuthorizing:  true,
        })
        if output.Gate.Status == jevExternalApplyCapabilityReviewRevisionCandidateUnknown {
            output.MissingStage = output.Gate.MissingStage
            return output
        }
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateReady
    case jevExternalApplyCapabilityReviewReplayOriginMutationCandidateHold:
        if input.Direction.Direction != jevExternalApplyCapabilityReviewImprovementDirectionHold {
            output.MissingStage = "candidate-direction"
            return output
        }
        output.Gate = GateJEVExternalApplyCapabilityReviewRevisionCandidate(JEVExternalApplyCapabilityReviewRevisionCandidateGateInput{
            Direction:       input.Direction,
            NonAuthorizing:  true,
        })
        if output.Gate.Status == jevExternalApplyCapabilityReviewRevisionCandidateUnknown {
            output.MissingStage = output.Gate.MissingStage
            return output
        }
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateHold
    case jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRejected:
        if input.Direction.Direction != jevExternalApplyCapabilityReviewImprovementDirectionReject {
            output.MissingStage = "candidate-direction"
            return output
        }
        output.Gate = GateJEVExternalApplyCapabilityReviewRevisionCandidate(JEVExternalApplyCapabilityReviewRevisionCandidateGateInput{
            Direction:       input.Direction,
            NonAuthorizing:  true,
        })
        if output.Gate.Status == jevExternalApplyCapabilityReviewRevisionCandidateUnknown {
            output.MissingStage = output.Gate.MissingStage
            return output
        }
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateRejected
    default:
        output.MissingStage = "mutation-candidate-status"
        return output
    }
    output.BridgeDigest = digestJEVExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateBridge(
        output.Status,
        output.CandidateStatus,
        output.CandidateGenerationDigest,
        output.Gate.GateDigest,
    )
    if err := output.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateUnknown
        output.MissingStage = "candidate-revision-gate-bridge-evidence"
        output.BridgeDigest = ""
    }
    return output
}

func digestJEVExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateBridge(status, candidateStatus, candidateGenerationDigest, gateDigest string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s", status, candidateStatus, candidateGenerationDigest, gateDigest)))
    return hex.EncodeToString(sum[:])
}
