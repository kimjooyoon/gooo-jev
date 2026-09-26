package gooo

import (
	"fmt"
	"strings"
)

// GenerationReceipt records generation and reverse observation without authorizing execution.
type GenerationReceipt struct {
	Status                string
	MissingStage          string
	SourceDigest          string
	GeneratedSourceDigest string
	GeneratedIRDigest     string
	StructureDigest       string
	ExactSourceMatch      bool
	StructureMatch        bool
	NonExecuting          bool
	NonAuthorizing        bool
	GeneratedSource       string
}

// Generate renders a canonical source form and parses it back into IR for observation.
func Generate(document DocumentIR) (GenerationReceipt, error) {
	receipt := GenerationReceipt{
		Status:         "UNKNOWN",
		MissingStage:   "generation",
		SourceDigest:   document.SourceDigest,
		NonExecuting:   true,
		NonAuthorizing: true,
	}
	if err := document.Validate(); err != nil {
		receipt.MissingStage = "generation-validation"
		return receipt, fmt.Errorf("gooo generation: %w", err)
	}

	source := renderCanonical(document)
	observed, err := Parse(source)
	if err != nil {
		receipt.MissingStage = "reverse-observation"
		return receipt, fmt.Errorf("gooo reverse observation: %w", err)
	}

	receipt.Status = "BOUND"
	receipt.MissingStage = ""
	receipt.GeneratedSource = source
	receipt.GeneratedSourceDigest = digestString(source)
	receipt.GeneratedIRDigest = observed.IRDigest
	receipt.StructureDigest = structureDigest(document)
	receipt.StructureMatch = receipt.StructureDigest == structureDigest(observed)
	receipt.ExactSourceMatch = receipt.GeneratedSourceDigest == document.SourceDigest
	if !receipt.StructureMatch {
		receipt.Status = "UNKNOWN"
		receipt.MissingStage = "reverse-observation"
		return receipt, fmt.Errorf("gooo reverse observation: generated structure does not match input IR")
	}
	return receipt, nil
}

func renderCanonical(document DocumentIR) string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "package %s\nnamespace %s\n", document.Package, document.Namespace)
	for _, entity := range document.Entities {
		fmt.Fprintf(&builder, "entity %s id \"%s\"\n", entity.Name, entity.ID)
	}
	for _, activity := range document.Activities {
		fmt.Fprintf(&builder, "activity %s(%s) -> %s\n", activity.Name, activity.Input, activity.Output)
	}
	return builder.String()
}

func structureDigest(document DocumentIR) string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "package=%s\nnamespace=%s\n", document.Package, document.Namespace)
	for _, entity := range document.Entities {
		fmt.Fprintf(&builder, "entity=%s|%s\n", entity.Name, entity.ID)
	}
	for _, activity := range document.Activities {
		fmt.Fprintf(&builder, "activity=%s|%s|%s\n", activity.Name, activity.Input, activity.Output)
	}
	return digestString(builder.String())
}
