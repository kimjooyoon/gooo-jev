package decision

import "testing"

func validExternalApplyCapabilityReviewImprovementDirectionForCandidateGate(direction string) JEVExternalApplyCapabilityReviewImprovementDirection {
    value := JEVExternalApplyCapabilityReviewImprovementDirection{
        Status:          jevExternalApplyCapabilityReviewImprovementDirectionBound,
        Direction:       direction,
        Target:          "decision/provenance",
        CandidateSource: "review-metric-trend",
        FeedbackDigest:  "feedback-digest",
        NonExecuting:    true,
        NonAuthorizing:  true,
    }
    value.DirectionDigest = digestJEVExternalApplyCapabilityReviewImprovementDirection(value.Status, value.Direction, value.Target, value.CandidateSource, value.FeedbackDigest)
    return value
}

func TestGateJEVExternalApplyCapabilityReviewRevisionCandidateGeneratesReadyCandidate(t *testing.T) {
    got := GateJEVExternalApplyCapabilityReviewRevisionCandidate(JEVExternalApplyCapabilityReviewRevisionCandidateGateInput{
        Direction:       validExternalApplyCapabilityReviewImprovementDirectionForCandidateGate(jevExternalApplyCapabilityReviewImprovementDirectionGenerate),
        CandidateDigest: "candidate-digest",
        RevisionSource: "revision-source",
        NonAuthorizing:  true,
    })
    if got.Decision != jevExternalApplyCapabilityReviewRevisionCandidateReady || got.CandidateDigest == "" || !got.NonExecuting || !got.NonAuthorizing {
        t.Fatalf("got %+v", got)
    }
    if err := got.Validate(); err != nil {
        t.Fatal(err)
    }
}

func TestGateJEVExternalApplyCapabilityReviewRevisionCandidatePreservesHoldAndReject(t *testing.T) {
    cases := []struct {
        direction string
        decision  string
    }{
        {jevExternalApplyCapabilityReviewImprovementDirectionHold, jevExternalApplyCapabilityReviewRevisionCandidateHold},
        {jevExternalApplyCapabilityReviewImprovementDirectionReject, jevExternalApplyCapabilityReviewRevisionCandidateRejected},
    }
    for _, testCase := range cases {
        got := GateJEVExternalApplyCapabilityReviewRevisionCandidate(JEVExternalApplyCapabilityReviewRevisionCandidateGateInput{
            Direction:      validExternalApplyCapabilityReviewImprovementDirectionForCandidateGate(testCase.direction),
            NonAuthorizing: true,
        })
        if got.Decision != testCase.decision || got.CandidateDigest != "" {
            t.Fatalf("direction %q got %+v", testCase.direction, got)
        }
        if err := got.Validate(); err != nil {
            t.Fatal(err)
        }
    }
}

func TestGateJEVExternalApplyCapabilityReviewRevisionCandidateRejectsMissingCandidate(t *testing.T) {
    got := GateJEVExternalApplyCapabilityReviewRevisionCandidate(JEVExternalApplyCapabilityReviewRevisionCandidateGateInput{
        Direction:      validExternalApplyCapabilityReviewImprovementDirectionForCandidateGate(jevExternalApplyCapabilityReviewImprovementDirectionGenerate),
        NonAuthorizing: true,
    })
    if got.Status != jevExternalApplyCapabilityReviewRevisionCandidateUnknown || got.MissingStage != "candidate-digest" {
        t.Fatalf("got %+v", got)
    }
}
