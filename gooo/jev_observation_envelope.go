package gooo

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

const JEVObservationEnvelopeVersion = "jev.observation.v1"

type JEVObservationState string

const (
	JEVObservationAvailable JEVObservationState = "AVAILABLE"
	JEVObservationDeferred  JEVObservationState = "DEFERRED"
	JEVObservationUnknown   JEVObservationState = "UNKNOWN"
)

type JEVCapabilityObservation struct {
	ID    string               `json:"id"`
	State JEVObservationState  `json:"state"`
}

type JEVObservationEnvelope struct {
	Version           string                       `json:"version"`
	CorrelationID     string                       `json:"correlation_id"`
	SourceDigest      string                       `json:"source_digest"`
	DeclarationDigest string                       `json:"declaration_digest"`
	CatalogDigest     string                       `json:"catalog_digest"`
	SchemaDigest      string                       `json:"schema_digest"`
	EvidenceDigest    string                       `json:"evidence_digest"`
	Capabilities      []JEVCapabilityObservation   `json:"capabilities"`
	FirstMissingStage int                          `json:"first_missing_stage"`
	NextOperation     string                       `json:"next_operation"`
	Executed          bool                         `json:"executed"`
	Authorizing       bool                         `json:"authorizing"`
}

func (e JEVObservationEnvelope) Validate() error {
	if e.Version != JEVObservationEnvelopeVersion {
		return fmt.Errorf("unsupported observation envelope version %q", e.Version)
	}
	for field, value := range map[string]string{
		"correlation_id":     e.CorrelationID,
		"source_digest":      e.SourceDigest,
		"declaration_digest": e.DeclarationDigest,
		"catalog_digest":     e.CatalogDigest,
		"schema_digest":      e.SchemaDigest,
		"evidence_digest":    e.EvidenceDigest,
	} {
		if value == "" {
			return fmt.Errorf("%s is required", field)
		}
	}
	if e.Executed {
		return fmt.Errorf("observation envelope cannot record execution")
	}
	if e.Authorizing {
		return fmt.Errorf("observation envelope cannot authorize an action")
	}
	if e.FirstMissingStage < -1 {
		return fmt.Errorf("first_missing_stage must be -1 or greater")
	}

	seen := make(map[string]struct{}, len(e.Capabilities))
	needsNextOperation := false
	for _, capability := range e.Capabilities {
		if capability.ID == "" {
			return fmt.Errorf("capability id is required")
		}
		if _, ok := seen[capability.ID]; ok {
			return fmt.Errorf("duplicate capability id %q", capability.ID)
		}
		seen[capability.ID] = struct{}{}
		switch capability.State {
		case JEVObservationAvailable:
		case JEVObservationDeferred, JEVObservationUnknown:
			needsNextOperation = true
		default:
			return fmt.Errorf("unsupported capability state %q", capability.State)
		}
	}
	if needsNextOperation {
		if e.FirstMissingStage < 0 {
			return fmt.Errorf("non-available observation requires first_missing_stage")
		}
		if e.NextOperation == "" {
			return fmt.Errorf("non-available observation requires next_operation")
		}
	}
	return nil
}

func (e JEVObservationEnvelope) CanonicalDigest() (string, error) {
	if err := e.Validate(); err != nil {
		return "", err
	}
	encoded, err := json.Marshal(e)
	if err != nil {
		return "", fmt.Errorf("marshal observation envelope: %w", err)
	}
	digest := sha256.Sum256(encoded)
	return hex.EncodeToString(digest[:]), nil
}
