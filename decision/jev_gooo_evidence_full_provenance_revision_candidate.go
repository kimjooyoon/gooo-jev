package decision

import "strings"

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionCandidateInput connects
// extended provenance and replay feedback to the existing revision generator.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionCandidateInput struct {
	EvidenceBinding       ExecutionEnvelopeGoooEvidenceFullProvenanceBinding
	Ledger                JEVImprovementReplayFeedbackLedger
	Directive             JEVImprovementDirectionDirective
	RevisionSource        string
	RevisionChangeDigest  string
	NonAuthorizing        bool
}

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionCandidateBinding
// preserves every evidence edge while keeping the candidate non-executing.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionCandidateBinding struct {
	Status                    string
	MissingStage              string
	EvidenceBindingDigest     string
	ProvenanceEvidenceDigest  string
	LedgerEvidenceDigest      string
	LedgerEntryDigest         string
	CandidateStatus           string
	CandidateDigest           string
	CandidateEvidenceDigest   string
	BoundRevisionChangeDigest string
	BindingDigest             string
	NonExecuting              bool
	NonAuthorizing            bool
}

func GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionCandidate(input ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionCandidateInput) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionCandidateBinding {
	output := ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionCandidateBinding{
		Status: "UNKNOWN", NonExecuting: true, NonAuthorizing: true,
	}
	if !input.NonAuthorizing || !input.EvidenceBinding.NonAuthorizing ||
		!input.Ledger.NonAuthorizing || !input.Directive.NonAuthorizing {
		if !input.NonAuthorizing {
			output.NonAuthorizing = false
		}
		output.MissingStage = "authorization-boundary"
		return output
	}
	if !input.EvidenceBinding.NonExecuting || !input.Ledger.NonExecuting || !input.Directive.NonExecuting {
		output.MissingStage = "execution-boundary"
		return output
	}
	if input.EvidenceBinding.Status != "complete" {
		output.MissingStage = input.EvidenceBinding.MissingStage
		if strings.TrimSpace(output.MissingStage) == "" {
			output.MissingStage = "evidence-full-provenance"
		}
		return output
	}
	if err := input.EvidenceBinding.Validate(); err != nil {
		output.MissingStage = "evidence-full-provenance-validation"
		return output
	}
	if err := input.Ledger.Validate(); err != nil {
		output.MissingStage = "replay-feedback-ledger"
		return output
	}
	if input.Ledger.FeedbackStatus == jevReplayFeedbackUnknown {
		output.MissingStage = "feedback-unknown"
		return output
	}
	if input.Directive.InputEvidenceDigest == "" ||
		input.Directive.InputEvidenceDigest != input.Ledger.EvidenceDigest {
		output.MissingStage = "feedback-direction-binding"
		return output
	}
	boundChangeDigest, err := Digest(struct {
		RevisionChangeDigest   string
		EvidenceBindingDigest  string
		ProvenanceEvidenceDigest string
		LedgerEvidenceDigest   string
		LedgerEntryDigest      string
		FeedbackStatus         string
	}{
		RevisionChangeDigest:    input.RevisionChangeDigest,
		EvidenceBindingDigest:   input.EvidenceBinding.EvidenceBindingDigest,
		ProvenanceEvidenceDigest: input.EvidenceBinding.EvidenceDigest,
		LedgerEvidenceDigest:    input.Ledger.EvidenceDigest,
		LedgerEntryDigest:       input.Ledger.EntryDigest,
		FeedbackStatus:          input.Ledger.FeedbackStatus,
	})
	if err != nil {
		output.MissingStage = "revision-change-digest"
		return output
	}
	candidate := GenerateJEVImprovementRevisionCandidate(JEVImprovementRevisionCandidateInput{
		Directive:            input.Directive,
		RevisionSource:      input.RevisionSource,
		RevisionChangeDigest: boundChangeDigest,
		NonAuthorizing:      true,
	})
	if err := candidate.Validate(); err != nil {
		output.MissingStage = "revision-candidate"
		return output
	}
	bindingDigest, err := Digest(struct {
		EvidenceBindingDigest     string
		ProvenanceEvidenceDigest  string
		LedgerEvidenceDigest      string
		LedgerEntryDigest         string
		CandidateDigest           string
		CandidateEvidenceDigest   string
		BoundRevisionChangeDigest string
	}{
		EvidenceBindingDigest:     input.EvidenceBinding.EvidenceBindingDigest,
		ProvenanceEvidenceDigest: input.EvidenceBinding.EvidenceDigest,
		LedgerEvidenceDigest:    input.Ledger.EvidenceDigest,
		LedgerEntryDigest:       input.Ledger.EntryDigest,
		CandidateDigest:          candidate.CandidateDigest,
		CandidateEvidenceDigest: candidate.EvidenceDigest,
		BoundRevisionChangeDigest: boundChangeDigest,
	})
	if err != nil {
		output.MissingStage = "revision-candidate-binding-digest"
		return output
	}
	output.Status = "bound"
	output.EvidenceBindingDigest = input.EvidenceBinding.EvidenceBindingDigest
	output.ProvenanceEvidenceDigest = input.EvidenceBinding.EvidenceDigest
	output.LedgerEvidenceDigest = input.Ledger.EvidenceDigest
	output.LedgerEntryDigest = input.Ledger.EntryDigest
	output.CandidateStatus = candidate.Status
	output.CandidateDigest = candidate.CandidateDigest
	output.CandidateEvidenceDigest = candidate.EvidenceDigest
	output.BoundRevisionChangeDigest = boundChangeDigest
	output.BindingDigest = bindingDigest
	return output
}
