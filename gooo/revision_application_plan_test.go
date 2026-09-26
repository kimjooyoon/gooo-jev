package gooo

import "testing"

func boundApplicationPlanInputs(t *testing.T) (RevisionCandidateBinding, SourceEdit) {
	t.Helper()
	selection, materialization := selectedCandidateInputs(t)
	binding, err := BindRevisionCandidateSelection(selection, materialization)
	if err != nil {
		t.Fatalf("BindRevisionCandidateSelection() error = %v", err)
	}
	edit := SourceEdit{
		Start:       Position{Line: 4, Column: 3},
		End:         Position{Line: 4, Column: 11},
		Replacement: "lineage",
	}
	edit.Digest = digestSourceEdit(edit)
	return binding, edit
}

func TestPlanRevisionApplicationBindsSourceAndIR(t *testing.T) {
	binding, edit := boundApplicationPlanInputs(t)
	plan, err := PlanRevisionApplication(validContract, binding, edit)
	if err != nil {
		t.Fatalf("PlanRevisionApplication() error = %v", err)
	}
	if plan.Status != "BOUND" || plan.InputIRDigest == "" ||
		plan.SourceDigest != binding.SourceDigest ||
		plan.CandidateDigest != binding.CandidateDigest {
		t.Fatalf("unexpected application plan: %#v", plan)
	}
	if err := plan.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	application, err := ApplyRevision(validContract, binding.Candidate, plan.Edit)
	if err != nil {
		t.Fatalf("ApplyRevision() error = %v", err)
	}
	if application.SourceDigest != plan.SourceDigest ||
		application.InputIRDigest != plan.InputIRDigest ||
		application.CandidateDigest != plan.CandidateDigest {
		t.Fatalf("application lost plan provenance: %#v", application)
	}
}

func TestPlanRevisionApplicationRetainsSourcePrecondition(t *testing.T) {
	binding, edit := boundApplicationPlanInputs(t)
	plan, err := PlanRevisionApplication(validContract+"\n", binding, edit)
	if err == nil {
		t.Fatal("PlanRevisionApplication() error = nil, want source precondition failure")
	}
	if plan.Status != "UNKNOWN" || plan.MissingStage != "revision-application-plan-source" {
		t.Fatalf("unexpected unknown application plan: %#v", plan)
	}
}

func TestPlanRevisionApplicationRetainsEditEvidenceStage(t *testing.T) {
	binding, edit := boundApplicationPlanInputs(t)
	edit.Digest = digestString("tampered")
	plan, err := PlanRevisionApplication(validContract, binding, edit)
	if err == nil {
		t.Fatal("PlanRevisionApplication() error = nil, want edit evidence failure")
	}
	if plan.Status != "UNKNOWN" || plan.MissingStage != "revision-application-plan-edit" {
		t.Fatalf("unexpected unknown application plan: %#v", plan)
	}
}

func TestPlanRevisionApplicationIsDeterministic(t *testing.T) {
	binding, edit := boundApplicationPlanInputs(t)
	first, err := PlanRevisionApplication(validContract, binding, edit)
	if err != nil {
		t.Fatalf("first PlanRevisionApplication() error = %v", err)
	}
	second, err := PlanRevisionApplication(validContract, binding, edit)
	if err != nil {
		t.Fatalf("second PlanRevisionApplication() error = %v", err)
	}
	if first.PlanDigest != second.PlanDigest {
		t.Fatal("same binding and edit produced different plan digest")
	}
}
