package decision

import (
	"errors"
	"fmt"
	"strings"
)

const DecisionBatchSchemaV1 = "gooo/jev-decision-batch/v1"

type BatchObservation struct {
	Spec   Spec
	Result Result
}

type DecisionBatchSpec struct {
	ID             string   `json:"id"`
	Question       string   `json:"question"`
	Kind           Kind     `json:"kind"`
	AllowedChoices []string `json:"allowed_choices"`
	Threshold      *float64 `json:"threshold"`
	PolicyDigest   string   `json:"policy_digest"`
}

type DecisionBatchReceipt struct {
	Schema         string              `json:"schema"`
	StateDigest    string              `json:"state_digest"`
	Specs          []DecisionBatchSpec `json:"specs"`
	SpecIDs        []string            `json:"spec_ids"`
	Receipts       []Receipt           `json:"receipts"`
	NonAuthorizing bool                `json:"non_authorizing"`
	BatchDigest    string              `json:"batch_digest"`
}

type DecisionBatchObservationJSON struct {
	Spec   DecisionSpecJSON   `json:"spec"`
	Result DecisionResultJSON `json:"result"`
}

type DecisionBatchDocumentJSON struct {
	State        DecisionStateJSON              `json:"state"`
	Observations []DecisionBatchObservationJSON `json:"observations"`
}

func ObserveBatch(state State, observations []BatchObservation) (DecisionBatchReceipt, error) {
	if strings.TrimSpace(state.Digest) == "" {
		return DecisionBatchReceipt{}, errors.New("decision batch state digest is required")
	}
	if len(observations) == 0 {
		return DecisionBatchReceipt{}, errors.New("decision batch requires at least one observation")
	}

	receipt := DecisionBatchReceipt{
		Schema:         DecisionBatchSchemaV1,
		StateDigest:    state.Digest,
		NonAuthorizing: true,
	}
	seen := make(map[string]struct{}, len(observations))
	for _, observation := range observations {
		if err := observation.Spec.Validate(); err != nil {
			return DecisionBatchReceipt{}, err
		}
		if _, exists := seen[observation.Spec.ID]; exists {
			return DecisionBatchReceipt{}, fmt.Errorf("decision batch contains duplicate spec id %q", observation.Spec.ID)
		}
		seen[observation.Spec.ID] = struct{}{}
		observed, err := Observe(observation.Spec, state, observation.Result)
		if err != nil {
			return DecisionBatchReceipt{}, err
		}
		receipt.Specs = append(receipt.Specs, batchSpec(observation.Spec))
		receipt.SpecIDs = append(receipt.SpecIDs, observation.Spec.ID)
		receipt.Receipts = append(receipt.Receipts, observed)
	}
	digest, err := receipt.computeDigest()
	if err != nil {
		return DecisionBatchReceipt{}, fmt.Errorf("digest decision batch: %w", err)
	}
	receipt.BatchDigest = digest
	if err := receipt.Validate(); err != nil {
		return DecisionBatchReceipt{}, err
	}
	return receipt, nil
}

func ObserveBatchJSON(data []byte) (DecisionBatchReceipt, error) {
	var document DecisionBatchDocumentJSON
	if err := decodeSingleJSON(data, &document); err != nil {
		return DecisionBatchReceipt{}, fmt.Errorf("decode decision batch: %w", err)
	}
	observations := make([]BatchObservation, 0, len(document.Observations))
	for _, item := range document.Observations {
		observations = append(observations, BatchObservation{
			Spec: Spec{
				ID:             item.Spec.ID,
				Question:       item.Spec.Question,
				Kind:           item.Spec.Kind,
				AllowedChoices: item.Spec.AllowedChoices,
				Threshold:      item.Spec.Threshold,
				PolicyDigest:   item.Spec.PolicyDigest,
			},
			Result: Result{
				SpecID: item.Result.SpecID,
				Kind:   item.Result.Kind,
				Value: Value{
					Choice: item.Result.Value.Choice,
					Score:  item.Result.Value.Score,
					Noul:   item.Result.Value.Noul,
				},
				Confidence:     item.Result.Confidence,
				EvidenceDigest: item.Result.EvidenceDigest,
				Provider:       item.Result.Provider,
				Model:          item.Result.Model,
				Status:         item.Result.Status,
				ObservedAt:     item.Result.ObservedAt,
			},
		})
	}
	return ObserveBatch(State{Digest: document.State.Digest}, observations)
}

func (receipt DecisionBatchReceipt) computeDigest() (string, error) {
	return Digest(struct {
		Schema         string
		StateDigest    string
		Specs          []DecisionBatchSpec
		SpecIDs        []string
		Receipts       []Receipt
		NonAuthorizing bool
	}{
		Schema:         receipt.Schema,
		StateDigest:    receipt.StateDigest,
		Specs:          receipt.Specs,
		SpecIDs:        receipt.SpecIDs,
		Receipts:       receipt.Receipts,
		NonAuthorizing: receipt.NonAuthorizing,
	})
}

func (receipt DecisionBatchReceipt) Validate() error {
	if receipt.Schema != DecisionBatchSchemaV1 {
		return errors.New("unsupported decision batch schema")
	}
	if strings.TrimSpace(receipt.StateDigest) == "" {
		return errors.New("decision batch state digest is required")
	}
	if len(receipt.SpecIDs) == 0 || len(receipt.SpecIDs) != len(receipt.Specs) || len(receipt.SpecIDs) != len(receipt.Receipts) {
		return errors.New("decision batch spec and receipt counts must match")
	}
	if !receipt.NonAuthorizing {
		return errors.New("decision batch must remain non-authorizing")
	}
	seen := make(map[string]struct{}, len(receipt.SpecIDs))
	for index, specID := range receipt.SpecIDs {
		if strings.TrimSpace(specID) == "" {
			return errors.New("decision batch spec id is required")
		}
		if _, exists := seen[specID]; exists {
			return fmt.Errorf("decision batch contains duplicate spec id %q", specID)
		}
		seen[specID] = struct{}{}
		spec := receipt.Specs[index].model()
		if spec.ID != specID {
			return errors.New("decision batch receipt spec id mismatch")
		}
		if err := spec.Validate(); err != nil {
			return fmt.Errorf("validate decision batch spec %d: %w", index, err)
		}
		specDigest, err := Digest(spec)
		if err != nil {
			return fmt.Errorf("digest decision batch spec %d: %w", index, err)
		}
		if receipt.Receipts[index].SpecDigest != specDigest || receipt.Receipts[index].PolicyDigest != spec.PolicyDigest {
			return errors.New("decision batch receipt does not match its spec")
		}
		if err := receipt.Receipts[index].Validate(); err != nil {
			return fmt.Errorf("validate decision batch receipt %d: %w", index, err)
		}
		if receipt.Receipts[index].StateDigest != receipt.StateDigest {
			return errors.New("decision batch receipt state digest mismatch")
		}
	}
	if strings.TrimSpace(receipt.BatchDigest) == "" {
		return errors.New("decision batch digest is required")
	}
	expected, err := receipt.computeDigest()
	if err != nil {
		return fmt.Errorf("digest decision batch: %w", err)
	}
	if receipt.BatchDigest != expected {
		return errors.New("decision batch digest mismatch")
	}
	return nil
}

func batchSpec(spec Spec) DecisionBatchSpec {
	var allowedChoices []string
	if spec.AllowedChoices != nil {
		allowedChoices = make([]string, len(spec.AllowedChoices))
		copy(allowedChoices, spec.AllowedChoices)
	}
	var threshold *float64
	if spec.Threshold != nil {
		value := *spec.Threshold
		threshold = &value
	}
	return DecisionBatchSpec{
		ID: spec.ID, Question: spec.Question, Kind: spec.Kind,
		AllowedChoices: allowedChoices, Threshold: threshold,
		PolicyDigest: spec.PolicyDigest,
	}
}

func (spec DecisionBatchSpec) model() Spec {
	return Spec{
		ID: spec.ID, Question: spec.Question, Kind: spec.Kind,
		AllowedChoices: spec.AllowedChoices, Threshold: spec.Threshold,
		PolicyDigest: spec.PolicyDigest,
	}
}
