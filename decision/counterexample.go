package decision

import (
	"errors"
	"strings"
	"time"
)

type CounterexampleObservation struct {
	AssessmentDigest   string
	ScenarioDigest     string
	ObservationDigest  string
	EvidenceDigest     string
	RecordedAt         time.Time
	NonAuthorizing     bool
	CounterexampleDigest string
}

func ObserveCounterexample(
	assessment DecisionAssessment,
	scenarioDigest string,
	observationDigest string,
	evidenceDigest string,
	recordedAt time.Time,
) (CounterexampleObservation, error) {
	if err := assessment.Validate(); err != nil {
		return CounterexampleObservation{}, err
	}
	if strings.TrimSpace(scenarioDigest) == "" ||
		strings.TrimSpace(observationDigest) == "" ||
		strings.TrimSpace(evidenceDigest) == "" {
		return CounterexampleObservation{}, errors.New("counterexample digests are required")
	}
	if recordedAt.IsZero() {
		return CounterexampleObservation{}, errors.New("counterexample recorded time is required")
	}
	counterexample := CounterexampleObservation{
		AssessmentDigest:  assessment.AssessmentDigest,
		ScenarioDigest:    scenarioDigest,
		ObservationDigest: observationDigest,
		EvidenceDigest:    evidenceDigest,
		RecordedAt:        recordedAt.UTC(),
		NonAuthorizing:    true,
	}
	var err error
	counterexample.CounterexampleDigest, err = Digest(counterexample)
	if err != nil {
		return CounterexampleObservation{}, err
	}
	return counterexample, nil
}

func (counterexample CounterexampleObservation) Validate() error {
	if strings.TrimSpace(counterexample.AssessmentDigest) == "" ||
		strings.TrimSpace(counterexample.ScenarioDigest) == "" ||
		strings.TrimSpace(counterexample.ObservationDigest) == "" ||
		strings.TrimSpace(counterexample.EvidenceDigest) == "" ||
		strings.TrimSpace(counterexample.CounterexampleDigest) == "" {
		return errors.New("counterexample observation is incomplete")
	}
	if counterexample.RecordedAt.IsZero() {
		return errors.New("counterexample recorded time is required")
	}
	if !counterexample.NonAuthorizing {
		return errors.New("counterexample observation must remain non-authorizing")
	}
	copy := counterexample
	copy.CounterexampleDigest = ""
	digest, err := Digest(copy)
	if err != nil {
		return err
	}
	if digest != counterexample.CounterexampleDigest {
		return errors.New("counterexample digest does not match its evidence")
	}
	return nil
}
