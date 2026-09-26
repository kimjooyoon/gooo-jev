package decision

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const DecisionConfidenceObservationReplaySchemaV1 = "gooo/jev-decision-confidence-observation-replay/v1"

type DecisionConfidenceObservationReplayStatus string

const (
	DecisionConfidenceObservationReplayed DecisionConfidenceObservationReplayStatus = "replayed"
	DecisionConfidenceObservationReplayUnknown DecisionConfidenceObservationReplayStatus = "unknown"
)

type DecisionConfidenceObservationReplay struct {
	Schema                   string
	ObservationDigest        string
	AssessmentDigest         string
	ReverseObservationDigest string
	AssessmentStatus         DecisionConfidenceAssessmentStatus
	ReverseStatus            ProvenanceObservationStatus
	ObservationStatus        DecisionConfidenceObservationStatus
	Status                   DecisionConfidenceObservationReplayStatus
	MissingStage             string
	ObservedAt               time.Time
	ReplayDigest             string
}

func ReplayDecisionConfidenceObservation(
	observation DecisionConfidenceObservation,
	observedAt time.Time,
) (DecisionConfidenceObservationReplay, error) {
	if observedAt.IsZero() {
		return DecisionConfidenceObservationReplay{}, errors.New("decision confidence observation replay time is required")
	}
	replay := DecisionConfidenceObservationReplay{
		Schema:                   DecisionConfidenceObservationReplaySchemaV1,
		ObservationDigest:        observation.ObservationDigest,
		AssessmentDigest:         observation.AssessmentDigest,
		ReverseObservationDigest: observation.ReverseObservationDigest,
		AssessmentStatus:         observation.AssessmentStatus,
		ReverseStatus:            observation.ReverseStatus,
		ObservationStatus:        observation.Status,
		Status:                   DecisionConfidenceObservationReplayed,
		ObservedAt:               observedAt.UTC(),
	}
	switch {
	case observation.Validate() != nil:
		replay.Status = DecisionConfidenceObservationReplayUnknown
		replay.MissingStage = "observation"
	case observation.Status == DecisionConfidenceObservationUnknown:
		replay.Status = DecisionConfidenceObservationReplayUnknown
		replay.MissingStage = "observation-status"
		if strings.TrimSpace(observation.MissingStage) != "" {
			replay.MissingStage = "observation:" + observation.MissingStage
		}
	}
	if err := replay.assignDigest(); err != nil {
		return DecisionConfidenceObservationReplay{}, fmt.Errorf("digest decision confidence observation replay: %w", err)
	}
	if err := replay.Validate(); err != nil {
		return DecisionConfidenceObservationReplay{}, err
	}
	return replay, nil
}

func (replay *DecisionConfidenceObservationReplay) assignDigest() error {
	digest, err := replay.computeDigest()
	if err != nil {
		return err
	}
	replay.ReplayDigest = digest
	return nil
}

func (replay DecisionConfidenceObservationReplay) Validate() error {
	if err := replay.validateShape(); err != nil {
		return err
	}
	expected, err := replay.computeDigest()
	if err != nil {
		return fmt.Errorf("digest decision confidence observation replay: %w", err)
	}
	if replay.ReplayDigest != expected {
		return errors.New("decision confidence observation replay digest mismatch")
	}
	return nil
}

func (replay DecisionConfidenceObservationReplay) validateShape() error {
	if replay.Schema != DecisionConfidenceObservationReplaySchemaV1 {
		return errors.New("unsupported decision confidence observation replay schema")
	}
	if replay.ObservedAt.IsZero() {
		return errors.New("decision confidence observation replay time is required")
	}
	if strings.TrimSpace(replay.ReplayDigest) == "" {
		return errors.New("decision confidence observation replay digest is required")
	}
	switch replay.Status {
	case DecisionConfidenceObservationReplayed:
		if strings.TrimSpace(replay.ObservationDigest) == "" ||
			strings.TrimSpace(replay.AssessmentDigest) == "" ||
			strings.TrimSpace(replay.ReverseObservationDigest) == "" ||
			replay.AssessmentStatus != DecisionConfidenceAssessed ||
			replay.ReverseStatus != ProvenanceObserved ||
			replay.ObservationStatus != DecisionConfidenceObserved ||
			strings.TrimSpace(replay.MissingStage) != "" {
			return errors.New("replayed decision confidence observation is incomplete")
		}
	case DecisionConfidenceObservationReplayUnknown:
		if strings.TrimSpace(replay.MissingStage) == "" {
			return errors.New("unknown decision confidence observation replay requires a missing stage")
		}
	default:
		return fmt.Errorf("unsupported decision confidence observation replay status %q", replay.Status)
	}
	return nil
}

func (replay DecisionConfidenceObservationReplay) computeDigest() (string, error) {
	return Digest(struct {
		Schema                   string
		ObservationDigest        string
		AssessmentDigest         string
		ReverseObservationDigest string
		AssessmentStatus         DecisionConfidenceAssessmentStatus
		ReverseStatus            ProvenanceObservationStatus
		ObservationStatus        DecisionConfidenceObservationStatus
		Status                   DecisionConfidenceObservationReplayStatus
		MissingStage             string
		ObservedAt               time.Time
	}{
		Schema:                   replay.Schema,
		ObservationDigest:        replay.ObservationDigest,
		AssessmentDigest:         replay.AssessmentDigest,
		ReverseObservationDigest: replay.ReverseObservationDigest,
		AssessmentStatus:         replay.AssessmentStatus,
		ReverseStatus:            replay.ReverseStatus,
		ObservationStatus:        replay.ObservationStatus,
		Status:                   replay.Status,
		MissingStage:             replay.MissingStage,
		ObservedAt:               replay.ObservedAt.UTC(),
	})
}
