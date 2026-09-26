package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

const (
    jevImprovementCandidateReady = "candidate-ready-for-external-review"
    jevImprovementCandidateUnknown = "UNKNOWN"
)

type JEVImprovementCandidateSelectionInput struct {
    Observation    JEVImprovementCycleObservation
    CandidateDigest string
    CandidateSource string
    NonAuthorizing bool
}

type JEVImprovementCandidateSelection struct {
    Status                 string
    MissingStage           string
    CandidateDigest        string
    CandidateSource        string
    ObservationEvidenceDigest string
    EvidenceDigest         string
    NonExecuting           bool
    NonAuthorizing         bool
}

func (s JEVImprovementCandidateSelection) Validate() error {
    if s.Status == "" || s.CandidateDigest == "" || s.CandidateSource == "" || s.ObservationEvidenceDigest == "" || s.EvidenceDigest == "" {
        return fmt.Errorf("incomplete JEV improvement candidate selection")
    }
    if !s.NonExecuting {
        return fmt.Errorf("JEV improvement candidate selection must be non-executing")
    }
    if !s.NonAuthorizing {
        return fmt.Errorf("JEV improvement candidate selection must be non-authorizing")
    }
    expected := digestJEVImprovementCandidateSelection(s.Status, s.CandidateDigest, s.CandidateSource, s.ObservationEvidenceDigest)
    if s.EvidenceDigest != expected {
        return fmt.Errorf("JEV improvement candidate selection evidence digest mismatch")
    }
    return nil
}

func SelectJEVImprovementCandidate(input JEVImprovementCandidateSelectionInput) JEVImprovementCandidateSelection {
    output := JEVImprovementCandidateSelection{
        Status: jevImprovementCandidateUnknown,
        NonExecuting: true,
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
        output.MissingStage = "improvement-cycle-observation"
        return output
    }
    if input.Observation.Status != changePlanFeedbackDispositionReady {
        output.MissingStage = "external-review-disposition"
        return output
    }
    if input.CandidateDigest == "" {
        output.MissingStage = "candidate-digest"
        return output
    }
    if input.CandidateSource == "" {
        output.MissingStage = "candidate-source"
        return output
    }
    output.Status = jevImprovementCandidateReady
    output.CandidateDigest = input.CandidateDigest
    output.CandidateSource = input.CandidateSource
    output.ObservationEvidenceDigest = input.Observation.EvidenceDigest
    output.EvidenceDigest = digestJEVImprovementCandidateSelection(output.Status, output.CandidateDigest, output.CandidateSource, output.ObservationEvidenceDigest)
    if err := output.Validate(); err != nil {
        output.Status = jevImprovementCandidateUnknown
        output.MissingStage = "candidate-selection-evidence"
        output.EvidenceDigest = ""
    }
    return output
}

func digestJEVImprovementCandidateSelection(status, candidateDigest, candidateSource, observationEvidenceDigest string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s", status, candidateDigest, candidateSource, observationEvidenceDigest)))
    return hex.EncodeToString(sum[:])
}
