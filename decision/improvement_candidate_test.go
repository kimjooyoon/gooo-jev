package decision

import "testing"

func TestImprovementCandidateRemainsNonExecuting(t *testing.T) {
    cycle, err := BuildImprovementCycle(
        "ledger-before-candidate",
        "ledger-after-candidate",
        "candidate://change/candidate",
        "gooo://jev/source/candidate",
        ImprovementObserved,
    )
    if err != nil {
        t.Fatalf("BuildImprovementCycle() error = %v", err)
    }
    replay, err := ReplayImprovementCycle(cycle, "ledger-after-candidate", "evidence-candidate")
    if err != nil {
        t.Fatalf("ReplayImprovementCycle() error = %v", err)
    }

    for _, decision := range []ImprovementCandidateDecision{
        CandidateProposed,
        CandidateReview,
        CandidateRejected,
    } {
        t.Run(string(decision), func(t *testing.T) {
            candidate, err := ProposeImprovementCandidate(
                cycle,
                replay,
                "candidate://change/candidate",
                "gooo://jev/source/candidate",
                decision,
            )
            if err != nil {
                t.Fatalf("ProposeImprovementCandidate() error = %v", err)
            }
            if !candidate.NonExecuting {
                t.Fatal("candidate unexpectedly became executable")
            }
            if err := candidate.Validate(); err != nil {
                t.Fatalf("Validate() error = %v", err)
            }
        })
    }

    mismatchCycle, err := BuildImprovementCycle(
        "ledger-other-before",
        "ledger-other-after",
        "candidate://change/other",
        "gooo://jev/source/other",
        ImprovementReview,
    )
    if err != nil {
        t.Fatalf("mismatch cycle error = %v", err)
    }
    if _, err := ProposeImprovementCandidate(
        mismatchCycle,
        replay,
        "candidate://change/other",
        "gooo://jev/source/other",
        CandidateReview,
    ); err == nil {
        t.Fatal("mismatched cycle and replay unexpectedly accepted")
    }

    tampered, err := ProposeImprovementCandidate(
        cycle,
        replay,
        "candidate://change/candidate",
        "gooo://jev/source/candidate",
        CandidateReview,
    )
    if err != nil {
        t.Fatalf("tampered fixture error = %v", err)
    }
    tampered.NonExecuting = false
    if err := tampered.Validate(); err == nil {
        t.Fatal("executing candidate unexpectedly validated")
    }
}
