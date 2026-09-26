package decision

import (
	"fmt"
	"strings"
)

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionInput
// connects candidate direction to the existing revision generator.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionInput struct {
	Feedback             ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateFeedbackBinding
	Direction            ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateDirectionBinding
	RevisionSource       string
	RevisionChangeDigest string
	NonAuthorizing       bool
}

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionBinding
// preserves source candidate lineage and generated revision provenance.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionBinding struct {
	Status                        string
	MissingStage                  string
	SourceCandidateID             string
	SourceRevisionCandidateDigest string
	SourceCandidateEvidenceDigest string
	MetricEvidenceDigest          string
	FeedbackDigest                string
	AggregationEvidenceDigest     string
	DirectionEvidenceDigest       string
	CandidateStatus               string
	CandidateDigest               string
	CandidateEvidenceDigest       string
	RevisionSource                string
	BoundRevisionChangeDigest     string
	EvidenceDigest                string
	NonExecuting                  bool
	NonAuthorizing                bool
}

func (b ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionBinding) Validate() error {
	if b.Status != "bound" ||
		b.MissingStage != "" ||
		b.SourceCandidateID == "" ||
		b.SourceRevisionCandidateDigest == "" ||
		b.SourceCandidateEvidenceDigest == "" ||
		b.MetricEvidenceDigest == "" ||
		b.FeedbackDigest == "" ||
		b.AggregationEvidenceDigest == "" ||
		b.DirectionEvidenceDigest == "" ||
		b.CandidateStatus != jevImprovementRevisionCandidateReady ||
		b.CandidateDigest == "" ||
		b.CandidateEvidenceDigest == "" ||
		b.RevisionSource == "" ||
		b.BoundRevisionChangeDigest == "" ||
		b.EvidenceDigest == "" {
		return fmt.Errorf("incomplete Gooo action-derived candidate revision binding")
	}
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("Gooo action-derived candidate revision must be non-executing and non-authorizing")
	}
	expected := digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevision(
		b.SourceCandidateID,
		b.SourceRevisionCandidateDigest,
		b.SourceCandidateEvidenceDigest,
		b.MetricEvidenceDigest,
		b.FeedbackDigest,
		b.AggregationEvidenceDigest,
		b.DirectionEvidenceDigest,
		b.CandidateStatus,
		b.CandidateDigest,
		b.CandidateEvidenceDigest,
		b.RevisionSource,
		b.BoundRevisionChangeDigest,
	)
	if b.EvidenceDigest != expected {
		return fmt.Errorf("Gooo action-derived candidate revision digest mismatch")
	}
	return nil
}

// GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevision
// generates only from a revision directive and never applies the candidate.
func GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevision(input ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionInput) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionBinding {
	unknown := func(stage string) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionBinding {
		if strings.TrimSpace(stage) == "" {
			stage = "revision-action-candidate-revision"
		}
		return ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionBinding{
			Status:         "UNKNOWN",
			MissingStage:   stage,
			NonExecuting:   true,
			NonAuthorizing: true,
		}
	}
	if !input.NonAuthorizing || !input.Feedback.NonAuthorizing || !input.Direction.NonAuthorizing {
		return unknown("authorization-boundary")
	}
	if !input.Feedback.NonExecuting || !input.Direction.NonExecuting {
		return unknown("execution-boundary")
	}
	if err := input.Feedback.Validate(); err != nil {
		return unknown("revision-action-candidate-feedback-validation")
	}
	if err := input.Direction.Validate(); err != nil {
		return unknown("revision-action-candidate-direction-validation")
	}
	if input.Direction.CandidateDigest != input.Feedback.CandidateDigest {
		return unknown("direction-feedback-candidate")
	}
	if input.Direction.InputEvidenceDigest != input.Feedback.AggregationEvidenceDigest {
		return unknown("direction-feedback-evidence")
	}
	if strings.TrimSpace(input.RevisionSource) == "" {
		return unknown("revision-source")
	}
	if strings.TrimSpace(input.RevisionChangeDigest) == "" {
		return unknown("revision-change")
	}
	boundChangeDigest, err := Digest(struct {
		RevisionChangeDigest      string
		MetricEvidenceDigest     string
		FeedbackDigest            string
		AggregationEvidenceDigest string
		DirectionEvidenceDigest  string
	}{
		RevisionChangeDigest:       input.RevisionChangeDigest,
		MetricEvidenceDigest:      input.Feedback.MetricEvidenceDigest,
		FeedbackDigest:             input.Feedback.FeedbackDigest,
		AggregationEvidenceDigest:  input.Feedback.AggregationEvidenceDigest,
		DirectionEvidenceDigest:    input.Direction.DirectionEvidenceDigest,
	})
	if err != nil {
		return unknown("revision-change-evidence")
	}
	direction := JEVImprovementDirectionDirective{
		Status:              input.Direction.DirectiveStatus,
		Directive:           input.Direction.Directive,
		CandidateDigest:     input.Direction.CandidateDigest,
		CandidateSource:     input.Direction.CandidateSource,
		InputEvidenceDigest: input.Direction.InputEvidenceDigest,
		EvidenceDigest:      input.Direction.DirectionEvidenceDigest,
		NonExecuting:        true,
		NonAuthorizing:      true,
	}
	candidate := GenerateJEVImprovementRevisionCandidate(JEVImprovementRevisionCandidateInput{
		Directive:            direction,
		RevisionSource:       input.RevisionSource,
		RevisionChangeDigest: boundChangeDigest,
		NonAuthorizing:       true,
	})
	if err := candidate.Validate(); err != nil {
		stage := candidate.MissingStage
		if strings.TrimSpace(stage) == "" {
			stage = "revision-candidate"
		}
		return unknown(stage)
	}
	output := ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionBinding{
		Status:                        "bound",
		SourceCandidateID:             input.Feedback.CandidateID,
		SourceRevisionCandidateDigest: input.Feedback.RevisionCandidateDigest,
		SourceCandidateEvidenceDigest: input.Feedback.CandidateEvidenceDigest,
		MetricEvidenceDigest:          input.Feedback.MetricEvidenceDigest,
		FeedbackDigest:                input.Feedback.FeedbackDigest,
		AggregationEvidenceDigest:     input.Feedback.AggregationEvidenceDigest,
		DirectionEvidenceDigest:       input.Direction.DirectionEvidenceDigest,
		CandidateStatus:               candidate.Status,
		CandidateDigest:               candidate.CandidateDigest,
		CandidateEvidenceDigest:       candidate.EvidenceDigest,
		RevisionSource:                candidate.RevisionSource,
		BoundRevisionChangeDigest:     boundChangeDigest,
		NonExecuting:                  true,
		NonAuthorizing:                true,
	}
	output.EvidenceDigest = digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevision(
		output.SourceCandidateID,
		output.SourceRevisionCandidateDigest,
		output.SourceCandidateEvidenceDigest,
		output.MetricEvidenceDigest,
		output.FeedbackDigest,
		output.AggregationEvidenceDigest,
		output.DirectionEvidenceDigest,
		output.CandidateStatus,
		output.CandidateDigest,
		output.CandidateEvidenceDigest,
		output.RevisionSource,
		output.BoundRevisionChangeDigest,
	)
	if err := output.Validate(); err != nil {
		return unknown("revision-action-candidate-revision-evidence")
	}
	return output
}

func digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevision(
	sourceCandidateID,
	sourceRevisionCandidateDigest,
	sourceCandidateEvidenceDigest,
	metricEvidenceDigest,
	feedbackDigest,
	aggregationEvidenceDigest,
	directionEvidenceDigest,
	candidateStatus,
	candidateDigest,
	candidateEvidenceDigest,
	revisionSource,
	boundRevisionChangeDigest string,
) string {
	digest, err := Digest(struct {
		SourceCandidateID             string
		SourceRevisionCandidateDigest string
		SourceCandidateEvidenceDigest string
		MetricEvidenceDigest          string
		FeedbackDigest                string
		AggregationEvidenceDigest     string
		DirectionEvidenceDigest       string
		CandidateStatus               string
		CandidateDigest               string
		CandidateEvidenceDigest       string
		RevisionSource                string
		BoundRevisionChangeDigest     string
	}{
		SourceCandidateID:             sourceCandidateID,
		SourceRevisionCandidateDigest: sourceRevisionCandidateDigest,
		SourceCandidateEvidenceDigest: sourceCandidateEvidenceDigest,
		MetricEvidenceDigest:          metricEvidenceDigest,
		FeedbackDigest:                feedbackDigest,
		AggregationEvidenceDigest:     aggregationEvidenceDigest,
		DirectionEvidenceDigest:       directionEvidenceDigest,
		CandidateStatus:               candidateStatus,
		CandidateDigest:               candidateDigest,
		CandidateEvidenceDigest:       candidateEvidenceDigest,
		RevisionSource:                revisionSource,
		BoundRevisionChangeDigest:     boundRevisionChangeDigest,
	})
	if err != nil {
		return ""
	}
	return digest
}