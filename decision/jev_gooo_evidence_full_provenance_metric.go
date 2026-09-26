package decision

import (
	"fmt"
	"strings"
)

// ExecutionEnvelopeGoooEvidenceFullProvenanceMetricInput binds one
// reverse-observation result to a validated extended provenance chain.
type ExecutionEnvelopeGoooEvidenceFullProvenanceMetricInput struct {
	Binding                ExecutionEnvelopeGoooEvidenceFullProvenanceBinding
	ObservedStatus         string
	ExpectedStatus         string
	ObservedMissingStage   string
	ExpectedMissingStage   string
	ObservedEvidenceDigest string
	ExpectedEvidenceDigest string
	MetricSource           string
	NonAuthorizing         bool
}

// ExecutionEnvelopeGoooEvidenceFullProvenanceMetric is a deterministic
// one-observation metric that never authorizes execution.
type ExecutionEnvelopeGoooEvidenceFullProvenanceMetric struct {
	Status               string
	MissingStage         string
	Classification       string
	ConfirmedCount       int
	RefutedCount         int
	UnknownCount         int
	ConfirmedPermille    int
	RefutedPermille      int
	UnknownPermille      int
	InputEvidenceDigest  string
	MetricSourceDigest   string
	EvidenceDigest       string
	NonExecuting         bool
	NonAuthorizing       bool
}

func (m ExecutionEnvelopeGoooEvidenceFullProvenanceMetric) Validate() error {
	if m.Status != "ready" ||
		(m.Classification != "CONFIRMED" && m.Classification != "REFUTED") ||
		m.MissingStage != "" ||
		m.ConfirmedCount+m.RefutedCount+m.UnknownCount != 1 ||
		m.ConfirmedPermille+m.RefutedPermille+m.UnknownPermille != 1000 ||
		m.InputEvidenceDigest == "" ||
		m.MetricSourceDigest == "" ||
		m.EvidenceDigest == "" {
		return fmt.Errorf("invalid Gooo evidence full provenance metric")
	}
	if m.Classification == "CONFIRMED" &&
		(m.ConfirmedCount != 1 || m.RefutedCount != 0 || m.UnknownCount != 0) {
		return fmt.Errorf("confirmed metric counts do not match classification")
	}
	if m.Classification == "REFUTED" &&
		(m.ConfirmedCount != 0 || m.RefutedCount != 1 || m.UnknownCount != 0) {
		return fmt.Errorf("refuted metric counts do not match classification")
	}
	if !m.NonExecuting {
		return fmt.Errorf("Gooo evidence full provenance metric must be non-executing")
	}
	if !m.NonAuthorizing {
		return fmt.Errorf("Gooo evidence full provenance metric must be non-authorizing")
	}
	return nil
}

func digestGoooEvidenceFullProvenanceMetricSource(source string) (string, error) {
	return Digest(struct {
		MetricSource string
	}{MetricSource: source})
}

// BindExecutionEnvelopeGoooEvidenceFullProvenanceMetric classifies one
// reverse observation only after the source, IR, generation, and provenance
// evidence chain has been replayed and validated.
func BindExecutionEnvelopeGoooEvidenceFullProvenanceMetric(input ExecutionEnvelopeGoooEvidenceFullProvenanceMetricInput) ExecutionEnvelopeGoooEvidenceFullProvenanceMetric {
	output := ExecutionEnvelopeGoooEvidenceFullProvenanceMetric{
		Status: "UNKNOWN", UnknownCount: 1, UnknownPermille: 1000,
		NonExecuting: true, NonAuthorizing: true,
	}
	if !input.NonAuthorizing || !input.Binding.NonAuthorizing {
		if !input.NonAuthorizing {
			output.NonAuthorizing = false
		}
		output.MissingStage = "authorization-boundary"
		return output
	}
	if !input.Binding.NonExecuting {
		output.MissingStage = "execution-boundary"
		return output
	}
	if input.Binding.Status != "complete" || strings.TrimSpace(input.Binding.MissingStage) != "" {
		output.MissingStage = input.Binding.MissingStage
		if strings.TrimSpace(output.MissingStage) == "" {
			output.MissingStage = "full-provenance"
		}
		return output
	}
	if err := input.Binding.Validate(); err != nil {
		output.MissingStage = "binding-validation"
		return output
	}
	if strings.TrimSpace(input.MetricSource) == "" {
		output.MissingStage = "metric-source"
		return output
	}
	if strings.TrimSpace(input.ObservedStatus) == "" ||
		strings.TrimSpace(input.ExpectedStatus) == "" ||
		strings.TrimSpace(input.ObservedEvidenceDigest) == "" ||
		strings.TrimSpace(input.ExpectedEvidenceDigest) == "" {
		output.MissingStage = "reverse-observation"
		return output
	}
	metricSourceDigest, err := digestGoooEvidenceFullProvenanceMetricSource(input.MetricSource)
	if err != nil {
		output.MissingStage = "metric-source-digest"
		return output
	}
	output.Status = "ready"
	output.UnknownCount = 0
	output.UnknownPermille = 0
	output.InputEvidenceDigest = input.Binding.EvidenceDigest
	output.MetricSourceDigest = metricSourceDigest
	if input.ObservedStatus == input.ExpectedStatus &&
		input.ObservedMissingStage == input.ExpectedMissingStage &&
		input.ObservedEvidenceDigest == input.ExpectedEvidenceDigest {
		output.Classification = "CONFIRMED"
		output.ConfirmedCount = 1
		output.ConfirmedPermille = 1000
	} else {
		output.Classification = "REFUTED"
		output.RefutedCount = 1
		output.RefutedPermille = 1000
	}
	output.EvidenceDigest, err = Digest(struct {
		Classification       string
		ConfirmedCount       int
		RefutedCount         int
		UnknownCount         int
		ConfirmedPermille    int
		RefutedPermille      int
		UnknownPermille      int
		InputEvidenceDigest  string
		MetricSourceDigest   string
		ObservedStatus        string
		ExpectedStatus        string
		ObservedMissingStage  string
		ExpectedMissingStage  string
		ObservedEvidenceDigest string
		ExpectedEvidenceDigest string
	}{
		Classification:         output.Classification,
		ConfirmedCount:         output.ConfirmedCount,
		RefutedCount:           output.RefutedCount,
		UnknownCount:           output.UnknownCount,
		ConfirmedPermille:      output.ConfirmedPermille,
		RefutedPermille:        output.RefutedPermille,
		UnknownPermille:        output.UnknownPermille,
		InputEvidenceDigest:    output.InputEvidenceDigest,
		MetricSourceDigest:     output.MetricSourceDigest,
		ObservedStatus:          input.ObservedStatus,
		ExpectedStatus:          input.ExpectedStatus,
		ObservedMissingStage:    input.ObservedMissingStage,
		ExpectedMissingStage:    input.ExpectedMissingStage,
		ObservedEvidenceDigest: input.ObservedEvidenceDigest,
		ExpectedEvidenceDigest: input.ExpectedEvidenceDigest,
	})
	if err != nil {
		output.Status = "UNKNOWN"
		output.MissingStage = "metric-binding"
		output.InputEvidenceDigest = ""
		output.MetricSourceDigest = ""
		output.EvidenceDigest = ""
		return output
	}
	if err := output.Validate(); err != nil {
		output.Status = "UNKNOWN"
		output.MissingStage = "metric-binding"
		output.EvidenceDigest = ""
	}
	return output
}
