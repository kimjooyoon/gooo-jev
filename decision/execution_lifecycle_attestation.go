package decision

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const ExecutionLifecycleAttestationSchemaV1 = "gooo/jev-execution-lifecycle-attestation/v1"

type ExecutionLifecycleAttestationStatus string

const (
	ExecutionLifecycleAttested          ExecutionLifecycleAttestationStatus = "attested"
	ExecutionLifecycleAttestationUnknown ExecutionLifecycleAttestationStatus = "unknown"
)

type ExecutionLifecycleAttestation struct {
	Schema              string
	AuthorizationDigest string
	ObservationDigest   string
	Status              ExecutionLifecycleAttestationStatus
	MissingStage        string
	ObservedAt          time.Time
	AttestationDigest   string
}

func NewExecutionLifecycleAttestation(
	authorization ExecutionLifecycleAuthorization,
	observationDigest string,
	observedAt time.Time,
) (ExecutionLifecycleAttestation, error) {
	if observedAt.IsZero() {
		return ExecutionLifecycleAttestation{}, errors.New("execution lifecycle attestation observation time is required")
	}
	attestation := ExecutionLifecycleAttestation{
		Schema:              ExecutionLifecycleAttestationSchemaV1,
		AuthorizationDigest: authorization.AuthorizationDigest,
		ObservationDigest:   strings.TrimSpace(observationDigest),
		Status:              ExecutionLifecycleAttested,
		ObservedAt:          observedAt.UTC(),
	}
	switch {
	case authorization.Validate() != nil:
		attestation.Status = ExecutionLifecycleAttestationUnknown
		attestation.MissingStage = "authorization"
	case authorization.Status != ExecutionLifecycleAuthorized:
		attestation.Status = ExecutionLifecycleAttestationUnknown
		if strings.TrimSpace(authorization.MissingStage) == "" {
			attestation.MissingStage = "authorization-status"
		} else {
			attestation.MissingStage = "authorization:" + authorization.MissingStage
		}
	case strings.TrimSpace(observationDigest) == "":
		attestation.Status = ExecutionLifecycleAttestationUnknown
		attestation.MissingStage = "observation"
	}
	if err := attestation.assignDigest(); err != nil {
		return ExecutionLifecycleAttestation{}, fmt.Errorf("digest execution lifecycle attestation: %w", err)
	}
	if err := attestation.Validate(); err != nil {
		return ExecutionLifecycleAttestation{}, err
	}
	return attestation, nil
}

func (attestation *ExecutionLifecycleAttestation) assignDigest() error {
	digest, err := attestation.computeDigest()
	if err != nil {
		return err
	}
	attestation.AttestationDigest = digest
	return nil
}

func (attestation ExecutionLifecycleAttestation) Validate() error {
	if err := attestation.validateShape(); err != nil {
		return err
	}
	expected, err := attestation.computeDigest()
	if err != nil {
		return fmt.Errorf("digest execution lifecycle attestation: %w", err)
	}
	if attestation.AttestationDigest != expected {
		return errors.New("execution lifecycle attestation digest mismatch")
	}
	return nil
}

func (attestation ExecutionLifecycleAttestation) validateShape() error {
	if attestation.Schema != ExecutionLifecycleAttestationSchemaV1 {
		return errors.New("unsupported execution lifecycle attestation schema")
	}
	if attestation.ObservedAt.IsZero() {
		return errors.New("execution lifecycle attestation observation time is required")
	}
	if strings.TrimSpace(attestation.AttestationDigest) == "" {
		return errors.New("execution lifecycle attestation digest is required")
	}
	switch attestation.Status {
	case ExecutionLifecycleAttested:
		if strings.TrimSpace(attestation.AuthorizationDigest) == "" ||
			strings.TrimSpace(attestation.ObservationDigest) == "" ||
			strings.TrimSpace(attestation.MissingStage) != "" {
			return errors.New("attested execution lifecycle is incomplete")
		}
	case ExecutionLifecycleAttestationUnknown:
		if strings.TrimSpace(attestation.MissingStage) == "" {
			return errors.New("unknown execution lifecycle attestation requires a missing stage")
		}
	default:
		return fmt.Errorf("unsupported execution lifecycle attestation status %q", attestation.Status)
	}
	return nil
}

func (attestation ExecutionLifecycleAttestation) computeDigest() (string, error) {
	return Digest(struct {
		Schema              string
		AuthorizationDigest string
		ObservationDigest   string
		Status              ExecutionLifecycleAttestationStatus
		MissingStage        string
		ObservedAt          time.Time
	}{
		Schema:              attestation.Schema,
		AuthorizationDigest: attestation.AuthorizationDigest,
		ObservationDigest:   attestation.ObservationDigest,
		Status:              attestation.Status,
		MissingStage:        attestation.MissingStage,
		ObservedAt:          attestation.ObservedAt.UTC(),
	})
}