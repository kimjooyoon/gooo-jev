package decision

import "testing"

func materializationForAdmissionTest() JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridge {
    materialization := JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridge{
        Status:                      jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridgeBound,
        CandidateDigest:              "candidate-digest",
        CandidateSource:              "candidate-source",
        CandidateGateDigest:          "candidate-gate-digest",
        FeedbackDirection:             jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackImprove,
        FeedbackDigest:                "feedback-digest",
        FeedbackSource:                "feedback-engine",
        FeedbackEvidenceDigest:         "feedback-evidence",
        FeedbackBridgeDigest:          "feedback-bridge-digest",
        ProposalDecision:              jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalRevise,
        ProposalTarget:                "revision-target",
        ProposalDigest:                "proposal-digest",
        ProposalSource:                "proposal-generator",
        ProposalEvidenceDigest:         "proposal-evidence",
        RevisionProposalBridgeDigest:  "revision-proposal-bridge-digest",
        MaterializationStatus:          jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterialized,
        GeneratedCandidateDigest:       "generated-candidate-digest",
        GeneratedCandidateSource:       "candidate-generator",
        GenerationInputDigest:          "generation-input",
        GenerationEvidenceDigest:       "generation-evidence",
        GeneratorIdentity:              "generator-v1",
        NonExecuting:                   true,
        NonAuthorizing:                 true,
    }
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
    return materialization
}

func TestBindJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeReady(t *testing.T) {
    bridge := BindJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridge(JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeInput{
        Materialization:             materializationForAdmissionTest(),
        AdmissionDecision:           jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionAdmit,
        AdmissionDigest:             "admission-digest",
        AdmissionSource:             "admission-review",
        AdmissionEvidenceDigest:     "admission-evidence",
        QualityMetricDigest:         "quality-metric-digest",
        QualityMetricSource:         "quality-metric",
        QualityMetricEvidenceDigest: "quality-metric-evidence",
        ReviewDigest:                "review-digest",
        ReviewSource:                "review-engine",
        ReviewEvidenceDigest:        "review-evidence",
        NonAuthorizing:               true,
    })
    if err := bridge.Validate(); err != nil {
        t.Fatalf("expected valid candidate admission bridge, got %v", err)
    }
    if bridge.AdmissionStatus != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBound ||
        bridge.AdmissionDecision != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionAdmit ||
        bridge.QualityMetricDigest != "quality-metric-digest" {
        t.Fatalf("candidate admission provenance was not preserved: %#v", bridge)
    }
}

func TestBindJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeRequiresReviewEvidence(t *testing.T) {
    bridge := BindJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridge(JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeInput{
        Materialization:             materializationForAdmissionTest(),
        AdmissionDecision:           jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionAdmit,
        AdmissionDigest:             "admission-digest",
        AdmissionSource:             "admission-review",
        AdmissionEvidenceDigest:     "admission-evidence",
        QualityMetricDigest:         "quality-metric-digest",
        QualityMetricSource:         "quality-metric",
        QualityMetricEvidenceDigest: "quality-metric-evidence",
        ReviewSource:                "review-engine",
        ReviewEvidenceDigest:        "review-evidence",
        NonAuthorizing:               true,
    })
    if bridge.Status != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeUnknown ||
        bridge.MissingStage != "review-digest" {
        t.Fatalf("missing review digest was not preserved: %#v", bridge)
    }
    if bridge.BridgeDigest != "" {
        t.Fatalf("incomplete candidate admission bridge unexpectedly has digest: %#v", bridge)
    }
}
