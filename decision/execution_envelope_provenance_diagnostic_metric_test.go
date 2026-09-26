package decision

import "testing"

func TestMeasureExecutionEnvelopeProvenanceDiagnosticMetric(t *testing.T) {
	output := MeasureExecutionEnvelopeProvenanceDiagnosticMetric(ExecutionEnvelopeProvenanceDiagnosticMetricInput{
		ClearCount: 3, DiagnosticCount: 2, UnknownCount: 1, NonAuthorizing: true,
	})
	if output.Status != "measured-with-unknown" || output.Total != 6 ||
		output.ClearCount != 3 || output.DiagnosticCount != 2 ||
		output.UnknownCount != 1 || output.NonAuthorizing != true {
		t.Fatalf("unexpected diagnostic metric: %#v", output)
	}

	output = MeasureExecutionEnvelopeProvenanceDiagnosticMetric(ExecutionEnvelopeProvenanceDiagnosticMetricInput{
		ClearCount: 3, DiagnosticCount: 2, UnknownCount: 0, NonAuthorizing: true,
	})
	if output.Status != "measured" || output.Total != 5 ||
		output.UnknownCount != 0 {
		t.Fatalf("unexpected clear metric: %#v", output)
	}

	output = MeasureExecutionEnvelopeProvenanceDiagnosticMetric(ExecutionEnvelopeProvenanceDiagnosticMetricInput{
		ClearCount: -1, DiagnosticCount: 2, UnknownCount: 0, NonAuthorizing: true,
	})
	if output.Status != "UNKNOWN" || output.Total != 0 {
		t.Fatalf("unexpected invalid metric: %#v", output)
	}

	output = MeasureExecutionEnvelopeProvenanceDiagnosticMetric(ExecutionEnvelopeProvenanceDiagnosticMetricInput{
		ClearCount: 1, DiagnosticCount: 1, UnknownCount: 1, NonAuthorizing: false,
	})
	if output.Status != "UNKNOWN" || output.NonAuthorizing != false {
		t.Fatalf("unexpected authorization metric: %#v", output)
	}
}