package main

import "testing"

func TestBuildFeedbackKeepsDiscoveryBoundary(t *testing.T) {
    feedback, err := buildFeedback(
        "What can gooo do with this declaration?",
        "Which constraints can be checked next?",
        "useful",
        "AVAILABLE",
        "",
        "evidence-1",
        "capability.discovery.feedback.v1",
    )
    if err != nil {
        t.Fatalf("buildFeedback() error = %v", err)
    }
    if !feedback.NonExecuting || !feedback.NonAuthorizing {
        t.Fatalf("feedback crossed the execution boundary: %#v", feedback)
    }
    if feedback.OriginalQuery == feedback.SuggestedQuery {
        t.Fatal("original and suggested questions must remain distinguishable")
    }
}

func TestBuildFeedbackRejectsMissingEvidence(t *testing.T) {
    if _, err := buildFeedback("question", "next", "useful", "AVAILABLE", "", "", "v1"); err == nil {
        t.Fatal("expected missing evidence to be rejected")
    }
}