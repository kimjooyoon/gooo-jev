package gooo

import (
	"fmt"
	"strings"
)

// NewJEVObservationEnvelopeFromCapabilityQueryHover binds an editor-facing
// capability observation to the JEV envelope without turning discovery into
// execution or authorization. The catalog and schema digests are supplied by
// the caller because this adapter does not own either registry.
func NewJEVObservationEnvelopeFromCapabilityQueryHover(
	hover CapabilityQueryHover,
	correlationID string,
	catalogDigest string,
	schemaDigest string,
	firstMissingStage int,
) (JEVObservationEnvelope, error) {
	if err := hover.Validate(); err != nil {
		return JEVObservationEnvelope{}, fmt.Errorf("capability query hover: %w", err)
	}
	for name, value := range map[string]string{
		"correlation_id": correlationID,
		"catalog_digest": catalogDigest,
		"schema_digest":  schemaDigest,
	} {
		if strings.TrimSpace(value) == "" {
			return JEVObservationEnvelope{}, fmt.Errorf("%s is required", name)
		}
	}

	stage := firstMissingStage
	nextOperation := ""
	switch hover.Status {
	case CapabilityQueryAvailable:
		stage = -1
		nextOperation = "inspect_next_operation"
	case CapabilityQueryDeferred:
		if stage < 0 {
			return JEVObservationEnvelope{}, fmt.Errorf("deferred hover requires first_missing_stage")
		}
		nextOperation = "provide_explicit_external_boundary"
	case CapabilityQueryUnknown:
		if stage < 0 {
			return JEVObservationEnvelope{}, fmt.Errorf("unknown hover requires first_missing_stage")
		}
		nextOperation = "ask_clarifying_question"
	default:
		return JEVObservationEnvelope{}, fmt.Errorf("unsupported hover status %q", hover.Status)
	}

	capabilities := make([]JEVCapabilityObservation, 0, len(hover.Capabilities))
	for _, capability := range hover.Capabilities {
		state := JEVObservationStateUnknown
		switch capability.State {
		case CapabilityQueryAvailable:
			state = JEVObservationAvailable
		case CapabilityQueryDeferred:
			state = JEVObservationDeferred
		case CapabilityQueryUnknown:
			state = JEVObservationUnknown
		default:
			return JEVObservationEnvelope{}, fmt.Errorf("unsupported capability state %q", capability.State)
		}
		capabilities = append(capabilities, JEVCapabilityObservation{ID: capability.ID, State: state})
	}
	if hover.Status == CapabilityQueryUnknown && len(capabilities) == 0 {
		capabilities = append(capabilities, JEVCapabilityObservation{
			ID:    "capability.query",
			State: JEVObservationUnknown,
		})
	}

	envelope := JEVObservationEnvelope{
		Version:           JEVObservationEnvelopeVersion,
		CorrelationID:     correlationID,
		SourceDigest:      hover.SourceDigest,
		DeclarationDigest: hover.SourceDigest,
		CatalogDigest:     catalogDigest,
		SchemaDigest:      schemaDigest,
		EvidenceDigest:    hover.EvidenceDigest,
		Capabilities:      capabilities,
		FirstMissingStage: stage,
		NextOperation:     nextOperation,
		Executed:          false,
		Authorizing:       false,
	}
	if _, err := envelope.CanonicalDigest(); err != nil {
		return JEVObservationEnvelope{}, fmt.Errorf("JEV observation envelope: %w", err)
	}
	return envelope, nil
}
