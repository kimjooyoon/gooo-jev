package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

const jevExternalApplyCapabilityReviewObservationMetricDirectionCandidateGateBridgeBound = "observation-metric-direction-candidate-gate-bound"
const jevExternalApplyCapabilityReviewObservationMetricDirectionCandidateGateBridgeUnknown = "UNKNOWN"

type JEVExternalApplyCapabilityReviewObservationMetricDirectionCandidateGateBridgeInput struct {
    ObservationMetricDirection JEVExternalApplyCapabilityReviewObservationMetricDirectionBridge
    CandidateDigest             string
    RevisionSource              string
    NonAuthorizing              bool
}

type JEVExternalApplyCapabilityReviewObservationMetricDirectionCandidateGateBridge struct {
    Status                        string
    MissingStage                  string
    ObservationStatus             string
    ObservationMetricDirectionDigest string
    CandidateGateStatus           string
    CandidateDecision             string
    CandidateDigest              string
    RevisionSource                string
    GateDigest                    string
    BridgeDigest                  string
    NonExecuting                  bool
    NonAuthorizing                bool
}

func (b JEVExternalApplyCapabilityReviewObservationMetricDirectionCandidateGateBridge) Validate() error {
    if b.Status == "" || b.ObservationStatus == "" || b.ObservationMetricDirectionDigest == "" ||
        b.CandidateGateStatus == "" || b.CandidateDecision == "" || b.GateDigest == "" ||
        b.BridgeDigest == "" {
        return fmt.Errorf("incomplete JEV observation metric direction candidate gate bridge")
    }
    if b.Status != jevExternalApplyCapabilityReviewObservationMetricDirectionCandidateGateBridgeBound {
        return fmt.Errorf("invalid JEV observation metric direction candidate gate bridge status")
    }
    if !b.NonExecuting {
        return fmt.Errorf("JEV observation metric direction candidate gate bridge must be non-executing")
    }
    if !b.NonAuthorizing {
        return fmt.Errorf("JEV observation metric direction candidate gate bridge must be non-authorizing")
    }
    if b.CandidateDecision == jevExternalApplyCapabilityReviewRevisionCandidateReady && b.CandidateDigest == "" {
        return fmt.Errorf("ready candidate gate bridge requires candidate digest")
    }
    expected := digestJEVExternalApplyCapabilityReviewObservationMetricDirectionCandidateGateBridge(
        b.Status,
        b.ObservationStatus,
        b.ObservationMetricDirectionDigest,
        b.CandidateGateStatus,
        b.CandidateDecision,
        b.CandidateDigest,
        b.RevisionSource,
        b.GateDigest,
    )
    if b.BridgeDigest != expected {
        return fmt.Errorf("JEV observation metric direction candidate gate bridge digest mismatch")
    }
    return nil
}

func BindJEVExternalApplyCapabilityReviewObservationMetricDirectionCandidateGate(input JEVExternalApplyCapabilityReviewObservationMetricDirectionCandidateGateBridgeInput) JEVExternalApplyCapabilityReviewObservationMetricDirectionCandidateGateBridge {
    output := JEVExternalApplyCapabilityReviewObservationMetricDirectionCandidateGateBridge{
        Status:           jevExternalApplyCapabilityReviewObservationMetricDirectionCandidateGateBridgeUnknown,
        NonExecuting:     true,
        NonAuthorizing:   true,
    }
    if !input.NonAuthorizing || !input.ObservationMetricDirection.NonAuthorizing {
        output.MissingStage = "authorization-boundary"
        return output
    }
    if !input.ObservationMetricDirection.NonExecuting {
        output.MissingStage = "execution-boundary"
        return output
    }
    if err := input.ObservationMetricDirection.Validate(); err != nil {
        output.MissingStage = "observation-metric-direction-bridge"
        return output
    }
    if input.ObservationMetricDirection.ObservationStatus == jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationMismatch &&
        input.ObservationMetricDirection.Direction == jevExternalApplyCapabilityReviewImprovementDirectionGenerate {
        output.MissingStage = "observation-direction-consistency"
        return output
    }
    direction := JEVExternalApplyCapabilityReviewImprovementDirection{
        Status:          jevExternalApplyCapabilityReviewImprovementDirectionBound,
        Direction:       input.ObservationMetricDirection.Direction,
        Target:          input.ObservationMetricDirection.Target,
        CandidateSource: input.ObservationMetricDirection.CandidateSource,
        FeedbackDigest:  input.ObservationMetricDirection.FeedbackDigest,
        DirectionDigest: input.ObservationMetricDirection.DirectionDigest,
        NonExecuting:    true,
        NonAuthorizing:  true,
    }
    gate := GateJEVExternalApplyCapabilityReviewRevisionCandidate(JEVExternalApplyCapabilityReviewRevisionCandidateGateInput{
        Direction:       direction,
        CandidateDigest: input.CandidateDigest,
        RevisionSource:  input.RevisionSource,
        NonAuthorizing:  true,
    })
    if gate.Status == jevExternalApplyCapabilityReviewRevisionCandidateUnknown {
        output.MissingStage = gate.MissingStage
        if output.MissingStage == "" {
            output.MissingStage = "revision-candidate-gate"
        }
        return output
    }
    output.Status = jevExternalApplyCapabilityReviewObservationMetricDirectionCandidateGateBridgeBound
    output.ObservationStatus = input.ObservationMetricDirection.ObservationStatus
    output.ObservationMetricDirectionDigest = input.ObservationMetricDirection.BridgeDigest
    output.CandidateGateStatus = gate.Status
    output.CandidateDecision = gate.Decision
    output.CandidateDigest = gate.CandidateDigest
    output.RevisionSource = input.RevisionSource
    output.GateDigest = gate.GateDigest
    output.BridgeDigest = digestJEVExternalApplyCapabilityReviewObservationMetricDirectionCandidateGateBridge(
        output.Status,
        output.ObservationStatus,
        output.ObservationMetricDirectionDigest,
        output.CandidateGateStatus,
        output.CandidateDecision,
        output.CandidateDigest,
        output.RevisionSource,
        output.GateDigest,
    )
    if err := output.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewObservationMetricDirectionCandidateGateBridgeUnknown
        output.MissingStage = "observation-metric-direction-candidate-gate-evidence"
        output.ObservationMetricDirectionDigest = ""
        output.GateDigest = ""
        output.BridgeDigest = ""
    }
    return output
}

func digestJEVExternalApplyCapabilityReviewObservationMetricDirectionCandidateGateBridge(status, observationStatus, observationMetricDirectionDigest, candidateGateStatus, candidateDecision, candidateDigest, revisionSource, gateDigest string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s", status, observationStatus, observationMetricDirectionDigest, candidateGateStatus, candidateDecision, candidateDigest, revisionSource, gateDigest)))
    return hex.EncodeToString(sum[:])
}
