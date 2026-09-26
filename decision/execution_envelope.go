package decision

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	ExecutionReceiptSchemaV1   = "gooo/jev-execution-receipt/v1"
	ReverseObservationSchemaV1 = "gooo/jev-reverse-observation/v1"
)

type ExecutionStatus string

const (
	ExecutionCompleted   ExecutionStatus = "completed"
	ExecutionUnknown     ExecutionStatus = "unknown"
	ExecutionReview      ExecutionStatus = "review"
	ExecutionNonTerminal ExecutionStatus = "non-terminal"
)

type WorkloadIdentityStatus string

const (
	WorkloadIdentityObserved WorkloadIdentityStatus = "observed"
	WorkloadIdentityUnknown  WorkloadIdentityStatus = "unknown"
	WorkloadIdentityReview   WorkloadIdentityStatus = "review"
)

type CapabilityRequest struct {
	Subject    string
	Audience   string
	Capability string
}

type CapabilityGrant struct {
	RequestDigest            string
	CapabilityBoundaryDigest string
	WorkloadIdentityDigest   string
	EvidenceDigest           string
	Granted                  bool
	GrantDigest              string
}

type WorkloadIdentityObservation struct {
	SPIFFEID          string
	EvidenceDigest    string
	Status            WorkloadIdentityStatus
	ObservedAt        time.Time
	ObservationDigest string
}

type ExecutionReceipt struct {
	Schema                    string
	GrantDigest               string
	CapabilityBoundaryDigest string
	DecisionReceiptDigest     string
	WorkloadIdentityDigest    string
	ResultDigest              string
	Status                    ExecutionStatus
	UnknownReason             string
	ObservedAt                time.Time
	ReceiptDigest             string
}

type ReverseObservation struct {
	Schema                 string
	ExecutionReceiptDigest string
	ObservedOutputDigest   string
	VerifierDigest         string
	Status                 ProvenanceObservationStatus
	MissingStage           string
	ObservedAt             time.Time
	ObservationDigest      string
}

func (request CapabilityRequest) Validate() error {
	if strings.TrimSpace(request.Subject) == "" {
		return errors.New("capability request subject is required")
	}
	if strings.TrimSpace(request.Audience) == "" {
		return errors.New("capability request audience is required")
	}
	if strings.TrimSpace(request.Capability) == "" {
		return errors.New("capability request capability is required")
	}
	return nil
}

func (observation WorkloadIdentityObservation) Validate() error {
	if err := observation.ValidateShape(); err != nil {
		return err
	}
	expected, err := observation.computeDigest()
	if err != nil {
		return fmt.Errorf("digest workload identity: %w", err)
	}
	if observation.ObservationDigest != expected {
		return errors.New("workload identity observation digest mismatch")
	}
	return nil
}

func (observation WorkloadIdentityObservation) ValidateShape() error {
	if strings.TrimSpace(observation.SPIFFEID) == "" {
		return errors.New("workload identity SPIFFE ID is required")
	}
	if strings.TrimSpace(observation.EvidenceDigest) == "" {
		return errors.New("workload identity evidence digest is required")
	}
	if observation.ObservedAt.IsZero() {
		return errors.New("workload identity observation time is required")
	}
	switch observation.Status {
	case WorkloadIdentityObserved, WorkloadIdentityUnknown, WorkloadIdentityReview:
		return nil
	default:
		return fmt.Errorf("unsupported workload identity status %q", observation.Status)
	}
}

func (observation WorkloadIdentityObservation) computeDigest() (string, error) {
	return Digest(struct {
		SPIFFEID       string
		EvidenceDigest string
		Status         WorkloadIdentityStatus
		ObservedAt     time.Time
	}{
		SPIFFEID:       observation.SPIFFEID,
		EvidenceDigest: observation.EvidenceDigest,
		Status:         observation.Status,
		ObservedAt:     observation.ObservedAt.UTC(),
	})
}

func NewWorkloadIdentityObservation(spiffeID, evidenceDigest string, status WorkloadIdentityStatus, observedAt time.Time) (WorkloadIdentityObservation, error) {
	observation := WorkloadIdentityObservation{
		SPIFFEID:       spiffeID,
		EvidenceDigest: evidenceDigest,
		Status:         status,
		ObservedAt:     observedAt.UTC(),
	}
	if err := observation.ValidateShape(); err != nil {
		return WorkloadIdentityObservation{}, err
	}
	digest, err := observation.computeDigest()
	if err != nil {
		return WorkloadIdentityObservation{}, fmt.Errorf("digest workload identity: %w", err)
	}
	observation.ObservationDigest = digest
	return observation, nil
}

func NewCapabilityGrant(request CapabilityRequest, boundary CapabilityBoundary, identity WorkloadIdentityObservation, evidenceDigest string) (CapabilityGrant, error) {
	if err := request.Validate(); err != nil {
		return CapabilityGrant{}, err
	}
	if err := boundary.Validate(); err != nil {
		return CapabilityGrant{}, err
	}
	if request.Subject != boundary.Subject || request.Audience != boundary.Audience || request.Capability != boundary.Capability {
		return CapabilityGrant{}, errors.New("capability request does not match boundary")
	}
	if err := identity.Validate(); err != nil {
		return CapabilityGrant{}, err
	}
	if identity.Status != WorkloadIdentityObserved {
		return CapabilityGrant{}, errors.New("capability grant requires observed workload identity")
	}
	if strings.TrimSpace(evidenceDigest) == "" {
		return CapabilityGrant{}, errors.New("capability grant evidence digest is required")
	}
	requestDigest, err := Digest(request)
	if err != nil {
		return CapabilityGrant{}, fmt.Errorf("digest capability request: %w", err)
	}
	grant := CapabilityGrant{
		RequestDigest:            requestDigest,
		CapabilityBoundaryDigest: boundary.BoundaryDigest,
		WorkloadIdentityDigest:   identity.ObservationDigest,
		EvidenceDigest:           evidenceDigest,
		Granted:                  true,
	}
	grant.GrantDigest, err = grant.computeDigest()
	if err != nil {
		return CapabilityGrant{}, fmt.Errorf("digest capability grant: %w", err)
	}
	return grant, nil
}

func (grant CapabilityGrant) computeDigest() (string, error) {
	return Digest(struct {
		RequestDigest            string
		CapabilityBoundaryDigest string
		WorkloadIdentityDigest   string
		EvidenceDigest           string
		Granted                  bool
	}{
		RequestDigest:            grant.RequestDigest,
		CapabilityBoundaryDigest: grant.CapabilityBoundaryDigest,
		WorkloadIdentityDigest:   grant.WorkloadIdentityDigest,
		EvidenceDigest:           grant.EvidenceDigest,
		Granted:                  grant.Granted,
	})
}

func (grant CapabilityGrant) Validate() error {
	if strings.TrimSpace(grant.RequestDigest) == "" ||
		strings.TrimSpace(grant.CapabilityBoundaryDigest) == "" ||
		strings.TrimSpace(grant.WorkloadIdentityDigest) == "" ||
		strings.TrimSpace(grant.EvidenceDigest) == "" ||
		strings.TrimSpace(grant.GrantDigest) == "" {
		return errors.New("capability grant is incomplete")
	}
	if !grant.Granted {
		return errors.New("capability grant is not granted")
	}
	digest, err := grant.computeDigest()
	if err != nil {
		return fmt.Errorf("digest capability grant: %w", err)
	}
	if grant.GrantDigest != digest {
		return errors.New("capability grant digest mismatch")
	}
	return nil
}

func NewExecutionReceipt(grant CapabilityGrant, boundary CapabilityBoundary, identity WorkloadIdentityObservation, decisionReceipt Receipt, resultDigest string, status ExecutionStatus, observedAt time.Time) (ExecutionReceipt, error) {
	receipt := ExecutionReceipt{
		Schema:                    ExecutionReceiptSchemaV1,
		GrantDigest:               grant.GrantDigest,
		CapabilityBoundaryDigest: boundary.BoundaryDigest,
		DecisionReceiptDigest:     decisionReceipt.DecisionDigest,
		WorkloadIdentityDigest:    identity.ObservationDigest,
		ResultDigest:              resultDigest,
		Status:                    status,
		ObservedAt:                observedAt.UTC(),
	}
	if err := grant.Validate(); err != nil {
		receipt.Status = ExecutionUnknown
		receipt.UnknownReason = "missing-or-invalid-grant"
	} else if err := boundary.Validate(); err != nil || boundary.StatusAt(observedAt) != CapabilityBoundaryActive {
		receipt.Status = ExecutionUnknown
		receipt.UnknownReason = "inactive-or-invalid-capability-boundary"
	} else if err := identity.Validate(); err != nil || identity.Status != WorkloadIdentityObserved {
		receipt.Status = ExecutionUnknown
		receipt.UnknownReason = "missing-or-invalid-workload-identity"
	} else if err := decisionReceipt.Validate(); err != nil {
		receipt.Status = ExecutionUnknown
		receipt.UnknownReason = "missing-or-invalid-decision-receipt"
	} else if strings.TrimSpace(resultDigest) == "" {
		receipt.Status = ExecutionUnknown
		receipt.UnknownReason = "missing-terminal-result"
	} else if status != ExecutionCompleted {
		receipt.Status = ExecutionUnknown
		receipt.UnknownReason = "execution-not-terminal"
	}
	if err := receipt.validateShape(); err != nil {
		return ExecutionReceipt{}, err
	}
	digest, err := receipt.computeDigest()
	if err != nil {
		return ExecutionReceipt{}, fmt.Errorf("digest execution receipt: %w", err)
	}
	receipt.ReceiptDigest = digest
	return receipt, nil
}

func (receipt ExecutionReceipt) Validate() error {
	if err := receipt.validateShape(); err != nil {
		return err
	}
	expected, err := receipt.computeDigest()
	if err != nil {
		return fmt.Errorf("digest execution receipt: %w", err)
	}
	if receipt.ReceiptDigest != expected {
		return errors.New("execution receipt digest mismatch")
	}
	return nil
}

func (receipt ExecutionReceipt) validateShape() error {
	if receipt.Schema != ExecutionReceiptSchemaV1 {
		return errors.New("unsupported execution receipt schema")
	}
	if receipt.ObservedAt.IsZero() {
		return errors.New("execution receipt observation time is required")
	}
	switch receipt.Status {
	case ExecutionCompleted:
		if strings.TrimSpace(receipt.GrantDigest) == "" ||
			strings.TrimSpace(receipt.CapabilityBoundaryDigest) == "" ||
			strings.TrimSpace(receipt.DecisionReceiptDigest) == "" ||
			strings.TrimSpace(receipt.WorkloadIdentityDigest) == "" ||
			strings.TrimSpace(receipt.ResultDigest) == "" {
			return errors.New("completed execution receipt is incomplete")
		}
	case ExecutionUnknown, ExecutionReview, ExecutionNonTerminal:
		if strings.TrimSpace(receipt.UnknownReason) == "" {
			return errors.New("non-completed execution receipt requires an explicit reason")
		}
	default:
		return fmt.Errorf("unsupported execution receipt status %q", receipt.Status)
	}
	if strings.TrimSpace(receipt.ReceiptDigest) == "" {
		return errors.New("execution receipt digest is required")
	}
	return nil
}

func (receipt ExecutionReceipt) computeDigest() (string, error) {
	return Digest(struct {
		Schema                    string
		GrantDigest               string
		CapabilityBoundaryDigest string
		DecisionReceiptDigest     string
		WorkloadIdentityDigest    string
		ResultDigest              string
		Status                    ExecutionStatus
		UnknownReason             string
		ObservedAt                time.Time
	}{
		Schema:                    receipt.Schema,
		GrantDigest:               receipt.GrantDigest,
		CapabilityBoundaryDigest: receipt.CapabilityBoundaryDigest,
		DecisionReceiptDigest:     receipt.DecisionReceiptDigest,
		WorkloadIdentityDigest:    receipt.WorkloadIdentityDigest,
		ResultDigest:              receipt.ResultDigest,
		Status:                    receipt.Status,
		UnknownReason:             receipt.UnknownReason,
		ObservedAt:                receipt.ObservedAt.UTC(),
	})
}

func NewReverseObservation(receipt ExecutionReceipt, observedOutputDigest, verifierDigest, missingStage string, status ProvenanceObservationStatus, observedAt time.Time) (ReverseObservation, error) {
	observation := ReverseObservation{
		Schema:                 ReverseObservationSchemaV1,
		ExecutionReceiptDigest: receipt.ReceiptDigest,
		ObservedOutputDigest:   observedOutputDigest,
		VerifierDigest:         verifierDigest,
		Status:                 status,
		MissingStage:           missingStage,
		ObservedAt:             observedAt.UTC(),
	}
	if err := receipt.Validate(); err != nil {
		observation.Status = ProvenanceUnknown
		observation.MissingStage = "execution-receipt"
	} else if receipt.Status != ExecutionCompleted {
		observation.Status = ProvenanceUnknown
		observation.MissingStage = "terminal-execution"
	} else if strings.TrimSpace(observedOutputDigest) == "" {
		observation.Status = ProvenanceUnknown
		observation.MissingStage = "reverse-output"
	} else if strings.TrimSpace(verifierDigest) == "" {
		observation.Status = ProvenanceUnknown
		observation.MissingStage = "reverse-verifier"
	}
	if err := observation.validateShape(); err != nil {
		return ReverseObservation{}, err
	}
	digest, err := observation.computeDigest()
	if err != nil {
		return ReverseObservation{}, fmt.Errorf("digest reverse observation: %w", err)
	}
	observation.ObservationDigest = digest
	return observation, nil
}

func (observation ReverseObservation) Validate() error {
	if err := observation.validateShape(); err != nil {
		return err
	}
	expected, err := observation.computeDigest()
	if err != nil {
		return fmt.Errorf("digest reverse observation: %w", err)
	}
	if observation.ObservationDigest != expected {
		return errors.New("reverse observation digest mismatch")
	}
	return nil
}

func (observation ReverseObservation) validateShape() error {
	if observation.Schema != ReverseObservationSchemaV1 {
		return errors.New("unsupported reverse observation schema")
	}
	if observation.ObservedAt.IsZero() {
		return errors.New("reverse observation time is required")
	}
	if strings.TrimSpace(observation.ExecutionReceiptDigest) == "" {
		return errors.New("reverse observation execution receipt digest is required")
	}
	switch observation.Status {
	case ProvenanceObserved, ProvenanceUnknown, ProvenanceReview:
	default:
		return fmt.Errorf("unsupported reverse observation status %q", observation.Status)
	}
	if observation.Status == ProvenanceObserved {
		if strings.TrimSpace(observation.ObservedOutputDigest) == "" || strings.TrimSpace(observation.VerifierDigest) == "" {
			return errors.New("observed reverse observation requires output and verifier digests")
		}
	} else if strings.TrimSpace(observation.MissingStage) == "" {
		return errors.New("unknown or review reverse observation requires missing stage")
	}
	if strings.TrimSpace(observation.ObservationDigest) == "" {
		return errors.New("reverse observation digest is required")
	}
	return nil
}

func (observation ReverseObservation) computeDigest() (string, error) {
	return Digest(struct {
		Schema                 string
		ExecutionReceiptDigest string
		ObservedOutputDigest   string
		VerifierDigest         string
		Status                 ProvenanceObservationStatus
		MissingStage           string
		ObservedAt             time.Time
	}{
		Schema:                 observation.Schema,
		ExecutionReceiptDigest: observation.ExecutionReceiptDigest,
		ObservedOutputDigest:   observation.ObservedOutputDigest,
		VerifierDigest:         observation.VerifierDigest,
		Status:                 observation.Status,
		MissingStage:           observation.MissingStage,
		ObservedAt:             observation.ObservedAt.UTC(),
	})
}

