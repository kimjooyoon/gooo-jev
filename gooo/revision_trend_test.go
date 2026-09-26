package gooo

import "testing"

func revisionTrendChainForReplacement(t *testing.T, replacement string) (RevisionEvidenceChain, RevisionMetrics) {
	t.Helper()
	assessment := choiceAssessmentForRevision(t)
	candidate, err := ProposeRevision(assessment, RepairRevision)
	if err != nil {
		t.Fatalf("ProposeRevision() error = %v", err)
	}
	edit := SourceEdit{
		Start:       Position{Line: 4, Column: 3},
		End:         Position{Line: 4, Column: 11},
		Replacement: replacement,
	}
	edit.Digest = digestSourceEdit(edit)
	application, err := ApplyRevision(validContract, digestString(validContract), candidate, edit)
	if err != nil {
		t.Fatalf("ApplyRevision() error = %v", err)
	}
	metrics, err := MeasureRevision(validContract, application)
	if err != nil {
		t.Fatalf("MeasureRevision() error = %v", err)
	}
	quality, err := EvaluateRevisionMetrics(metrics)
	if err != nil {
		t.Fatalf("EvaluateRevisionMetrics() error = %v", err)
	}
	feedback, err := DeriveRevisionFeedback(quality)
	if err != nil {
		t.Fatalf("DeriveRevisionFeedback() error = %v", err)
	}
	proposal, err := ProposeRevisionDirection(feedback)
	if err != nil {
		t.Fatalf("ProposeRevisionDirection() error = %v", err)
	}
	materialization, err := MaterializeRevisionCandidate(feedback, proposal, assessment)
	if err != nil {
		t.Fatalf("MaterializeRevisionCandidate() error = %v", err)
	}
	document, err := Parse(validContract)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	generation, err := ObserveCandidateGeneration(document, materialization)
	if err != nil {
		t.Fatalf("ObserveCandidateGeneration() error = %v", err)
	}
	chain, err := ObserveRevisionEvidenceChain(application, metrics, quality, feedback, materialization, generation)
	if err != nil {
		t.Fatalf("ObserveRevisionEvidenceChain() error = %v", err)
	}
	return chain, metrics
}

func TestObserveRevisionTrendBindsTwoEvidenceChains(t *testing.T) {
	previousChain, previousMetrics := revisionTrendChainForReplacement(t, "lineage")
	currentChain, currentMetrics := revisionTrendChainForReplacement(t, "lineage-expanded")
	trend, err := ObserveRevisionTrend(previousChain, previousMetrics, currentChain, currentMetrics)
	if err != nil {
		t.Fatalf("ObserveRevisionTrend() error = %v", err)
	}
	if trend.Status != "BOUND" || !trend.CandidateStable || !trend.SourceStable {
		t.Fatalf("unexpected revision trend: %#v", trend)
	}
	if !trend.IRChangeStable {
		t.Fatalf("IR stability was not retained: %#v", trend)
	}
	if err := trend.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestObserveRevisionTrendRetainsUnknownMetricStage(t *testing.T) {
	previousChain, previousMetrics := revisionTrendChainForReplacement(t, "lineage")
	currentChain, currentMetrics := revisionTrendChainForReplacement(t, "lineage-expanded")
	currentMetrics.MetricsDigest = digestString("tampered")
	trend, err := ObserveRevisionTrend(previousChain, previousMetrics, currentChain, currentMetrics)
	if err == nil {
		t.Fatal("ObserveRevisionTrend() error = nil, want current metrics failure")
	}
	if trend.Status != "UNKNOWN" || trend.MissingStage != "revision-trend-current-metrics" {
		t.Fatalf("unexpected unknown trend: %#v", trend)
	}
}

func TestObserveRevisionTrendIsDeterministic(t *testing.T) {
	previousChain, previousMetrics := revisionTrendChainForReplacement(t, "lineage")
	currentChain, currentMetrics := revisionTrendChainForReplacement(t, "lineage-expanded")
	first, err := ObserveRevisionTrend(previousChain, previousMetrics, currentChain, currentMetrics)
	if err != nil {
		t.Fatalf("first ObserveRevisionTrend() error = %v", err)
	}
	second, err := ObserveRevisionTrend(previousChain, previousMetrics, currentChain, currentMetrics)
	if err != nil {
		t.Fatalf("second ObserveRevisionTrend() error = %v", err)
	}
	if first.TrendDigest != second.TrendDigest {
		t.Fatal("same evidence pair produced different trend digest")
	}
}
