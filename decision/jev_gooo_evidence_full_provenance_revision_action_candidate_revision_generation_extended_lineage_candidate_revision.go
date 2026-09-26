package decision

import (
	"fmt"
	"strings"
)

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionInput
// connects candidate direction and feedback to the next revision candidate.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionInput struct {
	Feedback             ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateFeedbackBinding
	Direction            ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateDirectionBinding
	RevisionSource       string
	RevisionChangeDigest string
	NonAuthorizing       bool
}

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionBinding
// preserves prior candidate lineage and next revision candidate provenance.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionBinding struct {
	Status                          string
	MissingStage                    string
	CandidateID                     string
	SourceCandidateID               string
	SourceRevisionCandidateDigest   string
	SourceCandidateEvidenceDigest   string
	ParentCandidateDigest           string
	ParentCandidateEvidenceDigest   string
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
	MetricEvidenceDigest            string
	FeedbackDigest                  string
	AggregationEvidenceDigest       string
	DirectionEvidenceDigest         string
	EvidenceDigest                  string
	NonExecuting                    bool
	NonAuthorizing                  bool
}

func (b ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionBinding) Validate() error {
	if b.Status != "bound" ||
		b.MissingStage != "" ||
		b.CandidateID == "" ||
		b.SourceCandidateID == "" ||
		b.SourceRevisionCandidateDigest == "" ||
		b.SourceCandidateEvidenceDigest == "" ||
		b.ParentCandidateDigest == "" ||
		b.ParentCandidateEvidenceDigest == "" ||
		b.RevisionCandidateDigest == "" ||
		b.RevisionCandidateEvidenceDigest == "" ||
		b.CandidateStatus != jevImprovementRevisionCandidateReady ||
		b.CandidateDigest == "" ||
		b.CandidateEvidenceDigest == "" ||
		b.RevisionSource == "" ||
		b.BoundRevisionChangeDigest == "" ||
		b.GuardEvidenceDigest == "" ||
		b.VerificationEvidenceDigest == "" ||
		b.ReverseEvidenceDigest == "" ||
		b.EvidencePrefixDigest == "" ||
		b.MetricEvidenceDigest == "" ||
		b.FeedbackDigest == "" ||
		b.AggregationEvidenceDigest == "" ||
		b.DirectionEvidenceDigest == "" ||
		b.EvidenceDigest == "" {
		return fmt.Errorf("incomplete Gooo extended lineage candidate revision generation binding")
	}
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("Gooo extended lineage candidate revision generation must be non-executing and non-authorizing")
	}
	expected := digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevision(
		b.CandidateID,
		b.SourceCandidateID,
		b.SourceRevisionCandidateDigest,
		b.SourceCandidateEvidenceDigest,
		b.ParentCandidateDigest,
		b.ParentCandidateEvidenceDigest,
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
		b.MetricEvidenceDigest,
		b.FeedbackDigest,
		b.AggregationEvidenceDigest,
		b.DirectionEvidenceDigest,
	)
	if b.EvidenceDigest != expected {
		return fmt.Errorf("Gooo extended lineage candidate revision generation digest mismatch")
	}
	return nil
}

// GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevision
// generates only from candidate direction and never applies the revision.
func GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevision(input ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionInput) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionBinding {
	unknown := func(stage string) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionBinding {
		if strings.TrimSpace(stage) == "" {
			stage = "revision-action-candidate-generation-extended-lineage-candidate-revision"
		}
		return ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionBinding{
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
		return unknown("revision-action-candidate-generation-extended-lineage-candidate-feedback-validation")
	}
	if err := input.Direction.Validate(); err != nil {
		return unknown("revision-action-candidate-generation-extended-lineage-candidate-direction-validation")
	}
	if input.Direction.DirectionCandidateDigest != input.Feedback.CandidateDigest {
		return unknown("direction-feedback-candidate")
	}
	if input.Direction.InputEvidenceDigest == "" || input.Direction.InputEvidenceDigest != input.Feedback.AggregationEvidenceDigest {
		return unknown("direction-feedback-evidence")
	}
	if strings.TrimSpace(input.RevisionSource) == "" {
		return unknown("revision-source")
	}
	if strings.TrimSpace(input.RevisionChangeDigest) == "" {
		return unknown("revision-change")
	}
	boundChangeDigest, err := Digest(struct {
		RevisionChangeDigest       string
		VerificationEvidenceDigest string
		EvidencePrefixDigest       string
		MetricEvidenceDigest       string
		FeedbackDigest             string
		AggregationEvidenceDigest  string
		DirectionEvidenceDigest    string
	}{
		RevisionChangeDigest:       input.RevisionChangeDigest,
		VerificationEvidenceDigest: input.Feedback.VerificationEvidenceDigest,
		EvidencePrefixDigest:       input.Feedback.EvidencePrefixDigest,
		MetricEvidenceDigest:       input.Feedback.MetricEvidenceDigest,
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
		CandidateDigest:     input.Direction.DirectionCandidateDigest,
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
	output := ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionBinding{
		Status:                          "bound",
		CandidateID:                     input.Feedback.CandidateID,
		SourceCandidateID:               input.Feedback.SourceCandidateID,
		SourceRevisionCandidateDigest:   input.Feedback.SourceRevisionCandidateDigest,
		SourceCandidateEvidenceDigest:   input.Feedback.SourceCandidateEvidenceDigest,
		ParentCandidateDigest:           input.Feedback.CandidateDigest,
		ParentCandidateEvidenceDigest:   input.Feedback.CandidateEvidenceDigest,
		RevisionCandidateDigest:         input.Feedback.RevisionCandidateDigest,
		RevisionCandidateEvidenceDigest: input.Feedback.RevisionCandidateEvidenceDigest,
		CandidateStatus:                 candidate.Status,
		CandidateDigest:                 candidate.CandidateDigest,
		CandidateEvidenceDigest:         candidate.EvidenceDigest,
		RevisionSource:                  candidate.RevisionSource,
		BoundRevisionChangeDigest:       boundChangeDigest,
		GuardEvidenceDigest:             input.Feedback.GuardEvidenceDigest,
		VerificationEvidenceDigest:      input.Feedback.VerificationEvidenceDigest,
		ReverseEvidenceDigest:            input.Feedback.ReverseEvidenceDigest,
		EvidencePrefixDigest:             input.Feedback.EvidencePrefixDigest,
		MetricEvidenceDigest:             input.Feedback.MetricEvidenceDigest,
		FeedbackDigest:                  input.Feedback.FeedbackDigest,
		AggregationEvidenceDigest:        input.Feedback.AggregationEvidenceDigest,
		DirectionEvidenceDigest:          input.Direction.DirectionEvidenceDigest,
		NonExecuting:                    true,
		NonAuthorizing:                  true,
	}
	output.EvidenceDigest = digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevision(
		output.CandidateID,
		output.SourceCandidateID,
		output.SourceRevisionCandidateDigest,
		output.SourceCandidateEvidenceDigest,
		output.ParentCandidateDigest,
		output.ParentCandidateEvidenceDigest,
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
		output.MetricEvidenceDigest,
		output.FeedbackDigest,
		output.AggregationEvidenceDigest,
		output.DirectionEvidenceDigest,
	)
	if err := output.Validate(); err != nil {
		return unknown("revision-action-candidate-generation-extended-lineage-candidate-revision-evidence")
	}
	return output
}

func digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevision(
	candidateID,
	sourceCandidateID,
	sourceRevisionCandidateDigest,
	sourceCandidateEvidenceDigest,
	parentCandidateDigest,
	parentCandidateEvidenceDigest,
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
	evidencePrefixDigest,
	metricEvidenceDigest,
	feedbackDigest,
	aggregationEvidenceDigest,
	directionEvidenceDigest string,
) string {
	digest, err := Digest(struct {
		CandidateID                     string
		SourceCandidateID               string
		SourceRevisionCandidateDigest   string
		SourceCandidateEvidenceDigest   string
		ParentCandidateDigest           string
		ParentCandidateEvidenceDigest   string
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
		MetricEvidenceDigest            string
		FeedbackDigest                  string
		AggregationEvidenceDigest       string
		DirectionEvidenceDigest         string
	}{
		CandidateID:                     candidateID,
		SourceCandidateID:               sourceCandidateID,
		SourceRevisionCandidateDigest:   sourceRevisionCandidateDigest,
		SourceCandidateEvidenceDigest:   sourceCandidateEvidenceDigest,
		ParentCandidateDigest:           parentCandidateDigest,
		ParentCandidateEvidenceDigest:   parentCandidateEvidenceDigest,
		RevisionCandidateDigest:         revisionCandidateDigest,
		RevisionCandidateEvidenceDigest: revisionCandidateEvidenceDigest,
		CandidateStatus:                 candidateStatus,
		CandidateDigest:                 candidateDigest,
		CandidateEvidenceDigest:         candidateEvidenceDigest,
		RevisionSource:                  revisionSource,
		BoundRevisionChangeDigest:       boundRevisionChangeDigest,
		GuardEvidenceDigest:             guardEvidenceDigest,
		VerificationEvidenceDigest:      verificationEvidenceDigest,
		ReverseEvidenceDigest:           reverseEvidenceDigest,
		EvidencePrefixDigest:            evidencePrefixDigest,
		MetricEvidenceDigest:            metricEvidenceDigest,
		FeedbackDigest:                  feedbackDigest,
		AggregationEvidenceDigest:       aggregationEvidenceDigest,
		DirectionEvidenceDigest:          directionEvidenceDigest,
	})
	if err != nil {
		return ""
	}
	return digest
}