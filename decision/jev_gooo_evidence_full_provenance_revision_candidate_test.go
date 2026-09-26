package decision

import "testing"

func evidenceFullProvenanceRevisionCandidateInput(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionCandidateInput {
	t.Helper()
	bridge := evidenceFullProvenanceFeedbackBridgeInput(t)
	ledger := AppendExecutionEnvelopeGoooEvidenceFullProvenanceFeedbackBridge(bridge)
	directive := JEVImprovementDirectionDirective{
		Status:              jevImprovementDirectiveRevision,
		Directive:           jevImprovementDirectiveRevision,
		CandidateDigest:     "candidate-digest",
		CandidateSource:     "gooo://candidate/evidence",
		InputEvidenceDigest: ledger.EvidenceDigest,
		NonExecuting:        true,
		NonAuthorizing:      true,
	}
	directive.EvidenceDigest = digestJEVImprovementDirectionDirective(directive.Status, directive.Directive, directive.CandidateDigest, directive.CandidateSource, directive.InputEvidenceDigest)
	return ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionCandidateInput{
		EvidenceBinding:      BindExecutionEnvelopeGoooEvidenceFullProvenance(evidenceFullProvenanceInput(t)),
		Ledger:               ledger,
		Directive:            directive,
		RevisionSource:       "gooo://revision/source/extended-evidence",
		RevisionChangeDigest: "revision-change-input",
		NonAuthorizing:       true,
	}
}

func TestGenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionCandidate(t *testing.T) {
	input := evidenceFullProvenanceRevisionCandidateInput(t)
	binding := GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionCandidate(input)
	if binding.Status != "bound" || binding.CandidateStatus == "" ||
		binding.CandidateDigest == "" || binding.CandidateEvidenceDigest == "" ||
		binding.BoundRevisionChangeDigest == "" || binding.BindingDigest == "" {
		t.Fatalf("binding = %#v, want bound revision candidate", binding)
	}
	if !binding.NonExecuting || !binding.NonAuthorizing {
		t.Fatalf("missing safety boundary: %#v", binding)
	}
}

func TestGenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionCandidateRejectsDirectionMismatch(t *testing.T) {
	input := evidenceFullProvenanceRevisionCandidateInput(t)
	input.Directive.InputEvidenceDigest = "tampered"
	binding := GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionCandidate(input)
	if binding.Status != "UNKNOWN" || binding.MissingStage != "feedback-direction-binding" {
		t.Fatalf("binding = %#v, want feedback-direction-binding UNKNOWN", binding)
	}
}

func TestGenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionCandidateRejectsUnknownFeedback(t *testing.T) {
	input := evidenceFullProvenanceRevisionCandidateInput(t)
	input.Ledger.Status = "UNKNOWN"
	binding := GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionCandidate(input)
	if binding.Status != "UNKNOWN" || binding.MissingStage != "replay-feedback-ledger" {
		t.Fatalf("binding = %#v, want replay-feedback-ledger UNKNOWN", binding)
	}
}

func TestGenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionCandidateRejectsAuthorization(t *testing.T) {
	binding := GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionCandidate(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionCandidateInput{
		NonAuthorizing: false,
	})
	if binding.Status != "UNKNOWN" || binding.MissingStage != "authorization-boundary" ||
		binding.NonAuthorizing {
		t.Fatalf("binding = %#v, want authorization UNKNOWN", binding)
	}
}
