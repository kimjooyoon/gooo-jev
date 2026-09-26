package decision

import "fmt"

type JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateLSPDiagnostic struct {
    Severity             string
    Code                 string
    Message              string
    Status               string
    MissingStage         string
    PlanStatus           string
    Decision             string
    PlanDigest           string
    ReviewSource         string
    ReviewEvidenceDigest string
    GateDigest           string
    Publishable          bool
    NonExecuting         bool
    NonAuthorizing       bool
}

func (d JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateLSPDiagnostic) Validate() error {
    if d.Severity == "" || d.Code == "" || d.Message == "" || d.Status == "" {
        return fmt.Errorf("incomplete JEV revision plan review gate LSP diagnostic")
    }
    if !d.NonExecuting {
        return fmt.Errorf("JEV revision plan review gate LSP diagnostic must be non-executing")
    }
    if !d.NonAuthorizing {
        return fmt.Errorf("JEV revision plan review gate LSP diagnostic must be non-authorizing")
    }
    switch d.Status {
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateNotNeeded:
        if d.PlanStatus != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanNotNeeded ||
            d.Decision != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewNotNeeded ||
            d.ReviewSource != "" || d.ReviewEvidenceDigest != "" {
            return fmt.Errorf("not-needed review gate LSP diagnostic has inconsistent evidence")
        }
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateReady:
        if d.PlanStatus != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReady ||
            d.Decision != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewApprove ||
            d.PlanDigest == "" || d.ReviewSource == "" ||
            d.ReviewEvidenceDigest == "" || d.GateDigest == "" {
            return fmt.Errorf("ready review gate LSP diagnostic is incomplete")
        }
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateHold:
        if (d.PlanStatus != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReady &&
            d.PlanStatus != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanHold) ||
            d.Decision != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewHold ||
            d.PlanDigest == "" || d.ReviewSource == "" ||
            d.ReviewEvidenceDigest == "" || d.GateDigest == "" {
            return fmt.Errorf("held review gate LSP diagnostic is incomplete")
        }
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateRejected:
        if (d.PlanStatus != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReady &&
            d.PlanStatus != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanRejected) ||
            d.Decision != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewReject ||
            d.PlanDigest == "" || d.ReviewSource == "" ||
            d.ReviewEvidenceDigest == "" || d.GateDigest == "" {
            return fmt.Errorf("rejected review gate LSP diagnostic is incomplete")
        }
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateUnknown:
        if d.Publishable || d.MissingStage == "" {
            return fmt.Errorf("UNKNOWN review gate LSP diagnostic must remain non-publishable")
        }
    default:
        return fmt.Errorf("invalid JEV revision plan review gate LSP status")
    }
    return nil
}

func ProjectJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateLSP(input JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGate) JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateLSPDiagnostic {
    output := JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateLSPDiagnostic{
        Status:               input.Status,
        MissingStage:         input.MissingStage,
        PlanStatus:           input.PlanStatus,
        Decision:             input.Decision,
        PlanDigest:           input.PlanDigest,
        ReviewSource:         input.ReviewSource,
        ReviewEvidenceDigest: input.ReviewEvidenceDigest,
        GateDigest:           input.GateDigest,
        NonExecuting:         input.NonExecuting,
        NonAuthorizing:       input.NonAuthorizing,
    }
    if err := input.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateUnknown
        if output.MissingStage == "" {
            output.MissingStage = "revision-plan-review-evidence"
        }
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "JEV revision plan review gate is not publishable; evidence must be resolved"
        output.Publishable = false
        return output
    }
    output.Publishable = true
    switch input.Status {
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateNotNeeded:
        output.Severity = "info"
        output.Code = "jev.external-apply.revision-plan-review-not-needed"
        output.Message = "No revision plan review was needed; no execution or authorization occurred"
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateReady:
        output.Severity = "info"
        output.Code = "jev.external-apply.revision-plan-review-ready"
        output.Message = "Revision plan review is ready; no execution or authorization occurred"
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateHold:
        output.Severity = "info"
        output.Code = "jev.external-apply.revision-plan-review-hold"
        output.Message = "Revision plan review is held; no execution or authorization occurred"
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateRejected:
        output.Severity = "warning"
        output.Code = "jev.external-apply.revision-plan-review-rejected"
        output.Message = "Revision plan review was rejected; no execution or authorization occurred"
    default:
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateUnknown
        output.MissingStage = "revision-plan-review-status"
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "JEV revision plan review gate status is UNKNOWN; evidence must be resolved"
        output.Publishable = false
    }
    return output
}
