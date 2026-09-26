package decision

import (
	"fmt"
	"strconv"
	"strings"
)

// GoooEntityIR is the deterministic IR for one .gooo entity declaration.
type GoooEntityIR struct {
	Name       string
	ID         string
	Properties []GoooPropertyIR
}

// GoooPropertyIR is the deterministic IR for one entity property.
type GoooPropertyIR struct {
	Name string
	Type string
}

// GoooActivityParameterIR is the deterministic IR for one activity parameter.
type GoooActivityParameterIR struct {
	Name string
	Type string
}

// GoooActivityIR is the deterministic IR for one .gooo activity declaration.
type GoooActivityIR struct {
	Name       string
	Parameters []GoooActivityParameterIR
	ReturnType string
}

// GoooDeclarationIR is the normalized declaration representation used for
// generation and replay without executing any activity.
type GoooDeclarationIR struct {
	Package    string
	Namespace  string
	Entities   []GoooEntityIR
	Activities []GoooActivityIR
}

// GoooDeclarationIRGenerationInput supplies exact .gooo source for parsing.
type GoooDeclarationIRGenerationInput struct {
	SourceText     string
	NonAuthorizing bool
}

// GoooDeclarationIRGeneration is the source-to-IR-to-generation result. Every
// digest is evidence only and cannot authorize execution or mutation.
type GoooDeclarationIRGeneration struct {
	Status           string
	IR               GoooDeclarationIR
	SourceDigest     string
	IRDigest         string
	GeneratedSource  string
	GenerationDigest string
	MissingStage     string
	NonExecuting     bool
	NonAuthorizing   bool
}

// DeriveGoooDeclarationIRGeneration parses the supported .gooo declaration
// grammar, computes a deterministic IR, and generates canonical source.
func DeriveGoooDeclarationIRGeneration(input GoooDeclarationIRGenerationInput) GoooDeclarationIRGeneration {
	output := GoooDeclarationIRGeneration{
		Status: "UNKNOWN", NonExecuting: true, NonAuthorizing: true,
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
	ir, err := parseGoooDeclarationIR(input.SourceText)
	if err != nil {
		output.MissingStage = "syntax"
		return output
	}
	sourceDigest, err := Digest(struct {
		SourceText string
	}{SourceText: input.SourceText})
	if err != nil {
		output.MissingStage = "declaration-source-digest"
		return output
	}
	irDigest, err := Digest(ir)
	if err != nil {
		output.MissingStage = "ir-digest"
		return output
	}
	generated, err := GenerateGoooDeclaration(ir)
	if err != nil {
		output.MissingStage = "generation"
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
	output.IR = ir
	output.SourceDigest = sourceDigest
	output.IRDigest = irDigest
	output.GeneratedSource = generated
	output.GenerationDigest = generationDigest
	return output
}

// GenerateGoooDeclaration emits canonical source for a parsed declaration.
func GenerateGoooDeclaration(ir GoooDeclarationIR) (string, error) {
	if strings.TrimSpace(ir.Package) == "" || strings.TrimSpace(ir.Namespace) == "" {
		return "", fmt.Errorf("package and namespace are required")
	}
	var builder strings.Builder
	fmt.Fprintf(&builder, "package %s\nnamespace %s\n\n", ir.Package, ir.Namespace)
	for _, entity := range ir.Entities {
		if strings.TrimSpace(entity.Name) == "" || strings.TrimSpace(entity.ID) == "" {
			return "", fmt.Errorf("entity name and id are required")
		}
		fmt.Fprintf(&builder, "entity %s id %s\n", entity.Name, strconv.Quote(entity.ID))
		for _, property := range entity.Properties {
			if strings.TrimSpace(property.Name) == "" || strings.TrimSpace(property.Type) == "" {
				return "", fmt.Errorf("property name and type are required")
			}
			fmt.Fprintf(&builder, "  property %s %s\n", property.Name, property.Type)
		}
		builder.WriteString("\n")
	}
	for _, activity := range ir.Activities {
		if strings.TrimSpace(activity.Name) == "" || strings.TrimSpace(activity.ReturnType) == "" {
			return "", fmt.Errorf("activity name and return type are required")
		}
		parameters := make([]string, 0, len(activity.Parameters))
		for _, parameter := range activity.Parameters {
			if strings.TrimSpace(parameter.Type) == "" {
				return "", fmt.Errorf("activity parameter type is required")
			}
			if strings.TrimSpace(parameter.Name) == "" {
				parameters = append(parameters, parameter.Type)
			} else {
				parameters = append(parameters, parameter.Name+" "+parameter.Type)
			}
		}
		fmt.Fprintf(&builder, "activity %s(%s) -> %s\n", activity.Name, strings.Join(parameters, ", "), activity.ReturnType)
	}
	return builder.String(), nil
}

func parseGoooDeclarationIR(source string) (GoooDeclarationIR, error) {
	normalized := strings.ReplaceAll(source, "\r\n", "\n")
	lines := strings.Split(normalized, "\n")
	var ir GoooDeclarationIR
	currentEntity := -1
	seenPackage := false
	seenNamespace := false
	for lineNumber, rawLine := range lines {
		line := strings.TrimSpace(rawLine)
		if line == "" || strings.HasPrefix(line, "//") || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		switch fields[0] {
		case "package":
			if len(fields) != 2 || seenPackage {
				return GoooDeclarationIR{}, fmt.Errorf("invalid package at line %d", lineNumber+1)
			}
			ir.Package = fields[1]
			seenPackage = true
			currentEntity = -1
		case "namespace":
			if len(fields) != 2 || seenNamespace {
				return GoooDeclarationIR{}, fmt.Errorf("invalid namespace at line %d", lineNumber+1)
			}
			ir.Namespace = fields[1]
			seenNamespace = true
			currentEntity = -1
		case "entity":
			entity, err := parseGoooEntity(fields, lineNumber+1)
			if err != nil {
				return GoooDeclarationIR{}, err
			}
			for _, existing := range ir.Entities {
				if existing.Name == entity.Name || existing.ID == entity.ID {
					return GoooDeclarationIR{}, fmt.Errorf("duplicate entity at line %d", lineNumber+1)
				}
			}
			ir.Entities = append(ir.Entities, entity)
			currentEntity = len(ir.Entities) - 1
		case "property":
			if currentEntity < 0 {
				return GoooDeclarationIR{}, fmt.Errorf("property without entity at line %d", lineNumber+1)
			}
			if len(fields) != 3 {
				return GoooDeclarationIR{}, fmt.Errorf("invalid property at line %d", lineNumber+1)
			}
			property := GoooPropertyIR{Name: fields[1], Type: fields[2]}
			for _, existing := range ir.Entities[currentEntity].Properties {
				if existing.Name == property.Name {
					return GoooDeclarationIR{}, fmt.Errorf("duplicate property at lineNumber %d", lineNumber+1)
				}
			}
			ir.Entities[currentEntity].Properties = append(ir.Entities[currentEntity].Properties, property)
		case "activity":
			activity, err := parseGoooActivity(line, lineNumber+1)
			if err != nil {
				return GoooDeclarationIR{}, err
			}
			for _, existing := range ir.Activities {
				if existing.Name == activity.Name {
					return GoooDeclarationIR{}, fmt.Errorf("duplicate activity at line %d", lineNumber+1)
				}
			}
			ir.Activities = append(ir.Activities, activity)
			currentEntity = -1
		default:
			return GoooDeclarationIR{}, fmt.Errorf("unsupported declaration at line %d", lineNumber+1)
		}
	}
	if !seenPackage || !seenNamespace {
		return GoooDeclarationIR{}, fmt.Errorf("package and namespace are required")
	}
	if len(ir.Entities) == 0 && len(ir.Activities) == 0 {
		return GoooDeclarationIR{}, fmt.Errorf("declaration must contain an entity or activity")
	}
	return ir, nil
}

func parseGoooEntity(fields []string, lineNumber int) (GoooEntityIR, error) {
	if len(fields) != 4 || fields[2] != "id" {
		return GoooEntityIR{}, fmt.Errorf("invalid entity at line %d", lineNumber)
	}
	id, err := strconv.Unquote(fields[3])
	if err != nil || strings.TrimSpace(id) == "" {
		return GoooEntityIR{}, fmt.Errorf("invalid entity id at line %d", lineNumber)
	}
	return GoooEntityIR{Name: fields[1], ID: id}, nil
}

func parseGoooActivity(line string, lineNumber int) (GoooActivityIR, error) {
	rest := strings.TrimSpace(strings.TrimPrefix(line, "activity"))
	open := strings.Index(rest, "(")
	close := strings.LastIndex(rest, ")")
	if open <= 0 || close <= open {
		return GoooActivityIR{}, fmt.Errorf("invalid activity at line %d", lineNumber)
	}
	after := strings.TrimSpace(rest[close+1:])
	if !strings.HasPrefix(after, "->") {
		return GoooActivityIR{}, fmt.Errorf("invalid activity return at line %d", lineNumber)
	}
	name := strings.TrimSpace(rest[:open])
	returnType := strings.TrimSpace(strings.TrimPrefix(after, "->"))
	if name == "" || returnType == "" {
		return GoooActivityIR{}, fmt.Errorf("invalid activity signature at line %d", lineNumber)
	}
	parameters, err := parseGoooActivityParameters(rest[open+1:close], lineNumber)
	if err != nil {
		return GoooActivityIR{}, err
	}
	return GoooActivityIR{Name: name, Parameters: parameters, ReturnType: returnType}, nil
}

func parseGoooActivityParameters(source string, lineNumber int) ([]GoooActivityParameterIR, error) {
	if strings.TrimSpace(source) == "" {
		return nil, nil
	}
	rawParameters := strings.Split(source, ",")
	parameters := make([]GoooActivityParameterIR, 0, len(rawParameters))
	for _, rawParameter := range rawParameters {
		fields := strings.Fields(strings.TrimSpace(rawParameter))
		switch len(fields) {
		case 1:
			parameters = append(parameters, GoooActivityParameterIR{Type: fields[0]})
		case 2:
			parameters = append(parameters, GoooActivityParameterIR{Name: fields[0], Type: fields[1]})
		default:
			return nil, fmt.Errorf("invalid activity parameter at line %d", lineNumber)
		}
	}
	return parameters, nil
}
