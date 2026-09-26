package decision

import "testing"

func feedbackForNextCandidateBridge(t *testing.T, choice string) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationFeedbackBinding {
	t.Helper()
	return BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationFeedbackInput{
		Metric:                metricForObservationFeedback(t),
		FeedbackChoice:        choice,
		FeedbackEvidenceDigest: "feedback-evidence-next-bridge-" + choice,
		NonAuthorizing:        true,
	})
}

func TestDeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanFeedbackNextRefutedGeneratesCandidate(t *testing.T) {
	output := DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanFeedbackNext(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanFeedbackNextInput{
		Feedback:            feedbackForNextCandidateBridge(t, "refuted"),
		CandidateSource:     "gooo://candidate/execution-plan-feedback-next",
		RevisionSource:      "gooo://revision/execution-plan-feedback-next",
		RevisionChangeDigest: "change-from-refuted-plan-feedback",
		NonAuthorizing:      true,
	})
	if output.Status != "bound" || output.NextAction != "generate-revision-candidate" ||
		output.Directive != jevImprovementDirectiveRevision ||
		output.RevisionCandidateStatus != jevImprovementRevisionCandidateReady ||
		output.RevisionCandidateDigest == "" {
		t.Fatalf("output = %#v, want generated revision candidate", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("output should validate: %v", err)
	}
}

func TestDeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanFeedbackNextConfirmedRequestsReview(t *testing.T) {
	output := DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanFeedbackNext(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanFeedbackNextInput{
		Feedback:        feedbackForNextCandidateBridge(t, "confirmed"),
		CandidateSource:  "gooo://candidate/execution-plan-feedback-next-review",
		NonAuthorizing:   true,
	})
	if output.Status != "bound" || output.NextAction != "external-review-required" ||
		output.Directive != jevImprovementDirectiveExternalReview ||
		output.RevisionCandidateStatus != "" {
		t.Fatalf("output = %#v, want external review action", output)
	}
}

func TestDeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanFeedbackNextUnknownHolds(t *testing.T) {
	output := DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanFeedbackNext(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanFeedbackNextInput{
		Feedback:        feedbackForNextCandidateBridge(t, "unknown"),
		CandidateSource:  "gooo://candidate/execution-plan-feedback-next-hold",
		NonAuthorizing:   true,
	})
	if output.Status != "bound" || output.NextAction != "evidence-hold" ||
		output.Directive != jevImprovementDirectiveEvidence {
		t.Fatalf("output = %#v, want evidence hold action", output)
	}
}

func TestDeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanFeedbackNextRequiresRevisionSource(t *testing.T) {
	output := DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanFeedbackNext(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanFeedbackNextInput{
		Feedback:            feedbackForNextCandidateBridge(t, "refuted"),
		CandidateSource:     "gooo://candidate/execution-plan-feedback-next-missing-source",
		RevisionChangeDigest: "change-missing-revision-source",
		NonAuthorizing:      true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-source" {
		t.Fatalf("output = %#v, want revision-source UNKNOWN", output)
	}
}

func TestDeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanFeedbackNextRejectsFeedbackTampering(t *testing.T) {
	feedback := feedbackForNextCandidateBridge(t, "refuted")
	feedback.FeedbackEvidenceDigest = "tampered"
	output := DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanFeedbackNext(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanFeedbackNextInput{
		Feedback:            feedback,
		CandidateSource:     "gooo://candidate/execution-plan-feedback-next-tampered",
		RevisionSource:      "gooo://revision/execution-plan-feedback-next-tampered",
		RevisionChangeDigest: "change-tampered-feedback",
		NonAuthorizing:      true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "execution-plan-observation-feedback-validation" {
		t.Fatalf("output = %#v, want feedback validation UNKNOWN", output)
	}
}