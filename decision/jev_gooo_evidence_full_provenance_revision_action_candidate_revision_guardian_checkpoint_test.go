package decision

import "testing"

func lowConfidenceForGuardianCheckpoint(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionConfidenceBinding {
	t.Helper()
	return CalibrateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionConfidence(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionConfidenceInput{
		Review:                    supportRevisionCandidateReviewLSPForConfidence(t),
		ConfidenceBand:            "low",
		CalibrationEvidenceDigest: "calibration-evidence-guardian-low",
		NonAuthorizing:            true,
	})
}

func abstainConfidenceForGuardianCheckpoint(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionConfidenceBinding {
	t.Helper()
	return CalibrateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionConfidence(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionConfidenceInput{
		Review:                    abstainRevisionCandidateReviewLSPForConfidence(t),
		ConfidenceBand:            "high",
		CalibrationEvidenceDigest: "calibration-evidence-guardian-abstain",
		NonAuthorizing:            true,
	})
}

func TestEvaluateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGuardianCheckpointKeepsPermissionDenied(t *testing.T) {
	output := EvaluateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGuardianCheckpoint(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGuardianCheckpointInput{
		Confidence:                   lowConfidenceForGuardianCheckpoint(t),
		RequestedCapability:          "revision-candidate-apply",
		GuardianPolicyEvidenceDigest: "guardian-policy-evidence-low",
		NonAuthorizing:               true,
	})
	if output.Status != "bound" || output.PermissionState != "not-authorized" ||
		output.GuardianDecision != "shadow-observation-only" {
		t.Fatalf("output = %#v, want non-authorizing guardian checkpoint", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("output should validate: %v", err)
	}
}

func TestEvaluateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGuardianCheckpointRoutesAbstain(t *testing.T) {
	output := EvaluateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGuardianCheckpoint(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGuardianCheckpointInput{
		Confidence:                   abstainConfidenceForGuardianCheckpoint(t),
		RequestedCapability:          "revision-candidate-apply",
		GuardianPolicyEvidenceDigest: "guardian-policy-evidence-abstain",
		NonAuthorizing:               true,
	})
	if output.Status != "bound" || output.PermissionState != "not-authorized" ||
		output.GuardianDecision != "human-or-deterministic-policy" {
		t.Fatalf("output = %#v, want human fallback checkpoint", output)
	}
}

func TestEvaluateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGuardianCheckpointRequiresCapability(t *testing.T) {
	output := EvaluateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGuardianCheckpoint(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGuardianCheckpointInput{
		Confidence:                   lowConfidenceForGuardianCheckpoint(t),
		GuardianPolicyEvidenceDigest: "guardian-policy-evidence-missing-capability",
		NonAuthorizing:               true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "requested-capability" {
		t.Fatalf("output = %#v, want requested-capability UNKNOWN", output)
	}
}

func TestEvaluateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGuardianCheckpointRejectsConfidenceTampering(t *testing.T) {
	confidence := lowConfidenceForGuardianCheckpoint(t)
	confidence.AutomationMode = "tampered"
	output := EvaluateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGuardianCheckpoint(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGuardianCheckpointInput{
		Confidence:                   confidence,
		RequestedCapability:          "revision-candidate-apply",
		GuardianPolicyEvidenceDigest: "guardian-policy-evidence-tampered",
		NonAuthorizing:               true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-action-candidate-generation-extended-lineage-candidate-revision-confidence-validation" {
		t.Fatalf("output = %#v, want confidence validation UNKNOWN", output)
	}
}