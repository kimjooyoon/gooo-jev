package decision

import (
	"fmt"
	"strings"
)

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionDirectionInput
// adapts revision candidate feedback into the direction model.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionDirectionInput struct {
	Feedback       ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionFeedbackBinding
	CandidateSource string
	NonAuthorizing bool
}

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionDirectionBinding
// preserves revision candidate, aggregation, and direction evidence.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionDirectionBinding struct {
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
	DirectiveStatus                 string
	Directive                       string
	DirectionCandidateDigest        string
	CandidateSource                 string
	InputEvidenceDigest             string
	DirectionEvidenceDigest         string
	EvidenceDigest                  string
	NonExecuting                    bool
	NonAuthorizing                  bool
}

func (b ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionDirectionBinding) Validate() error {
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
		b.CandidateStatus == "" ||
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
		b.DirectiveStatus == "" ||
		b.Directive == "" ||
		b.DirectionCandidateDigest == "" ||
		b.CandidateSource == "" ||
		b.InputEvidenceDigest == "" ||
		b.DirectionEvidenceDigest == "" ||
		b.EvidenceDigest == "" {
		return fmt.Errorf("incomplete Gooo extended lineage candidate revision direction binding")
	}
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("Gooo extended lineage candidate revision direction must be non-executing and non-authorizing")
	}
	expected := digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionDirection(
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
		b.DirectiveStatus,
		b.Directive,
		b.DirectionCandidateDigest,
		b.CandidateSource,
		b.InputEvidenceDigest,
		b.DirectionEvidenceDigest,
	)
	if b.EvidenceDigest != expected {
		return fmt.Errorf("Gooo extended lineage candidate revision direction digest mismatch")
	}
	return nil
}

// DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionDirection
// derives a non-authorizing direction while preserving revision candidate lineage.
func DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionDirection(input ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionDirectionInput) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionDirectionBinding {
	unknown := func(stage string) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionDirectionBinding {
		if strings.TrimSpace(stage) == "" {
			stage = "revision-action-candidate-generation-extended-lineage-candidate-revision-direction"
		}
		return ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionDirectionBinding{
			Status:         "UNKNOWN",
			MissingStage:   stage,
			NonExecuting:   true,
			NonAuthorizing: true,
		}
	}
	if !input.NonAuthorizing || !input.Feedback.NonAuthorizing {
		return unknown("authorization-boundary")
	}
	if !input.Feedback.NonExecuting {
		return unknown("execution-boundary")
	}
	if err := input.Feedback.Validate(); err != nil {
		return unknown("revision-action-candidate-generation-extended-lineage-candidate-revision-feedback-validation")
	}
	if strings.TrimSpace(input.CandidateSource) == "" {
		return unknown("candidate-source")
	}
	aggregation := JEVImprovementFeedbackAggregation{
		Status:              input.Feedback.AggregationStatus,
		Total:               input.Feedback.Total,
		Confirmed:           input.Feedback.Confirmed,
		Refuted:             input.Feedback.Refuted,
		Unknown:             input.Feedback.Unknown,
		InputEvidenceDigest: digestJEVImprovementFeedbackInputs([]string{input.Feedback.FeedbackDigest}),
		EvidenceDigest:      input.Feedback.AggregationEvidenceDigest,
		NonExecuting:        true,
		NonAuthorizing:      true,
	}
	if err := aggregation.Validate(); err != nil {
		return unknown("feedback-aggregation")
	}
	direction := DeriveJEVImprovementDirectionDirective(JEVImprovementDirectionDirectiveInput{
		Aggregation:     aggregation,
		CandidateDigest: input.Feedback.CandidateDigest,
		CandidateSource: input.CandidateSource,
		NonAuthorizing:  true,
	})
	if err := direction.Validate(); err != nil {
		stage := direction.MissingStage
		if strings.TrimSpace(stage) == "" {
			stage = "feedback-direction"
		}
		return unknown(stage)
	}
	output := ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionDirectionBinding{
		Status:                          "bound",
		CandidateID:                     input.Feedback.CandidateID,
		SourceCandidateID:               input.Feedback.SourceCandidateID,
		SourceRevisionCandidateDigest:   input.Feedback.SourceRevisionCandidateDigest,
		SourceCandidateEvidenceDigest:   input.Feedback.SourceCandidateEvidenceDigest,
		ParentCandidateDigest:           input.Feedback.ParentCandidateDigest,
		ParentCandidateEvidenceDigest:   input.Feedback.ParentCandidateEvidenceDigest,
		RevisionCandidateDigest:         input.Feedback.RevisionCandidateDigest,
		RevisionCandidateEvidenceDigest: input.Feedback.RevisionCandidateEvidenceDigest,
		CandidateStatus:                 input.Feedback.CandidateStatus,
		CandidateDigest:                 input.Feedback.CandidateDigest,
		CandidateEvidenceDigest:         input.Feedback.CandidateEvidenceDigest,
		RevisionSource:                  input.Feedback.RevisionSource,
		BoundRevisionChangeDigest:       input.Feedback.BoundRevisionChangeDigest,
		GuardEvidenceDigest:             input.Feedback.GuardEvidenceDigest,
		VerificationEvidenceDigest:      input.Feedback.VerificationEvidenceDigest,
		ReverseEvidenceDigest:           input.Feedback.ReverseEvidenceDigest,
		EvidencePrefixDigest:            input.Feedback.EvidencePrefixDigest,
		MetricEvidenceDigest:            input.Feedback.MetricEvidenceDigest,
		FeedbackDigest:                  input.Feedback.FeedbackDigest,
		AggregationEvidenceDigest:       input.Feedback.AggregationEvidenceDigest,
		DirectiveStatus:                 direction.Status,
		Directive:                       direction.Directive,
		DirectionCandidateDigest:        direction.CandidateDigest,
		CandidateSource:                 direction.CandidateSource,
		InputEvidenceDigest:             direction.InputEvidenceDigest,
		DirectionEvidenceDigest:         direction.EvidenceDigest,
		NonExecuting:                    true,
		NonAuthorizing:                  true,
	}
	output.EvidenceDigest = digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionDirection(
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
		output.DirectiveStatus,
		output.Directive,
		output.DirectionCandidateDigest,
		output.CandidateSource,
		output.InputEvidenceDigest,
		output.DirectionEvidenceDigest,
	)
	if err := output.Validate(); err != nil {
		return unknown("revision-action-candidate-generation-extended-lineage-candidate-revision-direction-evidence")
	}
	return output
}

func digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionDirection(
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
	directiveStatus,
	directive,
	directionCandidateDigest,
	candidateSource,
	inputEvidenceDigest,
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
		DirectiveStatus                 string
		Directive                       string
		DirectionCandidateDigest        string
		CandidateSource                 string
		InputEvidenceDigest             string
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
		DirectiveStatus:                 directiveStatus,
		Directive:                       directive,
		DirectionCandidateDigest:        directionCandidateDigest,
		CandidateSource:                 candidateSource,
		InputEvidenceDigest:             inputEvidenceDigest,
		DirectionEvidenceDigest:          directionEvidenceDigest,
	})
	if err != nil {
		return ""
	}
	return digest
}