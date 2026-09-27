package gooo

import "testing"

func TestProjectJEVProvenanceStageCoverageLSPBound(t *testing.T) {
	projection := ProjectJEVProvenanceStageCoverageLSP(
		JEVProvenanceStageCoverageLSPInput{
			Status:              "BOUND",
			SourceVersion:       "src-v1",
			ContractVersion:     "contract-v1",
			ExpectedStages:      []string{"source", "ir", "generated", "reverse_observation"},
			ObservedStages:      []string{"source", "ir", "generated", "reverse_observation"},
			CoverageNumerator:   4,
			CoverageDenominator: 4,
			CoveragePercent:     100,
			TargetStage:         "provenance_stage_coverage",
			ObservationDigest:   "sha256:observation",
			ObservationalOnly:   true,
		},
	)

	if projection.Status != JEVProvenanceStageCoverageLSPInformation {
		t.Fatalf("status = %q, want Information", projection.Status)
	}
	if projection.Code != "jev.provenance.stage_coverage.bound" {
		t.Fatalf("code = %q, want bound code", projection.Code)
	}
	if projection.CoverageNumerator != 4 || projection.CoveragePercent != 100 {
		t.Fatalf("coverage = %d/%d (%d%%), want 4/4 (100%%)", projection.CoverageNumerator, projection.CoverageDenominator, projection.CoveragePercent)
	}
	if !projection.ObservationOnly || projection.ClaimsImprovement || projection.CanEdit || projection.CanExecute || projection.CanAuthorize {
		t.Fatal("coverage projection crossed a forbidden boundary")
	}
}

func TestProjectJEVProvenanceStageCoverageLSPDeferred(t *testing.T) {
	projection := ProjectJEVProvenanceStageCoverageLSP(
		JEVProvenanceStageCoverageLSPInput{
			Status:              "DEFERRED",
			CoverageDenominator: 4,
			CoverageNumerator:   2,
			CoveragePercent:     50,
			MissingStage:        "generated",
			TargetStage:         "observed_stage",
			ObservationDigest:   "sha256:observation",
			ObservationalOnly:   true,
		},
	)

	if projection.Status != JEVProvenanceStageCoverageLSPWarning {
		t.Fatalf("status = %q, want Warning", projection.Status)
	}
	if projection.Code != "jev.provenance.stage_coverage.deferred" {
		t.Fatalf("code = %q, want deferred code", projection.Code)
	}
}

func TestProjectJEVProvenanceStageCoverageLSPUnknownPreservesMissingStage(t *testing.T) {
	projection := ProjectJEVProvenanceStageCoverageLSP(
		JEVProvenanceStageCoverageLSPInput{
			Status:            "UNKNOWN",
			CoverageDenominator: 4,
			CoverageNumerator: 2,
			CoveragePercent:   50,
			MissingStage:      "generated",
			TargetStage:       "observed_stage",
			ObservationalOnly: true,
		},
	)

	if projection.Status != JEVProvenanceStageCoverageLSPError {
		t.Fatalf("status = %q, want Error", projection.Status)
	}
	if projection.MissingStage != "generated" || projection.TargetStage != "observed_stage" {
		t.Fatalf("missing stage was not preserved: %#v", projection)
	}
}

func TestProjectJEVProvenanceStageCoverageLSPRejectsImprovementClaim(t *testing.T) {
	projection := ProjectJEVProvenanceStageCoverageLSP(
		JEVProvenanceStageCoverageLSPInput{
			Status:              "BOUND",
			SourceVersion:       "src-v1",
			ContractVersion:     "contract-v1",
			CoverageNumerator:   4,
			CoverageDenominator: 4,
			CoveragePercent:     100,
			ObservationDigest:   "sha256:observation",
			ObservationalOnly:   true,
			ClaimsImprovement:   true,
		},
	)

	if projection.Status != JEVProvenanceStageCoverageLSPError {
		t.Fatalf("status = %q, want Error", projection.Status)
	}
	if projection.ClaimsImprovement {
		t.Fatal("projection must never carry an improvement claim")
	}
}

func TestProjectJEVProvenanceStageCoverageLSPRejectsInvalidMetric(t *testing.T) {
	projection := ProjectJEVProvenanceStageCoverageLSP(
		JEVProvenanceStageCoverageLSPInput{
			Status:              "BOUND",
			SourceVersion:       "src-v1",
			ContractVersion:     "contract-v1",
			CoverageNumerator:   5,
			CoverageDenominator: 4,
			CoveragePercent:     125,
			ObservationDigest:   "sha256:observation",
			ObservationalOnly:   true,
		},
	)

	if projection.Status != JEVProvenanceStageCoverageLSPError {
		t.Fatalf("status = %q, want Error", projection.Status)
	}
}
