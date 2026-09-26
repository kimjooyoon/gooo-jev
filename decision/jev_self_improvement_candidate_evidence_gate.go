package decision

import "strings"

// ExecutionEnvelopeJEVSelfImprovementCandidateEvidenceGateInput admits a
// revision candidate only after the complete observed evidence summary exists.
type ExecutionEnvelopeJEVSelfImprovementCandidateEvidenceGateInput struct {
	Summary        ExecutionEnvelopeJEVSelfImprovementEvidenceSummary
	Candidate      JEVImprovementRevisionCandidate
	NonAuthorizing bool
}

// ExecutionEnvelopeJEVSelfImprovementCandidateEvidenceGateBinding identifies a
// candidate together with the metric-backed summary that supports review.
type ExecutionEnvelopeJEVSelfImprovementCandidateEvidenceGateBinding struct {
	Status                  string
	CandidateStatus         string
	CandidateDigest         string
	CandidateEvidenceDigest string
	SummaryDigest           string
	MetricSourceDigest      string
	MetricValueDigest       string
	AdmissionDigest         string
	MissingStage            string
	NonExecuting            bool
	NonAuthorizing          bool
}

// BindExecutionEnvelopeJEVSelfImprovementSummaryToCandidate preserves the
// candidate as non-executing and non-authorizing review evidence.
func BindExecutionEnvelopeJEVSelfImprovementSummaryToCandidate(input ExecutionEnvelopeJEVSelfImprovementCandidateEvidenceGateInput) ExecutionEnvelopeJEVSelfImprovementCandidateEvidenceGateBinding {
	output := ExecutionEnvelopeJEVSelfImprovementCandidateEvidenceGateBinding{
		Status: "UNKNOWN", NonExecuting: true, NonAuthorizing: true,
	}
	if !input.NonAuthorizing || !input.Summary.NonAuthorizing || !input.Candidate.NonAuthorizing {
		if !input.NonAuthorizing {
			output.NonAuthorizing = false
		}
		output.MissingStage = "authorization-boundary"
		return output
	}
	if !input.Summary.NonExecuting {
		output.MissingStage = "execution-boundary"
		return output
	}
	if input.Summary.Status != "bound" || strings.TrimSpace(input.Summary.SummaryDigest) == "" {
		output.MissingStage = input.Summary.MissingStage
		if strings.TrimSpace(output.MissingStage) == "" {
			output.MissingStage = "self-improvement-evidence-summary"
		}
		return output
	}
	if strings.TrimSpace(input.Summary.MetricSourceDigest) == "" || strings.TrimSpace(input.Summary.MetricValueDigest) == "" {
		output.MissingStage = "metric-evidence"
		return output
	}
	if err := input.Candidate.Validate(); err != nil {
		output.MissingStage = "revision-candidate"
		return output
	}
	admissionDigest, err := Digest(struct {
		CandidateDigest         string
		CandidateEvidenceDigest string
		SummaryDigest           string
		MetricSourceDigest      string
		MetricValueDigest       string
	}{
		CandidateDigest:         input.Candidate.CandidateDigest,
		CandidateEvidenceDigest: input.Candidate.EvidenceDigest,
		SummaryDigest:           input.Summary.SummaryDigest,
		MetricSourceDigest:      input.Summary.MetricSourceDigest,
		MetricValueDigest:       input.Summary.MetricValueDigest,
	})
	if err != nil {
		output.MissingStage = "candidate-evidence-gate-digest"
		return output
	}
	output.Status = "bound"
	output.CandidateStatus = input.Candidate.Status
	output.CandidateDigest = input.Candidate.CandidateDigest
	output.CandidateEvidenceDigest = input.Candidate.EvidenceDigest
	output.SummaryDigest = input.Summary.SummaryDigest
	output.MetricSourceDigest = input.Summary.MetricSourceDigest
	output.MetricValueDigest = input.Summary.MetricValueDigest
	output.AdmissionDigest = admissionDigest
	return output
}
