package gooo

import (
	"fmt"
	"strings"
)

type CapabilityQueryOverviewMode string

const (
	CapabilityQueryOverviewCatalog        CapabilityQueryOverviewMode = "CATALOG_OVERVIEW"
	CapabilityQueryOverviewMatchedOptions CapabilityQueryOverviewMode = "MATCHED_OPTIONS"
	CapabilityQueryOverviewClarification  CapabilityQueryOverviewMode = "CLARIFICATION_REQUIRED"
)

type CapabilityQueryOverview struct {
	Trail              CapabilityQueryTrail       `json:"trail"`
	Mode               CapabilityQueryOverviewMode `json:"mode"`
	Summary            string                     `json:"summary"`
	SuggestedQuestions []string                   `json:"suggested_questions"`
	Constraints        []string                   `json:"constraints"`
	EvidenceDigest     string                     `json:"evidence_digest"`
}

var capabilityQueryOverviewConstraints = []string{
	"catalog_scope_only",
	"not_a_language_completeness_claim",
	"no_provider_invocation",
	"no_authorization_grant",
	"no_catalog_mutation",
}

func DiscoverCapabilityQueryOverview(query, declaration string) CapabilityQueryOverview {
	trail := DiscoverCapabilityQueryTrail(query, declaration)
	overview := CapabilityQueryOverview{
		Trail:       trail,
		Constraints: append([]string(nil), capabilityQueryOverviewConstraints...),
	}
	normalized := strings.ToLower(strings.TrimSpace(query))
	switch {
	case capabilityQueryIsOverview(normalized):
		overview.Mode = CapabilityQueryOverviewCatalog
		overview.Summary = fmt.Sprintf("gooo can describe %d safe catalog capabilities; this is not a language completeness claim.", len(trail.Response.Capabilities))
	case len(trail.Response.Capabilities) > 0:
		overview.Mode = CapabilityQueryOverviewMatchedOptions
		overview.Summary = fmt.Sprintf("gooo matched %d catalog capabilities; inspect each state before investing.", len(trail.Response.Capabilities))
	default:
		overview.Mode = CapabilityQueryOverviewClarification
		overview.Summary = "No explicit catalog capability matched; ask a narrower question so gooo can preserve the boundary."
	}
	overview.SuggestedQuestions = append([]string(nil), trail.NextQuestions...)
	overview.EvidenceDigest = digestCapabilityQueryOverview(overview)
	return overview
}

func (overview CapabilityQueryOverview) Validate() error {
	if err := overview.Trail.Validate(); err != nil {
		return fmt.Errorf("capability query overview trail: %w", err)
	}
	if strings.TrimSpace(overview.Summary) == "" || len(overview.SuggestedQuestions) == 0 {
		return fmt.Errorf("capability query overview is incomplete")
	}
	if len(overview.Constraints) != len(capabilityQueryOverviewConstraints) || !sameCapabilityOverviewStrings(overview.Constraints, capabilityQueryOverviewConstraints) {
		return fmt.Errorf("capability query overview constraints are incomplete")
	}
	if !validDigest(overview.EvidenceDigest) {
		return fmt.Errorf("capability query overview evidence digest is invalid")
	}
	switch overview.Mode {
	case CapabilityQueryOverviewCatalog:
		if !capabilityQueryIsOverview(strings.ToLower(strings.TrimSpace(overview.Trail.Response.Query))) || len(overview.Trail.Response.Capabilities) == 0 {
			return fmt.Errorf("catalog capability query overview is incomplete")
		}
	case CapabilityQueryOverviewMatchedOptions:
		if capabilityQueryIsOverview(strings.ToLower(strings.TrimSpace(overview.Trail.Response.Query))) || len(overview.Trail.Response.Capabilities) == 0 {
			return fmt.Errorf("matched capability query overview is incomplete")
		}
	case CapabilityQueryOverviewClarification:
		if overview.Trail.Response.Status != CapabilityQueryUnknown || len(overview.Trail.Response.Capabilities) != 0 {
			return fmt.Errorf("clarification capability query overview is incomplete")
		}
	default:
		return fmt.Errorf("capability query overview mode %q is invalid", overview.Mode)
	}
	if digestCapabilityQueryOverview(overview) != overview.EvidenceDigest {
		return fmt.Errorf("capability query overview evidence digest does not match")
	}
	return nil
}

func digestCapabilityQueryOverview(overview CapabilityQueryOverview) string {
	parts := []string{
		"gooo-capability-query-overview",
		overview.Trail.EvidenceDigest,
		string(overview.Mode),
		overview.Summary,
		strings.Join(overview.SuggestedQuestions, "|"),
		strings.Join(overview.Constraints, "|"),
	}
	return digestString(strings.Join(parts, "|"))
}

func sameCapabilityOverviewStrings(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}
