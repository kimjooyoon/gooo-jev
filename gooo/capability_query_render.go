package gooo

import (
	"fmt"
	"strings"
)

// RenderCapabilityQuery turns a validated capability-query trail into readable,
// provenance-bound guidance. It never executes a capability or grants authority.
func RenderCapabilityQuery(trail CapabilityQueryTrail, guide CapabilityQueryGuide) (string, error) {
	if err := trail.Validate(); err != nil {
		return "", fmt.Errorf("capability query trail: %w", err)
	}
	if err := guide.Validate(); err != nil {
		return "", fmt.Errorf("capability query guide: %w", err)
	}
	if err := validateCapabilityQueryGuideBinding(trail, guide); err != nil {
		return "", err
	}

	var builder strings.Builder
	fmt.Fprintf(&builder, "gooo capability discovery\nquestion: %s\nstatus: %s\n", trail.Response.Query, trail.Response.Status)
	switch trail.Response.Status {
	case CapabilityQueryAvailable:
		builder.WriteString("gooo can currently describe these read-only capabilities from the catalog.\n")
	case CapabilityQueryDeferred:
		builder.WriteString("The question reaches an external boundary. gooo will not execute or authorize it; bind the required boundary explicitly.\n")
	case CapabilityQueryUnknown:
		builder.WriteString("No catalog capability matched this question. Refine the question with one of the suggested examples.\n")
	}

	for _, capability := range trail.Response.Capabilities {
		fmt.Fprintf(
			&builder,
			"\n- [%s] %s: %s\n  next: %s\n  example: %s\n",
			capability.State,
			capability.ID,
			capability.Description,
			capability.NextOperation,
			capability.ExampleQuery,
		)
	}

	if len(guide.Questions) > 0 {
		builder.WriteString("\nnext questions:\n")
		for _, question := range guide.Questions {
			fmt.Fprintf(&builder, "- %s\n", question)
		}
	}

	if trail.Response.Declaration != nil {
		declaration := trail.Response.Declaration
		fmt.Fprintf(&builder, "\ndeclaration observation:\n- bound: %t\n- source digest: %s\n", declaration.Bound, declaration.SourceDigest)
		if len(declaration.ObservedSignals) > 0 {
			fmt.Fprintf(&builder, "- observed signals: %s\n", strings.Join(declaration.ObservedSignals, ", "))
		} else {
			builder.WriteString("- observed signals: none\n")
		}
	}

	builder.WriteString("\nconstraints:\n")
	for _, constraint := range guide.Constraints {
		fmt.Fprintf(&builder, "- %s\n", constraint)
	}
	fmt.Fprintf(
		&builder,
		"\nevidence:\n- query digest: %s\n- trail digest: %s\n- guide digest: %s\n- non-executing: true\n- non-authorizing: true\n",
		trail.Response.QueryDigest,
		trail.EvidenceDigest,
		guide.EvidenceDigest,
	)
	return builder.String(), nil
}

func validateCapabilityQueryGuideBinding(trail CapabilityQueryTrail, guide CapabilityQueryGuide) error {
	if guide.Status != trail.Response.Status {
		return fmt.Errorf("capability query guide status is not bound to the trail")
	}
	if len(guide.CapabilityIDs) != len(trail.Response.Capabilities) ||
		len(guide.NextOperations) != len(trail.Response.Capabilities) {
		return fmt.Errorf("capability query guide capabilities are not bound to the trail")
	}
	for index, capability := range trail.Response.Capabilities {
		if guide.CapabilityIDs[index] != capability.ID ||
			guide.NextOperations[index] != capability.NextOperation {
			return fmt.Errorf("capability query guide capability %d is not bound to the trail", index)
		}
	}
	return nil
}
