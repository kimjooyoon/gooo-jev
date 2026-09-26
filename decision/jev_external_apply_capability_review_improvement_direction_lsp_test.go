package decision

import "testing"

func validExternalApplyCapabilityReviewImprovementDirectionForLSP(direction string) JEVExternalApplyCapabilityReviewImprovementDirection {
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

func TestProjectJEVExternalApplyCapabilityReviewImprovementDirectionLSPPreservesDirections(t *testing.T) {
    cases := []struct {
        direction string
        severity  string
        code      string
    }{
        {jevExternalApplyCapabilityReviewImprovementDirectionGenerate, "info", "jev.external-apply.candidate-generation"},
        {jevExternalApplyCapabilityReviewImprovementDirectionHold, "info", "jev.external-apply.candidate-hold"},
        {jevExternalApplyCapabilityReviewImprovementDirectionReject, "warning", "jev.external-apply.candidate-rejection"},
    }
    for _, testCase := range cases {
        got := ProjectJEVExternalApplyCapabilityReviewImprovementDirectionLSP(validExternalApplyCapabilityReviewImprovementDirectionForLSP(testCase.direction))
        if got.Severity != testCase.severity || got.Code != testCase.code || !got.Publishable {
            t.Fatalf("direction %q got %+v", testCase.direction, got)
        }
        if err := got.Validate(); err != nil {
            t.Fatal(err)
        }
    }
}

func TestProjectJEVExternalApplyCapabilityReviewImprovementDirectionLSPUnknownIsNotPublishable(t *testing.T) {
    got := ProjectJEVExternalApplyCapabilityReviewImprovementDirectionLSP(JEVExternalApplyCapabilityReviewImprovementDirection{
        Status:         jevExternalApplyCapabilityReviewImprovementDirectionUnknown,
        MissingStage:   "feedback-bridge",
        NonExecuting:   true,
        NonAuthorizing: true,
    })
    if got.Publishable || got.Severity != "error" || got.Code != "jev.provenance.unknown" {
        t.Fatalf("got %+v", got)
    }
}
