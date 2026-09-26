package gooo

import "testing"

func TestObserveRevisionTrendWindowBindsOrderedHistory(t *testing.T) {
	previousChain, previousMetrics := revisionTrendChainForReplacement(t, "lineage")
	middleChain, middleMetrics := revisionTrendChainForReplacement(t, "lineage-expanded")
	currentChain, currentMetrics := revisionTrendChainForReplacement(t, "lineage-expanded-more")
	firstTrend, err := ObserveRevisionTrend(previousChain, previousMetrics, middleChain, middleMetrics)
	if err != nil {
		t.Fatalf("first ObserveRevisionTrend() error = %v", err)
	}
	secondTrend, err := ObserveRevisionTrend(middleChain, middleMetrics, currentChain, currentMetrics)
	if err != nil {
		t.Fatalf("second ObserveRevisionTrend() error = %v", err)
	}
	window, err := ObserveRevisionTrendWindow([]RevisionTrendObservation{firstTrend, secondTrend})
	if err != nil {
		t.Fatalf("ObserveRevisionTrendWindow() error = %v", err)
	}
	if window.Status != "BOUND" || window.ObservationCount != 2 {
		t.Fatalf("unexpected revision trend window: %#v", window)
	}
	if len(window.TrendDigests) != 2 ||
		window.TrendDigests[0] != firstTrend.TrendDigest ||
		window.TrendDigests[1] != secondTrend.TrendDigest {
		t.Fatalf("trend order was not retained: %#v", window)
	}
	if window.FirstTrendDigest != firstTrend.TrendDigest ||
		window.LastTrendDigest != secondTrend.TrendDigest {
		t.Fatalf("trend boundaries were not retained: %#v", window)
	}
	if window.NarrowerCount+window.WiderCount+window.StableCount+window.MixedCount != 2 {
		t.Fatalf("trend signal counts do not cover history: %#v", window)
	}
	if window.CandidateStableCount > window.ObservationCount ||
		window.IRChangeStableCount > window.ObservationCount ||
		window.SourceStableCount > window.ObservationCount {
		t.Fatalf("invalid stability counts: %#v", window)
	}
	if err := window.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestObserveRevisionTrendWindowRetainsUnknownIndex(t *testing.T) {
	previousChain, previousMetrics := revisionTrendChainForReplacement(t, "lineage")
	currentChain, currentMetrics := revisionTrendChainForReplacement(t, "lineage-expanded")
	trend, err := ObserveRevisionTrend(previousChain, previousMetrics, currentChain, currentMetrics)
	if err != nil {
		t.Fatalf("ObserveRevisionTrend() error = %v", err)
	}
	tampered := trend
	tampered.TrendDigest = digestString("tampered")
	window, err := ObserveRevisionTrendWindow([]RevisionTrendObservation{trend, tampered})
	if err == nil {
		t.Fatal("ObserveRevisionTrendWindow() error = nil, want tampered observation failure")
	}
	if window.Status != "UNKNOWN" || window.MissingStage != "revision-trend-window-observation-1" {
		t.Fatalf("unexpected unknown window: %#v", window)
	}
}

func TestObserveRevisionTrendWindowIsDeterministic(t *testing.T) {
	previousChain, previousMetrics := revisionTrendChainForReplacement(t, "lineage")
	currentChain, currentMetrics := revisionTrendChainForReplacement(t, "lineage-expanded")
	trend, err := ObserveRevisionTrend(previousChain, previousMetrics, currentChain, currentMetrics)
	if err != nil {
		t.Fatalf("ObserveRevisionTrend() error = %v", err)
	}
	first, err := ObserveRevisionTrendWindow([]RevisionTrendObservation{trend})
	if err != nil {
		t.Fatalf("first ObserveRevisionTrendWindow() error = %v", err)
	}
	second, err := ObserveRevisionTrendWindow([]RevisionTrendObservation{trend})
	if err != nil {
		t.Fatalf("second ObserveRevisionTrendWindow() error = %v", err)
	}
	if first.WindowDigest != second.WindowDigest {
		t.Fatal("same trend history produced different window digest")
	}
}
