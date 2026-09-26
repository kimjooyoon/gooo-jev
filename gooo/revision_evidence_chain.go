package gooo

import "fmt"

type RevisionEvidenceChain struct {
	Status                       string
	MissingStage                 string
	SourceDigest                 string
	ProposedSourceDigest        string
	InputIRDigest               string
	ProposedIRDigest            string
	CandidateDigest             string
	ApplicationDigest           string
	MetricsDigest               string
	QualityDigest               string
	FeedbackDigest              string
	MaterializationDigest       string
	GenerationObservationDigest string
	ExactSourceMatch             bool
	StructureMatch               bool
	ReverseObserved              bool
	ChainDigest                  string
	NonExecuting                 bool
	NonAuthorizing               bool
}

func ObserveRevisionEvidenceChain(application RevisionApplication, metrics RevisionMetrics, quality RevisionQuality, feedback RevisionFeedback, materialization RevisionCandidateMaterialization, generation CandidateGenerationObservation) (RevisionEvidenceChain, error) {
	chain := RevisionEvidenceChain{
		Status:                       "UNKNOWN",
		MissingStage:                 "revision-evidence-chain",
		SourceDigest:                 application.SourceDigest,
		ProposedSourceDigest:         application.ProposedSourceDigest,
		InputIRDigest:                application.InputIRDigest,
		ProposedIRDigest:             application.ProposedIRDigest,
		CandidateDigest:              application.CandidateDigest,
		ApplicationDigest:            application.ApplicationDigest,
		MetricsDigest:                metrics.MetricsDigest,
		QualityDigest:                quality.QualityDigest,
		FeedbackDigest:               feedback.FeedbackDigest,
		MaterializationDigest:        materialization.MaterializationDigest,
		GenerationObservationDigest:  generation.ObservationDigest,
		NonExecuting:                 true,
		NonAuthorizing:               true,
	}
	if err := application.Validate(); err != nil {
		chain.MissingStage = "revision-chain-application"
		return chain, fmt.Errorf("gooo revision evidence chain: application: %w", err)
	}
	if err := metrics.Validate(); err != nil {
		chain.MissingStage = "revision-chain-metrics"
		return chain, fmt.Errorf("gooo revision evidence chain: metrics: %w", err)
	}
	if err := quality.Validate(); err != nil {
		chain.MissingStage = "revision-chain-quality"
		return chain, fmt.Errorf("gooo revision evidence chain: quality: %w", err)
	}
	if err := feedback.Validate(); err != nil {
		chain.MissingStage = "revision-chain-feedback"
		return chain, fmt.Errorf("gooo revision evidence chain: feedback: %w", err)
	}
	if err := materialization.Validate(); err != nil {
		chain.MissingStage = "revision-chain-materialization"
		return chain, fmt.Errorf("gooo revision evidence chain: materialization: %w", err)
	}
	if err := generation.Validate(); err != nil {
		chain.MissingStage = "revision-chain-generation"
		return chain, fmt.Errorf("gooo revision evidence chain: generation: %w", err)
	}
	if metrics.SourceDigest != application.SourceDigest ||
		metrics.ProposedSourceDigest != application.ProposedSourceDigest ||
		metrics.InputIRDigest != application.InputIRDigest ||
		metrics.ProposedIRDigest != application.ProposedIRDigest ||
		metrics.ApplicationDigest != application.ApplicationDigest {
		chain.MissingStage = "revision-chain-link-metrics"
		return chain, fmt.Errorf("gooo revision evidence chain: metrics are not linked to application")
	}
	if quality.SourceDigest != application.SourceDigest ||
		quality.ProposedSourceDigest != application.ProposedSourceDigest ||
		quality.ApplicationDigest != application.ApplicationDigest ||
		quality.MetricsDigest != metrics.MetricsDigest {
		chain.MissingStage = "revision-chain-link-quality"
		return chain, fmt.Errorf("gooo revision evidence chain: quality is not linked to metrics")
	}
	if feedback.SourceDigest != application.SourceDigest ||
		feedback.ProposedSourceDigest != application.ProposedSourceDigest ||
		feedback.ApplicationDigest != application.ApplicationDigest ||
		feedback.MetricsDigest != metrics.MetricsDigest ||
		feedback.QualityDigest != quality.QualityDigest {
		chain.MissingStage = "revision-chain-link-feedback"
		return chain, fmt.Errorf("gooo revision evidence chain: feedback is not linked to quality")
	}
	if materialization.SourceDigest != application.SourceDigest ||
		materialization.CandidateDigest != application.CandidateDigest {
		chain.MissingStage = "revision-chain-link-materialization"
		return chain, fmt.Errorf("gooo revision evidence chain: materialization is not linked to application")
	}
	if generation.SourceDigest != application.SourceDigest ||
		generation.CandidateDigest != materialization.CandidateDigest ||
		generation.MaterializationDigest != materialization.MaterializationDigest {
		chain.MissingStage = "revision-chain-link-generation"
		return chain, fmt.Errorf("gooo revision evidence chain: generation is not linked to materialization")
	}
	chain.Status = "BOUND"
	chain.MissingStage = ""
	chain.ExactSourceMatch = generation.ExactSourceMatch
	chain.StructureMatch = generation.StructureMatch
	chain.ReverseObserved = generation.ReverseObserved
	chain.ChainDigest = digestRevisionEvidenceChain(chain)
	return chain, nil
}

func (c RevisionEvidenceChain) Validate() error {
	if c.Status != "BOUND" {
		return fmt.Errorf("revision evidence chain status must be BOUND")
	}
	if c.MissingStage != "" {
		return fmt.Errorf("revision evidence chain missing stage must be empty")
	}
	if !c.NonExecuting || !c.NonAuthorizing {
		return fmt.Errorf("revision evidence chain must remain non-executing and non-authorizing")
	}
	for name, digest := range map[string]string{
		"source": c.SourceDigest, "proposed source": c.ProposedSourceDigest,
		"input IR": c.InputIRDigest, "proposed IR": c.ProposedIRDigest,
		"candidate": c.CandidateDigest, "application": c.ApplicationDigest,
		"metrics": c.MetricsDigest, "quality": c.QualityDigest,
		"feedback": c.FeedbackDigest, "materialization": c.MaterializationDigest,
		"generation": c.GenerationObservationDigest, "chain": c.ChainDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("revision evidence chain %s digest is invalid", name)
		}
	}
	if !c.StructureMatch || !c.ReverseObserved {
		return fmt.Errorf("revision evidence chain lacks reverse observation evidence")
	}
	if digestRevisionEvidenceChain(c) != c.ChainDigest {
		return fmt.Errorf("revision evidence chain digest does not match its fields")
	}
	return nil
}

func digestRevisionEvidenceChain(chain RevisionEvidenceChain) string {
	return digestString(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%t|%t|%t|%t|%t",
		chain.Status,
		chain.MissingStage,
		chain.SourceDigest,
		chain.ProposedSourceDigest,
		chain.InputIRDigest,
		chain.ProposedIRDigest,
		chain.CandidateDigest,
		chain.ApplicationDigest,
		chain.MetricsDigest,
		chain.QualityDigest,
		chain.FeedbackDigest,
		chain.MaterializationDigest,
		chain.GenerationObservationDigest,
		chain.ExactSourceMatch,
		chain.StructureMatch,
		chain.ReverseObserved,
		chain.NonExecuting,
		chain.NonAuthorizing,
	))
}
