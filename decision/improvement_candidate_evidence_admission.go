package decision

import "strings"

// ImprovementCandidateEvidenceAdmissionInput flattens the validated
// improvement-cycle, replay, and reverse-transition evidence into one
// declaration-facing input.
type ImprovementCandidateEvidenceAdmissionInput struct {
	CycleBaselineLedgerDigest  string
	CycleResultLedgerDigest    string
	CycleCandidateReference    string
	CycleSourceReference       string
	CycleOutcome               ImprovementCycleOutcome
	CycleNonAuthorizing        bool
	CycleDigest                string
	ReplayRecordedResultLedgerDigest   string
	ReplayRecomputedResultLedgerDigest string
	ReplayEvidenceDigest       string
	ReplayStatus               ImprovementReplayStatus
	ReplayNonAuthorizing       bool
	ReplayDigest               string
	TransitionMetricStatus     string
	TransitionMetricTotal      uint64
	TransitionMetricStableCount uint64
	TransitionMetricChangedCount uint64
	TransitionMetricUnknownCount uint64
	TransitionMetricNonAuthorizing bool
	CandidateReference         string
	SourceReference            string
	NonAuthorizing             bool
}

// ImprovementCandidateEvidenceAdmission records whether a candidate may be
// proposed for review without ever becoming executable.
type ImprovementCandidateEvidenceAdmission struct {
	Status             string
	Decision           ImprovementCandidateDecision
	CycleDigest        string
	ReplayDigest       string
	EvidenceDigest     string
	CandidateDigest    string
	CandidateReference string
	SourceReference    string
	MissingStage       string
	NonExecuting       bool
	NonAuthorizing     bool
}

// AdmitImprovementCandidateEvidence requires a valid replay and complete
// reverse-transition metric before proposing a non-executing candidate.
func AdmitImprovementCandidateEvidence(input ImprovementCandidateEvidenceAdmissionInput) ImprovementCandidateEvidenceAdmission {
	output := ImprovementCandidateEvidenceAdmission{
		Status: "UNKNOWN", Decision: CandidateReview,
		NonExecuting: true, NonAuthorizing: true,
	}
	if !input.NonAuthorizing {
		output.NonAuthorizing = false
		output.MissingStage = "authorization-boundary"
		return output
	}
	cycle := ImprovementCycle{
		BaselineLedgerDigest: input.CycleBaselineLedgerDigest,
		ResultLedgerDigest:   input.CycleResultLedgerDigest,
		CandidateReference:   input.CycleCandidateReference,
		SourceReference:      input.CycleSourceReference,
		Outcome:              input.CycleOutcome,
		NonAuthorizing:       input.CycleNonAuthorizing,
		CycleDigest:          input.CycleDigest,
	}
	if err := cycle.Validate(); err != nil {
		output.MissingStage = "improvement-cycle"
		return output
	}
	replay := ImprovementCycleReplay{
		CycleDigest:                  input.CycleDigest,
		RecordedResultLedgerDigest:   input.ReplayRecordedResultLedgerDigest,
		RecomputedResultLedgerDigest: input.ReplayRecomputedResultLedgerDigest,
		EvidenceDigest:               input.ReplayEvidenceDigest,
		Status:                       input.ReplayStatus,
		NonAuthorizing:               input.ReplayNonAuthorizing,
		ReplayDigest:                 input.ReplayDigest,
	}
	if err := replay.Validate(); err != nil {
		output.MissingStage = "improvement-replay"
		return output
	}
	if replay.CycleDigest != cycle.CycleDigest {
		output.MissingStage = "cycle-replay-binding"
		return output
	}
	if !input.TransitionMetricNonAuthorizing {
		output.NonAuthorizing = false
		output.MissingStage = "authorization-boundary"
		return output
	}
	metricTotal := input.TransitionMetricStableCount +
		input.TransitionMetricChangedCount +
		input.TransitionMetricUnknownCount
	if input.TransitionMetricTotal == 0 || metricTotal != input.TransitionMetricTotal {
		output.MissingStage = "reverse-transition-metric"
		return output
	}
	if strings.TrimSpace(input.CandidateReference) == "" {
		output.MissingStage = "candidate-reference"
		return output
	}
	if strings.TrimSpace(input.SourceReference) == "" {
		output.MissingStage = "source-reference"
		return output
	}
	metric := ExecutionEnvelopeReverseObservationTransitionMetric{
		Status:       input.TransitionMetricStatus,
		Total:        input.TransitionMetricTotal,
		StableCount:  input.TransitionMetricStableCount,
		ChangedCount: input.TransitionMetricChangedCount,
		UnknownCount: input.TransitionMetricUnknownCount,
		NonAuthorizing: true,
	}
	if metric.Status == "measured" && metric.UnknownCount != 0 {
		output.MissingStage = "reverse-transition-metric-status"
		return output
	}
	if metric.Status == "measured-with-unknown" && metric.UnknownCount == 0 {
		output.MissingStage = "reverse-transition-metric-status"
		return output
	}
	if metric.Status != "measured" && metric.Status != "measured-with-unknown" {
		output.MissingStage = "reverse-transition-metric-status"
		return output
	}
	evidenceDigest, err := Digest(metric)
	if err != nil {
		output.MissingStage = "reverse-transition-metric-digest"
		return output
	}
	output.CycleDigest = cycle.CycleDigest
	output.ReplayDigest = replay.ReplayDigest
	output.EvidenceDigest = evidenceDigest
	output.CandidateReference = input.CandidateReference
	output.SourceReference = input.SourceReference
	if metric.Status == "measured-with-unknown" {
		output.Status = "review"
		output.MissingStage = "reverse-transition-unknown"
	} else {
		output.Status = "proposed"
		output.Decision = CandidateProposed
	}
	candidateDigest, err := Digest(output)
	if err != nil {
		output.Status = "UNKNOWN"
		output.Decision = CandidateReview
		output.CandidateDigest = ""
		output.MissingStage = "candidate-digest"
		return output
	}
	output.CandidateDigest = candidateDigest
	return output
}