package decision

import "fmt"

type JEVExternalApplyCapabilityReviewReplayObservationLSPDiagnostic struct {
    Severity                 string
    Code                    string
    Message                 string
    Status                   string
    MissingStage             string
    PreparationDecision      string
    ReplayScope              string
    ReplayEnvironmentDigest  string
    PreparationDigest        string
    CandidateDigest          string
    ObservedArtifactDigest   string
    ReverseObservationDigest string
    ObservationDigest        string
    Publishable              bool
    NonExecuting             bool
    NonAuthorizing           bool
}

func (d JEVExternalApplyCapabilityReviewReplayObservationLSPDiagnostic) Validate() error {
    if d.Severity == "" || d.Code == "" || d.Message == "" || d.Status == "" {
        return fmt.Errorf("incomplete external apply capability review replay observation LSP diagnostic")
    }
    if d.PreparationDecision != jevExternalApplyCapabilityReviewReplayReady &&
        d.PreparationDecision != jevExternalApplyCapabilityReviewReplayHold &&
        d.PreparationDecision != jevExternalApplyCapabilityReviewReplayRejected {
        return fmt.Errorf("invalid external apply capability review replay observation LSP decision")
    }
    if d.PreparationDecision == jevExternalApplyCapabilityReviewReplayReady {
        if d.Status != jevExternalApplyCapabilityReviewReplayObserved ||
            d.ObservedArtifactDigest == "" ||
            d.ReverseObservationDigest == "" {
            return fmt.Errorf("replay-observed LSP diagnostic requires forward and reverse observation evidence")
        }
    }
    if d.PreparationDecision == jevExternalApplyCapabilityReviewReplayHold &&
        d.Status != jevExternalApplyCapabilityReviewReplayObservationHold {
        return fmt.Errorf("replay-hold LSP diagnostic requires replay-observation-hold status")
    }
    if d.PreparationDecision == jevExternalApplyCapabilityReviewReplayRejected &&
        d.Status != jevExternalApplyCapabilityReviewReplayObservationRejected {
        return fmt.Errorf("replay-rejected LSP diagnostic requires replay-observation-rejected status")
    }
    if !d.NonExecuting {
        return fmt.Errorf("external apply capability review replay observation LSP diagnostic must be non-executing")
    }
    if !d.NonAuthorizing {
        return fmt.Errorf("external apply capability review replay observation LSP diagnostic must be non-authorizing")
    }
    if d.Publishable && (d.ReplayScope == "" || d.PreparationDigest == "" || d.ObservationDigest == "") {
        return fmt.Errorf("publishable external apply capability review replay observation LSP diagnostic requires evidence")
    }
    return nil
}

func ProjectJEVExternalApplyCapabilityReviewReplayObservationLSP(input JEVExternalApplyCapabilityReviewReplayObservation) JEVExternalApplyCapabilityReviewReplayObservationLSPDiagnostic {
    output := JEVExternalApplyCapabilityReviewReplayObservationLSPDiagnostic{
        Status:                   input.Status,
        MissingStage:             input.MissingStage,
        PreparationDecision:      input.PreparationDecision,
        ReplayScope:              input.ReplayScope,
        ReplayEnvironmentDigest:  input.ReplayEnvironmentDigest,
        PreparationDigest:        input.PreparationDigest,
        CandidateDigest:          input.CandidateDigest,
        ObservedArtifactDigest:   input.ObservedArtifactDigest,
        ReverseObservationDigest: input.ReverseObservationDigest,
        ObservationDigest:        input.ObservationDigest,
        NonExecuting:             input.NonExecuting,
        NonAuthorizing:           input.NonAuthorizing,
    }
    if err := input.Validate(); err != nil {
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "JEV replay observation is not publishable: forward or reverse evidence is incomplete"
        output.Publishable = false
        return output
    }
    output.Publishable = true
    switch input.Status {
    case jevExternalApplyCapabilityReviewReplayObserved:
        output.Severity = "info"
        output.Code = "jev.external-apply.replay-observed"
        output.Message = "Replay observation includes forward and reverse evidence; no execution or authorization occurred"
    case jevExternalApplyCapabilityReviewReplayObservationHold:
        output.Severity = "info"
        output.Code = "jev.external-apply.replay-observation-hold"
        output.Message = "Replay observation remains on hold; no execution or authorization occurred"
    case jevExternalApplyCapabilityReviewReplayObservationRejected:
        output.Severity = "warning"
        output.Code = "jev.external-apply.replay-observation-rejected"
        output.Message = "Replay observation is rejected; no execution or authorization occurred"
    default:
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "JEV replay observation is UNKNOWN; evidence must be resolved before publication"
        output.Publishable = false
    }
    return output
}
