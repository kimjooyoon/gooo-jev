package decision

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

const SchemaV1 = "gooo/jev-decision-receipt/v1"

type Kind string

const (
	KindChoice Kind = "choice"
	KindScore  Kind = "score"
	KindNoul   Kind = "noul"
)

type Status string

const (
	StatusObserved Status = "observed"
	StatusUnknown  Status = "unknown"
	StatusReview   Status = "review"
)

type Spec struct {
	ID             string
	Question       string
	Kind           Kind
	AllowedChoices []string
	Threshold      *float64
	PolicyDigest   string
}

type State struct {
	Digest string
}

type Value struct {
	Choice string
	Score  *float64
	Noul   bool
}

type Result struct {
	SpecID        string
	Kind          Kind
	Value         Value
	Confidence    *float64
	EvidenceDigest string
	Provider      string
	Model         string
	Status        Status
	ObservedAt    time.Time
}

type Receipt struct {
	Schema         string
	SpecDigest     string
	StateDigest    string
	ResultDigest   string
	PolicyDigest   string
	Provider       string
	Status         Status
	ObservedAt     time.Time
	NonAuthorizing bool
	DecisionDigest string
}

type Provider interface {
	Name() string
	Observe(Spec, State) (Result, error)
}

func (spec Spec) Validate() error {
	if strings.TrimSpace(spec.ID) == "" {
		return errors.New("decision spec id is required")
	}
	if strings.TrimSpace(spec.Question) == "" {
		return errors.New("decision spec question is required")
	}
	if strings.TrimSpace(spec.PolicyDigest) == "" {
		return errors.New("decision spec policy digest is required")
	}
	switch spec.Kind {
	case KindChoice:
		if len(spec.AllowedChoices) == 0 {
			return errors.New("choice decisions require allowed choices")
		}
	case KindScore:
		if spec.Threshold != nil && (*spec.Threshold < 0 || *spec.Threshold > 1) {
			return errors.New("score threshold must be between 0 and 1")
		}
	case KindNoul:
	default:
		return fmt.Errorf("unsupported decision kind %q", spec.Kind)
	}
	return nil
}

func (result Result) ValidateFor(spec Spec) error {
	if err := spec.Validate(); err != nil {
		return err
	}
	if result.SpecID != spec.ID {
		return errors.New("decision result spec id does not match the spec")
	}
	if result.Kind != spec.Kind {
		return errors.New("decision result kind does not match the spec")
	}
	if strings.TrimSpace(result.Provider) == "" {
		return errors.New("decision result provider is required")
	}
	if strings.TrimSpace(result.EvidenceDigest) == "" {
		return errors.New("decision result evidence digest is required")
	}
	if result.ObservedAt.IsZero() {
		return errors.New("decision result observation time is required")
	}
	switch result.Status {
	case StatusUnknown, StatusReview:
		return nil
	case StatusObserved:
	default:
		return fmt.Errorf("unsupported decision result status %q", result.Status)
	}
	switch spec.Kind {
	case KindChoice:
		for _, choice := range spec.AllowedChoices {
			if result.Value.Choice == choice {
				return nil
			}
		}
		return errors.New("choice result is not in the allowed set")
	case KindScore:
		if result.Value.Score == nil || *result.Value.Score < 0 || *result.Value.Score > 1 {
			return errors.New("score result must be between 0 and 1")
		}
	case KindNoul:
		if !result.Value.Noul {
			return errors.New("noul result must explicitly carry noul=true")
		}
	}
	return nil
}

func Observe(spec Spec, state State, result Result) (Receipt, error) {
	if err := spec.Validate(); err != nil {
		return Receipt{}, err
	}
	if strings.TrimSpace(state.Digest) == "" {
		return Receipt{}, errors.New("decision state digest is required")
	}
	if err := result.ValidateFor(spec); err != nil {
		return Receipt{}, err
	}
	specDigest, err := Digest(spec)
	if err != nil {
		return Receipt{}, fmt.Errorf("digest spec: %w", err)
	}
	resultDigest, err := Digest(result)
	if err != nil {
		return Receipt{}, fmt.Errorf("digest result: %w", err)
	}
	decisionDigest, err := Digest(struct {
		SpecDigest   string
		StateDigest  string
		ResultDigest string
	}{
		SpecDigest:   specDigest,
		StateDigest:  state.Digest,
		ResultDigest: resultDigest,
	})
	if err != nil {
		return Receipt{}, fmt.Errorf("digest decision: %w", err)
	}
	return Receipt{
		Schema:         SchemaV1,
		SpecDigest:     specDigest,
		StateDigest:    state.Digest,
		ResultDigest:   resultDigest,
		PolicyDigest:   spec.PolicyDigest,
		Provider:       result.Provider,
		Status:         result.Status,
		ObservedAt:     result.ObservedAt.UTC(),
		NonAuthorizing: true,
		DecisionDigest: decisionDigest,
	}, nil
}

func (receipt Receipt) Validate() error {
	if receipt.Schema != SchemaV1 {
		return errors.New("unsupported decision receipt schema")
	}
	if strings.TrimSpace(receipt.SpecDigest) == "" ||
		strings.TrimSpace(receipt.StateDigest) == "" ||
		strings.TrimSpace(receipt.ResultDigest) == "" ||
		strings.TrimSpace(receipt.PolicyDigest) == "" ||
		strings.TrimSpace(receipt.Provider) == "" ||
		strings.TrimSpace(receipt.DecisionDigest) == "" {
		return errors.New("decision receipt is incomplete")
	}
	if receipt.ObservedAt.IsZero() {
		return errors.New("decision receipt observation time is required")
	}
	if !receipt.NonAuthorizing {
		return errors.New("decision receipt must remain non-authorizing")
	}
	switch receipt.Status {
	case StatusObserved, StatusUnknown, StatusReview:
		return nil
	default:
		return fmt.Errorf("unsupported decision receipt status %q", receipt.Status)
	}
}

func Digest(value any) (string, error) {
	encoded, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(encoded)
	return hex.EncodeToString(sum[:]), nil
}
