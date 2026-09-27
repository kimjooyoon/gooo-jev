package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

type RevisionSelfImprovementCycleJEVSubsetConvergenceInput struct {
	SubsetName                 string
	DeclarationDigest          string
	IRDigest                   string
	GenerationDigest           string
	ReverseObservationDigest   string
	ConvergenceIteration       int64
	Stable                     bool
	ReverseObserved            bool
	SourceLineCount            int64
	GeneratedLineCount         int64
}

type RevisionSelfImprovementCycleJEVSubsetConvergenceObservation struct {
	Status                    string
	MissingStage              string
	SubsetName                string
	DeclarationDigest         string
	IRDigest                  string
	GenerationDigest          string
	ReverseObservationDigest  string
	ConvergenceIteration      int64
	StabilityStatus            string
	ReverseObservationStatus   string
	SourceLineCount            int64
	GeneratedLineCount         int64
	ConvergenceSignal          string
	ObservationDigest          string
	ReadOnly                   bool
	NonExecuting               bool
	NonAuthorizing             bool
}

func ObserveRevisionSelfImprovementCycleJEVSubsetConvergence(
	input RevisionSelfImprovementCycleJEVSubsetConvergenceInput,
) RevisionSelfImprovementCycleJEVSubsetConvergenceObservation {
	result := RevisionSelfImprovementCycleJEVSubsetConvergenceObservation{
		Status:                   "UNKNOWN",
		MissingStage:             "revision-self-improvement-cycle-jev-subset-convergence",
		SubsetName:               input.SubsetName,
		DeclarationDigest:        input.DeclarationDigest,
		IRDigest:                input.IRDigest,
		GenerationDigest:        input.GenerationDigest,
		ReverseObservationDigest: input.ReverseObservationDigest,
		ConvergenceIteration:    input.ConvergenceIteration,
		StabilityStatus:          "unknown",
		ReverseObservationStatus: "unknown",
		SourceLineCount:          input.SourceLineCount,
		GeneratedLineCount:       input.GeneratedLineCount,
		ConvergenceSignal:        "jev-subset-convergence-unknown",
		ReadOnly:                 true,
		NonExecuting:             true,
		NonAuthorizing:           true,
	}
	if input.Stable {
		result.StabilityStatus = "stable"
	}
	if input.ReverseObserved {
		result.ReverseObservationStatus = "observed"
	}
	setDigest := func() {
		result.ObservationDigest = digestRevisionSelfImprovementCycleJEVSubsetConvergence(result)
	}
	setDigest()

	if strings.TrimSpace(input.SubsetName) == "" {
		result.MissingStage = "revision-self-improvement-cycle-jev-subset-convergence-subset"
		setDigest()
		return result
	}
	if !validDigest(input.DeclarationDigest) ||
		!validDigest(input.IRDigest) ||
		!validDigest(input.GenerationDigest) ||
		!validDigest(input.ReverseObservationDigest) {
		result.MissingStage = "revision-self-improvement-cycle-jev-subset-convergence-lineage"
		setDigest()
		return result
	}
	if input.ConvergenceIteration < 1 {
		result.MissingStage = "revision-self-improvement-cycle-jev-subset-convergence-iteration"
		setDigest()
		return result
	}
	if input.SourceLineCount < 1 || input.GeneratedLineCount < 1 {
		result.MissingStage = "revision-self-improvement-cycle-jev-subset-convergence-counts"
		setDigest()
		return result
	}
	if !input.Stable {
		result.MissingStage = "revision-self-improvement-cycle-jev-subset-convergence-stability"
		setDigest()
		return result
	}
	if !input.ReverseObserved {
		result.MissingStage = "revision-self-improvement-cycle-jev-subset-convergence-reverse-observation"
		setDigest()
		return result
	}

	result.Status = "BOUND"
	result.MissingStage = ""
	result.ConvergenceSignal = "jev-subset-convergence-observed"
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "revision-self-improvement-cycle-jev-subset-convergence"
		result.ConvergenceSignal = "jev-subset-convergence-unknown"
		setDigest()
	}
	return result
}

func (value RevisionSelfImprovementCycleJEVSubsetConvergenceObservation) Validate() error {
	if value.Status != "BOUND" && value.Status != "UNKNOWN" {
		return fmt.Errorf("subset convergence status is invalid")
	}
	if value.Status == "BOUND" && value.MissingStage != "" {
		return fmt.Errorf("bound subset convergence has a missing stage")
	}
	if value.Status == "UNKNOWN" {
		if value.MissingStage == "" || value.ConvergenceSignal != "jev-subset-convergence-unknown" {
			return fmt.Errorf("unknown subset convergence must preserve incomplete evidence")
		}
	} else {
		if strings.TrimSpace(value.SubsetName) == "" ||
			!validDigest(value.DeclarationDigest) ||
			!validDigest(value.IRDigest) ||
			!validDigest(value.GenerationDigest) ||
			!validDigest(value.ReverseObservationDigest) ||
			value.ConvergenceIteration < 1 ||
			value.StabilityStatus != "stable" ||
			value.ReverseObservationStatus != "observed" ||
			value.SourceLineCount < 1 ||
			value.GeneratedLineCount < 1 ||
			value.ConvergenceSignal != "jev-subset-convergence-observed" {
			return fmt.Errorf("bound subset convergence evidence is incomplete")
		}
	}
	if !value.ReadOnly || !value.NonExecuting || !value.NonAuthorizing {
		return fmt.Errorf("subset convergence must remain read-only, non-executing, and non-authorizing")
	}
	if value.ObservationDigest != digestRevisionSelfImprovementCycleJEVSubsetConvergence(value) {
		return fmt.Errorf("subset convergence digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementCycleJEVSubsetConvergence(
	value RevisionSelfImprovementCycleJEVSubsetConvergenceObservation,
) string {
	return digestString(strings.Join([]string{
		value.Status,
		value.MissingStage,
		value.SubsetName,
		value.DeclarationDigest,
		value.IRDigest,
		value.GenerationDigest,
		value.ReverseObservationDigest,
		strconv.FormatInt(value.ConvergenceIteration, 10),
		value.StabilityStatus,
		value.ReverseObservationStatus,
		strconv.FormatInt(value.SourceLineCount, 10),
		strconv.FormatInt(value.GeneratedLineCount, 10),
		value.ConvergenceSignal,
		strconv.FormatBool(value.ReadOnly),
		strconv.FormatBool(value.NonExecuting),
		strconv.FormatBool(value.NonAuthorizing),
	}, "|"))
}