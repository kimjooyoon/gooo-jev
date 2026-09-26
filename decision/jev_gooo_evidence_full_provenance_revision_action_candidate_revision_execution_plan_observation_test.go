package decision

import "testing"

func executionPlanForObservation(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanBinding {
	t.Helper()
	return PlanExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlan(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanInput{
		Guardian:                  guardianCheckpointForExecutionPlan(t),
		WorkspaceSource:           "gooo://workspace/revision-candidate-observation",
		Gateway:                   "gateway://deterministic-policy",
		Model:                     "jev://review",
		NetworkAllowlistDigest:    "network-allowlist-observation",
		SuspendTokenDigest:        "suspend-token-observation",
		NonAuthorizing:            true,
	})
}

func TestObserveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationSuspended(t *testing.T) {
	output := ObserveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservation(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationInput{
		Plan:                     executionPlanForObservation(t),
		LifecycleEvent:           "suspended",
		ObservationEvidenceDigest: "observation-evidence-suspended",
		ReverseObservationDigest: "reverse-observation-suspended",
		NonAuthorizing:            true,
	})
	if output.Status != "bound" || output.ObservationStatus != "suspended" ||
		output.LifecycleDisposition != "resume-required" ||
		output.PermissionState != "not-authorized" {
		t.Fatalf("output = %#v, want suspended reverse observation", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("output should validate: %v", err)
	}
}

func TestObserveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationResumed(t *testing.T) {
	output := ObserveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservation(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationInput{
		Plan:                     executionPlanForObservation(t),
		LifecycleEvent:           "resumed",
		ObservationEvidenceDigest: "observation-evidence-resumed",
		ReverseObservationDigest: "reverse-observation-resumed",
		NonAuthorizing:            true,
	})
	if output.Status != "bound" || output.LifecycleDisposition != "resume-observed" {
		t.Fatalf("output = %#v, want resumed reverse observation", output)
	}
}

func TestObserveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationRejectsUnknownEvent(t *testing.T) {
	output := ObserveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservation(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationInput{
		Plan:                     executionPlanForObservation(t),
		LifecycleEvent:           "executed",
		ObservationEvidenceDigest: "observation-evidence-invalid",
		ReverseObservationDigest: "reverse-observation-invalid",
		NonAuthorizing:            true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "lifecycle-event" {
		t.Fatalf("output = %#v, want lifecycle-event UNKNOWN", output)
	}
}

func TestObserveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationRejectsPlanTampering(t *testing.T) {
	plan := executionPlanForObservation(t)
	plan.Model = "tampered-model"
	output := ObserveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservation(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationInput{
		Plan:                     plan,
		LifecycleEvent:           "not-started",
		ObservationEvidenceDigest: "observation-evidence-tampered",
		ReverseObservationDigest: "reverse-observation-tampered",
		NonAuthorizing:            true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-action-candidate-generation-extended-lineage-candidate-revision-execution-plan-validation" {
		t.Fatalf("output = %#v, want plan validation UNKNOWN", output)
	}
}