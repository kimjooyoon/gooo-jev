package gooo

// JEVProvenanceStageCoverageLSPStatus is the editor-facing severity for an
// evidence coverage observation that makes no improvement claim.
type JEVProvenanceStageCoverageLSPStatus string

const (
	JEVProvenanceStageCoverageLSPInformation JEVProvenanceStageCoverageLSPStatus = "Information"
	JEVProvenanceStageCoverageLSPWarning     JEVProvenanceStageCoverageLSPStatus = "Warning"
	JEVProvenanceStageCoverageLSPError       JEVProvenanceStageCoverageLSPStatus = "Error"
)

// JEVProvenanceStageCoverageLSPInput mirrors the runtime metric without
// coupling the core to a runtime package.
type JEVProvenanceStageCoverageLSPInput struct {
	Status              string
	SourceVersion       string
	ContractVersion     string
	ExpectedStages      []string
	ObservedStages      []string
	CoverageNumerator   int
	CoverageDenominator int
	CoveragePercent     int
	MissingStage        string
	TargetStage         string
	ObservationDigest   string
	ObservationalOnly   bool
	ClaimsImprovement   bool
}

// JEVProvenanceStageCoverageLSP is a diagnostic-only projection. It reports
// evidence coverage but never asserts improvement or grants authority.
type JEVProvenanceStageCoverageLSP struct {
	Status              JEVProvenanceStageCoverageLSPStatus
	Code                string
	Message             string
	SourceVersion       string
	ContractVersion     string
	ExpectedStages      []string
	ObservedStages      []string
	CoverageNumerator   int
	CoverageDenominator int
	CoveragePercent     int
	MissingStage        string
	TargetStage         string
	ObservationDigest   string
	ObservationOnly     bool
	ClaimsImprovement   bool
	CanEdit             bool
	CanExecute          bool
	CanAuthorize       bool
}

// ProjectJEVProvenanceStageCoverageLSP exposes coverage evidence without
// converting it into a self-improvement conclusion.
func ProjectJEVProvenanceStageCoverageLSP(
	input JEVProvenanceStageCoverageLSPInput,
) JEVProvenanceStageCoverageLSP {
	projection := JEVProvenanceStageCoverageLSP{
		Status:              JEVProvenanceStageCoverageLSPError,
		Code:                "jev.provenance.stage_coverage.unknown",
		Message:             "provenance stage coverage evidence is missing or invalid",
		SourceVersion:       input.SourceVersion,
		ContractVersion:     input.ContractVersion,
		ExpectedStages:      append([]string(nil), input.ExpectedStages...),
		ObservedStages:      append([]string(nil), input.ObservedStages...),
		CoverageNumerator:   input.CoverageNumerator,
		CoverageDenominator: input.CoverageDenominator,
		CoveragePercent:     input.CoveragePercent,
		MissingStage:        input.MissingStage,
		TargetStage:         input.TargetStage,
		ObservationDigest:   input.ObservationDigest,
		ObservationOnly:     true,
		ClaimsImprovement:   false,
		CanEdit:             false,
		CanExecute:          false,
		CanAuthorize:        false,
	}

	if input.CoverageDenominator < 0 ||
		input.CoverageNumerator < 0 ||
		input.CoverageNumerator > input.CoverageDenominator ||
		input.CoveragePercent < 0 ||
		input.CoveragePercent > 100 ||
		input.ClaimsImprovement ||
		!input.ObservationalOnly {
		projection.Message = "provenance coverage input crosses the observation-only boundary"
		return projection
	}

	switch input.Status {
	case "BOUND":
		if input.SourceVersion == "" ||
			input.ContractVersion == "" ||
			input.CoverageDenominator == 0 ||
			input.CoverageNumerator != input.CoverageDenominator ||
			input.CoveragePercent != 100 ||
			input.MissingStage != "" ||
			input.ObservationDigest == "" {
			projection.Message = "coverage is marked BOUND but evidence is incomplete"
			return projection
		}
		projection.Status = JEVProvenanceStageCoverageLSPInformation
		projection.Code = "jev.provenance.stage_coverage.bound"
		projection.Message = "provenance stage coverage is available for inspection"
	case "DEFERRED":
		projection.Status = JEVProvenanceStageCoverageLSPWarning
		projection.Code = "jev.provenance.stage_coverage.deferred"
		projection.Message = "provenance stage coverage is deferred"
	}

	return projection
}
