package decision

import (
	"fmt"
	"strings"
)

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateMetricInput
// adapts one action-derived candidate LSP projection into the metric model.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateMetricInput struct {
	Projection             ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateLSPBinding
	MetricName             string
	AdditionalUnknownCount uint64
	NonAuthorizing         bool
}

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateMetricBinding
// preserves candidate lineage and measured disposition counts.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateMetricBinding struct {
	Status                     string
	MissingStage               string
	CandidateID                string
	RevisionCandidateDigest    string
	CandidateEvidenceDigest    string
	GuardEvidenceDigest        string
	VerificationEvidenceDigest string
	EvidencePrefixDigest       string
	MetricName                 string
	Total                     uint64
	VerifiedCount             uint64
	ReviewCount               uint64
	CounterexampleCount       uint64
	UnknownCount               uint64
	ProjectionEvidenceDigest   string
	MetricEvidenceDigest       string
	EvidenceDigest             string
	NonExecuting               bool
	NonAuthorizing             bool
}

func (b ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateMetricBinding) Validate() error {
	switch b.Status {
	case "measured", "measured-with-counterexample", "measured-with-unknown":
	default:
		return fmt.Errorf("incomplete Gooo action-derived candidate metric binding")
	}
	if b.MissingStage != "" ||
		b.CandidateID == "" ||
		b.RevisionCandidateDigest == "" ||
		b.CandidateEvidenceDigest == "" ||
		b.GuardEvidenceDigest == "" ||
		b.VerificationEvidenceDigest == "" ||
		b.EvidencePrefixDigest == "" ||
		b.MetricName == "" ||
		b.Total == 0 ||
		b.VerifiedCount+b.ReviewCount+b.CounterexampleCount+b.UnknownCount != b.Total ||
		b.ProjectionEvidenceDigest == "" ||
		b.MetricEvidenceDigest == "" ||
		b.EvidenceDigest == "" {
		return fmt.Errorf("incomplete Gooo action-derived candidate metric evidence")
	}
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("Gooo action-derived candidate metric must be non-executing and non-authorizing")
	}
	expected := digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateMetric(
		b.Status,
		b.MetricName,
		b.CandidateID,
		b.RevisionCandidateDigest,
		b.CandidateEvidenceDigest,
		b.GuardEvidenceDigest,
		b.VerificationEvidenceDigest,
		b.EvidencePrefixDigest,
		b.Total,
		b.VerifiedCount,
		b.ReviewCount,
		b.CounterexampleCount,
		b.UnknownCount,
		b.ProjectionEvidenceDigest,
		b.MetricEvidenceDigest,
	)
	if b.EvidenceDigest != expected {
		return fmt.Errorf("Gooo action-derived candidate metric digest mismatch")
	}
	return nil
}

// MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateMetric
// classifies an LSP projection without claiming execution or improvement.
func MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateMetric(input ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateMetricInput) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateMetricBinding {
	unknown := func(stage string) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateMetricBinding {
		if strings.TrimSpace(stage) == "" {
			stage = "revision-action-candidate-metric"
		}
		return ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateMetricBinding{
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
		return unknown("revision-action-candidate-lsp-validation")
	}
	projectionEvidenceDigest, err := Digest(struct {
		ProjectionEvidenceDigest string
		CandidateID             string
		RevisionCandidateDigest string
		CandidateEvidenceDigest string
		EvidencePrefixDigest    string
		AdditionalUnknownCount  uint64
	}{
		ProjectionEvidenceDigest: input.Projection.EvidenceDigest,
		CandidateID:              input.Projection.CandidateID,
		RevisionCandidateDigest:  input.Projection.RevisionCandidateDigest,
		CandidateEvidenceDigest:  input.Projection.CandidateEvidenceDigest,
		EvidencePrefixDigest:     input.Projection.EvidencePrefixDigest,
		AdditionalUnknownCount:   input.AdditionalUnknownCount,
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
	output := ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateMetricBinding{
		Status:                     metric.Status,
		CandidateID:                input.Projection.CandidateID,
		RevisionCandidateDigest:    input.Projection.RevisionCandidateDigest,
		CandidateEvidenceDigest:    input.Projection.CandidateEvidenceDigest,
		GuardEvidenceDigest:        input.Projection.GuardEvidenceDigest,
		VerificationEvidenceDigest: input.Projection.VerificationEvidenceDigest,
		EvidencePrefixDigest:       input.Projection.EvidencePrefixDigest,
		MetricName:                 metric.MetricName,
		Total:                     metric.Total,
		VerifiedCount:             metric.VerifiedCount,
		ReviewCount:               metric.ReviewCount,
		CounterexampleCount:       metric.CounterexampleCount,
		UnknownCount:              metric.UnknownCount,
		ProjectionEvidenceDigest:   projectionEvidenceDigest,
		MetricEvidenceDigest:       metric.EvidenceDigest,
		NonExecuting:               true,
		NonAuthorizing:             true,
	}
	output.EvidenceDigest = digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateMetric(
		output.Status,
		output.MetricName,
		output.CandidateID,
		output.RevisionCandidateDigest,
		output.CandidateEvidenceDigest,
		output.GuardEvidenceDigest,
		output.VerificationEvidenceDigest,
		output.EvidencePrefixDigest,
		output.Total,
		output.VerifiedCount,
		output.ReviewCount,
		output.CounterexampleCount,
		output.UnknownCount,
		output.ProjectionEvidenceDigest,
		output.MetricEvidenceDigest,
	)
	if err := output.Validate(); err != nil {
		return unknown("revision-action-candidate-metric-evidence")
	}
	return output
}

func digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateMetric(
	status,
	metricName,
	candidateID,
	revisionCandidateDigest,
	candidateEvidenceDigest,
	guardEvidenceDigest,
	verificationEvidenceDigest,
	evidencePrefixDigest string,
	total,
	verifiedCount,
	reviewCount,
	counterexampleCount,
	unknownCount uint64,
	projectionEvidenceDigest,
	metricEvidenceDigest string,
) string {
	digest, err := Digest(struct {
		Status                     string
		MetricName                 string
		CandidateID                string
		RevisionCandidateDigest    string
		CandidateEvidenceDigest    string
		GuardEvidenceDigest        string
		VerificationEvidenceDigest string
		EvidencePrefixDigest       string
		Total                      uint64
		VerifiedCount              uint64
		ReviewCount                uint64
		CounterexampleCount        uint64
		UnknownCount               uint64
		ProjectionEvidenceDigest   string
		MetricEvidenceDigest       string
	}{
		Status:                     status,
		MetricName:                 metricName,
		CandidateID:                candidateID,
		RevisionCandidateDigest:    revisionCandidateDigest,
		CandidateEvidenceDigest:    candidateEvidenceDigest,
		GuardEvidenceDigest:        guardEvidenceDigest,
		VerificationEvidenceDigest: verificationEvidenceDigest,
		EvidencePrefixDigest:       evidencePrefixDigest,
		Total:                      total,
		VerifiedCount:              verifiedCount,
		ReviewCount:                reviewCount,
		CounterexampleCount:        counterexampleCount,
		UnknownCount:               unknownCount,
		ProjectionEvidenceDigest:   projectionEvidenceDigest,
		MetricEvidenceDigest:       metricEvidenceDigest,
	})
	if err != nil {
		return ""
	}
	return digest
}