package decision

import "testing"

func TestProjectJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridgeLSP(t *testing.T) {
    diagnostic := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridgeLSP(
        JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridge{
            Status:                      jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridgeBound,
            CandidateDigest:              "candidate-digest",
            CandidateSource:              "candidate-source",
            CandidateGateDigest:          "candidate-gate-digest",
            FeedbackDirection:             jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackImprove,
            FeedbackDigest:               "feedback-digest",
            FeedbackSource:               "feedback-engine",
            FeedbackEvidenceDigest:        "feedback-evidence",
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
            BridgeDigest:                   "materialization-bridge-digest",
            NonExecuting:                   true,
            NonAuthorizing:                 true,
        },
    )
    if err := diagnostic.Validate(); err != nil {
        t.Fatalf("expected valid candidate materialization LSP projection, got %v", err)
    }
    if !diagnostic.Publishable || diagnostic.Code != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridgeLSPBoundCode ||
        diagnostic.GeneratedCandidateDigest != "generated-candidate-digest" {
        t.Fatalf("candidate materialization evidence was not published: %#v", diagnostic)
    }
}

func TestProjectJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridgeLSPUnknown(t *testing.T) {
    diagnostic := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridgeLSP(
        JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridge{
            Status:        jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridgeUnknown,
            MissingStage:  "generation-evidence",
            NonExecuting: true,
            NonAuthorizing: true,
        },
    )
    if err := diagnostic.Validate(); err != nil {
        t.Fatalf("expected valid UNKNOWN candidate materialization projection, got %v", err)
    }
    if diagnostic.Publishable || diagnostic.MissingStage != "generation-evidence" {
        t.Fatalf("UNKNOWN candidate materialization projection lost missing stage: %#v", diagnostic)
    }
}
