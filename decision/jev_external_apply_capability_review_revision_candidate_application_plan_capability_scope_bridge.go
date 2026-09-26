package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "strings"
)

const jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridgeBound = "application-plan-capability-scope-bound"
const jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridgeUnknown = "UNKNOWN"
const jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBound = "capability-scope-bound"
const jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeHeld = "capability-scope-held"
const jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeRejected = "capability-scope-rejected"

type JEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridgeInput struct {
    ApplicationPlan        JEVExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridge
    SpiffeID                string
    Audience                string
    SandboxPolicyDigest     string
    NetworkAllowlistDigest  string
    ScopeEvidenceDigest     string
    NonAuthorizing          bool
}

type JEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridge struct {
    Status                      string
    MissingStage                string
    ApplicationPlanStatus       string
    PlanDigest                  string
    PlanSource                  string
    PlanEvidenceDigest          string
    ApplicationPlanBridgeDigest string
    SpiffeID                    string
    Audience                    string
    SandboxPolicyDigest         string
    NetworkAllowlistDigest      string
    ScopeEvidenceDigest         string
    ScopeStatus                 string
    BridgeDigest                string
    NonExecuting                bool
    NonAuthorizing               bool
}

func (b JEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridge) Validate() error {
    if b.Status == "" || b.ApplicationPlanStatus == "" ||
        b.ApplicationPlanBridgeDigest == "" || b.ScopeStatus == "" ||
        b.BridgeDigest == "" {
        return fmt.Errorf("incomplete JEV application plan capability scope bridge")
    }
    if b.Status != jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridgeBound {
        return fmt.Errorf("invalid JEV application plan capability scope bridge status")
    }
    switch b.ApplicationPlanStatus {
    case jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanReady:
        if b.PlanDigest == "" || b.PlanSource == "" || b.PlanEvidenceDigest == "" ||
            !strings.HasPrefix(b.SpiffeID, "spiffe://") || b.Audience == "" ||
            b.SandboxPolicyDigest == "" || b.NetworkAllowlistDigest == "" ||
            b.ScopeEvidenceDigest == "" ||
            b.ScopeStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBound {
            return fmt.Errorf("ready application plan capability scope bridge is incomplete")
        }
    case jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanHeld:
        if b.ScopeStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeHeld ||
            b.PlanDigest != "" || b.PlanSource != "" || b.PlanEvidenceDigest != "" ||
            b.SpiffeID != "" || b.Audience != "" || b.SandboxPolicyDigest != "" ||
            b.NetworkAllowlistDigest != "" || b.ScopeEvidenceDigest != "" {
            return fmt.Errorf("held application plan capability scope bridge has inconsistent evidence")
        }
    case jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanRejected:
        if b.ScopeStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeRejected ||
            b.PlanDigest != "" || b.PlanSource != "" || b.PlanEvidenceDigest != "" ||
            b.SpiffeID != "" || b.Audience != "" || b.SandboxPolicyDigest != "" ||
            b.NetworkAllowlistDigest != "" || b.ScopeEvidenceDigest != "" {
            return fmt.Errorf("rejected application plan capability scope bridge has inconsistent evidence")
        }
    default:
        return fmt.Errorf("invalid application plan status")
    }
    if !b.NonExecuting {
        return fmt.Errorf("JEV application plan capability scope bridge must be non-executing")
    }
    if !b.NonAuthorizing {
        return fmt.Errorf("JEV application plan capability scope bridge must be non-authorizing")
    }
    expected := digestJEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridge(
        b.Status,
        b.ApplicationPlanStatus,
        b.PlanDigest,
        b.PlanSource,
        b.PlanEvidenceDigest,
        b.ApplicationPlanBridgeDigest,
        b.SpiffeID,
        b.Audience,
        b.SandboxPolicyDigest,
        b.NetworkAllowlistDigest,
        b.ScopeEvidenceDigest,
        b.ScopeStatus,
    )
    if b.BridgeDigest != expected {
        return fmt.Errorf("JEV application plan capability scope bridge digest mismatch")
    }
    return nil
}

func BindJEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridge(input JEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridgeInput) JEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridge {
    output := JEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridge{
        Status:         jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridgeUnknown,
        NonExecuting:   true,
        NonAuthorizing: true,
    }
    if !input.NonAuthorizing || !input.ApplicationPlan.NonAuthorizing {
        output.MissingStage = "authorization-boundary"
        return output
    }
    if !input.ApplicationPlan.NonExecuting {
        output.MissingStage = "execution-boundary"
        return output
    }
    if err := input.ApplicationPlan.Validate(); err != nil {
        output.MissingStage = "application-plan"
        return output
    }
    output.Status = jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridgeBound
    output.ApplicationPlanStatus = input.ApplicationPlan.ApplicationPlanStatus
    output.PlanDigest = input.ApplicationPlan.PlanDigest
    output.PlanSource = input.ApplicationPlan.PlanSource
    output.PlanEvidenceDigest = input.ApplicationPlan.PlanEvidenceDigest
    output.ApplicationPlanBridgeDigest = input.ApplicationPlan.BridgeDigest
    switch input.ApplicationPlan.ApplicationPlanStatus {
    case jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanReady:
        if input.SpiffeID == "" {
            output.MissingStage = "spiffe-id"
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridgeUnknown
            return output
        }
        if !strings.HasPrefix(input.SpiffeID, "spiffe://") {
            output.MissingStage = "spiffe-id-format"
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridgeUnknown
            return output
        }
        if input.Audience == "" {
            output.MissingStage = "capability-audience"
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridgeUnknown
            return output
        }
        if input.SandboxPolicyDigest == "" {
            output.MissingStage = "sandbox-policy"
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridgeUnknown
            return output
        }
        if input.NetworkAllowlistDigest == "" {
            output.MissingStage = "network-allowlist"
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridgeUnknown
            return output
        }
        if input.ScopeEvidenceDigest == "" {
            output.MissingStage = "capability-scope-evidence"
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridgeUnknown
            return output
        }
        output.SpiffeID = input.SpiffeID
        output.Audience = input.Audience
        output.SandboxPolicyDigest = input.SandboxPolicyDigest
        output.NetworkAllowlistDigest = input.NetworkAllowlistDigest
        output.ScopeEvidenceDigest = input.ScopeEvidenceDigest
        output.ScopeStatus = jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBound
    case jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanHeld:
        output.ScopeStatus = jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeHeld
    case jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanRejected:
        output.ScopeStatus = jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeRejected
    default:
        output.MissingStage = "application-plan-status"
        output.Status = jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridgeUnknown
        return output
    }
    output.BridgeDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridge(
        output.Status,
        output.ApplicationPlanStatus,
        output.PlanDigest,
        output.PlanSource,
        output.PlanEvidenceDigest,
        output.ApplicationPlanBridgeDigest,
        output.SpiffeID,
        output.Audience,
        output.SandboxPolicyDigest,
        output.NetworkAllowlistDigest,
        output.ScopeEvidenceDigest,
        output.ScopeStatus,
    )
    if err := output.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridgeUnknown
        output.MissingStage = "application-plan-capability-scope-evidence"
        output.BridgeDigest = ""
    }
    return output
}

func digestJEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridge(status, applicationPlanStatus, planDigest, planSource, planEvidenceDigest, applicationPlanBridgeDigest, spiffeID, audience, sandboxPolicyDigest, networkAllowlistDigest, scopeEvidenceDigest, scopeStatus string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s", status, applicationPlanStatus, planDigest, planSource, planEvidenceDigest, applicationPlanBridgeDigest, spiffeID, audience, sandboxPolicyDigest, networkAllowlistDigest, scopeEvidenceDigest, scopeStatus)))
    return hex.EncodeToString(sum[:])
}
