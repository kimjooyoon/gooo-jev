package decision

import "testing"

func TestProjectJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalLSPReady(t *testing.T) {
    proposal := JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionProposal{
        Status:                    jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalReady,
        BridgeStatus:              jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateReady,
        CandidateGenerationDigest: "generation-digest",
        GateDigest:                "gate-digest",
        CandidateDigest:           "candidate-digest",
        RevisionSource:            "candidate-source",
        RevisionPlanSource:        "revision-plan-source",
        RevisionIntentDigest:      "revision-intent-digest",
        NonExecuting:              true,
        NonAuthorizing:            true,
    }
    proposal.ProposalDigest = digestJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionProposal(
        proposal.Status,
        proposal.BridgeStatus,
        proposal.CandidateGenerationDigest,
        proposal.GateDigest,
        proposal.CandidateDigest,
        proposal.RevisionSource,
        proposal.RevisionPlanSource,
        proposal.RevisionIntentDigest,
    )

    diagnostic := ProjectJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalLSP(proposal)
    if err := diagnostic.Validate(); err != nil {
        t.Fatalf("expected valid ready diagnostic, got %v", err)
    }
    if !diagnostic.Publishable || diagnostic.RevisionPlanSource != "revision-plan-source" ||
        diagnostic.RevisionIntentDigest != "revision-intent-digest" {
        t.Fatalf("ready diagnostic did not preserve proposal evidence: %#v", diagnostic)
    }
}

func TestProjectJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalLSPUnknown(t *testing.T) {
    proposal := JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionProposal{
        Status:         jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalReady,
        BridgeStatus:   jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateReady,
        NonExecuting:   true,
        NonAuthorizing: true,
    }

    diagnostic := ProjectJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalLSP(proposal)
    if err := diagnostic.Validate(); err != nil {
        t.Fatalf("expected valid UNKNOWN diagnostic, got %v", err)
    }
    if diagnostic.Publishable || diagnostic.Status != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalUnknown ||
        diagnostic.MissingStage != "revision-proposal-evidence" {
        t.Fatalf("incomplete proposal was not preserved as UNKNOWN: %#v", diagnostic)
    }
}
