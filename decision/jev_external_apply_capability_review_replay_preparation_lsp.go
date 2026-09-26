package decision

import "fmt"

type JEVExternalApplyCapabilityReviewReplayPreparationLSPDiagnostic struct {
    Severity                string
    Code                   string
    Message                string
    Status                 string
    MissingStage            string
    LifecycleState          string
    Decision                string
    ReplayScope             string
    ReplayEnvironmentDigest string
    TransitionDigest        string
    CandidateDigest         string
    PreparationDigest       string
    Publishable             bool
    NonExecuting            bool
    NonAuthorizing          bool
}

func (d JEVExternalApplyCapabilityReviewReplayPreparationLSPDiagnostic) Validate() error {
    if d.Severity == "" || d.Code == "" || d.Message == "" || d.Status == "" {
        return fmt.Errorf("incomplete external apply capability review replay preparation LSP diagnostic")
    }
    if d.Decision != jevExternalApplyCapabilityReviewReplayReady &&
        d.Decision != jevExternalApplyCapabilityReviewReplayHold &&
        d.Decision != jevExternalApplyCapabilityReviewReplayRejected {
        return fmt.Errorf("invalid external apply capability review replay preparation LSP decision")
    }
    if d.Decision == jevExternalApplyCapabilityReviewReplayReady &&
        d.ReplayEnvironmentDigest == "" {
        return fmt.Errorf("replay-ready LSP diagnostic requires replay environment digest")
    }
    if !d.NonExecuting {
        return fmt.Errorf("external apply capability review replay preparation LSP diagnostic must be non-executing")
    }
    if !d.NonAuthorizing {
        return fmt.Errorf("external apply capability review replay preparation LSP diagnostic must be non-authorizing")
    }
    if d.Publishable && (d.ReplayScope == "" || d.TransitionDigest == "" || d.PreparationDigest == "") {
        return fmt.Errorf("publishable external apply capability review replay preparation LSP diagnostic requires evidence")
    }
    return nil
}

func ProjectJEVExternalApplyCapabilityReviewReplayPreparationLSP(input JEVExternalApplyCapabilityReviewReplayPreparation) JEVExternalApplyCapabilityReviewReplayPreparationLSPDiagnostic {
    output := JEVExternalApplyCapabilityReviewReplayPreparationLSPDiagnostic{
        Status:                  input.Status,
        MissingStage:            input.MissingStage,
        LifecycleState:          input.LifecycleState,
        Decision:                input.Decision,
        ReplayScope:             input.ReplayScope,
        ReplayEnvironmentDigest: input.ReplayEnvironmentDigest,
        TransitionDigest:        input.TransitionDigest,
        CandidateDigest:         input.CandidateDigest,
        PreparationDigest:       input.PreparationDigest,
        NonExecuting:            input.NonExecuting,
        NonAuthorizing:          input.NonAuthorizing,
    }
    if err := input.Validate(); err != nil {
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "JEV replay preparation is not publishable: evidence is incomplete or invalid"
        output.Publishable = false
        return output
    }
    output.Publishable = true
    switch input.Decision {
    case jevExternalApplyCapabilityReviewReplayReady:
        output.Severity = "info"
        output.Code = "jev.external-apply.replay-ready"
        output.Message = "Replay preparation is ready for observation; no execution or authorization occurred"
    case jevExternalApplyCapabilityReviewReplayHold:
        output.Severity = "info"
        output.Code = "jev.external-apply.replay-hold"
        output.Message = "Replay preparation is on hold pending further evidence"
    case jevExternalApplyCapabilityReviewReplayRejected:
        output.Severity = "warning"
        output.Code = "jev.external-apply.replay-rejected"
        output.Message = "Replay preparation is rejected; no execution or authorization occurred"
    default:
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "JEV replay preparation is UNKNOWN; evidence must be resolved before observation"
        output.Publishable = false
    }
    return output
}
