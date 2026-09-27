package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

type LSPRevisionSelfImprovementCycleEvidenceCoverageMetric struct {
	Status                   string
	MissingStage             string
	BridgeDigest             string
	MetricName               string
	SourceDigest             string
	CandidateSourceDigest    string
	GeneratedIRDigest        string
	ReverseObservationDigest string
	LinkedEvidenceCount      int64
	RequiredEvidenceCount    int64
	CoverageSignal           string
	MetricDigest             string
	ProjectionDigest         string
	ReadOnly                 bool
	NonExecuting             bool
	NonAuthorizing           bool
}

func ObserveLSPRevisionSelfImprovementCycleEvidenceCoverageMetric(
	metric RevisionSelfImprovementCycleEvidenceCoverageMetricObservation,
) (LSPRevisionSelfImprovementCycleEvidenceCoverageMetric, error) {
	result := LSPRevisionSelfImprovementCycleEvidenceCoverageMetric{
		Status:                   "UNKNOWN",
		MissingStage:             "lsp-revision-self-improvement-cycle-evidence-coverage-metric",
		BridgeDigest:             metric.BridgeDigest,
		MetricName:               metric.MetricName,
		SourceDigest:              metric.SourceDigest,
		CandidateSourceDigest:     metric.CandidateSourceDigest,
		GeneratedIRDigest:         metric.GeneratedIRDigest,
		ReverseObservationDigest:  metric.ReverseObservationDigest,
		LinkedEvidenceCount:       metric.LinkedEvidenceCount,
		RequiredEvidenceCount:     metric.RequiredEvidenceCount,
		CoverageSignal:            metric.CoverageSignal,
		MetricDigest:              metric.MetricDigest,
		ReadOnly:                  true,
		NonExecuting:              true,
		NonAuthorizing:            true,
	}
	setDigest := func() {
		result.ProjectionDigest = digestLSPRevisionSelfImprovementCycleEvidenceCoverageMetric(result)
	}
	setDigest()

	if err := metric.Validate(); err != nil {
		result.MissingStage = "lsp-revision-self-improvement-cycle-evidence-coverage-metric-input"
		setDigest()
		return result, fmt.Errorf("cycle evidence coverage metric is not valid: %w", err)
	}

	result.Status = "BOUND"
	result.MissingStage = ""
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "lsp-revision-self-improvement-cycle-evidence-coverage-metric"
		setDigest()
		return result, fmt.Errorf("lsp evidence coverage metric projection is not valid: %w", err)
	}
	return result, nil
}

func (o LSPRevisionSelfImprovementCycleEvidenceCoverageMetric) Validate() error {
	if o.Status != "BOUND" && o.Status != "UNKNOWN" {
		return fmt.Errorf("lsp evidence coverage metric status is invalid")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound lsp evidence coverage metric has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown lsp evidence coverage metric has no missing stage")
	}
	for name, digest := range map[string]string{
		"bridge":              o.BridgeDigest,
		"source":              o.SourceDigest,
		"candidate source":    o.CandidateSourceDigest,
		"generated ir":        o.GeneratedIRDigest,
		"reverse observation": o.ReverseObservationDigest,
		"metric":              o.MetricDigest,
		"projection":          o.ProjectionDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("lsp evidence coverage metric %s digest is invalid", name)
		}
	}
	if o.MetricName != "source-ir-generation-reverse-observation-coverage" {
		return fmt.Errorf("lsp evidence coverage metric name is invalid")
	}
	if o.LinkedEvidenceCount < 0 || o.RequiredEvidenceCount < 0 {
		return fmt.Errorf("lsp evidence coverage metric counts are negative")
	}
	if o.CoverageSignal != "evidence-complete" && o.CoverageSignal != "evidence-incomplete" {
		return fmt.Errorf("lsp evidence coverage metric signal is invalid")
	}
	if o.Status == "BOUND" && (o.LinkedEvidenceCount != 4 || o.RequiredEvidenceCount != 4 || o.CoverageSignal != "evidence-complete") {
		return fmt.Errorf("bound lsp evidence coverage metric is incomplete")
	}
	if !o.ReadOnly || !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("lsp evidence coverage metric must remain read-only and non-executing")
	}
	if o.ProjectionDigest != digestLSPRevisionSelfImprovementCycleEvidenceCoverageMetric(o) {
		return fmt.Errorf("lsp evidence coverage metric projection digest does not match its fields")
	}
	return nil
}

func digestLSPRevisionSelfImprovementCycleEvidenceCoverageMetric(
	metric LSPRevisionSelfImprovementCycleEvidenceCoverageMetric,
) string {
	return digestString(strings.Join([]string{
		metric.Status,
		metric.MissingStage,
		metric.BridgeDigest,
		metric.MetricName,
		metric.SourceDigest,
		metric.CandidateSourceDigest,
		metric.GeneratedIRDigest,
		metric.ReverseObservationDigest,
		strconv.FormatInt(metric.LinkedEvidenceCount, 10),
		strconv.FormatInt(metric.RequiredEvidenceCount, 10),
		metric.CoverageSignal,
		metric.MetricDigest,
		strconv.FormatBool(metric.ReadOnly),
		strconv.FormatBool(metric.NonExecuting),
		strconv.FormatBool(metric.NonAuthorizing),
	}, "|"))
}
