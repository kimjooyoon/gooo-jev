package gooo

import "fmt"

type RevisionTrendObservation struct {
	Status                     string
	MissingStage               string
	PreviousChainDigest        string
	CurrentChainDigest         string
	PreviousMetricsDigest      string
	CurrentMetricsDigest       string
	PreviousApplicationDigest  string
	CurrentApplicationDigest   string
	PreviousSourceDigest       string
	CurrentSourceDigest        string
	ChangedByteDelta           int
	ChangedLineDelta           int
	IRChangeStable             bool
	CandidateStable            bool
	SourceStable               bool
	ChangeSignal               string
	TrendDigest                string
	NonExecuting               bool
	NonAuthorizing             bool
}

func ObserveRevisionTrend(previousChain RevisionEvidenceChain, previousMetrics RevisionMetrics, currentChain RevisionEvidenceChain, currentMetrics RevisionMetrics) (RevisionTrendObservation, error) {
	trend := RevisionTrendObservation{
		Status:                    "UNKNOWN",
		MissingStage:              "revision-trend",
		PreviousChainDigest:       previousChain.ChainDigest,
		CurrentChainDigest:        currentChain.ChainDigest,
		PreviousMetricsDigest:     previousMetrics.MetricsDigest,
		CurrentMetricsDigest:      currentMetrics.MetricsDigest,
		PreviousApplicationDigest: previousMetrics.ApplicationDigest,
		CurrentApplicationDigest:  currentMetrics.ApplicationDigest,
		PreviousSourceDigest:      previousMetrics.SourceDigest,
		CurrentSourceDigest:       currentMetrics.SourceDigest,
		NonExecuting:              true,
		NonAuthorizing:            true,
	}
	if err := previousChain.Validate(); err != nil {
		trend.MissingStage = "revision-trend-previous-chain"
		return trend, fmt.Errorf("gooo revision trend: previous chain: %w", err)
	}
	if err := previousMetrics.Validate(); err != nil {
		trend.MissingStage = "revision-trend-previous-metrics"
		return trend, fmt.Errorf("gooo revision trend: previous metrics: %w", err)
	}
	if err := currentChain.Validate(); err != nil {
		trend.MissingStage = "revision-trend-current-chain"
		return trend, fmt.Errorf("gooo revision trend: current chain: %w", err)
	}
	if err := currentMetrics.Validate(); err != nil {
		trend.MissingStage = "revision-trend-current-metrics"
		return trend, fmt.Errorf("gooo revision trend: current metrics: %w", err)
	}
	if previousChain.MetricsDigest != previousMetrics.MetricsDigest ||
		previousChain.ApplicationDigest != previousMetrics.ApplicationDigest ||
		previousChain.SourceDigest != previousMetrics.SourceDigest {
		trend.MissingStage = "revision-trend-previous-link"
		return trend, fmt.Errorf("gooo revision trend: previous chain is not linked to metrics")
	}
	if currentChain.MetricsDigest != currentMetrics.MetricsDigest ||
		currentChain.ApplicationDigest != currentMetrics.ApplicationDigest ||
		currentChain.SourceDigest != currentMetrics.SourceDigest {
		trend.MissingStage = "revision-trend-current-link"
		return trend, fmt.Errorf("gooo revision trend: current chain is not linked to metrics")
	}
	if previousChain.ChainDigest == currentChain.ChainDigest {
		trend.MissingStage = "revision-trend-identity"
		return trend, fmt.Errorf("gooo revision trend: previous and current chain are identical")
	}
	trend.Status = "BOUND"
	trend.MissingStage = ""
	trend.ChangedByteDelta = currentMetrics.ChangedByteCount - previousMetrics.ChangedByteCount
	trend.ChangedLineDelta = currentMetrics.ChangedLineCount - previousMetrics.ChangedLineCount
	trend.IRChangeStable = previousMetrics.IRChanged == currentMetrics.IRChanged
	trend.CandidateStable = previousChain.CandidateDigest == currentChain.CandidateDigest
	trend.SourceStable = previousMetrics.SourceDigest == currentMetrics.SourceDigest
	switch {
	case trend.ChangedByteDelta == 0 && trend.ChangedLineDelta == 0:
		trend.ChangeSignal = "stable"
	case trend.ChangedByteDelta <= 0 && trend.ChangedLineDelta <= 0:
		trend.ChangeSignal = "narrower"
	case trend.ChangedByteDelta >= 0 && trend.ChangedLineDelta >= 0:
		trend.ChangeSignal = "wider"
	default:
		trend.ChangeSignal = "mixed"
	}
	trend.TrendDigest = digestRevisionTrendObservation(trend)
	return trend, nil
}

func (t RevisionTrendObservation) Validate() error {
	if t.Status != "BOUND" {
		return fmt.Errorf("revision trend status must be BOUND")
	}
	if t.MissingStage != "" {
		return fmt.Errorf("revision trend missing stage must be empty")
	}
	if !t.NonExecuting || !t.NonAuthorizing {
		return fmt.Errorf("revision trend must remain non-executing and non-authorizing")
	}
	for name, digest := range map[string]string{
		"previous chain": t.PreviousChainDigest, "current chain": t.CurrentChainDigest,
		"previous metrics": t.PreviousMetricsDigest, "current metrics": t.CurrentMetricsDigest,
		"previous application": t.PreviousApplicationDigest, "current application": t.CurrentApplicationDigest,
		"previous source": t.PreviousSourceDigest, "current source": t.CurrentSourceDigest,
		"trend": t.TrendDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("revision trend %s digest is invalid", name)
		}
	}
	if t.ChangeSignal != "narrower" && t.ChangeSignal != "wider" && t.ChangeSignal != "stable" && t.ChangeSignal != "mixed" {
		return fmt.Errorf("revision trend change signal is invalid")
	}
	if digestRevisionTrendObservation(t) != t.TrendDigest {
		return fmt.Errorf("revision trend digest does not match its fields")
	}
	return nil
}

func digestRevisionTrendObservation(trend RevisionTrendObservation) string {
	return digestString(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%d|%d|%t|%t|%t|%s|%t|%t",
		trend.Status,
		trend.MissingStage,
		trend.PreviousChainDigest,
		trend.CurrentChainDigest,
		trend.PreviousMetricsDigest,
		trend.CurrentMetricsDigest,
		trend.PreviousApplicationDigest,
		trend.CurrentApplicationDigest,
		trend.PreviousSourceDigest,
		trend.CurrentSourceDigest,
		trend.ChangedByteDelta,
		trend.ChangedLineDelta,
		trend.IRChangeStable,
		trend.CandidateStable,
		trend.SourceStable,
		trend.ChangeSignal,
		trend.NonExecuting,
		trend.NonAuthorizing,
	))
}
