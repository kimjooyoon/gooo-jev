package decision

import "fmt"

type DecisionConfidenceChangePlanDispositionStatus string

const (
	DecisionConfidenceChangePlanReadyForExternalApplyReview DecisionConfidenceChangePlanDispositionStatus = "ready-for-external-apply-review"
	DecisionConfidenceChangePlanAbort                       DecisionConfidenceChangePlanDispositionStatus = "abort"
	DecisionConfidenceChangePlanHold                        DecisionConfidenceChangePlanDispositionStatus = "hold"
)

// DecisionConfidenceChangePlanDisposition records the next safe boundary
// after verification; it never applies or rolls back a write set itself.
type DecisionConfidenceChangePlanDisposition struct {
	ChangePlanDigest   string
	VerificationDigest string
	Status             DecisionConfidenceChangePlanDispositionStatus
	NonExecuting       bool
	DispositionDigest  string
}

func DispositionDecisionConfidenceChangePlan(
	plan DecisionConfidenceChangePlan,
	verification DecisionConfidenceChangePlanVerification,
) (DecisionConfidenceChangePlanDisposition, error) {
	if err := plan.Validate(); err != nil {
		return DecisionConfidenceChangePlanDisposition{}, fmt.Errorf("validate change plan: %w", err)
	}
	if err := verification.Validate(); err != nil {
		return DecisionConfidenceChangePlanDisposition{}, fmt.Errorf("validate change plan verification: %w", err)
	}
	if verification.ChangePlanDigest != plan.ChangePlanDigest {
		return DecisionConfidenceChangePlanDisposition{}, fmt.Errorf("verification does not reference change plan")
	}
	disposition := DecisionConfidenceChangePlanDisposition{
		ChangePlanDigest:   plan.ChangePlanDigest,
		VerificationDigest: verification.VerificationDigest,
		NonExecuting:       true,
	}
	switch verification.Status {
	case DecisionConfidenceChangePlanVerified:
		disposition.Status = DecisionConfidenceChangePlanReadyForExternalApplyReview
	case DecisionConfidenceChangePlanRejected:
		disposition.Status = DecisionConfidenceChangePlanAbort
	case DecisionConfidenceChangePlanUnknown:
		disposition.Status = DecisionConfidenceChangePlanHold
	default:
		return DecisionConfidenceChangePlanDisposition{}, fmt.Errorf("unsupported verification status %q", verification.Status)
	}
	digest, err := Digest(disposition)
	if err != nil {
		return DecisionConfidenceChangePlanDisposition{}, fmt.Errorf("digest change plan disposition: %w", err)
	}
	disposition.DispositionDigest = digest
	if err := disposition.Validate(); err != nil {
		return DecisionConfidenceChangePlanDisposition{}, fmt.Errorf("validate change plan disposition: %w", err)
	}
	return disposition, nil
}

func (d DecisionConfidenceChangePlanDisposition) Validate() error {
	if d.ChangePlanDigest == "" || d.VerificationDigest == "" {
		return fmt.Errorf("plan and verification digests are required")
	}
	if d.Status != DecisionConfidenceChangePlanReadyForExternalApplyReview && d.Status != DecisionConfidenceChangePlanAbort && d.Status != DecisionConfidenceChangePlanHold {
		return fmt.Errorf("unsupported disposition status %q", d.Status)
	}
	if !d.NonExecuting {
		return fmt.Errorf("disposition must be non-executing")
	}
	if d.DispositionDigest == "" {
		return fmt.Errorf("disposition digest is required")
	}
	withoutDigest := d
	withoutDigest.DispositionDigest = ""
	digest, err := Digest(withoutDigest)
	if err != nil {
		return fmt.Errorf("digest change plan disposition: %w", err)
	}
	if digest != d.DispositionDigest {
		return fmt.Errorf("disposition digest mismatch")
	}
	return nil
}
