package gooo

import (
	"fmt"
	"strings"
)

const (
	lspRepairActionDeferred = "lsp-repair-action-deferred"
	lspRepairActionComplete = "lsp-repair-action-complete"
	lspRepairActionUnknown = "lsp-repair-action-unknown"
)

type LSPRepairActionProjection struct {
	Status            string
	MissingStage      string
	TargetStage       string
	Title             string
	Kind              string
	SourceDigest      string
	IRDigest          string
	RepairPlanDigest  string
	PlanSignal        string
	ActionSignal      string
	ActionDigest      string
	IsPreferred       bool
	ReadOnly          bool
	NonExecuting      bool
	NonAuthorizing    bool
}

func ObserveLSPRepairActionProjection(
	snapshot LanguageSnapshot,
	plan RevisionSelfImprovementCycleJEVExecutionPlanCoverageLSPRepairPlan,
) LSPRepairActionProjection {
	result := LSPRepairActionProjection{
		Status:         "UNKNOWN",
		MissingStage:   "lsp-repair-action",
		SourceDigest:   snapshot.SourceDigest,
		IRDigest:       snapshot.IRDigest,
		ActionSignal:   lspRepairActionUnknown,
		ReadOnly:       true,
		NonExecuting:   true,
		NonAuthorizing: true,
	}
	setDigest := func() {
		result.ActionDigest = digestLSPRepairActionProjection(result)
	}
	setDigest()

	if err := snapshot.Validate(); err != nil {
		result.MissingStage = "lsp-repair-action-snapshot"
		setDigest()
		return result
	}
	if err := plan.Validate(); err != nil {
		result.MissingStage = "lsp-repair-action-plan"
		setDigest()
		return result
	}

	result.RepairPlanDigest = plan.ObservationDigest
	result.PlanSignal = plan.PlanSignal
	if snapshot.Status == "UNKNOWN" || plan.Status == "UNKNOWN" {
		result.TargetStage = plan.TargetStage
		if result.TargetStage == "" {
			result.TargetStage = snapshot.MissingStage
		}
		if result.TargetStage == "" {
			result.MissingStage = "lsp-repair-action-target"
			setDigest()
			return result
		}
		result.Status = "BOUND"
		result.MissingStage = ""
		result.Title = "Resolve missing provenance stage: " + result.TargetStage
		result.Kind = "quickfix"
		result.ActionSignal = lspRepairActionDeferred
		result.IsPreferred = true
		setDigest()
		return result
	}

	result.Status = "BOUND"
	result.MissingStage = ""
	result.Title = "No repair required"
	result.Kind = "quickfix"
	result.ActionSignal = lspRepairActionComplete
	result.IsPreferred = false
	setDigest()
	return result
}

func (value LSPRepairActionProjection) Validate() error {
	if value.Status != "BOUND" && value.Status != "UNKNOWN" {
		return fmt.Errorf("LSP repair action status is invalid")
	}
	if !value.ReadOnly || !value.NonExecuting || !value.NonAuthorizing {
		return fmt.Errorf("LSP repair action must remain read-only, non-executing, and non-authorizing")
	}
	if value.Status == "UNKNOWN" {
		if value.MissingStage == "" || value.ActionSignal != lspRepairActionUnknown {
			return fmt.Errorf("unknown LSP repair action must preserve its missing stage")
		}
	} else {
		if value.MissingStage != "" || value.SourceDigest == "" || !validDigest(value.RepairPlanDigest) {
			return fmt.Errorf("bound LSP repair action is missing provenance")
		}
		if value.Kind != "quickfix" {
			return fmt.Errorf("bound LSP repair action kind is invalid")
		}
		switch value.ActionSignal {
		case lspRepairActionDeferred:
			if value.TargetStage == "" ||
				!strings.HasPrefix(value.Title, "Resolve missing provenance stage: ") ||
				!value.IsPreferred {
				return fmt.Errorf("deferred LSP repair action is incomplete")
			}
		case lspRepairActionComplete:
			if value.TargetStage != "" || value.Title != "No repair required" || value.IsPreferred {
				return fmt.Errorf("complete LSP repair action is inconsistent")
			}
		default:
			return fmt.Errorf("bound LSP repair action signal is invalid")
		}
	}
	if value.ActionDigest != digestLSPRepairActionProjection(value) {
		return fmt.Errorf("LSP repair action digest does not match its fields")
	}
	return nil
}

func digestLSPRepairActionProjection(value LSPRepairActionProjection) string {
	return digestString(fmt.Sprintf(
		"lsp-repair-action|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%t|%t|%t|%t",
		value.Status,
		value.MissingStage,
		value.TargetStage,
		value.Title,
		value.Kind,
		value.SourceDigest,
		value.IRDigest,
		value.RepairPlanDigest,
		value.PlanSignal,
		value.ActionSignal,
		value.IsPreferred,
		value.ReadOnly,
		value.NonExecuting,
		value.NonAuthorizing,
	))
}