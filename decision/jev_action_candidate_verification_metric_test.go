package decision

import "testing"

func TestMeasureJEVActionCandidateVerificationMetricPreservesCounts(t *testing.T) {
	output := MeasureJEVActionCandidateVerificationMetric(JEVActionCandidateVerificationMetricInput{
		MetricName:           "jev-action-candidate-verification",
		VerifiedCount:        3,
		ReviewCount:          1,
		EvidenceSourceDigest: "source-digest",
		NonAuthorizing:       true,
	})
	if output.Status != "measured" || output.Total != 4 || output.VerifiedCount != 3 || output.ReviewCount != 1 {
		t.Fatalf("unexpected measured output: %+v", output)
	}
	if output.EvidenceDigest == "" || !output.NonExecuting || !output.NonAuthorizing {
		t.Fatalf("metric lost evidence or safety boundary: %+v", output)
	}
}

func TestMeasureJEVActionCandidateVerificationMetricRetainsUnknownAndCounterexamples(t *testing.T) {
	output := MeasureJEVActionCandidateVerificationMetric(JEVActionCandidateVerificationMetricInput{
		MetricName:           "jev-action-candidate-verification",
		VerifiedCount:        2,
		CounterexampleCount:  1,
		UnknownCount:         1,
		EvidenceSourceDigest: "source-digest",
		NonAuthorizing:       true,
	})
	if output.Status != "measured-with-unknown" || output.Total != 4 || output.CounterexampleCount != 1 || output.UnknownCount != 1 {
		t.Fatalf("unexpected retained disposition output: %+v", output)
	}
}

func TestMeasureJEVActionCandidateVerificationMetricFailsClosed(t *testing.T) {
	output := MeasureJEVActionCandidateVerificationMetric(JEVActionCandidateVerificationMetricInput{
		MetricName:     "jev-action-candidate-verification",
		NonAuthorizing: true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "evidence-source" {
		t.Fatalf("unexpected missing-source output: %+v", output)
	}

	output = MeasureJEVActionCandidateVerificationMetric(JEVActionCandidateVerificationMetricInput{
		MetricName:           "jev-action-candidate-verification",
		VerifiedCount:        1,
		EvidenceSourceDigest: "source-digest",
		NonAuthorizing:       false,
	})
	if output.Status != "UNKNOWN" || output.NonAuthorizing || output.MissingStage != "authorization-boundary" {
		t.Fatalf("unexpected authorization output: %+v", output)
	}

	output = MeasureJEVActionCandidateVerificationMetric(JEVActionCandidateVerificationMetricInput{
		MetricName:           "jev-action-candidate-verification",
		VerifiedCount:        ^uint64(0),
		ReviewCount:          1,
		EvidenceSourceDigest: "source-digest",
		NonAuthorizing:       true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "count-overflow" {
		t.Fatalf("unexpected overflow output: %+v", output)
	}
}
