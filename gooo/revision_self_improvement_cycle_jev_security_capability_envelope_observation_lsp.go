package gooo

// JEVCapabilityEnvelopeObservationLSPStatus is the editor-facing severity for
// a read-only capability envelope observation.
type JEVCapabilityEnvelopeObservationLSPStatus string

const (
	JEVCapabilityEnvelopeObservationLSPInformation JEVCapabilityEnvelopeObservationLSPStatus = "Information"
	JEVCapabilityEnvelopeObservationLSPWarning     JEVCapabilityEnvelopeObservationLSPStatus = "Warning"
	JEVCapabilityEnvelopeObservationLSPError       JEVCapabilityEnvelopeObservationLSPStatus = "Error"
)

// JEVCapabilityEnvelopeObservationLSPInput mirrors the runtime observation
// without coupling the core to the runtime package.
type JEVCapabilityEnvelopeObservationLSPInput struct {
	Status                   string
	WorkloadIdentityStatus   string
	EvidencePrefixDigest     string
	PlanDigest               string
	NetworkAllowlistDigest   string
	RequestedCapabilities    []string
	ObservedCapabilities     []string
	MissingCapability        string
	TargetStage              string
	ReverseObservationDigest string
	ObservationDigest        string
}

// JEVCapabilityEnvelopeObservationLSP is a diagnostic-only projection. It
// never authorizes capabilities or executes a workload.
type JEVCapabilityEnvelopeObservationLSP struct {
	Status                   JEVCapabilityEnvelopeObservationLSPStatus
	Code                     string
	Message                  string
	WorkloadIdentityStatus   string
	EvidencePrefixDigest     string
	PlanDigest               string
	NetworkAllowlistDigest   string
	RequestedCapabilities    []string
	ObservedCapabilities     []string
	MissingCapability        string
	TargetStage              string
	ReverseObservationDigest string
	ObservationDigest        string
	IsReadOnly               bool
	CanEdit                  bool
	CanExecute               bool
	CanAuthorize             bool
}

// ProjectJEVCapabilityEnvelopeObservationLSP preserves capability evidence
// while keeping authorization and execution outside the observation boundary.
func ProjectJEVCapabilityEnvelopeObservationLSP(
	input JEVCapabilityEnvelopeObservationLSPInput,
) JEVCapabilityEnvelopeObservationLSP {
	projection := JEVCapabilityEnvelopeObservationLSP{
		Status:                   JEVCapabilityEnvelopeObservationLSPError,
		Code:                     "jev.security.capability_envelope.unknown",
		Message:                  "capability envelope evidence is missing or incomplete",
		WorkloadIdentityStatus:   input.WorkloadIdentityStatus,
		EvidencePrefixDigest:     input.EvidencePrefixDigest,
		PlanDigest:               input.PlanDigest,
		NetworkAllowlistDigest:   input.NetworkAllowlistDigest,
		RequestedCapabilities:    append([]string(nil), input.RequestedCapabilities...),
		ObservedCapabilities:     append([]string(nil), input.ObservedCapabilities...),
		MissingCapability:        input.MissingCapability,
		TargetStage:              input.TargetStage,
		ReverseObservationDigest: input.ReverseObservationDigest,
		ObservationDigest:        input.ObservationDigest,
		IsReadOnly:               true,
		CanEdit:                  false,
		CanExecute:               false,
		CanAuthorize:             false,
	}

	switch input.Status {
	case "BOUND":
		if input.WorkloadIdentityStatus != "BOUND" ||
			input.EvidencePrefixDigest == "" ||
			input.PlanDigest == "" ||
			input.NetworkAllowlistDigest == "" ||
			input.MissingCapability != "" ||
			input.ReverseObservationDigest == "" ||
			input.ObservationDigest == "" {
			projection.Message = "capability envelope is marked BOUND but evidence is incomplete"
			return projection
		}
		projection.Status = JEVCapabilityEnvelopeObservationLSPInformation
		projection.Code = "jev.security.capability_envelope.bound"
		projection.Message = "capability envelope evidence is available for inspection"
	case "DEFERRED":
		projection.Status = JEVCapabilityEnvelopeObservationLSPWarning
		projection.Code = "jev.security.capability_envelope.deferred"
		projection.Message = "capability envelope evidence is deferred"
	}

	return projection
}
