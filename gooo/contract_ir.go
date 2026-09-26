package gooo

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
)

var (
	packagePattern   = regexp.MustCompile("^package[[:space:]]+([A-Za-z_][A-Za-z0-9_]*)$")
	namespacePattern = regexp.MustCompile("^namespace[[:space:]]+([A-Za-z_][A-Za-z0-9_]*)$")
	entityPattern    = regexp.MustCompile("^entity[[:space:]]+([A-Za-z_][A-Za-z0-9_]*)[[:space:]]+id[[:space:]]+\"([^\"]+)\"$")
	activityPattern  = regexp.MustCompile("^activity[[:space:]]+([A-Za-z_][A-Za-z0-9_]*)[[:space:]]*\\([[:space:]]*([A-Za-z_][A-Za-z0-9_]*)[[:space:]]*\\)[[:space:]]*->[[:space:]]*([A-Za-z_][A-Za-z0-9_]*)$")
)

// Position identifies the source location of a declaration.
type Position struct {
	Line   int
	Column int
}

// EntityDecl is a typed entity declaration in a .gooo document.
type EntityDecl struct {
	Name     string
	ID       string
	Position Position
	Digest   string
}

// ActivityDecl is a non-executing relationship between two entities.
type ActivityDecl struct {
	Name     string
	Input    string
	Output   string
	Position Position
	Digest   string
}

// DocumentIR is the deterministic intermediate representation of a .gooo document.
type DocumentIR struct {
	Status         string
	MissingStage   string
	Package        string
	Namespace      string
	Entities       []EntityDecl
	Activities     []ActivityDecl
	SourceDigest   string
	IRDigest       string
	NonExecuting   bool
	NonAuthorizing bool
}

// ParseError preserves the first stage at which a source document became UNKNOWN.
type ParseError struct {
	Stage  string
	Line   int
	Column int
	Cause  string
}

func (e ParseError) Error() string {
	return fmt.Sprintf("gooo parse %s at %d:%d: %s", e.Stage, e.Line, e.Column, e.Cause)
}

// Parse converts the supported .gooo declaration subset into a validated IR.
func Parse(source string) (DocumentIR, error) {
	document := DocumentIR{
		Status:         "UNKNOWN",
		MissingStage:   "parse",
		SourceDigest:   digestString(source),
		NonExecuting:   true,
		NonAuthorizing: true,
	}
	if source == "" {
		document.MissingStage = "empty-source"
		return document, ParseError{Stage: "empty-source", Line: 1, Column: 1, Cause: "source is empty"}
	}

	entityNames := map[string]struct{}{}
	entityIDs := map[string]struct{}{}
	activityNames := map[string]struct{}{}
	for lineIndex, line := range strings.Split(source, "\n") {
		text := strings.TrimSpace(line)
		if text == "" || strings.HasPrefix(text, "#") || strings.HasPrefix(text, "//") {
			continue
		}
		column := strings.Index(line, text) + 1
		if column < 1 {
			column = 1
		}
		position := Position{Line: lineIndex + 1, Column: column}

		switch {
		case packagePattern.MatchString(text):
			if document.Package != "" {
				document.MissingStage = "duplicate-package"
				return document, ParseError{Stage: document.MissingStage, Line: position.Line, Column: position.Column, Cause: "package is declared more than once"}
			}
			document.Package = packagePattern.FindStringSubmatch(text)[1]
		case namespacePattern.MatchString(text):
			if document.Namespace != "" {
				document.MissingStage = "duplicate-namespace"
				return document, ParseError{Stage: document.MissingStage, Line: position.Line, Column: position.Column, Cause: "namespace is declared more than once"}
			}
			document.Namespace = namespacePattern.FindStringSubmatch(text)[1]
		case entityPattern.MatchString(text):
			match := entityPattern.FindStringSubmatch(text)
			if _, exists := entityNames[match[1]]; exists {
				document.MissingStage = "duplicate-entity"
				return document, ParseError{Stage: document.MissingStage, Line: position.Line, Column: position.Column, Cause: "entity name is declared more than once"}
			}
			if _, exists := entityIDs[match[2]]; exists {
				document.MissingStage = "duplicate-entity-id"
				return document, ParseError{Stage: document.MissingStage, Line: position.Line, Column: position.Column, Cause: "entity id is declared more than once"}
			}
			entityNames[match[1]] = struct{}{}
			entityIDs[match[2]] = struct{}{}
			document.Entities = append(document.Entities, EntityDecl{
				Name: match[1], ID: match[2], Position: position,
				Digest: digestDeclaration("entity", match[1], match[2], "", position),
			})
		case activityPattern.MatchString(text):
			match := activityPattern.FindStringSubmatch(text)
			if _, exists := activityNames[match[1]]; exists {
				document.MissingStage = "duplicate-activity"
				return document, ParseError{Stage: document.MissingStage, Line: position.Line, Column: position.Column, Cause: "activity name is declared more than once"}
			}
			activityNames[match[1]] = struct{}{}
			document.Activities = append(document.Activities, ActivityDecl{
				Name: match[1], Input: match[2], Output: match[3], Position: position,
				Digest: digestDeclaration("activity", match[1], match[2], match[3], position),
			})
		default:
			document.MissingStage = "syntax"
			return document, ParseError{Stage: document.MissingStage, Line: position.Line, Column: position.Column, Cause: "unsupported declaration"}
		}
	}

	if document.Package == "" {
		document.MissingStage = "missing-package"
		return document, ParseError{Stage: document.MissingStage, Line: 1, Column: 1, Cause: "package declaration is required"}
	}
	if document.Namespace == "" {
		document.MissingStage = "missing-namespace"
		return document, ParseError{Stage: document.MissingStage, Line: 1, Column: 1, Cause: "namespace declaration is required"}
	}

	document.Status = "BOUND"
	document.MissingStage = ""
	document.IRDigest = digestDocument(document)
	if err := document.Validate(); err != nil {
		document.Status = "UNKNOWN"
		document.MissingStage = "ir-validation"
		return document, ParseError{Stage: document.MissingStage, Line: 1, Column: 1, Cause: err.Error()}
	}
	return document, nil
}

// Validate rejects incomplete, mutated, executing, or authorizing interpretations.
func (d DocumentIR) Validate() error {
	if d.Status != "BOUND" {
		return fmt.Errorf("status must be BOUND")
	}
	if d.MissingStage != "" {
		return fmt.Errorf("missing stage must be empty")
	}
	if d.Package == "" || d.Namespace == "" {
		return fmt.Errorf("package and namespace are required")
	}
	if d.SourceDigest == "" || d.IRDigest == "" {
		return fmt.Errorf("source and IR digests are required")
	}
	if len(d.Entities) == 0 || len(d.Activities) == 0 {
		return fmt.Errorf("at least one entity and activity are required")
	}
	if !d.NonExecuting || !d.NonAuthorizing {
		return fmt.Errorf("IR must remain non-executing and non-authorizing")
	}

	entityNames := map[string]struct{}{}
	entityIDs := map[string]struct{}{}
	for _, entity := range d.Entities {
		if entity.Name == "" || entity.ID == "" || entity.Digest == "" {
			return fmt.Errorf("entity declaration is incomplete")
		}
		if _, exists := entityNames[entity.Name]; exists {
			return fmt.Errorf("entity name %q is duplicated", entity.Name)
		}
		if _, exists := entityIDs[entity.ID]; exists {
			return fmt.Errorf("entity id %q is duplicated", entity.ID)
		}
		if expected := digestDeclaration("entity", entity.Name, entity.ID, "", entity.Position); expected != entity.Digest {
			return fmt.Errorf("entity %q digest does not match its declaration", entity.Name)
		}
		entityNames[entity.Name] = struct{}{}
		entityIDs[entity.ID] = struct{}{}
	}

	activityNames := map[string]struct{}{}
	for _, activity := range d.Activities {
		if activity.Name == "" || activity.Input == "" || activity.Output == "" || activity.Digest == "" {
			return fmt.Errorf("activity declaration is incomplete")
		}
		if _, exists := activityNames[activity.Name]; exists {
			return fmt.Errorf("activity name %q is duplicated", activity.Name)
		}
		if _, exists := entityNames[activity.Input]; !exists {
			return fmt.Errorf("activity %q input entity %q is undefined", activity.Name, activity.Input)
		}
		if _, exists := entityNames[activity.Output]; !exists {
			return fmt.Errorf("activity %q output entity %q is undefined", activity.Name, activity.Output)
		}
		if expected := digestDeclaration("activity", activity.Name, activity.Input, activity.Output, activity.Position); expected != activity.Digest {
			return fmt.Errorf("activity %q digest does not match its declaration", activity.Name)
		}
		activityNames[activity.Name] = struct{}{}
	}
	if expected := digestDocument(d); expected != d.IRDigest {
		return fmt.Errorf("IR digest does not match the document")
	}
	return nil
}

func digestDeclaration(kind, name, first, second string, position Position) string {
	return digestString(fmt.Sprintf("%s|%s|%s|%s|%d|%d", kind, name, first, second, position.Line, position.Column))
}

func digestDocument(d DocumentIR) string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "package=%s\nnamespace=%s\nsource=%s\n", d.Package, d.Namespace, d.SourceDigest)
	for _, entity := range d.Entities {
		fmt.Fprintf(&builder, "entity=%s|%s|%d|%d|%s\n", entity.Name, entity.ID, entity.Position.Line, entity.Position.Column, entity.Digest)
	}
	for _, activity := range d.Activities {
		fmt.Fprintf(&builder, "activity=%s|%s|%s|%d|%d|%s\n", activity.Name, activity.Input, activity.Output, activity.Position.Line, activity.Position.Column, activity.Digest)
	}
	return digestString(builder.String())
}

func digestString(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}
