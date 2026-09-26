package decision

import "testing"

func TestMeasureExecutionEnvelopeAuditDispositionMetric(t *testing.T) {
	output := MeasureExecutionEnvelopeAuditDispositionMetric(ExecutionEnvelopeAuditDispositionMetricInput{
		AuditOnlyCount: 4, ReviewOnlyCount: 2, BlockedCount: 1, NonAuthorizing: true,
	})
	if output.Status != "measured" || output.Total != 7 ||
		output.DominantMode != "audit-only" || output.NonAuthorizing != true {
		t.Fatalf("unexpected audit metric: %#v", output)
	}

	output = MeasureExecutionEnvelopeAuditDispositionMetric(ExecutionEnvelopeAuditDispositionMetricInput{
		AuditOnlyCount: 2, ReviewOnlyCount: 2, BlockedCount: 1, NonAuthorizing: true,
	})
	if output.Status != "measured" || output.DominantMode != "inconclusive" {
		t.Fatalf("unexpected tied metric: %#v", output)
	}

	output = MeasureExecutionEnvelopeAuditDispositionMetric(ExecutionEnvelopeAuditDispositionMetricInput{
		AuditOnlyCount: -1, ReviewOnlyCount: 2, BlockedCount: 1, NonAuthorizing: true,
	})
	if output.Status != "UNKNOWN" || output.Total != 0 || output.DominantMode != "" {
		t.Fatalf("unexpected invalid metric: %#v", output)
	}

	output = MeasureExecutionEnvelopeAuditDispositionMetric(ExecutionEnvelopeAuditDispositionMetricInput{
		AuditOnlyCount: 1, ReviewOnlyCount: 1, BlockedCount: 1, NonAuthorizing: false,
	})
	if output.Status != "UNKNOWN" || output.NonAuthorizing != false {
		t.Fatalf("unexpected authorization metric: %#v", output)
	}
}