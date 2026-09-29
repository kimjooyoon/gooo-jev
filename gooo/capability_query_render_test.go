package gooo

import (
	"strings"
	"testing"
)

func TestRenderCapabilityQueryAvailable(t *testing.T) {
	trail := DiscoverCapabilityQueryTrail("What can gooo do?", "entity Example")
	guide := DiscoverCapabilityQueryGuide("What can gooo do?", "entity Example")
	rendered, err := RenderCapabilityQuery(trail, guide)
	if err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{
		"gooo capability discovery",
		"declaration_analysis",
		"next questions:",
		"non-executing: true",
		"non-authorizing: true",
	} {
		if !strings.Contains(rendered, expected) {
			t.Fatalf("rendered guidance omitted %q: %s", expected, rendered)
		}
	}
}

func TestRenderCapabilityQueryDeferredExplainsBoundary(t *testing.T) {
	query := "Can gooo execute and authorize a workload?"
	trail := DiscoverCapabilityQueryTrail(query, "")
	guide := DiscoverCapabilityQueryGuide(query, "")
	rendered, err := RenderCapabilityQuery(trail, guide)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(rendered, "external boundary") {
		t.Fatalf("deferred guidance omitted external boundary: %s", rendered)
	}
}

func TestRenderCapabilityQueryRejectsUnboundGuide(t *testing.T) {
	trail := DiscoverCapabilityQueryTrail("What can gooo do?", "")
	guide := DiscoverCapabilityQueryGuide("Can gooo execute?", "")
	if _, err := RenderCapabilityQuery(trail, guide); err == nil {
		t.Fatal("renderer should reject a guide from another query")
	}
}

func TestRenderCapabilityQueryRejectsMismatchedCapabilityGuide(t *testing.T) {
	trail := DiscoverCapabilityQueryTrail("Can gooo inspect provenance?", "")
	guide := DiscoverCapabilityQueryGuide("Can gooo generate code?", "")
	if _, err := RenderCapabilityQuery(trail, guide); err == nil {
		t.Fatal("renderer should reject a guide with different capability evidence")
	}
}
