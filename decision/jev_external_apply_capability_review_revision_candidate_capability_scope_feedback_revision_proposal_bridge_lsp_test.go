package decision

import "testing"

func TestProjectJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridgeLSP(t *testing.T) {
    diagnostic := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridgeLSP(
        JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridge{
            Status:                 jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridgeBound,
            CandidateDigest:        "candidate-digest",
            CandidateSource:        "candidate-source",
            CandidateGateDigest:    "candidate-gate-digest",
            FeedbackStatus:         jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBound,
            FeedbackDirection:      jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackImprove,
            FeedbackDigest:         "feedback-digest",
            FeedbackSource:         "feedback-engine",
            FeedbackEvidenceDigest: "feedback-evidence",
            FeedbackBridgeDigest:   "feedback-bridge-digest",
            RevisionProposalStatus: jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBound,
            ProposalDecision:       jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalRevise,
            ProposalTarget:          "revision-target",
            ProposalDigest:          "proposal-digest",
            ProposalSource:          "proposal-generator",
            ProposalEvidenceDigest: "proposal-evidence",
            BridgeDigest:            "revision-bridge-digest",
            NonExecuting:            true,
            NonAuthorizing:          true,
        },
    )
    if err := diagnostic.Validate(); err != nil {
        t.Fatalf("expected valid revision proposal LSP projection, got %v", err)
    }
    if !diagnostic.Publishable || diagnostic.Code != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridgeLSPBoundCode ||
        diagnostic.ProposalDecision != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalRevise {
        t.Fatalf("revision proposal evidence was not published: %#v", diagnostic)
    }
}

func TestProjectJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridgeLSPUnknown(t *testing.T) {
    diagnostic := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridgeLSP(
        JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridge{
            Status:        jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridgeUnknown,
            MissingStage:  "proposal-evidence",
            NonExecuting:  true,
            NonAuthorizing: true,
        },
    )
    if err := diagnostic.Validate(); err != nil {
        t.Fatalf("expected valid UNKNOWN revision proposal projection, got %v", err)
    }
    if diagnostic.Publishable || diagnostic.MissingStage != "proposal-evidence" {
        t.Fatalf("UNKNOWN revision proposal projection lost missing stage: %#v", diagnostic)
    }
}
