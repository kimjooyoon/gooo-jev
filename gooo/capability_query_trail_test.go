package gooo

import "testing"

func TestDiscoverCapabilityQueryTrailBindsNaturalQuestionToDeclaration(t *testing.T) {
	trail := DiscoverCapabilityQueryTrail(
		"What can gooo do with this declaration?",
		"entity Invoice\noperation reconcile\nobserve receipt\n",
	)
	if err := trail.Validate(); err != nil {
		t.Fatalf("validate trail: %v", err)
	}
	if trail.Intent != "overview" {
		t.Fatalf("intent = %q, want overview", trail.Intent)
	}
	if trail.Response.Declaration == nil || !trail.Response.Declaration.Bound {
		t.Fatalf("declaration binding missing: %#v", trail.Response.Declaration)
	}
	if len(trail.NextQuestions) == 0 || len(trail.DiscoveryPath) < 4 {
		t.Fatalf("incomplete discovery trail: %#v", trail)
	}
}

func TestDiscoverCapabilityQueryTrailPreservesUnknownAndBoundary(t *testing.T) {
	unknown := DiscoverCapabilityQueryTrail("Can gooo predict the stock market tomorrow?", "")
	if err := unknown.Validate(); err != nil {
		t.Fatalf("validate unknown trail: %v", err)
	}
	if unknown.Response.Status != CapabilityQueryUnknown || unknown.Intent != "unknown" {
		t.Fatalf("unknown result = %#v", unknown)
	}

	boundary := DiscoverCapabilityQueryTrail("Can gooo execute and authorize this workload?", "workflow deploy")
	if err := boundary.Validate(); err != nil {
		t.Fatalf("validate boundary trail: %v", err)
	}
	if boundary.Response.Status != CapabilityQueryDeferred || boundary.Intent != "boundary_request" {
		t.Fatalf("boundary result = %#v", boundary)
	}
	if boundary.Response.NonExecuting != true || boundary.Response.NonAuthorizing != true {
		t.Fatalf("boundary crossed safety flags: %#v", boundary.Response)
	}
}
