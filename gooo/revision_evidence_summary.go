package gooo

import "fmt"

// RevisionEvidenceSummary is a durable, read-only projection of linked
// revision evidence. It contains no execution or authorization decision.
type RevisionEvidenceSummary struct {
	Status                       string
	MissingStage                 string
	ChainDigest                  string
	GenerationBindingDigest      string
	SourceDigest                 string
	ProposedSourceDigest         string
	InputIRDigest                string
	ProposedIRDigest             string
	CandidateDigest              string
	ApplicationDigest            string
	MetricsDigest                string
	QualityDigest                string
	FeedbackDigest               string
	MaterializationDigest        string
	GenerationObservationDigest  string
	GenerationSourceDigest       string
	GeneratedSourceDigest        string
	GeneratedIRDigest             string
	StructureDigest              string
	EvidenceStageCount            int
	ExactSourceMatch              bool
	StructureMatch                bool
	ReverseObserved               bool
	SummaryDigest                 string
	NonExecuting                  bool
	NonAuthorizing                bool
}

// SummarizeRevisionEvidence validates the existing chain and its generation
// binding before producing an exact digest-preserving stage summary.
func SummarizeRevisionEvidence(chain RevisionEvidenceChain, generationBinding RevisionEvidenceGenerationBinding, metrics RevisionMetrics) (RevisionEvidenceSummary, error) {
	summary := RevisionEvidenceSummary{
		Status:                      "UNKNOWN",
		MissingStage:                "revision-evidence-summary",
		ChainDigest:                 chain.ChainDigest,
		GenerationBindingDigest:     generationBinding.BindingDigest,
		SourceDigest:                chain.SourceDigest,
		ProposedSourceDigest:        chain.ProposedSourceDigest,
		InputIRDigest:               chain.InputIRDigest,
		ProposedIRDigest:            chain.ProposedIRDigest,
		CandidateDigest:             chain.CandidateDigest,
		ApplicationDigest:           chain.ApplicationDigest,
		MetricsDigest:               chain.MetricsDigest,
		QualityDigest:               chain.QualityDigest,
		FeedbackDigest:              chain.FeedbackDigest,
		MaterializationDigest:       chain.MaterializationDigest,
		GenerationObservationDigest: chain.GenerationObservationDigest,
		GenerationSourceDigest:      generationBinding.GenerationSourceDigest,
		GeneratedSourceDigest:       generationBinding.GeneratedSourceDigest,
		GeneratedIRDigest:           generationBinding.GeneratedIRDigest,
		StructureDigest:             generationBinding.StructureDigest,
		EvidenceStageCount:          6,
		ExactSourceMatch:            generationBinding.ExactSourceMatch,
		StructureMatch:              generationBinding.StructureMatch,
		ReverseObserved:             generationBinding.ReverseObserved,
		NonExecuting:                true,
		NonAuthorizing:              true,
	}
	setSummaryDigest := func() {
		summary.SummaryDigest = digestRevisionEvidenceSummary(summary)
	}
	setSummaryDigest()

	if err := chain.Validate(); err != nil {
		summary.MissingStage = "revision-evidence-summary-chain"
		setSummaryDigest()
		return summary, fmt.Errorf("revision evidence chain is not valid: %w", err)
	}
	if err := generationBinding.Validate(); err != nil {
		summary.MissingStage = "revision-evidence-summary-generation-binding"
		setSummaryDigest()
		return summary, fmt.Errorf("revision evidence generation binding is not valid: %w", err)
	}
	if err := metrics.Validate(); err != nil {
		summary.MissingStage = "revision-evidence-summary-metrics"
		setSummaryDigest()
		return summary, fmt.Errorf("revision metrics are not valid: %w", err)
	}
	if chain.ChainDigest != generationBinding.ChainDigest ||
		chain.MetricsDigest != metrics.MetricsDigest ||
		chain.SourceDigest != metrics.SourceDigest ||
		chain.ProposedSourceDigest != metrics.ProposedSourceDigest ||
		chain.InputIRDigest != metrics.InputIRDigest ||
		chain.ProposedIRDigest != metrics.ProposedIRDigest ||
		chain.ApplicationDigest != metrics.ApplicationDigest {
		summary.MissingStage = "revision-evidence-summary-link"
		setSummaryDigest()
		return summary, fmt.Errorf("revision evidence summary inputs are not linked")
	}

	summary.Status = "BOUND"
	summary.MissingStage = ""
	setSummaryDigest()
	if err := summary.Validate(); err != nil {
		summary.Status = "UNKNOWN"
		summary.MissingStage = "revision-evidence-summary"
		setSummaryDigest()
		return summary, fmt.Errorf("revision evidence summary is not valid: %w", err)
	}
	return summary, nil
}

func (s RevisionEvidenceSummary) Validate() error {
	if s.Status == "" {
		return fmt.Errorf("revision evidence summary status is empty")
	}
	if s.Status == "BOUND" && s.MissingStage != "" {
		return fmt.Errorf("bound revision evidence summary has a missing stage")
	}
	if s.Status == "UNKNOWN" && s.MissingStage == "" {
		return fmt.Errorf("unknown revision evidence summary has no missing stage")
	}
	if s.EvidenceStageCount != 6 {
		return fmt.Errorf("revision evidence summary stage count must be 6")
	}
	if !s.NonExecuting || !s.NonAuthorizing {
		return fmt.Errorf("revision evidence summary must remain non-executing and non-authorizing")
	}
	if s.Status == "BOUND" && (!s.StructureMatch || !s.ReverseObserved) {
		return fmt.Errorf("bound revision evidence summary lacks reverse observation")
	}
	if digestRevisionEvidenceSummary(s) != s.SummaryDigest {
		return fmt.Errorf("revision evidence summary digest does not match its fields")
	}
	return nil
}

func digestRevisionEvidenceSummary(summary RevisionEvidenceSummary) string {
	return digestString(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%d|%t|%t|%t|%t|%t",
		summary.Status,
		summary.MissingStage,
		summary.ChainDigest,
		summary.GenerationBindingDigest,
		summary.SourceDigest,
		summary.ProposedSourceDigest,
		summary.InputIRDigest,
		summary.ProposedIRDigest,
		summary.CandidateDigest,
		summary.ApplicationDigest,
		summary.MetricsDigest,
		summary.QualityDigest,
		summary.FeedbackDigest,
		summary.MaterializationDigest,
		summary.GenerationObservationDigest,
		summary.GenerationSourceDigest,
		summary.GeneratedSourceDigest,
		summary.GeneratedIRDigest,
		summary.StructureDigest,
		summary.EvidenceStageCount,
		summary.ExactSourceMatch,
		summary.StructureMatch,
		summary.ReverseObserved,
		summary.NonExecuting,
		summary.NonAuthorizing,
	))
}