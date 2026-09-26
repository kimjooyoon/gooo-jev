package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

const (
    jevPostApplyCycleReplayRequired = "replay-required"
    jevPostApplyCycleRevisionRequired = "revision-required"
    jevPostApplyCycleEvidenceRequired = "evidence-required"
    jevPostApplyCycleUnknown = "UNKNOWN"
)

type JEVPostApplyCycleDirectiveInput struct {
    Observation    JEVExternalApplyObservation
    NonAuthorizing bool
}

type JEVPostApplyCycleDirective struct {
    Status                   string
    NextStage                string
    MissingStage             string
    CandidateDigest          string
    ApplyObservationDigest   string
    EvidenceDigest           string
    NonExecuting             bool
    NonAuthorizing           bool
}

func (d JEVPostApplyCycleDirective) Validate() error {
    if d.Status == "" || d.NextStage == "" || d.CandidateDigest == "" || d.ApplyObservationDigest == "" || d.EvidenceDigest == "" {
        return fmt.Errorf("incomplete JEV post-apply cycle directive")
    }
    switch d.Status {
    case jevPostApplyCycleReplayRequired, jevPostApplyCycleRevisionRequired, jevPostApplyCycleEvidenceRequired:
    default:
        return fmt.Errorf("invalid JEV post-apply cycle directive status")
    }
    if !d.NonExecuting {
        return fmt.Errorf("JEV post-apply cycle directive must be non-executing")
    }
    if !d.NonAuthorizing {
        return fmt.Errorf("JEV post-apply cycle directive must be non-authorizing")
    }
    expected := digestJEVPostApplyCycleDirective(d.Status, d.NextStage, d.CandidateDigest, d.ApplyObservationDigest)
    if d.EvidenceDigest != expected {
        return fmt.Errorf("JEV post-apply cycle directive evidence mismatch")
    }
    return nil
}

func DeriveJEVPostApplyCycleDirective(input JEVPostApplyCycleDirectiveInput) JEVPostApplyCycleDirective {
    output := JEVPostApplyCycleDirective{
        Status: jevPostApplyCycleUnknown,
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
        output.MissingStage = "external-apply-observation"
        return output
    }
    output.CandidateDigest = input.Observation.CandidateDigest
    output.ApplyObservationDigest = input.Observation.EvidenceDigest
    switch input.Observation.Status {
    case jevExternalApplyObservedApplied:
        output.Status = jevPostApplyCycleReplayRequired
        output.NextStage = "replay-observation"
    case jevExternalApplyObservedRejected:
        output.Status = jevPostApplyCycleRevisionRequired
        output.NextStage = "revision-candidate"
    case jevExternalApplyObservedUnknown:
        output.Status = jevPostApplyCycleEvidenceRequired
        output.NextStage = "apply-evidence"
    default:
        output.MissingStage = "apply-observation-status"
        return output
    }
    output.EvidenceDigest = digestJEVPostApplyCycleDirective(output.Status, output.NextStage, output.CandidateDigest, output.ApplyObservationDigest)
    if err := output.Validate(); err != nil {
        output.Status = jevPostApplyCycleUnknown
        output.MissingStage = "post-apply-cycle-evidence"
        output.EvidenceDigest = ""
    }
    return output
}

func digestJEVPostApplyCycleDirective(status, nextStage, candidateDigest, applyObservationDigest string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s", status, nextStage, candidateDigest, applyObservationDigest)))
    return hex.EncodeToString(sum[:])
}
