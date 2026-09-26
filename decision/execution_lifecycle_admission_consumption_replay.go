package decision

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const ExecutionLifecycleAdmissionConsumptionReplaySchemaV1 = "gooo/jev-execution-lifecycle-admission-consumption-replay/v1"

type ExecutionLifecycleAdmissionConsumptionReplayStatus string

const (
	ExecutionLifecycleAdmissionConsumptionReplayed        ExecutionLifecycleAdmissionConsumptionReplayStatus = "replayed"
	ExecutionLifecycleAdmissionConsumptionReplayUnknown   ExecutionLifecycleAdmissionConsumptionReplayStatus = "unknown"
)

type ExecutionLifecycleAdmissionConsumptionReplay struct {
	Schema            string
	ConsumptionDigest string
	AdmissionDigest   string
	ObservationDigest string
	UseDigest         string
	MaxUses           int
	ConsumedUses      int
	Status            ExecutionLifecycleAdmissionConsumptionReplayStatus
	MissingStage      string
	ObservedAt        time.Time
	ReplayDigest      string
}

func ReplayExecutionLifecycleAdmissionConsumption(
	consumption ExecutionLifecycleAdmissionConsumption,
	observedAt time.Time,
) (ExecutionLifecycleAdmissionConsumptionReplay, error) {
	if observedAt.IsZero() {
		return ExecutionLifecycleAdmissionConsumptionReplay{}, errors.New("execution lifecycle admission consumption replay observation time is required")
	}
	replay := ExecutionLifecycleAdmissionConsumptionReplay{
		Schema:            ExecutionLifecycleAdmissionConsumptionReplaySchemaV1,
		ConsumptionDigest: consumption.ConsumptionDigest,
		AdmissionDigest:   consumption.AdmissionDigest,
		ObservationDigest: consumption.ObservationDigest,
		UseDigest:         consumption.UseDigest,
		MaxUses:           consumption.MaxUses,
		ConsumedUses:      consumption.ConsumedUses,
		Status:            ExecutionLifecycleAdmissionConsumptionReplayed,
		ObservedAt:        observedAt.UTC(),
	}
	switch {
	case consumption.Validate() != nil:
		replay.Status = ExecutionLifecycleAdmissionConsumptionReplayUnknown
		replay.MissingStage = "consumption"
	case consumption.Status != ExecutionLifecycleAdmissionConsumed:
		replay.Status = ExecutionLifecycleAdmissionConsumptionReplayUnknown
		if strings.TrimSpace(consumption.MissingStage) == "" {
			replay.MissingStage = "consumption-status"
		} else {
			replay.MissingStage = "consumption:" + consumption.MissingStage
		}
	}
	if err := replay.assignDigest(); err != nil {
		return ExecutionLifecycleAdmissionConsumptionReplay{}, fmt.Errorf("digest execution lifecycle admission consumption replay: %w", err)
	}
	if err := replay.Validate(); err != nil {
		return ExecutionLifecycleAdmissionConsumptionReplay{}, err
	}
	return replay, nil
}

func (replay *ExecutionLifecycleAdmissionConsumptionReplay) assignDigest() error {
	digest, err := replay.computeDigest()
	if err != nil {
		return err
	}
	replay.ReplayDigest = digest
	return nil
}

func (replay ExecutionLifecycleAdmissionConsumptionReplay) Validate() error {
	if err := replay.validateShape(); err != nil {
		return err
	}
	expected, err := replay.computeDigest()
	if err != nil {
		return fmt.Errorf("digest execution lifecycle admission consumption replay: %w", err)
	}
	if replay.ReplayDigest != expected {
		return errors.New("execution lifecycle admission consumption replay digest mismatch")
	}
	return nil
}

func (replay ExecutionLifecycleAdmissionConsumptionReplay) validateShape() error {
	if replay.Schema != ExecutionLifecycleAdmissionConsumptionReplaySchemaV1 {
		return errors.New("unsupported execution lifecycle admission consumption replay schema")
	}
	if replay.ObservedAt.IsZero() {
		return errors.New("execution lifecycle admission consumption replay observation time is required")
	}
	if strings.TrimSpace(replay.ReplayDigest) == "" {
		return errors.New("execution lifecycle admission consumption replay digest is required")
	}
	switch replay.Status {
	case ExecutionLifecycleAdmissionConsumptionReplayed:
		if strings.TrimSpace(replay.ConsumptionDigest) == "" ||
			strings.TrimSpace(replay.AdmissionDigest) == "" ||
			strings.TrimSpace(replay.ObservationDigest) == "" ||
			strings.TrimSpace(replay.UseDigest) == "" ||
			replay.MaxUses <= 0 ||
			replay.ConsumedUses <= 0 ||
			replay.ConsumedUses > replay.MaxUses ||
			strings.TrimSpace(replay.MissingStage) != "" {
			return errors.New("replayed execution lifecycle admission consumption is incomplete")
		}
	case ExecutionLifecycleAdmissionConsumptionReplayUnknown:
		if strings.TrimSpace(replay.MissingStage) == "" {
			return errors.New("unknown execution lifecycle admission consumption replay requires a missing stage")
		}
	default:
		return fmt.Errorf("unsupported execution lifecycle admission consumption replay status %q", replay.Status)
	}
	return nil
}

func (replay ExecutionLifecycleAdmissionConsumptionReplay) computeDigest() (string, error) {
	return Digest(struct {
		Schema            string
		ConsumptionDigest string
		AdmissionDigest   string
		ObservationDigest string
		UseDigest         string
		MaxUses           int
		ConsumedUses      int
		Status            ExecutionLifecycleAdmissionConsumptionReplayStatus
		MissingStage      string
		ObservedAt        time.Time
	}{
		Schema:            replay.Schema,
		ConsumptionDigest: replay.ConsumptionDigest,
		AdmissionDigest:   replay.AdmissionDigest,
		ObservationDigest: replay.ObservationDigest,
		UseDigest:         replay.UseDigest,
		MaxUses:           replay.MaxUses,
		ConsumedUses:      replay.ConsumedUses,
		Status:            replay.Status,
		MissingStage:      replay.MissingStage,
		ObservedAt:        replay.ObservedAt.UTC(),
	})
}
