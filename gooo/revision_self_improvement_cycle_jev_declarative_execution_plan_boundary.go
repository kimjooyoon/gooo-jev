package gooo

import (
	"fmt"
	"strings"
)

// RevisionSelfImprovementCycleJEVDeclarativeExecutionPlanBoundaryInput describes
// the four declaration boundaries needed to observe a JEV execution plan.
// Observation is deliberately separate from execution, authentication, and
// authorization so the value can be projected into other tooling safely.
type RevisionSelfImprovementCycleJEVDeclarativeExecutionPlanBoundaryInput struct {
	TaskID       string
	WorkspaceID  string
	GatewayID    string
	ModelID      string
	PolicyDigest string
}

// RevisionSelfImprovementCycleJEVDeclarativeExecutionPlanBoundaryObservation
// is the provenance-bearing observation of a declarative JEV execution plan.
type RevisionSelfImprovementCycleJEVDeclarativeExecutionPlanBoundaryObservation struct {
	Status         string
	MissingStage   string
	TaskID         string
	WorkspaceID    string
	GatewayID      string
	ModelID        string
	PolicyDigest   string
	BoundarySignal string
	BoundaryDigest string
	NonExecuting   bool
	NonAuthorizing bool
}

func ObserveRevisionSelfImprovementCycleJEVDeclarativeExecutionPlanBoundary(
	input RevisionSelfImprovementCycleJEVDeclarativeExecutionPlanBoundaryInput,
) (RevisionSelfImprovementCycleJEVDeclarativeExecutionPlanBoundaryObservation, error) {
	result := RevisionSelfImprovementCycleJEVDeclarativeExecutionPlanBoundaryObservation{
		Status:         "UNKNOWN",
		MissingStage:   "revision-self-improvement-cycle-jev-declarative-execution-plan-boundary",
		TaskID:         input.TaskID,
		WorkspaceID:    input.WorkspaceID,
		GatewayID:      input.GatewayID,
		ModelID:        input.ModelID,
		PolicyDigest:   input.PolicyDigest,
		BoundarySignal: "jev-declarative-plan-boundary-unknown",
		NonExecuting:   true,
		NonAuthorizing: true,
	}
	setDigest := func() {
		result.BoundaryDigest = digestRevisionSelfImprovementCycleJEVDeclarativeExecutionPlanBoundary(result)
	}
	setDigest()

	if err := input.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-cycle-jev-declarative-execution-plan-boundary-input"
		setDigest()
		return result, fmt.Errorf("jev declarative execution plan boundary is not valid: %w", err)
	}

	result.Status = "BOUND"
	result.MissingStage = ""
	result.BoundarySignal = "jev-declarative-plan-boundary-observed"
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "revision-self-improvement-cycle-jev-declarative-execution-plan-boundary"
		result.BoundarySignal = "jev-declarative-plan-boundary-unknown"
		setDigest()
		return result, fmt.Errorf("jev declarative execution plan boundary observation is not valid: %w", err)
	}
	return result, nil
}

func (i RevisionSelfImprovementCycleJEVDeclarativeExecutionPlanBoundaryInput) Validate() error {
	for name, value := range map[string]string{
		"task":      i.TaskID,
		"workspace": i.WorkspaceID,
		"gateway":   i.GatewayID,
		"model":     i.ModelID,
	} {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("jev declarative execution plan %s id is empty", name)
		}
		if strings.Contains(value, "|") {
			return fmt.Errorf("jev declarative execution plan %s id contains a digest separator", name)
		}
	}
	if !validDigest(i.PolicyDigest) {
		return fmt.Errorf("jev declarative execution plan policy digest is invalid")
	}
	return nil
}

func (o RevisionSelfImprovementCycleJEVDeclarativeExecutionPlanBoundaryObservation) Validate() error {
	if o.Status != "BOUND" && o.Status != "UNKNOWN" {
		return fmt.Errorf("jev declarative execution plan boundary status is invalid")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound jev declarative execution plan boundary has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown jev declarative execution plan boundary has no missing stage")
	}
	for name, value := range map[string]string{
		"task":      o.TaskID,
		"workspace": o.WorkspaceID,
		"gateway":   o.GatewayID,
		"model":     o.ModelID,
	} {
		if strings.TrimSpace(value) == "" || strings.Contains(value, "|") {
			return fmt.Errorf("jev declarative execution plan boundary %s id is invalid", name)
		}
	}
	if !validDigest(o.PolicyDigest) || !validDigest(o.BoundaryDigest) {
		return fmt.Errorf("jev declarative execution plan boundary digest is invalid")
	}
	if o.BoundarySignal != "jev-declarative-plan-boundary-observed" &&
		o.BoundarySignal != "jev-declarative-plan-boundary-unknown" {
		return fmt.Errorf("jev declarative execution plan boundary signal is invalid")
	}
	if o.Status == "BOUND" && o.BoundarySignal != "jev-declarative-plan-boundary-observed" {
		return fmt.Errorf("bound jev declarative execution plan boundary has an unknown signal")
	}
	if !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("jev declarative execution plan boundary must remain non-executing and non-authorizing")
	}
	if o.BoundaryDigest != digestRevisionSelfImprovementCycleJEVDeclarativeExecutionPlanBoundary(o) {
		return fmt.Errorf("jev declarative execution plan boundary digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementCycleJEVDeclarativeExecutionPlanBoundary(
	observation RevisionSelfImprovementCycleJEVDeclarativeExecutionPlanBoundaryObservation,
) string {
	return digestString(strings.Join([]string{
		observation.Status,
		observation.MissingStage,
		observation.TaskID,
		observation.WorkspaceID,
		observation.GatewayID,
		observation.ModelID,
		observation.PolicyDigest,
		observation.BoundarySignal,
		fmt.Sprintf("%t", observation.NonExecuting),
		fmt.Sprintf("%t", observation.NonAuthorizing),
	}, "|"))
}
