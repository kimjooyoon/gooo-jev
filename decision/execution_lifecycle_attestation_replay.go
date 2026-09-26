package decision

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const ExecutionLifecycleAttestationReplaySchemaV1 = "gooo/jev-execution-lifecycle-attestation-replay/v1"

type ExecutionLifecycleAttestationReplayStatus string

const (
	ExecutionLifecycleAttestationReplayed ExecutionLifecycleAttestationReplayStatus = "replayed"
	ExecutionLifecycleAttestationReplayUnknown ExecutionLifecycleAttestationReplayStatus = "unknown"
)

type ExecutionLifecycleAttestationReplay struct {
	Schema            string
	AttestationDigest string
	ObservationDigest string
	Status            ExecutionLifecycleAttestationReplayStatus
	MissingStage      string
	ObservedAt        time.Time
	ReplayDigest      string
}

func ReplayExecutionLifecycleAttestation(
	attestation ExecutionLifecycleAttestation,
	observationDigest string,
	observedAt time.Time,
) (ExecutionLifecycleAttestationReplay, error) {
	if observedAt.IsZero() {
		return ExecutionLifecycleAttestationReplay{}, errors.New("execution lifecycle attestation replay observation time is required")
	}
	replay := ExecutionLifecycleAttestationReplay{
		Schema:            ExecutionLifecycleAttestationReplaySchemaV1,
		AttestationDigest: attestation.AttestationDigest,
		ObservationDigest: strings.TrimSpace(observationDigest),
		Status:            ExecutionLifecycleAttestationReplayed,
		ObservedAt:        observedAt.UTC(),
	}
	switch {
	case attestation.Validate() != nil:
		replay.Status = ExecutionLifecycleAttestationReplayUnknown
		replay.MissingStage = "attestation"
	case attestation.Status != ExecutionLifecycleAttested:
		replay.Status = ExecutionLifecycleAttestationReplayUnknown
		if strings.TrimSpace(attestation.MissingStage) == "" {
			replay.MissingStage = "attestation-status"
		} else {
			replay.MissingStage = "attestation:" + attestation.MissingStage
		}
	case strings.TrimSpace(observationDigest) == "":
		replay.Status = ExecutionLifecycleAttestationReplayUnknown
		replay.MissingStage = "observation"
	case attestation.ObservationDigest != strings.TrimSpace(observationDigest):
		replay.Status = ExecutionLifecycleAttestationReplayUnknown
		replay.MissingStage = "observation-replay"
	}
	if err := replay.assignDigest(); err != nil {
		return ExecutionLifecycleAttestationReplay{}, fmt.Errorf("digest execution lifecycle attestation replay: %w", err)
	}
	if err := replay.Validate(); err != nil {
		return ExecutionLifecycleAttestationReplay{}, err
	}
	return replay, nil
}

func (replay *ExecutionLifecycleAttestationReplay) assignDigest() error {
	digest, err := replay.computeDigest()
	if err != nil {
		return err
	}
	replay.ReplayDigest = digest
	return nil
}

func (replay ExecutionLifecycleAttestationReplay) Validate() error {
	if err := replay.validateShape(); err != nil {
		return err
	}
	expected, err := replay.computeDigest()
	if err != nil {
		return fmt.Errorf("digest execution lifecycle attestation replay: %w", err)
	}
	if replay.ReplayDigest != expected {
		return errors.New("execution lifecycle attestation replay digest mismatch")
	}
	return nil
}

func (replay ExecutionLifecycleAttestationReplay) validateShape() error {
	if replay.Schema != ExecutionLifecycleAttestationReplaySchemaV1 {
		return errors.New("unsupported execution lifecycle attestation replay schema")
	}
	if replay.ObservedAt.IsZero() {
		return errors.New("execution lifecycle attestation replay observation time is required")
	}
	if strings.TrimSpace(replay.ReplayDigest) == "" {
		return errors.New("execution lifecycle attestation replay digest is required")
	}
	switch replay.Status {
	case ExecutionLifecycleAttestationReplayed:
		if strings.TrimSpace(replay.AttestationDigest) == "" ||
			strings.TrimSpace(replay.ObservationDigest) == "" ||
			strings.TrimSpace(replay.MissingStage) != "" {
			return errors.New("replayed execution lifecycle attestation is incomplete")
		}
	case ExecutionLifecycleAttestationReplayUnknown:
		if strings.TrimSpace(replay.MissingStage) == "" {
			return errors.New("unknown execution lifecycle attestation replay requires a missing stage")
		}
	default:
		return fmt.Errorf("unsupported execution lifecycle attestation replay status %q", replay.Status)
	}
	return nil
}

func (replay ExecutionLifecycleAttestationReplay) computeDigest() (string, error) {
	return Digest(struct {
		Schema            string
		AttestationDigest string
		ObservationDigest string
		Status            ExecutionLifecycleAttestationReplayStatus
		MissingStage      string
		ObservedAt        time.Time
	}{
		Schema:            replay.Schema,
		AttestationDigest: replay.AttestationDigest,
		ObservationDigest: replay.ObservationDigest,
		Status:            replay.Status,
		MissingStage:      replay.MissingStage,
		ObservedAt:        replay.ObservedAt.UTC(),
	})
}