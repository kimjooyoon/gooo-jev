package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

const (
    jevExternalApplyObservedApplied = "applied"
    jevExternalApplyObservedRejected = "rejected"
    jevExternalApplyObservedUnknown = "UNKNOWN"
)

type JEVExternalApplyObservationInput struct {
    Request                JEVExternalApplyRequest
    AppliedTargetReference string
    Outcome                string
    ApplyEvidenceDigest    string
    NonAuthorizing         bool
}

type JEVExternalApplyObservation struct {
    Status                string
    MissingStage          string
    CandidateDigest       string
    RequestEvidenceDigest string
    AppliedTargetReference string
    ApplyEvidenceDigest   string
    EvidenceDigest        string
    NonExecuting          bool
    NonAuthorizing        bool
}

func (o JEVExternalApplyObservation) Validate() error {
    if o.Status == "" || o.CandidateDigest == "" || o.RequestEvidenceDigest == "" || o.AppliedTargetReference == "" || o.ApplyEvidenceDigest == "" || o.EvidenceDigest == "" {
        return fmt.Errorf("incomplete JEV external apply observation")
    }
    switch o.Status {
    case jevExternalApplyObservedApplied, jevExternalApplyObservedRejected, jevExternalApplyObservedUnknown:
    default:
        return fmt.Errorf("invalid JEV external apply observation status")
    }
    if !o.NonExecuting {
        return fmt.Errorf("JEV external apply observation must be non-executing")
    }
    if !o.NonAuthorizing {
        return fmt.Errorf("JEV external apply observation must be non-authorizing")
    }
    expected := digestJEVExternalApplyObservation(o.Status, o.CandidateDigest, o.RequestEvidenceDigest, o.AppliedTargetReference, o.ApplyEvidenceDigest)
    if o.EvidenceDigest != expected {
        return fmt.Errorf("JEV external apply observation evidence mismatch")
    }
    return nil
}

func ObserveJEVExternalApply(input JEVExternalApplyObservationInput) JEVExternalApplyObservation {
    output := JEVExternalApplyObservation{
        Status: jevExternalApplyObservedUnknown,
        NonExecuting: true,
        NonAuthorizing: true,
    }
    if !input.NonAuthorizing || !input.Request.NonAuthorizing {
        output.MissingStage = "authorization-boundary"
        return output
    }
    if !input.Request.NonExecuting {
        output.MissingStage = "execution-boundary"
        return output
    }
    if err := input.Request.Validate(); err != nil {
        output.MissingStage = "external-apply-request"
        return output
    }
    if input.AppliedTargetReference == "" || input.AppliedTargetReference != input.Request.TargetReference {
        output.MissingStage = "target-binding"
        return output
    }
    if input.ApplyEvidenceDigest == "" {
        output.MissingStage = "apply-evidence"
        return output
    }
    switch input.Outcome {
    case jevExternalApplyObservedApplied, jevExternalApplyObservedRejected, jevExternalApplyObservedUnknown:
    default:
        output.MissingStage = "apply-outcome"
        return output
    }
    output.Status = input.Outcome
    output.CandidateDigest = input.Request.CandidateDigest
    output.RequestEvidenceDigest = input.Request.EvidenceDigest
    output.AppliedTargetReference = input.AppliedTargetReference
    output.ApplyEvidenceDigest = input.ApplyEvidenceDigest
    output.EvidenceDigest = digestJEVExternalApplyObservation(output.Status, output.CandidateDigest, output.RequestEvidenceDigest, output.AppliedTargetReference, output.ApplyEvidenceDigest)
    if err := output.Validate(); err != nil {
        output.Status = jevExternalApplyObservedUnknown
        output.MissingStage = "apply-observation-evidence"
        output.EvidenceDigest = ""
    }
    return output
}

func digestJEVExternalApplyObservation(status, candidateDigest, requestEvidenceDigest, appliedTargetReference, applyEvidenceDigest string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s|%s", status, candidateDigest, requestEvidenceDigest, appliedTargetReference, applyEvidenceDigest)))
    return hex.EncodeToString(sum[:])
}
