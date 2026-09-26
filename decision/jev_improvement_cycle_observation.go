package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

const jevImprovementCycleObservationUnknown = "UNKNOWN"

type JEVImprovementCycleObservationInput struct {
    DeclarationDigest              string
    IRDigest                       string
    GenerationDigest               string
    ReverseObservationDigest       string
    MetricDigest                   string
    Disposition                    DecisionConfidenceChangePlanFeedbackDisposition
    LSPDiagnostic                  DecisionConfidenceChangePlanFeedbackDispositionLSPDiagnostic
    NonAuthorizing                 bool
}

type JEVImprovementCycleObservation struct {
    Status                     string
    MissingStage               string
    DeclarationDigest          string
    IRDigest                   string
    GenerationDigest           string
    ReverseObservationDigest   string
    MetricDigest               string
    ChangePlanDigest           string
    DispositionStatus          string
    LSPCode                    string
    EvidenceDigest             string
    NonExecuting               bool
    NonAuthorizing             bool
}

func (o JEVImprovementCycleObservation) Validate() error {
    if o.Status == "" || o.DeclarationDigest == "" || o.IRDigest == "" || o.GenerationDigest == "" || o.ReverseObservationDigest == "" || o.MetricDigest == "" || o.ChangePlanDigest == "" || o.DispositionStatus == "" || o.LSPCode == "" || o.EvidenceDigest == "" {
        return fmt.Errorf("incomplete JEV improvement cycle observation")
    }
    if !o.NonExecuting {
        return fmt.Errorf("JEV improvement cycle observation must be non-executing")
    }
    if !o.NonAuthorizing {
        return fmt.Errorf("JEV improvement cycle observation must be non-authorizing")
    }
    expected := digestJEVImprovementCycleObservation(o.DeclarationDigest, o.IRDigest, o.GenerationDigest, o.ReverseObservationDigest, o.MetricDigest, o.ChangePlanDigest, o.DispositionStatus, o.LSPCode)
    if o.EvidenceDigest != expected {
        return fmt.Errorf("JEV improvement cycle observation evidence digest mismatch")
    }
    return nil
}

func ObserveJEVImprovementCycle(input JEVImprovementCycleObservationInput) JEVImprovementCycleObservation {
    output := JEVImprovementCycleObservation{
        Status: jevImprovementCycleObservationUnknown,
        NonExecuting: true,
        NonAuthorizing: true,
    }
    if !input.NonAuthorizing || !input.Disposition.NonAuthorizing || !input.LSPDiagnostic.NonAuthorizing {
        output.MissingStage = "authorization-boundary"
        return output
    }
    if !input.Disposition.NonExecuting {
        output.MissingStage = "execution-boundary"
        return output
    }
    stages := []struct {
        name  string
        value string
    }{
        {"declaration", input.DeclarationDigest},
        {"ir", input.IRDigest},
        {"generation", input.GenerationDigest},
        {"reverse-observation", input.ReverseObservationDigest},
        {"metric", input.MetricDigest},
    }
    for _, stage := range stages {
        if stage.value == "" {
            output.MissingStage = stage.name
            return output
        }
    }
    if err := input.Disposition.Validate(); err != nil {
        output.MissingStage = "feedback-disposition"
        return output
    }
    if err := input.LSPDiagnostic.Validate(); err != nil || !input.LSPDiagnostic.Publishable {
        output.MissingStage = "lsp-diagnostic"
        return output
    }
    if input.LSPDiagnostic.Status != input.Disposition.Status {
        output.MissingStage = "lsp-status-binding"
        return output
    }
    output.DeclarationDigest = input.DeclarationDigest
    output.IRDigest = input.IRDigest
    output.GenerationDigest = input.GenerationDigest
    output.ReverseObservationDigest = input.ReverseObservationDigest
    output.MetricDigest = input.MetricDigest
    output.ChangePlanDigest = input.Disposition.ChangePlanDigest
    output.DispositionStatus = input.Disposition.Status
    output.LSPCode = input.LSPDiagnostic.Code
    output.EvidenceDigest = digestJEVImprovementCycleObservation(output.DeclarationDigest, output.IRDigest, output.GenerationDigest, output.ReverseObservationDigest, output.MetricDigest, output.ChangePlanDigest, output.DispositionStatus, output.LSPCode)
    if err := output.Validate(); err != nil {
        output.Status = jevImprovementCycleObservationUnknown
        output.MissingStage = "cycle-observation-evidence"
        output.EvidenceDigest = ""
    } else {
        output.Status = input.Disposition.Status
    }
    return output
}

func digestJEVImprovementCycleObservation(declarationDigest, irDigest, generationDigest, reverseObservationDigest, metricDigest, changePlanDigest, dispositionStatus, lspCode string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s", declarationDigest, irDigest, generationDigest, reverseObservationDigest, metricDigest, changePlanDigest, dispositionStatus, lspCode)))
    return hex.EncodeToString(sum[:])
}
