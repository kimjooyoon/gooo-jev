package decision

import (
	"fmt"
	"strings"
)

// ExecutionEnvelopeGoooEvidenceFullProvenanceCandidateGateInput binds a
// proposed source revision to an evidence-linked observation metric.
type ExecutionEnvelopeGoooEvidenceFullProvenanceCandidateGateInput struct {
	Metric                 ExecutionEnvelopeGoooEvidenceFullProvenanceMetric
	CandidateSource        string
	ProposedCandidateDigest string
	NonAuthorizing         bool
}

// ExecutionEnvelopeGoooEvidenceFullProvenanceCandidateGate is an admission
// decision only; ADMIT never means execute.
type ExecutionEnvelopeGoooEvidenceFullProvenanceCandidateGate struct {
	Status              string
	MissingStage        string
	AdmissionStatus     string
	CandidateSourceDigest string
	MetricEvidenceDigest string
	EvidenceDigest      string
	NonExecuting        bool
	NonAuthorizing      bool
}

func (g ExecutionEnvelopeGoooEvidenceFullProvenanceCandidateGate) Validate() error {
	if g.Status != "ready" ||
		(g.AdmissionStatus != "ADMIT" && g.AdmissionStatus != "HOLD") ||
		g.MissingStage != "" ||
		g.CandidateSourceDigest == "" ||
		g.MetricEvidenceDigest == "" ||
		g.EvidenceDigest == "" {
		return fmt.Errorf("invalid Gooo evidence full provenance candidate gate")
	}
	if !g.NonExecuting {
		return fmt.Errorf("Gooo evidence full provenance candidate gate must be non-executing")
	}
	if !g.NonAuthorizing {
		return fmt.Errorf("Gooo evidence full provenance candidate gate must be non-authorizing")
	}
	return nil
}

func digestGoooEvidenceFullProvenanceCandidateSource(source string) (string, error) {
	return Digest(struct {
		CandidateSource string
	}{CandidateSource: source})
}

// BindExecutionEnvelopeGoooEvidenceFullProvenanceCandidateGate admits a
// candidate only when its source digest and observation metric are bound.
func BindExecutionEnvelopeGoooEvidenceFullProvenanceCandidateGate(input ExecutionEnvelopeGoooEvidenceFullProvenanceCandidateGateInput) ExecutionEnvelopeGoooEvidenceFullProvenanceCandidateGate {
	output := ExecutionEnvelopeGoooEvidenceFullProvenanceCandidateGate{
		Status: "UNKNOWN", NonExecuting: true, NonAuthorizing: true,
	}
	if !input.NonAuthorizing || !input.Metric.NonAuthorizing {
		if !input.NonAuthorizing {
			output.NonAuthorizing = false
		}
		output.MissingStage = "authorization-boundary"
		return output
	}
	if input.Metric.Status != "ready" {
		output.MissingStage = "metric"
		return output
	}
	if err := input.Metric.Validate(); err != nil {
		output.MissingStage = "metric-validation"
		return output
	}
	if strings.TrimSpace(input.CandidateSource) == "" {
		output.MissingStage = "candidate-source"
		return output
	}
	candidateDigest, err := digestGoooEvidenceFullProvenanceCandidateSource(input.CandidateSource)
	if err != nil {
		output.MissingStage = "candidate-source-digest"
		return output
	}
	if strings.TrimSpace(input.ProposedCandidateDigest) == "" ||
		input.ProposedCandidateDigest != candidateDigest {
		output.MissingStage = "candidate-source-digest"
		return output
	}
	output.Status = "ready"
	output.CandidateSourceDigest = candidateDigest
	output.MetricEvidenceDigest = input.Metric.EvidenceDigest
	if input.Metric.Classification == "CONFIRMED" {
		output.AdmissionStatus = "ADMIT"
	} else {
		output.AdmissionStatus = "HOLD"
	}
	output.EvidenceDigest, err = Digest(struct {
		AdmissionStatus       string
		CandidateSourceDigest string
		MetricEvidenceDigest  string
	}{
		AdmissionStatus:       output.AdmissionStatus,
		CandidateSourceDigest: output.CandidateSourceDigest,
		MetricEvidenceDigest:  output.MetricEvidenceDigest,
	})
	if err != nil {
		output.Status = "UNKNOWN"
		output.MissingStage = "candidate-gate"
		output.CandidateSourceDigest = ""
		output.MetricEvidenceDigest = ""
		output.EvidenceDigest = ""
		return output
	}
	if err := output.Validate(); err != nil {
		output.Status = "UNKNOWN"
		output.MissingStage = "candidate-gate"
		output.EvidenceDigest = ""
	}
	return output
}
