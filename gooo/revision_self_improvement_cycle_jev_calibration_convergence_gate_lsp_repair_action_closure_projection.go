package gooo

// JEVCalibrationConvergenceGateLSPRepairActionClosureProjectionStatus is the
// status shown by a read-only LSP projection.
type JEVCalibrationConvergenceGateLSPRepairActionClosureProjectionStatus string

const (
	JEVCalibrationConvergenceGateLSPRepairActionClosureProjectionBound    JEVCalibrationConvergenceGateLSPRepairActionClosureProjectionStatus = "BOUND"
	JEVCalibrationConvergenceGateLSPRepairActionClosureProjectionUnknown  JEVCalibrationConvergenceGateLSPRepairActionClosureProjectionStatus = "UNKNOWN"
	JEVCalibrationConvergenceGateLSPRepairActionClosureProjectionDeferred JEVCalibrationConvergenceGateLSPRepairActionClosureProjectionStatus = "DEFERRED"
)

// JEVCalibrationConvergenceGateLSPRepairActionClosureProjection is an LSP-safe
// view of a closure metric. It deliberately contains no edit or command.
type JEVCalibrationConvergenceGateLSPRepairActionClosureProjection struct {
	Status           JEVCalibrationConvergenceGateLSPRepairActionClosureProjectionStatus
	Code             string
	Title            string
	Kind             string
	MetricDigest     string
	MissingStage     string
	MissingStageIndex int
	IsActionable     bool
	IsReadOnly       bool
	CanExecute       bool
	CanAuthorize     bool
	Edits            []string
	Command          string
}

// ProjectJEVCalibrationConvergenceGateLSPRepairActionClosure projects a
// closure metric without turning provenance evidence into execution authority.
func ProjectJEVCalibrationConvergenceGateLSPRepairActionClosure(metric JEVCalibrationConvergenceGateLSPRepairActionClosureMetric) JEVCalibrationConvergenceGateLSPRepairActionClosureProjection {
	projection := JEVCalibrationConvergenceGateLSPRepairActionClosureProjection{
		Status:            JEVCalibrationConvergenceGateLSPRepairActionClosureProjectionUnknown,
		Code:              "gooo.provenance.unknown",
		Title:             "Inspect provenance closure",
		Kind:              "quickfix",
		MetricDigest:      metric.MetricDigest,
		MissingStage:      metric.MissingStage,
		MissingStageIndex: metric.MissingStageIndex,
		IsActionable:      true,
		IsReadOnly:        true,
		CanExecute:        false,
		CanAuthorize:      false,
	}

	switch metric.Status {
	case JEVCalibrationConvergenceGateLSPRepairActionClosureBound:
		projection.Status = JEVCalibrationConvergenceGateLSPRepairActionClosureProjectionBound
		projection.Code = "gooo.provenance.bound"
		projection.Title = "Provenance closure is bound"
		projection.Kind = "info"
		projection.IsActionable = false
	case JEVCalibrationConvergenceGateLSPRepairActionClosureDeferred:
		projection.Status = JEVCalibrationConvergenceGateLSPRepairActionClosureProjectionDeferred
		projection.Code = "gooo.provenance.deferred"
		projection.Title = "Provenance closure is deferred"
	case JEVCalibrationConvergenceGateLSPRepairActionClosureUnknown:
		projection.Status = JEVCalibrationConvergenceGateLSPRepairActionClosureProjectionUnknown
	}

	return projection
}