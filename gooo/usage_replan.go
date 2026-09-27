package gooo

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// UsageReplanDisposition describes a read-only next step derived from observed usage.
type UsageReplanDisposition string

const (
	UsageReplanProposeNextIteration     UsageReplanDisposition = "PROPOSE_NEXT_ITERATION"
	UsageReplanHoldForEvidence          UsageReplanDisposition = "HOLD_FOR_EVIDENCE"
	UsageReplanReviewNoConfirmedOutcome UsageReplanDisposition = "REVIEW_NO_CONFIRMED_OUTCOME"
)

// UsageReplanObservation is a non-executing, non-authorizing proposal derived from one usage window.
type UsageReplanObservation struct {
	Status            string                 `json:"status"`
	WindowDigest      string                 `json:"window_digest"`
	ObservationCount  int                    `json:"observation_count"`
	KnownCount        int                    `json:"known_count"`
	ConfirmedCount    int                    `json:"confirmed_count"`
	RefutedCount      int                    `json:"refuted_count"`
	UnknownCount      int                    `json:"unknown_count"`
	ObservedCoverage  float64                `json:"observed_coverage"`
	Disposition       UsageReplanDisposition  `json:"disposition"`
	NextOperation     string                 `json:"next_operation"`
	NonExecuting      bool                   `json:"non_executing"`
	NonAuthorizing    bool                   `json:"non_authorizing"`
	FirstMismatch     string                 `json:"first_mismatch,omitempty"`
	MissingStageIndex int                    `json:"missing_stage_index"`
	EvidenceDigest    string                 `json:"evidence_digest"`
}

// ProposeUsageReplan converts measured usage into a bounded proposal without executing it.
// Invalid or incomplete windows are preserved as UNKNOWN rather than promoted to a decision.
func ProposeUsageReplan(window UsageObservationWindow) UsageReplanObservation {
	observation := UsageReplanObservation{
		Status:            "BOUND",
		WindowDigest:      window.EvidenceDigest,
		ObservationCount:  window.ObservationCount,
		KnownCount:        window.KnownCount,
		ConfirmedCount:    window.ConfirmedCount,
		RefutedCount:      window.RefutedCount,
		UnknownCount:      window.UnknownCount,
		ObservedCoverage:  window.ObservedCoverage,
		NonExecuting:      true,
		NonAuthorizing:    true,
		MissingStageIndex: -1,
	}

	if err := window.Validate(); err != nil {
		observation.Status = "UNKNOWN"
		observation.Disposition = UsageReplanHoldForEvidence
		observation.NextOperation = "repair_usage_window_provenance"
		observation.FirstMismatch = "WINDOW_INTEGRITY"
		observation.MissingStageIndex = 0
	} else if window.UnknownCount > 0 || window.KnownCount == 0 {
		observation.Status = "UNKNOWN"
		observation.Disposition = UsageReplanHoldForEvidence
		observation.NextOperation = "collect_more_usage_evidence"
		observation.FirstMismatch = "INSUFFICIENT_KNOWN_USAGE_EVIDENCE"
		observation.MissingStageIndex = 0
	} else if window.ConfirmedCount > 0 {
		observation.Disposition = UsageReplanProposeNextIteration
		observation.NextOperation = "measure_next_usage_window"
	} else {
		observation.Disposition = UsageReplanReviewNoConfirmedOutcome
		observation.NextOperation = "review_refuted_or_inconclusive_usage"
	}

	observation.EvidenceDigest = usageReplanDigest(observation)
	return observation
}

// Validate verifies the proposal's own provenance and preserves UNKNOWN boundaries.
func (observation UsageReplanObservation) Validate() error {
	if observation.Status != "BOUND" && observation.Status != "UNKNOWN" {
		return fmt.Errorf("invalid usage replan status %q", observation.Status)
	}
	if observation.ObservationCount < 0 || observation.KnownCount < 0 || observation.ConfirmedCount < 0 || observation.RefutedCount < 0 || observation.UnknownCount < 0 {
		return fmt.Errorf("negative usage replan count")
	}
	if observation.KnownCount+observation.UnknownCount != observation.ObservationCount {
		return fmt.Errorf("usage replan counts do not reconcile")
	}
	if !observation.NonExecuting || !observation.NonAuthorizing {
		return fmt.Errorf("usage replan must remain non-executing and non-authorizing")
	}
	if observation.EvidenceDigest == "" {
		return fmt.Errorf("missing usage replan evidence digest")
	}
	if observation.Status == "BOUND" {
		if observation.WindowDigest == "" || observation.FirstMismatch != "" || observation.MissingStageIndex != -1 {
			return fmt.Errorf("bound usage replan has incomplete provenance")
		}
	} else if observation.FirstMismatch == "" || observation.MissingStageIndex < 0 {
		return fmt.Errorf("unknown usage replan must preserve first mismatch")
	}
	if expected := usageReplanDigest(observation); observation.EvidenceDigest != expected {
		return fmt.Errorf("usage replan evidence digest mismatch")
	}
	return nil
}

func usageReplanDigest(observation UsageReplanObservation) string {
	canonical := struct {
		Status            string                 `json:"status"`
		WindowDigest      string                 `json:"window_digest"`
		ObservationCount  int                    `json:"observation_count"`
		KnownCount        int                    `json:"known_count"`
		ConfirmedCount    int                    `json:"confirmed_count"`
		RefutedCount      int                    `json:"refuted_count"`
		UnknownCount      int                    `json:"unknown_count"`
		ObservedCoverage  float64                `json:"observed_coverage"`
		Disposition       UsageReplanDisposition  `json:"disposition"`
		NextOperation     string                 `json:"next_operation"`
		NonExecuting      bool                   `json:"non_executing"`
		NonAuthorizing    bool                   `json:"non_authorizing"`
		FirstMismatch     string                 `json:"first_mismatch,omitempty"`
		MissingStageIndex int                    `json:"missing_stage_index"`
	}{
		Status: observation.Status, WindowDigest: observation.WindowDigest,
		ObservationCount: observation.ObservationCount, KnownCount: observation.KnownCount,
		ConfirmedCount: observation.ConfirmedCount, RefutedCount: observation.RefutedCount,
		UnknownCount: observation.UnknownCount, ObservedCoverage: observation.ObservedCoverage,
		Disposition: observation.Disposition, NextOperation: observation.NextOperation,
		NonExecuting: observation.NonExecuting, NonAuthorizing: observation.NonAuthorizing,
		FirstMismatch: observation.FirstMismatch, MissingStageIndex: observation.MissingStageIndex,
	}
	payload, _ := json.Marshal(canonical)
	digest := sha256.Sum256(payload)
	return hex.EncodeToString(digest[:])
}