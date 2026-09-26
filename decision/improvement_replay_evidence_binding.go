package decision

// ImprovementReplayEvidenceBindingInput links a recorded cycle binding to
// its replay receipt without authorizing execution.
type ImprovementReplayEvidenceBindingInput struct {
	CycleBindingStatus              string
	CycleDigest                     string
	CycleEvidenceDigest             string
	CandidateDigest                 string
	AdmissionDigest                 string
	ReplayCycleDigest               string
	ReplayRecordedResultLedgerDigest string
	ReplayRecomputedResultLedgerDigest string
	ReplayEvidenceDigest            string
	ReplayStatus                    ImprovementReplayStatus
	ReplayDigest                    string
	NonAuthorizing                  bool
}

// ImprovementReplayEvidenceBinding preserves replay direction and the
// provenance that led to the replay.
type ImprovementReplayEvidenceBinding struct {
	Status             string
	ReplayStatus       ImprovementReplayStatus
	CycleDigest        string
	ReplayDigest       string
	CandidateDigest    string
	AdmissionDigest    string
	EvidenceDigest     string
	MissingStage       string
	NonExecuting       bool
	NonAuthorizing     bool
}

// BindImprovementReplayEvidence requires a recorded cycle and valid replay
// before classifying replayed or diverged evidence.
func BindImprovementReplayEvidence(input ImprovementReplayEvidenceBindingInput) ImprovementReplayEvidenceBinding {
	output := ImprovementReplayEvidenceBinding{
		Status: "UNKNOWN", ReplayStatus: ImprovementReplayUnknown,
		NonExecuting: true, NonAuthorizing: true,
	}
	if !input.NonAuthorizing {
		output.NonAuthorizing = false
		output.MissingStage = "authorization-boundary"
		return output
	}
	if input.CycleBindingStatus == "review" {
		output.Status = "review"
		output.MissingStage = "cycle-outcome"
		return output
	}
	if input.CycleBindingStatus != "recorded" {
		output.MissingStage = "cycle-evidence-binding"
		return output
	}
	if input.CycleDigest == "" || input.CycleEvidenceDigest == "" ||
		input.CandidateDigest == "" || input.AdmissionDigest == "" {
		output.MissingStage = "cycle-evidence"
		return output
	}
	replay := ImprovementCycleReplay{
		CycleDigest:                  input.ReplayCycleDigest,
		RecordedResultLedgerDigest:   input.ReplayRecordedResultLedgerDigest,
		RecomputedResultLedgerDigest: input.ReplayRecomputedResultLedgerDigest,
		EvidenceDigest:               input.ReplayEvidenceDigest,
		Status:                       input.ReplayStatus,
		NonAuthorizing:               true,
		ReplayDigest:                 input.ReplayDigest,
	}
	if err := replay.Validate(); err != nil {
		output.MissingStage = "improvement-replay"
		return output
	}
	if replay.CycleDigest != input.CycleDigest {
		output.MissingStage = "cycle-replay-binding"
		return output
	}
	evidence, err := Digest(struct {
		CycleEvidenceDigest string
		ReplayDigest        string
		CandidateDigest     string
		AdmissionDigest     string
	}{
		CycleEvidenceDigest: input.CycleEvidenceDigest,
		ReplayDigest:        replay.ReplayDigest,
		CandidateDigest:     input.CandidateDigest,
		AdmissionDigest:     input.AdmissionDigest,
	})
	if err != nil {
		output.MissingStage = "replay-evidence"
		return output
	}
	output.ReplayStatus = replay.Status
	output.CycleDigest = input.CycleDigest
	output.ReplayDigest = replay.ReplayDigest
	output.CandidateDigest = input.CandidateDigest
	output.AdmissionDigest = input.AdmissionDigest
	output.EvidenceDigest = evidence
	switch replay.Status {
	case ImprovementReplayed:
		output.Status = "replayed"
	case ImprovementDiverged:
		output.Status = "diverged"
	case ImprovementReplayUnknown:
		output.Status = "review"
		output.MissingStage = "replay-unknown"
	default:
		output.Status = "UNKNOWN"
		output.MissingStage = "replay-status"
	}
	return output
}