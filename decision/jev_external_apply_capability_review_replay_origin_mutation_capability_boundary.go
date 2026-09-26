package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "strings"
)

const jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityNotNeeded = "mutation-capability-not-needed"
const jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityRecorded = "mutation-capability-recorded"
const jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityHold = "mutation-capability-hold"
const jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityRejected = "mutation-capability-rejected"
const jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityUnknown = "UNKNOWN"

type JEVExternalApplyCapabilityReviewReplayOriginMutationCapabilityBoundaryInput struct {
    ReviewGate                JEVExternalApplyCapabilityReviewReplayOriginMutationReviewGate
    Workspace                 string
    PrincipalURI              string
    Audience                 string
    NetworkAllowlistDigest    string
    NonAuthorizing            bool
}

type JEVExternalApplyCapabilityReviewReplayOriginMutationCapabilityBoundary struct {
    Status                    string
    MissingStage              string
    ReviewStatus              string
    Workspace                 string
    PrincipalURI              string
    Audience                 string
    NetworkAllowlistDigest    string
    ReviewDecisionDigest      string
    CapabilityDigest          string
    NonExecuting              bool
    NonAuthorizing            bool
}

func (b JEVExternalApplyCapabilityReviewReplayOriginMutationCapabilityBoundary) Validate() error {
    if b.Status == "" || b.ReviewStatus == "" || b.ReviewDecisionDigest == "" || b.CapabilityDigest == "" {
        return fmt.Errorf("incomplete JEV external apply capability review replay origin mutation capability boundary")
    }
    if b.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityNotNeeded &&
        b.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityRecorded &&
        b.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityHold &&
        b.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityRejected {
        return fmt.Errorf("invalid JEV external apply capability review replay origin mutation capability status")
    }
    if b.ReviewStatus != jevExternalApplyCapabilityReviewReplayOriginMutationReviewNotNeeded &&
        b.ReviewStatus != jevExternalApplyCapabilityReviewReplayOriginMutationReviewApproved &&
        b.ReviewStatus != jevExternalApplyCapabilityReviewReplayOriginMutationReviewHold &&
        b.ReviewStatus != jevExternalApplyCapabilityReviewReplayOriginMutationReviewRejected {
        return fmt.Errorf("invalid JEV external apply capability review status")
    }
    if b.ReviewStatus == jevExternalApplyCapabilityReviewReplayOriginMutationReviewNotNeeded &&
        b.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityNotNeeded {
        return fmt.Errorf("mutation-not-needed review requires capability-not-needed status")
    }
    if b.ReviewStatus == jevExternalApplyCapabilityReviewReplayOriginMutationReviewApproved {
        if b.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityRecorded {
            return fmt.Errorf("approved review requires recorded capability status")
        }
        if b.Workspace == "" || !strings.HasPrefix(b.PrincipalURI, "spiffe://") || b.Audience == "" || b.NetworkAllowlistDigest == "" {
            return fmt.Errorf("recorded capability requires workspace, SPIFFE principal, audience, and network allowlist digest")
        }
    }
    if b.ReviewStatus == jevExternalApplyCapabilityReviewReplayOriginMutationReviewHold &&
        b.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityHold {
        return fmt.Errorf("held review requires capability-hold status")
    }
    if b.ReviewStatus == jevExternalApplyCapabilityReviewReplayOriginMutationReviewRejected &&
        b.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityRejected {
        return fmt.Errorf("rejected review requires capability-rejected status")
    }
    if !b.NonExecuting {
        return fmt.Errorf("JEV external apply capability review replay origin mutation capability boundary must be non-executing")
    }
    if !b.NonAuthorizing {
        return fmt.Errorf("JEV external apply capability review replay origin mutation capability boundary must be non-authorizing")
    }
    expected := digestJEVExternalApplyCapabilityReviewReplayOriginMutationCapabilityBoundary(
        b.Status,
        b.ReviewStatus,
        b.Workspace,
        b.PrincipalURI,
        b.Audience,
        b.NetworkAllowlistDigest,
        b.ReviewDecisionDigest,
    )
    if b.CapabilityDigest != expected {
        return fmt.Errorf("JEV external apply capability review replay origin mutation capability digest mismatch")
    }
    return nil
}

func RecordJEVExternalApplyCapabilityReviewReplayOriginMutationCapability(input JEVExternalApplyCapabilityReviewReplayOriginMutationCapabilityBoundaryInput) JEVExternalApplyCapabilityReviewReplayOriginMutationCapabilityBoundary {
    output := JEVExternalApplyCapabilityReviewReplayOriginMutationCapabilityBoundary{
        Status:         jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityUnknown,
        NonExecuting:   true,
        NonAuthorizing: true,
    }
    if !input.NonAuthorizing || !input.ReviewGate.NonAuthorizing {
        output.MissingStage = "authorization-boundary"
        return output
    }
    if !input.ReviewGate.NonExecuting {
        output.MissingStage = "execution-boundary"
        return output
    }
    if err := input.ReviewGate.Validate(); err != nil {
        output.MissingStage = "mutation-review-gate"
        return output
    }
    output.ReviewStatus = input.ReviewGate.Status
    output.ReviewDecisionDigest = input.ReviewGate.DecisionDigest
    switch input.ReviewGate.Status {
    case jevExternalApplyCapabilityReviewReplayOriginMutationReviewNotNeeded:
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityNotNeeded
    case jevExternalApplyCapabilityReviewReplayOriginMutationReviewApproved:
        if input.Workspace == "" {
            output.MissingStage = "workspace-boundary"
            return output
        }
        if !strings.HasPrefix(input.PrincipalURI, "spiffe://") {
            output.MissingStage = "spiffe-principal"
            return output
        }
        if input.Audience == "" {
            output.MissingStage = "capability-audience"
            return output
        }
        if input.NetworkAllowlistDigest == "" {
            output.MissingStage = "network-allowlist"
            return output
        }
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityRecorded
        output.Workspace = input.Workspace
        output.PrincipalURI = input.PrincipalURI
        output.Audience = input.Audience
        output.NetworkAllowlistDigest = input.NetworkAllowlistDigest
    case jevExternalApplyCapabilityReviewReplayOriginMutationReviewHold:
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityHold
    case jevExternalApplyCapabilityReviewReplayOriginMutationReviewRejected:
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityRejected
    default:
        output.MissingStage = "mutation-review-status"
        return output
    }
    output.CapabilityDigest = digestJEVExternalApplyCapabilityReviewReplayOriginMutationCapabilityBoundary(
        output.Status,
        output.ReviewStatus,
        output.Workspace,
        output.PrincipalURI,
        output.Audience,
        output.NetworkAllowlistDigest,
        output.ReviewDecisionDigest,
    )
    if err := output.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityUnknown
        output.MissingStage = "mutation-capability-boundary-evidence"
        output.CapabilityDigest = ""
    }
    return output
}

func digestJEVExternalApplyCapabilityReviewReplayOriginMutationCapabilityBoundary(status, reviewStatus, workspace, principalURI, audience, networkAllowlistDigest, reviewDecisionDigest string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s", status, reviewStatus, workspace, principalURI, audience, networkAllowlistDigest, reviewDecisionDigest)))
    return hex.EncodeToString(sum[:])
}
