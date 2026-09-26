package decision

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const ExecutionLifecycleAttestationWindowSchemaV1 = "gooo/jev-execution-lifecycle-attestation-window/v1"

type ExecutionLifecycleAttestationWindowStatus string

const (
	ExecutionLifecycleAttestationWithinWindow ExecutionLifecycleAttestationWindowStatus = "within-window"
	ExecutionLifecycleAttestationWindowUnknown ExecutionLifecycleAttestationWindowStatus = "unknown"
)

type ExecutionLifecycleAttestationWindow struct {
	Schema              string
	AttestationDigest   string
	ObservationDigest   string
	NotBefore           time.Time
	NotAfter            time.Time
	Status              ExecutionLifecycleAttestationWindowStatus
	MissingStage        string
	ObservedAt          time.Time
	WindowDigest        string
}

func ObserveExecutionLifecycleAttestationWindow(
	attestation ExecutionLifecycleAttestation,
	notBefore time.Time,
	notAfter time.Time,
) (ExecutionLifecycleAttestationWindow, error) {
	if notBefore.IsZero() || notAfter.IsZero() {
		return ExecutionLifecycleAttestationWindow{}, errors.New("execution lifecycle attestation window bounds are required")
	}
	window := ExecutionLifecycleAttestationWindow{
		Schema:            ExecutionLifecycleAttestationWindowSchemaV1,
		AttestationDigest: attestation.AttestationDigest,
		ObservationDigest: attestation.ObservationDigest,
		NotBefore:         notBefore.UTC(),
		NotAfter:          notAfter.UTC(),
		Status:            ExecutionLifecycleAttestationWithinWindow,
		ObservedAt:        attestation.ObservedAt.UTC(),
	}
	switch {
	case notAfter.Before(notBefore):
		window.Status = ExecutionLifecycleAttestationWindowUnknown
		window.MissingStage = "window"
	case attestation.Validate() != nil:
		window.Status = ExecutionLifecycleAttestationWindowUnknown
		window.MissingStage = "attestation"
	case attestation.Status != ExecutionLifecycleAttested:
		window.Status = ExecutionLifecycleAttestationWindowUnknown
		if strings.TrimSpace(attestation.MissingStage) == "" {
			window.MissingStage = "attestation-status"
		} else {
			window.MissingStage = "attestation:" + attestation.MissingStage
		}
	case attestation.ObservedAt.Before(notBefore):
		window.Status = ExecutionLifecycleAttestationWindowUnknown
		window.MissingStage = "before-window"
	case attestation.ObservedAt.After(notAfter):
		window.Status = ExecutionLifecycleAttestationWindowUnknown
		window.MissingStage = "after-window"
	}
	if err := window.assignDigest(); err != nil {
		return ExecutionLifecycleAttestationWindow{}, fmt.Errorf("digest execution lifecycle attestation window: %w", err)
	}
	if err := window.Validate(); err != nil {
		return ExecutionLifecycleAttestationWindow{}, err
	}
	return window, nil
}

func (window *ExecutionLifecycleAttestationWindow) assignDigest() error {
	digest, err := window.computeDigest()
	if err != nil {
		return err
	}
	window.WindowDigest = digest
	return nil
}

func (window ExecutionLifecycleAttestationWindow) Validate() error {
	if err := window.validateShape(); err != nil {
		return err
	}
	expected, err := window.computeDigest()
	if err != nil {
		return fmt.Errorf("digest execution lifecycle attestation window: %w", err)
	}
	if window.WindowDigest != expected {
		return errors.New("execution lifecycle attestation window digest mismatch")
	}
	return nil
}

func (window ExecutionLifecycleAttestationWindow) validateShape() error {
	if window.Schema != ExecutionLifecycleAttestationWindowSchemaV1 {
		return errors.New("unsupported execution lifecycle attestation window schema")
	}
	if window.NotBefore.IsZero() || window.NotAfter.IsZero() {
		return errors.New("execution lifecycle attestation window bounds are required")
	}
	if window.ObservedAt.IsZero() {
		return errors.New("execution lifecycle attestation window observation time is required")
	}
	if strings.TrimSpace(window.WindowDigest) == "" {
		return errors.New("execution lifecycle attestation window digest is required")
	}
	switch window.Status {
	case ExecutionLifecycleAttestationWithinWindow:
		if strings.TrimSpace(window.AttestationDigest) == "" ||
			strings.TrimSpace(window.ObservationDigest) == "" ||
			window.NotAfter.Before(window.NotBefore) ||
			strings.TrimSpace(window.MissingStage) != "" {
			return errors.New("within-window execution lifecycle attestation is incomplete")
		}
	case ExecutionLifecycleAttestationWindowUnknown:
		if strings.TrimSpace(window.MissingStage) == "" {
			return errors.New("unknown execution lifecycle attestation window requires a missing stage")
		}
	default:
		return fmt.Errorf("unsupported execution lifecycle attestation window status %q", window.Status)
	}
	return nil
}

func (window ExecutionLifecycleAttestationWindow) computeDigest() (string, error) {
	return Digest(struct {
		Schema              string
		AttestationDigest   string
		ObservationDigest   string
		NotBefore           time.Time
		NotAfter            time.Time
		Status              ExecutionLifecycleAttestationWindowStatus
		MissingStage        string
		ObservedAt          time.Time
	}{
		Schema:            window.Schema,
		AttestationDigest: window.AttestationDigest,
		ObservationDigest: window.ObservationDigest,
		NotBefore:         window.NotBefore.UTC(),
		NotAfter:          window.NotAfter.UTC(),
		Status:            window.Status,
		MissingStage:      window.MissingStage,
		ObservedAt:        window.ObservedAt.UTC(),
	})
}