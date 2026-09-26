package gooo

import "testing"

func evidenceSummaryInputs(t *testing.T) (RevisionEvidenceChain, RevisionEvidenceGenerationBinding, RevisionMetrics) {
	t.Helper()
	chain, metrics := revisionTrendChainForReplacement(t, "lineage")
	document, err := Parse(validContract)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	generation, err := Generate(document)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	binding, err := BindRevisionEvidenceGeneration(chain, generation)
	if err != nil {
		t.Fatalf("BindRevisionEvidenceGeneration() error = %v", err)
	}
	return chain, binding, metrics
}

func TestSummarizeRevisionEvidencePreservesAllStages(t *testing.T) {
	chain, binding, metrics := evidenceSummaryInputs(t)
	summary, err := SummarizeRevisionEvidence(chain, binding, metrics)
	if err != nil {
		t.Fatalf("SummarizeRevisionEvidence() error = %v", err)
	}
	if summary.Status != "BOUND" || summary.EvidenceStageCount != 6 || !summary.StructureMatch || !summary.ReverseObserved {
		t.Fatalf("unexpected revision evidence summary: %#v", summary)
	}
	if summary.ChainDigest != chain.ChainDigest ||
		summary.GenerationBindingDigest != binding.BindingDigest ||
		summary.MetricsDigest != metrics.MetricsDigest {
		t.Fatalf("summary lost stage links: %#v", summary)
	}
	if err := summary.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestSummarizeRevisionEvidenceRetainsChainFailure(t *testing.T) {
	chain, binding, metrics := evidenceSummaryInputs(t)
	chain.ChainDigest = digestString("tampered")
	summary, err := SummarizeRevisionEvidence(chain, binding, metrics)
	if err == nil {
		t.Fatal("SummarizeRevisionEvidence() error = nil, want chain failure")
	}
	if summary.Status != "UNKNOWN" || summary.MissingStage != "revision-evidence-summary-chain" {
		t.Fatalf("unexpected unknown summary: %#v", summary)
	}
}

func TestSummarizeRevisionEvidenceRetainsMetricsLinkFailure(t *testing.T) {
	chain, binding, metrics := evidenceSummaryInputs(t)
	metrics.MetricsDigest = digestString("tampered")
	summary, err := SummarizeRevisionEvidence(chain, binding, metrics)
	if err == nil {
		t.Fatal("SummarizeRevisionEvidence() error = nil, want metrics failure")
	}
	if summary.Status != "UNKNOWN" || summary.MissingStage != "revision-evidence-summary-metrics" {
		t.Fatalf("unexpected unknown metrics summary: %#v", summary)
	}
}

func TestSummarizeRevisionEvidenceRetainsLinkFailure(t *testing.T) {
	chain, binding, metrics := evidenceSummaryInputs(t)
	binding.ChainDigest = digestString("other-chain")
	binding.BindingDigest = digestRevisionEvidenceGenerationBinding(binding)
	summary, err := SummarizeRevisionEvidence(chain, binding, metrics)
	if err == nil {
		t.Fatal("SummarizeRevisionEvidence() error = nil, want link failure")
	}
	if summary.Status != "UNKNOWN" || summary.MissingStage != "revision-evidence-summary-link" {
		t.Fatalf("unexpected unknown link summary: %#v", summary)
	}
}

func TestSummarizeRevisionEvidenceIsDeterministic(t *testing.T) {
	chain, binding, metrics := evidenceSummaryInputs(t)
	first, err := SummarizeRevisionEvidence(chain, binding, metrics)
	if err != nil {
		t.Fatalf("first SummarizeRevisionEvidence() error = %v", err)
	}
	second, err := SummarizeRevisionEvidence(chain, binding, metrics)
	if err != nil {
		t.Fatalf("second SummarizeRevisionEvidence() error = %v", err)
	}
	if first.SummaryDigest != second.SummaryDigest {
		t.Fatal("same linked stages produced different summary digest")
	}
}