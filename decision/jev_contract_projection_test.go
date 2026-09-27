package decision

import (
	"strings"
	"testing"
)

func completeGoooContractSource() string {
	return strings.Join([]string{
		"module support_triage_observation",
		"input ticket_state",
		"input typed_decision_signal",
		"input confidence",
		"input confidence_method",
		"input review_threshold",
		"status ACCEPT observation confidence meets a reviewable routing threshold",
		"status REVIEW observation remains below threshold or needs human review",
		"field status",
		"field route",
		"field confidence",
		"field review_threshold",
		"constraint observation_only",
		"constraint human_review_boundary",
		"constraint confidence_method_explicit",
		"constraint no_execution",
		"constraint no_authorization",
	}, "\n")
}

func TestDeriveGoooContractProjectionBindsDeclarationEvidence(t *testing.T) {
	projection := DeriveGoooContractProjection(GoooContractInput{
		SourceText:     completeGoooContractSource(),
		NonAuthorizing: true,
	})
	if projection.Status != "ready" {
		t.Fatalf("status = %q, want ready (missing %q)", projection.Status, projection.MissingStage)
	}
	if err := projection.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if projection.ContractDigest == "" || projection.Evidence.ReverseObservationDigest == "" {
		t.Fatal("projection did not retain contract and reverse-observation evidence")
	}
	if !strings.Contains(projection.DeclarationSource, "activity observe_support_triage_observation") {
		t.Fatal("projection did not render the support-triage observation activity")
	}
}

func TestDeriveGoooContractProjectionRequiresCapabilityConstraints(t *testing.T) {
	source := completeGoooContractSource()
	source = strings.Replace(source, "constraint no_authorization", "constraint human_review_boundary", 1)
	projection := DeriveGoooContractProjection(GoooContractInput{
		SourceText:     source,
		NonAuthorizing: true,
	})
	if projection.Status != "UNKNOWN" || projection.MissingStage != "contract-boundary" {
		t.Fatalf("projection = %#v, want UNKNOWN at contract-boundary", projection)
	}
	if err := projection.Validate(); err != nil {
		t.Fatalf("unknown projection Validate() error = %v", err)
	}
}
