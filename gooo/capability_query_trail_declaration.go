package gooo

import "strings"

// DiscoverCapabilityQueryTrailForDeclaration focuses overview questions with
// declaration-aware suggestions while preserving explicit question matching.
// The returned trail remains descriptive and non-executing.
func DiscoverCapabilityQueryTrailForDeclaration(query, declaration string) CapabilityQueryTrail {
	trail := DiscoverCapabilityQueryTrail(query, declaration)
	if strings.TrimSpace(declaration) == "" || !capabilityQueryIsOverview(normalizeCapabilityQuery(query)) {
		return trail
	}

	focused := capabilityQueryDeclarationSuggestions(trail.Response.Declaration.ObservedSignals)
	if len(focused) == 0 {
		return trail
	}
	trail.Response.SuggestedQueries = focused
	trail.Response.QueryDigest = digestCapabilityQuery(trail.Response)
	trail.NextQuestions = capabilityQueryTrailNextQuestions(trail.Response)
	trail.EvidenceDigest = digestCapabilityQueryTrail(trail)
	return trail
}