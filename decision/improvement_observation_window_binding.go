package decision

// ImprovementObservationWindowBindingInput binds an accepted review gate to
// a validated lifecycle attestation window.
type ImprovementObservationWindowBindingInput struct {
	ReviewGateStatus string
	CandidateDigest string
	ReviewDigest    string
	Window          ExecutionLifecycleAttestationWindow
	ExecutionGranted bool
	NonAuthorizing  bool
}

// ImprovementObservationWindowBinding admits observation only; it never
// creates an execution grant.
type ImprovementObservationWindowBinding struct {
	Status            string
	ObservationOnly   bool
	CandidateDigest   string
	ReviewDigest      string
	WindowDigest      string
	ObservationDigest string
	AdmissionDigest   string
	MissingStage      string
	ExecutionGranted  bool
	NonAuthorizing    bool
}

// BindImprovementObservationWindow requires an accepted review and a valid
// attestation window before admitting observation.
func BindImprovementObservationWindow(input ImprovementObservationWindowBindingInput) ImprovementObservationWindowBinding {
	output := ImprovementObservationWindowBinding{
		Status: "UNKNOWN", ExecutionGranted: false, NonAuthorizing: true,
	}
	if !input.NonAuthorizing || input.ExecutionGranted {
		output.NonAuthorizing = false
		output.MissingStage = "authorization-boundary"
		return output
	}
	if input.ReviewGateStatus == "rejected" {
		output.Status = "rejected"
		output.CandidateDigest = input.CandidateDigest
		output.ReviewDigest = input.ReviewDigest
		return output
	}
	if input.ReviewGateStatus == "review" {
		output.Status = "review"
		output.MissingStage = "review-gate"
		return output
	}
	if input.ReviewGateStatus != "accepted-for-observation" {
		output.MissingStage = "review-gate"
		return output
	}
	if input.CandidateDigest == "" || input.ReviewDigest == "" {
		output.MissingStage = "review-evidence"
		return output
	}
	if err := input.Window.Validate(); err != nil {
		output.MissingStage = "attestation-window"
		return output
	}
	if input.Window.Status != ExecutionLifecycleAttestationWithinWindow {
		output.Status = "hold"
		output.MissingStage = input.Window.MissingStage
		if output.MissingStage == "" {
			output.MissingStage = "attestation-window"
		} else {
			output.MissingStage = "attestation-window:" + output.MissingStage
		}
		return output
	}
	output.Status = "observation-admitted"
	output.ObservationOnly = true
	output.CandidateDigest = input.CandidateDigest
	output.ReviewDigest = input.ReviewDigest
	output.WindowDigest = input.Window.WindowDigest
	output.ObservationDigest = input.Window.ObservationDigest
	admissionDigest, err := Digest(output)
	if err != nil {
		output.Status = "UNKNOWN"
		output.ObservationOnly = false
		output.MissingStage = "observation-admission-digest"
		return output
	}
	output.AdmissionDigest = admissionDigest
	return output
}