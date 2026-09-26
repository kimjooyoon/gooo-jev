package decision

import (
	"errors"
	"strings"
	"time"
)

type RuleProvider struct {
	ProviderName   string
	Values         map[string]Value
	EvidencePrefix string
	Clock          func() time.Time
}

func (provider RuleProvider) Name() string {
	if strings.TrimSpace(provider.ProviderName) == "" {
		return "rule"
	}
	return provider.ProviderName
}

func (provider RuleProvider) Observe(spec Spec, state State) (Result, error) {
	if err := spec.Validate(); err != nil {
		return Result{}, err
	}
	if strings.TrimSpace(state.Digest) == "" {
		return Result{}, errors.New("decision state digest is required")
	}
	observedAt := time.Now().UTC()
	if provider.Clock != nil {
		observedAt = provider.Clock().UTC()
	}
	if observedAt.IsZero() {
		return Result{}, errors.New("rule provider observation time is required")
	}
	result := Result{
		SpecID:         spec.ID,
		Kind:           spec.Kind,
		Provider:       provider.Name(),
		Status:         StatusUnknown,
		EvidenceDigest: provider.evidenceDigest(spec, state),
		ObservedAt:     observedAt,
	}
	if value, ok := provider.Values[spec.ID]; ok {
		result.Value = value
		result.Status = StatusObserved
	}
	if err := result.ValidateFor(spec); err != nil {
		return Result{}, err
	}
	return result, nil
}

func (provider RuleProvider) evidenceDigest(spec Spec, state State) string {
	digest, _ := Digest(struct {
		SpecID        string
		StateDigest   string
		Provider      string
		EvidencePrefix string
	}{
		SpecID:         spec.ID,
		StateDigest:    state.Digest,
		Provider:       provider.Name(),
		EvidencePrefix: provider.EvidencePrefix,
	})
	return digest
}
