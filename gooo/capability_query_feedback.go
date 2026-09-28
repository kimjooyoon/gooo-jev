package gooo

import (
	"fmt"
	"sort"
	"strings"
)

type CapabilityQueryFeedbackDisposition string

const (
	CapabilityQueryFeedbackPendingEvidence     CapabilityQueryFeedbackDisposition = "PENDING_REQUIRED_EVIDENCE"
	CapabilityQueryFeedbackPreserveUnknown     CapabilityQueryFeedbackDisposition = "PRESERVE_UNKNOWN"
	CapabilityQueryFeedbackPreserveDeferred    CapabilityQueryFeedbackDisposition = "PRESERVE_DEFERRED"
	CapabilityQueryFeedbackRequiresDeclaration CapabilityQueryFeedbackDisposition = "REQUIRES_DECLARATION_BINDING"
	CapabilityQueryFeedbackEligibleForReview   CapabilityQueryFeedbackDisposition = "ELIGIBLE_FOR_CATALOG_REVIEW"
)

var capabilityQueryFeedbackRequiredStages = []string{
	"declaration",
	"identity",
	"generation",
	"boundary",
	"receipt",
	"reverse_observation",
}

// CapabilityQueryFeedback is an evidence-bound observation for the language's
// self-improvement loop. It never promotes a capability by itself.
type CapabilityQueryFeedback struct {
	Query                   string                             `json:"query"`
	ResponseDigest          string                             `json:"response_digest"`
	DeclarationSourceDigest string                             `json:"declaration_source_digest,omitempty"`
	CapabilityIDs           []string                           `json:"capability_ids"`
	VerifiedStages          []string                           `json:"verified_stages"`
	MissingStages           []string                           `json:"missing_stages"`
	Disposition              CapabilityQueryFeedbackDisposition `json:"disposition"`
	EvidenceDigest          string                             `json:"evidence_digest"`
	NonExecuting            bool                               `json:"non_executing"`
	NonAuthorizing          bool                               `json:"non_authorizing"`
}

// ObserveCapabilityQueryFeedback records externally verified evidence for a
// capability question. The caller must supply stages observed outside this
// function; this function does not execute, authorize, or infer them.
func ObserveCapabilityQueryFeedback(response CapabilityQueryResponse, verifiedStages []string) (CapabilityQueryFeedback, error) {
	if err := response.Validate(); err != nil {
		return CapabilityQueryFeedback{}, fmt.Errorf("capability query feedback response: %w", err)
	}
	stages, err := normalizeCapabilityQueryFeedbackStages(verifiedStages)
	if err != nil {
		return CapabilityQueryFeedback{}, err
	}
	capabilityIDs := make([]string, 0, len(response.Capabilities))
	for _, capability := range response.Capabilities {
		capabilityIDs = append(capabilityIDs, capability.ID)
	}
	sort.Strings(capabilityIDs)
	missing := missingCapabilityQueryFeedbackStages(response, stages)
	disposition := CapabilityQueryFeedbackPendingEvidence
	switch response.Status {
	case CapabilityQueryUnknown:
		disposition = CapabilityQueryFeedbackPreserveUnknown
	case CapabilityQueryDeferred:
		disposition = CapabilityQueryFeedbackPreserveDeferred
	case CapabilityQueryAvailable:
		switch {
		case response.Declaration == nil:
			disposition = CapabilityQueryFeedbackRequiresDeclaration
		case !response.Declaration.Bound:
			disposition = CapabilityQueryFeedbackRequiresDeclaration
		case len(missing) == 0:
			disposition = CapabilityQueryFeedbackEligibleForReview
		}
	}
	feedback := CapabilityQueryFeedback{
		Query:           response.Query,
		ResponseDigest:  response.QueryDigest,
		CapabilityIDs:   capabilityIDs,
		VerifiedStages:  stages,
		MissingStages:   missing,
		Disposition:     disposition,
		NonExecuting:   true,
		NonAuthorizing: true,
	}
	if response.Declaration != nil {
		feedback.DeclarationSourceDigest = response.Declaration.SourceDigest
	}
	feedback.EvidenceDigest = digestCapabilityQueryFeedback(feedback)
	return feedback, nil
}

func normalizeCapabilityQueryFeedbackStages(stages []string) ([]string, error) {
	allowed := make(map[string]struct{}, len(capabilityQueryFeedbackRequiredStages))
	for _, stage := range capabilityQueryFeedbackRequiredStages {
		allowed[stage] = struct{}{}
	}
	seen := make(map[string]struct{}, len(stages))
	normalized := make([]string, 0, len(stages))
	for _, raw := range stages {
		stage := strings.ToLower(strings.TrimSpace(raw))
		if stage == "" {
			continue
		}
		if _, ok := allowed[stage]; !ok {
			return nil, fmt.Errorf("capability query feedback stage %q is unknown", stage)
		}
		if _, ok := seen[stage]; ok {
			continue
		}
		seen[stage] = struct{}{}
		normalized = append(normalized, stage)
	}
	sort.Strings(normalized)
	return normalized, nil
}

func missingCapabilityQueryFeedbackStages(response CapabilityQueryResponse, verified []string) []string {
	seen := make(map[string]struct{}, len(verified))
	for _, stage := range verified {
		seen[stage] = struct{}{}
	}
	missing := make([]string, 0, len(capabilityQueryFeedbackRequiredStages))
	for _, stage := range capabilityQueryFeedbackRequiredStages {
		if stage == "declaration" && response.Declaration != nil && response.Declaration.Bound {
			seen[stage] = struct{}{}
		}
		if _, ok := seen[stage]; !ok {
			missing = append(missing, stage)
		}
	}
	return missing
}

func (feedback CapabilityQueryFeedback) Validate() error {
	if strings.TrimSpace(feedback.Query) == "" {
		return fmt.Errorf("capability query feedback has no query")
	}
	if !validDigest(feedback.ResponseDigest) || !validDigest(feedback.EvidenceDigest) {
		return fmt.Errorf("capability query feedback digest is invalid")
	}
	if !feedback.NonExecuting || !feedback.NonAuthorizing {
		return fmt.Errorf("capability query feedback crossed an execution or authorization boundary")
	}
	if feedback.Disposition != CapabilityQueryFeedbackPendingEvidence && feedback.Disposition != CapabilityQueryFeedbackPreserveUnknown && feedback.Disposition != CapabilityQueryFeedbackPreserveDeferred && feedback.Disposition != CapabilityQueryFeedbackRequiresDeclaration && feedback.Disposition != CapabilityQueryFeedbackEligibleForReview {
		return fmt.Errorf("capability query feedback disposition %q is invalid", feedback.Disposition)
	}
	if feedback.Disposition == CapabilityQueryFeedbackEligibleForReview && len(feedback.MissingStages) != 0 {
		return fmt.Errorf("catalog review feedback still has missing stages")
	}
	if digestCapabilityQueryFeedback(feedback) != feedback.EvidenceDigest {
		return fmt.Errorf("capability query feedback evidence digest does not match its evidence")
	}
	return nil
}

func digestCapabilityQueryFeedback(feedback CapabilityQueryFeedback) string {
	parts := []string{
		"gooo-capability-query-feedback",
		strings.ToLower(strings.TrimSpace(feedback.Query)),
		feedback.ResponseDigest,
		feedback.DeclarationSourceDigest,
		string(feedback.Disposition),
	}
	parts = append(parts, feedback.CapabilityIDs...)
	parts = append(parts, feedback.VerifiedStages...)
	parts = append(parts, feedback.MissingStages...)
	return digestString(strings.Join(parts, "|"))
}
