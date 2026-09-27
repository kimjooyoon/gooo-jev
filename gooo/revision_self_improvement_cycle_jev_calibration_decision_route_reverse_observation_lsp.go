package gooo

// JEVDecisionRouteReverseObservationLSPStatus is the editor-facing severity
// for a read-only decision-route reverse observation.
type JEVDecisionRouteReverseObservationLSPStatus string

const (
	JEVDecisionRouteReverseObservationLSPInformation JEVDecisionRouteReverseObservationLSPStatus = "Information"
	JEVDecisionRouteReverseObservationLSPWarning     JEVDecisionRouteReverseObservationLSPStatus = "Warning"
	JEVDecisionRouteReverseObservationLSPError       JEVDecisionRouteReverseObservationLSPStatus = "Error"
)

// JEVDecisionRouteReverseObservationLSPInput mirrors the runtime reverse
// observation without coupling the core to a runtime package.
type JEVDecisionRouteReverseObservationLSPInput struct {
	Status               string
	SourceVersion        string
	ContractVersion      string
	RouteStatus          string
	RouteRecordDigest    string
	DecisionDigest       string
	ExpectedRouteStatus  string
	ObservedRouteStatus  string
	OutcomeDigest        string
	TargetStage          string
	ObservationDigest    string
	ReviewRequired       bool
}

// JEVDecisionRouteReverseObservationLSP is a diagnostic-only projection. It
// never judges correctness or authorizes a route.
type JEVDecisionRouteReverseObservationLSP struct {
	Status               JEVDecisionRouteReverseObservationLSPStatus
	Code                 string
	Message              string
	SourceVersion        string
	ContractVersion      string
	RouteStatus          string
	RouteRecordDigest    string
	DecisionDigest       string
	ExpectedRouteStatus  string
	ObservedRouteStatus  string
	OutcomeDigest        string
	TargetStage           string
	ObservationDigest     string
	ReviewRequired        bool
	IsReadOnly            bool
	CanEdit               bool
	CanExecute            bool
	CanAuthorize          bool
}

// ProjectJEVDecisionRouteReverseObservationLSP preserves route and outcome
// evidence while keeping correctness and authority outside the projection.
func ProjectJEVDecisionRouteReverseObservationLSP(
	input JEVDecisionRouteReverseObservationLSPInput,
) JEVDecisionRouteReverseObservationLSP {
	projection := JEVDecisionRouteReverseObservationLSP{
		Status:              JEVDecisionRouteReverseObservationLSPError,
		Code:                "jev.route.reverse.unknown",
		Message:             "decision route reverse evidence is missing or unresolved",
		SourceVersion:       input.SourceVersion,
		ContractVersion:     input.ContractVersion,
		RouteStatus:         input.RouteStatus,
		RouteRecordDigest:   input.RouteRecordDigest,
		DecisionDigest:       input.DecisionDigest,
		ExpectedRouteStatus:  input.ExpectedRouteStatus,
		ObservedRouteStatus:  input.ObservedRouteStatus,
		OutcomeDigest:        input.OutcomeDigest,
		TargetStage:         input.TargetStage,
		ObservationDigest:   input.ObservationDigest,
		ReviewRequired:      input.ReviewRequired,
		IsReadOnly:           true,
		CanEdit:              false,
		CanExecute:           false,
		CanAuthorize:         false,
	}

	if input.Status == "DEFERRED" {
		projection.Status = JEVDecisionRouteReverseObservationLSPWarning
		projection.Code = "jev.route.reverse.deferred"
		projection.Message = "decision route reverse evidence is deferred"
		return projection
	}

	if input.Status != "BOUND" ||
		input.SourceVersion == "" ||
		input.ContractVersion == "" ||
		input.RouteRecordDigest == "" ||
		input.DecisionDigest == "" ||
		input.ExpectedRouteStatus == "" ||
		input.ObservedRouteStatus == "" ||
		input.ExpectedRouteStatus != input.ObservedRouteStatus ||
		input.OutcomeDigest == "" ||
		input.ObservationDigest == "" ||
		!validJEVDecisionRouteStatus(input.ExpectedRouteStatus) {
		projection.TargetStage = firstMissingJEVDecisionRouteReverseStage(input)
		projection.Message = "decision route reverse evidence is incomplete or mismatched"
		return projection
	}

	projection.Status = JEVDecisionRouteReverseObservationLSPInformation
	projection.Code = "jev.route.reverse.bound"
	projection.Message = "decision route and outcome evidence are available for inspection"
	switch input.ObservedRouteStatus {
	case "REVIEW_CANDIDATE":
		projection.Status = JEVDecisionRouteReverseObservationLSPWarning
		projection.Code = "jev.route.reverse.review"
		projection.Message = "review route and outcome evidence are available for human inspection"
	case "DEFERRED":
		projection.Status = JEVDecisionRouteReverseObservationLSPWarning
		projection.Code = "jev.route.reverse.deferred"
		projection.Message = "deferred route evidence is available for inspection"
	}
	return projection
}

func validJEVDecisionRouteStatus(status string) bool {
	switch status {
	case "ACCEPT_CANDIDATE", "REVIEW_CANDIDATE", "ABSTAIN_CANDIDATE", "DEFERRED":
		return true
	default:
		return false
	}
}

func firstMissingJEVDecisionRouteReverseStage(input JEVDecisionRouteReverseObservationLSPInput) string {
	switch {
	case input.SourceVersion == "" || input.ContractVersion == "":
		return "identity"
	case input.RouteRecordDigest == "":
		return "route_record"
	case input.DecisionDigest == "":
		return "decision_digest"
	case input.ExpectedRouteStatus == "":
		return "expected_route_status"
	case input.ObservedRouteStatus == "":
		return "observed_route_status"
	case input.ExpectedRouteStatus != input.ObservedRouteStatus:
		return "route_status_match"
	case input.OutcomeDigest == "":
		return "outcome"
	case input.ObservationDigest == "":
		return "reverse_observation"
	case !validJEVDecisionRouteStatus(input.ExpectedRouteStatus):
		return "route_status"
	default:
		return "reverse_observation"
	}
}
