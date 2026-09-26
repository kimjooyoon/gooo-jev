package decision

import "testing"

func TestImprovementCycleReplayClassifiesReplayedDivergedAndUnknown(t *testing.T) {
    cycle, err := BuildImprovementCycle(
        "ledger-before",
        "ledger-after",
        "candidate://change/replay",
        "gooo://jev/source/replay",
        ImprovementObserved,
    )
    if err != nil {
        t.Fatalf("BuildImprovementCycle() error = %v", err)
    }

    replayed, err := ReplayImprovementCycle(cycle, "ledger-after", "evidence-replayed")
    if err != nil {
        t.Fatalf("replayed cycle error = %v", err)
    }
    if replayed.Status != ImprovementReplayed {
        t.Fatalf("replayed status = %q", replayed.Status)
    }
    if err := replayed.Validate(); err != nil {
        t.Fatalf("replayed Validate() error = %v", err)
    }

    diverged, err := ReplayImprovementCycle(cycle, "ledger-changed", "evidence-diverged")
    if err != nil {
        t.Fatalf("diverged cycle error = %v", err)
    }
    if diverged.Status != ImprovementDiverged {
        t.Fatalf("diverged status = %q", diverged.Status)
    }
    if err := diverged.Validate(); err != nil {
        t.Fatalf("diverged Validate() error = %v", err)
    }

    unknown, err := ReplayImprovementCycle(cycle, "", "evidence-unknown-replay")
    if err != nil {
        t.Fatalf("unknown cycle error = %v", err)
    }
    if unknown.Status != ImprovementReplayUnknown {
        t.Fatalf("unknown status = %q", unknown.Status)
    }
    if err := unknown.Validate(); err != nil {
        t.Fatalf("unknown Validate() error = %v", err)
    }

    tampered := replayed
    tampered.NonAuthorizing = false
    if err := tampered.Validate(); err == nil {
        t.Fatal("authorizing replay unexpectedly validated")
    }
}
