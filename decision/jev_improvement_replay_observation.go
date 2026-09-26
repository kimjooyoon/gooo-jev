package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

const (
    jevImprovementReplayReproduced = "reproduced"
    jevImprovementReplayCounterexample = "counterexample"
    jevImprovementReplayObservationUnknown = "UNKNOWN"
)

type JEVImprovementReplayObservationInput struct {
    Preparation             JEVImprovementReplayPreparation
    ObservedCandidateDigest string
    Outcome                 string
    ObservationEvidenceDigest string
    NonAuthorizing          bool
}

type JEVImprovementReplayObservation struct {
    Status                    string
    MissingStage              string
    CandidateDigest           string
    PreparationEvidenceDigest string
    ObservationEvidenceDigest string
    EvidenceDigest            string
    NonExecuting              bool
    NonAuthorizing            bool
}

func (o JEVImprovementReplayObservation) Validate() error {
    if o.Status == "" || o.CandidateDigest == "" || o.PreparationEvidenceDigest == "" || o.ObservationEvidenceDigest == "" || o.EvidenceDigest == "" {
        return fmt.Errorf("incomplete JEV improvement replay observation")
    }
    switch o.Status {
    case jevImprovementReplayReproduced, jevImprovementReplayCounterexample, jevImprovementReplayObservationUnknown:
    default:
        return fmt.Errorf("invalid JEV improvement replay observation status")
    }
    if !o.NonExecuting {
        return fmt.Errorf("JEV improvement replay observation must be non-executing")
    }
    if !o.NonAuthorizing {
        return fmt.Errorf("JEV improvement replay observation must be non-authorizing")
    }
    expected := digestJEVImprovementReplayObservation(o.Status, o.CandidateDigest, o.PreparationEvidenceDigest, o.ObservationEvidenceDigest)
    if o.EvidenceDigest != expected {
        return fmt.Errorf("JEV improvement replay observation evidence digest mismatch")
    }
    return nil
}

func ObserveJEVImprovementReplay(input JEVImprovementReplayObservationInput) JEVImprovementReplayObservation {
    output := JEVImprovementReplayObservation{
        Status: jevImprovementReplayObservationUnknown,
        NonExecuting: true,
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
    if input.Preparation.Status != jevImprovementReplayReady {
        output.MissingStage = "replay-readiness"
        return output
    }
    if input.ObservedCandidateDigest == "" || input.ObservedCandidateDigest != input.Preparation.CandidateDigest {
        output.MissingStage = "candidate-binding"
        return output
    }
    if input.ObservationEvidenceDigest == "" {
        output.MissingStage = "reverse-observation-evidence"
        return output
    }
    switch input.Outcome {
    case jevImprovementReplayReproduced, jevImprovementReplayCounterexample, jevImprovementReplayObservationUnknown:
    default:
        output.MissingStage = "reverse-observation-outcome"
        return output
    }
    output.Status = input.Outcome
    output.CandidateDigest = input.Preparation.CandidateDigest
    output.PreparationEvidenceDigest = input.Preparation.EvidenceDigest
    output.ObservationEvidenceDigest = input.ObservationEvidenceDigest
    output.EvidenceDigest = digestJEVImprovementReplayObservation(output.Status, output.CandidateDigest, output.PreparationEvidenceDigest, output.ObservationEvidenceDigest)
    if err := output.Validate(); err != nil {
        output.Status = jevImprovementReplayObservationUnknown
        output.MissingStage = "reverse-observation-binding"
        output.EvidenceDigest = ""
    }
    return output
}

func digestJEVImprovementReplayObservation(status, candidateDigest, preparationEvidenceDigest, observationEvidenceDigest string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s", status, candidateDigest, preparationEvidenceDigest, observationEvidenceDigest)))
    return hex.EncodeToString(sum[:])
}
