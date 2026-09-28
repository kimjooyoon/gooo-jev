package gooo

import "sort"

// DiscoverCapabilityQueryForDeclaration keeps natural-language discovery
// declaration-aware without executing or authorizing any operation. The
// catalog remains the source of capability meaning; declaration signals only
// focus the next safe questions.
func DiscoverCapabilityQueryForDeclaration(declaration string) CapabilityQueryResponse {
	response := DiscoverCapabilityQueryWithDeclaration(
		"What can gooo do with this declaration?",
		declaration,
	)
	if response.Declaration == nil || !response.Declaration.Bound {
		return response
	}
	focused := capabilityQueryDeclarationSuggestions(response.Declaration.ObservedSignals)
	if len(focused) > 0 {
		response.SuggestedQueries = focused
		response.QueryDigest = digestCapabilityQuery(response)
	}
	return response
}

func capabilityQueryDeclarationSuggestions(signals []string) []string {
	ids := make(map[string]struct{})
	for _, signal := range signals {
		for _, id := range capabilityQueryDeclarationCapabilityIDs[signal] {
			ids[id] = struct{}{}
		}
	}

	queries := make([]string, 0, len(ids))
	for _, entry := range capabilityQueryCatalog {
		if !entry.Safe {
			continue
		}
		if _, ok := ids[entry.ID]; ok {
			queries = append(queries, entry.ExampleQuery)
		}
	}
	sort.Strings(queries)
	return queries
}

var capabilityQueryDeclarationCapabilityIDs = map[string][]string{
	"module":     {"syntax_completion", "declaration_analysis", "ir_generation", "provenance"},
	"input":      {"declaration_analysis", "ir_generation", "canonical_generation", "round_trip_observation"},
	"output":     {"declaration_analysis", "canonical_generation", "round_trip_observation", "provenance"},
	"status":     {"declaration_analysis", "provenance", "feedback_trend"},
	"field":      {"declaration_analysis", "ir_generation", "canonical_generation"},
	"constraint": {"declaration_analysis", "provenance", "feedback_trend"},
	"entity":     {"declaration_analysis", "ir_generation", "provenance"},
	"operation":  {"ir_generation", "canonical_generation", "round_trip_observation"},
	"observe":    {"reverse_observation", "provenance", "feedback_trend"},
	"transform":  {"ir_generation", "canonical_generation", "round_trip_observation"},
	"contract":   {"declaration_analysis", "provenance", "support_triage"},
	"policy":     {"declaration_analysis", "provenance", "support_triage"},
	"workflow":   {"support_triage", "provenance", "feedback_trend"},
}