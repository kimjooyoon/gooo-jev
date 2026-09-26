package gooo

import "testing"

func TestCompleteReturnsProvenanceLinkedSymbols(t *testing.T) {
	response := Complete(validContract, "Decision")
	if response.Status != "BOUND" || response.MissingStage != "" {
		t.Fatalf("unexpected completion status: %#v", response)
	}
	if len(response.Items) != 2 {
		t.Fatalf("completion count = %d, want 2: %#v", len(response.Items), response.Items)
	}
	if response.Items[0].Label != "DecisionSpec" || response.Items[1].Label != "DecisionReceipt" {
		t.Fatalf("unexpected completion order: %#v", response.Items)
	}
	if response.IRDigest == "" || response.ItemsDigest == "" ||
		response.Items[0].Kind != EntitySymbol ||
		response.Items[0].Digest == "" ||
		response.Items[0].SymbolDigest != response.Items[0].Digest ||
		response.Items[0].ItemDigest == "" {
		t.Fatalf("missing completion provenance: %#v", response)
	}
	if err := response.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestCompleteIncludesActivitiesByPrefix(t *testing.T) {
	response := Complete(validContract, "Observe")
	if len(response.Items) != 1 || response.Items[0].Kind != ActivitySymbol {
		t.Fatalf("unexpected activity completions: %#v", response.Items)
	}
}

func TestCompleteRetainsUnknownDiagnostics(t *testing.T) {
	response := Complete("package jevdecision
namespace jevdecision
activity broken
", "")
	if response.Status != "UNKNOWN" || response.MissingStage != "syntax" {
		t.Fatalf("unexpected unknown completion: %#v", response)
	}
	if len(response.Items) != 0 || len(response.Diagnostics) != 1 {
		t.Fatalf("unexpected unknown completion contents: %#v", response)
	}
	if err := response.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestValidateRejectsTamperedCompletionDigest(t *testing.T) {
	response := Complete(validContract, "")
	if len(response.Items) == 0 {
		t.Fatal("completion returned no items")
	}
	response.Items[0].ItemDigest = "tampered"
	if err := response.Validate(); err == nil {
		t.Fatal("Validate() error = nil, want tamper rejection")
	}
}

func TestCompleteIsDeterministic(t *testing.T) {
	first := Complete(validContract, "D")
	second := Complete(validContract, "D")
	if len(first.Items) != len(second.Items) || first.ItemsDigest != second.ItemsDigest {
		t.Fatal("same source produced different completion result")
	}
	for index := range first.Items {
		if first.Items[index] != second.Items[index] {
			t.Fatalf("same source produced different completion at %d", index)
		}
	}
}
