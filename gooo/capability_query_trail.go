package gooo

import (
	"fmt"
	"strings"
)

// CapabilityQueryTrail explains how a natural-language question was mapped to
// the read-only capability catalog. It is evidence about discovery, not an
// execution plan, authorization grant, or semantic proof.
type CapabilityQueryTrail struct {
	Response        CapabilityQueryResponse `json:"response"`
	Intent          string                 `json:"intent"`
	DiscoveryPath   []string               `json:"discovery_path"`
	NextQuestions   []string               `json:"next_questions"`
	EvidenceDigest  string                 `json:"evidence_digest"`
}

// DiscoverCapabilityQueryTrail adds a deterministic explanation to the
// existing capability-query response. Declaration binding is included only
// when source text is supplied.
func DiscoverCapabilityQueryTrail(query, declaration string) CapabilityQueryTrail {
	response := DiscoverCapabilityQuery(query)
	if strings.TrimSpace(declaration) != "" {
		response = DiscoverCapabilityQueryWithDeclaration(query, declaration)
	}

	path := []string{"normalize_question"}
	intent := "unknown"
	switch {
	case strings.TrimSpace(query) == "":
		intent = "missing_query"
		path = append(path, "missing_question")
	case len(response.Capabilities) == 0:
		path = append(path, "no_catalog_match", "preserve_unknown")
	case capabilityQueryIsOverview(strings.ToLower(strings.TrimSpace(query))):
		intent = "overview"
		path = append(path, "classify_overview", "catalog_match")
	case hasDeferredCapabilityQuery(response.Capabilities):
		intent = "boundary_request"
		path = append(path, "catalog_match", "classify_external_boundary", "preserve_deferred")
	default:
		intent = "capability_request"
		path = append(path, "catalog_match", "classify_available")
	}
	if response.Declaration != nil {
		path = append(path, "bind_declaration")
		if response.Declaration.Bound {
			path = append(path, "observe_declaration_signals")
		} else {
			path = append(path, "declaration_missing")
		}
	}

	trail := CapabilityQueryTrail{
		Response:       response,
		Intent:         intent,
		DiscoveryPath:  path,
		NextQuestions:  capabilityQueryTrailNextQuestions(response),
	}
	trail.EvidenceDigest = digestCapabilityQueryTrail(trail)
	return trail
}

func capabilityQueryTrailNextQuestions(response CapabilityQueryResponse) []string {
	questions := make([]string, 0, 3)
	seen := make(map[string]struct{})
	appendQuestion := func(question string) {
		question = strings.TrimSpace(question)
		if question == "" || len(questions) >= 3 {
			return
		}
		if _, ok := seen[question]; ok {
			return
		}
		seen[question] = struct{}{}
		questions = append(questions, question)
	}
	for _, capability := range response.Capabilities {
		appendQuestion(capability.ExampleQuery)
	}
	for _, question := range response.SuggestedQueries {
		appendQuestion(question)
	}
	if response.Status == CapabilityQueryDeferred {
		appendQuestion("What explicit external boundary is required before execution?")
	}
	return questions
}

func (trail CapabilityQueryTrail) Validate() error {
	if err := trail.Response.Validate(); err != nil {
		return fmt.Errorf("capability query trail response: %w", err)
	}
	if strings.TrimSpace(trail.Intent) == "" {
		return fmt.Errorf("capability query trail has no intent")
	}
	if len(trail.DiscoveryPath) == 0 || trail.DiscoveryPath[0] != "normalize_question" {
		return fmt.Errorf("capability query trail has no normalization stage")
	}
	if len(trail.NextQuestions) == 0 {
		return fmt.Errorf("capability query trail has no next questions")
	}
	if !validDigest(trail.EvidenceDigest) {
		return fmt.Errorf("capability query trail evidence digest is invalid")
	}
	if digestCapabilityQueryTrail(trail) != trail.EvidenceDigest {
		return fmt.Errorf("capability query trail evidence digest does not match its evidence")
	}
	return nil
}

func digestCapabilityQueryTrail(trail CapabilityQueryTrail) string {
	parts := []string{"gooo-capability-query-trail", trail.Response.QueryDigest, trail.Intent}
	parts = append(parts, trail.DiscoveryPath...)
	parts = append(parts, trail.NextQuestions...)
	return digestString(strings.Join(parts, "|"))
}
