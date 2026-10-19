package gooo

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
)

var (
	identifierPattern = regexp.MustCompile("^[A-Za-z_][A-Za-z0-9_]*$")
	packagePattern    = regexp.MustCompile("^package[[:space:]]+([A-Za-z_][A-Za-z0-9_]*)$")
	namespacePattern  = regexp.MustCompile("^namespace[[:space:]]+([A-Za-z_][A-Za-z0-9_]*)$")
	entityPattern     = regexp.MustCompile("^entity[[:space:]]+([A-Za-z_][A-Za-z0-9_]*)[[:space:]]+id[[:space:]]+\"([^\"]+)\"$")
	propertyPattern   = regexp.MustCompile("^property[[:space:]]+([A-Za-z_][A-Za-z0-9_]*)[[:space:]]+([A-Za-z_][A-Za-z0-9_]*)$")
	activityPattern   = regexp.MustCompile("^activity[[:space:]]+([A-Za-z_][A-Za-z0-9_]*)[[:space:]]*\\(([^()]*)\\)[[:space:]]*->[[:space:]]*([A-Za-z_][A-Za-z0-9_]*)$")
)

// Position identifies the source location of a declaration.
type Position struct {
	Line   int
	Column int
}

// EntityDecl is a typed entity declaration in a .gooo document.
type EntityDecl struct {
	Name       string
	ID         string
	Properties []PropertyDecl
	Position   Position
	Digest     string
}

// PropertyDecl keeps an entity property bound to its source location.
type PropertyDecl struct {
	Name     string
	Type     string
	Position Position
}

// ActivityParameterDecl preserves the order and optional name of an input.
type ActivityParameterDecl struct {
	Name     string
	Type     string
	Position Position
}

// ActivityDecl is a non-executing declaration with typed inputs and an entity result.
type ActivityDecl struct {
	Name       string
	Input      string
	Parameters []ActivityParameterDecl
	Output     string
	Position   Position
	Digest     string
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
	currentEntity := -1
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
			currentEntity = -1
			if document.Package != "" {
				document.MissingStage = "duplicate-package"
				return document, ParseError{Stage: document.MissingStage, Line: position.Line, Column: position.Column, Cause: "package is declared more than once"}
			}
			document.Package = packagePattern.FindStringSubmatch(text)[1]
		case namespacePattern.MatchString(text):
			currentEntity = -1
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
			currentEntity = len(document.Entities) - 1
		case propertyPattern.MatchString(text):
			if currentEntity < 0 {
				document.MissingStage = "property-owner"
				return document, ParseError{Stage: document.MissingStage, Line: position.Line, Column: position.Column, Cause: "property must follow its entity"}
			}
			match := propertyPattern.FindStringSubmatch(text)
			entity := &document.Entities[currentEntity]
			for _, property := range entity.Properties {
				if property.Name == match[1] {
					document.MissingStage = "duplicate-property"
					return document, ParseError{Stage: document.MissingStage, Line: position.Line, Column: position.Column, Cause: "property name is declared more than once"}
				}
			}
			entity.Properties = append(entity.Properties, PropertyDecl{Name: match[1], Type: match[2], Position: position})
			entity.Digest = digestEntityDeclaration(*entity)
		case activityPattern.MatchString(text):
			currentEntity = -1
			match := activityPattern.FindStringSubmatch(text)
			parameterPosition := Position{Line: position.Line, Column: position.Column + strings.Index(text, "(") + 1}
			parameters, err := parseActivityParameters(match[2], parameterPosition)
			if err != nil {
				document.MissingStage = "syntax"
				return document, ParseError{Stage: document.MissingStage, Line: position.Line, Column: position.Column, Cause: err.Error()}
			}
			if _, exists := activityNames[match[1]]; exists {
				document.MissingStage = "duplicate-activity"
				return document, ParseError{Stage: document.MissingStage, Line: position.Line, Column: position.Column, Cause: "activity name is declared more than once"}
			}
			activityNames[match[1]] = struct{}{}
			activity := ActivityDecl{Name: match[1], Output: match[3], Position: position, Parameters: parameters}
			if len(parameters) > 0 {
				activity.Input = parameters[0].Type
			}
			// Retain the original one-unnamed-input IR and digests.
			if len(parameters) == 1 && parameters[0].Name == "" {
				activity.Parameters = nil
			}
			activity.Digest = digestActivityDeclaration(activity)
			document.Activities = append(document.Activities, activity)
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
		entityNames[entity.Name] = struct{}{}
		entityIDs[entity.ID] = struct{}{}
	}
	for _, entity := range d.Entities {
		properties := map[string]struct{}{}
		for _, property := range entity.Properties {
			if !identifierPattern.MatchString(property.Name) || !validType(property.Type, entityNames) || property.Position.Line <= entity.Position.Line {
				return fmt.Errorf("entity %q property %q is invalid", entity.Name, property.Name)
			}
			if _, exists := properties[property.Name]; exists {
				return fmt.Errorf("entity %q property %q is duplicated", entity.Name, property.Name)
			}
			properties[property.Name] = struct{}{}
		}
		if expected := digestEntityDeclaration(entity); expected != entity.Digest {
			return fmt.Errorf("entity %q digest does not match its declaration", entity.Name)
		}
	}

	activityNames := map[string]struct{}{}
	for _, activity := range d.Activities {
		if activity.Name == "" || activity.Output == "" || activity.Digest == "" {
			return fmt.Errorf("activity declaration is incomplete")
		}
		if _, exists := activityNames[activity.Name]; exists {
			return fmt.Errorf("activity name %q is duplicated", activity.Name)
		}
		if activity.Parameters == nil {
			if _, exists := entityNames[activity.Input]; !exists {
				return fmt.Errorf("activity %q input entity %q is undefined", activity.Name, activity.Input)
			}
		} else {
			first := ""
			if len(activity.Parameters) > 0 {
				first = activity.Parameters[0].Type
			}
			if activity.Input != first {
				return fmt.Errorf("activity %q first input does not match its parameters", activity.Name)
			}
			parameterNames := map[string]struct{}{}
			for _, parameter := range activity.Parameters {
				if !validType(parameter.Type, entityNames) || (parameter.Name != "" && !identifierPattern.MatchString(parameter.Name)) {
					return fmt.Errorf("activity %q parameter is invalid", activity.Name)
				}
				if parameter.Name != "" {
					if _, exists := parameterNames[parameter.Name]; exists {
						return fmt.Errorf("activity %q parameter %q is duplicated", activity.Name, parameter.Name)
					}
					parameterNames[parameter.Name] = struct{}{}
				}
			}
		}
		if _, exists := entityNames[activity.Output]; !exists {
			return fmt.Errorf("activity %q output entity %q is undefined", activity.Name, activity.Output)
		}
		if expected := digestActivityDeclaration(activity); expected != activity.Digest {
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

func parseActivityParameters(raw string, position Position) ([]ActivityParameterDecl, error) {
	parameters := make([]ActivityParameterDecl, 0)
	if strings.TrimSpace(raw) == "" {
		return parameters, nil
	}
	offset := 0
	for _, item := range strings.Split(raw, ",") {
		fields := strings.Fields(item)
		parameter := ActivityParameterDecl{Position: Position{Line: position.Line, Column: position.Column + offset + len(item) - len(strings.TrimLeft(item, " \t"))}}
		switch len(fields) {
		case 1:
			parameter.Type = fields[0]
		case 2:
			parameter.Name, parameter.Type = fields[0], fields[1]
		default:
			return nil, fmt.Errorf("activity parameter must have a type and optional name")
		}
		if !identifierPattern.MatchString(parameter.Type) || (parameter.Name != "" && !identifierPattern.MatchString(parameter.Name)) {
			return nil, fmt.Errorf("activity parameter has an invalid identifier")
		}
		parameters = append(parameters, parameter)
		offset += len(item) + 1
	}
	return parameters, nil
}

func validType(name string, entities map[string]struct{}) bool {
	if !identifierPattern.MatchString(name) {
		return false
	}
	switch name {
	case "text", "number", "boolean":
		return true
	}
	_, exists := entities[name]
	return exists
}

func digestEntityDeclaration(entity EntityDecl) string {
	base := digestDeclaration("entity", entity.Name, entity.ID, "", entity.Position)
	if len(entity.Properties) == 0 {
		return base
	}
	var builder strings.Builder
	builder.WriteString(base)
	for _, property := range entity.Properties {
		fmt.Fprintf(&builder, "|%s|%s|%d|%d", property.Name, property.Type, property.Position.Line, property.Position.Column)
	}
	return digestString(builder.String())
}

func digestActivityDeclaration(activity ActivityDecl) string {
	base := digestDeclaration("activity", activity.Name, activity.Input, activity.Output, activity.Position)
	if activity.Parameters == nil {
		return base
	}
	var builder strings.Builder
	fmt.Fprintf(&builder, "%s|parameters=%d", base, len(activity.Parameters))
	for _, parameter := range activity.Parameters {
		fmt.Fprintf(&builder, "|%s|%s|%d|%d", parameter.Name, parameter.Type, parameter.Position.Line, parameter.Position.Column)
	}
	return digestString(builder.String())
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
