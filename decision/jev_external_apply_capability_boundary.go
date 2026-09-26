package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "sort"
    "strings"
)

const jevExternalApplyCapabilityDescribed = "capability-described"
const jevExternalApplyCapabilityUnknown = "UNKNOWN"

type JEVExternalApplyCapabilityBoundaryInput struct {
    Binding                 JEVExternalApplyOriginBinding
    Principal               string
    Audience                string
    Workspace               string
    NetworkAllowlistDigest  string
    Scopes                  []string
    NonAuthorizing          bool
}

type JEVExternalApplyCapabilityBoundary struct {
    Status                   string
    MissingStage             string
    BindingDigest            string
    Principal                string
    Audience                 string
    Workspace                string
    NetworkAllowlistDigest   string
    ScopeDigest              string
    CapabilityDigest         string
    NonExecuting             bool
    NonAuthorizing           bool
}

func (c JEVExternalApplyCapabilityBoundary) Validate() error {
    if c.Status == "" || c.BindingDigest == "" || c.Principal == "" || c.Audience == "" || c.Workspace == "" || c.NetworkAllowlistDigest == "" || c.ScopeDigest == "" || c.CapabilityDigest == "" {
        return fmt.Errorf("incomplete JEV external apply capability boundary")
    }
    if c.Status != jevExternalApplyCapabilityDescribed {
        return fmt.Errorf("invalid JEV external apply capability boundary status")
    }
    if !strings.HasPrefix(c.Principal, "spiffe://") {
        return fmt.Errorf("JEV external apply capability principal must use SPIFFE URI form")
    }
    if !c.NonExecuting {
        return fmt.Errorf("JEV external apply capability boundary must be non-executing")
    }
    if !c.NonAuthorizing {
        return fmt.Errorf("JEV external apply capability boundary must be non-authorizing")
    }
    expected := digestJEVExternalApplyCapability(c.BindingDigest, c.Principal, c.Audience, c.Workspace, c.NetworkAllowlistDigest, c.ScopeDigest)
    if c.CapabilityDigest != expected {
        return fmt.Errorf("JEV external apply capability digest mismatch")
    }
    return nil
}

func DescribeJEVExternalApplyCapabilityBoundary(input JEVExternalApplyCapabilityBoundaryInput) JEVExternalApplyCapabilityBoundary {
    output := JEVExternalApplyCapabilityBoundary{
        Status:         jevExternalApplyCapabilityUnknown,
        NonExecuting:   true,
        NonAuthorizing: true,
    }
    if !input.NonAuthorizing || !input.Binding.NonAuthorizing {
        output.MissingStage = "authorization-boundary"
        return output
    }
    if !input.Binding.NonExecuting {
        output.MissingStage = "execution-boundary"
        return output
    }
    if err := input.Binding.Validate(); err != nil {
        output.MissingStage = "external-apply-origin-binding"
        return output
    }
    if !strings.HasPrefix(input.Principal, "spiffe://") {
        output.MissingStage = "spiffe-principal"
        return output
    }
    if input.Audience == "" {
        output.MissingStage = "audience"
        return output
    }
    if input.Workspace == "" {
        output.MissingStage = "workspace"
        return output
    }
    if input.NetworkAllowlistDigest == "" {
        output.MissingStage = "network-allowlist"
        return output
    }
    scopeDigest, ok := digestJEVExternalApplyCapabilityScopes(input.Scopes)
    if !ok {
        output.MissingStage = "capability-scopes"
        return output
    }
    output.Status = jevExternalApplyCapabilityDescribed
    output.BindingDigest = input.Binding.EvidenceDigest
    output.Principal = input.Principal
    output.Audience = input.Audience
    output.Workspace = input.Workspace
    output.NetworkAllowlistDigest = input.NetworkAllowlistDigest
    output.ScopeDigest = scopeDigest
    output.CapabilityDigest = digestJEVExternalApplyCapability(output.BindingDigest, output.Principal, output.Audience, output.Workspace, output.NetworkAllowlistDigest, output.ScopeDigest)
    if err := output.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityUnknown
        output.MissingStage = "capability-boundary-evidence"
        output.BindingDigest = ""
        output.CapabilityDigest = ""
    }
    return output
}

func digestJEVExternalApplyCapabilityScopes(scopes []string) (string, bool) {
    if len(scopes) == 0 {
        return "", false
    }
    normalized := append([]string(nil), scopes...)
    sort.Strings(normalized)
    for index, scope := range normalized {
        if scope == "" || (index > 0 && normalized[index-1] == scope) {
            return "", false
        }
    }
    sum := sha256.Sum256([]byte(strings.Join(normalized, "|")))
    return hex.EncodeToString(sum[:]), true
}

func digestJEVExternalApplyCapability(bindingDigest, principal, audience, workspace, networkAllowlistDigest, scopeDigest string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s|%s|%s", bindingDigest, principal, audience, workspace, networkAllowlistDigest, scopeDigest)))
    return hex.EncodeToString(sum[:])
}
