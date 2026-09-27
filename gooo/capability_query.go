package gooo

import (
	"fmt"
	"sort"
	"strings"
)

type CapabilityQueryState string

const (
	CapabilityQueryAvailable CapabilityQueryState = "AVAILABLE"
	CapabilityQueryDeferred  CapabilityQueryState = "DEFERRED"
	CapabilityQueryUnknown   CapabilityQueryState = "UNKNOWN"
)

type CapabilityQueryCapability struct {
	ID             string              `json:"id"`
	State          CapabilityQueryState `json:"state"`
	Stage          string              `json:"stage"`
	Description    string              `json:"description"`
	NextOperation  string              `json:"next_operation"`
	ExampleQuery   string              `json:"example_query"`
	NonExecuting   bool                `json:"non_executing"`
	NonAuthorizing bool                `json:"non_authorizing"`
}

type CapabilityQueryResponse struct {
	Status          CapabilityQueryState       `json:"status"`
	Query           string                     `json:"query"`
	Capabilities    []CapabilityQueryCapability `json:"capabilities"`
	Declaration     *CapabilityQueryDeclaration `json:"declaration,omitempty"`
	Suggestions     []string                   `json:"suggestions"`
	SuggestedQueries []string                  `json:"suggested_queries"`
	FirstMismatch   string                     `json:"first_mismatch"`
	MissingStage    string                     `json:"missing_stage"`
	QueryDigest     string                     `json:"query_digest"`
	NonExecuting    bool                       `json:"non_executing"`
	NonAuthorizing  bool                       `json:"non_authorizing"`
}

// CapabilityQueryDeclaration records only source-bound observations from a
// declaration. It is not an execution plan, a semantic proof, or an
// authorization grant.
type CapabilityQueryDeclaration struct {
	Bound           bool     `json:"bound"`
	SourceDigest    string   `json:"source_digest"`
	ObservedSignals []string `json:"observed_signals"`
}

type capabilityQueryEntry struct {
	ID            string
	Stage         string
	Description   string
	NextOperation string
	ExampleQuery  string
	Aliases       []string
	Safe          bool
}

var capabilityQueryCatalog = []capabilityQueryEntry{
	{ID: "syntax_completion", Stage: "SYNTAX", Description: "suggest keywords and symbols for a declaration prefix", NextOperation: "edit_declaration", ExampleQuery: "How do I complete a .gooo declaration?", Aliases: []string{"syntax", "completion", "autocomplete", "lsp", "문법", "완성"}, Safe: true},
	{ID: "declaration_analysis", Stage: "ANALYSIS", Description: "inspect source-bound declaration diagnostics without execution", NextOperation: "inspect_diagnostics", ExampleQuery: "Can gooo analyze this .gooo declaration?", Aliases: []string{"analysis", "diagnostic", "diagnostics", "분석", "진단"}, Safe: true},
	{ID: "ir_generation", Stage: "GENERATION", Description: "derive an intermediate representation from a bound declaration", NextOperation: "generate_ir", ExampleQuery: "What IR can this .gooo declaration produce?", Aliases: []string{"ir", "intermediate representation", "representation", "중간 표현"}, Safe: true},
	{ID: "canonical_generation", Stage: "GENERATION", Description: "produce a canonical .gooo declaration from its representation", NextOperation: "write_generated_declaration", ExampleQuery: "How do I generate a canonical .gooo declaration?", Aliases: []string{"generate", "generation", "codegen", "code generation", "생성", "코드 생성"}, Safe: true},
	{ID: "round_trip_observation", Stage: "REVERSE_OBSERVATION", Description: "reparse generated source and compare its representation", NextOperation: "compare_round_trip_ir", ExampleQuery: "Can gooo compare a generated declaration round trip?", Aliases: []string{"round trip", "roundtrip", "reparse", "재파싱"}, Safe: true},
	{ID: "reverse_observation", Stage: "REVERSE_OBSERVATION", Description: "inspect source and generation evidence at a reverse boundary", NextOperation: "inspect_reverse_digest", ExampleQuery: "How do I inspect reverse observation evidence?", Aliases: []string{"reverse observation", "reverse", "provenance", "origin", "기원", "역관찰"}, Safe: true},
	{ID: "provenance", Stage: "PROVENANCE", Description: "trace source, generation, and reverse-observation evidence", NextOperation: "inspect_provenance_chain", ExampleQuery: "Where did this .gooo declaration come from?", Aliases: []string{"provenance chain", "source lineage", "기원 추적"}, Safe: true},
	{ID: "feedback_trend", Stage: "FEEDBACK", Description: "compare feedback and calibration windows with evidence lineage", NextOperation: "compare_feedback_window", ExampleQuery: "How has gooo feedback changed over time?", Aliases: []string{"feedback", "trend", "calibration", "피드백", "추세"}, Safe: true},
	{ID: "support_triage", Stage: "WORKFLOW", Description: "structure support-triage workflows and their next observations", NextOperation: "inspect_support_route", ExampleQuery: "How should I triage this gooo support request?", Aliases: []string{"support", "support triage", "triage", "workflow", "지원", "분류"}, Safe: true},
	{ID: "security_boundary", Stage: "SECURITY_BOUNDARY", Description: "observe workload identity and network capability boundaries", NextOperation: "bind_external_security_evidence", ExampleQuery: "What external security boundary is required?", Aliases: []string{"security", "spiffe", "workload identity", "network", "credential", "보안"}, Safe: false},
	{ID: "execution_authorization", Stage: "EXECUTION_AUTHORIZATION_BOUNDARY", Description: "execution and authorization require an explicit external boundary", NextOperation: "provide_explicit_external_boundary", ExampleQuery: "What explicit boundary is needed before execution?", Aliases: []string{"execute", "execution", "run", "authorize", "authorization", "permission", "실행", "권한"}, Safe: false},
}

// DiscoverCapabilityQuery explains what the language can currently expose from
// a natural-language question. It never executes a plan or grants authority.
func DiscoverCapabilityQuery(query string) CapabilityQueryResponse {
	response := CapabilityQueryResponse{
		Status:         CapabilityQueryUnknown,
		Query:          strings.TrimSpace(query),
		Suggestions:    capabilityQuerySuggestions(),
		SuggestedQueries: capabilityQueryExampleSuggestions(),
		FirstMismatch:  "query",
		MissingStage:   "capability_catalog",
		NonExecuting:   true,
		NonAuthorizing: true,
	}
	if response.Query == "" {
		response.MissingStage = "capability_query"
		response.QueryDigest = digestString("gooo-capability-query|" + response.Query + "|UNKNOWN|query|capability_query")
		return response
	}
	response.Capabilities = capabilityQueryMatches(response.Query)
	if len(response.Capabilities) == 0 {
		response.QueryDigest = digestCapabilityQuery(response)
		return response
	}
	if hasDeferredCapabilityQuery(response.Capabilities) {
		response.Status = CapabilityQueryDeferred
		response.FirstMismatch = "external_boundary"
		response.MissingStage = "explicit_external_boundary"
		response.QueryDigest = digestCapabilityQuery(response)
		return response
	}
	response.Status = CapabilityQueryAvailable
	response.FirstMismatch = ""
	response.MissingStage = ""
	response.QueryDigest = digestCapabilityQuery(response)
	return response
}

// DiscoverCapabilityQueryWithDeclaration adds a deterministic, read-only
// observation of the declaration that motivated the question. The query
// result remains catalog-driven; declaration signals only explain what was
// actually present in the supplied source.
func DiscoverCapabilityQueryWithDeclaration(query, declaration string) CapabilityQueryResponse {
	response := DiscoverCapabilityQuery(query)
	response.Declaration = inspectCapabilityQueryDeclaration(declaration)
	response.QueryDigest = digestCapabilityQuery(response)
	return response
}

func inspectCapabilityQueryDeclaration(declaration string) *CapabilityQueryDeclaration {
	raw := strings.TrimSpace(declaration)
	signals := make([]string, 0)
	for _, line := range strings.Split(raw, "\n") {
		normalized := strings.ToLower(strings.TrimSpace(line))
		for _, candidate := range []struct {
			prefix string
			id     string
		}{
			{prefix: "entity ", id: "entity"},
			{prefix: "operation ", id: "operation"},
			{prefix: "observe ", id: "observe"},
			{prefix: "transform ", id: "transform"},
			{prefix: "contract ", id: "contract"},
			{prefix: "policy ", id: "policy"},
			{prefix: "workflow ", id: "workflow"},
		} {
			if !strings.HasPrefix(normalized, candidate.prefix) {
				continue
			}
			seen := false
			for _, signal := range signals {
				if signal == candidate.id {
					seen = true
					break
				}
			}
			if !seen {
				signals = append(signals, candidate.id)
			}
		}
	}
	sort.Strings(signals)
	return &CapabilityQueryDeclaration{
		Bound:           raw != "",
		SourceDigest:    digestString("gooo-capability-declaration|" + raw),
		ObservedSignals: signals,
	}
}

func capabilityQueryMatches(query string) []CapabilityQueryCapability {
	normalized := strings.ToLower(strings.TrimSpace(query))
	matches := make([]CapabilityQueryCapability, 0)
	if capabilityQueryIsOverview(normalized) {
		for _, entry := range capabilityQueryCatalog {
			if entry.Safe {
				matches = append(matches, capabilityQueryMatch(entry))
			}
		}
	} else {
		for _, entry := range capabilityQueryCatalog {
			for _, alias := range entry.Aliases {
				if strings.Contains(normalized, strings.ToLower(alias)) {
					matches = append(matches, capabilityQueryMatch(entry))
					break
				}
			}
		}
	}
	sort.Slice(matches, func(i, j int) bool { return matches[i].ID < matches[j].ID })
	seen := make(map[string]struct{}, len(matches))
	unique := matches[:0]
	for _, match := range matches {
		if _, ok := seen[match.ID]; ok {
			continue
		}
		seen[match.ID] = struct{}{}
		unique = append(unique, match)
	}
	return unique
}

func capabilityQueryMatch(entry capabilityQueryEntry) CapabilityQueryCapability {
	state := CapabilityQueryAvailable
	if !entry.Safe {
		state = CapabilityQueryDeferred
	}
	return CapabilityQueryCapability{ID: entry.ID, State: state, Stage: entry.Stage, Description: entry.Description, NextOperation: entry.NextOperation, ExampleQuery: entry.ExampleQuery, NonExecuting: true, NonAuthorizing: true}
}

func capabilityQueryIsOverview(query string) bool {
	for _, phrase := range []string{"what can", "capabilities", "show examples", "discover", "help", "무엇을", "가능", "할 수"} {
		if strings.Contains(query, phrase) {
			return true
		}
	}
	return false
}

func capabilityQuerySuggestions() []string {
	suggestions := make([]string, 0)
	for _, entry := range capabilityQueryCatalog {
		if entry.Safe {
			suggestions = append(suggestions, entry.ID)
		}
	}
	sort.Strings(suggestions)
	return suggestions
}

func capabilityQueryExampleSuggestions() []string {
	suggestions := make([]string, 0)
	for _, entry := range capabilityQueryCatalog {
		if entry.Safe && strings.TrimSpace(entry.ExampleQuery) != "" {
			suggestions = append(suggestions, entry.ExampleQuery)
		}
	}
	sort.Strings(suggestions)
	return suggestions
}

func hasDeferredCapabilityQuery(capabilities []CapabilityQueryCapability) bool {
	for _, capability := range capabilities {
		if capability.State == CapabilityQueryDeferred {
			return true
		}
	}
	return false
}

func (response CapabilityQueryResponse) Validate() error {
	if response.Status != CapabilityQueryAvailable && response.Status != CapabilityQueryDeferred && response.Status != CapabilityQueryUnknown {
		return fmt.Errorf("capability query status %q is invalid", response.Status)
	}
	if len(response.SuggestedQueries) == 0 {
		return fmt.Errorf("capability query has no natural-language follow-up suggestions")
	}
	if !response.NonExecuting || !response.NonAuthorizing {
		return fmt.Errorf("capability query crossed an execution or authorization boundary")
	}
	if !validDigest(response.QueryDigest) {
		return fmt.Errorf("capability query digest is invalid")
	}
	if response.Declaration != nil && !validDigest(response.Declaration.SourceDigest) {
		return fmt.Errorf("capability query declaration digest is invalid")
	}
	if response.Status == CapabilityQueryAvailable && (len(response.Capabilities) == 0 || response.FirstMismatch != "" || response.MissingStage != "") {
		return fmt.Errorf("available capability query is incomplete")
	}
	if response.Status != CapabilityQueryAvailable && (response.FirstMismatch == "" || response.MissingStage == "") {
		return fmt.Errorf("unresolved capability query lost its boundary")
	}
	for _, capability := range response.Capabilities {
		if strings.TrimSpace(capability.ID) == "" || strings.TrimSpace(capability.Stage) == "" || strings.TrimSpace(capability.Description) == "" || strings.TrimSpace(capability.NextOperation) == "" || strings.TrimSpace(capability.ExampleQuery) == "" || capability.State == "" || !capability.NonExecuting || !capability.NonAuthorizing {
			return fmt.Errorf("capability query match is incomplete")
		}
	}
	if digestCapabilityQuery(response) != response.QueryDigest {
		return fmt.Errorf("capability query digest does not match its evidence")
	}
	return nil
}

func digestCapabilityQuery(response CapabilityQueryResponse) string {
	parts := []string{"gooo-capability-query", strings.ToLower(strings.TrimSpace(response.Query)), string(response.Status), response.FirstMismatch, response.MissingStage}
	for _, capability := range response.Capabilities {
		parts = append(parts, capability.ID, string(capability.State), capability.Stage, capability.Description, capability.NextOperation, capability.ExampleQuery)
	}
	parts = append(parts, response.Suggestions...)
	parts = append(parts, response.SuggestedQueries...)
	if response.Declaration != nil {
		parts = append(parts, "declaration", fmt.Sprintf("%t", response.Declaration.Bound), response.Declaration.SourceDigest)
		parts = append(parts, response.Declaration.ObservedSignals...)
	}
	return digestString(strings.Join(parts, "|"))
}
