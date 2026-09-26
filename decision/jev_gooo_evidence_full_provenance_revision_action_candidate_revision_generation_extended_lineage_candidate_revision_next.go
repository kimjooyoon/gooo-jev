package decision

import (
	"fmt"
	"strings"
)

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionNextInput
// advances one non-executing revision candidate from a previously observed direction.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionNextInput struct {
	Direction                       ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionDirectionBinding
	PreviousCandidateDigest         string
	PreviousCandidateEvidenceDigest string
	CandidateSource                 string
	RevisionSource                 string
	RevisionChangeDigest            string
	NonAuthorizing                  bool
}

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionNextBinding
// preserves the complete prior evidence chain while binding the next revision candidate.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionNextBinding struct {
	Status                       string
	MissingStage                 string
	CandidateStatus              string
	PreviousCandidateDigest      string
	PreviousCandidateEvidenceDigest string
	ParentCandidateDigest        string
	CandidateSource              string
	RevisionSource               string
	SourceRevisionChangeDigest   string
	BoundRevisionChangeDigest    string
	DirectiveEvidenceDigest      string
	CandidateDigest              string
	CandidateEvidenceDigest      string
	DirectionCandidateDigest     string
	InputEvidenceDigest          string
	DirectionEvidenceDigest      string
	GuardEvidenceDigest          string
	VerificationEvidenceDigest  string
	ReverseEvidenceDigest        string
	EvidencePrefixDigest         string
	MetricEvidenceDigest         string
	FeedbackDigest               string
	AggregationEvidenceDigest    string
	EvidenceDigest               string
	NonExecuting                 bool
	NonAuthorizing               bool
}

func (b ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionNextBinding) Validate() error {
	if b.Status != "bound" ||
		b.MissingStage != "" ||
		b.CandidateStatus != jevImprovementRevisionCandidateReady ||
		b.PreviousCandidateDigest == "" ||
		b.PreviousCandidateEvidenceDigest == "" ||
		b.ParentCandidateDigest == "" ||
		b.CandidateSource == "" ||
		b.RevisionSource == "" ||
		b.SourceRevisionChangeDigest == "" ||
		b.BoundRevisionChangeDigest == "" ||
		b.DirectiveEvidenceDigest == "" ||
		b.CandidateDigest == "" ||
		b.CandidateEvidenceDigest == "" ||
		b.DirectionCandidateDigest == "" ||
		b.InputEvidenceDigest == "" ||
		b.DirectionEvidenceDigest == "" ||
		b.GuardEvidenceDigest == "" ||
		b.VerificationEvidenceDigest == "" ||
		b.ReverseEvidenceDigest == "" ||
		b.EvidencePrefixDigest == "" ||
		b.MetricEvidenceDigest == "" ||
		b.FeedbackDigest == "" ||
		b.AggregationEvidenceDigest == "" ||
		b.EvidenceDigest == "" {
		return fmt.Errorf("incomplete Gooo extended lineage next candidate revision binding")
	}
	if b.PreviousCandidateDigest != b.DirectionCandidateDigest ||
		b.PreviousCandidateDigest != b.ParentCandidateDigest {
		return fmt.Errorf("Gooo extended lineage next candidate parent mismatch")
	}
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("Gooo extended lineage next candidate revision must be non-executing and non-authorizing")
	}
	expectedBoundChangeDigest := digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionNextChange(
		b.SourceRevisionChangeDigest,
		b.PreviousCandidateDigest,
		b.PreviousCandidateEvidenceDigest,
		b.InputEvidenceDigest,
		b.DirectionEvidenceDigest,
		b.EvidencePrefixDigest,
		b.MetricEvidenceDigest,
		b.FeedbackDigest,
		b.AggregationEvidenceDigest,
	)
	if b.BoundRevisionChangeDigest != expectedBoundChangeDigest {
		return fmt.Errorf("Gooo extended lineage next candidate revision change digest mismatch")
	}
	expectedCandidateDigest := digestJEVImprovementRevisionCandidate(
		b.ParentCandidateDigest,
		b.RevisionSource,
		b.BoundRevisionChangeDigest,
		b.DirectiveEvidenceDigest,
	)
	if b.CandidateDigest != expectedCandidateDigest {
		return fmt.Errorf("Gooo extended lineage next candidate digest mismatch")
	}
	expectedCandidateEvidenceDigest := digestJEVImprovementRevisionCandidateEvidence(
		b.CandidateStatus,
		b.CandidateDigest,
		b.DirectiveEvidenceDigest,
		b.BoundRevisionChangeDigest,
	)
	if b.CandidateEvidenceDigest != expectedCandidateEvidenceDigest {
		return fmt.Errorf("Gooo extended lineage next candidate evidence digest mismatch")
	}
	expectedEvidenceDigest := digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionNext(b)
	if b.EvidenceDigest != expectedEvidenceDigest {
		return fmt.Errorf("Gooo extended lineage next candidate evidence digest mismatch")
	}
	return nil
}

// GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionNext
// generates the next candidate only; it never executes or authorizes the proposed revision.
func GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionNext(input ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionNextInput) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionNextBinding {
	unknown := func(stage string) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionNextBinding {
		if strings.TrimSpace(stage) == "" {
			stage = "revision-action-candidate-generation-extended-lineage-candidate-revision-next"
		}
		return ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionNextBinding{
			Status:         "UNKNOWN",
			MissingStage:   stage,
			NonExecuting:   true,
			NonAuthorizing: true,
		}
	}
	if !input.NonAuthorizing || !input.Direction.NonAuthorizing {
		return unknown("authorization-boundary")
	}
	if !input.Direction.NonExecuting {
		return unknown("execution-boundary")
	}
	if err := input.Direction.Validate(); err != nil {
		return unknown("revision-action-candidate-generation-extended-lineage-candidate-revision-direction-validation")
	}
	if strings.TrimSpace(input.PreviousCandidateDigest) == "" ||
		input.PreviousCandidateDigest != input.Direction.DirectionCandidateDigest {
		return unknown("previous-candidate")
	}
	if strings.TrimSpace(input.PreviousCandidateEvidenceDigest) == "" ||
		input.PreviousCandidateEvidenceDigest != input.Direction.CandidateEvidenceDigest {
		return unknown("previous-candidate-evidence")
	}
	if strings.TrimSpace(input.CandidateSource) == "" {
		return unknown("candidate-source")
	}
	if strings.TrimSpace(input.RevisionSource) == "" {
		return unknown("revision-source")
	}
	if strings.TrimSpace(input.RevisionChangeDigest) == "" {
		return unknown("revision-change")
	}
	boundChangeDigest := digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionNextChange(
		input.RevisionChangeDigest,
		input.PreviousCandidateDigest,
		input.PreviousCandidateEvidenceDigest,
		input.Direction.InputEvidenceDigest,
		input.Direction.DirectionEvidenceDigest,
		input.Direction.EvidencePrefixDigest,
		input.Direction.MetricEvidenceDigest,
		input.Direction.FeedbackDigest,
		input.Direction.AggregationEvidenceDigest,
	)
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
	output := ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionNextBinding{
		Status:                       "bound",
		CandidateStatus:              candidate.Status,
		PreviousCandidateDigest:      input.PreviousCandidateDigest,
		PreviousCandidateEvidenceDigest: input.PreviousCandidateEvidenceDigest,
		ParentCandidateDigest:        candidate.ParentCandidateDigest,
		CandidateSource:              input.CandidateSource,
		RevisionSource:               candidate.RevisionSource,
		SourceRevisionChangeDigest:   input.RevisionChangeDigest,
		BoundRevisionChangeDigest:    candidate.RevisionChangeDigest,
		DirectiveEvidenceDigest:      candidate.DirectiveEvidenceDigest,
		CandidateDigest:              candidate.CandidateDigest,
		CandidateEvidenceDigest:      candidate.EvidenceDigest,
		DirectionCandidateDigest:     input.Direction.DirectionCandidateDigest,
		InputEvidenceDigest:          input.Direction.InputEvidenceDigest,
		DirectionEvidenceDigest:      input.Direction.DirectionEvidenceDigest,
		GuardEvidenceDigest:          input.Direction.GuardEvidenceDigest,
		VerificationEvidenceDigest:  input.Direction.VerificationEvidenceDigest,
		ReverseEvidenceDigest:        input.Direction.ReverseEvidenceDigest,
		EvidencePrefixDigest:         input.Direction.EvidencePrefixDigest,
		MetricEvidenceDigest:         input.Direction.MetricEvidenceDigest,
		FeedbackDigest:               input.Direction.FeedbackDigest,
		AggregationEvidenceDigest:    input.Direction.AggregationEvidenceDigest,
		NonExecuting:                 true,
		NonAuthorizing:               true,
	}
	output.EvidenceDigest = digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionNext(output)
	if err := output.Validate(); err != nil {
		return unknown("revision-action-candidate-generation-extended-lineage-candidate-revision-next-evidence")
	}
	return output
}

func digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionNextChange(
	revisionChangeDigest,
	previousCandidateDigest,
	previousCandidateEvidenceDigest,
	inputEvidenceDigest,
	directionEvidenceDigest,
	evidencePrefixDigest,
	metricEvidenceDigest,
	feedbackDigest,
	aggregationEvidenceDigest string,
) string {
	digest, err := Digest(struct {
		RevisionChangeDigest         string
		PreviousCandidateDigest      string
		PreviousCandidateEvidenceDigest string
		InputEvidenceDigest          string
		DirectionEvidenceDigest      string
		EvidencePrefixDigest         string
		MetricEvidenceDigest         string
		FeedbackDigest               string
		AggregationEvidenceDigest    string
	}{
		RevisionChangeDigest:            revisionChangeDigest,
		PreviousCandidateDigest:         previousCandidateDigest,
		PreviousCandidateEvidenceDigest: previousCandidateEvidenceDigest,
		InputEvidenceDigest:             inputEvidenceDigest,
		DirectionEvidenceDigest:         directionEvidenceDigest,
		EvidencePrefixDigest:            evidencePrefixDigest,
		MetricEvidenceDigest:            metricEvidenceDigest,
		FeedbackDigest:                  feedbackDigest,
		AggregationEvidenceDigest:       aggregationEvidenceDigest,
	})
	if err != nil {
		return ""
	}
	return digest
}

func digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionNext(b ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionNextBinding) string {
	digest, err := Digest(struct {
		Status                       string
		CandidateStatus              string
		PreviousCandidateDigest      string
		PreviousCandidateEvidenceDigest string
		ParentCandidateDigest        string
		CandidateSource              string
		RevisionSource               string
		SourceRevisionChangeDigest   string
		BoundRevisionChangeDigest    string
		DirectiveEvidenceDigest      string
		CandidateDigest              string
		CandidateEvidenceDigest      string
		DirectionCandidateDigest     string
		InputEvidenceDigest          string
		DirectionEvidenceDigest      string
		GuardEvidenceDigest          string
		VerificationEvidenceDigest  string
		ReverseEvidenceDigest        string
		EvidencePrefixDigest         string
		MetricEvidenceDigest         string
		FeedbackDigest               string
		AggregationEvidenceDigest    string
		NonExecuting                 bool
		NonAuthorizing               bool
	}{
		Status:                       b.Status,
		CandidateStatus:              b.CandidateStatus,
		PreviousCandidateDigest:      b.PreviousCandidateDigest,
		PreviousCandidateEvidenceDigest: b.PreviousCandidateEvidenceDigest,
		ParentCandidateDigest:        b.ParentCandidateDigest,
		CandidateSource:              b.CandidateSource,
		RevisionSource:               b.RevisionSource,
		SourceRevisionChangeDigest:   b.SourceRevisionChangeDigest,
		BoundRevisionChangeDigest:    b.BoundRevisionChangeDigest,
		DirectiveEvidenceDigest:      b.DirectiveEvidenceDigest,
		CandidateDigest:              b.CandidateDigest,
		CandidateEvidenceDigest:      b.CandidateEvidenceDigest,
		DirectionCandidateDigest:     b.DirectionCandidateDigest,
		InputEvidenceDigest:          b.InputEvidenceDigest,
		DirectionEvidenceDigest:      b.DirectionEvidenceDigest,
		GuardEvidenceDigest:           b.GuardEvidenceDigest,
		VerificationEvidenceDigest:   b.VerificationEvidenceDigest,
		ReverseEvidenceDigest:        b.ReverseEvidenceDigest,
		EvidencePrefixDigest:         b.EvidencePrefixDigest,
		MetricEvidenceDigest:         b.MetricEvidenceDigest,
		FeedbackDigest:               b.FeedbackDigest,
		AggregationEvidenceDigest:    b.AggregationEvidenceDigest,
		NonExecuting:                 b.NonExecuting,
		NonAuthorizing:               b.NonAuthorizing,
	})
	if err != nil {
		return ""
	}
	return digest
}