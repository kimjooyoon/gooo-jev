package decision

import "testing"

func validChangePlanForProvenance(t *testing.T) DecisionConfidenceChangePlan {
	t.Helper()
	plan := DecisionConfidenceChangePlan{
		ProposalDigest:  "proposal-digest",
		SourceReference: "gooo://gooo-jev/declaration/example",
		ChangeDigest:    "change-digest",
		WriteSetDigest:  "write-set-digest",
		NonExecuting:    true,
	}
	digest, err := Digest(plan)
	if err != nil {
		t.Fatalf("digest change plan: %v", err)
	}
	plan.ChangePlanDigest = digest
	if err := plan.Validate(); err != nil {
		t.Fatalf("validate change plan: %v", err)
	}
	return plan
}

func TestBindDecisionConfidenceChangePlanProvenanceBindsSourceAndGeneration(t *testing.T) {
	binding := BindDecisionConfidenceChangePlanProvenance(DecisionConfidenceChangePlanProvenanceBindingInput{
		ChangePlan:              validChangePlanForProvenance(t),
		DeclarationIRGeneration: validDeclarationIRGenerationBinding(),
		NonAuthorizing:          true,
	})
	if binding.Status != "bound" || binding.MissingStage != "" || binding.EvidenceDigest == "" {
		t.Fatalf("unexpected bound plan: %+v", binding)
	}
	if err := binding.Validate(); err != nil {
		t.Fatalf("bound plan should validate: %v", err)
	}
	if binding.SourceReference != binding.DeclarationID || binding.GenerationDigest == "" {
		t.Fatalf("plan origin was not preserved: %+v", binding)
	}
}

func TestBindDecisionConfidenceChangePlanProvenanceReviewsSourceMismatch(t *testing.T) {
	plan := validChangePlanForProvenance(t)
	plan.SourceReference = "gooo://gooo-jev/declaration/other"
	plan.ChangePlanDigest = ""
	digest, err := Digest(plan)
	if err != nil {
		t.Fatalf("digest mismatched source plan: %v", err)
	}
	plan.ChangePlanDigest = digest
	binding := BindDecisionConfidenceChangePlanProvenance(DecisionConfidenceChangePlanProvenanceBindingInput{
		ChangePlan:              plan,
		DeclarationIRGeneration: validDeclarationIRGenerationBinding(),
		NonAuthorizing:          true,
	})
	if binding.Status != "review" || binding.MissingStage != "source-reference" {
		t.Fatalf("unexpected source mismatch: %+v", binding)
	}
}

func TestBindDecisionConfidenceChangePlanProvenanceFailsClosed(t *testing.T) {
	plan := validChangePlanForProvenance(t)
	plan.ChangePlanDigest = "tampered"
	binding := BindDecisionConfidenceChangePlanProvenance(DecisionConfidenceChangePlanProvenanceBindingInput{
		ChangePlan:              plan,
		DeclarationIRGeneration: validDeclarationIRGenerationBinding(),
		NonAuthorizing:          true,
	})
	if binding.Status != "UNKNOWN" || binding.MissingStage != "change-plan" {
		t.Fatalf("unexpected tampered plan: %+v", binding)
	}

	output := BindDecisionConfidenceChangePlanProvenance(DecisionConfidenceChangePlanProvenanceBindingInput{
		ChangePlan:              validChangePlanForProvenance(t),
		DeclarationIRGeneration: validDeclarationIRGenerationBinding(),
		NonAuthorizing:          false,
	})
	if output.Status != "UNKNOWN" || output.NonAuthorizing || output.MissingStage != "authorization-boundary" {
		t.Fatalf("unexpected authorization output: %+v", output)
	}
}
