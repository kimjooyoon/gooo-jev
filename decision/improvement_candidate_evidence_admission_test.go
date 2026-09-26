package decision

import "testing"

func makeImprovementCandidateEvidenceAdmissionInput(t *testing.T, unknownCount uint64) ImprovementCandidateEvidenceAdmissionInput {
	t.Helper()
	cycle, err := BuildImprovementCycle(
		"ledger-before-admission",
		"ledger-after-admission",
		"candidate://change/admission",
		"gooo://jev/source/admission",
		ImprovementObserved,
	)
	if err != nil {
		t.Fatalf("BuildImprovementCycle() error = %v", err)
	}
	replay, err := ReplayImprovementCycle(cycle, "ledger-after-admission", "replay-evidence")
	if err != nil {
		t.Fatalf("ReplayImprovementCycle() error = %v", err)
	}
	status := "measured"
	if unknownCount > 0 {
		status = "measured-with-unknown"
	}
	return ImprovementCandidateEvidenceAdmissionInput{
		CycleBaselineLedgerDigest: cycle.BaselineLedgerDigest,
		CycleResultLedgerDigest: cycle.ResultLedgerDigest,
		CycleCandidateReference: cycle.CandidateReference,
		CycleSourceReference: cycle.SourceReference,
		CycleOutcome: cycle.Outcome,
		CycleNonAuthorizing: cycle.NonAuthorizing,
		CycleDigest: cycle.CycleDigest,
		ReplayRecordedResultLedgerDigest: replay.RecordedResultLedgerDigest,
		ReplayRecomputedResultLedgerDigest: replay.RecomputedResultLedgerDigest,
		ReplayEvidenceDigest: replay.EvidenceDigest,
		ReplayStatus: replay.Status,
		ReplayNonAuthorizing: replay.NonAuthorizing,
		ReplayDigest: replay.ReplayDigest,
		TransitionMetricStatus: status,
		TransitionMetricTotal: 4,
		TransitionMetricStableCount: 2,
		TransitionMetricChangedCount: 2 - unknownCount,
		TransitionMetricUnknownCount: unknownCount,
		TransitionMetricNonAuthorizing: true,
		CandidateReference: "candidate://change/admission",
		SourceReference: "gooo://jev/source/admission",
		NonAuthorizing: true,
	}
}

func TestAdmitImprovementCandidateEvidence(t *testing.T) {
	input := makeImprovementCandidateEvidenceAdmissionInput(t, 0)
	output := AdmitImprovementCandidateEvidence(input)
	if output.Status != "proposed" || output.Decision != CandidateProposed ||
		output.CandidateDigest == "" || output.EvidenceDigest == "" ||
		!output.NonExecuting || !output.NonAuthorizing ||
		output.MissingStage != "" {
		t.Fatalf("unexpected proposed admission: %#v", output)
	}

	input = makeImprovementCandidateEvidenceAdmissionInput(t, 1)
	output = AdmitImprovementCandidateEvidence(input)
	if output.Status != "review" || output.Decision != CandidateReview ||
		output.CandidateDigest == "" ||
		output.MissingStage != "reverse-transition-unknown" {
		t.Fatalf("unexpected review admission: %#v", output)
	}

	input = makeImprovementCandidateEvidenceAdmissionInput(t, 0)
	input.TransitionMetricTotal = 5
	output = AdmitImprovementCandidateEvidence(input)
	if output.Status != "UNKNOWN" || output.CandidateDigest != "" ||
		output.MissingStage != "reverse-transition-metric" {
		t.Fatalf("unexpected metric admission: %#v", output)
	}

	input = makeImprovementCandidateEvidenceAdmissionInput(t, 0)
	input.CycleDigest = "tampered"
	output = AdmitImprovementCandidateEvidence(input)
	if output.Status != "UNKNOWN" || output.MissingStage != "improvement-cycle" {
		t.Fatalf("unexpected cycle admission: %#v", output)
	}

	input = makeImprovementCandidateEvidenceAdmissionInput(t, 0)
	input.NonAuthorizing = false
	output = AdmitImprovementCandidateEvidence(input)
	if output.Status != "UNKNOWN" || output.NonAuthorizing ||
		output.MissingStage != "authorization-boundary" {
		t.Fatalf("unexpected authorization admission: %#v", output)
	}
}