package gooo

import (
	"strings"
	"testing"
)

func TestExplainUsageRendersBoundedGuidance(t *testing.T) {
	source := "package support\nnamespace triage\nentity ticket\nproperty\n"
	discovery := DiscoverUsage(source, "pro")
	plan, err := PlanUsageActions(discovery)
	if err != nil {
		t.Fatalf("plan usage actions: %v", err)
	}
	assistance, err := ExplainUsage(discovery, plan)
	if err != nil {
		t.Fatalf("explain usage: %v", err)
	}
	if err := assistance.Validate(); err != nil {
		t.Fatalf("validate assistance: %v", err)
	}
	text, err := RenderUsageAssistance(assistance)
	if err != nil {
		t.Fatalf("render assistance: %v", err)
	}
	if !strings.Contains(assistance.Summary, "syntax_completion") ||
		!strings.Contains(assistance.Summary, "canonical_generation") ||
		!strings.Contains(text, "non-authorizing: true") {
		t.Fatalf("assistance did not explain bounded capabilities: %#v\n%s", assistance, text)
	}
}

func TestUsageAssistanceTamperingInvalidatesDigest(t *testing.T) {
	discovery := DiscoverUsage("package support\nnamespace triage\nentity ticket\nproperty\n", "")
	plan, err := PlanUsageActions(discovery)
	if err != nil {
		t.Fatalf("plan usage actions: %v", err)
	}
	assistance, err := ExplainUsage(discovery, plan)
	if err != nil {
		t.Fatalf("explain usage: %v", err)
	}
	assistance.Summary = "tampered"
	if err := assistance.Validate(); err == nil {
		t.Fatal("expected tampered assistance to invalidate digest")
	}
}
