package gooo

import "testing"

func TestObserveLSPRevisionSelfImprovementCycleJEVDeclarativeExecutionPlanBoundary(t *testing.T) {
	boundary, err := ObserveRevisionSelfImprovementCycleJEVDeclarativeExecutionPlanBoundary(
		RevisionSelfImprovementCycleJEVDeclarativeExecutionPlanBoundaryInput{
			TaskID:       "task:compile",
			WorkspaceID:  "workspace:isolated",
			GatewayID:    "gateway:egress",
			ModelID:      "model:planner",
			PolicyDigest: digestString("policy:v1"),
		},
	)
	if err != nil {
		t.Fatalf("boundary observation returned an error: %v", err)
	}
	projection, err := ObserveLSPRevisionSelfImprovementCycleJEVDeclarativeExecutionPlanBoundary(boundary)
	if err != nil {
		t.Fatalf("lsp projection returned an error: %v", err)
	}
	if projection.Status != "BOUND" || projection.MissingStage != "" {
		t.Fatalf("expected a bound projection, got %#v", projection)
	}
	if projection.ProjectionSignal != "jev-declarative-plan-boundary-lsp-projected" {
		t.Fatalf("unexpected projection signal: %q", projection.ProjectionSignal)
	}
	if err := projection.Validate(); err != nil {
		t.Fatalf("projection did not validate: %v", err)
	}
}

func TestObserveLSPRevisionSelfImprovementCycleJEVDeclarativeExecutionPlanBoundaryRejectsUnknownInput(t *testing.T) {
	boundary := RevisionSelfImprovementCycleJEVDeclarativeExecutionPlanBoundaryObservation{
		Status:         "UNKNOWN",
		MissingStage:   "revision-self-improvement-cycle-jev-declarative-execution-plan-boundary-input",
		TaskID:         "task:compile",
		WorkspaceID:    "workspace:isolated",
		GatewayID:      "gateway:egress",
		ModelID:        "model:planner",
		PolicyDigest:   digestString("policy:v1"),
		BoundarySignal: "jev-declarative-plan-boundary-unknown",
		BoundaryDigest: digestString("boundary"),
		NonExecuting:   true,
		NonAuthorizing: true,
	}
	projection, err := ObserveLSPRevisionSelfImprovementCycleJEVDeclarativeExecutionPlanBoundary(boundary)
	if err == nil {
		t.Fatal("expected unknown boundary input to be rejected")
	}
	if projection.Status != "UNKNOWN" || projection.MissingStage == "" {
		t.Fatalf("expected unknown projection with a missing stage, got %#v", projection)
	}
}

func TestLSPRevisionSelfImprovementCycleJEVDeclarativeExecutionPlanBoundaryRejectsTampering(t *testing.T) {
	boundary, err := ObserveRevisionSelfImprovementCycleJEVDeclarativeExecutionPlanBoundary(
		RevisionSelfImprovementCycleJEVDeclarativeExecutionPlanBoundaryInput{
			TaskID:       "task:compile",
			WorkspaceID:  "workspace:isolated",
			GatewayID:    "gateway:egress",
			ModelID:      "model:planner",
			PolicyDigest: digestString("policy:v1"),
		},
	)
	if err != nil {
		t.Fatalf("boundary observation returned an error: %v", err)
	}
	projection, err := ObserveLSPRevisionSelfImprovementCycleJEVDeclarativeExecutionPlanBoundary(boundary)
	if err != nil {
		t.Fatalf("lsp projection returned an error: %v", err)
	}
	projection.GatewayID = "gateway:tampered"
	if err := projection.Validate(); err == nil {
		t.Fatal("expected tampered lsp projection to be rejected")
	}
}
