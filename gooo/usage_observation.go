package gooo

import (
	"errors"
	"math"
	"strings"
	"time"
)

type UsageObservationKind string

const (
	UsageObservationConfirmed UsageObservationKind = "confirmed"
	UsageObservationRefuted   UsageObservationKind = "refuted"
	UsageObservationUnknown   UsageObservationKind = "unknown"
)

type UsageObservation struct {
	SourceDigest      string                 `json:"source_digest"`
	DiscoveryDigest   string                 `json:"discovery_digest"`
	PlanDigest        string                 `json:"plan_digest"`
	ActionID          string                 `json:"action_id"`
	CapabilityID     string                 `json:"capability_id"`
	Kind             UsageObservationKind   `json:"kind"`
	MetricName       string                 `json:"metric_name"`
	MetricValue      float64                `json:"metric_value"`
	EvidenceDigest   string                 `json:"evidence_digest"`
	RecordedAt       time.Time              `json:"recorded_at"`
	NonAuthorizing   bool                   `json:"non_authorizing"`
	ObservationDigest string                `json:"observation_digest"`
}

// ObserveUsageAction records an evidence-bound result without executing or authorizing an action.
// Deferred actions may only produce UNKNOWN observations.
func ObserveUsageAction(plan UsageActionPlan, actionID string, kind UsageObservationKind, metricName string, metricValue float64, evidenceDigest string, recordedAt time.Time) (UsageObservation, error) {
	if err := plan.Validate(); err != nil {
		return UsageObservation{}, err
	}
	action, found := findUsageAction(plan, actionID)
	if !found {
		return UsageObservation{}, errors.New("usage observation action is not present in plan")
	}
	if action.State == UsageActionDeferred && kind != UsageObservationUnknown {
		return UsageObservation{}, errors.New("deferred usage action requires unknown observation")
	}
	observation := UsageObservation{
		SourceDigest:    plan.SourceDigest,
		DiscoveryDigest: plan.DiscoveryDigest,
		PlanDigest:      plan.PlanDigest,
		ActionID:        action.ID,
		CapabilityID:    action.CapabilityID,
		Kind:            kind,
		MetricName:      metricName,
		MetricValue:     metricValue,
		EvidenceDigest:  evidenceDigest,
		RecordedAt:      recordedAt.UTC(),
		NonAuthorizing:  true,
	}
	observation.ObservationDigest = digestUsageObservation(observation)
	if err := observation.ValidateAgainst(plan); err != nil {
		return UsageObservation{}, err
	}
	return observation, nil
}

func (observation UsageObservation) Validate() error {
	if !validDigest(observation.SourceDigest) || !validDigest(observation.DiscoveryDigest) ||
		!validDigest(observation.PlanDigest) || !validDigest(observation.ObservationDigest) ||
		strings.TrimSpace(observation.ActionID) == "" || strings.TrimSpace(observation.CapabilityID) == "" ||
		strings.TrimSpace(observation.MetricName) == "" || strings.TrimSpace(observation.EvidenceDigest) == "" {
		return errors.New("usage observation is incomplete")
	}
	if math.IsNaN(observation.MetricValue) || math.IsInf(observation.MetricValue, 0) {
		return errors.New("usage observation metric value must be finite")
	}
	if observation.RecordedAt.IsZero() {
		return errors.New("usage observation recorded time is required")
	}
	if !observation.NonAuthorizing {
		return errors.New("usage observation must remain non-authorizing")
	}
	switch observation.Kind {
	case UsageObservationConfirmed, UsageObservationRefuted, UsageObservationUnknown:
	default:
		return errors.New("unsupported usage observation kind")
	}
	if digestUsageObservation(observation) != observation.ObservationDigest {
		return errors.New("usage observation digest does not match its evidence")
	}
	return nil
}

func (observation UsageObservation) ValidateAgainst(plan UsageActionPlan) error {
	if err := plan.Validate(); err != nil {
		return err
	}
	if err := observation.Validate(); err != nil {
		return err
	}
	if observation.SourceDigest != plan.SourceDigest || observation.DiscoveryDigest != plan.DiscoveryDigest ||
		observation.PlanDigest != plan.PlanDigest {
		return errors.New("usage observation provenance is not plan-bound")
	}
	action, found := findUsageAction(plan, observation.ActionID)
	if !found || action.CapabilityID != observation.CapabilityID {
		return errors.New("usage observation action is not plan-bound")
	}
	if action.State == UsageActionDeferred && observation.Kind != UsageObservationUnknown {
		return errors.New("deferred usage action requires unknown observation")
	}
	return nil
}

func findUsageAction(plan UsageActionPlan, actionID string) (UsageAction, bool) {
	for _, action := range plan.Actions {
		if action.ID == actionID {
			return action, true
		}
	}
	return UsageAction{}, false
}

func digestUsageObservation(observation UsageObservation) string {
	return digestString(observation.SourceDigest + "|" + observation.DiscoveryDigest + "|" +
		observation.PlanDigest + "|" + observation.ActionID + "|" + observation.CapabilityID + "|" +
		string(observation.Kind) + "|" + observation.MetricName + "|" +
		formatUsageObservationMetric(observation.MetricValue) + "|" + observation.EvidenceDigest + "|" +
		observation.RecordedAt.UTC().Format(time.RFC3339Nano))
}

