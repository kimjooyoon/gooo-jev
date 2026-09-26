package gooo

import (
	"fmt"
	"strings"
)

type RevisionApplicationPlan struct {
	Status          string
	MissingStage    string
	SourceDigest    string
	InputIRDigest   string
	BindingDigest   string
	CandidateDigest string
	Edit            SourceEdit
	EditDigest      string
	PlanDigest      string
	NonExecuting    bool
	NonAuthorizing  bool
}

// PlanRevisionApplication records the validated inputs for one bounded revision.
func PlanRevisionApplication(source string, binding RevisionCandidateBinding, edit SourceEdit) (RevisionApplicationPlan, error) {
	plan := RevisionApplicationPlan{
		Status:          "UNKNOWN",
		MissingStage:    "revision-application-plan",
		SourceDigest:    digestString(source),
		BindingDigest:   binding.BindingDigest,
		CandidateDigest: binding.CandidateDigest,
		Edit:            edit,
		EditDigest:      edit.Digest,
		NonExecuting:    true,
		NonAuthorizing:  true,
	}
	if err := binding.Validate(); err != nil {
		plan.MissingStage = "revision-application-plan-binding"
		return plan, fmt.Errorf("gooo revision application plan: binding: %w", err)
	}
	if plan.SourceDigest != binding.SourceDigest {
		plan.MissingStage = "revision-application-plan-source"
		return plan, fmt.Errorf("gooo revision application plan: source digest precondition failed")
	}
	if edit.Digest == "" || digestSourceEdit(edit) != edit.Digest {
		plan.MissingStage = "revision-application-plan-edit"
		return plan, fmt.Errorf("gooo revision application plan: edit digest does not match its fields")
	}
	if edit.Start.Line != edit.End.Line || strings.Contains(edit.Replacement, "\n") {
		plan.MissingStage = "revision-application-plan-range"
		return plan, fmt.Errorf("gooo revision application plan: edit must stay on one line")
	}
	if _, err := positionOffset(source, edit.Start); err != nil {
		plan.MissingStage = "revision-application-plan-range"
		return plan, fmt.Errorf("gooo revision application plan: start: %w", err)
	}
	if _, err := positionOffset(source, edit.End); err != nil {
		plan.MissingStage = "revision-application-plan-range"
		return plan, fmt.Errorf("gooo revision application plan: end: %w", err)
	}
	inputIRDigest, err := parseRevisionSource(source)
	if err != nil {
		plan.MissingStage = "revision-application-plan-input-ir"
		return plan, fmt.Errorf("gooo revision application plan: source: %w", err)
	}
	plan.Status = "BOUND"
	plan.MissingStage = ""
	plan.InputIRDigest = inputIRDigest
	plan.PlanDigest = digestRevisionApplicationPlan(plan)
	return plan, nil
}

func (p RevisionApplicationPlan) Validate() error {
	if p.Status != "BOUND" {
		return fmt.Errorf("revision application plan status must be BOUND")
	}
	if p.MissingStage != "" {
		return fmt.Errorf("revision application plan missing stage must be empty")
	}
	if !p.NonExecuting || !p.NonAuthorizing {
		return fmt.Errorf("revision application plan must remain non-executing and non-authorizing")
	}
	for name, digest := range map[string]string{
		"source": p.SourceDigest, "input IR": p.InputIRDigest,
		"binding": p.BindingDigest, "candidate": p.CandidateDigest,
		"edit": p.EditDigest, "plan": p.PlanDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("revision application plan %s digest is invalid", name)
		}
	}
	if p.Edit.Digest != p.EditDigest || digestSourceEdit(p.Edit) != p.EditDigest {
		return fmt.Errorf("revision application plan edit is not linked")
	}
	if p.Edit.Start.Line != p.Edit.End.Line || strings.Contains(p.Edit.Replacement, "\n") {
		return fmt.Errorf("revision application plan range is invalid")
	}
	if err := p.Candidate.Validate(); err != nil {
		return fmt.Errorf("revision application plan candidate: %w", err)
	}
	if p.Candidate.CandidateDigest != p.CandidateDigest {
		return fmt.Errorf("revision application plan candidate digest is not linked")
	}
	if digestRevisionApplicationPlan(p) != p.PlanDigest {
		return fmt.Errorf("revision application plan digest does not match its fields")
	}
	return nil
}

func digestRevisionApplicationPlan(plan RevisionApplicationPlan) string {
	return digestString(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%t|%t",
		plan.Status,
		plan.MissingStage,
		plan.SourceDigest,
		plan.InputIRDigest,
		plan.BindingDigest,
		plan.CandidateDigest,
		digestSourceEdit(plan.Edit),
		plan.NonExecuting,
		plan.NonAuthorizing,
	))
}
