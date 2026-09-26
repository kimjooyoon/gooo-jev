package decision

import "testing"

func validMetricBinding() ExecutionEnvelopeJEVReviewLedgerMetricBinding {
	return ExecutionEnvelopeJEVReviewLedgerMetricBinding{
		Status: "bound", MetricName: "review-stability", MetricSourceDigest: "metric-source-digest", MetricValueDigest: "metric-value-digest", AggregationEvidenceDigest: "aggregation-evidence-digest", BindingDigest: "binding-digest", NonExecuting: true, NonAuthorizing: true,
	}
}

func TestProjectExecutionEnvelopeJEVReviewLedgerMetricBindingLSP(t *testing.T) {
	got := ProjectExecutionEnvelopeJEVReviewLedgerMetricBindingLSP(ExecutionEnvelopeJEVReviewLedgerMetricBindingLSPInput{Binding: validMetricBinding(), NonAuthorizing: true})
	if got.Status != "bound" || got.Code != "JEV_METRIC_PROVENANCE_BOUND" || got.Severity != "info" || got.EvidenceDigest == "" {
		t.Fatalf("got %+v", got)
	}
	if got.MetricName != "review-stability" || got.MetricSourceDigest != "metric-source-digest" || got.MetricValueDigest != "metric-value-digest" || !got.NonExecuting || !got.NonAuthorizing {
		t.Fatalf("missing projection %+v", got)
	}
}

func TestProjectExecutionEnvelopeJEVReviewLedgerMetricBindingLSPPreservesUnknown(t *testing.T) {
	got := ProjectExecutionEnvelopeJEVReviewLedgerMetricBindingLSP(ExecutionEnvelopeJEVReviewLedgerMetricBindingLSPInput{
		Binding:        ExecutionEnvelopeJEVReviewLedgerMetricBinding{Status: "UNKNOWN", MissingStage: "metric-source", NonExecuting: true, NonAuthorizing: true},
		NonAuthorizing: true,
	})
	if got.Status != "UNKNOWN" || got.MissingStage != "metric-source" || got.EvidenceDigest == "" {
		t.Fatalf("got %+v", got)
	}
}

func TestProjectExecutionEnvelopeJEVReviewLedgerMetricBindingLSPRequiresAggregationEvidence(t *testing.T) {
	binding := validMetricBinding()
	binding.AggregationEvidenceDigest = ""
	got := ProjectExecutionEnvelopeJEVReviewLedgerMetricBindingLSP(ExecutionEnvelopeJEVReviewLedgerMetricBindingLSPInput{Binding: binding, NonAuthorizing: true})
	if got.Status != "UNKNOWN" || got.MissingStage != "ledger-aggregation-evidence" || got.EvidenceDigest == "" {
		t.Fatalf("got %+v", got)
	}
}

func TestProjectExecutionEnvelopeJEVReviewLedgerMetricBindingLSPRejectsAuthorization(t *testing.T) {
	got := ProjectExecutionEnvelopeJEVReviewLedgerMetricBindingLSP(ExecutionEnvelopeJEVReviewLedgerMetricBindingLSPInput{NonAuthorizing: false})
	if got.Status != "UNKNOWN" || got.MissingStage != "authorization-boundary" || got.NonAuthorizing || got.EvidenceDigest == "" {
		t.Fatalf("got %+v", got)
	}
}
