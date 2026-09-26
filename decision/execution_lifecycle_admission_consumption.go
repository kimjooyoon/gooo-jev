package decision

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const ExecutionLifecycleAdmissionConsumptionSchemaV1 = "gooo/jev-execution-lifecycle-admission-consumption/v1"

type ExecutionLifecycleAdmissionConsumptionStatus string

const (
	ExecutionLifecycleAdmissionConsumed        ExecutionLifecycleAdmissionConsumptionStatus = "consumed"
	ExecutionLifecycleAdmissionConsumptionUnknown ExecutionLifecycleAdmissionConsumptionStatus = "unknown"
)

type ExecutionLifecycleAdmissionConsumption struct {
	Schema            string
	AdmissionDigest   string
	ObservationDigest string
	UseDigest         string
	MaxUses           int
	ConsumedUses      int
	Status            ExecutionLifecycleAdmissionConsumptionStatus
	MissingStage      string
	ConsumedAt        time.Time
	ConsumptionDigest string
}

func ConsumeExecutionLifecycleAdmission(
	admission ExecutionLifecycleAdmission,
	useDigest string,
	maxUses int,
	consumedUses int,
	consumedAt time.Time,
) (ExecutionLifecycleAdmissionConsumption, error) {
	if consumedAt.IsZero() {
		return ExecutionLifecycleAdmissionConsumption{}, errors.New("execution lifecycle admission consumption time is required")
	}
	consumption := ExecutionLifecycleAdmissionConsumption{
		Schema:            ExecutionLifecycleAdmissionConsumptionSchemaV1,
		AdmissionDigest:   admission.AdmissionDigest,
		ObservationDigest: admission.ObservationDigest,
		UseDigest:         strings.TrimSpace(useDigest),
		MaxUses:           maxUses,
		ConsumedUses:      consumedUses,
		Status:            ExecutionLifecycleAdmissionConsumed,
		ConsumedAt:        consumedAt.UTC(),
	}
	switch {
	case admission.Validate() != nil:
		consumption.Status = ExecutionLifecycleAdmissionConsumptionUnknown
		consumption.MissingStage = "admission"
	case admission.Status != ExecutionLifecycleAdmitted:
		consumption.Status = ExecutionLifecycleAdmissionConsumptionUnknown
		if strings.TrimSpace(admission.MissingStage) == "" {
			consumption.MissingStage = "admission-status"
		} else {
			consumption.MissingStage = "admission:" + admission.MissingStage
		}
	case strings.TrimSpace(useDigest) == "":
		consumption.Status = ExecutionLifecycleAdmissionConsumptionUnknown
		consumption.MissingStage = "use"
	case maxUses <= 0:
		consumption.Status = ExecutionLifecycleAdmissionConsumptionUnknown
		consumption.MissingStage = "use-limit"
	case consumedUses <= 0:
		consumption.Status = ExecutionLifecycleAdmissionConsumptionUnknown
		consumption.MissingStage = "use-count"
	case consumedUses > maxUses:
		consumption.Status = ExecutionLifecycleAdmissionConsumptionUnknown
		consumption.MissingStage = "use-limit"
	}
	if err := consumption.assignDigest(); err != nil {
		return ExecutionLifecycleAdmissionConsumption{}, fmt.Errorf("digest execution lifecycle admission consumption: %w", err)
	}
	if err := consumption.Validate(); err != nil {
		return ExecutionLifecycleAdmissionConsumption{}, err
	}
	return consumption, nil
}

func (consumption *ExecutionLifecycleAdmissionConsumption) assignDigest() error {
	digest, err := consumption.computeDigest()
	if err != nil {
		return err
	}
	consumption.ConsumptionDigest = digest
	return nil
}

func (consumption ExecutionLifecycleAdmissionConsumption) Validate() error {
	if err := consumption.validateShape(); err != nil {
		return err
	}
	expected, err := consumption.computeDigest()
	if err != nil {
		return fmt.Errorf("digest execution lifecycle admission consumption: %w", err)
	}
	if consumption.ConsumptionDigest != expected {
		return errors.New("execution lifecycle admission consumption digest mismatch")
	}
	return nil
}

func (consumption ExecutionLifecycleAdmissionConsumption) validateShape() error {
	if consumption.Schema != ExecutionLifecycleAdmissionConsumptionSchemaV1 {
		return errors.New("unsupported execution lifecycle admission consumption schema")
	}
	if consumption.ConsumedAt.IsZero() {
		return errors.New("execution lifecycle admission consumption time is required")
	}
	if strings.TrimSpace(consumption.ConsumptionDigest) == "" {
		return errors.New("execution lifecycle admission consumption digest is required")
	}
	switch consumption.Status {
	case ExecutionLifecycleAdmissionConsumed:
		if strings.TrimSpace(consumption.AdmissionDigest) == "" ||
			strings.TrimSpace(consumption.ObservationDigest) == "" ||
			strings.TrimSpace(consumption.UseDigest) == "" ||
			consumption.MaxUses <= 0 ||
			consumption.ConsumedUses <= 0 ||
			consumption.ConsumedUses > consumption.MaxUses ||
			strings.TrimSpace(consumption.MissingStage) != "" {
			return errors.New("consumed execution lifecycle admission is incomplete")
		}
	case ExecutionLifecycleAdmissionConsumptionUnknown:
		if strings.TrimSpace(consumption.MissingStage) == "" {
			return errors.New("unknown execution lifecycle admission consumption requires a missing stage")
		}
	default:
		return fmt.Errorf("unsupported execution lifecycle admission consumption status %q", consumption.Status)
	}
	return nil
}

func (consumption ExecutionLifecycleAdmissionConsumption) computeDigest() (string, error) {
	return Digest(struct {
		Schema            string
		AdmissionDigest   string
		ObservationDigest string
		UseDigest         string
		MaxUses           int
		ConsumedUses      int
		Status            ExecutionLifecycleAdmissionConsumptionStatus
		MissingStage      string
		ConsumedAt        time.Time
	}{
		Schema:            consumption.Schema,
		AdmissionDigest:   consumption.AdmissionDigest,
		ObservationDigest: consumption.ObservationDigest,
		UseDigest:         consumption.UseDigest,
		MaxUses:           consumption.MaxUses,
		ConsumedUses:      consumption.ConsumedUses,
		Status:            consumption.Status,
		MissingStage:      consumption.MissingStage,
		ConsumedAt:        consumption.ConsumedAt.UTC(),
	})
}
