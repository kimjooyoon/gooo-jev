package decision

import "testing"

func cycleObservationDisposition() DecisionConfidenceChangePlanFeedbackDisposition {
    d := DecisionConfidenceChangePlanFeedbackDisposition{
        Status: changePlanFeedbackDispositionReady,
        FeedbackStatus: changePlanFeedbackStatusConfirmed,
        ChangePlanDigest: "plan-digest",
        VerificationDigest: "verification-digest",
        DispositionDigest: "disposition-digest",
        FeedbackEvidenceDigest: "feedback-evidence-digest",
        NonExecuting: true,
        NonAuthorizing: true,
    }
    d.EvidenceDigest = digestChangePlanFeedbackDisposition(d.ChangePlanDigest, d.VerificationDigest, d.DispositionDigest, d.FeedbackEvidenceDigest, d.FeedbackStatus)
    return d
}

func TestObserveJEVImprovementCycle(t *testing.T) {
    disposition := cycleObservationDisposition()
    lsp := ProjectDecisionConfidenceChangePlanFeedbackDispositionLSP(disposition)
    got := ObserveJEVImprovementCycle(JEVImprovementCycleObservationInput{
        DeclarationDigest: "declaration-digest",
        IRDigest: "ir-digest",
        GenerationDigest: "generation-digest",
        ReverseObservationDigest: "reverse-observation-digest",
        MetricDigest: "metric-digest",
        Disposition: disposition,
        LSPDiagnostic: lsp,
        NonAuthorizing: true,
    })
    if got.Status != changePlanFeedbackDispositionReady || got.MissingStage != "" || got.EvidenceDigest == "" {
        t.Fatalf("got %+v", got)
    }
    if err := got.Validate(); err != nil {
        t.Fatal(err)
    }
}

func TestObserveJEVImprovementCyclePreservesMissingReverseObservation(t *testing.T) {
    disposition := cycleObservationDisposition()
    lsp := ProjectDecisionConfidenceChangePlanFeedbackDispositionLSP(disposition)
    got := ObserveJEVImprovementCycle(JEVImprovementCycleObservationInput{
        DeclarationDigest: "declaration-digest",
        IRDigest: "ir-digest",
        GenerationDigest: "generation-digest",
        MetricDigest: "metric-digest",
        Disposition: disposition,
        LSPDiagnostic: lsp,
        NonAuthorizing: true,
    })
    if got.Status != jevImprovementCycleObservationUnknown || got.MissingStage != "reverse-observation" {
        t.Fatalf("got %+v", got)
    }
}
