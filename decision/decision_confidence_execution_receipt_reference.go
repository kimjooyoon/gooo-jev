package decision

import "fmt"

// DecisionConfidenceExecutionReceiptReference binds a disposition outcome
// to an existing execution receipt by digest without issuing a receipt.
type DecisionConfidenceExecutionReceiptReference struct {
	OutcomeObservationDigest string
	ExecutionReceiptDigest  string
	NonAuthorizing          bool
	ReferenceDigest         string
}

func ReferenceDecisionConfidenceExecutionReceipt(
	observation DecisionConfidenceChangePlanDispositionObservation,
	receipt ExecutionReceipt,
) (DecisionConfidenceExecutionReceiptReference, error) {
	if err := observation.Validate(); err != nil {
		return DecisionConfidenceExecutionReceiptReference{}, fmt.Errorf("validate disposition observation: %w", err)
	}
	receiptDigest, err := Digest(receipt)
	if err != nil {
		return DecisionConfidenceExecutionReceiptReference{}, fmt.Errorf("digest execution receipt: %w", err)
	}
	reference := DecisionConfidenceExecutionReceiptReference{
		OutcomeObservationDigest: observation.ObservationDigest,
		ExecutionReceiptDigest:  receiptDigest,
		NonAuthorizing:          true,
	}
	digest, err := Digest(reference)
	if err != nil {
		return DecisionConfidenceExecutionReceiptReference{}, fmt.Errorf("digest receipt reference: %w", err)
	}
	reference.ReferenceDigest = digest
	if err := reference.Validate(); err != nil {
		return DecisionConfidenceExecutionReceiptReference{}, fmt.Errorf("validate receipt reference: %w", err)
	}
	return reference, nil
}

func (r DecisionConfidenceExecutionReceiptReference) Validate() error {
	if r.OutcomeObservationDigest == "" || r.ExecutionReceiptDigest == "" {
		return fmt.Errorf("observation and receipt digests are required")
	}
	if !r.NonAuthorizing {
		return fmt.Errorf("receipt reference must be non-authorizing")
	}
	if r.ReferenceDigest == "" {
		return fmt.Errorf("reference digest is required")
	}
	withoutDigest := r
	withoutDigest.ReferenceDigest = ""
	digest, err := Digest(withoutDigest)
	if err != nil {
		return fmt.Errorf("digest receipt reference: %w", err)
	}
	if digest != r.ReferenceDigest {
		return fmt.Errorf("reference digest mismatch")
	}
	return nil
}
