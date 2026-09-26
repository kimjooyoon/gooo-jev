package decision

import "testing"

func TestImprovementReviewNeverGrantsExecution(t *testing.T) {
    cycle, err := BuildImprovementCycle(
        "ledger-review-before",
        "ledger-review-after",
        "candidate://change/review",
        "gooo://jev/source/review",
        ImprovementObserved,
    )
    if err != nil {
        t.Fatalf("BuildImprovementCycle() error = %v", err)
    }
    replay, err := ReplayImprovementCycle(cycle, "ledger-review-after", "evidence-review")
    if err != nil {
        t.Fatalf("ReplayImprovementCycle() error = %v", err)
    }
    candidate, err := ProposeImprovementCandidate(
        cycle,
        replay,
        "candidate://change/review",
        "gooo://jev/source/review",
        CandidateReview,
    )
    if err != nil {
        t.Fatalf("ProposeImprovementCandidate() error = %v", err)
    }

    for _, decision := range []ImprovementReviewDecision{
        ReviewRequired,
        ReviewPassed,
        ReviewFailed,
        ReviewUnknown,
    } {
        t.Run(string(decision), func(t *testing.T) {
            review, err := ReviewImprovementCandidate(
                candidate,
                "reviewer://human-or-policy",
                "review-evidence-"+string(decision),
                decision,
            )
            if err != nil {
                t.Fatalf("ReviewImprovementCandidate() error = %v", err)
            }
            if review.ExecutionGranted {
                t.Fatal("review unexpectedly granted execution")
            }
            if err := review.Validate(); err != nil {
                t.Fatalf("Validate() error = %v", err)
            }
        })
    }

    tampered, err := ReviewImprovementCandidate(
        candidate,
        "reviewer://human-or-policy",
        "review-evidence-tampered",
        ReviewPassed,
    )
    if err != nil {
        t.Fatalf("tampered fixture error = %v", err)
    }
    tampered.ExecutionGranted = true
    if err := tampered.Validate(); err == nil {
        t.Fatal("execution-granting review unexpectedly validated")
    }
}
