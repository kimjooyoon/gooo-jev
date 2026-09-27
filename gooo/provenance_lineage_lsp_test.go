package gooo

import "testing"

func TestProjectProvenanceLineageLSPProjectsBoundReceipt(t *testing.T) {
	document, assessment, candidate := lineageInputs(t)
	receipt, err := ObserveLineage(document, assessment, candidate)
	if err != nil {
		t.Fatalf("ObserveLineage() error = %v", err)
	}

	projection := ProjectProvenanceLineageLSP(receipt)
	if projection.Status != ProvenanceLineageLSPBound ||
		projection.Code != provenanceLineageLSPComplete ||
		projection.MissingStage != "" ||
		projection.EdgeCount != 4 ||
		projection.ChainDigest != receipt.ChainDigest ||
		projection.EvidenceDigest != receipt.EvidenceDigest ||
		projection.Publishable {
		t.Fatalf("unexpected bound projection: %#v", projection)
	}
	if err := projection.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestProjectProvenanceLineageLSPPreservesUnknownStage(t *testing.T) {
	document, assessment, candidate := lineageInputs(t)
	candidate.DecisionName = "DifferentDecision"
	candidate.CandidateDigest = digestRevisionCandidate(candidate)
	receipt, err := ObserveLineage(document, assessment, candidate)
	if err == nil {
		t.Fatal("ObserveLineage() error = nil, want binding failure")
	}

	projection := ProjectProvenanceLineageLSP(receipt)
	if projection.Status != ProvenanceLineageLSPUnknown ||
		projection.Code != provenanceLineageLSPUnknown ||
		projection.MissingStage != "lineage-binding" ||
		projection.ChainDigest != "" {
		t.Fatalf("unexpected unknown projection: %#v", projection)
	}
	if err := projection.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestProjectProvenanceLineageLSPRejectsTamperedReceipt(t *testing.T) {
	document, assessment, candidate := lineageInputs(t)
	receipt, err := ObserveLineage(document, assessment, candidate)
	if err != nil {
		t.Fatalf("ObserveLineage() error = %v", err)
	}
	receipt.Edges[0].ToDigest = digestString("tampered-edge")

	projection := ProjectProvenanceLineageLSP(receipt)
	if projection.Status != ProvenanceLineageLSPUnknown ||
		projection.Code != provenanceLineageLSPIntegrity ||
		projection.MissingStage != "lineage-validation" {
		t.Fatalf("tampered receipt was projected: %#v", projection)
	}
	if err := projection.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestProjectProvenanceLineageLSPFailsClosedOnCapabilityBoundary(t *testing.T) {
	document, assessment, candidate := lineageInputs(t)
	receipt, err := ObserveLineage(document, assessment, candidate)
	if err != nil {
		t.Fatalf("ObserveLineage() error = %v", err)
	}
	receipt.NonAuthorizing = false

	projection := ProjectProvenanceLineageLSP(receipt)
	if projection.Status != ProvenanceLineageLSPUnknown ||
		projection.Code != provenanceLineageLSPBoundary ||
		projection.NonAuthorizing ||
		projection.MissingStage != "capability-boundary" {
		t.Fatalf("capability boundary was not preserved: %#v", projection)
	}
}

func TestProvenanceLineageLSPRejectsTamperedProjection(t *testing.T) {
	document, assessment, candidate := lineageInputs(t)
	receipt, err := ObserveLineage(document, assessment, candidate)
	if err != nil {
		t.Fatalf("ObserveLineage() error = %v", err)
	}
	projection := ProjectProvenanceLineageLSP(receipt)
	projection.ChainDigest = digestString("tampered-chain")
	if err := projection.Validate(); err == nil {
		t.Fatal("Validate() error = nil, want tamper rejection")
	}
}