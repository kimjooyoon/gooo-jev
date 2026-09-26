package decision

import "testing"

func validExternalApplyCapabilityReviewRevisionCandidateGateForLSP(decision string) JEVExternalApplyCapabilityReviewRevisionCandidateGate {
    value := JEVExternalApplyCapabilityReviewRevisionCandidateGate{
        Status:          jevExternalApplyCapabilityReviewRevisionCandidateGated,
        Decision:        decision,
        Direction:       jevExternalApplyCapabilityReviewImprovementDirectionGenerate,
        Target:          "decision/provenance",
        CandidateSource: "review-metric-trend",
        DirectionDigest: "direction-digest",
        NonExecuting:    true,
        NonAuthorizing:  true,
    }
    if decision == jevExternalApplyCapabilityReviewRevisionCandidateReady {
        value.CandidateDigest = "candidate-digest"
    }
    value.GateDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateGate(value.Status, value.Decision, value.Direction, value.Target, value.CandidateSource, value.CandidateDigest, value.DirectionDigest)
    return value
}

func TestProjectJEVExternalApplyCapabilityReviewRevisionCandidateGateLSPPreservesDecisions(t *testing.T) {
    cases := []struct {
        decision string
        severity string
        code     string
    }{
        {jevExternalApplyCapabilityReviewRevisionCandidateReady, "info", "jev.external-apply.candidate-ready"},
        {jevExternalApplyCapabilityReviewRevisionCandidateHold, "info", "jev.external-apply.candidate-hold"},
        {jevExternalApplyCapabilityReviewRevisionCandidateRejected, "warning", "jev.external-apply.candidate-rejected"},
    }
    for _, testCase := range cases {
        got := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateGateLSP(validExternalApplyCapabilityReviewRevisionCandidateGateForLSP(testCase.decision))
        if got.Severity != testCase.severity || got.Code != testCase.code || !got.Publishable {
            t.Fatalf("decision %q got %+v", testCase.decision, got)
        }
        if err := got.Validate(); err != nil {
            t.Fatal(err)
        }
    }
}

func TestProjectJEVExternalApplyCapabilityReviewRevisionCandidateGateLSPUnknownIsNotPublishable(t *testing.T) {
    got := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateGateLSP(JEVExternalApplyCapabilityReviewRevisionCandidateGate{
        Status:         jevExternalApplyCapabilityReviewRevisionCandidateUnknown,
        MissingStage:   "candidate-digest",
        NonExecuting:   true,
        NonAuthorizing: true,
    })
    if got.Publishable || got.Severity != "error" || got.Code != "jev.provenance.unknown" {
        t.Fatalf("got %+v", got)
    }
}
