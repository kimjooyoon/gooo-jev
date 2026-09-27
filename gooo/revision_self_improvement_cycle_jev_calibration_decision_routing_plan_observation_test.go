package gooo

import "testing"

func TestObserveJEVCalibrationDecisionRoutingPlanBindsEvidenceChain(t *testing.T) {
	observation := ObserveJEVCalibrationDecisionRoutingPlan(JEVCalibrationDecisionRoutingPlanObservationInput{
		SourceVersion:             "source-v1",
		ContractVersion:           "contract-v1",
		DeclarationDigest:         "sha256:declaration-v1",
		IRDigest:                  "sha256:ir-v1",
		GeneratedDigest:            "sha256:generated-v1",
		ReverseObservationStatus:  "BOUND",
		ReverseObservationDigest:  "sha256:reverse-v1",
		ProposalStatus:            "PROPOSED",
		PlanDigest:                "sha256:plan-v1",
		ExpectedPlanDigest:        "sha256:plan-v1",
		OutcomeWindowDigest:       "sha256:window-v1",
		OutcomeCount:              3,
	})
	if observation.Status != JEVCalibrationDecisionRoutingPlanObservationBound || observation.ObservationDigest == "" || observation.TargetStage != "plan_application_observation" {
		t.Fatalf("unexpected bound observation: %+v", observation)
	}
	if !observation.IsReadOnly || observation.CanExecute || observation.CanAuthorize || len(observation.Edits) != 0 || observation.Command != "" {
		t.Fatalf("observation crossed authority boundary: %+v", observation)
	}
}

func TestObserveJEVCalibrationDecisionRoutingPlanPreservesMissingIR(t *testing.T) {
	observation := ObserveJEVCalibrationDecisionRoutingPlan(JEVCalibrationDecisionRoutingPlanObservationInput{
		SourceVersion:     "source-v1",
		ContractVersion:   "contract-v1",
		DeclarationDigest: "sha256:declaration-v1",
	})
	if observation.Status != JEVCalibrationDecisionRoutingPlanObservationUnknown || observation.TargetStage != "ir" || observation.ObservationDigest == "" {
		t.Fatalf("missing IR evidence was not preserved: %+v", observation)
	}
}

func TestObserveJEVCalibrationDecisionRoutingPlanDefersReverseProducer(t *testing.T) {
	observation := ObserveJEVCalibrationDecisionRoutingPlan(JEVCalibrationDecisionRoutingPlanObservationInput{
		SourceVersion:             "source-v1",
		ContractVersion:           "contract-v1",
		DeclarationDigest:         "sha256:declaration-v1",
		IRDigest:                  "sha256:ir-v1",
		GeneratedDigest:            "sha256:generated-v1",
		ReverseObservationStatus:  "DEFERRED",
		ReverseObservationDigest:  "sha256:reverse-v1",
		ProposalStatus:            "PROPOSED",
		PlanDigest:                "sha256:plan-v1",
		ExpectedPlanDigest:        "sha256:plan-v1",
		OutcomeWindowDigest:       "sha256:window-v1",
		OutcomeCount:              3,
	})
	if observation.Status != JEVCalibrationDecisionRoutingPlanObservationDeferred || observation.Reason != "reverse observation is deferred" {
		t.Fatalf("deferred reverse producer was not preserved: %+v", observation)
	}
}

func TestObserveJEVCalibrationDecisionRoutingPlanRejectsPlanMismatch(t *testing.T) {
	observation := ObserveJEVCalibrationDecisionRoutingPlan(JEVCalibrationDecisionRoutingPlanObservationInput{
		SourceVersion:             "source-v1",
		ContractVersion:           "contract-v1",
		DeclarationDigest:         "sha256:declaration-v1",
		IRDigest:                  "sha256:ir-v1",
		GeneratedDigest:            "sha256:generated-v1",
		ReverseObservationStatus:  "BOUND",
		ReverseObservationDigest:  "sha256:reverse-v1",
		ProposalStatus:            "PROPOSED",
		PlanDigest:                "sha256:plan-v1",
		ExpectedPlanDigest:        "sha256:plan-v2",
		OutcomeWindowDigest:       "sha256:window-v1",
		OutcomeCount:              3,
	})
	if observation.Status != JEVCalibrationDecisionRoutingPlanObservationUnknown || observation.TargetStage != "plan_digest_match" {
		t.Fatalf("plan digest mismatch was not preserved: %+v", observation)
	}
}

func TestObserveJEVCalibrationDecisionRoutingPlanRequiresPositiveOutcome(t *testing.T) {
	observation := ObserveJEVCalibrationDecisionRoutingPlan(JEVCalibrationDecisionRoutingPlanObservationInput{
		SourceVersion:             "source-v1",
		ContractVersion:           "contract-v1",
		DeclarationDigest:         "sha256:declaration-v1",
		IRDigest:                  "sha256:ir-v1",
		GeneratedDigest:            "sha256:generated-v1",
		ReverseObservationStatus:  "BOUND",
		ReverseObservationDigest:  "sha256:reverse-v1",
		ProposalStatus:            "PROPOSED",
		PlanDigest:                "sha256:plan-v1",
		ExpectedPlanDigest:        "sha256:plan-v1",
		OutcomeWindowDigest:       "sha256:window-v1",
	})
	if observation.Status != JEVCalibrationDecisionRoutingPlanObservationUnknown || observation.TargetStage != "outcome_count" {
		t.Fatalf("empty outcome window was not preserved: %+v", observation)
	}
}