package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

type RevisionSelfImprovementCycleSecurityWorkloadIdentityBoundaryInput struct {
	WorkloadID string
	IssuerID   string
	AudienceID string
	PolicyDigest string
}

type RevisionSelfImprovementCycleSecurityWorkloadIdentityBoundaryObservation struct {
	Status          string
	MissingStage    string
	WorkloadID      string
	IssuerID        string
	AudienceID      string
	PolicyDigest    string
	BoundarySignal  string
	BoundaryDigest  string
	NonExecuting    bool
	NonAuthorizing  bool
}

func ObserveRevisionSelfImprovementCycleSecurityWorkloadIdentityBoundary(
	input RevisionSelfImprovementCycleSecurityWorkloadIdentityBoundaryInput,
) (RevisionSelfImprovementCycleSecurityWorkloadIdentityBoundaryObservation, error) {
	result := RevisionSelfImprovementCycleSecurityWorkloadIdentityBoundaryObservation{
		Status:         "UNKNOWN",
		MissingStage:   "revision-self-improvement-cycle-security-workload-identity-boundary",
		WorkloadID:     input.WorkloadID,
		IssuerID:       input.IssuerID,
		AudienceID:     input.AudienceID,
		PolicyDigest:   input.PolicyDigest,
		BoundarySignal: "workload-identity-boundary-unknown",
		NonExecuting:   true,
		NonAuthorizing: true,
	}
	setDigest := func() {
		result.BoundaryDigest = digestRevisionSelfImprovementCycleSecurityWorkloadIdentityBoundary(result)
	}
	setDigest()

	if err := input.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-cycle-security-workload-identity-boundary-input"
		setDigest()
		return result, fmt.Errorf("security workload identity boundary input is not valid: %w", err)
	}

	result.Status = "BOUND"
	result.MissingStage = ""
	result.BoundarySignal = "workload-identity-boundary-observed"
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "revision-self-improvement-cycle-security-workload-identity-boundary"
		result.BoundarySignal = "workload-identity-boundary-unknown"
		setDigest()
		return result, fmt.Errorf("security workload identity boundary observation is not valid: %w", err)
	}
	return result, nil
}

func (i RevisionSelfImprovementCycleSecurityWorkloadIdentityBoundaryInput) Validate() error {
	for name, id := range map[string]string{
		"workload": i.WorkloadID,
		"issuer":   i.IssuerID,
		"audience": i.AudienceID,
	} {
		if !validJEVSPIFFEWorkloadIdentity(id) {
			return fmt.Errorf("security workload identity %s is not a valid SPIFFE ID", name)
		}
	}
	if !validDigest(i.PolicyDigest) {
		return fmt.Errorf("security workload identity policy digest is invalid")
	}
	return nil
}

func (o RevisionSelfImprovementCycleSecurityWorkloadIdentityBoundaryObservation) Validate() error {
	if o.Status != "BOUND" && o.Status != "UNKNOWN" {
		return fmt.Errorf("security workload identity boundary status is invalid")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound security workload identity boundary has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown security workload identity boundary has no missing stage")
	}
	for name, id := range map[string]string{
		"workload": o.WorkloadID,
		"issuer":   o.IssuerID,
		"audience": o.AudienceID,
	} {
		if !validJEVSPIFFEWorkloadIdentity(id) {
			return fmt.Errorf("security workload identity %s is not a valid SPIFFE ID", name)
		}
	}
	if !validDigest(o.PolicyDigest) || !validDigest(o.BoundaryDigest) {
		return fmt.Errorf("security workload identity boundary digest is invalid")
	}
	if o.BoundarySignal != "workload-identity-boundary-observed" &&
		o.BoundarySignal != "workload-identity-boundary-unknown" {
		return fmt.Errorf("security workload identity boundary signal is invalid")
	}
	if o.Status == "BOUND" && o.BoundarySignal != "workload-identity-boundary-observed" {
		return fmt.Errorf("bound security workload identity boundary is not observed")
	}
	if !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("security workload identity boundary must remain non-executing and non-authorizing")
	}
	if o.BoundaryDigest != digestRevisionSelfImprovementCycleSecurityWorkloadIdentityBoundary(o) {
		return fmt.Errorf("security workload identity boundary digest does not match its fields")
	}
	return nil
}

func validJEVSPIFFEWorkloadIdentity(id string) bool {
	if !strings.HasPrefix(id, "spiffe://") {
		return false
	}
	authorityAndPath := strings.TrimPrefix(id, "spiffe://")
	parts := strings.SplitN(authorityAndPath, "/", 2)
	return len(parts) == 2 && parts[0] != "" && parts[1] != ""
}

func digestRevisionSelfImprovementCycleSecurityWorkloadIdentityBoundary(
	observation RevisionSelfImprovementCycleSecurityWorkloadIdentityBoundaryObservation,
) string {
	return digestString(strings.Join([]string{
		observation.Status,
		observation.MissingStage,
		observation.WorkloadID,
		observation.IssuerID,
		observation.AudienceID,
		observation.PolicyDigest,
		observation.BoundarySignal,
		strconv.FormatBool(observation.NonExecuting),
		strconv.FormatBool(observation.NonAuthorizing),
	}, "|"))
}
