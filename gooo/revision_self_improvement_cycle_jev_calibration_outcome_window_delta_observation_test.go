package gooo

import "testing"

func TestObserveJEVCalibrationOutcomeWindowDeltaPreservesPositiveSignedDelta(t *testing.T) {
	observation := ObserveJEVCalibrationOutcomeWindowDelta(JEVCalibrationOutcomeWindowDeltaObservationInput{
		SourceVersion:        "source-v1",
		ContractVersion:      "contract-v1",
		EvidencePrefixDigest: "sha256:prefix-v1",
		ObservationDigest:    "sha256:observation-v1",
		ExpectedOutcomeCount: 4,
		ObservedOutcomeCount: 7,
	})
	if observation.Status != JEVCalibrationOutcomeWindowDeltaBound || observation.Delta != 3 || observation.RecordDigest == "" {
		t.Fatalf("unexpected positive delta observation: %+v", observation)
	}
	if !observation.IsReadOnly || observation.CanExecute || observation.CanAuthorize || len(observation.Edits) != 0 || observation.Command != "" {
		t.Fatalf("observation crossed authority boundary: %+v", observation)
	}
}

func TestObserveJEVCalibrationOutcomeWindowDeltaPreservesNegativeSignedDelta(t *testing.T) {
	observation := ObserveJEVCalibrationOutcomeWindowDelta(JEVCalibrationOutcomeWindowDeltaObservationInput{
		SourceVersion:        "source-v1",
		ContractVersion:      "contract-v1",
		EvidencePrefixDigest: "sha256:prefix-v1",
		ObservationDigest:    "sha256:observation-v1",
		ExpectedOutcomeCount: 7,
		ObservedOutcomeCount: 4,
	})
	if observation.Status != JEVCalibrationOutcomeWindowDeltaBound || observation.Delta != -3 {
		t.Fatalf("negative delta was misclassified: %+v", observation)
	}
	if observation.Reason != "signed outcome-count delta is recorded without an improvement claim" {
		t.Fatalf("negative delta acquired an implicit improvement claim: %+v", observation)
	}
}

func TestObserveJEVCalibrationOutcomeWindowDeltaPreservesDeferredProducer(t *testing.T) {
	observation := ObserveJEVCalibrationOutcomeWindowDelta(JEVCalibrationOutcomeWindowDeltaObservationInput{
		SourceVersion:        "source-v1",
		ContractVersion:      "contract-v1",
		EvidencePrefixDigest: "sha256:prefix-v1",
		ObservationDigest:    "sha256:observation-v1",
		ExpectedOutcomeCount: 4,
		ObservedOutcomeCount: 7,
		ProducerDeferred:     true,
	})
	if observation.Status != JEVCalibrationOutcomeWindowDeltaDeferred || observation.Reason != "outcome window delta producer is deferred" {
		t.Fatalf("deferred producer was not preserved: %+v", observation)
	}
}

func TestObserveJEVCalibrationOutcomeWindowDeltaPreservesFirstMissingStage(t *testing.T) {
	observation := ObserveJEVCalibrationOutcomeWindowDelta(JEVCalibrationOutcomeWindowDeltaObservationInput{
		SourceVersion:   "source-v1",
		ContractVersion: "contract-v1",
	})
	if observation.Status != JEVCalibrationOutcomeWindowDeltaUnknown || observation.TargetStage != "evidence_prefix" {
		t.Fatalf("first missing stage was not preserved: %+v", observation)
	}
}
