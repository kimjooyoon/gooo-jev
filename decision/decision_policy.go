package decision

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

const DecisionPolicySchemaV1 = "gooo/jev-decision-policy/v1"

type PolicyDisposition string

const (
	PolicyEligible PolicyDisposition = "eligible"
	PolicyEscalate PolicyDisposition = "escalate"
	PolicyAbstain  PolicyDisposition = "abstain"
	PolicyUnknown  PolicyDisposition = "unknown"
)

type DecisionPolicy struct {
	PolicyDigest       string  `json:"policy_digest"`
	MinimumConfidence  float64 `json:"minimum_confidence"`
}

func NewDecisionPolicy(policyDigest string, minimumConfidence float64) (DecisionPolicy, error) {
	policy := DecisionPolicy{
		PolicyDigest:      policyDigest,
		MinimumConfidence: minimumConfidence,
	}
	if err := policy.Validate(); err != nil {
		return DecisionPolicy{}, err
	}
	return policy, nil
}

func (policy DecisionPolicy) Validate() error {
	if strings.TrimSpace(policy.PolicyDigest) == "" {
		return errors.New("decision policy digest is required")
	}
	if math.IsNaN(policy.MinimumConfidence) ||
		math.IsInf(policy.MinimumConfidence, 0) ||
		policy.MinimumConfidence < 0 ||
		policy.MinimumConfidence > 1 {
		return errors.New("decision policy confidence threshold must be between 0 and 1")
	}
	return nil
}

type PolicyOutcome struct {
	Schema         string            `json:"schema"`
	PolicyDigest   string            `json:"policy_digest"`
	DecisionDigest string            `json:"decision_digest"`
	Disposition    PolicyDisposition `json:"disposition"`
	Reason         string            `json:"reason"`
	NonAuthorizing bool              `json:"non_authorizing"`
	OutcomeDigest  string            `json:"outcome_digest"`
}

func (outcome PolicyOutcome) computeDigest() (string, error) {
	return Digest(struct {
		Schema         string
		PolicyDigest   string
		DecisionDigest string
		Disposition    PolicyDisposition
		Reason         string
		NonAuthorizing bool
	}{
		Schema:         outcome.Schema,
		PolicyDigest:   outcome.PolicyDigest,
		DecisionDigest: outcome.DecisionDigest,
		Disposition:    outcome.Disposition,
		Reason:         outcome.Reason,
		NonAuthorizing: outcome.NonAuthorizing,
	})
}

func (outcome PolicyOutcome) Validate() error {
	if outcome.Schema != DecisionPolicySchemaV1 {
		return errors.New("unsupported decision policy outcome schema")
	}
	if strings.TrimSpace(outcome.PolicyDigest) == "" ||
		strings.TrimSpace(outcome.DecisionDigest) == "" ||
		strings.TrimSpace(outcome.Reason) == "" ||
		strings.TrimSpace(outcome.OutcomeDigest) == "" {
		return errors.New("decision policy outcome is incomplete")
	}
	if !outcome.NonAuthorizing {
		return errors.New("decision policy outcome must remain non-authorizing")
	}
	switch outcome.Disposition {
	case PolicyEligible, PolicyEscalate, PolicyAbstain, PolicyUnknown:
	default:
		return fmt.Errorf("unsupported decision policy disposition %q", outcome.Disposition)
	}
	expected, err := outcome.computeDigest()
	if err != nil {
		return fmt.Errorf("digest decision policy outcome: %w", err)
	}
	if outcome.OutcomeDigest != expected {
		return errors.New("decision policy outcome digest mismatch")
	}
	return nil
}

func newPolicyOutcome(policy DecisionPolicy, receipt Receipt, disposition PolicyDisposition, reason string) (PolicyOutcome, error) {
	outcome := PolicyOutcome{
		Schema:         DecisionPolicySchemaV1,
		PolicyDigest:   policy.PolicyDigest,
		DecisionDigest: receipt.DecisionDigest,
		Disposition:    disposition,
		Reason:         reason,
		NonAuthorizing: true,
	}
	digest, err := outcome.computeDigest()
	if err != nil {
		return PolicyOutcome{}, fmt.Errorf("digest decision policy outcome: %w", err)
	}
	outcome.OutcomeDigest = digest
	if err := outcome.Validate(); err != nil {
		return PolicyOutcome{}, err
	}
	return outcome, nil
}

// EvaluateDecisionPolicy applies a deterministic confidence threshold to a
// validated decision receipt. Its result is a routing observation only: it
// never grants execution or changes source.
func EvaluateDecisionPolicy(policy DecisionPolicy, spec Spec, state State, result Result, receipt Receipt) (PolicyOutcome, error) {
	if err := policy.Validate(); err != nil {
		return PolicyOutcome{}, err
	}
	if err := spec.Validate(); err != nil {
		return PolicyOutcome{}, err
	}
	if strings.TrimSpace(state.Digest) == "" {
		return PolicyOutcome{}, errors.New("decision policy state digest is required")
	}
	if err := result.ValidateFor(spec); err != nil {
		return PolicyOutcome{}, err
	}
	if err := receipt.Validate(); err != nil {
		return PolicyOutcome{}, err
	}
	expected, err := Observe(spec, state, result)
	if err != nil {
		return PolicyOutcome{}, err
	}
	if receipt.SpecDigest != expected.SpecDigest ||
		receipt.StateDigest != expected.StateDigest ||
		receipt.ResultDigest != expected.ResultDigest ||
		receipt.DecisionDigest != expected.DecisionDigest ||
		receipt.PolicyDigest != expected.PolicyDigest ||
		receipt.PolicyDigest != policy.PolicyDigest ||
		spec.PolicyDigest != policy.PolicyDigest {
		return newPolicyOutcome(policy, receipt, PolicyUnknown, "decision-receipt-binding-mismatch")
	}
	if result.Status == StatusUnknown || receipt.Status == StatusUnknown {
		return newPolicyOutcome(policy, receipt, PolicyUnknown, "decision-unknown")
	}
	if result.Status == StatusReview || receipt.Status == StatusReview {
		return newPolicyOutcome(policy, receipt, PolicyEscalate, "decision-requires-review")
	}
	if result.Confidence == nil ||
		math.IsNaN(*result.Confidence) ||
		math.IsInf(*result.Confidence, 0) ||
		*result.Confidence < 0 ||
		*result.Confidence > 1 {
		return newPolicyOutcome(policy, receipt, PolicyUnknown, "confidence-missing-or-invalid")
	}
	if *result.Confidence < policy.MinimumConfidence {
		return newPolicyOutcome(policy, receipt, PolicyEscalate, "confidence-below-policy-threshold")
	}
	return newPolicyOutcome(policy, receipt, PolicyEligible, "confidence-at-or-above-policy-threshold")
}

func EvaluateJSONPolicy(data []byte, policy DecisionPolicy) (PolicyOutcome, error) {
	var document DecisionObservationDocument
	if err := decodeSingleJSON(data, &document); err != nil {
		return PolicyOutcome{}, fmt.Errorf("decode decision observation: %w", err)
	}
	receipt, err := ObserveJSON(data)
	if err != nil {
		return PolicyOutcome{}, err
	}
	spec := Spec{
		ID:             document.Spec.ID,
		Question:       document.Spec.Question,
		Kind:           document.Spec.Kind,
		AllowedChoices: document.Spec.AllowedChoices,
		Threshold:      document.Spec.Threshold,
		PolicyDigest:   document.Spec.PolicyDigest,
	}
	state := State{Digest: document.State.Digest}
	result := Result{
		SpecID: document.Result.SpecID,
		Kind:   document.Result.Kind,
		Value: Value{
			Choice: document.Result.Value.Choice,
			Score:  document.Result.Value.Score,
			Noul:   document.Result.Value.Noul,
		},
		Confidence:     document.Result.Confidence,
		EvidenceDigest: document.Result.EvidenceDigest,
		Provider:       document.Result.Provider,
		Model:          document.Result.Model,
		Status:         document.Result.Status,
		ObservedAt:     document.Result.ObservedAt,
	}
	return EvaluateDecisionPolicy(policy, spec, state, result, receipt)
}
