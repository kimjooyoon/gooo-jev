package decision

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const ExecutionLifecycleReplaySchemaV1 = "gooo/jev-execution-lifecycle-replay/v1"

type ExecutionLifecycleReplayStatus string

const (
	ExecutionLifecycleReplayObserved ExecutionLifecycleReplayStatus = "observed"
	ExecutionLifecycleReplayUnknown  ExecutionLifecycleReplayStatus = "unknown"
)

type ExecutionLifecycleReplayObservation struct {
	Schema                    string
	SuspensionDigest          string
	ResumptionDigest          string
	ReverseObservationDigest  string
	Status                    ExecutionLifecycleReplayStatus
	MissingStage              string
	ObservedAt                time.Time
	ReplayDigest              string
}

func NewExecutionLifecycleReplay(
	suspension ExecutionLifecycleReceipt,
	resumption ExecutionLifecycleReceipt,
	reverse ReverseObservation,
	observedAt time.Time,
) (ExecutionLifecycleReplayObservation, error) {
	if observedAt.IsZero() {
		return ExecutionLifecycleReplayObservation{}, errors.New("execution lifecycle replay observation time is required")
	}
	observation := ExecutionLifecycleReplayObservation{
		Schema:                   ExecutionLifecycleReplaySchemaV1,
		SuspensionDigest:         suspension.LifecycleDigest,
		ResumptionDigest:         resumption.LifecycleDigest,
		ReverseObservationDigest: reverse.ObservationDigest,
		Status:                   ExecutionLifecycleReplayObserved,
		ObservedAt:               observedAt.UTC(),
	}
	switch {
	case suspension.Validate() != nil:
		observation.Status = ExecutionLifecycleReplayUnknown
		observation.MissingStage = "suspension"
	case resumption.Validate() != nil:
		observation.Status = ExecutionLifecycleReplayUnknown
		observation.MissingStage = "resumption"
	case resumption.ParentLifecycleDigest != suspension.LifecycleDigest:
		observation.Status = ExecutionLifecycleReplayUnknown
		observation.MissingStage = "lifecycle-binding"
	case reverse.Validate() != nil:
		observation.Status = ExecutionLifecycleReplayUnknown
		observation.MissingStage = "reverse-observation"
	case reverse.Status != ProvenanceObserved:
		observation.Status = ExecutionLifecycleReplayUnknown
		observation.MissingStage = "reverse-observation-status"
	}
	if err := observation.assignDigest(); err != nil {
		return ExecutionLifecycleReplayObservation{}, fmt.Errorf("digest execution lifecycle replay: %w", err)
	}
	if err := observation.Validate(); err != nil {
		return ExecutionLifecycleReplayObservation{}, err
	}
	return observation, nil
}

func (observation *ExecutionLifecycleReplayObservation) assignDigest() error {
	digest, err := observation.computeDigest()
	if err != nil {
		return err
	}
	observation.ReplayDigest = digest
	return nil
}

func (observation ExecutionLifecycleReplayObservation) Validate() error {
	if err := observation.validateShape(); err != nil {
		return err
	}
	expected, err := observation.computeDigest()
	if err != nil {
		return fmt.Errorf("digest execution lifecycle replay: %w", err)
	}
	if observation.ReplayDigest != expected {
		return errors.New("execution lifecycle replay digest mismatch")
	}
	return nil
}

func (observation ExecutionLifecycleReplayObservation) validateShape() error {
	if observation.Schema != ExecutionLifecycleReplaySchemaV1 {
		return errors.New("unsupported execution lifecycle replay schema")
	}
	if observation.ObservedAt.IsZero() {
		return errors.New("execution lifecycle replay observation time is required")
	}
	if strings.TrimSpace(observation.ReplayDigest) == "" {
		return errors.New("execution lifecycle replay digest is required")
	}
	switch observation.Status {
	case ExecutionLifecycleReplayObserved:
		if strings.TrimSpace(observation.SuspensionDigest) == "" ||
			strings.TrimSpace(observation.ResumptionDigest) == "" ||
			strings.TrimSpace(observation.ReverseObservationDigest) == "" ||
			strings.TrimSpace(observation.MissingStage) != "" {
			return errors.New("observed execution lifecycle replay is incomplete")
		}
	case ExecutionLifecycleReplayUnknown:
		if strings.TrimSpace(observation.MissingStage) == "" {
			return errors.New("unknown execution lifecycle replay requires a missing stage")
		}
	default:
		return fmt.Errorf("unsupported execution lifecycle replay status %q", observation.Status)
	}
	return nil
}

func (observation ExecutionLifecycleReplayObservation) computeDigest() (string, error) {
	return Digest(struct {
		Schema                   string
		SuspensionDigest         string
		ResumptionDigest         string
		ReverseObservationDigest string
		Status                   ExecutionLifecycleReplayStatus
		MissingStage             string
		ObservedAt               time.Time
	}{
		Schema:                   observation.Schema,
		SuspensionDigest:         observation.SuspensionDigest,
		ResumptionDigest:         observation.ResumptionDigest,
		ReverseObservationDigest: observation.ReverseObservationDigest,
		Status:                   observation.Status,
		MissingStage:             observation.MissingStage,
		ObservedAt:               observation.ObservedAt.UTC(),
	})
}
