package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

const jevExternalApplyCapabilityReviewRevisionCandidateGated = "revision-candidate-gated"
const jevExternalApplyCapabilityReviewRevisionCandidateReady = "candidate-ready"
const jevExternalApplyCapabilityReviewRevisionCandidateHold = "candidate-hold"
const jevExternalApplyCapabilityReviewRevisionCandidateRejected = "candidate-rejected"
const jevExternalApplyCapabilityReviewRevisionCandidateUnknown = "UNKNOWN"

type JEVExternalApplyCapabilityReviewRevisionCandidateGateInput struct {
    Direction      JEVExternalApplyCapabilityReviewImprovementDirection
    CandidateDigest string
    RevisionSource string
    NonAuthorizing bool
}

type JEVExternalApplyCapabilityReviewRevisionCandidateGate struct {
    Status          string
    MissingStage    string
    Decision        string
    Direction       string
    Target          string
    CandidateSource string
    CandidateDigest string
    DirectionDigest string
    GateDigest      string
    NonExecuting    bool
    NonAuthorizing  bool
}

func (g JEVExternalApplyCapabilityReviewRevisionCandidateGate) Validate() error {
    if g.Status == "" || g.Decision == "" || g.Direction == "" || g.Target == "" || g.CandidateSource == "" || g.DirectionDigest == "" || g.GateDigest == "" {
        return fmt.Errorf("incomplete JEV external apply capability review revision candidate gate")
    }
    if g.Status != jevExternalApplyCapabilityReviewRevisionCandidateGated {
        return fmt.Errorf("invalid JEV external apply capability review revision candidate gate status")
    }
    if g.Decision != jevExternalApplyCapabilityReviewRevisionCandidateReady && g.Decision != jevExternalApplyCapabilityReviewRevisionCandidateHold && g.Decision != jevExternalApplyCapabilityReviewRevisionCandidateRejected {
        return fmt.Errorf("invalid JEV external apply capability review revision candidate gate decision")
    }
    if g.Decision == jevExternalApplyCapabilityReviewRevisionCandidateReady && g.CandidateDigest == "" {
        return fmt.Errorf("ready revision candidate gate requires candidate digest")
    }
    if !g.NonExecuting {
        return fmt.Errorf("JEV external apply capability review revision candidate gate must be non-executing")
    }
    if !g.NonAuthorizing {
        return fmt.Errorf("JEV external apply capability review revision candidate gate must be non-authorizing")
    }
    expected := digestJEVExternalApplyCapabilityReviewRevisionCandidateGate(g.Status, g.Decision, g.Direction, g.Target, g.CandidateSource, g.CandidateDigest, g.DirectionDigest)
    if g.GateDigest != expected {
        return fmt.Errorf("JEV external apply capability review revision candidate gate digest mismatch")
    }
    return nil
}

func GateJEVExternalApplyCapabilityReviewRevisionCandidate(input JEVExternalApplyCapabilityReviewRevisionCandidateGateInput) JEVExternalApplyCapabilityReviewRevisionCandidateGate {
    output := JEVExternalApplyCapabilityReviewRevisionCandidateGate{
        Status:         jevExternalApplyCapabilityReviewRevisionCandidateUnknown,
        NonExecuting:   true,
        NonAuthorizing: true,
    }
    if !input.NonAuthorizing || !input.Direction.NonAuthorizing {
        output.MissingStage = "authorization-boundary"
        return output
    }
    if !input.Direction.NonExecuting {
        output.MissingStage = "execution-boundary"
        return output
    }
    if err := input.Direction.Validate(); err != nil {
        output.MissingStage = "improvement-direction"
        return output
    }
    decision := jevExternalApplyCapabilityReviewRevisionCandidateHold
    candidateSource := input.Direction.CandidateSource
    candidateDigest := ""
    switch input.Direction.Direction {
    case jevExternalApplyCapabilityReviewImprovementDirectionGenerate:
        if input.CandidateDigest == "" {
            output.MissingStage = "candidate-digest"
            return output
        }
        if input.RevisionSource == "" {
            output.MissingStage = "revision-source"
            return output
        }
        decision = jevExternalApplyCapabilityReviewRevisionCandidateReady
        candidateSource = input.RevisionSource
        candidateDigest = input.CandidateDigest
    case jevExternalApplyCapabilityReviewImprovementDirectionHold:
        decision = jevExternalApplyCapabilityReviewRevisionCandidateHold
    case jevExternalApplyCapabilityReviewImprovementDirectionReject:
        decision = jevExternalApplyCapabilityReviewRevisionCandidateRejected
    default:
        output.MissingStage = "direction"
        return output
    }
    output.Status = jevExternalApplyCapabilityReviewRevisionCandidateGated
    output.Decision = decision
    output.Direction = input.Direction.Direction
    output.Target = input.Direction.Target
    output.CandidateSource = candidateSource
    output.CandidateDigest = candidateDigest
    output.DirectionDigest = input.Direction.DirectionDigest
    output.GateDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateGate(output.Status, output.Decision, output.Direction, output.Target, output.CandidateSource, output.CandidateDigest, output.DirectionDigest)
    if err := output.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewRevisionCandidateUnknown
        output.MissingStage = "revision-candidate-gate-evidence"
        output.CandidateDigest = ""
        output.DirectionDigest = ""
        output.GateDigest = ""
    }
    return output
}

func digestJEVExternalApplyCapabilityReviewRevisionCandidateGate(status, decision, direction, target, candidateSource, candidateDigest, directionDigest string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s", status, decision, direction, target, candidateSource, candidateDigest, directionDigest)))
    return hex.EncodeToString(sum[:])
}
