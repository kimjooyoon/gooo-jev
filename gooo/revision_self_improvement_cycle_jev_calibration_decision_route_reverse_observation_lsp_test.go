package gooo

import "testing"

func TestProjectJEVDecisionRouteReverseObservationLSPReview(t *testing.T) {
	projection := ProjectJEVDecisionRouteReverseObservationLSP(
		JEVDecisionRouteReverseObservationLSPInput{
			Status:               "BOUND",
			SourceVersion:        "src-v1",
			ContractVersion:      "contract-v1",
			RouteStatus:          "REVIEW_CANDIDATE",
			RouteRecordDigest:    "sha256:record",
			DecisionDigest:       "sha256:decision",
			ExpectedRouteStatus:  "REVIEW_CANDIDATE",
			ObservedRouteStatus:  "REVIEW_CANDIDATE",
			OutcomeDigest:        "sha256:outcome",
			TargetStage:          "decision_route_reverse_observation",
			ObservationDigest:    "sha256:observation",
			ReviewRequired:       true,
		},
	)

	if projection.Status != JEVDecisionRouteReverseObservationLSPWarning {
		t.Fatalf("status = %q, want Warning", projection.Status)
	}
	if projection.Code != "jev.route.reverse.review" {
		t.Fatalf("code = %q, want review code", projection.Code)
	}
	if !projection.ReviewRequired || !projection.IsReadOnly || projection.CanEdit || projection.CanExecute || projection.CanAuthorize {
		t.Fatal("review reverse projection crossed a forbidden boundary")
	}
}

func TestProjectJEVDecisionRouteReverseObservationLSPBound(t *testing.T) {
	projection := ProjectJEVDecisionRouteReverseObservationLSP(
		JEVDecisionRouteReverseObservationLSPInput{
			Status:               "BOUND",
			SourceVersion:        "src-v1",
			ContractVersion:      "contract-v1",
			RouteStatus:          "ACCEPT_CANDIDATE",
			RouteRecordDigest:    "sha256:record",
			DecisionDigest:       "sha256:decision",
			ExpectedRouteStatus:  "ACCEPT_CANDIDATE",
			ObservedRouteStatus:  "ACCEPT_CANDIDATE",
			OutcomeDigest:        "sha256:outcome",
			TargetStage:          "decision_route_reverse_observation",
			ObservationDigest:    "sha256:observation",
		},
	)

	if projection.Status != JEVDecisionRouteReverseObservationLSPInformation {
		t.Fatalf("status = %q, want Information", projection.Status)
	}
	if projection.Code != "jev.route.reverse.bound" {
		t.Fatalf("code = %q, want bound code", projection.Code)
	}
}

func TestProjectJEVDecisionRouteReverseObservationLSPDeferred(t *testing.T) {
	projection := ProjectJEVDecisionRouteReverseObservationLSP(
		JEVDecisionRouteReverseObservationLSPInput{
			Status:      "DEFERRED",
			TargetStage: "decision_route_reverse_observation",
		},
	)

	if projection.Status != JEVDecisionRouteReverseObservationLSPWarning {
		t.Fatalf("status = %q, want Warning", projection.Status)
	}
	if projection.Code != "jev.route.reverse.deferred" {
		t.Fatalf("code = %q, want deferred code", projection.Code)
	}
}

func TestProjectJEVDecisionRouteReverseObservationLSPRejectsMismatch(t *testing.T) {
	projection := ProjectJEVDecisionRouteReverseObservationLSP(
		JEVDecisionRouteReverseObservationLSPInput{
			Status:               "BOUND",
			SourceVersion:        "src-v1",
			ContractVersion:      "contract-v1",
			RouteRecordDigest:    "sha256:record",
			DecisionDigest:       "sha256:decision",
			ExpectedRouteStatus:  "ACCEPT_CANDIDATE",
			ObservedRouteStatus:  "REVIEW_CANDIDATE",
			OutcomeDigest:        "sha256:outcome",
			ObservationDigest:    "sha256:observation",
		},
	)

	if projection.Status != JEVDecisionRouteReverseObservationLSPError {
		t.Fatalf("status = %q, want Error", projection.Status)
	}
	if projection.TargetStage != "route_status_match" {
		t.Fatalf("target stage = %q, want route_status_match", projection.TargetStage)
	}
}

func TestProjectJEVDecisionRouteReverseObservationLSPRejectsMissingOutcome(t *testing.T) {
	projection := ProjectJEVDecisionRouteReverseObservationLSP(
		JEVDecisionRouteReverseObservationLSPInput{
			Status:               "BOUND",
			SourceVersion:        "src-v1",
			ContractVersion:      "contract-v1",
			RouteRecordDigest:    "sha256:record",
			DecisionDigest:       "sha256:decision",
			ExpectedRouteStatus:  "ABSTAIN_CANDIDATE",
			ObservedRouteStatus:  "ABSTAIN_CANDIDATE",
			ObservationDigest:    "sha256:observation",
		},
	)

	if projection.Status != JEVDecisionRouteReverseObservationLSPError {
		t.Fatalf("status = %q, want Error", projection.Status)
	}
	if projection.TargetStage != "outcome" {
		t.Fatalf("target stage = %q, want outcome", projection.TargetStage)
	}
}
