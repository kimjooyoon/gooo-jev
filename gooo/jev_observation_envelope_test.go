package gooo

import "testing"

func TestJEVObservationEnvelopeRejectsExecutionAndAuthorization(t *testing.T) {
	base := JEVObservationEnvelope{
		Version:           JEVObservationEnvelopeVersion,
		CorrelationID:     "corr-1",
		SourceDigest:      "source-1",
		DeclarationDigest: "decl-1",
		CatalogDigest:     "catalog-1",
		SchemaDigest:      "schema-1",
		EvidenceDigest:    "evidence-1",
		Capabilities: []JEVCapabilityObservation{{
			ID:    "capability.query",
			State: JEVObservationAvailable,
		}},
		FirstMissingStage: -1,
	}
	for name, mutate := range map[string]func(*JEVObservationEnvelope){
		"execution":   func(value *JEVObservationEnvelope) { value.Executed = true },
		"authorization": func(value *JEVObservationEnvelope) { value.Authorizing = true },
	} {
		t.Run(name, func(t *testing.T) {
			value := base
			mutate(&value)
			if err := value.Validate(); err == nil {
				t.Fatalf("expected validation failure")
			}
		})
	}
}

func TestJEVObservationEnvelopeRequiresBoundaryForUnknown(t *testing.T) {
	value := JEVObservationEnvelope{
		Version:           JEVObservationEnvelopeVersion,
		CorrelationID:     "corr-2",
		SourceDigest:      "source-2",
		DeclarationDigest: "decl-2",
		CatalogDigest:     "catalog-2",
		SchemaDigest:      "schema-2",
		EvidenceDigest:    "evidence-2",
		Capabilities: []JEVCapabilityObservation{{
			ID:    "capability.network",
			State: JEVObservationUnknown,
		}},
		FirstMissingStage: 2,
		NextOperation:     "inspect-provider-boundary",
	}
	if _, err := value.CanonicalDigest(); err != nil {
		t.Fatalf("expected valid unknown observation: %v", err)
	}
}
