package decision

import "testing"

func candidateEvidenceGateSummary() ExecutionEnvelopeJEVSelfImprovementEvidenceSummary {
	return ExecutionEnvelopeJEVSelfImprovementEvidenceSummary{
		Status: "bound", CycleStatus: "stable-for-review", CycleEvidenceDigest: "cycle-evidence-digest", LedgerStatus: "stable-for-review", LedgerEvidenceDigest: "ledger-evidence-digest", MetricName: "review-stability", MetricSourceDigest: "metric-source-digest", MetricValueDigest: "metric-value-digest", MetricBindingDigest: "metric-binding-digest", SummaryDigest: "summary-digest", NonExecuting: true, NonAuthorizing: true,
	}
}

func candidateEvidenceGateCandidate() JEVImprovementRevisionCandidate {
	return GenerateJEVImprovementRevisionCandidate(JEVImprovementRevisionCandidateInput{
		Directive:           revisionCandidateDirective(),
		RevisionSource:      "gooo://revision/source/one",
		RevisionChangeDigest: "revision-change-digest",
		NonAuthorizing:      true,
	})
}

func TestBindExecutionEnvelopeJEVSelfImprovementSummaryToCandidate(t *testing.T) {
	candidate := candidateEvidenceGateCandidate()
	got := BindExecutionEnvelopeJEVSelfImprovementSummaryToCandidate(ExecutionEnvelopeJEVSelfImprovementCandidateEvidenceGateInput{
		Summary:        candidateEvidenceGateSummary(),
		Candidate:      candidate,
		NonAuthorizing: true,
	})
	if got.Status != "bound" || got.CandidateStatus != candidate.Status || got.CandidateDigest != candidate.CandidateDigest || got.SummaryDigest != "summary-digest" || got.AdmissionDigest == "" {
		t.Fatalf("got %+v", got)
	}
	if !got.NonExecuting || !got.NonAuthorizing || got.MetricSourceDigest != "metric-source-digest" {
		t.Fatalf("missing boundary %+v", got)
	}
}

func TestBindExecutionEnvelopeJEVSelfImprovementSummaryToCandidatePreservesSummaryStage(t *testing.T) {
	summary := candidateEvidenceGateSummary()
	summary.Status = "UNKNOWN"
	summary.MissingStage = "reverse-observation"
	got := BindExecutionEnvelopeJEVSelfImprovementSummaryToCandidate(ExecutionEnvelopeJEVSelfImprovementCandidateEvidenceGateInput{
		Summary:        summary,
		Candidate:      candidateEvidenceGateCandidate(),
		NonAuthorizing: true,
	})
	if got.Status != "UNKNOWN" || got.MissingStage != "reverse-observation" || got.AdmissionDigest != "" {
		t.Fatalf("got %+v", got)
	}
}

func TestBindExecutionEnvelopeJEVSelfImprovementSummaryToCandidateRejectsTampering(t *testing.T) {
	candidate := candidateEvidenceGateCandidate()
	candidate.EvidenceDigest = "tampered"
	got := BindExecutionEnvelopeJEVSelfImprovementSummaryToCandidate(ExecutionEnvelopeJEVSelfImprovementCandidateEvidenceGateInput{
		Summary:        candidateEvidenceGateSummary(),
		Candidate:      candidate,
		NonAuthorizing: true,
	})
	if got.Status != "UNKNOWN" || got.MissingStage != "revision-candidate" || got.AdmissionDigest != "" {
		t.Fatalf("got %+v", got)
	}
}

func TestBindExecutionEnvelopeJEVSelfImprovementSummaryToCandidateRejectsAuthorization(t *testing.T) {
	got := BindExecutionEnvelopeJEVSelfImprovementSummaryToCandidate(ExecutionEnvelopeJEVSelfImprovementCandidateEvidenceGateInput{NonAuthorizing: false})
	if got.Status != "UNKNOWN" || got.MissingStage != "authorization-boundary" || got.NonAuthorizing || got.AdmissionDigest != "" {
		t.Fatalf("got %+v", got)
	}
}
