package decision

import "fmt"

type JEVExternalApplyCapabilityReviewReplayOriginMutationReviewGateLSPDiagnostic struct {
    Severity             string
    Code                 string
    Message              string
    Status               string
    MissingStage         string
    ProposalStatus       string
    ReviewDecision       string
    ProposalDigest       string
    ReviewEvidenceDigest string
    DecisionDigest       string
    Publishable          bool
    NonExecuting         bool
    NonAuthorizing       bool
}

func (d JEVExternalApplyCapabilityReviewReplayOriginMutationReviewGateLSPDiagnostic) Validate() error {
    if d.Severity == "" || d.Code == "" || d.Message == "" || d.Status == "" {
        return fmt.Errorf("incomplete external apply capability review replay origin mutation review gate LSP diagnostic")
    }
    if d.Status != jevExternalApplyCapabilityReviewReplayOriginMutationReviewNotNeeded &&
        d.Status != jevExternalApplyCapabilityReviewReplayOriginMutationReviewApproved &&
        d.Status != jevExternalApplyCapabilityReviewReplayOriginMutationReviewHold &&
        d.Status != jevExternalApplyCapabilityReviewReplayOriginMutationReviewRejected {
        return fmt.Errorf("invalid external apply capability review replay origin mutation review gate LSP status")
    }
    if d.ProposalStatus != jevExternalApplyCapabilityReviewReplayOriginMutationNotNeeded &&
        d.ProposalStatus != jevExternalApplyCapabilityReviewReplayOriginMutationProposed {
        return fmt.Errorf("invalid external apply capability review replay origin mutation proposal LSP status")
    }
    if d.ProposalStatus == jevExternalApplyCapabilityReviewReplayOriginMutationNotNeeded {
        if d.Status != jevExternalApplyCapabilityReviewReplayOriginMutationReviewNotNeeded ||
            d.ReviewDecision != "" ||
            d.ReviewEvidenceDigest != "" {
            return fmt.Errorf("mutation-not-needed LSP diagnostic requires no review decision")
        }
    }
    if d.ProposalStatus == jevExternalApplyCapabilityReviewReplayOriginMutationProposed {
        if d.ReviewDecision != jevExternalApplyCapabilityReviewReplayOriginMutationReviewDecisionApprove &&
            d.ReviewDecision != jevExternalApplyCapabilityReviewReplayOriginMutationReviewDecisionHold &&
            d.ReviewDecision != jevExternalApplyCapabilityReviewReplayOriginMutationReviewDecisionReject {
            return fmt.Errorf("invalid review decision in mutation review gate LSP diagnostic")
        }
        if d.ReviewEvidenceDigest == "" {
            return fmt.Errorf("review decision LSP diagnostic requires review evidence")
        }
        if d.ReviewDecision == jevExternalApplyCapabilityReviewReplayOriginMutationReviewDecisionApprove &&
            d.Status != jevExternalApplyCapabilityReviewReplayOriginMutationReviewApproved {
            return fmt.Errorf("approve LSP diagnostic requires approved status")
        }
        if d.ReviewDecision == jevExternalApplyCapabilityReviewReplayOriginMutationReviewDecisionHold &&
            d.Status != jevExternalApplyCapabilityReviewReplayOriginMutationReviewHold {
            return fmt.Errorf("hold LSP diagnostic requires hold status")
        }
        if d.ReviewDecision == jevExternalApplyCapabilityReviewReplayOriginMutationReviewDecisionReject &&
            d.Status != jevExternalApplyCapabilityReviewReplayOriginMutationReviewRejected {
            return fmt.Errorf("reject LSP diagnostic requires rejected status")
        }
    }
    if d.ProposalDigest == "" || d.DecisionDigest == "" {
        return fmt.Errorf("mutation review gate LSP diagnostic requires complete digest evidence")
    }
    if !d.NonExecuting {
        return fmt.Errorf("external apply capability review replay origin mutation review gate LSP diagnostic must be non-executing")
    }
    if !d.NonAuthorizing {
        return fmt.Errorf("external apply capability review replay origin mutation review gate LSP diagnostic must be non-authorizing")
    }
    return nil
}

func ProjectJEVExternalApplyCapabilityReviewReplayOriginMutationReviewGateLSP(input JEVExternalApplyCapabilityReviewReplayOriginMutationReviewGate) JEVExternalApplyCapabilityReviewReplayOriginMutationReviewGateLSPDiagnostic {
    output := JEVExternalApplyCapabilityReviewReplayOriginMutationReviewGateLSPDiagnostic{
        Status:               input.Status,
        MissingStage:         input.MissingStage,
        ProposalStatus:       input.ProposalStatus,
        ReviewDecision:       input.ReviewDecision,
        ProposalDigest:       input.ProposalDigest,
        ReviewEvidenceDigest: input.ReviewEvidenceDigest,
        DecisionDigest:       input.DecisionDigest,
        NonExecuting:         input.NonExecuting,
        NonAuthorizing:       input.NonAuthorizing,
    }
    if err := input.Validate(); err != nil {
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "JEV mutation review gate is not publishable: review evidence is incomplete or invalid"
        output.Publishable = false
        return output
    }
    output.Publishable = true
    switch input.Status {
    case jevExternalApplyCapabilityReviewReplayOriginMutationReviewNotNeeded:
        output.Severity = "info"
        output.Code = "jev.external-apply.origin-mutation-not-needed"
        output.Message = "No origin mutation was needed; no automatic change or authorization occurred"
    case jevExternalApplyCapabilityReviewReplayOriginMutationReviewApproved:
        output.Severity = "info"
        output.Code = "jev.external-apply.origin-mutation-review-approved"
        output.Message = "Origin mutation proposal was approved for the next controlled step; no automatic change occurred"
    case jevExternalApplyCapabilityReviewReplayOriginMutationReviewHold:
        output.Severity = "info"
        output.Code = "jev.external-apply.origin-mutation-review-hold"
        output.Message = "Origin mutation proposal is held pending more evidence; no automatic change occurred"
    case jevExternalApplyCapabilityReviewReplayOriginMutationReviewRejected:
        output.Severity = "warning"
        output.Code = "jev.external-apply.origin-mutation-review-rejected"
        output.Message = "Origin mutation proposal was rejected; no automatic change occurred"
    default:
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "JEV mutation review gate is UNKNOWN; evidence must be resolved before the next step"
        output.Publishable = false
    }
    return output
}
