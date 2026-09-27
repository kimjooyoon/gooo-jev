package gooo

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

type JEVCalibrationOutcomeWindowObservationStatus string

const (
	JEVCalibrationOutcomeWindowObservationBound    JEVCalibrationOutcomeWindowObservationStatus = "BOUND"
	JEVCalibrationOutcomeWindowObservationDeferred JEVCalibrationOutcomeWindowObservationStatus = "DEFERRED"
	JEVCalibrationOutcomeWindowObservationUnknown  JEVCalibrationOutcomeWindowObservationStatus = "UNKNOWN"
)

// JEVCalibrationOutcomeWindowObservationInput is the evidence required to
// observe a calibration outcome window without mutating its metrics.
type JEVCalibrationOutcomeWindowObservationInput struct {
	SourceVersion           string
	ContractVersion         string
	ChoiceSetDigest         string
	ExpectedChoiceSetDigest string
	EvidencePrefixDigest    string
	ObservationDigest       string
	OutcomeCount            int
	ProducerDeferred        bool
}

// JEVCalibrationOutcomeWindowObservation is a read-only metric observation.
// It is not a score update, execution decision, or authorization result.
type JEVCalibrationOutcomeWindowObservation struct {
	Status               JEVCalibrationOutcomeWindowObservationStatus
	SourceVersion        string
	ContractVersion      string
	ChoiceSetDigest      string
	EvidencePrefixDigest string
	ObservationDigest    string
	OutcomeCount         int
	TargetStage          string
	Reason               string
	RecordDigest          string
	IsReadOnly            bool
	CanExecute            bool
	CanAuthorize          bool
	Edits                []string
	Command              string
}

func ObserveJEVCalibrationOutcomeWindow(
	input JEVCalibrationOutcomeWindowObservationInput,
) JEVCalibrationOutcomeWindowObservation {
	observation := JEVCalibrationOutcomeWindowObservation{
		Status:               JEVCalibrationOutcomeWindowObservationUnknown,
		SourceVersion:        input.SourceVersion,
		ContractVersion:      input.ContractVersion,
		ChoiceSetDigest:      input.ChoiceSetDigest,
		EvidencePrefixDigest: input.EvidencePrefixDigest,
		ObservationDigest:    input.ObservationDigest,
		OutcomeCount:         input.OutcomeCount,
		TargetStage:          "outcome_window",
		IsReadOnly:           true,
		CanExecute:           false,
		CanAuthorize:         false,
	}

	switch {
	case input.SourceVersion == "" || input.ContractVersion == "":
		observation.TargetStage = "identity"
		observation.Reason = "source or contract identity is missing"
	case input.ChoiceSetDigest == "":
		observation.TargetStage = "choice_set"
		observation.Reason = "choice-set digest is missing"
	case input.ExpectedChoiceSetDigest == "":
		observation.TargetStage = "expected_choice_set"
		observation.Reason = "expected choice-set digest is missing"
	case input.EvidencePrefixDigest == "":
		observation.TargetStage = "evidence_prefix"
		observation.Reason = "evidence prefix digest is missing"
	case input.ObservationDigest == "":
		observation.TargetStage = "observation"
		observation.Reason = "outcome observation digest is missing"
	case input.ProducerDeferred:
		observation.Status = JEVCalibrationOutcomeWindowObservationDeferred
		observation.Reason = "outcome window producer is deferred"
	case input.ChoiceSetDigest != input.ExpectedChoiceSetDigest:
		observation.TargetStage = "choice_set_match"
		observation.Reason = "choice-set digest does not match expected digest"
	case input.OutcomeCount <= 0:
		observation.TargetStage = "outcome_count"
		observation.Reason = "outcome window has no positive observations"
	default:
		observation.Status = JEVCalibrationOutcomeWindowObservationBound
		observation.Reason = "choice-set and evidence-prefix metrics are aligned"
	}

	observation.RecordDigest = digestJEVCalibrationOutcomeWindow(
		string(observation.Status),
		observation.SourceVersion,
		observation.ContractVersion,
		observation.ChoiceSetDigest,
		input.ExpectedChoiceSetDigest,
		observation.EvidencePrefixDigest,
		observation.ObservationDigest,
		strconv.Itoa(observation.OutcomeCount),
		observation.TargetStage,
		observation.Reason,
	)
	return observation
}

func digestJEVCalibrationOutcomeWindow(parts ...string) string {
	sum := sha256.Sum256([]byte(strings.Join(parts, "|")))
	return fmt.Sprintf("sha256:%s", hex.EncodeToString(sum[:]))
}