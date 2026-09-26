package gooo

import "testing"

func TestObserveRevisionCandidateSelectionBindsSingleCandidate(t *testing.T) {
	firstChain, _ := revisionTrendChainForReplacement(t, "lineage")
	secondChain, _ := revisionTrendChainForReplacement(t, "lineage-expanded")
	observation, err := ObserveRevisionCandidateSelection([]RevisionEvidenceChain{firstChain, secondChain})
	if err != nil {
		t.Fatalf("ObserveRevisionCandidateSelection() error = %v", err)
	}
	if observation.Status != "BOUND" || observation.SelectionStatus != "single" {
		t.Fatalf("unexpected candidate selection: %#v", observation)
	}
	if observation.SelectedCandidateDigest != firstChain.CandidateDigest ||
		len(observation.UniqueCandidateDigests) != 1 {
		t.Fatalf("single candidate was not retained: %#v", observation)
	}
	if err := observation.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestObserveRevisionCandidateSelectionDoesNotChooseAmbiguousCandidate(t *testing.T) {
	firstChain, _ := revisionTrendChainForReplacement(t, "lineage")
	secondChain, _ := revisionTrendChainForReplacement(t, "lineage-expanded")
	secondChain.CandidateDigest = digestString("alternative-candidate")
	secondChain.ChainDigest = digestRevisionEvidenceChain(secondChain)
	observation, err := ObserveRevisionCandidateSelection([]RevisionEvidenceChain{firstChain, secondChain})
	if err != nil {
		t.Fatalf("ObserveRevisionCandidateSelection() error = %v", err)
	}
	if observation.SelectionStatus != "ambiguous" || observation.SelectedCandidateDigest != "" {
		t.Fatalf("ambiguous candidate was selected: %#v", observation)
	}
	if err := observation.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestObserveRevisionCandidateSelectionRetainsUnknownStage(t *testing.T) {
	firstChain, _ := revisionTrendChainForReplacement(t, "lineage")
	secondChain, _ := revisionTrendChainForReplacement(t, "lineage-expanded")
	secondChain.ChainDigest = digestString("tampered")
	observation, err := ObserveRevisionCandidateSelection([]RevisionEvidenceChain{firstChain, secondChain})
	if err == nil {
		t.Fatal("ObserveRevisionCandidateSelection() error = nil, want tampered chain failure")
	}
	if observation.Status != "UNKNOWN" || observation.MissingStage != "revision-candidate-selection-chain-1" {
		t.Fatalf("unexpected unknown candidate selection: %#v", observation)
	}
}
