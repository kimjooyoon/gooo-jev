package decision

import (
	"errors"
	"strings"
)

type CounterexampleSet struct {
	AssessmentDigest string
	Observations    []CounterexampleObservation
	NonAuthorizing  bool
	SetDigest       string
}

func BuildCounterexampleSet(
	assessment DecisionAssessment,
	observations []CounterexampleObservation,
) (CounterexampleSet, error) {
	if err := assessment.Validate(); err != nil {
		return CounterexampleSet{}, err
	}
	if len(observations) == 0 {
		return CounterexampleSet{}, errors.New("counterexample set requires observations")
	}
	for _, observation := range observations {
		if err := observation.Validate(); err != nil {
			return CounterexampleSet{}, err
		}
		if observation.AssessmentDigest != assessment.AssessmentDigest {
			return CounterexampleSet{}, errors.New("counterexample assessment digests do not match")
		}
	}
	set := CounterexampleSet{
		AssessmentDigest: assessment.AssessmentDigest,
		Observations:     append([]CounterexampleObservation(nil), observations...),
		NonAuthorizing:   true,
	}
	var err error
	set.SetDigest, err = Digest(set)
	if err != nil {
		return CounterexampleSet{}, err
	}
	return set, nil
}

func (set CounterexampleSet) Validate() error {
	if strings.TrimSpace(set.AssessmentDigest) == "" ||
		len(set.Observations) == 0 ||
		strings.TrimSpace(set.SetDigest) == "" {
		return errors.New("counterexample set is incomplete")
	}
	if !set.NonAuthorizing {
		return errors.New("counterexample set must remain non-authorizing")
	}
	for _, observation := range set.Observations {
		if err := observation.Validate(); err != nil {
			return err
		}
		if observation.AssessmentDigest != set.AssessmentDigest {
			return errors.New("counterexample set contains another assessment")
		}
	}
	copy := set
	copy.SetDigest = ""
	digest, err := Digest(copy)
	if err != nil {
		return err
	}
	if digest != set.SetDigest {
		return errors.New("counterexample set digest does not match its evidence")
	}
	return nil
}
