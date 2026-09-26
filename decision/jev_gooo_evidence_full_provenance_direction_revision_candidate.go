package decision

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

const (
	jevGoooEvidenceFullProvenanceDirectionRevisionCandidateBound   = "bound"
	jevGoooEvidenceFullProvenanceDirectionRevisionCandidateUnknown = "UNKNOWN"
)

// ExecutionEnvelopeGoooEvidenceFullProvenanceDirectionRevisionCandidateInput
// binds feedback and direction evidence before producing a revision candidate.
type ExecutionEnvelopeGoooEvidenceFullProvenanceDirectionRevisionCandidateInput struct {
	Feedback            ExecutionEnvelopeGoooEvidenceFullProvenanceEvaluationFeedbackBinding
	Direction           JEVImprovementDirectionDirective
	RevisionSource      string
	RevisionChangeDigest string
	NonAuthorizing      bool
}

// ExecutionEnvelopeGoooEvidenceFullProvenanceDirectionRevisionCandidateBinding
// preserves the complete evidence path from evaluation feedback to revision.
type ExecutionEnvelopeGoooEvidenceFullProvenanceDirectionRevisionCandidateBinding struct {
	Status                            string
	MissingStage                      string
	CandidateDigest                   string
	FeedbackAggregationEvidenceDigest string
	DirectionEvidenceDigest           string
	RevisionCandidateStatus           string
	RevisionSource                    string
	RevisionChangeDigest              string
	RevisionCandidateDigest            string
	RevisionCandidateEvidenceDigest   string
	EvidenceDigest                    string
	NonExecuting                      bool
	NonAuthorizing                    bool
}

func (b ExecutionEnvelopeGoooEvidenceFullProvenanceDirectionRevisionCandidateBinding) Validate() error {
	if b.Status != jevGoooEvidenceFullProvenanceDirectionRevisionCandidateBound ||
		b.MissingStage != "" ||
		b.CandidateDigest == "" ||
		b.FeedbackAggregationEvidenceDigest == "" ||
		b.DirectionEvidenceDigest == "" ||
		b.RevisionCandidateStatus != jevImprovementRevisionCandidateReady ||
		b.RevisionSource == "" ||
		b.RevisionChangeDigest == "" ||
		b.RevisionCandidateDigest == "" ||
		b.RevisionCandidateEvidenceDigest == "" ||
		b.EvidenceDigest == "" {
		return fmt.Errorf("incomplete Gooo direction revision candidate binding")
	}
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("Gooo direction revision candidate binding must be non-executing and non-authorizing")
	}
	expected := digestJEVGoooEvidenceFullProvenanceDirectionRevisionCandidateBinding(
		b.Status,
		b.CandidateDigest,
		b.FeedbackAggregationEvidenceDigest,
		b.DirectionEvidenceDigest,
		b.RevisionCandidateStatus,
		b.RevisionSource,
		b.RevisionChangeDigest,
		b.RevisionCandidateDigest,
		b.RevisionCandidateEvidenceDigest,
	)
	if b.EvidenceDigest != expected {
		return fmt.Errorf("Gooo direction revision candidate binding digest mismatch")
	}
	return nil
}

// BindExecutionEnvelopeGoooEvidenceFullProvenanceDirectionRevisionCandidate
// calls the existing revision generator only after cross-stage evidence agrees.
func BindExecutionEnvelopeGoooEvidenceFullProvenanceDirectionRevisionCandidate(input ExecutionEnvelopeGoooEvidenceFullProvenanceDirectionRevisionCandidateInput) ExecutionEnvelopeGoooEvidenceFullProvenanceDirectionRevisionCandidateBinding {
	unknown := func(stage string) ExecutionEnvelopeGoooEvidenceFullProvenanceDirectionRevisionCandidateBinding {
		if strings.TrimSpace(stage) == "" {
			stage = "direction-revision-candidate"
		}
		return ExecutionEnvelopeGoooEvidenceFullProvenanceDirectionRevisionCandidateBinding{
			Status:                  jevGoooEvidenceFullProvenanceDirectionRevisionCandidateUnknown,
			MissingStage:            stage,
			RevisionCandidateStatus: jevImprovementRevisionCandidateUnknown,
			NonExecuting:            true,
			NonAuthorizing:          true,
		}
	}
	if !input.NonAuthorizing || !input.Feedback.NonAuthorizing || !input.Direction.NonAuthorizing {
		return unknown("authorization-boundary")
	}
	if !input.Feedback.NonExecuting || !input.Direction.NonExecuting {
		return unknown("execution-boundary")
	}
	if err := input.Feedback.Validate(); err != nil {
		return unknown("evaluation-feedback-validation")
	}
	if err := input.Direction.Validate(); err != nil {
		return unknown("direction-validation")
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
	candidate := GenerateJEVImprovementRevisionCandidate(JEVImprovementRevisionCandidateInput{
		Directive:            input.Direction,
		RevisionSource:       input.RevisionSource,
		RevisionChangeDigest: input.RevisionChangeDigest,
		NonAuthorizing:       true,
	})
	if err := candidate.Validate(); err != nil {
		stage := candidate.MissingStage
		if strings.TrimSpace(stage) == "" {
			stage = "revision-candidate-validation"
		}
		return unknown(stage)
	}
	output := ExecutionEnvelopeGoooEvidenceFullProvenanceDirectionRevisionCandidateBinding{
		Status:                            jevGoooEvidenceFullProvenanceDirectionRevisionCandidateBound,
		CandidateDigest:                   input.Feedback.CandidateDigest,
		FeedbackAggregationEvidenceDigest: input.Feedback.AggregationEvidenceDigest,
		DirectionEvidenceDigest:           input.Direction.EvidenceDigest,
		RevisionCandidateStatus:            candidate.Status,
		RevisionSource:                     candidate.RevisionSource,
		RevisionChangeDigest:              candidate.RevisionChangeDigest,
		RevisionCandidateDigest:            candidate.CandidateDigest,
		RevisionCandidateEvidenceDigest:    candidate.EvidenceDigest,
		NonExecuting:                      true,
		NonAuthorizing:                    true,
	}
	output.EvidenceDigest = digestJEVGoooEvidenceFullProvenanceDirectionRevisionCandidateBinding(
		output.Status,
		output.CandidateDigest,
		output.FeedbackAggregationEvidenceDigest,
		output.DirectionEvidenceDigest,
		output.RevisionCandidateStatus,
		output.RevisionSource,
		output.RevisionChangeDigest,
		output.RevisionCandidateDigest,
		output.RevisionCandidateEvidenceDigest,
	)
	if err := output.Validate(); err != nil {
		return unknown("direction-revision-candidate-evidence")
	}
	return output
}

func digestJEVGoooEvidenceFullProvenanceDirectionRevisionCandidateBinding(
	status,
	candidateDigest,
	feedbackAggregationEvidenceDigest,
	directionEvidenceDigest,
	revisionCandidateStatus,
	revisionSource,
	revisionChangeDigest,
	revisionCandidateDigest,
	revisionCandidateEvidenceDigest string,
) string {
	sum := sha256.Sum256([]byte(fmt.Sprintf(
		"%s|%s|%s|%s|%s|%s|%s|%s|%s",
		status,
		candidateDigest,
		feedbackAggregationEvidenceDigest,
		directionEvidenceDigest,
		revisionCandidateStatus,
		revisionSource,
		revisionChangeDigest,
		revisionCandidateDigest,
		revisionCandidateEvidenceDigest,
	)))
	return hex.EncodeToString(sum[:])
}