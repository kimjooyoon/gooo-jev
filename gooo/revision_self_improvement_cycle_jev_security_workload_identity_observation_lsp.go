package gooo

// JEVWorkloadIdentityObservationLSPStatus is the editor-facing severity for a
// read-only workload identity observation.
type JEVWorkloadIdentityObservationLSPStatus string

const (
	JEVWorkloadIdentityObservationLSPInformation JEVWorkloadIdentityObservationLSPStatus = "Information"
	JEVWorkloadIdentityObservationLSPError       JEVWorkloadIdentityObservationLSPStatus = "Error"
)

// JEVWorkloadIdentityObservationLSPInput mirrors the runtime observation
// without coupling the core to a runtime package.
type JEVWorkloadIdentityObservationLSPInput struct {
	Status               string
	SpiffeID             string
	Audience             string
	Capabilities         []string
	EvidencePrefixDigest string
}

// JEVWorkloadIdentityObservationLSP is a diagnostic-only projection. It never
// authorizes capabilities or executes a workload.
type JEVWorkloadIdentityObservationLSP struct {
	Status               JEVWorkloadIdentityObservationLSPStatus
	Code                 string
	Message              string
	SpiffeID             string
	Audience             string
	Capabilities         []string
	EvidencePrefixDigest string
	IsReadOnly           bool
	CanEdit              bool
	CanExecute           bool
	CanAuthorize         bool
}

// ProjectJEVWorkloadIdentityObservationLSP preserves identity evidence while
// keeping authorization and execution outside the observation boundary.
func ProjectJEVWorkloadIdentityObservationLSP(
	input JEVWorkloadIdentityObservationLSPInput,
) JEVWorkloadIdentityObservationLSP {
	projection := JEVWorkloadIdentityObservationLSP{
		Status:               JEVWorkloadIdentityObservationLSPError,
		Code:                 "jev.security.workload_identity.unknown",
		Message:              "workload identity evidence is missing or incomplete",
		SpiffeID:             input.SpiffeID,
		Audience:             input.Audience,
		Capabilities:         append([]string(nil), input.Capabilities...),
		EvidencePrefixDigest: input.EvidencePrefixDigest,
		IsReadOnly:           true,
		CanEdit:              false,
		CanExecute:           false,
		CanAuthorize:         false,
	}

	if input.Status == "BOUND" {
		projection.Status = JEVWorkloadIdentityObservationLSPInformation
		projection.Code = "jev.security.workload_identity.bound"
		projection.Message = "workload identity evidence is available for inspection"
	}

	return projection
}