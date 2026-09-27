package decision

import "testing"

func testDecisionConfidenceChangeOutcomeMetric(
	t *testing.T,
	status DecisionConfidenceChangePlanDispositionObservationStatus,
	evidence string,
) DecisionConfidenceChangeOutcomeMetric {
	t.Helper()
	disposition := DecisionConfidenceChangePlanDisposition{
		ChangePlanDigest:   "plan",
		VerificationDigest: "verification",
		Status:             DecisionConfidenceChangePlanReadyForExternalApplyReview,
		NonExecuting:       true,
	}
	digest, err := Digest(disposition)
	if err != nil {
		t.Fatal(err)
	}
	disposition.DispositionDigest = digest
	observation, err := ObserveDecisionConfidenceChangePlanDisposition(
		disposition,
		status,
		evidence,
	)
	if err != nil {
		t.Fatal(err)
	}
	metric, err := ObserveDecisionConfidenceChangeOutcomeMetric(observation)
	if err != nil {
		t.Fatal(err)
	}
	return metric
}

func TestProjectDecisionConfidenceChangeOutcomeMetricLSPApplied(t *testing.T) {
	metric := testDecisionConfidenceChangeOutcomeMetric(
		t,
		DecisionConfidenceChangePlanObservedApplied,
		"evidence",
	)
	diagnostic := ProjectDecisionConfidenceChangeOutcomeMetricLSP(metric)
	if diagnostic.Status != "BOUND" ||
		diagnostic.Outcome != string(DecisionConfidenceChangePlanObservedApplied) ||
		diagnostic.Code != "JEV_DECISION_OUTCOME_APPLIED" ||
		diagnostic.Severity != "info" {
		t.Fatalf("unexpected applied diagnostic: %#v", diagnostic)
	}
	if err := diagnostic.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestProjectDecisionConfidenceChangeOutcomeMetricLSPPreservesUnknown(t *testing.T) {
	metric := testDecisionConfidenceChangeOutcomeMetric(
		t,
		DecisionConfidenceChangePlanObservedUnknown,
		"",
	)
	diagnostic := ProjectDecisionConfidenceChangeOutcomeMetricLSP(metric)
	if diagnostic.Status != "BOUND" ||
		diagnostic.Outcome != string(DecisionConfidenceChangePlanObservedUnknown) ||
		diagnostic.Code != "JEV_DECISION_OUTCOME_UNKNOWN" ||
		diagnostic.Severity != "warning" {
		t.Fatalf("unexpected unknown diagnostic: %#v", diagnostic)
	}
	if err := diagnostic.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestProjectDecisionConfidenceChangeOutcomeMetricLSPRejectsTampering(t *testing.T) {
	metric := DecisionConfidenceChangeOutcomeMetric{
		ObservationDigest: "observation",
		AppliedCount:      1,
		NonAuthorizing:    true,
		MetricDigest:      "tampered",
	}
	diagnostic := ProjectDecisionConfidenceChangeOutcomeMetricLSP(metric)
	if diagnostic.Status != "UNKNOWN" || diagnostic.MissingStage != "outcome-metric" {
		t.Fatalf("expected unknown tampered diagnostic: %#v", diagnostic)
	}
	if err := diagnostic.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestProjectDecisionConfidenceChangeOutcomeMetricLSPRejectsAuthorizationBoundary(t *testing.T) {
	metric := DecisionConfidenceChangeOutcomeMetric{
		ObservationDigest: "observation",
		AppliedCount:      1,
		NonAuthorizing:    false,
		MetricDigest:      "tampered",
	}
	diagnostic := ProjectDecisionConfidenceChangeOutcomeMetricLSP(metric)
	if diagnostic.Status != "UNKNOWN" ||
		diagnostic.MissingStage != "authorization-boundary" ||
		diagnostic.NonAuthorizing {
		t.Fatalf("expected authorization boundary diagnostic: %#v", diagnostic)
	}
}
