package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

type RevisionSelfImprovementCycleJEVSecurityWorkloadIdentityBoundaryInput struct {
	WorkloadIdentity                  string
	WorkloadIdentityDigest            string
	ExternalEvidenceDigest            string
	Audience                          string
	ExpiresAtUnix                     int64
	EvidenceStatus                    string
	PermissionState                   string
	GenerationTraceDigest             string
	ReverseObservationCoverageDigest string
}

type RevisionSelfImprovementCycleJEVSecurityWorkloadIdentityBoundaryObservation struct {
	Status                          string
	MissingStage                    string
	WorkloadIdentityDigest          string
	ExternalEvidenceDigest          string
	Audience                        string
	ExpiresAtUnix                   int64
	EvidenceStatus                  string
	PermissionState                 string
	GenerationTraceDigest           string
	ReverseObservationCoverageDigest string
	BoundarySignal                  string
	ObservationDigest               string
	ReadOnly                        bool
	NonExecuting                    bool
	NonAuthorizing                  bool
}

func ObserveRevisionSelfImprovementCycleJEVSecurityWorkloadIdentityBoundary(
	input RevisionSelfImprovementCycleJEVSecurityWorkloadIdentityBoundaryInput,
) RevisionSelfImprovementCycleJEVSecurityWorkloadIdentityBoundaryObservation {
	result := RevisionSelfImprovementCycleJEVSecurityWorkloadIdentityBoundaryObservation{
		Status:                           "UNKNOWN",
		MissingStage:                     "revision-self-improvement-cycle-jev-security-workload-identity-boundary",
		WorkloadIdentityDigest:           input.WorkloadIdentityDigest,
		ExternalEvidenceDigest:            input.ExternalEvidenceDigest,
		Audience:                         input.Audience,
		ExpiresAtUnix:                    input.ExpiresAtUnix,
		EvidenceStatus:                  input.EvidenceStatus,
		PermissionState:                 input.PermissionState,
		GenerationTraceDigest:             input.GenerationTraceDigest,
		ReverseObservationCoverageDigest: input.ReverseObservationCoverageDigest,
		BoundarySignal:                   "jev-security-workload-identity-unknown",
		ReadOnly:                         true,
		NonExecuting:                     true,
		NonAuthorizing:                   true,
	}
	setDigest := func() {
		result.ObservationDigest = digestRevisionSelfImprovementCycleJEVSecurityWorkloadIdentityBoundary(result)
	}
	setDigest()

	if strings.TrimSpace(input.WorkloadIdentity) == "" {
		result.MissingStage = "revision-self-improvement-cycle-jev-security-workload-identity-input"
		setDigest()
		return result
	}
	if !validDigest(input.WorkloadIdentityDigest) ||
		!validDigest(input.ExternalEvidenceDigest) ||
		!validDigest(input.GenerationTraceDigest) ||
		!validDigest(input.ReverseObservationCoverageDigest) {
		result.MissingStage = "revision-self-improvement-cycle-jev-security-workload-identity-lineage"
		setDigest()
		return result
	}
	if strings.TrimSpace(input.Audience) == "" || input.ExpiresAtUnix <= 0 {
		result.MissingStage = "revision-self-improvement-cycle-jev-security-workload-identity-boundary-metadata"
		setDigest()
		return result
	}
	if input.EvidenceStatus != "observed" || input.PermissionState != "DEFER" {
		result.MissingStage = "revision-self-improvement-cycle-jev-security-workload-identity-permission"
		setDigest()
		return result
	}

	result.Status = "BOUND"
	result.MissingStage = ""
	result.BoundarySignal = "jev-security-workload-identity-observed"
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "revision-self-improvement-cycle-jev-security-workload-identity-boundary"
		result.BoundarySignal = "jev-security-workload-identity-unknown"
		setDigest()
	}
	return result
}

func (o RevisionSelfImprovementCycleJEVSecurityWorkloadIdentityBoundaryObservation) Validate() error {
	if o.Status != "BOUND" && o.Status != "UNKNOWN" {
		return fmt.Errorf("security workload identity status is invalid")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound security workload identity has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown security workload identity has no missing stage")
	}
	if o.Status == "UNKNOWN" {
		if o.BoundarySignal != "jev-security-workload-identity-unknown" {
			return fmt.Errorf("unknown security workload identity must preserve incomplete evidence")
		}
	} else {
		if !validDigest(o.WorkloadIdentityDigest) ||
			!validDigest(o.ExternalEvidenceDigest) ||
			!validDigest(o.GenerationTraceDigest) ||
			!validDigest(o.ReverseObservationCoverageDigest) ||
			strings.TrimSpace(o.Audience) == "" ||
			o.ExpiresAtUnix <= 0 ||
			o.EvidenceStatus != "observed" ||
			o.PermissionState != "DEFER" ||
			o.BoundarySignal != "jev-security-workload-identity-observed" {
			return fmt.Errorf("bound security workload identity has invalid evidence boundary")
		}
	}
	if !o.ReadOnly || !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("security workload identity must remain read-only, non-executing, and non-authorizing")
	}
	if o.ObservationDigest != digestRevisionSelfImprovementCycleJEVSecurityWorkloadIdentityBoundary(o) {
		return fmt.Errorf("security workload identity digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementCycleJEVSecurityWorkloadIdentityBoundary(
	o RevisionSelfImprovementCycleJEVSecurityWorkloadIdentityBoundaryObservation,
) string {
	return digestString(strings.Join([]string{
		o.Status,
		o.MissingStage,
		o.WorkloadIdentityDigest,
		o.ExternalEvidenceDigest,
		o.Audience,
		strconv.FormatInt(o.ExpiresAtUnix, 10),
		o.EvidenceStatus,
		o.PermissionState,
		o.GenerationTraceDigest,
		o.ReverseObservationCoverageDigest,
		o.BoundarySignal,
		strconv.FormatBool(o.ReadOnly),
		strconv.FormatBool(o.NonExecuting),
		strconv.FormatBool(o.NonAuthorizing),
	}, "|"))
}