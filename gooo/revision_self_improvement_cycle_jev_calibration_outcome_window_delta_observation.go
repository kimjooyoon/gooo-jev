package gooo

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

type JEVCalibrationOutcomeWindowDeltaObservationStatus string

const (
	JEVCalibrationOutcomeWindowDeltaBound    JEVCalibrationOutcomeWindowDeltaObservationStatus = "BOUND"
	JEVCalibrationOutcomeWindowDeltaDeferred JEVCalibrationOutcomeWindowDeltaObservationStatus = "DEFERRED"
	JEVCalibrationOutcomeWindowDeltaUnknown  JEVCalibrationOutcomeWindowDeltaObservationStatus = "UNKNOWN"
)

// JEVCalibrationOutcomeWindowDeltaObservationInput records a signed metric
// change without claiming that an increase is success or a decrease is failure.
type JEVCalibrationOutcomeWindowDeltaObservationInput struct {
	SourceVersion        string
	ContractVersion      string
	EvidencePrefixDigest string
	ObservationDigest    string
	ExpectedOutcomeCount int
	ObservedOutcomeCount int
	ProducerDeferred     bool
}

// JEVCalibrationOutcomeWindowDeltaObservation is a read-only signed metric
// observation. It cannot edit, execute, or authorize.
type JEVCalibrationOutcomeWindowDeltaObservation struct {
	Status               JEVCalibrationOutcomeWindowDeltaObservationStatus
	SourceVersion        string
	ContractVersion      string
	EvidencePrefixDigest string
	ObservationDigest    string
	ExpectedOutcomeCount int
	ObservedOutcomeCount int
	Delta                int
	TargetStage          string
	Reason               string
	RecordDigest         string
	IsReadOnly           bool
	CanExecute           bool
	CanAuthorize         bool
	Edits                []string
	Command              string
}

func ObserveJEVCalibrationOutcomeWindowDelta(
	input JEVCalibrationOutcomeWindowDeltaObservationInput,
) JEVCalibrationOutcomeWindowDeltaObservation {
	observation := JEVCalibrationOutcomeWindowDeltaObservation{
		Status:               JEVCalibrationOutcomeWindowDeltaUnknown,
		SourceVersion:        input.SourceVersion,
		ContractVersion:      input.ContractVersion,
		EvidencePrefixDigest: input.EvidencePrefixDigest,
		ObservationDigest:    input.ObservationDigest,
		ExpectedOutcomeCount: input.ExpectedOutcomeCount,
		ObservedOutcomeCount: input.ObservedOutcomeCount,
		Delta:                input.ObservedOutcomeCount - input.ExpectedOutcomeCount,
		TargetStage:          "outcome_window_delta",
		IsReadOnly:           true,
		CanExecute:           false,
		CanAuthorize:         false,
	}

	switch {
	case input.SourceVersion == "" || input.ContractVersion == "":
		observation.TargetStage = "identity"
		observation.Reason = "source or contract identity is missing"
	case input.EvidencePrefixDigest == "":
		observation.TargetStage = "evidence_prefix"
		observation.Reason = "evidence prefix digest is missing"
	case input.ObservationDigest == "":
		observation.TargetStage = "observation"
		observation.Reason = "outcome observation digest is missing"
	case input.ProducerDeferred:
		observation.Status = JEVCalibrationOutcomeWindowDeltaDeferred
		observation.Reason = "outcome window delta producer is deferred"
	case input.ExpectedOutcomeCount < 0 || input.ObservedOutcomeCount < 0:
		observation.TargetStage = "outcome_count"
		observation.Reason = "outcome counts must be non-negative"
	default:
		observation.Status = JEVCalibrationOutcomeWindowDeltaBound
		observation.Reason = "signed outcome-count delta is recorded without an improvement claim"
	}

	observation.RecordDigest = digestJEVCalibrationOutcomeWindowDelta(
		string(observation.Status),
		observation.SourceVersion,
		observation.ContractVersion,
		observation.EvidencePrefixDigest,
		observation.ObservationDigest,
		strconv.Itoa(observation.ExpectedOutcomeCount),
		strconv.Itoa(observation.ObservedOutcomeCount),
		strconv.Itoa(observation.Delta),
		observation.TargetStage,
		observation.Reason,
	)
	return observation
}

func digestJEVCalibrationOutcomeWindowDelta(parts ...string) string {
	sum := sha256.Sum256([]byte(strings.Join(parts, "|")))
	return fmt.Sprintf("sha256:%s", hex.EncodeToString(sum[:]))
}
