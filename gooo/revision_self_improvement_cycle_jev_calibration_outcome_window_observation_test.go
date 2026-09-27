package gooo

import "testing"

func TestObserveJEVCalibrationOutcomeWindowBindsMetricEvidence(t *testing.T) {
	observation := ObserveJEVCalibrationOutcomeWindow(JEVCalibrationOutcomeWindowObservationInput{
		SourceVersion:           "source-v1",
		ContractVersion:         "contract-v1",
		ChoiceSetDigest:         "sha256:choices-v1",
		ExpectedChoiceSetDigest: "sha256:choices-v1",
		EvidencePrefixDigest:    "sha256:prefix-v1",
		ObservationDigest:       "sha256:observation-v1",
		OutcomeCount:            4,
	})
	if observation.Status != JEVCalibrationOutcomeWindowObservationBound || observation.TargetStage != "outcome_window" || observation.RecordDigest == "" {
		t.Fatalf("unexpected outcome window observation: %+v", observation)
	}
	if !observation.IsReadOnly || observation.CanExecute || observation.CanAuthorize || len(observation.Edits) != 0 || observation.Command != "" {
		t.Fatalf("observation crossed authority boundary: %+v", observation)
	}
}

func TestObserveJEVCalibrationOutcomeWindowPreservesChoiceSetMismatch(t *testing.T) {
	observation := ObserveJEVCalibrationOutcomeWindow(JEVCalibrationOutcomeWindowObservationInput{
		SourceVersion:           "source-v1",
		ContractVersion:         "contract-v1",
		ChoiceSetDigest:         "sha256:choices-v1",
		ExpectedChoiceSetDigest: "sha256:choices-v2",
		EvidencePrefixDigest:    "sha256:prefix-v1",
		ObservationDigest:       "sha256:observation-v1",
		OutcomeCount:            4,
	})
	if observation.Status != JEVCalibrationOutcomeWindowObservationUnknown || observation.TargetStage != "choice_set_match" {
		t.Fatalf("choice-set mismatch was not preserved: %+v", observation)
	}
}

func TestObserveJEVCalibrationOutcomeWindowPreservesDeferredProducer(t *testing.T) {
	observation := ObserveJEVCalibrationOutcomeWindow(JEVCalibrationOutcomeWindowObservationInput{
		SourceVersion:           "source-v1",
		ContractVersion:         "contract-v1",
		ChoiceSetDigest:         "sha256:choices-v1",
		ExpectedChoiceSetDigest: "sha256:choices-v1",
		EvidencePrefixDigest:    "sha256:prefix-v1",
		ObservationDigest:       "sha256:observation-v1",
		ProducerDeferred:        true,
	})
	if observation.Status != JEVCalibrationOutcomeWindowObservationDeferred || observation.Reason != "outcome window producer is deferred" {
		t.Fatalf("deferred producer was not preserved: %+v", observation)
	}
}

func TestObserveJEVCalibrationOutcomeWindowRequiresPositiveCount(t *testing.T) {
	observation := ObserveJEVCalibrationOutcomeWindow(JEVCalibrationOutcomeWindowObservationInput{
		SourceVersion:           "source-v1",
		ContractVersion:         "contract-v1",
		ChoiceSetDigest:         "sha256:choices-v1",
		ExpectedChoiceSetDigest: "sha256:choices-v1",
		EvidencePrefixDigest:    "sha256:prefix-v1",
		ObservationDigest:       "sha256:observation-v1",
	})
	if observation.Status != JEVCalibrationOutcomeWindowObservationUnknown || observation.TargetStage != "outcome_count" {
		t.Fatalf("empty outcome window was not preserved: %+v", observation)
	}
}