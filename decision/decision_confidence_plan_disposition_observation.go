package decision

import "fmt"

type DecisionConfidenceChangePlanDispositionObservationStatus string

const (
	DecisionConfidenceChangePlanObservedApplied    DecisionConfidenceChangePlanDispositionObservationStatus = "applied"
	DecisionConfidenceChangePlanObservedAborted    DecisionConfidenceChangePlanDispositionObservationStatus = "aborted"
	DecisionConfidenceChangePlanObservedRolledBack DecisionConfidenceChangePlanDispositionObservationStatus = "rolled-back"
	DecisionConfidenceChangePlanObservedNotApplied DecisionConfidenceChangePlanDispositionObservationStatus = "not-applied"
	DecisionConfidenceChangePlanObservedUnknown   DecisionConfidenceChangePlanDispositionObservationStatus = "unknown"
)

// DecisionConfidenceChangePlanDispositionObservation records an external
// outcome without performing the outcome itself.
type DecisionConfidenceChangePlanDispositionObservation struct {
	DispositionDigest string
	Status            DecisionConfidenceChangePlanDispositionObservationStatus
	EvidenceDigest    string
	NonAuthorizing    bool
	ObservationDigest string
}

func ObserveDecisionConfidenceChangePlanDisposition(
	disposition DecisionConfidenceChangePlanDisposition,
	status DecisionConfidenceChangePlanDispositionObservationStatus,
	evidenceDigest string,
) (DecisionConfidenceChangePlanDispositionObservation, error) {
	if err := disposition.Validate(); err != nil {
		return DecisionConfidenceChangePlanDispositionObservation{}, fmt.Errorf("validate disposition: %w", err)
	}
	if status != DecisionConfidenceChangePlanObservedApplied && status != DecisionConfidenceChangePlanObservedAborted && status != DecisionConfidenceChangePlanObservedRolledBack && status != DecisionConfidenceChangePlanObservedNotApplied && status != DecisionConfidenceChangePlanObservedUnknown {
		return DecisionConfidenceChangePlanDispositionObservation{}, fmt.Errorf("unsupported disposition observation status %q", status)
	}
	if status != DecisionConfidenceChangePlanObservedUnknown && evidenceDigest == "" {
		return DecisionConfidenceChangePlanDispositionObservation{}, fmt.Errorf("non-unknown disposition observation requires evidence digest")
	}
	if disposition.Status == DecisionConfidenceChangePlanAbort && status == DecisionConfidenceChangePlanObservedApplied {
		return DecisionConfidenceChangePlanDispositionObservation{}, fmt.Errorf("aborted disposition cannot be observed as applied")
	}
	if disposition.Status == DecisionConfidenceChangePlanHold && status == DecisionConfidenceChangePlanObservedApplied {
		return DecisionConfidenceChangePlanDispositionObservation{}, fmt.Errorf("held disposition cannot be observed as applied")
	}
	observation := DecisionConfidenceChangePlanDispositionObservation{
		DispositionDigest: disposition.DispositionDigest,
		Status:            status,
		EvidenceDigest:    evidenceDigest,
		NonAuthorizing:    true,
	}
	digest, err := Digest(observation)
	if err != nil {
		return DecisionConfidenceChangePlanDispositionObservation{}, fmt.Errorf("digest disposition observation: %w", err)
	}
	observation.ObservationDigest = digest
	if err := observation.Validate(); err != nil {
		return DecisionConfidenceChangePlanDispositionObservation{}, fmt.Errorf("validate disposition observation: %w", err)
	}
	return observation, nil
}

func (o DecisionConfidenceChangePlanDispositionObservation) Validate() error {
	if o.DispositionDigest == "" {
		return fmt.Errorf("disposition digest is required")
	}
	if o.Status != DecisionConfidenceChangePlanObservedApplied && o.Status != DecisionConfidenceChangePlanObservedAborted && o.Status != DecisionConfidenceChangePlanObservedRolledBack && o.Status != DecisionConfidenceChangePlanObservedNotApplied && o.Status != DecisionConfidenceChangePlanObservedUnknown {
		return fmt.Errorf("unsupported disposition observation status %q", o.Status)
	}
	if o.Status != DecisionConfidenceChangePlanObservedUnknown && o.EvidenceDigest == "" {
		return fmt.Errorf("non-unknown disposition observation requires evidence digest")
	}
	if !o.NonAuthorizing {
		return fmt.Errorf("disposition observation must be non-authorizing")
	}
	if o.ObservationDigest == "" {
		return fmt.Errorf("observation digest is required")
	}
	withoutDigest := o
	withoutDigest.ObservationDigest = ""
	digest, err := Digest(withoutDigest)
	if err != nil {
		return fmt.Errorf("digest disposition observation: %w", err)
	}
	if digest != o.ObservationDigest {
		return fmt.Errorf("observation digest mismatch")
	}
	return nil
}
