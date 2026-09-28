package decision

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"
)

// DecisionObservationDocument is the provider-neutral JSON input accepted by
// ObserveJSON. It contains only the state and typed decision observation; it
// cannot authorize execution.
type DecisionObservationDocument struct {
	Spec   DecisionSpecJSON  `json:"spec"`
	State  DecisionStateJSON `json:"state"`
	Result DecisionResultJSON `json:"result"`
}

type DecisionSpecJSON struct {
	ID             string   `json:"id"`
	Question       string   `json:"question"`
	Kind           Kind     `json:"kind"`
	AllowedChoices []string `json:"allowed_choices,omitempty"`
	Threshold      *float64 `json:"threshold,omitempty"`
	PolicyDigest   string   `json:"policy_digest"`
}

type DecisionStateJSON struct {
	Digest string `json:"digest"`
}

type DecisionValueJSON struct {
	Choice string   `json:"choice,omitempty"`
	Score  *float64 `json:"score,omitempty"`
	Noul   bool     `json:"noul,omitempty"`
}

type DecisionResultJSON struct {
	SpecID        string            `json:"spec_id"`
	Kind          Kind              `json:"kind"`
	Value         DecisionValueJSON `json:"value"`
	Confidence    *float64          `json:"confidence,omitempty"`
	EvidenceDigest string            `json:"evidence_digest"`
	Provider      string            `json:"provider"`
	Model         string            `json:"model,omitempty"`
	Status        Status            `json:"status"`
	ObservedAt    time.Time         `json:"observed_at"`
}

// DecisionReceiptJSON is the stable, non-authorizing JSON representation of a
// validated Receipt.
type DecisionReceiptJSON struct {
	Schema         string    `json:"schema"`
	SpecDigest     string    `json:"spec_digest"`
	StateDigest    string    `json:"state_digest"`
	ResultDigest   string    `json:"result_digest"`
	PolicyDigest   string    `json:"policy_digest"`
	Provider       string    `json:"provider"`
	Status         Status    `json:"status"`
	ObservedAt     time.Time `json:"observed_at"`
	NonAuthorizing bool      `json:"non_authorizing"`
	DecisionDigest string    `json:"decision_digest"`
}

// ObserveJSON decodes exactly one provider-neutral observation document and
// produces the same non-authorizing receipt as Observe. Unknown fields and
// trailing JSON are rejected so an input cannot silently change meaning.
func ObserveJSON(data []byte) (Receipt, error) {
	var document DecisionObservationDocument
	if err := decodeSingleJSON(data, &document); err != nil {
		return Receipt{}, fmt.Errorf("decode decision observation: %w", err)
	}
	return Observe(
		Spec{
			ID:             document.Spec.ID,
			Question:       document.Spec.Question,
			Kind:           document.Spec.Kind,
			AllowedChoices: document.Spec.AllowedChoices,
			Threshold:      document.Spec.Threshold,
			PolicyDigest:   document.Spec.PolicyDigest,
		},
		State{Digest: document.State.Digest},
		Result{
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
		},
	)
}

// MarshalReceiptJSON validates receipt integrity before exposing a JSON
// representation. It never emits an authorizing field or execution grant.
func MarshalReceiptJSON(receipt Receipt) ([]byte, error) {
	if err := receipt.Validate(); err != nil {
		return nil, fmt.Errorf("validate decision receipt: %w", err)
	}
	return json.MarshalIndent(DecisionReceiptJSON{
		Schema:         receipt.Schema,
		SpecDigest:     receipt.SpecDigest,
		StateDigest:    receipt.StateDigest,
		ResultDigest:   receipt.ResultDigest,
		PolicyDigest:   receipt.PolicyDigest,
		Provider:       receipt.Provider,
		Status:         receipt.Status,
		ObservedAt:     receipt.ObservedAt.UTC(),
		NonAuthorizing: receipt.NonAuthorizing,
		DecisionDigest: receipt.DecisionDigest,
	}, "", "  ")
}

func decodeSingleJSON(data []byte, destination any) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		return err
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		if err == nil {
			return errors.New("multiple JSON values are not allowed")
		}
		return fmt.Errorf("trailing JSON: %w", err)
	}
	return nil
}
