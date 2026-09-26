package decision

import (
	"fmt"
	"strings"
)

// ExecutionEnvelopeJEVFullProvenanceReplayMetricInput binds the source
// provenance metric and the replay-ledger metric without conflating them.
type ExecutionEnvelopeJEVFullProvenanceReplayMetricInput struct {
	Provenance     ExecutionEnvelopeFullProvenanceSourceBinding
	ReplayMetric   JEVReplayLedgerMetric
	NonAuthorizing bool
}

// JEVFullProvenanceReplayMetricBinding preserves both metric origins.
type JEVFullProvenanceReplayMetricBinding struct {
	Status                    string
	MissingStage              string
	DeclarationDigest         string
	IRDigest                  string
	GenerationDigest          string
	ReverseObservationDigest  string
	ProvenanceMetricDigest    string
	ReplayMetricDigest        string
	ProvenanceEvidenceDigest  string
	ReplayMetricEvidenceDigest string
	MetricBindingDigest       string
	EvidenceDigest            string
	NonExecuting              bool
	NonAuthorizing            bool
}

func (b JEVFullProvenanceReplayMetricBinding) Validate() error {
	if b.Status != "complete" ||
		b.DeclarationDigest == "" ||
		b.IRDigest == "" ||
		b.GenerationDigest == "" ||
		b.ReverseObservationDigest == "" ||
		b.ProvenanceMetricDigest == "" ||
		b.ReplayMetricDigest == "" ||
		b.ProvenanceEvidenceDigest == "" ||
		b.ReplayMetricEvidenceDigest == "" ||
		b.MetricBindingDigest == "" ||
		b.EvidenceDigest == "" {
		return fmt.Errorf("incomplete JEV full provenance replay metric binding")
	}
	if !b.NonExecuting {
		return fmt.Errorf("JEV full provenance replay metric binding must be non-executing")
	}
	if !b.NonAuthorizing {
		return fmt.Errorf("JEV full provenance replay metric binding must be non-authorizing")
	}
	expectedBinding, err := Digest(struct {
		ProvenanceMetricDigest     string
		ReplayMetricDigest         string
		ProvenanceEvidenceDigest   string
		ReplayMetricEvidenceDigest string
	}{
		ProvenanceMetricDigest:     b.ProvenanceMetricDigest,
		ReplayMetricDigest:         b.ReplayMetricDigest,
		ProvenanceEvidenceDigest:   b.ProvenanceEvidenceDigest,
		ReplayMetricEvidenceDigest: b.ReplayMetricEvidenceDigest,
	})
	if err != nil || b.MetricBindingDigest != expectedBinding {
		return fmt.Errorf("JEV full provenance replay metric binding digest mismatch")
	}
	expectedEvidence, err := Digest(struct {
		DeclarationDigest        string
		IRDigest                 string
		GenerationDigest         string
		ReverseObservationDigest string
		ProvenanceMetricDigest   string
		ReplayMetricDigest       string
		MetricBindingDigest      string
	}{
		DeclarationDigest:        b.DeclarationDigest,
		IRDigest:                 b.IRDigest,
		GenerationDigest:         b.GenerationDigest,
		ReverseObservationDigest: b.ReverseObservationDigest,
		ProvenanceMetricDigest:   b.ProvenanceMetricDigest,
		ReplayMetricDigest:       b.ReplayMetricDigest,
		MetricBindingDigest:      b.MetricBindingDigest,
	})
	if err != nil || b.EvidenceDigest != expectedEvidence {
		return fmt.Errorf("JEV full provenance replay metric evidence mismatch")
	}
	return nil
}

// BindJEVFullProvenanceReplayMetric preserves the first incomplete boundary.
func BindJEVFullProvenanceReplayMetric(input ExecutionEnvelopeJEVFullProvenanceReplayMetricInput) JEVFullProvenanceReplayMetricBinding {
	output := JEVFullProvenanceReplayMetricBinding{
		Status:         "UNKNOWN",
		NonExecuting:   true,
		NonAuthorizing: true,
	}
	if !input.NonAuthorizing || !input.Provenance.NonAuthorizing || !input.ReplayMetric.NonAuthorizing {
		if !input.NonAuthorizing {
			output.NonAuthorizing = false
		}
		output.MissingStage = "authorization-boundary"
		return output
	}
	if !input.Provenance.NonExecuting || !input.ReplayMetric.NonExecuting {
		output.MissingStage = "execution-boundary"
		return output
	}
	if input.Provenance.Status != "complete" {
		output.MissingStage = input.Provenance.MissingStage
		if strings.TrimSpace(output.MissingStage) == "" {
			output.MissingStage = "full-provenance"
		}
		return output
	}
	provenanceStages := []struct {
		name  string
		value string
	}{
		{"declaration", input.Provenance.DeclarationDigest},
		{"ir", input.Provenance.IRDigest},
		{"generation", input.Provenance.GenerationDigest},
		{"reverse-observation", input.Provenance.ReverseObservationDigest},
		{"provenance-metric", input.Provenance.MetricDigest},
		{"provenance-evidence", input.Provenance.EvidenceDigest},
	}
	for _, stage := range provenanceStages {
		if strings.TrimSpace(stage.value) == "" {
			output.MissingStage = stage.name
			return output
		}
	}
	if err := input.ReplayMetric.Validate(); err != nil {
		output.MissingStage = "replay-metric"
		return output
	}
	output.Status = "complete"
	output.DeclarationDigest = input.Provenance.DeclarationDigest
	output.IRDigest = input.Provenance.IRDigest
	output.GenerationDigest = input.Provenance.GenerationDigest
	output.ReverseObservationDigest = input.Provenance.ReverseObservationDigest
	output.ProvenanceMetricDigest = input.Provenance.MetricDigest
	output.ReplayMetricDigest = input.ReplayMetric.MetricDigest
	output.ProvenanceEvidenceDigest = input.Provenance.EvidenceDigest
	output.ReplayMetricEvidenceDigest = input.ReplayMetric.EvidenceDigest
	var err error
	output.MetricBindingDigest, err = Digest(struct {
		ProvenanceMetricDigest     string
		ReplayMetricDigest         string
		ProvenanceEvidenceDigest   string
		ReplayMetricEvidenceDigest string
	}{
		ProvenanceMetricDigest:     output.ProvenanceMetricDigest,
		ReplayMetricDigest:         output.ReplayMetricDigest,
		ProvenanceEvidenceDigest:   output.ProvenanceEvidenceDigest,
		ReplayMetricEvidenceDigest: output.ReplayMetricEvidenceDigest,
	})
	if err != nil {
		output.Status = "UNKNOWN"
		output.MissingStage = "metric-binding"
		return output
	}
	output.EvidenceDigest, err = Digest(struct {
		DeclarationDigest        string
		IRDigest                 string
		GenerationDigest         string
		ReverseObservationDigest string
		ProvenanceMetricDigest   string
		ReplayMetricDigest       string
		MetricBindingDigest      string
	}{
		DeclarationDigest:        output.DeclarationDigest,
		IRDigest:                 output.IRDigest,
		GenerationDigest:         output.GenerationDigest,
		ReverseObservationDigest: output.ReverseObservationDigest,
		ProvenanceMetricDigest:   output.ProvenanceMetricDigest,
		ReplayMetricDigest:       output.ReplayMetricDigest,
		MetricBindingDigest:      output.MetricBindingDigest,
	})
	if err != nil {
		output.Status = "UNKNOWN"
		output.MissingStage = "metric-binding-evidence"
		output.MetricBindingDigest = ""
		return output
	}
	if err := output.Validate(); err != nil {
		output.Status = "UNKNOWN"
		output.MissingStage = "metric-binding-evidence"
		output.MetricBindingDigest = ""
		output.EvidenceDigest = ""
	}
	return output
}