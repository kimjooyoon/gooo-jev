package gooo

import "testing"

func TestCapabilityDiscoveryFeedbackCanonicalDigest(t *testing.T) {
    feedback := CapabilityDiscoveryFeedback{
        Version:           "capability.discovery.feedback.v1",
        OriginalQuery:     "What can gooo do with this declaration?",
        SuggestedQuery:    "Which constraints can be checked next?",
        Outcome:           "useful",
        CapabilityStatus:  "AVAILABLE",
        EvidenceDigest:    "evidence-1",
        NonExecuting:      true,
        NonAuthorizing:    true,
    }

    first, err := feedback.CanonicalDigest()
    if err != nil {
        t.Fatalf("canonical digest: %v", err)
    }
    second, err := feedback.CanonicalDigest()
    if err != nil {
        t.Fatalf("canonical digest repeat: %v", err)
    }
    if first == "" || first != second {
        t.Fatalf("digest is not stable: %q %q", first, second)
    }
}

func TestCapabilityDiscoveryFeedbackRejectsAuthorizingObservation(t *testing.T) {
    feedback := CapabilityDiscoveryFeedback{
        Version:          "capability.discovery.feedback.v1",
        OriginalQuery:    "What can gooo do?",
        SuggestedQuery:   "What can be inspected?",
        Outcome:          "unresolved",
        CapabilityStatus: "UNKNOWN",
        EvidenceDigest:   "evidence-2",
        NonExecuting:     true,
        NonAuthorizing:   false,
    }

    if err := feedback.Validate(); err == nil {
        t.Fatal("expected authorizing feedback to be rejected")
    }
}