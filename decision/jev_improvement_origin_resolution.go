package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

const jevImprovementOriginResolved = "origin-resolved"
const jevImprovementOriginUnknown = "UNKNOWN"

type JEVImprovementOriginResolutionInput struct {
    Observation    JEVImprovementCycleObservation
    Selection      JEVImprovementCandidateSelection
    NonAuthorizing bool
}

type JEVImprovementOriginResolution struct {
    Status                    string
    MissingStage              string
    DeclarationDigest         string
    IRDigest                  string
    GenerationDigest          string
    ReverseObservationDigest  string
    MetricDigest              string
    CandidateDigest           string
    CandidateSource           string
    ObservationEvidenceDigest string
    SelectionEvidenceDigest   string
    OriginDigest              string
    EvidenceDigest            string
    NonExecuting              bool
    NonAuthorizing            bool
}

func (o JEVImprovementOriginResolution) Validate() error {
    if o.Status == "" || o.DeclarationDigest == "" || o.IRDigest == "" || o.GenerationDigest == "" || o.ReverseObservationDigest == "" || o.MetricDigest == "" || o.CandidateDigest == "" || o.CandidateSource == "" || o.ObservationEvidenceDigest == "" || o.SelectionEvidenceDigest == "" || o.OriginDigest == "" || o.EvidenceDigest == "" {
        return fmt.Errorf("incomplete JEV improvement origin resolution")
    }
    if o.Status != jevImprovementOriginResolved {
        return fmt.Errorf("invalid JEV improvement origin resolution status")
    }
    if !o.NonExecuting {
        return fmt.Errorf("JEV improvement origin resolution must be non-executing")
    }
    if !o.NonAuthorizing {
        return fmt.Errorf("JEV improvement origin resolution must be non-authorizing")
    }
    expectedOrigin := digestJEVImprovementOrigin(o.DeclarationDigest, o.IRDigest, o.GenerationDigest, o.ReverseObservationDigest, o.MetricDigest, o.CandidateDigest, o.CandidateSource)
    if o.OriginDigest != expectedOrigin {
        return fmt.Errorf("JEV improvement origin digest mismatch")
    }
    expected := digestJEVImprovementOriginEvidence(o.Status, o.OriginDigest, o.ObservationEvidenceDigest, o.SelectionEvidenceDigest)
    if o.EvidenceDigest != expected {
        return fmt.Errorf("JEV improvement origin evidence mismatch")
    }
    return nil
}

func ResolveJEVImprovementOrigin(input JEVImprovementOriginResolutionInput) JEVImprovementOriginResolution {
    output := JEVImprovementOriginResolution{
        Status: jevImprovementOriginUnknown,
        NonExecuting: true,
        NonAuthorizing: true,
    }
    if !input.NonAuthorizing || !input.Observation.NonAuthorizing || !input.Selection.NonAuthorizing {
        output.MissingStage = "authorization-boundary"
        return output
    }
    if !input.Observation.NonExecuting || !input.Selection.NonExecuting {
        output.MissingStage = "execution-boundary"
        return output
    }
    if err := input.Observation.Validate(); err != nil {
        output.MissingStage = "improvement-cycle-observation"
        return output
    }
    if err := input.Selection.Validate(); err != nil {
        output.MissingStage = "candidate-selection"
        return output
    }
    if input.Selection.ObservationEvidenceDigest != input.Observation.EvidenceDigest {
        output.MissingStage = "observation-selection-binding"
        return output
    }
    output.Status = jevImprovementOriginResolved
    output.DeclarationDigest = input.Observation.DeclarationDigest
    output.IRDigest = input.Observation.IRDigest
    output.GenerationDigest = input.Observation.GenerationDigest
    output.ReverseObservationDigest = input.Observation.ReverseObservationDigest
    output.MetricDigest = input.Observation.MetricDigest
    output.CandidateDigest = input.Selection.CandidateDigest
    output.CandidateSource = input.Selection.CandidateSource
    output.ObservationEvidenceDigest = input.Observation.EvidenceDigest
    output.SelectionEvidenceDigest = input.Selection.EvidenceDigest
    output.OriginDigest = digestJEVImprovementOrigin(output.DeclarationDigest, output.IRDigest, output.GenerationDigest, output.ReverseObservationDigest, output.MetricDigest, output.CandidateDigest, output.CandidateSource)
    output.EvidenceDigest = digestJEVImprovementOriginEvidence(output.Status, output.OriginDigest, output.ObservationEvidenceDigest, output.SelectionEvidenceDigest)
    if err := output.Validate(); err != nil {
        output.Status = jevImprovementOriginUnknown
        output.MissingStage = "origin-resolution-evidence"
        output.EvidenceDigest = ""
    }
    return output
}

func digestJEVImprovementOrigin(declarationDigest, irDigest, generationDigest, reverseObservationDigest, metricDigest, candidateDigest, candidateSource string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s", declarationDigest, irDigest, generationDigest, reverseObservationDigest, metricDigest, candidateDigest, candidateSource)))
    return hex.EncodeToString(sum[:])
}

func digestJEVImprovementOriginEvidence(status, originDigest, observationEvidenceDigest, selectionEvidenceDigest string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s", status, originDigest, observationEvidenceDigest, selectionEvidenceDigest)))
    return hex.EncodeToString(sum[:])
}
