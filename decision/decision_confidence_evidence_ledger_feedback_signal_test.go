package decision

import "testing"

func TestDecisionConfidenceEvidenceLedgerFeedbackSignalDoesNotSelectChange(t *testing.T) {
	signal, err := DeriveDecisionConfidenceEvidenceLedgerFeedbackSignal(DecisionConfidenceEvidenceLedgerCandidateReviewMetricComparison{
		PreviousSignalDigest: "previous",
		CurrentSignalDigest:  "current",
		Delta:                "declined",
		NonAuthorizing:       true,
	})
	if err != nil {
		t.Fatalf("derive declined feedback: %v", err)
	}
	if signal.Status != "review-required" || signal.ComparisonDigest == "" || !signal.NonAuthorizing {
		t.Fatalf("unexpected declined feedback: %#v", signal)
	}

	signal, err = DeriveDecisionConfidenceEvidenceLedgerFeedbackSignal(DecisionConfidenceEvidenceLedgerCandidateReviewMetricComparison{
		PreviousSignalDigest: "previous",
		CurrentSignalDigest:  "current",
		Delta:                "increased",
		NonAuthorizing:       true,
	})
	if err != nil {
		t.Fatalf("derive increased feedback: %v", err)
	}
	if signal.Status != "observation-only" || !signal.NonAuthorizing {
		t.Fatalf("increased observation was promoted: %#v", signal)
	}
}
