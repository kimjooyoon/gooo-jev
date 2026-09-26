package decision

import "testing"

func guardianCheckpointForExecutionPlan(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGuardianCheckpointBinding {
	t.Helper()
	return EvaluateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGuardianCheckpoint(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGuardianCheckpointInput{
		Confidence:                   lowConfidenceForGuardianCheckpoint(t),
		RequestedCapability:          "revision-candidate-apply",
		GuardianPolicyEvidenceDigest: "guardian-policy-evidence-plan",
		NonAuthorizing:               true,
	})
}

func TestPlanExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanPreservesBoundary(t *testing.T) {
	output := PlanExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlan(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanInput{
		Guardian:                  guardianCheckpointForExecutionPlan(t),
		WorkspaceSource:           "gooo://workspace/revision-candidate",
		Gateway:                   "gateway://deterministic-policy",
		Model:                     "jev://review",
		NetworkAllowlistDigest:    "network-allowlist-digest",
		SuspendTokenDigest:        "suspend-token-digest",
		NonAuthorizing:            true,
	})
	if output.Status != "bound" || output.PlanStatus != "planned" ||
		output.PermissionState != "not-authorized" || output.ResumeTokenDigest == "" {
		t.Fatalf("output = %#v, want non-authorizing execution plan", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("output should validate: %v", err)
	}
}

func TestPlanExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanRequiresWorkspace(t *testing.T) {
	output := PlanExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlan(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanInput{
		Guardian:               guardianCheckpointForExecutionPlan(t),
		Gateway:                "gateway://deterministic-policy",
		Model:                  "jev://review",
		NetworkAllowlistDigest: "network-allowlist-digest",
		SuspendTokenDigest:     "suspend-token-digest",
		NonAuthorizing:         true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "workspace-source" {
		t.Fatalf("output = %#v, want workspace-source UNKNOWN", output)
	}
}

func TestPlanExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanRejectsGuardianTampering(t *testing.T) {
	guardian := guardianCheckpointForExecutionPlan(t)
	guardian.GuardianPolicyEvidenceDigest = "tampered"
	output := PlanExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlan(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanInput{
		Guardian:                  guardian,
		WorkspaceSource:           "gooo://workspace/revision-candidate-tampered",
		Gateway:                   "gateway://deterministic-policy",
		Model:                     "jev://review",
		NetworkAllowlistDigest:    "network-allowlist-digest-tampered",
		SuspendTokenDigest:        "suspend-token-digest-tampered",
		NonAuthorizing:            true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-action-candidate-generation-extended-lineage-candidate-revision-guardian-validation" {
		t.Fatalf("output = %#v, want guardian validation UNKNOWN", output)
	}
}