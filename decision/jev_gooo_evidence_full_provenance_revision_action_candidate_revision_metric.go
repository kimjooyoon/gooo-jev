package decision

import (
	"fmt"
	"strings"
)

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionMetricInput
// adapts a candidate revision LSP projection into the metric model.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionMetricInput struct {
	Projection             ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionLSPBinding
	MetricName             string
	AdditionalUnknownCount uint64
	NonAuthorizing         bool
}

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionMetricBinding
// preserves candidate revision lineage and measured disposition counts.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionMetricBinding struct {
	Status                          string
	MissingStage                    string
	CandidateID                     string
	SourceCandidateID               string
	SourceRevisionCandidateDigest   string
	SourceCandidateEvidenceDigest   string
	RevisionCandidateDigest         string
	RevisionCandidateEvidenceDigest string
	VerificationEvidenceDigest      string
	EvidencePrefixDigest            string
	MetricName                      string
	Total                           uint64
	VerifiedCount                   uint64
	ReviewCount                     uint64
	CounterexampleCount             uint64
	UnknownCount                    uint64
	ProjectionEvidenceDigest        string
	MetricEvidenceDigest            string
	EvidenceDigest                  string
	NonExecuting                    bool
	NonAuthorizing                  bool
}

func (b ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionMetricBinding) Validate() error {
	switch b.Status {
	case "measured", "measured-with-counterexample", "measured-with-unknown":
	default:
		return fmt.Errorf("incomplete Gooo action-derived candidate revision metric binding")
	}
	if b.MissingStage != "" ||
		b.CandidateID == "" ||
		b.SourceCandidateID == "" ||
		b.SourceRevisionCandidateDigest == "" ||
		b.SourceCandidateEvidenceDigest == "" ||
		b.RevisionCandidateDigest == "" ||
		b.RevisionCandidateEvidenceDigest == "" ||
		b.VerificationEvidenceDigest == "" ||
		b.EvidencePrefixDigest == "" ||
		b.MetricName == "" ||
		b.Total == 0 ||
		b.VerifiedCount+b.ReviewCount+b.CounterexampleCount+b.UnknownCount != b.Total ||
		b.ProjectionEvidenceDigest == "" ||
		b.MetricEvidenceDigest == "" ||
		b.EvidenceDigest == "" {
		return fmt.Errorf("incomplete Gooo action-derived candidate revision metric evidence")
	}
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("Gooo action-derived candidate revision metric must be non-executing and non-authorizing")
	}
	expected := digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionMetric(
		b.Status,
		b.MetricName,
		b.CandidateID,
		b.SourceCandidateID,
		b.SourceRevisionCandidateDigest,
		b.SourceCandidateEvidenceDigest,
		b.RevisionCandidateDigest,
		b.RevisionCandidateEvidenceDigest,
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
		return fmt.Errorf("Gooo action-derived candidate revision metric digest mismatch")
	}
	return nil
}

// MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionMetric
// classifies an LSP projection without claiming execution or improvement.
func MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionMetric(input ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionMetricInput) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionMetricBinding {
	unknown := func(stage string) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionMetricBinding {
		if strings.TrimSpace(stage) == "" {
			stage = "revision-action-candidate-revision-metric"
		}
		return ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionMetricBinding{
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
		return unknown("revision-action-candidate-revision-lsp-validation")
	}
	projectionEvidenceDigest, err := Digest(struct {
		ProjectionEvidenceDigest        string
		CandidateID                     string
		SourceCandidateID               string
		SourceRevisionCandidateDigest   string
		RevisionCandidateDigest         string
		RevisionCandidateEvidenceDigest string
		EvidencePrefixDigest            string
		AdditionalUnknownCount          uint64
	}{
		ProjectionEvidenceDigest:        input.Projection.EvidenceDigest,
		CandidateID:                     input.Projection.CandidateID,
		SourceCandidateID:               input.Projection.SourceCandidateID,
		SourceRevisionCandidateDigest:   input.Projection.SourceRevisionCandidateDigest,
		RevisionCandidateDigest:         input.Projection.RevisionCandidateDigest,
		RevisionCandidateEvidenceDigest: input.Projection.RevisionCandidateEvidenceDigest,
		EvidencePrefixDigest:             input.Projection.EvidencePrefixDigest,
		AdditionalUnknownCount:           input.AdditionalUnknownCount,
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
	output := ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionMetricBinding{
		Status:                          metric.Status,
		CandidateID:                     input.Projection.CandidateID,
		SourceCandidateID:               input.Projection.SourceCandidateID,
		SourceRevisionCandidateDigest:   input.Projection.SourceRevisionCandidateDigest,
		SourceCandidateEvidenceDigest:   input.Projection.SourceCandidateEvidenceDigest,
		RevisionCandidateDigest:         input.Projection.RevisionCandidateDigest,
		RevisionCandidateEvidenceDigest: input.Projection.RevisionCandidateEvidenceDigest,
		VerificationEvidenceDigest:      input.Projection.VerificationEvidenceDigest,
		EvidencePrefixDigest:            input.Projection.EvidencePrefixDigest,
		MetricName:                      metric.MetricName,
		Total:                           metric.Total,
		VerifiedCount:                   metric.VerifiedCount,
		ReviewCount:                     metric.ReviewCount,
		CounterexampleCount:             metric.CounterexampleCount,
		UnknownCount:                   metric.UnknownCount,
		ProjectionEvidenceDigest:        projectionEvidenceDigest,
		MetricEvidenceDigest:            metric.EvidenceDigest,
		NonExecuting:                    true,
		NonAuthorizing:                  true,
	}
	output.EvidenceDigest = digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionMetric(
		output.Status,
		output.MetricName,
		output.CandidateID,
		output.SourceCandidateID,
		output.SourceRevisionCandidateDigest,
		output.SourceCandidateEvidenceDigest,
		output.RevisionCandidateDigest,
		output.RevisionCandidateEvidenceDigest,
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
		return unknown("revision-action-candidate-revision-metric-evidence")
	}
	return output
}

func digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionMetric(
	status,
	metricName,
	candidateID,
	sourceCandidateID,
	sourceRevisionCandidateDigest,
	sourceCandidateEvidenceDigest,
	revisionCandidateDigest,
	revisionCandidateEvidenceDigest,
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
		Status                          string
		MetricName                      string
		CandidateID                     string
		SourceCandidateID               string
		SourceRevisionCandidateDigest   string
		SourceCandidateEvidenceDigest   string
		RevisionCandidateDigest         string
		RevisionCandidateEvidenceDigest string
		VerificationEvidenceDigest      string
		EvidencePrefixDigest            string
		Total                           uint64
		VerifiedCount                   uint64
		ReviewCount                     uint64
		CounterexampleCount             uint64
		UnknownCount                    uint64
		ProjectionEvidenceDigest        string
		MetricEvidenceDigest            string
	}{
		Status:                          status,
		MetricName:                      metricName,
		CandidateID:                     candidateID,
		SourceCandidateID:               sourceCandidateID,
		SourceRevisionCandidateDigest:   sourceRevisionCandidateDigest,
		SourceCandidateEvidenceDigest:   sourceCandidateEvidenceDigest,
		RevisionCandidateDigest:         revisionCandidateDigest,
		RevisionCandidateEvidenceDigest: revisionCandidateEvidenceDigest,
		VerificationEvidenceDigest:      verificationEvidenceDigest,
		EvidencePrefixDigest:             evidencePrefixDigest,
		Total:                           total,
		VerifiedCount:                   verifiedCount,
		ReviewCount:                     reviewCount,
		CounterexampleCount:             counterexampleCount,
		UnknownCount:                    unknownCount,
		ProjectionEvidenceDigest:        projectionEvidenceDigest,
		MetricEvidenceDigest:            metricEvidenceDigest,
	})
	if err != nil {
		return ""
	}
	return digest
}