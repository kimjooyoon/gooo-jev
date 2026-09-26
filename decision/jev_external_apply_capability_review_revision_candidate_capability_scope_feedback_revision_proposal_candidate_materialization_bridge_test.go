package decision

import "testing"

func revisionProposalForMaterializationTest() JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridge {
    proposal := JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridge{
        Status:                 jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridgeBound,
        CandidateDigest:        "candidate-digest",
        CandidateSource:        "candidate-source",
        CandidateGateDigest:    "candidate-gate-digest",
        FeedbackStatus:         jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBound,
        FeedbackDirection:       jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackImprove,
        FeedbackDigest:          "feedback-digest",
        FeedbackSource:          "feedback-engine",
        FeedbackEvidenceDigest:  "feedback-evidence",
        FeedbackBridgeDigest:    "feedback-bridge-digest",
        RevisionProposalStatus:  jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBound,
        ProposalDecision:        jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalRevise,
        ProposalTarget:           "revision-target",
        ProposalDigest:           "proposal-digest",
        ProposalSource:           "proposal-generator",
        ProposalEvidenceDigest:   "proposal-evidence",
        NonExecuting:             true,
        NonAuthorizing:           true,
    }
    proposal.BridgeDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridge(
        proposal.Status,
        proposal.CandidateDigest,
        proposal.CandidateSource,
        proposal.CandidateGateDigest,
        proposal.FeedbackStatus,
        proposal.FeedbackDirection,
        proposal.FeedbackDigest,
        proposal.FeedbackSource,
        proposal.FeedbackEvidenceDigest,
        proposal.FeedbackBridgeDigest,
        proposal.RevisionProposalStatus,
        proposal.ProposalDecision,
        proposal.ProposalTarget,
        proposal.ProposalDigest,
        proposal.ProposalSource,
        proposal.ProposalEvidenceDigest,
    )
    return proposal
}

func TestBindJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridgeReady(t *testing.T) {
    bridge := BindJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridge(JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridgeInput{
        RevisionProposal:         revisionProposalForMaterializationTest(),
        GeneratedCandidateDigest: "generated-candidate-digest",
        GeneratedCandidateSource: "candidate-generator",
        GenerationInputDigest:    "generation-input",
        GenerationEvidenceDigest: "generation-evidence",
        GeneratorIdentity:        "generator-v1",
        NonAuthorizing:            true,
    })
    if err := bridge.Validate(); err != nil {
        t.Fatalf("expected valid candidate materialization bridge, got %v", err)
    }
    if bridge.MaterializationStatus != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterialized ||
        bridge.GeneratedCandidateDigest != "generated-candidate-digest" ||
        bridge.ProposalDecision != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalRevise {
        t.Fatalf("candidate materialization provenance was not preserved: %#v", bridge)
    }
}

func TestBindJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridgeRequiresGenerationEvidence(t *testing.T) {
    bridge := BindJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridge(JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridgeInput{
        RevisionProposal:        revisionProposalForMaterializationTest(),
        GeneratedCandidateSource: "candidate-generator",
        GenerationInputDigest:   "generation-input",
        GenerationEvidenceDigest: "generation-evidence",
        GeneratorIdentity:       "generator-v1",
        NonAuthorizing:           true,
    })
    if bridge.Status != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridgeUnknown ||
        bridge.MissingStage != "generated-candidate-digest" {
        t.Fatalf("missing generated candidate digest was not preserved: %#v", bridge)
    }
    if bridge.BridgeDigest != "" {
        t.Fatalf("incomplete candidate materialization bridge unexpectedly has digest: %#v", bridge)
    }
}
