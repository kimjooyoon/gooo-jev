package gooo

import "testing"

func TestNewJEVObservationEnvelopeFromAvailableHover(t *testing.T) {
	hover := DiscoverCapabilityQueryHover("What can gooo do?", "entity user")
	envelope, err := NewJEVObservationEnvelopeFromCapabilityQueryHover(hover, "corr-available", "catalog-1", "schema-1", 9)
	if err != nil {
		t.Fatalf("project available hover: %v", err)
	}
	if err := envelope.Validate(); err != nil {
		t.Fatalf("validate available envelope: %v", err)
	}
	if envelope.FirstMissingStage != -1 || envelope.NextOperation != "inspect_next_operation" {
		t.Fatalf("available boundary = %d/%q, want -1/inspect_next_operation", envelope.FirstMissingStage, envelope.NextOperation)
	}
	if len(envelope.Capabilities) != 1 || envelope.Capabilities[0].State != JEVObservationAvailable {
		t.Fatalf("available capabilities = %#v", envelope.Capabilities)
	}
	if envelope.Executed || envelope.Authorizing || envelope.EvidenceDigest != hover.EvidenceDigest {
		t.Fatalf("available envelope lost safety or hover evidence")
	}
}

func TestNewJEVObservationEnvelopeFromDeferredHover(t *testing.T) {
	hover := DiscoverCapabilityQueryHover("What security boundary does gooo need?", "policy network")
	envelope, err := NewJEVObservationEnvelopeFromCapabilityQueryHover(hover, "corr-deferred", "catalog-2", "schema-2", 3)
	if err != nil {
		t.Fatalf("project deferred hover: %v", err)
	}
	if envelope.Capabilities[0].State != JEVObservationDeferred {
		t.Fatalf("deferred capability state = %q, want %q", envelope.Capabilities[0].State, JEVObservationDeferred)
	}
	if envelope.FirstMissingStage != 3 || envelope.NextOperation != "provide_explicit_external_boundary" {
		t.Fatalf("deferred boundary = %d/%q", envelope.FirstMissingStage, envelope.NextOperation)
	}
}

func TestNewJEVObservationEnvelopePreservesUnknownHover(t *testing.T) {
	hover := DiscoverCapabilityQueryHover("What is the weather?", "entity user")
	envelope, err := NewJEVObservationEnvelopeFromCapabilityQueryHover(hover, "corr-unknown", "catalog-3", "schema-3", 0)
	if err != nil {
		t.Fatalf("project unknown hover: %v", err)
	}
	if len(envelope.Capabilities) != 1 || envelope.Capabilities[0].ID != "capability.query" || envelope.Capabilities[0].State != JEVObservationUnknown {
		t.Fatalf("unknown capabilities = %#v", envelope.Capabilities)
	}
	if envelope.NextOperation != "ask_clarifying_question" {
		t.Fatalf("unknown next operation = %q, want ask_clarifying_question", envelope.NextOperation)
	}
}

func TestNewJEVObservationEnvelopeRequiresUnknownBoundary(t *testing.T) {
	hover := DiscoverCapabilityQueryHover("What is the weather?", "entity user")
	if _, err := NewJEVObservationEnvelopeFromCapabilityQueryHover(hover, "corr-missing", "catalog-4", "schema-4", -1); err == nil {
		t.Fatal("unknown hover without a missing stage should be rejected")
	}
}
