package decision

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const ExecutionLifecycleAuthorizationSchemaV1 = "gooo/jev-execution-lifecycle-authorization/v1"

type ExecutionLifecycleAuthorizationStatus string

const (
	ExecutionLifecycleAuthorized          ExecutionLifecycleAuthorizationStatus = "authorized"
	ExecutionLifecycleAuthorizationUnknown ExecutionLifecycleAuthorizationStatus = "unknown"
)

type ExecutionLifecycleAuthorization struct {
	Schema                   string
	ReplayDigest             string
	ResumptionDigest         string
	GrantDigest              string
	CapabilityBoundaryDigest string
	WorkloadIdentityDigest   string
	DecisionReceiptDigest    string
	Status                   ExecutionLifecycleAuthorizationStatus
	MissingStage             string
	ObservedAt               time.Time
	AuthorizationDigest      string
}

func NewExecutionLifecycleAuthorization(
	replay ExecutionLifecycleReplayObservation,
	resumption ExecutionLifecycleReceipt,
	grant CapabilityGrant,
	boundary CapabilityBoundary,
	identity WorkloadIdentityObservation,
	decisionReceipt Receipt,
	observedAt time.Time,
) (ExecutionLifecycleAuthorization, error) {
	if observedAt.IsZero() {
		return ExecutionLifecycleAuthorization{}, errors.New("execution lifecycle authorization observation time is required")
	}
	authorization := ExecutionLifecycleAuthorization{
		Schema:                   ExecutionLifecycleAuthorizationSchemaV1,
		ReplayDigest:             replay.ReplayDigest,
		ResumptionDigest:         resumption.LifecycleDigest,
		GrantDigest:              grant.GrantDigest,
		CapabilityBoundaryDigest: boundary.BoundaryDigest,
		WorkloadIdentityDigest:   identity.ObservationDigest,
		DecisionReceiptDigest:    decisionReceipt.DecisionDigest,
		Status:                   ExecutionLifecycleAuthorized,
		ObservedAt:               observedAt.UTC(),
	}
	switch {
	case replay.Validate() != nil:
		authorization.Status = ExecutionLifecycleAuthorizationUnknown
		authorization.MissingStage = "replay"
	case replay.Status != ExecutionLifecycleReplayObserved:
		authorization.Status = ExecutionLifecycleAuthorizationUnknown
		authorization.MissingStage = "replay-status"
	case resumption.Validate() != nil:
		authorization.Status = ExecutionLifecycleAuthorizationUnknown
		authorization.MissingStage = "resumption"
	case resumption.Status != ExecutionLifecycleResumed:
		authorization.Status = ExecutionLifecycleAuthorizationUnknown
		authorization.MissingStage = "resumption-status"
	case replay.ResumptionDigest != resumption.LifecycleDigest ||
		replay.SuspensionDigest != resumption.ParentLifecycleDigest:
		authorization.Status = ExecutionLifecycleAuthorizationUnknown
		authorization.MissingStage = "lifecycle-binding"
	case grant.Validate() != nil:
		authorization.Status = ExecutionLifecycleAuthorizationUnknown
		authorization.MissingStage = "grant"
	case boundary.Validate() != nil || boundary.StatusAt(observedAt) != CapabilityBoundaryActive:
		authorization.Status = ExecutionLifecycleAuthorizationUnknown
		authorization.MissingStage = "capability-boundary"
	case identity.Validate() != nil || identity.Status != WorkloadIdentityObserved:
		authorization.Status = ExecutionLifecycleAuthorizationUnknown
		authorization.MissingStage = "workload-identity"
	case decisionReceipt.Validate() != nil:
		authorization.Status = ExecutionLifecycleAuthorizationUnknown
		authorization.MissingStage = "decision-receipt"
	case resumption.GrantDigest != grant.GrantDigest ||
		resumption.CapabilityBoundaryDigest != boundary.BoundaryDigest ||
		resumption.WorkloadIdentityDigest != identity.ObservationDigest ||
		resumption.DecisionReceiptDigest != decisionReceipt.DecisionDigest:
		authorization.Status = ExecutionLifecycleAuthorizationUnknown
		authorization.MissingStage = "authorization-binding"
	}
	if err := authorization.assignDigest(); err != nil {
		return ExecutionLifecycleAuthorization{}, fmt.Errorf("digest execution lifecycle authorization: %w", err)
	}
	if err := authorization.Validate(); err != nil {
		return ExecutionLifecycleAuthorization{}, err
	}
	return authorization, nil
}

func (authorization *ExecutionLifecycleAuthorization) assignDigest() error {
	digest, err := authorization.computeDigest()
	if err != nil {
		return err
	}
	authorization.AuthorizationDigest = digest
	return nil
}

func (authorization ExecutionLifecycleAuthorization) Validate() error {
	if err := authorization.validateShape(); err != nil {
		return err
	}
	expected, err := authorization.computeDigest()
	if err != nil {
		return fmt.Errorf("digest execution lifecycle authorization: %w", err)
	}
	if authorization.AuthorizationDigest != expected {
		return errors.New("execution lifecycle authorization digest mismatch")
	}
	return nil
}

func (authorization ExecutionLifecycleAuthorization) validateShape() error {
	if authorization.Schema != ExecutionLifecycleAuthorizationSchemaV1 {
		return errors.New("unsupported execution lifecycle authorization schema")
	}
	if authorization.ObservedAt.IsZero() {
		return errors.New("execution lifecycle authorization observation time is required")
	}
	if strings.TrimSpace(authorization.AuthorizationDigest) == "" {
		return errors.New("execution lifecycle authorization digest is required")
	}
	switch authorization.Status {
	case ExecutionLifecycleAuthorized:
		if strings.TrimSpace(authorization.ReplayDigest) == "" ||
			strings.TrimSpace(authorization.ResumptionDigest) == "" ||
			strings.TrimSpace(authorization.GrantDigest) == "" ||
			strings.TrimSpace(authorization.CapabilityBoundaryDigest) == "" ||
			strings.TrimSpace(authorization.WorkloadIdentityDigest) == "" ||
			strings.TrimSpace(authorization.DecisionReceiptDigest) == "" ||
			strings.TrimSpace(authorization.MissingStage) != "" {
			return errors.New("authorized execution lifecycle is incomplete")
		}
	case ExecutionLifecycleAuthorizationUnknown:
		if strings.TrimSpace(authorization.MissingStage) == "" {
			return errors.New("unknown execution lifecycle authorization requires a missing stage")
		}
	default:
		return fmt.Errorf("unsupported execution lifecycle authorization status %q", authorization.Status)
	}
	return nil
}

func (authorization ExecutionLifecycleAuthorization) computeDigest() (string, error) {
	return Digest(struct {
		Schema                   string
		ReplayDigest             string
		ResumptionDigest         string
		GrantDigest              string
		CapabilityBoundaryDigest string
		WorkloadIdentityDigest   string
		DecisionReceiptDigest    string
		Status                   ExecutionLifecycleAuthorizationStatus
		MissingStage             string
		ObservedAt               time.Time
	}{
		Schema:                   authorization.Schema,
		ReplayDigest:             authorization.ReplayDigest,
		ResumptionDigest:         authorization.ResumptionDigest,
		GrantDigest:              authorization.GrantDigest,
		CapabilityBoundaryDigest: authorization.CapabilityBoundaryDigest,
		WorkloadIdentityDigest:   authorization.WorkloadIdentityDigest,
		DecisionReceiptDigest:    authorization.DecisionReceiptDigest,
		Status:                   authorization.Status,
		MissingStage:             authorization.MissingStage,
		ObservedAt:               authorization.ObservedAt.UTC(),
	})
}
