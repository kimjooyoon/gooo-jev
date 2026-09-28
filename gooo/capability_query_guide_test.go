package gooo

import (
	"strings"
	"testing"
)

func TestDiscoverCapabilityQueryGuideAvailable(t *testing.T) {
	guide := DiscoverCapabilityQueryGuide("What can gooo do?", "entity example")
	if err := guide.Validate(); err != nil {
		t.Fatal(err)
	}
	if guide.Status != CapabilityQueryAvailable || guide.Action != CapabilityQueryGuideInspectNextOperation {
		t.Fatalf("guide = %#v", guide)
	}
	if len(guide.CapabilityIDs) == 0 || len(guide.NextOperations) != len(guide.CapabilityIDs) {
		t.Fatalf("guide operations = %#v", guide)
	}
}

func TestDiscoverCapabilityQueryGuideDeferredPreservesBoundary(t *testing.T) {
	guide := DiscoverCapabilityQueryGuide("Can gooo execute and authorize a workload?", "workflow deploy")
	if err := guide.Validate(); err != nil {
		t.Fatal(err)
	}
	if guide.Status != CapabilityQueryDeferred || guide.Action != CapabilityQueryGuideRequireExternalBoundary {
		t.Fatalf("guide = %#v", guide)
	}
	found := false
	for _, question := range guide.Questions {
		if strings.Contains(question, "external boundary") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("guide lost external boundary question: %#v", guide.Questions)
	}
}

func TestDiscoverCapabilityQueryGuideUnknownPreservesClarification(t *testing.T) {
	guide := DiscoverCapabilityQueryGuide("Can gooo model a database migration?", "")
	if err := guide.Validate(); err != nil {
		t.Fatal(err)
	}
	if guide.Status != CapabilityQueryUnknown || guide.Action != CapabilityQueryGuideAskClarifyingQuestion || len(guide.CapabilityIDs) != 0 {
		t.Fatalf("guide = %#v", guide)
	}
}

func TestDiscoverCapabilityQueryGuideRejectsTampering(t *testing.T) {
	guide := DiscoverCapabilityQueryGuide("What can gooo do?", "")
	guide.EvidenceDigest = strings.Repeat("0", 64)
	if err := guide.Validate(); err == nil {
		t.Fatal("tampered capability query guide should fail validation")
	}
}
