package gooo

import "sort"

// capabilityQueryDeclarationBindings is the versioned, explicit bridge from
// source observations to catalog capabilities. Natural-language wording never
// creates a binding that is absent from this table.
var capabilityQueryDeclarationBindings = map[string][]string{
	"module":     {"declaration_analysis"},
	"input":      {"declaration_analysis"},
	"output":     {"canonical_generation", "ir_generation"},
	"status":     {"declaration_analysis"},
	"field":      {"declaration_analysis"},
	"constraint": {"declaration_analysis"},
	"entity":    {"declaration_analysis"},
	"operation": {"canonical_generation", "ir_generation"},
	"observe":   {"provenance", "reverse_observation", "round_trip_observation"},
	"transform": {"canonical_generation", "ir_generation"},
	"contract":  {"declaration_analysis", "provenance"},
	"policy":    {"security_boundary"},
	"workflow":  {"support_triage"},
}

func capabilityQueryObservedCapabilities(signals []string) []string {
	seen := make(map[string]struct{})
	for _, signal := range signals {
		for _, capabilityID := range capabilityQueryDeclarationBindings[signal] {
			seen[capabilityID] = struct{}{}
		}
	}
	capabilityIDs := make([]string, 0, len(seen))
	for capabilityID := range seen {
		capabilityIDs = append(capabilityIDs, capabilityID)
	}
	sort.Strings(capabilityIDs)
	return capabilityIDs
}

func bindCapabilityQueryDeclaration(response CapabilityQueryResponse) CapabilityQueryResponse {
	declaration := response.Declaration
	if declaration == nil || !declaration.Bound || len(response.Capabilities) == 0 {
		return response
	}
	if response.Status == CapabilityQueryDeferred {
		return response
	}
	allowed := make(map[string]struct{}, len(declaration.ObservedCapabilities))
	for _, capabilityID := range declaration.ObservedCapabilities {
		allowed[capabilityID] = struct{}{}
	}
	filtered := make([]CapabilityQueryCapability, 0, len(response.Capabilities))
	for _, capability := range response.Capabilities {
		if _, ok := allowed[capability.ID]; ok {
			filtered = append(filtered, capability)
		}
	}
	if len(filtered) == 0 {
		response.Capabilities = nil
		response.Status = CapabilityQueryUnknown
		response.FirstMismatch = "declaration_capability_binding"
		response.MissingStage = "declaration_capability_binding"
		return response
	}
	response.Capabilities = filtered
	response.Suggestions = capabilityQueryIDs(filtered)
	response.SuggestedQueries = capabilityQueryExamples(filtered)
	if hasDeferredCapabilityQuery(filtered) {
		response.Status = CapabilityQueryDeferred
		response.FirstMismatch = "external_boundary"
		response.MissingStage = "explicit_external_boundary"
		return response
	}
	response.Status = CapabilityQueryAvailable
	response.FirstMismatch = ""
	response.MissingStage = ""
	return response
}

func capabilityQueryIDs(capabilities []CapabilityQueryCapability) []string {
	ids := make([]string, 0, len(capabilities))
	for _, capability := range capabilities {
		ids = append(ids, capability.ID)
	}
	return ids
}

func capabilityQueryExamples(capabilities []CapabilityQueryCapability) []string {
	examples := make([]string, 0, len(capabilities))
	for _, capability := range capabilities {
		if capability.ExampleQuery != "" {
			examples = append(examples, capability.ExampleQuery)
		}
	}
	sort.Strings(examples)
	return examples
}