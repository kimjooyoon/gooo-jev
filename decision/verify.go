package decision

import (
	"errors"
	"fmt"
)

func (receipt Receipt) Verify(spec Spec, state State, result Result) error {
	if err := receipt.Validate(); err != nil {
		return err
	}
	expected, err := Observe(spec, state, result)
	if err != nil {
		return err
	}
	if receipt.SpecDigest != expected.SpecDigest {
		return errors.New("decision receipt spec digest does not match")
	}
	if receipt.StateDigest != expected.StateDigest {
		return errors.New("decision receipt state digest does not match")
	}
	if receipt.ResultDigest != expected.ResultDigest {
		return errors.New("decision receipt result digest does not match")
	}
	if receipt.PolicyDigest != expected.PolicyDigest {
		return errors.New("decision receipt policy digest does not match")
	}
	if receipt.Provider != expected.Provider {
		return errors.New("decision receipt provider does not match")
	}
	if receipt.Status != expected.Status {
		return errors.New("decision receipt status does not match")
	}
	if !receipt.ObservedAt.Equal(expected.ObservedAt) {
		return errors.New("decision receipt observation time does not match")
	}
	if receipt.DecisionDigest != expected.DecisionDigest {
		return fmt.Errorf("decision receipt digest does not match")
	}
	return nil
}
