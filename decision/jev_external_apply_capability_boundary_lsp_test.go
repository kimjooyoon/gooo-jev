package decision

import "testing"

func validExternalApplyCapabilityBoundaryForLSP() JEVExternalApplyCapabilityBoundary {
    scopeDigest, ok := digestJEVExternalApplyCapabilityScopes([]string{"review"})
    if !ok {
        panic("test scope digest must be valid")
    }
    boundary := JEVExternalApplyCapabilityBoundary{
        Status:                 jevExternalApplyCapabilityDescribed,
        BindingDigest:          "binding-evidence-digest",
        Principal:              "spiffe://example.org/ns/prod/sa/jev-reviewer",
        Audience:               "gooo-jev",
        Workspace:              "gooo-jev",
        NetworkAllowlistDigest: "network-allowlist-digest",
        ScopeDigest:             scopeDigest,
        NonExecuting:            true,
        NonAuthorizing:          true,
    }
    boundary.CapabilityDigest = digestJEVExternalApplyCapability(boundary.BindingDigest, boundary.Principal, boundary.Audience, boundary.Workspace, boundary.NetworkAllowlistDigest, boundary.ScopeDigest)
    return boundary
}

func TestProjectJEVExternalApplyCapabilityBoundaryLSPIsReviewOnly(t *testing.T) {
    got := ProjectJEVExternalApplyCapabilityBoundaryLSP(validExternalApplyCapabilityBoundaryForLSP())
    if got.Severity != "info" || got.Code != "jev.external-apply.capability-described" || !got.Publishable {
        t.Fatalf("got %+v", got)
    }
    if err := got.Validate(); err != nil {
        t.Fatal(err)
    }
}

func TestProjectJEVExternalApplyCapabilityBoundaryLSPUnknownIsNotPublishable(t *testing.T) {
    got := ProjectJEVExternalApplyCapabilityBoundaryLSP(JEVExternalApplyCapabilityBoundary{
        Status:         jevExternalApplyCapabilityUnknown,
        MissingStage:   "spiffe-principal",
        NonExecuting:   true,
        NonAuthorizing: true,
    })
    if got.Publishable || got.Severity != "error" || got.Code != "jev.provenance.unknown" {
        t.Fatalf("got %+v", got)
    }
}
