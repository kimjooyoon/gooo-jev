package gooo

import (
	"strings"
	"testing"
)

func TestCapabilityQueryOverviewBroadQuestion(t *testing.T) {
	overview := DiscoverCapabilityQueryOverview("What can gooo do?", "entity example")
	if err := overview.Validate(); err != nil {
		t.Fatal(err)
	}
	if overview.Mode != CapabilityQueryOverviewCatalog || len(overview.Trail.Response.Capabilities) == 0 {
		t.Fatalf("overview = %#v", overview)
	}
	if !strings.Contains(overview.Summary, "not a language completeness claim") {
		t.Fatalf("summary = %q", overview.Summary)
	}
}

func TestCapabilityQueryOverviewKeepsScopedQuestionMatched(t *testing.T) {
	overview := DiscoverCapabilityQueryOverview("How do I inspect reverse observation provenance?", "observe origin")
	if err := overview.Validate(); err != nil {
		t.Fatal(err)
	}
	if overview.Mode != CapabilityQueryOverviewMatchedOptions {
		t.Fatalf("overview = %#v", overview)
	}
}

func TestCapabilityQueryOverviewPreservesClarification(t *testing.T) {
	overview := DiscoverCapabilityQueryOverview("Can gooo model a database migration?", "")
	if err := overview.Validate(); err != nil {
		t.Fatal(err)
	}
	if overview.Mode != CapabilityQueryOverviewClarification || len(overview.SuggestedQuestions) == 0 {
		t.Fatalf("overview = %#v", overview)
	}
}

func TestCapabilityQueryOverviewRejectsTampering(t *testing.T) {
	overview := DiscoverCapabilityQueryOverview("What can gooo do?", "")
	overview.Summary += " tampered"
	if err := overview.Validate(); err == nil {
		t.Fatal("expected tampered overview to fail validation")
	}
}
