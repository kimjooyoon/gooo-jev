package decision

import (
	"fmt"
	"strings"
)

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateMetricInput
// adapts extended lineage candidate LSP state into measured evidence.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateMetricInput struct {
	Projection             ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateLSPBinding
	MetricName             string
	AdditionalUnknownCount uint64
	NonAuthorizing         bool
}

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateMetricBinding
// preserves extended lineage candidate and measured disposition counts.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateMetricBinding struct {
	Status                          string
	MissingStage                    string
	CandidateID                     string
	SourceCandidateID               string
	SourceRevisionCandidateDigest   string
	SourceCandidateEvidenceDigest   string
	RevisionCandidateDigest         string
	RevisionCandidateEvidenceDigest string
	CandidateStatus                 string
	CandidateDigest                 string
	CandidateEvidenceDigest         string
	RevisionSource                  string
	BoundRevisionChangeDigest       string
	GuardEvidenceDigest             string
	VerificationEvidenceDigest      string
	ReverseEvidenceDigest           string
	EvidencePrefixDigest            string
	MetricName                      string
	Total                           uint64
	VerifiedCount                   uint64
	ReviewCount                     uint64
	CounterexampleCount             uint64
	UnknownCount                    uint64
	ProjectionEvidenceDigest        string
	MetricEvidenceDigest            string
	EvidenceDigest                 string
	NonExecuting                    bool
	NonAuthorizing                  bool
}

func (b ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateMetricBinding) Validate() error {
	switch b.Status {
	case "measured", "measured-with-counterexample", "measured-with-unknown":
	default:
		return fmt.Errorf("incomplete Gooo extended lineage candidate metric binding")
	}
	if b.MissingStage != "" ||
		b.CandidateID == "" ||
		b.SourceCandidateID == "" ||
		b.SourceRevisionCandidateDigest == "" ||
		b.SourceCandidateEvidenceDigest == "" ||
		b.RevisionCandidateDigest == "" ||
		b.RevisionCandidateEvidenceDigest == "" ||
		b.CandidateStatus == "" ||
		b.CandidateDigest == "" ||
		b.CandidateEvidenceDigest == "" ||
		b.RevisionSource == "" ||
		b.BoundRevisionChangeDigest == "" ||
		b.GuardEvidenceDigest == "" ||
		b.VerificationEvidenceDigest == "" ||
		b.ReverseEvidenceDigest == "" ||
		b.EvidencePrefixDigest == "" ||
		b.MetricName == "" ||
		b.Total == 0 ||
		b.VerifiedCount+b.ReviewCount+b.CounterexampleCount+b.UnknownCount != b.Total ||
		b.ProjectionEvidenceDigest == "" ||
		b.MetricEvidenceDigest == "" ||
		b.EvidenceDigest == "" {
		return fmt.Errorf("incomplete Gooo extended lineage candidate metric evidence")
	}
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("Gooo extended lineage candidate metric must be non-executing and non-authorizing")
	}
	expected := digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateMetric(
		b.Status,
		b.MetricName,
		b.CandidateID,
		b.SourceCandidateID,
		b.SourceRevisionCandidateDigest,
		b.SourceCandidateEvidenceDigest,
		b.RevisionCandidateDigest,
		b.RevisionCandidateEvidenceDigest,
		b.CandidateStatus,
		b.CandidateDigest,
		b.CandidateEvidenceDigest,
		b.RevisionSource,
		b.BoundRevisionChangeDigest,
		b.GuardEvidenceDigest,
		b.VerificationEvidenceDigest,
		b.ReverseEvidenceDigest,
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
		return fmt.Errorf("Gooo extended lineage candidate metric digest mismatch")
	}
	return nil
}

// MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateMetric
// classifies extended lineage candidate LSP state without claiming execution or improvement.
func MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateMetric(input ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateMetricInput) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateMetricBinding {
	unknown := func(stage string) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateMetricBinding {
		if strings.TrimSpace(stage) == "" {
			stage = "revision-action-candidate-generation-extended-lineage-candidate-metric"
		}
		return ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateMetricBinding{
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
		return unknown("revision-action-candidate-generation-extended-lineage-candidate-lsp-validation")
	}
	projectionEvidenceDigest, err := Digest(struct {
		ProjectionEvidenceDigest        string
		GuardEvidenceDigest             string
		VerificationEvidenceDigest      string
		ReverseEvidenceDigest           string
		CandidateID                     string
		SourceCandidateID               string
		SourceRevisionCandidateDigest   string
		RevisionCandidateDigest         string
		RevisionCandidateEvidenceDigest string
		CandidateDigest                 string
		CandidateEvidenceDigest         string
		EvidencePrefixDigest            string
		AdditionalUnknownCount          uint64
	}{
		ProjectionEvidenceDigest:        input.Projection.EvidenceDigest,
		GuardEvidenceDigest:             input.Projection.GuardEvidenceDigest,
		VerificationEvidenceDigest:      input.Projection.VerificationEvidenceDigest,
		ReverseEvidenceDigest:            input.Projection.ReverseEvidenceDigest,
		CandidateID:                     input.Projection.CandidateID,
		SourceCandidateID:               input.Projection.SourceCandidateID,
		SourceRevisionCandidateDigest:   input.Projection.SourceRevisionCandidateDigest,
		RevisionCandidateDigest:         input.Projection.RevisionCandidateDigest,
		RevisionCandidateEvidenceDigest: input.Projection.RevisionCandidateEvidenceDigest,
		CandidateDigest:                 input.Projection.CandidateDigest,
		CandidateEvidenceDigest:         input.Projection.CandidateEvidenceDigest,
		EvidencePrefixDigest:            input.Projection.EvidencePrefixDigest,
		AdditionalUnknownCount:          input.AdditionalUnknownCount,
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
	output := ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateMetricBinding{
		Status:                          metric.Status,
		CandidateID:                     input.Projection.CandidateID,
		SourceCandidateID:               input.Projection.SourceCandidateID,
		SourceRevisionCandidateDigest:   input.Projection.SourceRevisionCandidateDigest,
		SourceCandidateEvidenceDigest:   input.Projection.SourceCandidateEvidenceDigest,
		RevisionCandidateDigest:         input.Projection.RevisionCandidateDigest,
		RevisionCandidateEvidenceDigest: input.Projection.RevisionCandidateEvidenceDigest,
		CandidateStatus:                 input.Projection.CandidateStatus,
		CandidateDigest:                 input.Projection.CandidateDigest,
		CandidateEvidenceDigest:         input.Projection.CandidateEvidenceDigest,
		RevisionSource:                  input.Projection.RevisionSource,
		BoundRevisionChangeDigest:       input.Projection.BoundRevisionChangeDigest,
		GuardEvidenceDigest:             input.Projection.GuardEvidenceDigest,
		VerificationEvidenceDigest:      input.Projection.VerificationEvidenceDigest,
		ReverseEvidenceDigest:           input.Projection.ReverseEvidenceDigest,
		EvidencePrefixDigest:            input.Projection.EvidencePrefixDigest,
		MetricName:                      metric.MetricName,
		Total:                           metric.Total,
		VerifiedCount:                   metric.VerifiedCount,
		ReviewCount:                     metric.ReviewCount,
		CounterexampleCount:             metric.CounterexampleCount,
		UnknownCount:                    metric.UnknownCount,
		ProjectionEvidenceDigest:        projectionEvidenceDigest,
		MetricEvidenceDigest:            metric.EvidenceDigest,
		NonExecuting:                    true,
		NonAuthorizing:                  true,
	}
	output.EvidenceDigest = digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateMetric(
		output.Status,
		output.MetricName,
		output.CandidateID,
		output.SourceCandidateID,
		output.SourceRevisionCandidateDigest,
		output.SourceCandidateEvidenceDigest,
		output.RevisionCandidateDigest,
		output.RevisionCandidateEvidenceDigest,
		output.CandidateStatus,
		output.CandidateDigest,
		output.CandidateEvidenceDigest,
		output.RevisionSource,
		output.BoundRevisionChangeDigest,
		output.GuardEvidenceDigest,
		output.VerificationEvidenceDigest,
		output.ReverseEvidenceDigest,
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
		return unknown("revision-action-candidate-generation-extended-lineage-candidate-metric-evidence")
	}
	return output
}

func digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateMetric(
	status,
	metricName,
	candidateID,
	sourceCandidateID,
	sourceRevisionCandidateDigest,
	sourceCandidateEvidenceDigest,
	revisionCandidateDigest,
	revisionCandidateEvidenceDigest,
	candidateStatus,
	candidateDigest,
	candidateEvidenceDigest,
	revisionSource,
	boundRevisionChangeDigest,
	guardEvidenceDigest,
	verificationEvidenceDigest,
	reverseEvidenceDigest,
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
		CandidateStatus                 string
		CandidateDigest                 string
		CandidateEvidenceDigest         string
		RevisionSource                  string
		BoundRevisionChangeDigest       string
		GuardEvidenceDigest             string
		VerificationEvidenceDigest      string
		ReverseEvidenceDigest           string
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
		CandidateStatus:                 candidateStatus,
		CandidateDigest:                 candidateDigest,
		CandidateEvidenceDigest:         candidateEvidenceDigest,
		RevisionSource:                  revisionSource,
		BoundRevisionChangeDigest:       boundRevisionChangeDigest,
		GuardEvidenceDigest:             guardEvidenceDigest,
		VerificationEvidenceDigest:      verificationEvidenceDigest,
		ReverseEvidenceDigest:            reverseEvidenceDigest,
		EvidencePrefixDigest:            evidencePrefixDigest,
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