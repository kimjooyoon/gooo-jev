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

	focused := DiscoverCapabilityQueryForDeclaration(declaration)
	trail.Response = focused
	trail.NextQuestions = capabilityQueryTrailNextQuestions(focused)
	trail.EvidenceDigest = digestCapabilityQueryTrail(trail)
	return trail
}