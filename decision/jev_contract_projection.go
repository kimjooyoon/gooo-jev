package decision

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

// GoooContractInput is the experimental contract-sketch input used by the
// usability adapter. It is not an execution or authorization request.
type GoooContractInput struct {
	SourceText     string
	NonAuthorizing bool
}

// GoooContractProjection binds a contract sketch to the supported declaration
// pipeline without claiming that the sketch is a complete language grammar.
type GoooContractProjection struct {
	Status             string
	MissingStage       string
	ContractDigest     string
	DeclarationSource  string
	Evidence           GoooDeclarationIRGeneration
	NonExecuting       bool
	NonAuthorizing     bool
}

// DeriveGoooContractProjection adapts the small contract sketch used by the
// support-triage example into a canonical .gooo declaration. It never executes
// an activity and never grants authorization.
func DeriveGoooContractProjection(input GoooContractInput) GoooContractProjection {
	output := GoooContractProjection{
		Status:         "UNKNOWN",
		NonExecuting:   true,
		NonAuthorizing: true,
	}
	if !input.NonAuthorizing {
		output.NonAuthorizing = false
		output.MissingStage = "authorization-boundary"
		return output
	}
	source := strings.TrimSpace(input.SourceText)
	if source == "" {
		output.MissingStage = "contract-source"
		return output
	}
	output.ContractDigest = digestContractSource(source)
	shape, missingStage := parseGoooContractShape(source)
	if missingStage != "" {
		output.MissingStage = missingStage
		return output
	}
	declaration, err := renderGoooContractDeclaration(shape)
	if err != nil {
		output.MissingStage = "contract-generation"
		return output
	}
	output.DeclarationSource = declaration
	output.Evidence = DeriveGoooDeclarationIRGeneration(
		GoooDeclarationIRGenerationInput{
			SourceText:     declaration,
			NonAuthorizing: true,
		},
	)
	output.Status = output.Evidence.Status
	output.MissingStage = output.Evidence.MissingStage
	return output
}

// Validate checks the adapter boundary and any available declaration evidence.
func (projection GoooContractProjection) Validate() error {
	if projection.Status != "UNKNOWN" && projection.Status != "ready" {
		return fmt.Errorf("invalid gooo contract projection status %q", projection.Status)
	}
	if !projection.NonExecuting || !projection.NonAuthorizing {
		return fmt.Errorf("gooo contract projection crossed a capability boundary")
	}
	if projection.Status == "UNKNOWN" {
		if projection.MissingStage == "" {
			return fmt.Errorf("unknown gooo contract projection lost its first missing stage")
		}
		return nil
	}
	if projection.ContractDigest == "" || projection.DeclarationSource == "" {
		return fmt.Errorf("ready gooo contract projection is missing source evidence")
	}
	if err := projection.Evidence.Validate(); err != nil {
		return err
	}
	if projection.Evidence.Status != "ready" {
		return fmt.Errorf("ready gooo contract projection has incomplete declaration evidence")
	}
	return nil
}

type goooContractShape struct {
	module      string
	inputs      []string
	statuses    []string
	fields      []string
	constraints map[string]bool
}

func parseGoooContractShape(source string) (goooContractShape, string) {
	shape := goooContractShape{constraints: make(map[string]bool)}
	for _, rawLine := range strings.Split(source, "\n") {
		line := strings.TrimSpace(rawLine)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Fields(line)
		switch parts[0] {
		case "module":
			if len(parts) != 2 || !validContractIdentifier(parts[1]) || shape.module != "" {
				return goooContractShape{}, "contract-syntax"
			}
			shape.module = parts[1]
		case "input":
			if len(parts) != 2 || !validContractIdentifier(parts[1]) || containsString(shape.inputs, parts[1]) {
				return goooContractShape{}, "contract-syntax"
			}
			shape.inputs = append(shape.inputs, parts[1])
		case "status":
			if len(parts) < 2 || !validContractIdentifier(parts[1]) || containsString(shape.statuses, parts[1]) {
				return goooContractShape{}, "contract-syntax"
			}
			shape.statuses = append(shape.statuses, parts[1])
		case "field":
			if len(parts) != 2 || !validContractIdentifier(parts[1]) || containsString(shape.fields, parts[1]) {
				return goooContractShape{}, "contract-syntax"
			}
			shape.fields = append(shape.fields, parts[1])
		case "constraint":
			if len(parts) != 2 || shape.constraints[parts[1]] {
				return goooContractShape{}, "contract-syntax"
			}
			shape.constraints[parts[1]] = true
		default:
			return goooContractShape{}, "contract-syntax"
		}
	}
	switch {
	case shape.module == "":
		return goooContractShape{}, "contract-module"
	case len(shape.inputs) == 0:
		return goooContractShape{}, "contract-input"
	case len(shape.statuses) == 0:
		return goooContractShape{}, "contract-status"
	case len(shape.fields) == 0:
		return goooContractShape{}, "contract-field"
	case !shape.constraints["no_execution"] || !shape.constraints["no_authorization"]:
		return goooContractShape{}, "contract-boundary"
	default:
		return shape, ""
	}
}

func renderGoooContractDeclaration(shape goooContractShape) (string, error) {
	if !validContractIdentifier(shape.module) {
		return "", fmt.Errorf("invalid contract module")
	}
	var builder strings.Builder
	fmt.Fprintf(&builder, "package %s\nnamespace contract\n\n", shape.module)
	builder.WriteString("entity ContractInput id \"contract-input\"\n")
	for _, input := range shape.inputs {
		fmt.Fprintf(&builder, "  property %s %s\n", input, contractInputType(input))
	}
	builder.WriteString("\nentity ContractReceipt id \"contract-receipt\"\n")
	for _, field := range shape.fields {
		fmt.Fprintf(&builder, "  property %s %s\n", field, contractFieldType(field))
	}
	builder.WriteString("\nactivity observe_")
	builder.WriteString(shape.module)
	builder.WriteString("(")
	parameters := make([]string, 0, len(shape.inputs))
	for _, input := range shape.inputs {
		parameters = append(parameters, input+" "+contractInputType(input))
	}
	builder.WriteString(strings.Join(parameters, ", "))
	builder.WriteString(") -> ContractReceipt\n")
	return builder.String(), nil
}

func contractInputType(name string) string {
	switch name {
	case "confidence", "review_threshold":
		return "number"
	case "confidence_method":
		return "text"
	default:
		return "text"
	}
}

func contractFieldType(name string) string {
	if name == "confidence" || name == "review_threshold" {
		return "number"
	}
	return "text"
}

func validContractIdentifier(value string) bool {
	if value == "" {
		return false
	}
	for index, character := range value {
		if (character >= 'a' && character <= 'z') ||
			(character >= 'A' && character <= 'Z') ||
			(character >= '0' && character <= '9' && index > 0) ||
			(character == '_' && index > 0) {
			continue
		}
		return false
	}
	return true
}

func containsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func digestContractSource(source string) string {
	sum := sha256.Sum256([]byte(source))
	return "sha256:" + hex.EncodeToString(sum[:])
}
