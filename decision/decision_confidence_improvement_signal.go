package decision

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const DecisionConfidenceImprovementSignalSchemaV1 = "gooo/jev-decision-confidence-improvement-signal/v1"

type DecisionConfidenceImprovementSignalStatus string

const (
	DecisionConfidenceNoChange  DecisionConfidenceImprovementSignalStatus = "no-change"
	DecisionConfidenceCandidate DecisionConfidenceImprovementSignalStatus = "candidate"
	DecisionConfidenceSignalUnknown DecisionConfidenceImprovementSignalStatus = "unknown"
)

type DecisionConfidenceImprovementSignal struct {
	Schema            string
	AssessmentDigest  string
	ObservationDigest string
	ReplayDigest      string
	Status            DecisionConfidenceImprovementSignalStatus
	Action            string
	MissingStage      string
	NonAuthorizing    bool
	ObservedAt        time.Time
	SignalDigest      string
}

func DeriveDecisionConfidenceImprovementSignal(
	assessment DecisionConfidenceAssessment,
	observation DecisionConfidenceObservation,
	replay DecisionConfidenceObservationReplay,
	observedAt time.Time,
) (DecisionConfidenceImprovementSignal, error) {
	if observedAt.IsZero() {
		return DecisionConfidenceImprovementSignal{}, errors.New("decision confidence improvement observation time is required")
	}
	signal := DecisionConfidenceImprovementSignal{
		Schema:            DecisionConfidenceImprovementSignalSchemaV1,
		AssessmentDigest:  assessment.AssessmentDigest,
		ObservationDigest: observation.ObservationDigest,
		ReplayDigest:      replay.ReplayDigest,
		Status:            DecisionConfidenceNoChange,
		Action:            "retain-observed-provenance",
		NonAuthorizing:    true,
		ObservedAt:        observedAt.UTC(),
	}
	switch {
	case assessment.Validate() != nil:
		signal.Status = DecisionConfidenceSignalUnknown
		signal.Action = "repair-assessment-evidence"
		signal.MissingStage = "assessment"
	case observation.Validate() != nil:
		signal.Status = DecisionConfidenceSignalUnknown
		signal.Action = "repair-observation-evidence"
		signal.MissingStage = "observation"
	case replay.Validate() != nil:
		signal.Status = DecisionConfidenceSignalUnknown
		signal.Action = "repair-replay-evidence"
		signal.MissingStage = "replay"
	case assessment.Status == DecisionConfidenceAssessmentReview:
		signal.Status = DecisionConfidenceCandidate
		signal.Action = "request-additional-confidence-evidence"
		signal.MissingStage = assessment.MissingStage
	case observation.Status == DecisionConfidenceObservationUnknown:
		signal.Status = DecisionConfidenceCandidate
		signal.Action = "reobserve-reverse-output"
		signal.MissingStage = observation.MissingStage
	case replay.Status == DecisionConfidenceObservationReplayUnknown:
		signal.Status = DecisionConfidenceCandidate
		signal.Action = "replay-provenance-before-change"
		signal.MissingStage = replay.MissingStage
	}
	if err := signal.assignDigest(); err != nil {
		return DecisionConfidenceImprovementSignal{}, fmt.Errorf("digest decision confidence improvement signal: %w", err)
	}
	if err := signal.Validate(); err != nil {
		return DecisionConfidenceImprovementSignal{}, err
	}
	return signal, nil
}

func (signal *DecisionConfidenceImprovementSignal) assignDigest() error {
	digest, err := signal.computeDigest()
	if err != nil {
		return err
	}
	signal.SignalDigest = digest
	return nil
}

func (signal DecisionConfidenceImprovementSignal) Validate() error {
	if err := signal.validateShape(); err != nil {
		return err
	}
	expected, err := signal.computeDigest()
	if err != nil {
		return fmt.Errorf("digest decision confidence improvement signal: %w", err)
	}
	if signal.SignalDigest != expected {
		return errors.New("decision confidence improvement signal digest mismatch")
	}
	return nil
}

func (signal DecisionConfidenceImprovementSignal) validateShape() error {
	if signal.Schema != DecisionConfidenceImprovementSignalSchemaV1 {
		return errors.New("unsupported decision confidence improvement signal schema")
	}
	if signal.ObservedAt.IsZero() {
		return errors.New("decision confidence improvement observation time is required")
	}
	if strings.TrimSpace(signal.SignalDigest) == "" {
		return errors.New("decision confidence improvement signal digest is required")
	}
	if !signal.NonAuthorizing {
		return errors.New("decision confidence improvement signal must remain non-authorizing")
	}
	if strings.TrimSpace(signal.Action) == "" {
		return errors.New("decision confidence improvement signal action is required")
	}
	switch signal.Status {
	case DecisionConfidenceNoChange:
		if strings.TrimSpace(signal.AssessmentDigest) == "" ||
			strings.TrimSpace(signal.ObservationDigest) == "" ||
			strings.TrimSpace(signal.ReplayDigest) == "" ||
			strings.TrimSpace(signal.MissingStage) != "" {
			return errors.New("no-change decision confidence signal is incomplete")
		}
	case DecisionConfidenceCandidate:
		if strings.TrimSpace(signal.AssessmentDigest) == "" ||
			strings.TrimSpace(signal.ObservationDigest) == "" ||
			strings.TrimSpace(signal.ReplayDigest) == "" ||
			strings.TrimSpace(signal.MissingStage) == "" {
			return errors.New("candidate decision confidence signal is incomplete")
		}
	case DecisionConfidenceSignalUnknown:
		if strings.TrimSpace(signal.MissingStage) == "" {
			return errors.New("unknown decision confidence signal requires a missing stage")
		}
	default:
		return fmt.Errorf("unsupported decision confidence improvement signal status %q", signal.Status)
	}
	return nil
}

func (signal DecisionConfidenceImprovementSignal) computeDigest() (string, error) {
	return Digest(struct {
		Schema            string
		AssessmentDigest  string
		ObservationDigest string
		ReplayDigest      string
		Status            DecisionConfidenceImprovementSignalStatus
		Action            string
		MissingStage      string
		NonAuthorizing    bool
		ObservedAt        time.Time
	}{
		Schema:            signal.Schema,
		AssessmentDigest:  signal.AssessmentDigest,
		ObservationDigest: signal.ObservationDigest,
		ReplayDigest:      signal.ReplayDigest,
		Status:            signal.Status,
		Action:            signal.Action,
		MissingStage:      signal.MissingStage,
		NonAuthorizing:    signal.NonAuthorizing,
		ObservedAt:        signal.ObservedAt.UTC(),
	})
}
