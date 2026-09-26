package decision

import "testing"

func validReviewLedgerAggregation() ExecutionEnvelopeJEVImprovementReviewLedgerAggregation {
	aggregation := ExecutionEnvelopeJEVImprovementReviewLedgerAggregation{
		Status: "stable-for-review", Total: 1, Confirmed: 1, InputLedgerDigest: "ledger-input-digest", NonExecuting: true, NonAuthorizing: true,
	}
	aggregation.EvidenceDigest, _ = digestExecutionEnvelopeJEVImprovementReviewLedgerAggregationEvidence(aggregation)
	return aggregation
}

func TestBindExecutionEnvelopeJEVReviewLedgerToMetric(t *testing.T) {
	got := BindExecutionEnvelopeJEVReviewLedgerToMetric(ExecutionEnvelopeJEVReviewLedgerMetricBindingInput{
		Aggregation:        validReviewLedgerAggregation(),
		MetricName:         "review-stability",
		MetricValueDigest:  "metric-value-digest",
		MetricSourceDigest: "metric-source-digest",
		NonAuthorizing:     true,
	})
	if got.Status != "bound" || got.MetricName != "review-stability" || got.BindingDigest == "" || got.AggregationEvidenceDigest == "" {
		t.Fatalf("got %+v", got)
	}
	if !got.NonExecuting || !got.NonAuthorizing || got.InputLedgerDigest != "ledger-input-digest" {
		t.Fatalf("missing boundary %+v", got)
	}
}

func TestBindExecutionEnvelopeJEVReviewLedgerToMetricRequiresSource(t *testing.T) {
	got := BindExecutionEnvelopeJEVReviewLedgerToMetric(ExecutionEnvelopeJEVReviewLedgerMetricBindingInput{
		Aggregation:       validReviewLedgerAggregation(),
		MetricName:        "review-stability",
		MetricValueDigest: "metric-value-digest",
		NonAuthorizing:    true,
	})
	if got.Status != "UNKNOWN" || got.MissingStage != "metric-source" || got.BindingDigest != "" {
		t.Fatalf("got %+v", got)
	}
}

func TestBindExecutionEnvelopeJEVReviewLedgerToMetricPreservesAggregationStage(t *testing.T) {
	aggregation := validReviewLedgerAggregation()
	aggregation.Status = "UNKNOWN"
	aggregation.MissingStage = "ledger[1]-predecessor"
	aggregation.EvidenceDigest = ""
	got := BindExecutionEnvelopeJEVReviewLedgerToMetric(ExecutionEnvelopeJEVReviewLedgerMetricBindingInput{
		Aggregation:        aggregation,
		MetricName:         "review-stability",
		MetricValueDigest:  "metric-value-digest",
		MetricSourceDigest: "metric-source-digest",
		NonAuthorizing:     true,
	})
	if got.Status != "UNKNOWN" || got.MissingStage != "ledger[1]-predecessor" {
		t.Fatalf("got %+v", got)
	}
}

func TestBindExecutionEnvelopeJEVReviewLedgerToMetricRejectsAuthorization(t *testing.T) {
	got := BindExecutionEnvelopeJEVReviewLedgerToMetric(ExecutionEnvelopeJEVReviewLedgerMetricBindingInput{NonAuthorizing: false})
	if got.Status != "UNKNOWN" || got.MissingStage != "authorization-boundary" || got.NonAuthorizing {
		t.Fatalf("got %+v", got)
	}
}
