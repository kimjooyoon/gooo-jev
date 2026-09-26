package gooo

import "testing"

func selectedCandidateInputs(t *testing.T) (RevisionCandidateSelectionObservation, RevisionCandidateMaterialization) {
	t.Helper()
	firstChain, _ := revisionTrendChainForReplacement(t, "lineage")
	secondChain, _ := revisionTrendChainForReplacement(t, "lineage-expanded")
	selection, err := ObserveRevisionCandidateSelection([]RevisionEvidenceChain{firstChain, secondChain})
	if err != nil {
		t.Fatalf("ObserveRevisionCandidateSelection() error = %v", err)
	}
	materialization := candidateMaterializationForGeneration(t)
	if selection.SelectedCandidateDigest != materialization.CandidateDigest {
		t.Fatalf("selection and materialization candidate digests differ: %#v %#v", selection, materialization)
	}
	return selection, materialization
}

func TestBindRevisionCandidateSelectionBindsMaterialization(t *testing.T) {
	selection, materialization := selectedCandidateInputs(t)
	binding, err := BindRevisionCandidateSelection(selection, materialization)
	if err != nil {
		t.Fatalf("BindRevisionCandidateSelection() error = %v", err)
	}
	if binding.Status != "BOUND" || binding.SelectionStatus != "single" {
		t.Fatalf("unexpected candidate binding: %#v", binding)
	}
	if binding.CandidateDigest != materialization.CandidateDigest ||
		binding.SourceDigest != materialization.SourceDigest {
		t.Fatalf("candidate binding lost provenance: %#v", binding)
	}
	if err := binding.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestBindRevisionCandidateSelectionStopsOnAmbiguousSelection(t *testing.T) {
	firstChain, _ := revisionTrendChainForReplacement(t, "lineage")
	secondChain, _ := revisionTrendChainForReplacement(t, "lineage-expanded")
	secondChain.CandidateDigest = digestString("alternative-candidate")
	secondChain.ChainDigest = digestRevisionEvidenceChain(secondChain)
	selection, err := ObserveRevisionCandidateSelection([]RevisionEvidenceChain{firstChain, secondChain})
	if err != nil {
		t.Fatalf("ObserveRevisionCandidateSelection() error = %v", err)
	}
	materialization := candidateMaterializationForGeneration(t)
	binding, err := BindRevisionCandidateSelection(selection, materialization)
	if err == nil {
		t.Fatal("BindRevisionCandidateSelection() error = nil, want ambiguous selection failure")
	}
	if binding.Status != "UNKNOWN" || binding.MissingStage != "revision-candidate-binding-selection" {
		t.Fatalf("unexpected unknown binding: %#v", binding)
	}
}

func TestBindRevisionCandidateSelectionRetainsMaterializationFailure(t *testing.T) {
	selection, materialization := selectedCandidateInputs(t)
	materialization.MaterializationDigest = digestString("tampered")
	binding, err := BindRevisionCandidateSelection(selection, materialization)
	if err == nil {
		t.Fatal("BindRevisionCandidateSelection() error = nil, want materialization failure")
	}
	if binding.Status != "UNKNOWN" || binding.MissingStage != "revision-candidate-binding-materialization" {
		t.Fatalf("unexpected unknown binding: %#v", binding)
	}
}

func TestBindRevisionCandidateSelectionIsDeterministic(t *testing.T) {
	selection, materialization := selectedCandidateInputs(t)
	first, err := BindRevisionCandidateSelection(selection, materialization)
	if err != nil {
		t.Fatalf("first BindRevisionCandidateSelection() error = %v", err)
	}
	second, err := BindRevisionCandidateSelection(selection, materialization)
	if err != nil {
		t.Fatalf("second BindRevisionCandidateSelection() error = %v", err)
	}
	if first.BindingDigest != second.BindingDigest {
		t.Fatal("same selection and materialization produced different binding digest")
	}
}
