package gooo

import (
	"fmt"
	"strings"
)

// CapabilityQueryHover is the editor-facing projection of declaration-bound
// capability discovery. It is explanatory evidence, never an execution plan
// or an authorization grant.
type CapabilityQueryHover struct {
	SourceDigest             string                       `json:"source_digest"`
	QueryDigest              string                       `json:"query_digest"`
	Query                    string                       `json:"query"`
	Status                   CapabilityQueryState         `json:"status"`
	Capabilities             []CapabilityQueryCapability `json:"capabilities"`
	Declaration              *CapabilityQueryDeclaration  `json:"declaration,omitempty"`
	Suggestions              []string                     `json:"suggestions"`
	SuggestedQueries         []string                     `json:"suggested_queries"`
	FirstMismatch            string                       `json:"first_mismatch"`
	MissingStage             string                       `json:"missing_stage"`
	Intent                   string                       `json:"intent"`
	DiscoveryPath            []string                     `json:"discovery_path"`
	NextQuestions            []string                     `json:"next_questions"`
	Constraints              []string                     `json:"constraints"`
	DiscoveryEvidenceDigest  string                       `json:"discovery_evidence_digest"`
	EvidenceDigest           string                       `json:"evidence_digest"`
	NonExecuting             bool                         `json:"non_executing"`
	NonAuthorizing           bool                         `json:"non_authorizing"`
}

// DiscoverCapabilityQueryHover projects the existing capability trail into a
// deterministic LSP-friendly payload without creating a second catalog.
func DiscoverCapabilityQueryHover(query, declaration string) CapabilityQueryHover {
	trail := DiscoverCapabilityQueryTrail(query, declaration)
	sourceDigest := digestString("gooo-capability-declaration|" + strings.TrimSpace(declaration))
	if trail.Response.Declaration != nil {
		sourceDigest = trail.Response.Declaration.SourceDigest
	}
	hover := CapabilityQueryHover{
		SourceDigest:            sourceDigest,
		QueryDigest:             trail.Response.QueryDigest,
		Query:                   trail.Response.Query,
		Status:                  trail.Response.Status,
		Capabilities:            append([]CapabilityQueryCapability(nil), trail.Response.Capabilities...),
		Declaration:             trail.Response.Declaration,
		Suggestions:             append([]string(nil), trail.Response.Suggestions...),
		SuggestedQueries:        append([]string(nil), trail.Response.SuggestedQueries...),
		FirstMismatch:           trail.Response.FirstMismatch,
		MissingStage:            trail.Response.MissingStage,
		Intent:                  trail.Intent,
		DiscoveryPath:           append([]string(nil), trail.DiscoveryPath...),
		NextQuestions:           append([]string(nil), trail.NextQuestions...),
		Constraints:             append([]string(nil), capabilityQueryOverviewConstraints...),
		DiscoveryEvidenceDigest: trail.EvidenceDigest,
		NonExecuting:            trail.Response.NonExecuting,
		NonAuthorizing:          trail.Response.NonAuthorizing,
	}
	hover.EvidenceDigest = digestCapabilityQueryHover(hover)
	return hover
}

func (hover CapabilityQueryHover) Validate() error {
	response := CapabilityQueryResponse{
		Status:           hover.Status,
		Query:            hover.Query,
		Capabilities:     hover.Capabilities,
		Declaration:      hover.Declaration,
		Suggestions:      hover.Suggestions,
		SuggestedQueries: hover.SuggestedQueries,
		FirstMismatch:    hover.FirstMismatch,
		MissingStage:     hover.MissingStage,
		QueryDigest:      hover.QueryDigest,
		NonExecuting:     hover.NonExecuting,
		NonAuthorizing:   hover.NonAuthorizing,
	}
	if err := response.Validate(); err != nil {
		return fmt.Errorf("capability query hover response: %w", err)
	}
	if !validDigest(hover.SourceDigest) || !validDigest(hover.DiscoveryEvidenceDigest) || !validDigest(hover.EvidenceDigest) {
		return fmt.Errorf("capability query hover digests are invalid")
	}
	if hover.Declaration != nil && hover.SourceDigest != hover.Declaration.SourceDigest {
		return fmt.Errorf("capability query hover source digest is not declaration-bound")
	}
	if len(hover.DiscoveryPath) == 0 || hover.DiscoveryPath[0] != "normalize_question" || len(hover.NextQuestions) == 0 {
		return fmt.Errorf("capability query hover discovery evidence is incomplete")
	}
	if len(hover.Constraints) != len(capabilityQueryOverviewConstraints) || !sameCapabilityOverviewStrings(hover.Constraints, capabilityQueryOverviewConstraints) {
		return fmt.Errorf("capability query hover constraints are incomplete")
	}
	trail := CapabilityQueryTrail{
		Response:      response,
		Intent:        hover.Intent,
		DiscoveryPath: hover.DiscoveryPath,
		NextQuestions: hover.NextQuestions,
	}
	if digestCapabilityQueryTrail(trail) != hover.DiscoveryEvidenceDigest {
		return fmt.Errorf("capability query hover discovery digest does not match")
	}
	if digestCapabilityQueryHover(hover) != hover.EvidenceDigest {
		return fmt.Errorf("capability query hover evidence digest does not match")
	}
	if !hover.NonExecuting || !hover.NonAuthorizing {
		return fmt.Errorf("capability query hover crossed an execution or authorization boundary")
	}
	return nil
}

func digestCapabilityQueryHover(hover CapabilityQueryHover) string {
	parts := []string{
		"gooo-capability-query-hover",
		hover.SourceDigest,
		hover.QueryDigest,
		hover.Query,
		string(hover.Status),
		hover.FirstMismatch,
		hover.MissingStage,
		hover.Intent,
		strings.Join(hover.DiscoveryPath, "|"),
		strings.Join(hover.NextQuestions, "|"),
		strings.Join(hover.Suggestions, "|"),
		strings.Join(hover.SuggestedQueries, "|"),
		strings.Join(hover.Constraints, "|"),
		hover.DiscoveryEvidenceDigest,
		fmt.Sprintf("%t|%t", hover.NonExecuting, hover.NonAuthorizing),
	}
	for _, capability := range hover.Capabilities {
		parts = append(parts, capability.ID, string(capability.State), capability.Stage, capability.Description, capability.NextOperation, capability.ExampleQuery)
	}
	if hover.Declaration != nil {
		parts = append(parts, "declaration", fmt.Sprintf("%t", hover.Declaration.Bound), hover.Declaration.SourceDigest)
		parts = append(parts, hover.Declaration.ObservedSignals...)
		parts = append(parts, hover.Declaration.ObservedCapabilities...)
	}
	return digestString(strings.Join(parts, "|"))
}
