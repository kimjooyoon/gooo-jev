package decision

import (
	"fmt"
	"strings"
)

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionMetricInput
// adapts one evidence-backed LSP projection into the existing metric model.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionMetricInput struct {
	Projection             ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionLSPBinding
	MetricName             string
	AdditionalUnknownCount uint64
	NonAuthorizing         bool
}

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionMetricBinding
// preserves projection evidence and the measured disposition counts.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionMetricBinding struct {
	Status                    string
	MissingStage              string
	MetricName                string
	Total                     uint64
	VerifiedCount             uint64
	ReviewCount               uint64
	CounterexampleCount       uint64
	UnknownCount              uint64
	ProjectionEvidenceDigest  string
	MetricEvidenceDigest      string
	EvidenceDigest             string
	NonExecuting              bool
	NonAuthorizing            bool
}

func (b ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionMetricBinding) Validate() error {
	switch b.Status {
	case "measured", "measured-with-counterexample", "measured-with-unknown":
	default:
		return fmt.Errorf("incomplete Gooo revision action metric binding")
	}
	if b.MissingStage != "" ||
		b.MetricName == "" ||
		b.Total == 0 ||
		b.VerifiedCount+b.ReviewCount+b.CounterexampleCount+b.UnknownCount != b.Total ||
		b.ProjectionEvidenceDigest == "" ||
		b.MetricEvidenceDigest == "" ||
		b.EvidenceDigest == "" {
		return fmt.Errorf("incomplete Gooo revision action metric evidence")
	}
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("Gooo revision action metric must be non-executing and non-authorizing")
	}
	expected := digestJEVGoooEvidenceFullProvenanceRevisionActionMetric(
		b.Status,
		b.MetricName,
		b.Total,
		b.VerifiedCount,
		b.ReviewCount,
		b.CounterexampleCount,
		b.UnknownCount,
		b.ProjectionEvidenceDigest,
		b.MetricEvidenceDigest,
	)
	if b.EvidenceDigest != expected {
		return fmt.Errorf("Gooo revision action metric digest mismatch")
	}
	return nil
}

// MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionMetric
// classifies an LSP projection without claiming execution or improvement.
func MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionMetric(input ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionMetricInput) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionMetricBinding {
	unknown := func(stage string) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionMetricBinding {
		if strings.TrimSpace(stage) == "" {
			stage = "revision-action-metric"
		}
		return ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionMetricBinding{
			Status:         "UNKNOWN",
			MissingStage:   stage,
			NonExecuting:   true,
			NonAuthorizing: true,
		}
	}
	if !input.NonAuthorizing || !input.Projection.NonAuthorizing {
		return unknown("authorization-boundary")
	}
	if !input.Projection.NonExecuting {
		return unknown("execution-boundary")
	}
	if err := input.Projection.Validate(); err != nil {
		return unknown("revision-action-lsp-validation")
	}
	projectionEvidenceDigest, err := Digest(struct {
		ProjectionEvidenceDigest string
		EvidencePrefixDigest    string
		AdditionalUnknownCount  uint64
	}{
		ProjectionEvidenceDigest: input.Projection.EvidenceDigest,
		EvidencePrefixDigest:    input.Projection.EvidencePrefixDigest,
		AdditionalUnknownCount:  input.AdditionalUnknownCount,
	})
	if err != nil {
		return unknown("metric-source-evidence")
	}
	var verifiedCount uint64
	var reviewCount uint64
	var counterexampleCount uint64
	switch {
	case input.Projection.Status == "clear":
		verifiedCount = 1
	case input.Projection.Status == "publishable" && input.Projection.Code == "jev-reverse-counterexample":
		counterexampleCount = 1
	case input.Projection.Status == "publishable":
		reviewCount = 1
	default:
		return unknown("lsp-projection-status")
	}
	metric := MeasureJEVActionCandidateVerificationMetric(JEVActionCandidateVerificationMetricInput{
		MetricName:           input.MetricName,
		VerifiedCount:        verifiedCount,
		ReviewCount:          reviewCount,
		CounterexampleCount:  counterexampleCount,
		UnknownCount:         input.AdditionalUnknownCount,
		EvidenceSourceDigest: projectionEvidenceDigest,
		NonAuthorizing:       true,
	})
	if metric.Status == "UNKNOWN" {
		return unknown(metric.MissingStage)
	}
	output := ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionMetricBinding{
		Status:                   metric.Status,
		MetricName:               metric.MetricName,
		Total:                    metric.Total,
		VerifiedCount:            metric.VerifiedCount,
		ReviewCount:              metric.ReviewCount,
		CounterexampleCount:      metric.CounterexampleCount,
		UnknownCount:             metric.UnknownCount,
		ProjectionEvidenceDigest: projectionEvidenceDigest,
		MetricEvidenceDigest:     metric.EvidenceDigest,
		NonExecuting:             true,
		NonAuthorizing:           true,
	}
	output.EvidenceDigest = digestJEVGoooEvidenceFullProvenanceRevisionActionMetric(
		output.Status,
		output.MetricName,
		output.Total,
		output.VerifiedCount,
		output.ReviewCount,
		output.CounterexampleCount,
		output.UnknownCount,
		output.ProjectionEvidenceDigest,
		output.MetricEvidenceDigest,
	)
	if err := output.Validate(); err != nil {
		return unknown("revision-action-metric-evidence")
	}
	return output
}

func digestJEVGoooEvidenceFullProvenanceRevisionActionMetric(
	status,
	metricName string,
	total,
	verifiedCount,
	reviewCount,
	counterexampleCount,
	unknownCount uint64,
	projectionEvidenceDigest,
	metricEvidenceDigest string,
) string {
	digest, err := Digest(struct {
		Status                   string
		MetricName               string
		Total                    uint64
		VerifiedCount            uint64
		ReviewCount              uint64
		CounterexampleCount      uint64
		UnknownCount             uint64
		ProjectionEvidenceDigest string
		MetricEvidenceDigest     string
	}{
		Status:                   status,
		MetricName:               metricName,
		Total:                    total,
		VerifiedCount:            verifiedCount,
		ReviewCount:              reviewCount,
		CounterexampleCount:      counterexampleCount,
		UnknownCount:             unknownCount,
		ProjectionEvidenceDigest: projectionEvidenceDigest,
		MetricEvidenceDigest:     metricEvidenceDigest,
	})
	if err != nil {
		return ""
	}
	return digest
}