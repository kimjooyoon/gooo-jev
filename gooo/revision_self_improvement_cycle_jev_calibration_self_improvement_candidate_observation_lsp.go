package gooo

// JEVSelfImprovementCandidateObservationLSPStatus is the editor-facing
// severity for a source-to-reverse-observation evidence boundary.
type JEVSelfImprovementCandidateObservationLSPStatus string

const (
	JEVSelfImprovementCandidateObservationLSPInformation JEVSelfImprovementCandidateObservationLSPStatus = "Information"
	JEVSelfImprovementCandidateObservationLSPError       JEVSelfImprovementCandidateObservationLSPStatus = "Error"
)

// JEVSelfImprovementCandidateObservationLSPInput mirrors the runtime
// observation without coupling the core to a particular runtime package.
type JEVSelfImprovementCandidateObservationLSPInput struct {
	Status                   string
	SourceDigest             string
	IRDigest                 string
	GeneratedDigest          string
	ReverseObservationDigest string
	ExpectedDelta            float64
	ObservedDelta            float64
	MissingStageIndex        int
}

// JEVSelfImprovementCandidateObservationLSP is a read-only projection that
// keeps evidence and signed deltas visible without claiming improvement.
type JEVSelfImprovementCandidateObservationLSP struct {
	Status                   JEVSelfImprovementCandidateObservationLSPStatus
	Code                     string
	Message                  string
	SourceDigest             string
	IRDigest                 string
	GeneratedDigest          string
	ReverseObservationDigest string
	ExpectedDelta            float64
	ObservedDelta            float64
	MissingStageIndex        int
	IsReadOnly               bool
	CanEdit                  bool
	CanExecute               bool
	CanAuthorize             bool
}

// ProjectJEVSelfImprovementCandidateObservationLSP maps complete or incomplete
// evidence to diagnostics while preserving the first unresolved stage.
func ProjectJEVSelfImprovementCandidateObservationLSP(
	input JEVSelfImprovementCandidateObservationLSPInput,
) JEVSelfImprovementCandidateObservationLSP {
	projection := JEVSelfImprovementCandidateObservationLSP{
		Status:                   JEVSelfImprovementCandidateObservationLSPError,
		Code:                     "jev.self_improvement_candidate.unknown",
		Message:                  "source-to-reverse-observation evidence is incomplete",
		SourceDigest:             input.SourceDigest,
		IRDigest:                 input.IRDigest,
		GeneratedDigest:          input.GeneratedDigest,
		ReverseObservationDigest: input.ReverseObservationDigest,
		ExpectedDelta:            input.ExpectedDelta,
		ObservedDelta:            input.ObservedDelta,
		MissingStageIndex:        input.MissingStageIndex,
		IsReadOnly:               true,
		CanEdit:                  false,
		CanExecute:               false,
		CanAuthorize:             false,
	}

	if input.Status == "BOUND" {
		projection.Status = JEVSelfImprovementCandidateObservationLSPInformation
		projection.Code = "jev.self_improvement_candidate.bound"
		projection.Message = "source-to-reverse-observation evidence chain is complete"
	}

	return projection
}