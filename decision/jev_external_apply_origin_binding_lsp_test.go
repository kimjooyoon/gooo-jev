package decision

import "testing"

func validExternalApplyOriginBindingForLSP(status string) JEVExternalApplyOriginBinding {
    binding := JEVExternalApplyOriginBinding{
        Status:                  status,
        OriginDigest:            "origin-digest",
        CandidateDigest:         "candidate-digest",
        RequestEvidenceDigest:   "request-evidence-digest",
        NonExecuting:            true,
        NonAuthorizing:          true,
    }
    binding.BindingEvidenceDigest = digestJEVExternalApplyOriginBinding(binding.OriginDigest, binding.CandidateDigest, binding.RequestEvidenceDigest)
    binding.EvidenceDigest = digestJEVExternalApplyOriginBindingEvidence(binding.Status, binding.OriginDigest, binding.BindingEvidenceDigest)
    return binding
}

func TestProjectJEVExternalApplyOriginBindingLSPReadyIsReviewOnly(t *testing.T) {
    got := ProjectJEVExternalApplyOriginBindingLSP(validExternalApplyOriginBindingForLSP(jevExternalApplyOriginBindingReady))
    if got.Severity != "info" || got.Code != "jev.external-apply.origin-bound" || !got.Publishable {
        t.Fatalf("got %+v", got)
    }
    if err := got.Validate(); err != nil {
        t.Fatal(err)
    }
}

func TestProjectJEVExternalApplyOriginBindingLSPUnknownIsNotPublishable(t *testing.T) {
    got := ProjectJEVExternalApplyOriginBindingLSP(JEVExternalApplyOriginBinding{
        Status:         jevExternalApplyOriginBindingUnknown,
        MissingStage:   "origin-request-binding",
        NonExecuting:   true,
        NonAuthorizing: true,
    })
    if got.Publishable || got.Severity != "error" || got.Code != "jev.provenance.unknown" {
        t.Fatalf("got %+v", got)
    }
}
