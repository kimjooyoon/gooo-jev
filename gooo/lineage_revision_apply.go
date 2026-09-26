package gooo

import "fmt"

// LineageRevisionReceipt binds a revision application to its complete provenance chain.
type LineageRevisionReceipt struct {
	Status               string
	MissingStage         string
	SourceDigest         string
	LineageDigest        string
	CandidateDigest      string
	ApplicationDigest    string
	ProposedSourceDigest string
	ResultDigest         string
	ProposedSource       string
	Application          RevisionApplication
	NonExecuting         bool
	NonAuthorizing       bool
}

// ApplyLineageRevision applies a revision only when lineage and source bindings agree.
func ApplyLineageRevision(source string, lineage LineageReceipt, candidate RevisionCandidate, edit SourceEdit) (LineageRevisionReceipt, error) {
	receipt := LineageRevisionReceipt{
		Status:          "UNKNOWN",
		MissingStage:    "lineage-revision-application",
		SourceDigest:    digestString(source),
		LineageDigest:   lineage.ChainDigest,
		CandidateDigest: candidate.CandidateDigest,
		NonExecuting:    true,
		NonAuthorizing:  true,
	}
	if err := lineage.Validate(); err != nil {
		receipt.MissingStage = "lineage-application-lineage"
		return receipt, fmt.Errorf("gooo lineage revision: lineage: %w", err)
	}
	if err := candidate.Validate(); err != nil {
		receipt.MissingStage = "lineage-application-candidate"
		return receipt, fmt.Errorf("gooo lineage revision: candidate: %w", err)
	}
	if lineage.SourceDigest != receipt.SourceDigest {
		receipt.MissingStage = "lineage-application-source"
		return receipt, fmt.Errorf("gooo lineage revision: source is not the lineage source")
	}
	if lineage.CandidateDigest != candidate.CandidateDigest {
		receipt.MissingStage = "lineage-application-binding"
		return receipt, fmt.Errorf("gooo lineage revision: candidate is not the lineage candidate")
	}

	application, err := ApplyRevision(source, lineage.SourceDigest, candidate, edit)
	receipt.Application = application
	receipt.ApplicationDigest = application.ApplicationDigest
	receipt.ProposedSource = application.ProposedSource
	receipt.ProposedSourceDigest = application.ProposedSourceDigest
	if err != nil {
		receipt.MissingStage = application.MissingStage
		return receipt, fmt.Errorf("gooo lineage revision: application: %w", err)
	}

	receipt.Status = "BOUND"
	receipt.MissingStage = ""
	receipt.ResultDigest = digestLineageRevision(receipt)
	return receipt, nil
}

// Validate verifies nested application evidence and the complete lineage binding.
func (r LineageRevisionReceipt) Validate() error {
	if r.Status != "BOUND" {
		return fmt.Errorf("lineage revision status must be BOUND")
	}
	if r.MissingStage != "" {
		return fmt.Errorf("lineage revision missing stage must be empty")
	}
	if !r.NonExecuting || !r.NonAuthorizing {
		return fmt.Errorf("lineage revision must remain non-executing and non-authorizing")
	}
	for name, digest := range map[string]string{
		"source": r.SourceDigest, "lineage": r.LineageDigest, "candidate": r.CandidateDigest,
		"application": r.ApplicationDigest, "proposed source": r.ProposedSourceDigest, "result": r.ResultDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("lineage revision %s digest is invalid", name)
		}
	}
	if r.ProposedSource == "" {
		return fmt.Errorf("lineage revision proposed source is required")
	}
	if err := r.Application.Validate(); err != nil {
		return fmt.Errorf("lineage revision application: %w", err)
	}
	if r.ApplicationDigest != r.Application.ApplicationDigest || r.ProposedSourceDigest != r.Application.ProposedSourceDigest {
		return fmt.Errorf("lineage revision nested application binding does not match")
	}
	if expected := digestLineageRevision(r); expected != r.ResultDigest {
		return fmt.Errorf("lineage revision result digest does not match")
	}
	return nil
}

func digestLineageRevision(receipt LineageRevisionReceipt) string {
	return digestString(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s",
		receipt.SourceDigest,
		receipt.LineageDigest,
		receipt.CandidateDigest,
		receipt.ApplicationDigest,
		receipt.ProposedSourceDigest,
		receipt.ProposedSource,
		receipt.NonExecuting,
	))
}
