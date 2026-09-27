package gooo

import (
	"fmt"
	"strings"
)

// LSPRevisionSelfImprovementCycleJEVDeclarativeExecutionPlanBoundaryObservation
// carries a JEV plan boundary into an editor-facing, read-only projection.
type LSPRevisionSelfImprovementCycleJEVDeclarativeExecutionPlanBoundaryObservation struct {
	Status           string
	MissingStage     string
	TaskID           string
	WorkspaceID      string
	GatewayID        string
	ModelID          string
	PolicyDigest     string
	BoundarySignal   string
	BoundaryDigest   string
	ProjectionSignal string
	ProjectionDigest string
	ReadOnly         bool
	NonExecuting     bool
	NonAuthorizing   bool
}

func ObserveLSPRevisionSelfImprovementCycleJEVDeclarativeExecutionPlanBoundary(
	boundary RevisionSelfImprovementCycleJEVDeclarativeExecutionPlanBoundaryObservation,
) (LSPRevisionSelfImprovementCycleJEVDeclarativeExecutionPlanBoundaryObservation, error) {
	result := LSPRevisionSelfImprovementCycleJEVDeclarativeExecutionPlanBoundaryObservation{
		Status:           "UNKNOWN",
		MissingStage:     "lsp-revision-self-improvement-cycle-jev-declarative-execution-plan-boundary",
		TaskID:           boundary.TaskID,
		WorkspaceID:      boundary.WorkspaceID,
		GatewayID:        boundary.GatewayID,
		ModelID:          boundary.ModelID,
		PolicyDigest:     boundary.PolicyDigest,
		BoundarySignal:   boundary.BoundarySignal,
		BoundaryDigest:   boundary.BoundaryDigest,
		ProjectionSignal: "jev-declarative-plan-boundary-lsp-unknown",
		ReadOnly:         true,
		NonExecuting:     true,
		NonAuthorizing:   true,
	}
	setDigest := func() {
		result.ProjectionDigest = digestLSPRevisionSelfImprovementCycleJEVDeclarativeExecutionPlanBoundary(result)
	}
	setDigest()

	if err := boundary.Validate(); err != nil {
		result.MissingStage = "lsp-revision-self-improvement-cycle-jev-declarative-execution-plan-boundary-input"
		setDigest()
		return result, fmt.Errorf("jev declarative execution plan boundary is not valid for lsp projection: %w", err)
	}
	if boundary.Status != "BOUND" || boundary.MissingStage != "" {
		result.MissingStage = "lsp-revision-self-improvement-cycle-jev-declarative-execution-plan-boundary-input-status"
		setDigest()
		return result, fmt.Errorf("jev declarative execution plan boundary is not BOUND")
	}

	result.Status = "BOUND"
	result.MissingStage = ""
	result.ProjectionSignal = "jev-declarative-plan-boundary-lsp-projected"
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "lsp-revision-self-improvement-cycle-jev-declarative-execution-plan-boundary"
		result.ProjectionSignal = "jev-declarative-plan-boundary-lsp-unknown"
		setDigest()
		return result, fmt.Errorf("jev declarative execution plan boundary lsp projection is not valid: %w", err)
	}
	return result, nil
}

func (o LSPRevisionSelfImprovementCycleJEVDeclarativeExecutionPlanBoundaryObservation) Validate() error {
	if o.Status != "BOUND" && o.Status != "UNKNOWN" {
		return fmt.Errorf("jev declarative execution plan boundary lsp status is invalid")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound jev declarative execution plan boundary lsp projection has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown jev declarative execution plan boundary lsp projection has no missing stage")
	}
	for name, value := range map[string]string{
		"task":      o.TaskID,
		"workspace": o.WorkspaceID,
		"gateway":   o.GatewayID,
		"model":     o.ModelID,
	} {
		if strings.TrimSpace(value) == "" || strings.Contains(value, "|") {
			return fmt.Errorf("jev declarative execution plan boundary lsp %s id is invalid", name)
		}
	}
	for name, value := range map[string]string{
		"policy":     o.PolicyDigest,
		"boundary":   o.BoundaryDigest,
		"projection": o.ProjectionDigest,
	} {
		if !validDigest(value) {
			return fmt.Errorf("jev declarative execution plan boundary lsp %s digest is invalid", name)
		}
	}
	if o.BoundarySignal != "jev-declarative-plan-boundary-observed" &&
		o.BoundarySignal != "jev-declarative-plan-boundary-unknown" {
		return fmt.Errorf("jev declarative execution plan boundary lsp boundary signal is invalid")
	}
	if o.ProjectionSignal != "jev-declarative-plan-boundary-lsp-projected" &&
		o.ProjectionSignal != "jev-declarative-plan-boundary-lsp-unknown" {
		return fmt.Errorf("jev declarative execution plan boundary lsp projection signal is invalid")
	}
	if o.Status == "BOUND" && o.ProjectionSignal != "jev-declarative-plan-boundary-lsp-projected" {
		return fmt.Errorf("bound jev declarative execution plan boundary lsp projection has an unknown signal")
	}
	if !o.ReadOnly || !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("jev declarative execution plan boundary lsp projection must remain read-only, non-executing, and non-authorizing")
	}
	if o.ProjectionDigest != digestLSPRevisionSelfImprovementCycleJEVDeclarativeExecutionPlanBoundary(o) {
		return fmt.Errorf("jev declarative execution plan boundary lsp projection digest does not match its fields")
	}
	return nil
}

func digestLSPRevisionSelfImprovementCycleJEVDeclarativeExecutionPlanBoundary(
	observation LSPRevisionSelfImprovementCycleJEVDeclarativeExecutionPlanBoundaryObservation,
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
		observation.BoundaryDigest,
		observation.ProjectionSignal,
		fmt.Sprintf("%t", observation.ReadOnly),
		fmt.Sprintf("%t", observation.NonExecuting),
		fmt.Sprintf("%t", observation.NonAuthorizing),
	}, "|"))
}
