package decision

// ImprovementCycleEvidenceBindingInput links observation admission to one
// self-improvement cycle without authorizing execution.
type ImprovementCycleEvidenceBindingInput struct {
	ObservationAdmissionStatus string
	CandidateDigest            string
	AdmissionDigest            string
	BaselineLedgerDigest       string
	ResultLedgerDigest         string
	CandidateReference         string
	SourceReference            string
	Outcome                    ImprovementCycleOutcome
	NonAuthorizing             bool
}

// ImprovementCycleEvidenceBinding records a cycle and the admission evidence
// that permitted observation of its result.
type ImprovementCycleEvidenceBinding struct {
	Status          string
	Outcome         ImprovementCycleOutcome
	CycleDigest     string
	CandidateDigest string
	AdmissionDigest string
	EvidenceDigest  string
	MissingStage    string
	NonExecuting    bool
	NonAuthorizing  bool
}

// BindImprovementCycleEvidence builds a cycle only after observation
// admission and all ledger references are present.
func BindImprovementCycleEvidence(input ImprovementCycleEvidenceBindingInput) ImprovementCycleEvidenceBinding {
	output := ImprovementCycleEvidenceBinding{
		Status: "UNKNOWN", Outcome: ImprovementUnknown,
		NonExecuting: true, NonAuthorizing: true,
	}
	if !input.NonAuthorizing {
		output.NonAuthorizing = false
		output.MissingStage = "authorization-boundary"
		return output
	}
	if input.ObservationAdmissionStatus == "review" {
		output.Status = "review"
		output.MissingStage = "observation-admission"
		return output
	}
	if input.ObservationAdmissionStatus == "rejected" {
		output.Status = "rejected"
		output.MissingStage = "observation-admission-rejected"
		return output
	}
	if input.ObservationAdmissionStatus != "observation-admitted" {
		output.MissingStage = "observation-admission"
		return output
	}
	if input.CandidateDigest == "" || input.AdmissionDigest == "" {
		output.MissingStage = "observation-evidence"
		return output
	}
	cycle, err := BuildImprovementCycle(
		input.BaselineLedgerDigest,
		input.ResultLedgerDigest,
		input.CandidateReference,
		input.SourceReference,
		input.Outcome,
	)
	if err != nil {
		output.MissingStage = "improvement-cycle"
		return output
	}
	evidence, err := Digest(struct {
		CycleDigest     string
		CandidateDigest string
		AdmissionDigest string
	}{
		CycleDigest:     cycle.CycleDigest,
		CandidateDigest: input.CandidateDigest,
		AdmissionDigest: input.AdmissionDigest,
	})
	if err != nil {
		output.MissingStage = "cycle-evidence"
		return output
	}
	output.Outcome = cycle.Outcome
	output.CycleDigest = cycle.CycleDigest
	output.CandidateDigest = input.CandidateDigest
	output.AdmissionDigest = input.AdmissionDigest
	output.EvidenceDigest = evidence
	switch cycle.Outcome {
	case ImprovementObserved, ImprovementRegressed:
		output.Status = "recorded"
	case ImprovementUnknown, ImprovementReview:
		output.Status = "review"
		output.MissingStage = "improvement-outcome"
	default:
		output.Status = "UNKNOWN"
		output.MissingStage = "improvement-outcome"
	}
	return output
}