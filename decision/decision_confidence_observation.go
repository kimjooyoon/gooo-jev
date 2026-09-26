package decision

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const DecisionConfidenceObservationSchemaV1 = "gooo/jev-decision-confidence-observation/v1"

type DecisionConfidenceObservationStatus string

const (
	DecisionConfidenceObserved DecisionConfidenceObservationStatus = "observed"
	DecisionConfidenceObservationUnknown DecisionConfidenceObservationStatus = "unknown"
)

type DecisionConfidenceObservation struct {
	Schema                   string
	AssessmentDigest         string
	ReverseObservationDigest string
	AssessmentStatus         DecisionConfidenceAssessmentStatus
	ReverseStatus            ProvenanceObservationStatus
	Status                   DecisionConfidenceObservationStatus
	MissingStage             string
	ObservedAt               time.Time
	ObservationDigest        string
}

func ObserveDecisionConfidence(
	assessment DecisionConfidenceAssessment,
	reverse ReverseObservation,
	observedAt time.Time,
) (DecisionConfidenceObservation, error) {
	if observedAt.IsZero() {
		return DecisionConfidenceObservation{}, errors.New("decision confidence observation time is required")
	}
	observation := DecisionConfidenceObservation{
		Schema:                   DecisionConfidenceObservationSchemaV1,
		AssessmentDigest:         assessment.AssessmentDigest,
		ReverseObservationDigest: reverse.ObservationDigest,
		AssessmentStatus:         assessment.Status,
		ReverseStatus:            reverse.Status,
		Status:                   DecisionConfidenceObserved,
		ObservedAt:               observedAt.UTC(),
	}
	switch {
	case assessment.Validate() != nil:
		observation.Status = DecisionConfidenceObservationUnknown
		observation.MissingStage = "assessment"
	case assessment.Status != DecisionConfidenceAssessed:
		observation.Status = DecisionConfidenceObservationUnknown
		observation.MissingStage = "assessment-status"
		if strings.TrimSpace(assessment.MissingStage) != "" {
			observation.MissingStage = "assessment:" + assessment.MissingStage
		}
	case reverse.Validate() != nil:
		observation.Status = DecisionConfidenceObservationUnknown
		observation.MissingStage = "reverse-observation"
	case reverse.Status != ProvenanceObserved:
		observation.Status = DecisionConfidenceObservationUnknown
		observation.MissingStage = "reverse-observation-status"
		if strings.TrimSpace(reverse.MissingStage) != "" {
			observation.MissingStage = "reverse-observation:" + reverse.MissingStage
		}
	}
	if err := observation.assignDigest(); err != nil {
		return DecisionConfidenceObservation{}, fmt.Errorf("digest decision confidence observation: %w", err)
	}
	if err := observation.Validate(); err != nil {
		return DecisionConfidenceObservation{}, err
	}
	return observation, nil
}

func (observation *DecisionConfidenceObservation) assignDigest() error {
	digest, err := observation.computeDigest()
	if err != nil {
		return err
	}
	observation.ObservationDigest = digest
	return nil
}

func (observation DecisionConfidenceObservation) Validate() error {
	if err := observation.validateShape(); err != nil {
		return err
	}
	expected, err := observation.computeDigest()
	if err != nil {
		return fmt.Errorf("digest decision confidence observation: %w", err)
	}
	if observation.ObservationDigest != expected {
		return errors.New("decision confidence observation digest mismatch")
	}
	return nil
}

func (observation DecisionConfidenceObservation) validateShape() error {
	if observation.Schema != DecisionConfidenceObservationSchemaV1 {
		return errors.New("unsupported decision confidence observation schema")
	}
	if observation.ObservedAt.IsZero() {
		return errors.New("decision confidence observation time is required")
	}
	if strings.TrimSpace(observation.ObservationDigest) == "" {
		return errors.New("decision confidence observation digest is required")
	}
	switch observation.Status {
	case DecisionConfidenceObserved:
		if strings.TrimSpace(observation.AssessmentDigest) == "" ||
			strings.TrimSpace(observation.ReverseObservationDigest) == "" ||
			observation.AssessmentStatus != DecisionConfidenceAssessed ||
			observation.ReverseStatus != ProvenanceObserved ||
			strings.TrimSpace(observation.MissingStage) != "" {
			return errors.New("observed decision confidence is incomplete")
		}
	case DecisionConfidenceObservationUnknown:
		if strings.TrimSpace(observation.MissingStage) == "" {
			return errors.New("unknown decision confidence observation requires a missing stage")
		}
	default:
		return fmt.Errorf("unsupported decision confidence observation status %q", observation.Status)
	}
	return nil
}

func (observation DecisionConfidenceObservation) computeDigest() (string, error) {
	return Digest(struct {
		Schema                   string
		AssessmentDigest         string
		ReverseObservationDigest string
		AssessmentStatus         DecisionConfidenceAssessmentStatus
		ReverseStatus            ProvenanceObservationStatus
		Status                   DecisionConfidenceObservationStatus
		MissingStage             string
		ObservedAt               time.Time
	}{
		Schema:                   observation.Schema,
		AssessmentDigest:         observation.AssessmentDigest,
		ReverseObservationDigest: observation.ReverseObservationDigest,
		AssessmentStatus:         observation.AssessmentStatus,
		ReverseStatus:            observation.ReverseStatus,
		Status:                   observation.Status,
		MissingStage:             observation.MissingStage,
		ObservedAt:               observation.ObservedAt.UTC(),
	})
}
