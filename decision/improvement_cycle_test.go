package decision

import "testing"

func TestImprovementCycleOutcomesRemainEvidenceBound(t *testing.T) {
    outcomes := []ImprovementCycleOutcome{
        ImprovementObserved,
        ImprovementRegressed,
        ImprovementUnknown,
        ImprovementReview,
    }
    for _, outcome := range outcomes {
        t.Run(string(outcome), func(t *testing.T) {
            cycle, err := BuildImprovementCycle(
                "ledger-before",
                "ledger-after",
                "candidate://change/1",
                "gooo://jev/source/1",
                outcome,
            )
            if err != nil {
                t.Fatalf("BuildImprovementCycle() error = %v", err)
            }
            if !cycle.NonAuthorizing {
                t.Fatal("cycle unexpectedly authorizes a change")
            }
            if err := cycle.Validate(); err != nil {
                t.Fatalf("Validate() error = %v", err)
            }
        })
    }

    tampered, err := BuildImprovementCycle(
        "ledger-before",
        "ledger-after",
        "candidate://change/2",
        "gooo://jev/source/2",
        ImprovementObserved,
    )
    if err != nil {
        t.Fatalf("tampered fixture construction error = %v", err)
    }
    tampered.NonAuthorizing = false
    if err := tampered.Validate(); err == nil {
        t.Fatal("authorizing cycle unexpectedly validated")
    }

    tampered = cycleWithDifferentCandidate(t, tampered)
    if err := tampered.Validate(); err == nil {
        t.Fatal("tampered candidate unexpectedly validated")
    }
}

func cycleWithDifferentCandidate(t *testing.T, cycle ImprovementCycle) ImprovementCycle {
    t.Helper()
    cycle.CandidateReference = "candidate://change/tampered"
    return cycle
}
