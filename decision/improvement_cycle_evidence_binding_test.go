package decision

import "testing"

func makeImprovementCycleEvidenceBindingInput(outcome ImprovementCycleOutcome) ImprovementCycleEvidenceBindingInput {
	return ImprovementCycleEvidenceBindingInput{
		ObservationAdmissionStatus: "observation-admitted",
		CandidateDigest:            "candidate-digest",
		AdmissionDigest:            "admission-digest",
		BaselineLedgerDigest:       "ledger-before",
		ResultLedgerDigest:         "ledger-after",
		CandidateReference:         "candidate://cycle",
		SourceReference:            "gooo://jev/source/cycle",
		Outcome:                    outcome,
		NonAuthorizing:             true,
	}
}

func TestBindImprovementCycleEvidence(t *testing.T) {
	output := BindImprovementCycleEvidence(makeImprovementCycleEvidenceBindingInput(ImprovementObserved))
	if output.Status != "recorded" || output.Outcome != ImprovementObserved ||
		output.CycleDigest == "" || output.EvidenceDigest == "" ||
		!output.NonExecuting || !output.NonAuthorizing {
		t.Fatalf("unexpected observed cycle binding: %#v", output)
	}

	output = BindImprovementCycleEvidence(makeImprovementCycleEvidenceBindingInput(ImprovementRegressed))
	if output.Status != "recorded" || output.Outcome != ImprovementRegressed {
		t.Fatalf("unexpected regressed cycle binding: %#v", output)
	}

	output = BindImprovementCycleEvidence(makeImprovementCycleEvidenceBindingInput(ImprovementReview))
	if output.Status != "review" || output.MissingStage != "improvement-outcome" ||
		output.Outcome != ImprovementReview {
		t.Fatalf("unexpected review cycle binding: %#v", output)
	}

	input := makeImprovementCycleEvidenceBindingInput(ImprovementObserved)
	input.ObservationAdmissionStatus = "review"
	output = BindImprovementCycleEvidence(input)
	if output.Status != "review" || output.MissingStage != "observation-admission" {
		t.Fatalf("unexpected admission review: %#v", output)
	}

	input = makeImprovementCycleEvidenceBindingInput(ImprovementObserved)
	input.ResultLedgerDigest = ""
	output = BindImprovementCycleEvidence(input)
	if output.Status != "UNKNOWN" || output.MissingStage != "improvement-cycle" {
		t.Fatalf("unexpected incomplete cycle: %#v", output)
	}

	input = makeImprovementCycleEvidenceBindingInput(ImprovementObserved)
	input.NonAuthorizing = false
	output = BindImprovementCycleEvidence(input)
	if output.Status != "UNKNOWN" || output.NonAuthorizing ||
		output.MissingStage != "authorization-boundary" {
		t.Fatalf("unexpected authorization cycle: %#v", output)
	}
}