package decision

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const ExecutionLifecycleAdmissionSchemaV1 = "gooo/jev-execution-lifecycle-admission/v1"

type ExecutionLifecycleAdmissionStatus string

const (
	ExecutionLifecycleAdmitted        ExecutionLifecycleAdmissionStatus = "admitted"
	ExecutionLifecycleAdmissionUnknown ExecutionLifecycleAdmissionStatus = "unknown"
)

type ExecutionLifecycleAdmission struct {
	Schema            string
	WindowDigest      string
	AttestationDigest string
	ObservationDigest string
	Status            ExecutionLifecycleAdmissionStatus
	MissingStage      string
	AdmittedAt        time.Time
	AdmissionDigest   string
}

func AdmitExecutionLifecycleAttestationWindow(
	window ExecutionLifecycleAttestationWindow,
	admittedAt time.Time,
) (ExecutionLifecycleAdmission, error) {
	if admittedAt.IsZero() {
		return ExecutionLifecycleAdmission{}, errors.New("execution lifecycle admission time is required")
	}
	admission := ExecutionLifecycleAdmission{
		Schema:            ExecutionLifecycleAdmissionSchemaV1,
		WindowDigest:      window.WindowDigest,
		AttestationDigest: window.AttestationDigest,
		ObservationDigest: window.ObservationDigest,
		Status:            ExecutionLifecycleAdmitted,
		AdmittedAt:        admittedAt.UTC(),
	}
	switch {
	case window.Validate() != nil:
		admission.Status = ExecutionLifecycleAdmissionUnknown
		admission.MissingStage = "window"
	case window.Status != ExecutionLifecycleAttestationWithinWindow:
		admission.Status = ExecutionLifecycleAdmissionUnknown
		if strings.TrimSpace(window.MissingStage) == "" {
			admission.MissingStage = "window-status"
		} else {
			admission.MissingStage = "window:" + window.MissingStage
		}
	}
	if err := admission.assignDigest(); err != nil {
		return ExecutionLifecycleAdmission{}, fmt.Errorf("digest execution lifecycle admission: %w", err)
	}
	if err := admission.Validate(); err != nil {
		return ExecutionLifecycleAdmission{}, err
	}
	return admission, nil
}

func (admission *ExecutionLifecycleAdmission) assignDigest() error {
	digest, err := admission.computeDigest()
	if err != nil {
		return err
	}
	admission.AdmissionDigest = digest
	return nil
}

func (admission ExecutionLifecycleAdmission) Validate() error {
	if err := admission.validateShape(); err != nil {
		return err
	}
	expected, err := admission.computeDigest()
	if err != nil {
		return fmt.Errorf("digest execution lifecycle admission: %w", err)
	}
	if admission.AdmissionDigest != expected {
		return errors.New("execution lifecycle admission digest mismatch")
	}
	return nil
}

func (admission ExecutionLifecycleAdmission) validateShape() error {
	if admission.Schema != ExecutionLifecycleAdmissionSchemaV1 {
		return errors.New("unsupported execution lifecycle admission schema")
	}
	if admission.AdmittedAt.IsZero() {
		return errors.New("execution lifecycle admission time is required")
	}
	if strings.TrimSpace(admission.AdmissionDigest) == "" {
		return errors.New("execution lifecycle admission digest is required")
	}
	switch admission.Status {
	case ExecutionLifecycleAdmitted:
		if strings.TrimSpace(admission.WindowDigest) == "" ||
			strings.TrimSpace(admission.AttestationDigest) == "" ||
			strings.TrimSpace(admission.ObservationDigest) == "" ||
			strings.TrimSpace(admission.MissingStage) != "" {
			return errors.New("admitted execution lifecycle is incomplete")
		}
	case ExecutionLifecycleAdmissionUnknown:
		if strings.TrimSpace(admission.MissingStage) == "" {
			return errors.New("unknown execution lifecycle admission requires a missing stage")
		}
	default:
		return fmt.Errorf("unsupported execution lifecycle admission status %q", admission.Status)
	}
	return nil
}

func (admission ExecutionLifecycleAdmission) computeDigest() (string, error) {
	return Digest(struct {
		Schema            string
		WindowDigest      string
		AttestationDigest string
		ObservationDigest string
		Status            ExecutionLifecycleAdmissionStatus
		MissingStage      string
		AdmittedAt        time.Time
	}{
		Schema:            admission.Schema,
		WindowDigest:      admission.WindowDigest,
		AttestationDigest: admission.AttestationDigest,
		ObservationDigest: admission.ObservationDigest,
		Status:            admission.Status,
		MissingStage:      admission.MissingStage,
		AdmittedAt:        admission.AdmittedAt.UTC(),
	})
}
