package gooo

import (
	"strings"
	"testing"
)

func TestCapabilityQueryHoverBindsDeclarationAndPreservesEvidence(t *testing.T) {
	hover := DiscoverCapabilityQueryHover("What can gooo do?", "entity user")
	if err := hover.Validate(); err != nil {
		t.Fatalf("hover should validate: %v", err)
	}
	if hover.Status != CapabilityQueryAvailable {
		t.Fatalf("hover status = %q, want %q", hover.Status, CapabilityQueryAvailable)
	}
	if hover.Declaration == nil || hover.SourceDigest != hover.Declaration.SourceDigest {
		t.Fatalf("hover is not bound to the declaration digest")
	}
	if len(hover.Capabilities) != 1 || hover.Capabilities[0].ID != "declaration_analysis" {
		t.Fatalf("hover capabilities = %#v, want declaration_analysis", hover.Capabilities)
	}
	if !hover.NonExecuting || !hover.NonAuthorizing || hover.EvidenceDigest == "" {
		t.Fatalf("hover lost safety or evidence binding")
	}
}

func TestCapabilityQueryHoverPreservesDeferredExternalBoundary(t *testing.T) {
	hover := DiscoverCapabilityQueryHover("What security boundary does gooo need?", "policy network")
	if err := hover.Validate(); err != nil {
		t.Fatalf("deferred hover should validate: %v", err)
	}
	if hover.Status != CapabilityQueryDeferred {
		t.Fatalf("hover status = %q, want %q", hover.Status, CapabilityQueryDeferred)
	}
	if hover.MissingStage != "explicit_external_boundary" || hover.FirstMismatch != "external_boundary" {
		t.Fatalf("hover boundary = %q/%q, want explicit external boundary", hover.MissingStage, hover.FirstMismatch)
	}
	if len(hover.Capabilities) != 1 || hover.Capabilities[0].ID != "security_boundary" {
		t.Fatalf("hover capabilities = %#v, want security_boundary", hover.Capabilities)
	}
}

func TestCapabilityQueryHoverRejectsTamperedEvidence(t *testing.T) {
	hover := DiscoverCapabilityQueryHover("Can gooo inspect provenance?", "observe user")
	if err := hover.Validate(); err != nil {
		t.Fatalf("hover should validate before tampering: %v", err)
	}
	hover.EvidenceDigest = strings.Repeat("0", 64)
	if err := hover.Validate(); err == nil {
		t.Fatal("tampered hover evidence should be rejected")
	}
}
