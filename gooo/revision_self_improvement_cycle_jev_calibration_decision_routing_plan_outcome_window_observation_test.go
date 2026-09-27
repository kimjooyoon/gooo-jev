package gooo

import "testing"

func TestObserveJEVCalibrationDecisionRoutingPlanOutcomeWindowBindsChain(t *testing.T) {
	observation := ObserveJEVCalibrationDecisionRoutingPlanOutcomeWindow(JEVCalibrationDecisionRoutingPlanOutcomeWindowObservationInput{
		SourceVersion:                          "source-v1",
		ContractVersion:                        "contract-v1",
		PlanObservationDigest:                  "sha256:plan-v1",
		ExpectedPlanObservationDigest:          "sha256:plan-v1",
		OutcomeWindowObservationDigest:         "sha256:window-v1",
		ExpectedOutcomeWindowObservationDigest: "sha256:window-v1",
	})
	if observation.Status != JEVCalibrationDecisionRoutingPlanOutcomeWindowObservationBound || observation.TargetStage != "plan_outcome_window" || observation.ChainDigest == "" {
		t.Fatalf("unexpected plan-outcome chain: %+v", observation)
	}
	if !observation.IsReadOnly || observation.CanExecute || observation.CanAuthorize || len(observation.Edits) != 0 || observation.Command != "" {
		t.Fatalf("observation crossed authority boundary: %+v", observation)
	}
}

func TestObserveJEVCalibrationDecisionRoutingPlanOutcomeWindowPreservesDeferredProducer(t *testing.T) {
	observation := ObserveJEVCalibrationDecisionRoutingPlanOutcomeWindow(JEVCalibrationDecisionRoutingPlanOutcomeWindowObservationInput{
		SourceVersion:                          "source-v1",
		ContractVersion:                        "contract-v1",
		PlanObservationDigest:                  "sha256:plan-v1",
		ExpectedPlanObservationDigest:          "sha256:plan-v1",
		OutcomeWindowObservationDigest:         "sha256:window-v1",
		ExpectedOutcomeWindowObservationDigest: "sha256:window-v1",
		ProducerDeferred:                       true,
	})
	if observation.Status != JEVCalibrationDecisionRoutingPlanOutcomeWindowObservationDeferred || observation.Reason != "plan-outcome producer is deferred" {
		t.Fatalf("deferred producer was not preserved: %+v", observation)
	}
}

func TestObserveJEVCalibrationDecisionRoutingPlanOutcomeWindowPreservesMismatch(t *testing.T) {
	observation := ObserveJEVCalibrationDecisionRoutingPlanOutcomeWindow(JEVCalibrationDecisionRoutingPlanOutcomeWindowObservationInput{
		SourceVersion:                          "source-v1",
		ContractVersion:                        "contract-v1",
		PlanObservationDigest:                  "sha256:plan-v1",
		ExpectedPlanObservationDigest:          "sha256:plan-v1",
		OutcomeWindowObservationDigest:         "sha256:window-v1",
		ExpectedOutcomeWindowObservationDigest: "sha256:window-v2",
	})
	if observation.Status != JEVCalibrationDecisionRoutingPlanOutcomeWindowObservationUnknown || observation.TargetStage != "outcome_window_observation_digest_match" {
		t.Fatalf("outcome-window mismatch was not preserved: %+v", observation)
	}
}

func TestObserveJEVCalibrationDecisionRoutingPlanOutcomeWindowPreservesFirstMissingStage(t *testing.T) {
	observation := ObserveJEVCalibrationDecisionRoutingPlanOutcomeWindow(JEVCalibrationDecisionRoutingPlanOutcomeWindowObservationInput{
		SourceVersion:         "source-v1",
		ContractVersion:       "contract-v1",
		PlanObservationDigest: "sha256:plan-v1",
	})
	if observation.Status != JEVCalibrationDecisionRoutingPlanOutcomeWindowObservationUnknown || observation.TargetStage != "expected_plan_observation_digest" {
		t.Fatalf("first missing stage was not preserved: %+v", observation)
	}
}
