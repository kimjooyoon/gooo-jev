package decision

import "testing"

func validExternalApplyOriginBindingForCapability() JEVExternalApplyOriginBinding {
    binding := JEVExternalApplyOriginBinding{
        Status:                jevExternalApplyOriginBindingReady,
        OriginDigest:          "origin-digest",
        CandidateDigest:       "candidate-digest",
        RequestEvidenceDigest: "request-evidence-digest",
        NonExecuting:          true,
        NonAuthorizing:        true,
    }
    binding.BindingEvidenceDigest = digestJEVExternalApplyOriginBinding(binding.OriginDigest, binding.CandidateDigest, binding.RequestEvidenceDigest)
    binding.EvidenceDigest = digestJEVExternalApplyOriginBindingEvidence(binding.Status, binding.OriginDigest, binding.BindingEvidenceDigest)
    return binding
}

func TestDescribeJEVExternalApplyCapabilityBoundaryIsReviewOnly(t *testing.T) {
    got := DescribeJEVExternalApplyCapabilityBoundary(JEVExternalApplyCapabilityBoundaryInput{
        Binding:                validExternalApplyOriginBindingForCapability(),
        Principal:              "spiffe://example.org/ns/prod/sa/jev-reviewer",
        Audience:               "gooo-jev",
        Workspace:              "gooo-jev",
        NetworkAllowlistDigest: "network-allowlist-digest",
        Scopes:                 []string{"provenance.read", "review"},
        NonAuthorizing:         true,
    })
    if got.Status != jevExternalApplyCapabilityDescribed || !got.NonExecuting || !got.NonAuthorizing {
        t.Fatalf("got %+v", got)
    }
    if err := got.Validate(); err != nil {
        t.Fatal(err)
    }
}

func TestDescribeJEVExternalApplyCapabilityBoundaryRejectsUnboundPrincipal(t *testing.T) {
    got := DescribeJEVExternalApplyCapabilityBoundary(JEVExternalApplyCapabilityBoundaryInput{
        Binding:                validExternalApplyOriginBindingForCapability(),
        Principal:              "service-account:jev-reviewer",
        Audience:               "gooo-jev",
        Workspace:              "gooo-jev",
        NetworkAllowlistDigest: "network-allowlist-digest",
        Scopes:                 []string{"review"},
        NonAuthorizing:         true,
    })
    if got.Status != jevExternalApplyCapabilityUnknown || got.MissingStage != "spiffe-principal" {
        t.Fatalf("got %+v", got)
    }
}

func TestDescribeJEVExternalApplyCapabilityBoundaryRejectsDuplicateScopes(t *testing.T) {
    got := DescribeJEVExternalApplyCapabilityBoundary(JEVExternalApplyCapabilityBoundaryInput{
        Binding:                validExternalApplyOriginBindingForCapability(),
        Principal:              "spiffe://example.org/ns/prod/sa/jev-reviewer",
        Audience:               "gooo-jev",
        Workspace:              "gooo-jev",
        NetworkAllowlistDigest: "network-allowlist-digest",
        Scopes:                 []string{"review", "review"},
        NonAuthorizing:         true,
    })
    if got.Status != jevExternalApplyCapabilityUnknown || got.MissingStage != "capability-scopes" {
        t.Fatalf("got %+v", got)
    }
}
