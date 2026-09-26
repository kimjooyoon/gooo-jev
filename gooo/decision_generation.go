package gooo

import (
	"fmt"
	"strings"
)

// DecisionGenerationReceipt records typed-decision generation and reverse observation.
type DecisionGenerationReceipt struct {
	Status                string
	MissingStage          string
	SourceDigest          string
	GeneratedSourceDigest string
	GeneratedIRDigest     string
	DecisionIRDigest      string
	StructureDigest       string
	ExactSourceMatch      bool
	StructureMatch        bool
	NonExecuting          bool
	NonAuthorizing        bool
	GeneratedSource       string
}

// GenerateDecision renders a canonical typed-decision source form and parses it back.
func GenerateDecision(document DecisionDocumentIR) (DecisionGenerationReceipt, error) {
	receipt := DecisionGenerationReceipt{
		Status:           "UNKNOWN",
		MissingStage:     "decision-generation",
		SourceDigest:     document.SourceDigest,
		DecisionIRDigest: document.IRDigest,
		NonExecuting:     true,
		NonAuthorizing:   true,
	}
	if err := document.Validate(); err != nil {
		receipt.MissingStage = "decision-generation-validation"
		return receipt, fmt.Errorf("gooo decision generation: %w", err)
	}

	source := renderDecisionCanonical(document)
	observed, err := ParseDecision(source)
	if err != nil {
		receipt.MissingStage = "decision-reverse-observation"
		return receipt, fmt.Errorf("gooo decision reverse observation: %w", err)
	}

	receipt.Status = "BOUND"
	receipt.MissingStage = ""
	receipt.GeneratedSource = source
	receipt.GeneratedSourceDigest = digestString(source)
	receipt.GeneratedIRDigest = observed.IRDigest
	receipt.StructureDigest = decisionStructureDigest(document)
	receipt.StructureMatch = receipt.StructureDigest == decisionStructureDigest(observed)
	receipt.ExactSourceMatch = receipt.GeneratedSourceDigest == document.SourceDigest
	if !receipt.StructureMatch {
		receipt.Status = "UNKNOWN"
		receipt.MissingStage = "decision-reverse-observation"
		return receipt, fmt.Errorf("gooo decision reverse observation: generated structure does not match input IR")
	}
	return receipt, nil
}

func renderDecisionCanonical(document DecisionDocumentIR) string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "package %s\nnamespace %s\n", document.Base.Package, document.Base.Namespace)
	for _, entity := range document.Base.Entities {
		fmt.Fprintf(&builder, "entity %s id \"%s\"\n", entity.Name, entity.ID)
	}
	for _, decision := range document.Decisions {
		fmt.Fprintf(&builder, "decision %s kind %s id \"%s\"\n", decision.Name, decision.Kind, decision.ID)
	}
	for _, activity := range document.Base.Activities {
		fmt.Fprintf(&builder, "activity %s(%s) -> %s\n", activity.Name, activity.Input, activity.Output)
	}
	return builder.String()
}

func decisionStructureDigest(document DecisionDocumentIR) string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "package=%s\nnamespace=%s\n", document.Base.Package, document.Base.Namespace)
	for _, entity := range document.Base.Entities {
		fmt.Fprintf(&builder, "entity=%s|%s\n", entity.Name, entity.ID)
	}
	for _, decision := range document.Decisions {
		fmt.Fprintf(&builder, "decision=%s|%s|%s\n", decision.Name, decision.Kind, decision.ID)
	}
	for _, activity := range document.Base.Activities {
		fmt.Fprintf(&builder, "activity=%s|%s|%s\n", activity.Name, activity.Input, activity.Output)
	}
	return digestString(builder.String())
}
