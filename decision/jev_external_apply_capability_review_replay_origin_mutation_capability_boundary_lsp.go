package decision

import (
    "fmt"
    "strings"
)

type JEVExternalApplyCapabilityReviewReplayOriginMutationCapabilityBoundaryLSPDiagnostic struct {
    Severity                 string
    Code                    string
    Message                  string
    Status                   string
    MissingStage             string
    ReviewStatus             string
    Workspace                string
    PrincipalURI             string
    Audience                 string
    NetworkAllowlistDigest   string
    ReviewDecisionDigest     string
    CapabilityDigest         string
    Publishable              bool
    NonExecuting             bool
    NonAuthorizing           bool
}

func (d JEVExternalApplyCapabilityReviewReplayOriginMutationCapabilityBoundaryLSPDiagnostic) Validate() error {
    if d.Severity == "" || d.Code == "" || d.Message == "" || d.Status == "" {
        return fmt.Errorf("incomplete external apply capability review replay origin mutation capability boundary LSP diagnostic")
    }
    if d.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityNotNeeded &&
        d.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityRecorded &&
        d.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityHold &&
        d.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityRejected {
        return fmt.Errorf("invalid external apply capability review replay origin mutation capability LSP status")
    }
    if d.ReviewStatus != jevExternalApplyCapabilityReviewReplayOriginMutationReviewNotNeeded &&
        d.ReviewStatus != jevExternalApplyCapabilityReviewReplayOriginMutationReviewApproved &&
        d.ReviewStatus != jevExternalApplyCapabilityReviewReplayOriginMutationReviewHold &&
        d.ReviewStatus != jevExternalApplyCapabilityReviewReplayOriginMutationReviewRejected {
        return fmt.Errorf("invalid capability boundary LSP review status")
    }
    if d.ReviewStatus == jevExternalApplyCapabilityReviewReplayOriginMutationReviewApproved {
        if d.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityRecorded ||
            d.Workspace == "" ||
            !strings.HasPrefix(d.PrincipalURI, "spiffe://") ||
            d.Audience == "" ||
            d.NetworkAllowlistDigest == "" {
            return fmt.Errorf("recorded capability LSP diagnostic requires complete security boundary evidence")
        }
    }
    if d.ReviewStatus == jevExternalApplyCapabilityReviewReplayOriginMutationReviewNotNeeded &&
        d.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityNotNeeded {
        return fmt.Errorf("not-needed review LSP diagnostic requires capability-not-needed status")
    }
    if d.ReviewStatus == jevExternalApplyCapabilityReviewReplayOriginMutationReviewHold &&
        d.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityHold {
        return fmt.Errorf("held review LSP diagnostic requires capability-hold status")
    }
    if d.ReviewStatus == jevExternalApplyCapabilityReviewReplayOriginMutationReviewRejected &&
        d.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityRejected {
        return fmt.Errorf("rejected review LSP diagnostic requires capability-rejected status")
    }
    if d.ReviewDecisionDigest == "" || d.CapabilityDigest == "" {
        return fmt.Errorf("capability boundary LSP diagnostic requires digest evidence")
    }
    if !d.NonExecuting {
        return fmt.Errorf("external apply capability review replay origin mutation capability LSP diagnostic must be non-executing")
    }
    if !d.NonAuthorizing {
        return fmt.Errorf("external apply capability review replay origin mutation capability LSP diagnostic must be non-authorizing")
    }
    return nil
}

func ProjectJEVExternalApplyCapabilityReviewReplayOriginMutationCapabilityBoundaryLSP(input JEVExternalApplyCapabilityReviewReplayOriginMutationCapabilityBoundary) JEVExternalApplyCapabilityReviewReplayOriginMutationCapabilityBoundaryLSPDiagnostic {
    output := JEVExternalApplyCapabilityReviewReplayOriginMutationCapabilityBoundaryLSPDiagnostic{
        Status:                 input.Status,
        MissingStage:           input.MissingStage,
        ReviewStatus:            input.ReviewStatus,
        Workspace:               input.Workspace,
        PrincipalURI:            input.PrincipalURI,
        Audience:               input.Audience,
        NetworkAllowlistDigest: input.NetworkAllowlistDigest,
        ReviewDecisionDigest:    input.ReviewDecisionDigest,
        CapabilityDigest:        input.CapabilityDigest,
        NonExecuting:            input.NonExecuting,
        NonAuthorizing:          input.NonAuthorizing,
    }
    if err := input.Validate(); err != nil {
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "JEV capability boundary is not publishable: security evidence is incomplete or invalid"
        output.Publishable = false
        return output
    }
    output.Publishable = true
    switch input.Status {
    case jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityNotNeeded:
        output.Severity = "info"
        output.Code = "jev.external-apply.mutation-capability-not-needed"
        output.Message = "No mutation capability was needed; no authorization or execution occurred"
    case jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityRecorded:
        output.Severity = "info"
        output.Code = "jev.external-apply.mutation-capability-recorded"
        output.Message = "Mutation capability boundary was recorded for controlled review; no authorization or execution occurred"
    case jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityHold:
        output.Severity = "info"
        output.Code = "jev.external-apply.mutation-capability-hold"
        output.Message = "Mutation capability remains on hold; no authorization or execution occurred"
    case jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityRejected:
        output.Severity = "warning"
        output.Code = "jev.external-apply.mutation-capability-rejected"
        output.Message = "Mutation capability was rejected; no authorization or execution occurred"
    default:
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "JEV capability boundary is UNKNOWN; security evidence must be resolved"
        output.Publishable = false
    }
    return output
}
