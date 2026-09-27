package gooo

import "testing"

func TestProjectJEVSelfImprovementCandidateObservationLSPBound(t *testing.T) {
	projection := ProjectJEVSelfImprovementCandidateObservationLSP(
		JEVSelfImprovementCandidateObservationLSPInput{
			Status:                   "BOUND",
			SourceDigest:             "source-digest",
			IRDigest:                 "ir-digest",
			GeneratedDigest:          "generated-digest",
			ReverseObservationDigest: "reverse-digest",
			ExpectedDelta:            0.25,
			ObservedDelta:            -0.10,
			MissingStageIndex:        -1,
		},
	)

	if projection.Status != JEVSelfImprovementCandidateObservationLSPInformation {
		t.Fatalf("status = %q, want Information", projection.Status)
	}
	if projection.Code != "jev.self_improvement_candidate.bound" {
		t.Fatalf("code = %q, want jev.self_improvement_candidate.bound", projection.Code)
	}
	if projection.SourceDigest != "source-digest" ||
		projection.IRDigest != "ir-digest" ||
		projection.GeneratedDigest != "generated-digest" ||
		projection.ReverseObservationDigest != "reverse-digest" {
		t.Fatal("projection did not preserve the evidence chain")
	}
	if projection.ExpectedDelta != 0.25 || projection.ObservedDelta != -0.10 {
		t.Fatal("projection did not preserve signed deltas")
	}
	if !projection.IsReadOnly || projection.CanEdit || projection.CanExecute || projection.CanAuthorize {
		t.Fatal("LSP projection must remain read-only and non-authorizing")
	}
}

func TestProjectJEVSelfImprovementCandidateObservationLSPUnknown(t *testing.T) {
	projection := ProjectJEVSelfImprovementCandidateObservationLSP(
		JEVSelfImprovementCandidateObservationLSPInput{
			Status:            "UNKNOWN",
			SourceDigest:      "source-digest",
			IRDigest:          "ir-digest",
			MissingStageIndex: 2,
			ObservedDelta:     0.10,
		},
	)

	if projection.Status != JEVSelfImprovementCandidateObservationLSPError {
		t.Fatalf("status = %q, want Error", projection.Status)
	}
	if projection.MissingStageIndex != 2 {
		t.Fatalf("missing stage = %d, want 2", projection.MissingStageIndex)
	}
}