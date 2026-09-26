package decision

import "testing"

func TestMeasureExecutionEnvelopeReverseObservationTransitionMetric(t *testing.T) {
	output := MeasureExecutionEnvelopeReverseObservationTransitionMetric(ExecutionEnvelopeReverseObservationTransitionMetricInput{
		StableCount: 4, ChangedCount: 3, UnknownCount: 1, NonAuthorizing: true,
	})
	if output.Status != "measured-with-unknown" || output.Total != 8 ||
		output.StableCount != 4 || output.ChangedCount != 3 ||
		output.UnknownCount != 1 || output.NonAuthorizing != true {
		t.Fatalf("unexpected transition metric: %#v", output)
	}

	output = MeasureExecutionEnvelopeReverseObservationTransitionMetric(ExecutionEnvelopeReverseObservationTransitionMetricInput{
		StableCount: 2, ChangedCount: 1, UnknownCount: 0, NonAuthorizing: true,
	})
	if output.Status != "measured" || output.Total != 3 ||
		output.UnknownCount != 0 {
		t.Fatalf("unexpected measured transition metric: %#v", output)
	}

	output = MeasureExecutionEnvelopeReverseObservationTransitionMetric(ExecutionEnvelopeReverseObservationTransitionMetricInput{
		StableCount: 1, ChangedCount: 1, UnknownCount: 1, NonAuthorizing: false,
	})
	if output.Status != "UNKNOWN" || output.Total != 0 ||
		output.NonAuthorizing != false {
		t.Fatalf("unexpected authorization metric: %#v", output)
	}
}