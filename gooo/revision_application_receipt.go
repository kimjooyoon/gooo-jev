package gooo

import "fmt"

type RevisionApplicationReceipt struct {
	Status               string
	MissingStage         string
	PlanDigest           string
	ApplicationDigest    string
	SourceDigest         string
	ProposedSourceDigest string
	InputIRDigest        string
	ProposedIRDigest     string
	CandidateDigest      string
	EditDigest           string
	ReceiptDigest        string
	NonExecuting         bool
	NonAuthorizing       bool
}

func ObserveRevisionApplicationReceipt(plan RevisionApplicationPlan, application RevisionApplication) (RevisionApplicationReceipt, error) {
	receipt := RevisionApplicationReceipt{
		Status:               "UNKNOWN",
		MissingStage:         "revision-application-receipt",
		PlanDigest:           plan.PlanDigest,
		ApplicationDigest:    application.ApplicationDigest,
		SourceDigest:         application.SourceDigest,
		ProposedSourceDigest: application.ProposedSourceDigest,
		InputIRDigest:        application.InputIRDigest,
		ProposedIRDigest:     application.ProposedIRDigest,
		CandidateDigest:      application.CandidateDigest,
		EditDigest:           application.EditDigest,
		NonExecuting:         true,
		NonAuthorizing:       true,
	}
	if err := plan.Validate(); err != nil {
		receipt.MissingStage = "revision-application-receipt-plan"
		return receipt, fmt.Errorf("gooo revision application receipt: plan: %w", err)
	}
	if err := application.Validate(); err != nil {
		receipt.MissingStage = "revision-application-receipt-application"
		return receipt, fmt.Errorf("gooo revision application receipt: application: %w", err)
	}
	if application.SourceDigest != plan.SourceDigest ||
		application.InputIRDigest != plan.InputIRDigest ||
		application.CandidateDigest != plan.CandidateDigest ||
		application.EditDigest != plan.EditDigest {
		receipt.MissingStage = "revision-application-receipt-link"
		return receipt, fmt.Errorf("gooo revision application receipt: application is not linked to plan")
	}
	receipt.Status = "BOUND"
	receipt.MissingStage = ""
	receipt.PlanDigest = plan.PlanDigest
	receipt.ReceiptDigest = digestRevisionApplicationReceipt(receipt)
	return receipt, nil
}

func (r RevisionApplicationReceipt) Validate() error {
	if r.Status != "BOUND" {
		return fmt.Errorf("revision application receipt status must be BOUND")
	}
	if r.MissingStage != "" {
		return fmt.Errorf("revision application receipt missing stage must be empty")
	}
	if !r.NonExecuting || !r.NonAuthorizing {
		return fmt.Errorf("revision application receipt must remain non-executing and non-authorizing")
	}
	for name, digest := range map[string]string{
		"plan": r.PlanDigest, "application": r.ApplicationDigest,
		"source": r.SourceDigest, "proposed source": r.ProposedSourceDigest,
		"input IR": r.InputIRDigest, "proposed IR": r.ProposedIRDigest,
		"candidate": r.CandidateDigest, "edit": r.EditDigest,
		"receipt": r.ReceiptDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("revision application receipt %s digest is invalid", name)
		}
	}
	if digestRevisionApplicationReceipt(r) != r.ReceiptDigest {
		return fmt.Errorf("revision application receipt digest does not match its fields")
	}
	return nil
}

func digestRevisionApplicationReceipt(receipt RevisionApplicationReceipt) string {
	return digestString(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%t|%t",
		receipt.Status,
		receipt.MissingStage,
		receipt.PlanDigest,
		receipt.ApplicationDigest,
		receipt.SourceDigest,
		receipt.ProposedSourceDigest,
		receipt.InputIRDigest,
		receipt.ProposedIRDigest,
		receipt.CandidateDigest,
		receipt.EditDigest,
		receipt.NonExecuting,
		receipt.NonAuthorizing,
	))
}
