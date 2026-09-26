package gooo

import "testing"

func TestObserveRevisionEvidenceHistoryBindsOrderedChains(t *testing.T) {
	firstChain, _ := revisionTrendChainForReplacement(t, "lineage")
	secondChain, _ := revisionTrendChainForReplacement(t, "lineage-expanded")
	thirdChain, _ := revisionTrendChainForReplacement(t, "lineage-expanded-more")
	history, err := ObserveRevisionEvidenceHistory([]RevisionEvidenceChain{firstChain, secondChain, thirdChain})
	if err != nil {
		t.Fatalf("ObserveRevisionEvidenceHistory() error = %v", err)
	}
	if history.Status != "BOUND" || history.ObservationCount != 3 {
		t.Fatalf("unexpected revision evidence history: %#v", history)
	}
	if len(history.ChainDigests) != 3 ||
		history.ChainDigests[0] != firstChain.ChainDigest ||
		history.ChainDigests[1] != secondChain.ChainDigest ||
		history.ChainDigests[2] != thirdChain.ChainDigest {
		t.Fatalf("chain order was not retained: %#v", history)
	}
	if history.FirstChainDigest != firstChain.ChainDigest ||
		history.LastChainDigest != thirdChain.ChainDigest {
		t.Fatalf("chain boundaries were not retained: %#v", history)
	}
	if history.SourceChangeCount != 2 || history.ReverseObservedCount != 3 {
		t.Fatalf("unexpected history counts: %#v", history)
	}
	if err := history.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestObserveRevisionEvidenceHistoryRetainsUnknownIndex(t *testing.T) {
	firstChain, _ := revisionTrendChainForReplacement(t, "lineage")
	secondChain, _ := revisionTrendChainForReplacement(t, "lineage-expanded")
	secondChain.ChainDigest = digestString("tampered")
	history, err := ObserveRevisionEvidenceHistory([]RevisionEvidenceChain{firstChain, secondChain})
	if err == nil {
		t.Fatal("ObserveRevisionEvidenceHistory() error = nil, want tampered chain failure")
	}
	if history.Status != "UNKNOWN" || history.MissingStage != "revision-evidence-history-chain-1" {
		t.Fatalf("unexpected unknown history: %#v", history)
	}
}

func TestObserveRevisionEvidenceHistoryIsDeterministic(t *testing.T) {
	firstChain, _ := revisionTrendChainForReplacement(t, "lineage")
	secondChain, _ := revisionTrendChainForReplacement(t, "lineage-expanded")
	first, err := ObserveRevisionEvidenceHistory([]RevisionEvidenceChain{firstChain, secondChain})
	if err != nil {
		t.Fatalf("first ObserveRevisionEvidenceHistory() error = %v", err)
	}
	second, err := ObserveRevisionEvidenceHistory([]RevisionEvidenceChain{firstChain, secondChain})
	if err != nil {
		t.Fatalf("second ObserveRevisionEvidenceHistory() error = %v", err)
	}
	if first.HistoryDigest != second.HistoryDigest {
		t.Fatal("same chain history produced different history digest")
	}
}
