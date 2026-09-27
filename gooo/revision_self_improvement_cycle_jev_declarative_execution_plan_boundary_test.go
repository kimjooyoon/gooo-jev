package gooo

import "testing"

func TestObserveRevisionSelfImprovementCycleJEVDeclarativeExecutionPlanBoundary(t *testing.T) {
	input := RevisionSelfImprovementCycleJEVDeclarativeExecutionPlanBoundaryInput{
		TaskID:       "task:compile",
		WorkspaceID:  "workspace:isolated",
		GatewayID:    "gateway:egress",
		ModelID:      "model:planner",
		PolicyDigest: digestString("policy:v1"),
	}

	first, err := ObserveRevisionSelfImprovementCycleJEVDeclarativeExecutionPlanBoundary(input)
	if err != nil {
		t.Fatalf("observe returned an error: %v", err)
	}
	second, err := ObserveRevisionSelfImprovementCycleJEVDeclarativeExecutionPlanBoundary(input)
	if err != nil {
		t.Fatalf("second observe returned an error: %v", err)
	}
	if first.Status != "BOUND" || first.MissingStage != "" {
		t.Fatalf("expected a bound observation, got %#v", first)
	}
	if first.BoundaryDigest != second.BoundaryDigest {
		t.Fatalf("expected deterministic boundary digest, got %q and %q", first.BoundaryDigest, second.BoundaryDigest)
	}
	if err := first.Validate(); err != nil {
		t.Fatalf("bound observation did not validate: %v", err)
	}
}

func TestObserveRevisionSelfImprovementCycleJEVDeclarativeExecutionPlanBoundaryUnknown(t *testing.T) {
	input := RevisionSelfImprovementCycleJEVDeclarativeExecutionPlanBoundaryInput{
		WorkspaceID:  "workspace:isolated",
		GatewayID:    "gateway:egress",
		ModelID:      "model:planner",
		PolicyDigest: digestString("policy:v1"),
	}
	observation, err := ObserveRevisionSelfImprovementCycleJEVDeclarativeExecutionPlanBoundary(input)
	if err == nil {
		t.Fatal("expected invalid input to return an error")
	}
	if observation.Status != "UNKNOWN" || observation.MissingStage == "" {
		t.Fatalf("expected unknown observation with a missing stage, got %#v", observation)
	}
}

func TestRevisionSelfImprovementCycleJEVDeclarativeExecutionPlanBoundaryRejectsTampering(t *testing.T) {
	input := RevisionSelfImprovementCycleJEVDeclarativeExecutionPlanBoundaryInput{
		TaskID:       "task:compile",
		WorkspaceID:  "workspace:isolated",
		GatewayID:    "gateway:egress",
		ModelID:      "model:planner",
		PolicyDigest: digestString("policy:v1"),
	}
	observation, err := ObserveRevisionSelfImprovementCycleJEVDeclarativeExecutionPlanBoundary(input)
	if err != nil {
		t.Fatalf("observe returned an error: %v", err)
	}
	observation.ModelID = "model:tampered"
	if err := observation.Validate(); err == nil {
		t.Fatal("expected tampered observation to be rejected")
	}
}
