package decision

import "testing"

func validExternalApplyCapabilityReviewTrendSignalForLSP(signal string) JEVExternalApplyCapabilityReviewTrendSignal {
    value := JEVExternalApplyCapabilityReviewTrendSignal{
        Status:         jevExternalApplyCapabilityReviewTrendSignalRecorded,
        Signal:         signal,
        MetricName:     "review.confidence",
        TrendDigest:    "trend-digest",
        MeanValue:      0.75,
        Threshold:      0.7,
        EvidenceDigest: "trend-evidence-digest",
        NonExecuting:   true,
        NonAuthorizing: true,
    }
    value.SignalDigest = digestJEVExternalApplyCapabilityReviewTrendSignal(value.Status, value.Signal, value.MetricName, value.TrendDigest, value.MeanValue, value.Threshold, value.EvidenceDigest)
    return value
}

func TestProjectJEVExternalApplyCapabilityReviewTrendSignalLSPPreservesSignals(t *testing.T) {
    cases := []struct {
        signal   string
        severity string
        code     string
    }{
        {jevExternalApplyCapabilityReviewTrendSignalImprove, "info", "jev.external-apply.trend.improve"},
        {jevExternalApplyCapabilityReviewTrendSignalHold, "info", "jev.external-apply.trend.hold"},
        {jevExternalApplyCapabilityReviewTrendSignalRollback, "warning", "jev.external-apply.trend.rollback"},
    }
    for _, testCase := range cases {
        got := ProjectJEVExternalApplyCapabilityReviewTrendSignalLSP(validExternalApplyCapabilityReviewTrendSignalForLSP(testCase.signal))
        if got.Severity != testCase.severity || got.Code != testCase.code || !got.Publishable {
            t.Fatalf("signal %q got %+v", testCase.signal, got)
        }
        if err := got.Validate(); err != nil {
            t.Fatal(err)
        }
    }
}

func TestProjectJEVExternalApplyCapabilityReviewTrendSignalLSPUnknownIsNotPublishable(t *testing.T) {
    got := ProjectJEVExternalApplyCapabilityReviewTrendSignalLSP(JEVExternalApplyCapabilityReviewTrendSignal{
        Status:         jevExternalApplyCapabilityReviewTrendSignalUnknown,
        MissingStage:   "signal-threshold",
        NonExecuting:   true,
        NonAuthorizing: true,
    })
    if got.Publishable || got.Severity != "error" || got.Code != "jev.provenance.unknown" {
        t.Fatalf("got %+v", got)
    }
}
