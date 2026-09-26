package gooo

import "testing"

func revisionEvidenceChainInputs(t *testing.T) (RevisionApplication, RevisionMetrics, RevisionQuality, RevisionFeedback, RevisionCandidateMaterialization, CandidateGenerationObservation) {
	t.Helper()
	assessment := choiceAssessmentForRevision(t)
	candidate, err := ProposeRevision(assessment, RepairRevision)
	if err != nil {
		t.Fatalf("ProposeRevision() error = %v", err)
	}
	edit := SourceEdit{
		Start:       Position{Line: 4, Column: 3},
		End:         Position{Line: 4, Column: 11},
		Replacement: "lineage",
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
	return application, metrics, quality, feedback, materialization, generation
}

func TestObserveRevisionEvidenceChainBindsAllStages(t *testing.T) {
	application, metrics, quality, feedback, materialization, generation := revisionEvidenceChainInputs(t)
	chain, err := ObserveRevisionEvidenceChain(application, metrics, quality, feedback, materialization, generation)
	if err != nil {
		t.Fatalf("ObserveRevisionEvidenceChain() error = %v", err)
	}
	if chain.Status != "BOUND" || !chain.StructureMatch || !chain.ReverseObserved {
		t.Fatalf("unexpected evidence chain: %#v", chain)
	}
	if chain.SourceDigest != application.SourceDigest ||
		chain.CandidateDigest != materialization.CandidateDigest ||
		chain.GenerationObservationDigest != generation.ObservationDigest {
		t.Fatalf("evidence chain provenance mismatch: %#v", chain)
	}
	if err := chain.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestObserveRevisionEvidenceChainRetainsMetricLinkStage(t *testing.T) {
	_, metrics, quality, feedback, materialization, generation := revisionEvidenceChainInputs(t)
	otherCandidate, err := ProposeRevision(choiceAssessmentForRevision(t), ClarifyRevision)
	if err != nil {
		t.Fatalf("ProposeRevision() error = %v", err)
	}
	edit := SourceEdit{Start: Position{Line: 4, Column: 3}, End: Position{Line: 4, Column: 11}, Replacement: "lineage"}
	edit.Digest = digestSourceEdit(edit)
	otherApplication, err := ApplyRevision(validContract, digestString(validContract), otherCandidate, edit)
	if err != nil {
		t.Fatalf("ApplyRevision() error = %v", err)
	}
	chain, chainErr := ObserveRevisionEvidenceChain(otherApplication, metrics, quality, feedback, materialization, generation)
	if chainErr == nil || chain.Status != "UNKNOWN" || chain.MissingStage != "revision-chain-link-metrics" {
		t.Fatalf("unexpected metric link failure: %#v, %v", chain, chainErr)
	}
}

func TestObserveRevisionEvidenceChainRejectsTamperedDigest(t *testing.T) {
	application, metrics, quality, feedback, materialization, generation := revisionEvidenceChainInputs(t)
	chain, err := ObserveRevisionEvidenceChain(application, metrics, quality, feedback, materialization, generation)
	if err != nil {
		t.Fatalf("ObserveRevisionEvidenceChain() error = %v", err)
	}
	chain.ChainDigest = digestString("tampered")
	if err := chain.Validate(); err == nil {
		t.Fatal("Validate() error = nil, want chain digest rejection")
	}
}

func TestObserveRevisionEvidenceChainIsDeterministic(t *testing.T) {
	application, metrics, quality, feedback, materialization, generation := revisionEvidenceChainInputs(t)
	first, err := ObserveRevisionEvidenceChain(application, metrics, quality, feedback, materialization, generation)
	if err != nil {
		t.Fatalf("first ObserveRevisionEvidenceChain() error = %v", err)
	}
	second, err := ObserveRevisionEvidenceChain(application, metrics, quality, feedback, materialization, generation)
	if err != nil {
		t.Fatalf("second ObserveRevisionEvidenceChain() error = %v", err)
	}
	if first.ChainDigest != second.ChainDigest {
		t.Fatal("same evidence produced different chain digest")
	}
}
