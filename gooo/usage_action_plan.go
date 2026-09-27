package gooo

import (
	"fmt"
	"strings"
)

type UsageActionState string

const (
	UsageActionReady    UsageActionState = "READY"
	UsageActionDeferred UsageActionState = "DEFERRED"
)

type UsageAction struct {
	ID             string           `json:"id"`
	CapabilityID   string           `json:"capability_id"`
	State          UsageActionState `json:"state"`
	Operation      string           `json:"operation"`
	Description    string           `json:"description"`
	Reason         string           `json:"reason"`
	SourceDigest   string           `json:"source_digest"`
	DiscoveryDigest string          `json:"discovery_digest"`
	NonExecuting   bool             `json:"non_executing"`
	NonAuthorizing bool             `json:"non_authorizing"`
}

type UsageActionPlan struct {
	SourceDigest    string        `json:"source_digest"`
	DiscoveryDigest string        `json:"discovery_digest"`
	Actions         []UsageAction `json:"actions"`
	PlanDigest      string        `json:"plan_digest"`
	NonExecuting    bool          `json:"non_executing"`
	NonAuthorizing  bool          `json:"non_authorizing"`
}

// PlanUsageActions turns capability evidence into non-executing next steps.
func PlanUsageActions(discovery UsageDiscoveryResponse) (UsageActionPlan, error) {
	if err := discovery.Validate(); err != nil {
		return UsageActionPlan{}, err
	}
	plan := UsageActionPlan{
		SourceDigest:    discovery.SourceDigest,
		DiscoveryDigest: discovery.DiscoveryDigest,
		NonExecuting:    true,
		NonAuthorizing:  true,
	}
	for _, capability := range discovery.Capabilities {
		state := UsageActionDeferred
		reason := capability.Reason
		if capability.State == CapabilityAvailable {
			state = UsageActionReady
			reason = "capability_available"
		}
		plan.Actions = append(plan.Actions, UsageAction{
			ID:              capability.ID + ":" + capability.NextOperation,
			CapabilityID:    capability.ID,
			State:           state,
			Operation:       capability.NextOperation,
			Description:     capability.Description,
			Reason:          reason,
			SourceDigest:    discovery.SourceDigest,
			DiscoveryDigest: discovery.DiscoveryDigest,
			NonExecuting:    true,
			NonAuthorizing:  true,
		})
	}
	plan.PlanDigest = digestUsageActionPlan(plan)
	return plan, nil
}

func (plan UsageActionPlan) Validate() error {
	if !validDigest(plan.SourceDigest) || !validDigest(plan.DiscoveryDigest) || !validDigest(plan.PlanDigest) {
		return fmt.Errorf("usage action plan digests are invalid")
	}
	if len(plan.Actions) == 0 || !plan.NonExecuting || !plan.NonAuthorizing {
		return fmt.Errorf("usage action plan is incomplete or crossed a capability boundary")
	}
	seen := map[string]struct{}{}
	for _, action := range plan.Actions {
		if strings.TrimSpace(action.ID) == "" || strings.TrimSpace(action.CapabilityID) == "" ||
			strings.TrimSpace(action.Operation) == "" || strings.TrimSpace(action.Description) == "" ||
			strings.TrimSpace(action.Reason) == "" {
			return fmt.Errorf("usage action is incomplete")
		}
		if action.State != UsageActionReady && action.State != UsageActionDeferred {
			return fmt.Errorf("usage action state is invalid")
		}
		if action.SourceDigest != plan.SourceDigest || action.DiscoveryDigest != plan.DiscoveryDigest {
			return fmt.Errorf("usage action provenance is not plan-bound")
		}
		if !action.NonExecuting || !action.NonAuthorizing {
			return fmt.Errorf("usage action crossed a capability boundary")
		}
		if _, exists := seen[action.ID]; exists {
			return fmt.Errorf("usage action %q is duplicated", action.ID)
		}
		seen[action.ID] = struct{}{}
	}
	if digestUsageActionPlan(plan) != plan.PlanDigest {
		return fmt.Errorf("usage action plan digest does not match its evidence")
	}
	return nil
}

func digestUsageActionPlan(plan UsageActionPlan) string {
	value := plan.SourceDigest + "|" + plan.DiscoveryDigest
	for _, action := range plan.Actions {
		value += "|" + action.ID + "|" + string(action.State) + "|" + action.Operation + "|" + action.Description + "|" + action.Reason
	}
	return digestString(value)
}

