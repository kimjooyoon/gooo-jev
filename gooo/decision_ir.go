package gooo

import (
	"fmt"
	"regexp"
	"strings"
)

var decisionPattern = regexp.MustCompile("^decision[[:space:]]+([A-Za-z_][A-Za-z0-9_]*)[[:space:]]+kind[[:space:]]+(choice|score|boolean)[[:space:]]+id[[:space:]]+\"([^\"]+)\"$")

// DecisionKind is the typed output shape used by the Jev decision layer.
type DecisionKind string

const (
	ChoiceDecision  DecisionKind = "choice"
	ScoreDecision   DecisionKind = "score"
	BooleanDecision DecisionKind = "boolean"
)

// DecisionDecl is a source-bound typed decision declaration.
type DecisionDecl struct {
	Name     string
	Kind     DecisionKind
	ID       string
	Position Position
	Digest   string
}

// DecisionDocumentIR extends the base .gooo IR without changing execution authority.
type DecisionDocumentIR struct {
	Status         string
	MissingStage   string
	Base           DocumentIR
	Decisions      []DecisionDecl
	SourceDigest   string
	IRDigest       string
	NonExecuting   bool
	NonAuthorizing bool
}

// ParseDecision parses typed Jev decision declarations alongside base .gooo declarations.
func ParseDecision(source string) (DecisionDocumentIR, error) {
	document := DecisionDocumentIR{
		Status:         "UNKNOWN",
		MissingStage:   "decision-parse",
		SourceDigest:   digestString(source),
		NonExecuting:   true,
		NonAuthorizing: true,
	}
	if source == "" {
		document.MissingStage = "empty-source"
		return document, ParseError{Stage: document.MissingStage, Line: 1, Column: 1, Cause: "source is empty"}
	}

	lines := strings.Split(source, "\n")
	normalized := append([]string(nil), lines...)
	for lineIndex, line := range lines {
		text := strings.TrimSpace(line)
		if text == "" || strings.HasPrefix(text, "#") || strings.HasPrefix(text, "//") {
			continue
		}
		if !strings.HasPrefix(text, "decision") {
			continue
		}

		column := strings.Index(line, text) + 1
		if column < 1 {
			column = 1
		}
		position := Position{Line: lineIndex + 1, Column: column}
		match := decisionPattern.FindStringSubmatch(text)
		if match == nil {
			document.MissingStage = "decision-syntax"
			return document, ParseError{Stage: document.MissingStage, Line: position.Line, Column: position.Column, Cause: "decision must declare kind choice, score, or boolean"}
		}
		document.Decisions = append(document.Decisions, DecisionDecl{
			Name: match[1], Kind: DecisionKind(match[2]), ID: match[3], Position: position,
			Digest: digestDeclaration("decision", match[1], match[2], match[3], position),
		})
		normalized[lineIndex] = "# decision declaration omitted from base parser"
	}

	if len(document.Decisions) == 0 {
		document.MissingStage = "missing-decision"
		return document, ParseError{Stage: document.MissingStage, Line: 1, Column: 1, Cause: "at least one decision declaration is required"}
	}

	base, err := Parse(strings.Join(normalized, "\n"))
	document.Base = base
	if err != nil {
		document.MissingStage = "decision-base-parse"
		return document, ParseError{Stage: document.MissingStage, Line: 1, Column: 1, Cause: err.Error()}
	}

	document.Status = "BOUND"
	document.MissingStage = ""
	document.IRDigest = digestDecisionDocument(document)
	if err := document.Validate(); err != nil {
		document.Status = "UNKNOWN"
		document.MissingStage = "decision-ir-validation"
		return document, ParseError{Stage: document.MissingStage, Line: 1, Column: 1, Cause: err.Error()}
	}
	return document, nil
}

// Validate checks the typed decision extension and its base IR.
func (d DecisionDocumentIR) Validate() error {
	if d.Status != "BOUND" {
		return fmt.Errorf("status must be BOUND")
	}
	if d.MissingStage != "" {
		return fmt.Errorf("missing stage must be empty")
	}
	if err := d.Base.Validate(); err != nil {
		return fmt.Errorf("base IR: %w", err)
	}
	if d.SourceDigest == "" || d.IRDigest == "" {
		return fmt.Errorf("decision source and IR digests are required")
	}
	if len(d.Decisions) == 0 {
		return fmt.Errorf("at least one decision is required")
	}
	if !d.NonExecuting || !d.NonAuthorizing {
		return fmt.Errorf("decision IR must remain non-executing and non-authorizing")
	}

	names := map[string]struct{}{}
	ids := map[string]struct{}{}
	for _, decision := range d.Decisions {
		if decision.Name == "" || decision.ID == "" || decision.Digest == "" {
			return fmt.Errorf("decision declaration is incomplete")
		}
		if decision.Kind != ChoiceDecision && decision.Kind != ScoreDecision && decision.Kind != BooleanDecision {
			return fmt.Errorf("decision %q has unsupported kind %q", decision.Name, decision.Kind)
		}
		if _, exists := names[decision.Name]; exists {
			return fmt.Errorf("decision name %q is duplicated", decision.Name)
		}
		if _, exists := ids[decision.ID]; exists {
			return fmt.Errorf("decision id %q is duplicated", decision.ID)
		}
		expected := digestDeclaration("decision", decision.Name, string(decision.Kind), decision.ID, decision.Position)
		if expected != decision.Digest {
			return fmt.Errorf("decision %q digest does not match its declaration", decision.Name)
		}
		names[decision.Name] = struct{}{}
		ids[decision.ID] = struct{}{}
	}
	if expected := digestDecisionDocument(d); expected != d.IRDigest {
		return fmt.Errorf("decision IR digest does not match the document")
	}
	return nil
}

func digestDecisionDocument(document DecisionDocumentIR) string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "source=%s\nbase=%s\n", document.SourceDigest, document.Base.IRDigest)
	for _, decision := range document.Decisions {
		fmt.Fprintf(&builder, "decision=%s|%s|%s|%d|%d|%s\n", decision.Name, decision.Kind, decision.ID, decision.Position.Line, decision.Position.Column, decision.Digest)
	}
	return digestString(builder.String())
}
