package decision

import "testing"

func confirmedActionDerivedCandidateRevisionFeedback(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionFeedbackBinding {
	t.Helper()
	metric := confirmedActionDerivedCandidateRevisionMetric(t)
	return BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionFeedbackInput{
		Metric:                  metric,
		CandidateDigest:         metric.RevisionCandidateDigest,
		ReplayObservationDigest: "candidate-revision-direction-confirmed",
		NonAuthorizing:          true,
	})
}

func refutedActionDerivedCandidateRevisionFeedback(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionFeedbackBinding {
	t.Helper()
	metric := refutedActionDerivedCandidateRevisionMetric(t)
	return BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionFeedbackInput{
		Metric:                  metric,
		CandidateDigest:         metric.RevisionCandidateDigest,
		ReplayObservationDigest: "candidate-revision-direction-refuted",
		NonAuthorizing:          true,
	})
}

func unknownActionDerivedCandidateRevisionFeedback(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionFeedbackBinding {
	t.Helper()
	metric := unknownActionDerivedCandidateRevisionMetric(t)
	return BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionFeedbackInput{
		Metric:                  metric,
		CandidateDigest:         metric.RevisionCandidateDigest,
		ReplayObservationDigest: "candidate-revision-direction-unknown",
		NonAuthorizing:          true,
	})
}

func TestDeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionDirectionReview(t *testing.T) {
	output := DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionDirection(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionDirectionInput{
		Feedback:       confirmedActionDerivedCandidateRevisionFeedback(t),
		CandidateSource: "gooo://candidate/revision-confirmed",
		NonAuthorizing: true,
	})
	if output.Status != "bound" || output.Directive != jevImprovementDirectiveExternalReview ||
		output.SourceCandidateID == "" || output.RevisionCandidateDigest == "" {
		t.Fatalf("output = %#v, want external review direction", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("output should validate: %v", err)
	}
}

func TestDeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionDirectionRevision(t *testing.T) {
	output := DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionDirection(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionDirectionInput{
		Feedback:       refutedActionDerivedCandidateRevisionFeedback(t),
		CandidateSource: "gooo://candidate/revision-refuted",
		NonAuthorizing: true,
	})
	if output.Status != "bound" || output.Directive != jevImprovementDirectiveRevision ||
		output.DirectiveStatus != jevImprovementDirectiveRevision {
		t.Fatalf("output = %#v, want revision direction", output)
	}
}

func TestDeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionDirectionHold(t *testing.T) {
	output := DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionDirection(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionDirectionInput{
		Feedback:       unknownActionDerivedCandidateRevisionFeedback(t),
		CandidateSource: "gooo://candidate/revision-unknown",
		NonAuthorizing: true,
	})
	if output.Status != "bound" || output.Directive != jevImprovementDirectiveEvidence {
		t.Fatalf("output = %#v, want evidence hold direction", output)
	}
}

func TestDeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionDirectionRejectsTampering(t *testing.T) {
	feedback := confirmedActionDerivedCandidateRevisionFeedback(t)
	feedback.FeedbackDigest = "tampered"
	output := DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionDirection(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionDirectionInput{
		Feedback:       feedback,
		CandidateSource: "gooo://candidate/revision-tampered",
		NonAuthorizing: true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-action-candidate-revision-feedback-validation" {
		t.Fatalf("output = %#v, want feedback validation UNKNOWN", output)
	}
}

func TestDeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionDirectionRequiresSource(t *testing.T) {
	output := DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionDirection(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionDirectionInput{
		Feedback:       confirmedActionDerivedCandidateRevisionFeedback(t),
		NonAuthorizing: true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "candidate-source" {
		t.Fatalf("output = %#v, want candidate-source UNKNOWN", output)
	}
}