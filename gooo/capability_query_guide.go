package gooo

import "fmt"

type CapabilityQueryGuideAction string

const (
	CapabilityQueryGuideAskClarifyingQuestion  CapabilityQueryGuideAction = "ASK_CLARIFYING_QUESTION"
	CapabilityQueryGuideInspectNextOperation   CapabilityQueryGuideAction = "INSPECT_NEXT_OPERATION"
	CapabilityQueryGuideRequireExternalBoundary CapabilityQueryGuideAction = "REQUIRE_EXTERNAL_BOUNDARY"
)

type CapabilityQueryGuide struct {
	OverviewDigest  string                     `json:"overview_digest"`
	Status          CapabilityQueryState       `json:"status"`
	Action          CapabilityQueryGuideAction  `json:"action"`
	CapabilityIDs   []string                   `json:"capability_ids"`
	NextOperations  []string                   `json:"next_operations"`
	Questions       []string                   `json:"questions"`
	Constraints     []string                   `json:"constraints"`
	EvidenceDigest  string                     `json:"evidence_digest"`
}

var capabilityQueryGuideConstraints = []string{
	"catalog_scope_only",
	"unknown_is_preserved",
	"no_provider_invocation",
	"no_authorization_grant",
	"no_catalog_mutation",
	"cache_presence_is_not_semantic_evidence",
}

func DiscoverCapabilityQueryGuide(query, declaration string) CapabilityQueryGuide {
	overview := DiscoverCapabilityQueryOverview(query, declaration)
	guide := CapabilityQueryGuide{
		OverviewDigest: overview.EvidenceDigest,
		Status:         overview.Trail.Response.Status,
		Questions:      append([]string(nil), overview.SuggestedQuestions...),
		Constraints:    append([]string(nil), capabilityQueryGuideConstraints...),
	}
	for _, capability := range overview.Trail.Response.Capabilities {
		guide.CapabilityIDs = append(guide.CapabilityIDs, capability.ID)
		guide.NextOperations = append(guide.NextOperations, capability.NextOperation)
	}
	switch guide.Status {
	case CapabilityQueryUnknown:
		guide.Action = CapabilityQueryGuideAskClarifyingQuestion
	case CapabilityQueryDeferred:
		guide.Action = CapabilityQueryGuideRequireExternalBoundary
	case CapabilityQueryAvailable:
		guide.Action = CapabilityQueryGuideInspectNextOperation
	}
	guide.EvidenceDigest = digestCapabilityQueryGuide(guide)
	return guide
}

func (guide CapabilityQueryGuide) Validate() error {
	if err := validCapabilityQueryGuideDigest(guide.OverviewDigest); err != nil {
		return err
	}
	if !validDigest(guide.EvidenceDigest) {
		return fmt.Errorf("capability query guide evidence digest is invalid")
	}
	if len(guide.Constraints) != len(capabilityQueryGuideConstraints) || !sameCapabilityOverviewStrings(guide.Constraints, capabilityQueryGuideConstraints) {
		return fmt.Errorf("capability query guide constraints are incomplete")
	}
	if len(guide.Questions) == 0 {
		return fmt.Errorf("capability query guide has no questions")
	}
	if len(guide.CapabilityIDs) != len(guide.NextOperations) {
		return fmt.Errorf("capability query guide capability operations are misaligned")
	}
	switch guide.Status {
	case CapabilityQueryUnknown:
		if guide.Action != CapabilityQueryGuideAskClarifyingQuestion || len(guide.CapabilityIDs) != 0 {
			return fmt.Errorf("unknown capability query guide is incomplete")
		}
	case CapabilityQueryDeferred:
		if guide.Action != CapabilityQueryGuideRequireExternalBoundary || len(guide.CapabilityIDs) == 0 {
			return fmt.Errorf("deferred capability query guide is incomplete")
		}
	case CapabilityQueryAvailable:
		if guide.Action != CapabilityQueryGuideInspectNextOperation || len(guide.CapabilityIDs) == 0 {
			return fmt.Errorf("available capability query guide is incomplete")
		}
	default:
		return fmt.Errorf("capability query guide status %q is invalid", guide.Status)
	}
	for index := range guide.CapabilityIDs {
		if guide.CapabilityIDs[index] == "" || guide.NextOperations[index] == "" {
			return fmt.Errorf("capability query guide operation %d is incomplete", index)
		}
	}
	if digestCapabilityQueryGuide(guide) != guide.EvidenceDigest {
		return fmt.Errorf("capability query guide evidence digest does not match")
	}
	return nil
}

func validCapabilityQueryGuideDigest(value string) error {
	if !validDigest(value) {
		return fmt.Errorf("capability query guide overview digest is invalid")
	}
	return nil
}

func digestCapabilityQueryGuide(guide CapabilityQueryGuide) string {
	parts := []string{
		"gooo-capability-query-guide",
		guide.OverviewDigest,
		string(guide.Status),
		string(guide.Action),
	}
	parts = append(parts, guide.CapabilityIDs...)
	parts = append(parts, guide.NextOperations...)
	parts = append(parts, guide.Questions...)
	parts = append(parts, guide.Constraints...)
	return digestString(joinCapabilityQueryGuideParts(parts))
}

func joinCapabilityQueryGuideParts(parts []string) string {
	result := ""
	for index, part := range parts {
		if index > 0 {
			result += "|"
		}
		result += part
	}
	return result
}
