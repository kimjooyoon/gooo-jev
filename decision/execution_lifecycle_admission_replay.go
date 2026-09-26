package decision

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const ExecutionLifecycleAdmissionReplaySchemaV1 = "gooo/jev-execution-lifecycle-admission-replay/v1"

type ExecutionLifecycleAdmissionReplayStatus string

const (
	ExecutionLifecycleAdmissionReplayed        ExecutionLifecycleAdmissionReplayStatus = "replayed"
	ExecutionLifecycleAdmissionReplayUnknown   ExecutionLifecycleAdmissionReplayStatus = "unknown"
)

type ExecutionLifecycleAdmissionReplay struct {
	Schema            string
	AdmissionDigest   string
	WindowDigest      string
	AttestationDigest string
	ObservationDigest string
	Status            ExecutionLifecycleAdmissionReplayStatus
	MissingStage      string
	ObservedAt        time.Time
	ReplayDigest      string
}

func ReplayExecutionLifecycleAdmission(
	admission ExecutionLifecycleAdmission,
	observationDigest string,
	observedAt time.Time,
) (ExecutionLifecycleAdmissionReplay, error) {
	if observedAt.IsZero() {
		return ExecutionLifecycleAdmissionReplay{}, errors.New("execution lifecycle admission replay observation time is required")
	}
	replay := ExecutionLifecycleAdmissionReplay{
		Schema:            ExecutionLifecycleAdmissionReplaySchemaV1,
		AdmissionDigest:   admission.AdmissionDigest,
		WindowDigest:      admission.WindowDigest,
		AttestationDigest: admission.AttestationDigest,
		ObservationDigest: strings.TrimSpace(observationDigest),
		Status:            ExecutionLifecycleAdmissionReplayed,
		ObservedAt:        observedAt.UTC(),
	}
	switch {
	case admission.Validate() != nil:
		replay.Status = ExecutionLifecycleAdmissionReplayUnknown
		replay.MissingStage = "admission"
	case admission.Status != ExecutionLifecycleAdmitted:
		replay.Status = ExecutionLifecycleAdmissionReplayUnknown
		if strings.TrimSpace(admission.MissingStage) == "" {
			replay.MissingStage = "admission-status"
		} else {
			replay.MissingStage = "admission:" + admission.MissingStage
		}
	case strings.TrimSpace(observationDigest) == "":
		replay.Status = ExecutionLifecycleAdmissionReplayUnknown
		replay.MissingStage = "observation"
	case admission.ObservationDigest != strings.TrimSpace(observationDigest):
		replay.Status = ExecutionLifecycleAdmissionReplayUnknown
		replay.MissingStage = "observation-replay"
	}
	if err := replay.assignDigest(); err != nil {
		return ExecutionLifecycleAdmissionReplay{}, fmt.Errorf("digest execution lifecycle admission replay: %w", err)
	}
	if err := replay.Validate(); err != nil {
		return ExecutionLifecycleAdmissionReplay{}, err
	}
	return replay, nil
}

func (replay *ExecutionLifecycleAdmissionReplay) assignDigest() error {
	digest, err := replay.computeDigest()
	if err != nil {
		return err
	}
	replay.ReplayDigest = digest
	return nil
}

func (replay ExecutionLifecycleAdmissionReplay) Validate() error {
	if err := replay.validateShape(); err != nil {
		return err
	}
	expected, err := replay.computeDigest()
	if err != nil {
		return fmt.Errorf("digest execution lifecycle admission replay: %w", err)
	}
	if replay.ReplayDigest != expected {
		return errors.New("execution lifecycle admission replay digest mismatch")
	}
	return nil
}

func (replay ExecutionLifecycleAdmissionReplay) validateShape() error {
	if replay.Schema != ExecutionLifecycleAdmissionReplaySchemaV1 {
		return errors.New("unsupported execution lifecycle admission replay schema")
	}
	if replay.ObservedAt.IsZero() {
		return errors.New("execution lifecycle admission replay observation time is required")
	}
	if strings.TrimSpace(replay.ReplayDigest) == "" {
		return errors.New("execution lifecycle admission replay digest is required")
	}
	switch replay.Status {
	case ExecutionLifecycleAdmissionReplayed:
		if strings.TrimSpace(replay.AdmissionDigest) == "" ||
			strings.TrimSpace(replay.WindowDigest) == "" ||
			strings.TrimSpace(replay.AttestationDigest) == "" ||
			strings.TrimSpace(replay.ObservationDigest) == "" ||
			strings.TrimSpace(replay.MissingStage) != "" {
			return errors.New("replayed execution lifecycle admission is incomplete")
		}
	case ExecutionLifecycleAdmissionReplayUnknown:
		if strings.TrimSpace(replay.MissingStage) == "" {
			return errors.New("unknown execution lifecycle admission replay requires a missing stage")
		}
	default:
		return fmt.Errorf("unsupported execution lifecycle admission replay status %q", replay.Status)
	}
	return nil
}

func (replay ExecutionLifecycleAdmissionReplay) computeDigest() (string, error) {
	return Digest(struct {
		Schema            string
		AdmissionDigest   string
		WindowDigest      string
		AttestationDigest string
		ObservationDigest string
		Status            ExecutionLifecycleAdmissionReplayStatus
		MissingStage      string
		ObservedAt        time.Time
	}{
		Schema:            replay.Schema,
		AdmissionDigest:   replay.AdmissionDigest,
		WindowDigest:      replay.WindowDigest,
		AttestationDigest: replay.AttestationDigest,
		ObservationDigest: replay.ObservationDigest,
		Status:            replay.Status,
		MissingStage:      replay.MissingStage,
		ObservedAt:        replay.ObservedAt.UTC(),
	})
}
