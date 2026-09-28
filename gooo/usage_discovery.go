package gooo

import (
	"fmt"
	"strings"
)

type CapabilityState string

const (
	CapabilityAvailable CapabilityState = "AVAILABLE"
	CapabilityUnknown   CapabilityState = "UNKNOWN"
)

type UsageCapability struct {
	ID             string         `json:"id"`
	State          CapabilityState `json:"state"`
	Stage          string         `json:"stage"`
	Description    string         `json:"description"`
	Reason         string         `json:"reason"`
	NextOperation  string         `json:"next_operation"`
	NonExecuting   bool           `json:"non_executing"`
	NonAuthorizing bool           `json:"non_authorizing"`
}

type UsageDiscoveryResponse struct {
	Status          string                    `json:"status"`
	MissingStage    string                    `json:"missing_stage"`
	SourceDigest    string                    `json:"source_digest"`
	IRDigest        string                    `json:"ir_digest"`
	Prefix          string                    `json:"prefix"`
	Completion      SyntaxCompletionResponse  `json:"completion"`
	Capabilities    []UsageCapability         `json:"capabilities"`
	DiscoveryDigest string                    `json:"discovery_digest"`
	NonExecuting    bool                      `json:"non_executing"`
	NonAuthorizing  bool                      `json:"non_authorizing"`
}

// DiscoverUsage reports available and not-yet-bound language capabilities.
// UNKNOWN capabilities retain the next operation instead of guessing success.
func DiscoverUsage(source, prefix string) UsageDiscoveryResponse {
	completion := CompleteSyntax(source, prefix)
	response := UsageDiscoveryResponse{
		Status:         completion.Status,
		MissingStage:   completion.MissingStage,
		SourceDigest:   completion.SourceDigest,
		IRDigest:       completion.IRDigest,
		Prefix:         prefix,
		Completion:     completion,
		NonExecuting:   true,
		NonAuthorizing: true,
	}
	add := func(id string, state CapabilityState, stage, description, reason, next string) {
		response.Capabilities = append(response.Capabilities, UsageCapability{
			ID:             id,
			State:          state,
			Stage:          stage,
			Description:    description,
			Reason:         reason,
			NextOperation:  next,
			NonExecuting:   true,
			NonAuthorizing: true,
		})
	}
	add("syntax_completion", CapabilityAvailable, "SYNTAX", "Find keyword and symbol candidates for the current prefix.", "keyword_and_symbol_candidates_are_available", "edit_declaration")
	add("declaration_analysis", CapabilityAvailable, "ANALYSIS", "Inspect source-bound declaration diagnostics without execution.", "source_bound_analysis_is_available", "inspect_diagnostics")
	if completion.Status == "BOUND" {
		add("ir_generation", CapabilityAvailable, "GENERATION", "Derive the declaration's intermediate representation.", "declaration_is_bound_to_ir", "generate_canonical_declaration")
		add("canonical_generation", CapabilityAvailable, "GENERATION", "Produce a canonical .gooo declaration from the bound representation.", "canonical_gooo_generation_is_available", "write_generated_declaration")
		add("round_trip_observation", CapabilityAvailable, "REVERSE_OBSERVATION", "Reparse generated source and compare its representation.", "generated_declaration_can_be_reparsed", "compare_round_trip_ir")
		add("reverse_observation", CapabilityAvailable, "REVERSE_OBSERVATION", "Inspect the digest at the generation boundary.", "generation_boundary_is_observable", "inspect_reverse_digest")
		add("lsp_symbol_completion", CapabilityAvailable, "SYMBOLS", "Complete symbols from the bound declaration in the editor.", "bound_symbols_can_be_completed", "request_symbol_completion")
	} else {
		missing := completion.MissingStage
		if strings.TrimSpace(missing) == "" {
			missing = "SYNTAX"
		}
		add("ir_generation", CapabilityUnknown, missing, "Derive an intermediate representation after the declaration is complete.", "declaration_is_not_bound", "complete_declaration")
		add("canonical_generation", CapabilityUnknown, missing, "Produce canonical source after the generation boundary is bound.", "generation_boundary_is_not_bound", "complete_declaration")
		add("round_trip_observation", CapabilityUnknown, missing, "Compare generated and reparsed representations after binding.", "reverse_observation_is_not_bound", "complete_declaration")
		add("reverse_observation", CapabilityUnknown, missing, "Inspect generation evidence after the declaration is bound.", "reverse_observation_is_not_bound", "complete_declaration")
		add("lsp_symbol_completion", CapabilityUnknown, missing, "Complete bound symbols after the declaration is complete.", "symbols_are_not_bound", "complete_declaration")
	}
	response.DiscoveryDigest = digestUsageDiscovery(response)
	return response
}

func (response UsageDiscoveryResponse) Validate() error {
	if response.Status != "BOUND" && response.Status != "UNKNOWN" {
		return fmt.Errorf("usage discovery status %q is invalid", response.Status)
	}
	if !validDigest(response.SourceDigest) {
		return fmt.Errorf("usage discovery source digest is invalid")
	}
	if err := response.Completion.Validate(); err != nil {
		return err
	}
	if !response.NonExecuting || !response.NonAuthorizing {
		return fmt.Errorf("usage discovery crossed a capability boundary")
	}
	if len(response.Capabilities) == 0 || !validDigest(response.DiscoveryDigest) {
		return fmt.Errorf("usage discovery capabilities are incomplete")
	}
	seen := map[string]struct{}{}
	for _, capability := range response.Capabilities {
		if strings.TrimSpace(capability.ID) == "" || strings.TrimSpace(capability.Stage) == "" ||
			strings.TrimSpace(capability.Description) == "" || strings.TrimSpace(capability.Reason) == "" ||
			strings.TrimSpace(capability.NextOperation) == "" {
			return fmt.Errorf("usage capability is incomplete")
		}
		if capability.State != CapabilityAvailable && capability.State != CapabilityUnknown {
			return fmt.Errorf("usage capability state is invalid")
		}
		if !capability.NonExecuting || !capability.NonAuthorizing {
			return fmt.Errorf("usage capability crossed a capability boundary")
		}
		if _, exists := seen[capability.ID]; exists {
			return fmt.Errorf("usage capability %q is duplicated", capability.ID)
		}
		seen[capability.ID] = struct{}{}
	}
	if digestUsageDiscovery(response) != response.DiscoveryDigest {
		return fmt.Errorf("usage discovery digest does not match its evidence")
	}
	return nil
}

func digestUsageDiscovery(response UsageDiscoveryResponse) string {
	value := fmt.Sprintf("%s|%s|%s|%s|%s", response.Status, response.SourceDigest, response.IRDigest, response.Prefix, response.Completion.ItemsDigest)
	for _, capability := range response.Capabilities {
		value += fmt.Sprintf("|%s|%s|%s|%s|%s|%s", capability.ID, capability.State, capability.Stage, capability.Description, capability.Reason, capability.NextOperation)
	}
	return digestString(value)
}

