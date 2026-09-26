package decision

import (
	"fmt"
	"strings"
)

// GoooEvidenceDeclarationIR is an optional provenance declaration that does
// not alter the legacy GoooDeclarationIR digest.
type GoooEvidenceDeclarationIR struct {
	Name   string
	Source string
}

// GoooEvidenceDeclarationIRGenerationInput supplies extended .gooo source.
type GoooEvidenceDeclarationIRGenerationInput struct {
	SourceText     string
	NonAuthorizing bool
}

// GoooEvidenceDeclarationIRGeneration preserves legacy IR plus evidence
// declarations in a separate, deterministic generation path.
type GoooEvidenceDeclarationIRGeneration struct {
	Status           string
	BaseIR           GoooDeclarationIR
	Evidence         []GoooEvidenceDeclarationIR
	SourceDigest     string
	BaseIRDigest     string
	EvidenceDigest   string
	GeneratedSource  string
	GenerationDigest string
	MissingStage     string
	NonExecuting     bool
	NonAuthorizing   bool
}

// DeriveGoooEvidenceDeclarationIRGeneration parses evidence declarations,
// reuses the legacy parser, and never executes an activity.
func DeriveGoooEvidenceDeclarationIRGeneration(input GoooEvidenceDeclarationIRGenerationInput) GoooEvidenceDeclarationIRGeneration {
	output := GoooEvidenceDeclarationIRGeneration{
		Status:         "UNKNOWN",
		NonExecuting:   true,
		NonAuthorizing: true,
	}
	if !input.NonAuthorizing {
		output.NonAuthorizing = false
		output.MissingStage = "authorization-boundary"
		return output
	}
	if strings.TrimSpace(input.SourceText) == "" {
		output.MissingStage = "declaration-source"
		return output
	}
	baseSource, evidence, err := parseGoooEvidenceDeclarationSource(input.SourceText)
	if err != nil {
		output.MissingStage = "syntax"
		return output
	}
	base := DeriveGoooDeclarationIRGeneration(GoooDeclarationIRGenerationInput{
		SourceText:     baseSource,
		NonAuthorizing: true,
	})
	if base.Status != "ready" {
		output.MissingStage = base.MissingStage
		if strings.TrimSpace(output.MissingStage) == "" {
			output.MissingStage = "legacy-declaration"
		}
		return output
	}
	evidenceDigest, err := Digest(struct {
		Evidence []GoooEvidenceDeclarationIR
	}{Evidence: evidence})
	if err != nil {
		output.MissingStage = "evidence-digest"
		return output
	}
	generated, err := GenerateGoooEvidenceDeclaration(base.IR, evidence)
	if err != nil {
		output.MissingStage = "generation"
		return output
	}
	sourceDigest, err := Digest(struct {
		SourceText string
	}{SourceText: input.SourceText})
	if err != nil {
		output.MissingStage = "declaration-source-digest"
		return output
	}
	generationDigest, err := Digest(struct {
		SourceText string
	}{SourceText: generated})
	if err != nil {
		output.MissingStage = "generation-digest"
		return output
	}
	output.Status = "ready"
	output.BaseIR = base.IR
	output.Evidence = evidence
	output.SourceDigest = sourceDigest
	output.BaseIRDigest = base.IRDigest
	output.EvidenceDigest = evidenceDigest
	output.GeneratedSource = generated
	output.GenerationDigest = generationDigest
	return output
}

// GenerateGoooEvidenceDeclaration emits legacy canonical source followed by
// canonical evidence declarations.
func GenerateGoooEvidenceDeclaration(ir GoooDeclarationIR, evidence []GoooEvidenceDeclarationIR) (string, error) {
	base, err := GenerateGoooDeclaration(ir)
	if err != nil {
		return "", err
	}
	if len(evidence) == 0 {
		return base, nil
	}
	var builder strings.Builder
	builder.WriteString(base)
	if !strings.HasSuffix(base, "\n") {
		builder.WriteString("\n")
	}
	for _, declaration := range evidence {
		if strings.TrimSpace(declaration.Name) == "" || strings.TrimSpace(declaration.Source) == "" {
			return "", fmt.Errorf("evidence name and source are required")
		}
		fmt.Fprintf(&builder, "evidence %s from %s\n", declaration.Name, declaration.Source)
	}
	return builder.String(), nil
}

func parseGoooEvidenceDeclarationSource(source string) (string, []GoooEvidenceDeclarationIR, error) {
	normalized := strings.ReplaceAll(source, "\r\n", "\n")
	lines := strings.Split(normalized, "\n")
	baseLines := make([]string, 0, len(lines))
	evidence := make([]GoooEvidenceDeclarationIR, 0)
	seen := make(map[string]struct{})
	for lineNumber, rawLine := range lines {
		line := strings.TrimSpace(rawLine)
		if line == "" || strings.HasPrefix(line, "//") || strings.HasPrefix(line, "#") {
			baseLines = append(baseLines, rawLine)
			continue
		}
		fields := strings.Fields(line)
		if fields[0] != "evidence" {
			baseLines = append(baseLines, rawLine)
			continue
		}
		if len(fields) != 4 || fields[2] != "from" ||
			strings.TrimSpace(fields[1]) == "" || strings.TrimSpace(fields[3]) == "" {
			return "", nil, fmt.Errorf("invalid evidence at line %d", lineNumber+1)
		}
		if _, exists := seen[fields[1]]; exists {
			return "", nil, fmt.Errorf("duplicate evidence at line %d", lineNumber+1)
		}
		seen[fields[1]] = struct{}{}
		evidence = append(evidence, GoooEvidenceDeclarationIR{
			Name:   fields[1],
			Source: fields[3],
		})
	}
	return strings.Join(baseLines, "\n"), evidence, nil
}
