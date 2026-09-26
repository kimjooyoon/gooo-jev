package decision

import "testing"

func summaryCycleBinding() ExecutionEnvelopeJEVPlanImprovementCycleBinding {
	return ExecutionEnvelopeJEVPlanImprovementCycleBinding{
		Status: "bound", PlanID: "triage-plan", CycleStatus: "stable-for-review", CycleEvidenceDigest: "cycle-evidence-digest", BindingDigest: "cycle-binding-digest", NonExecuting: true, NonAuthorizing: true,
	}
}

func summaryLedgerAggregation() ExecutionEnvelopeJEVImprovementReviewLedgerAggregation {
	aggregation := ExecutionEnvelopeJEVImprovementReviewLedgerAggregation{
		Status: "stable-for-review", Total: 1, Confirmed: 1, InputLedgerDigest: "ledger-input-digest", NonExecuting: true, NonAuthorizing: true,
	}
	aggregation.EvidenceDigest, _ = digestExecutionEnvelopeJEVImprovementReviewLedgerAggregationEvidence(aggregation)
	return aggregation
}

func summaryMetricBinding() ExecutionEnvelopeJEVReviewLedgerMetricBinding {
	return ExecutionEnvelopeJEVReviewLedgerMetricBinding{
		Status: "bound", MetricName: "review-stability", MetricSourceDigest: "metric-source-digest", MetricValueDigest: "metric-value-digest", AggregationEvidenceDigest: "aggregation-evidence-digest", BindingDigest: "metric-binding-digest", NonExecuting: true, NonAuthorizing: true,
	}
}

func TestBindExecutionEnvelopeJEVSelfImprovementEvidenceSummary(t *testing.T) {
	got := BindExecutionEnvelopeJEVSelfImprovementEvidenceSummary(ExecutionEnvelopeJEVSelfImprovementEvidenceSummaryInput{
		CycleBinding:      summaryCycleBinding(),
		LedgerAggregation: summaryLedgerAggregation(),
		MetricBinding:     summaryMetricBinding(),
		NonAuthorizing:    true,
	})
	if got.Status != "bound" || got.CycleStatus != "stable-for-review" || got.LedgerStatus != "stable-for-review" || got.MetricName != "review-stability" || got.SummaryDigest == "" {
		t.Fatalf("got %+v", got)
	}
	if !got.NonExecuting || !got.NonAuthorizing || got.CycleEvidenceDigest == "" || got.LedgerEvidenceDigest == "" || got.MetricBindingDigest == "" {
		t.Fatalf("missing chain evidence %+v", got)
	}
}

func TestBindExecutionEnvelopeJEVSelfImprovementEvidenceSummaryPreservesCycleStage(t *testing.T) {
	cycle := summaryCycleBinding()
	cycle.Status = "UNKNOWN"
	cycle.MissingStage = "reverse-observation"
	got := BindExecutionEnvelopeJEVSelfImprovementEvidenceSummary(ExecutionEnvelopeJEVSelfImprovementEvidenceSummaryInput{
		CycleBinding:      cycle,
		LedgerAggregation: summaryLedgerAggregation(),
		MetricBinding:     summaryMetricBinding(),
		NonAuthorizing:    true,
	})
	if got.Status != "UNKNOWN" || got.MissingStage != "reverse-observation" || got.SummaryDigest == "" {
		t.Fatalf("got %+v", got)
	}
}

func TestBindExecutionEnvelopeJEVSelfImprovementEvidenceSummaryPreservesLedgerStage(t *testing.T) {
	ledger := summaryLedgerAggregation()
	ledger.Status = "UNKNOWN"
	ledger.MissingStage = "ledger[1]-predecessor"
	ledger.EvidenceDigest = ""
	got := BindExecutionEnvelopeJEVSelfImprovementEvidenceSummary(ExecutionEnvelopeJEVSelfImprovementEvidenceSummaryInput{
		CycleBinding:      summaryCycleBinding(),
		LedgerAggregation: ledger,
		MetricBinding:     summaryMetricBinding(),
		NonAuthorizing:    true,
	})
	if got.Status != "UNKNOWN" || got.MissingStage != "ledger[1]-predecessor" || got.SummaryDigest == "" {
		t.Fatalf("got %+v", got)
	}
}

func TestBindExecutionEnvelopeJEVSelfImprovementEvidenceSummaryPreservesMetricStage(t *testing.T) {
	metric := summaryMetricBinding()
	metric.Status = "UNKNOWN"
	metric.MissingStage = "metric-source"
	got := BindExecutionEnvelopeJEVSelfImprovementEvidenceSummary(ExecutionEnvelopeJEVSelfImprovementEvidenceSummaryInput{
		CycleBinding:      summaryCycleBinding(),
		LedgerAggregation: summaryLedgerAggregation(),
		MetricBinding:     metric,
		NonAuthorizing:    true,
	})
	if got.Status != "UNKNOWN" || got.MissingStage != "metric-source" || got.SummaryDigest == "" {
		t.Fatalf("got %+v", got)
	}
}

func TestBindExecutionEnvelopeJEVSelfImprovementEvidenceSummaryRejectsAuthorization(t *testing.T) {
	got := BindExecutionEnvelopeJEVSelfImprovementEvidenceSummary(ExecutionEnvelopeJEVSelfImprovementEvidenceSummaryInput{NonAuthorizing: false})
	if got.Status != "UNKNOWN" || got.MissingStage != "authorization-boundary" || got.NonAuthorizing || got.SummaryDigest == "" {
		t.Fatalf("got %+v", got)
	}
}
