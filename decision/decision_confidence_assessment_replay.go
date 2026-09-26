package decision

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const DecisionConfidenceAssessmentReplaySchemaV1 = "gooo/jev-decision-confidence-assessment-replay/v1"

type DecisionConfidenceAssessmentReplayStatus string

const (
	DecisionConfidenceAssessmentReplayed      DecisionConfidenceAssessmentReplayStatus = "replayed"
	DecisionConfidenceAssessmentReplayUnknown DecisionConfidenceAssessmentReplayStatus = "unknown"
)

type DecisionConfidenceAssessmentReplay struct {
	Schema             string
	AssessmentDigest   string
	DecisionDigest     string
	ResultDigest       string
	Confidence         float64
	MinimumConfidence  float64
	AssessmentStatus   DecisionConfidenceAssessmentStatus
	Status             DecisionConfidenceAssessmentReplayStatus
	MissingStage       string
	ObservedAt         time.Time
	ReplayDigest       string
}

func ReplayDecisionConfidenceAssessment(
	assessment DecisionConfidenceAssessment,
	observedAt time.Time,
) (DecisionConfidenceAssessmentReplay, error) {
	if observedAt.IsZero() {
		return DecisionConfidenceAssessmentReplay{}, errors.New("decision confidence replay observation time is required")
	}
	replay := DecisionConfidenceAssessmentReplay{
		Schema:            DecisionConfidenceAssessmentReplaySchemaV1,
		AssessmentDigest:  assessment.AssessmentDigest,
		DecisionDigest:    assessment.DecisionDigest,
		ResultDigest:      assessment.ResultDigest,
		Confidence:        assessment.Confidence,
		MinimumConfidence: assessment.MinimumConfidence,
		AssessmentStatus:  assessment.Status,
		Status:            DecisionConfidenceAssessmentReplayed,
		ObservedAt:        observedAt.UTC(),
	}
	switch {
	case assessment.Validate() != nil:
		replay.Status = DecisionConfidenceAssessmentReplayUnknown
		replay.MissingStage = "assessment"
	case assessment.Status == DecisionConfidenceAssessmentUnknown:
		replay.Status = DecisionConfidenceAssessmentReplayUnknown
		replay.MissingStage = "assessment-status"
		if strings.TrimSpace(assessment.MissingStage) != "" {
			replay.MissingStage = "assessment:" + assessment.MissingStage
		}
	}
	if err := replay.assignDigest(); err != nil {
		return DecisionConfidenceAssessmentReplay{}, fmt.Errorf("digest decision confidence replay: %w", err)
	}
	if err := replay.Validate(); err != nil {
		return DecisionConfidenceAssessmentReplay{}, err
	}
	return replay, nil
}

func (replay *DecisionConfidenceAssessmentReplay) assignDigest() error {
	digest, err := replay.computeDigest()
	if err != nil {
		return err
	}
	replay.ReplayDigest = digest
	return nil
}

func (replay DecisionConfidenceAssessmentReplay) Validate() error {
	if err := replay.validateShape(); err != nil {
		return err
	}
	expected, err := replay.computeDigest()
	if err != nil {
		return fmt.Errorf("digest decision confidence replay: %w", err)
	}
	if replay.ReplayDigest != expected {
		return errors.New("decision confidence replay digest mismatch")
	}
	return nil
}

func (replay DecisionConfidenceAssessmentReplay) validateShape() error {
	if replay.Schema != DecisionConfidenceAssessmentReplaySchemaV1 {
		return errors.New("unsupported decision confidence replay schema")
	}
	if replay.ObservedAt.IsZero() {
		return errors.New("decision confidence replay observation time is required")
	}
	if strings.TrimSpace(replay.ReplayDigest) == "" {
		return errors.New("decision confidence replay digest is required")
	}
	switch replay.Status {
	case DecisionConfidenceAssessmentReplayed:
		if strings.TrimSpace(replay.AssessmentDigest) == "" ||
			strings.TrimSpace(replay.DecisionDigest) == "" ||
			strings.TrimSpace(replay.ResultDigest) == "" ||
			replay.Confidence < 0 || replay.Confidence > 1 ||
			replay.MinimumConfidence < 0 || replay.MinimumConfidence > 1 ||
			strings.TrimSpace(replay.MissingStage) != "" {
			return errors.New("replayed decision confidence is incomplete")
		}
	case DecisionConfidenceAssessmentReplayUnknown:
		if strings.TrimSpace(replay.MissingStage) == "" {
			return errors.New("unknown decision confidence replay requires a missing stage")
		}
	default:
		return fmt.Errorf("unsupported decision confidence replay status %q", replay.Status)
	}
	return nil
}

func (replay DecisionConfidenceAssessmentReplay) computeDigest() (string, error) {
	return Digest(struct {
		Schema             string
		AssessmentDigest   string
		DecisionDigest     string
		ResultDigest       string
		Confidence         float64
		MinimumConfidence  float64
		AssessmentStatus   DecisionConfidenceAssessmentStatus
		Status             DecisionConfidenceAssessmentReplayStatus
		MissingStage       string
		ObservedAt         time.Time
	}{
		Schema:             replay.Schema,
		AssessmentDigest:   replay.AssessmentDigest,
		DecisionDigest:     replay.DecisionDigest,
		ResultDigest:       replay.ResultDigest,
		Confidence:         replay.Confidence,
		MinimumConfidence:  replay.MinimumConfidence,
		AssessmentStatus:   replay.AssessmentStatus,
		Status:             replay.Status,
		MissingStage:       replay.MissingStage,
		ObservedAt:         replay.ObservedAt.UTC(),
	})
}
