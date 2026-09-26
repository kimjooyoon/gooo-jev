package decision

import "testing"

func evidenceFullProvenanceMaterializationEvaluationInput(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceMaterializationAdmissionEvaluationInput {
	t.Helper()
	revision := GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionCandidate(evidenceFullProvenanceRevisionCandidateInput(t))
	materialization := materializationForAdmissionTest()
	materialization.CandidateDigest = revision.CandidateDigest
	materialization.BridgeDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridge(
		materialization.Status,
		materialization.CandidateDigest,
		materialization.CandidateSource,
		materialization.CandidateGateDigest,
		materialization.FeedbackDirection,
		materialization.FeedbackDigest,
		materialization.FeedbackSource,
		materialization.FeedbackEvidenceDigest,
		materialization.FeedbackBridgeDigest,
		materialization.ProposalDecision,
		materialization.ProposalTarget,
		materialization.ProposalDigest,
		materialization.ProposalSource,
		materialization.ProposalEvidenceDigest,
		materialization.RevisionProposalBridgeDigest,
		materialization.MaterializationStatus,
		materialization.GeneratedCandidateDigest,
		materialization.GeneratedCandidateSource,
		materialization.GenerationInputDigest,
		materialization.GenerationEvidenceDigest,
		materialization.GeneratorIdentity,
	)
	admission := BindJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridge(JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeInput{
		Materialization:             materialization,
		AdmissionDecision:           jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionAdmit,
		AdmissionDigest:             "admission-digest",
		AdmissionSource:              "admission-review",
		AdmissionEvidenceDigest:     "admission-evidence",
		QualityMetricDigest:          "quality-metric-digest",
		QualityMetricSource:          "quality-metric",
		QualityMetricEvidenceDigest: "quality-metric-evidence",
		ReviewDigest:                "review-digest",
		ReviewSource:                "review-engine",
		ReviewEvidenceDigest:        "review-evidence",
		NonAuthorizing:               true,
	})
	return ExecutionEnvelopeGoooEvidenceFullProvenanceMaterializationAdmissionEvaluationInput{
		RevisionCandidate:                  revision,
		Admission:                          admission,
		EvaluationMode:                     "reverse-observation",
		ReverseObservationDigest:           "evaluation-reverse-observation",
		ReverseObservationSource:           "evaluation-engine",
		ReverseObservationEvidenceDigest:   "evaluation-observation-evidence",
		EvaluationMetricDigest:             "evaluation-metric",
		EvaluationMetricSource:             "evaluation-metric-source",
		EvaluationMetricEvidenceDigest:     "evaluation-metric-evidence",
		NonAuthorizing:                     true,
	}
}

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceMaterializationAdmissionEvaluation(t *testing.T) {
	binding := BindExecutionEnvelopeGoooEvidenceFullProvenanceMaterializationAdmissionEvaluation(evidenceFullProvenanceMaterializationEvaluationInput(t))
	if binding.Status != "bound" || binding.AdmissionDecision != "admit" ||
		binding.EvaluationStatus != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeEvaluationBound ||
		binding.EvidenceDigest == "" {
		t.Fatalf("binding = %#v, want bound admitted evaluation", binding)
	}
	if err := binding.Validate(); err != nil {
		t.Fatalf("binding should validate: %v", err)
	}
}

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceMaterializationAdmissionEvaluationRejectsCandidateMismatch(t *testing.T) {
	input := evidenceFullProvenanceMaterializationEvaluationInput(t)
	input.RevisionCandidate.CandidateDigest = "tampered"
	binding := BindExecutionEnvelopeGoooEvidenceFullProvenanceMaterializationAdmissionEvaluation(input)
	if binding.Status != "UNKNOWN" || binding.MissingStage != "candidate-digest-binding" {
		t.Fatalf("binding = %#v, want candidate-digest-binding UNKNOWN", binding)
	}
}

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceMaterializationAdmissionEvaluationPreservesEvaluationFailure(t *testing.T) {
	input := evidenceFullProvenanceMaterializationEvaluationInput(t)
	input.ReverseObservationSource = ""
	binding := BindExecutionEnvelopeGoooEvidenceFullProvenanceMaterializationAdmissionEvaluation(input)
	if binding.Status != "UNKNOWN" || binding.MissingStage != "reverse-observation-source" {
		t.Fatalf("binding = %#v, want reverse-observation-source UNKNOWN", binding)
	}
}

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceMaterializationAdmissionEvaluationRejectsAuthorization(t *testing.T) {
	binding := BindExecutionEnvelopeGoooEvidenceFullProvenanceMaterializationAdmissionEvaluation(ExecutionEnvelopeGoooEvidenceFullProvenanceMaterializationAdmissionEvaluationInput{
		NonAuthorizing: false,
	})
	if binding.Status != "UNKNOWN" || binding.MissingStage != "authorization-boundary" {
		t.Fatalf("binding = %#v, want authorization UNKNOWN", binding)
	}
}
