package gooo

import (
	"fmt"
	"strings"
)

type UsageAssistanceItem struct {
	ID          string           `json:"id"`
	State       UsageActionState `json:"state"`
	Description string           `json:"description"`
	Operation   string           `json:"operation"`
	Reason      string           `json:"reason"`
}

type UsageAssistance struct {
	SourceDigest     string               `json:"source_digest"`
	DiscoveryDigest  string               `json:"discovery_digest"`
	PlanDigest       string               `json:"plan_digest"`
	Summary          string               `json:"summary"`
	Items            []UsageAssistanceItem `json:"items"`
	NextOperations   []string             `json:"next_operations"`
	AssistanceDigest string               `json:"assistance_digest"`
	NonExecuting     bool                 `json:"non_executing"`
	NonAuthorizing   bool                 `json:"non_authorizing"`
}

// ExplainUsage turns capability and action evidence into bounded, human-readable guidance.
func ExplainUsage(discovery UsageDiscoveryResponse, plan UsageActionPlan) (UsageAssistance, error) {
	if err := discovery.Validate(); err != nil {
		return UsageAssistance{}, err
	}
	if err := plan.Validate(); err != nil {
		return UsageAssistance{}, err
	}
	if plan.SourceDigest != discovery.SourceDigest || plan.DiscoveryDigest != discovery.DiscoveryDigest {
		return UsageAssistance{}, fmt.Errorf("usage assistance provenance is not discovery-bound")
	}

	assistance := UsageAssistance{
		SourceDigest:    discovery.SourceDigest,
		DiscoveryDigest: discovery.DiscoveryDigest,
		PlanDigest:      plan.PlanDigest,
		NonExecuting:    true,
		NonAuthorizing:  true,
	}
	var readyIDs []string
	var deferredIDs []string
	appendUnique := func(values *[]string, value string) {
		for _, existing := range *values {
			if existing == value {
				return
			}
		}
		*values = append(*values, value)
	}
	for _, action := range plan.Actions {
		assistance.Items = append(assistance.Items, UsageAssistanceItem{
			ID:          action.ID,
			State:       action.State,
			Description: action.Description,
			Operation:   action.Operation,
			Reason:      action.Reason,
		})
		appendUnique(&assistance.NextOperations, action.Operation)
		if action.State == UsageActionReady {
			appendUnique(&readyIDs, action.CapabilityID)
		} else {
			appendUnique(&deferredIDs, action.CapabilityID)
		}
	}
	switch {
	case len(readyIDs) > 0 && len(deferredIDs) > 0:
		assistance.Summary = fmt.Sprintf(
			"gooo can currently help with %s; it keeps %s deferred until their evidence boundary is bound.",
			strings.Join(readyIDs, ", "),
			strings.Join(deferredIDs, ", "),
		)
	case len(readyIDs) > 0:
		assistance.Summary = fmt.Sprintf(
			"gooo can currently help with %s.",
			strings.Join(readyIDs, ", "),
		)
	default:
		assistance.Summary = "gooo has no bound capability to claim yet; complete the declaration before taking the next step."
	}
	assistance.AssistanceDigest = digestUsageAssistance(assistance)
	return assistance, nil
}

func (assistance UsageAssistance) Validate() error {
	if !validDigest(assistance.SourceDigest) ||
		!validDigest(assistance.DiscoveryDigest) ||
		!validDigest(assistance.PlanDigest) ||
		!validDigest(assistance.AssistanceDigest) {
		return fmt.Errorf("usage assistance digests are invalid")
	}
	if strings.TrimSpace(assistance.Summary) == "" ||
		len(assistance.Items) == 0 ||
		len(assistance.NextOperations) == 0 ||
		!assistance.NonExecuting ||
		!assistance.NonAuthorizing {
		return fmt.Errorf("usage assistance is incomplete or crossed a capability boundary")
	}
	seen := map[string]struct{}{}
	for _, item := range assistance.Items {
		if strings.TrimSpace(item.ID) == "" ||
			strings.TrimSpace(item.Description) == "" ||
			strings.TrimSpace(item.Operation) == "" ||
			strings.TrimSpace(item.Reason) == "" {
			return fmt.Errorf("usage assistance item is incomplete")
		}
		if item.State != UsageActionReady && item.State != UsageActionDeferred {
			return fmt.Errorf("usage assistance item state is invalid")
		}
		if _, exists := seen[item.ID]; exists {
			return fmt.Errorf("usage assistance item %q is duplicated", item.ID)
		}
		seen[item.ID] = struct{}{}
	}
	for _, operation := range assistance.NextOperations {
		if strings.TrimSpace(operation) == "" {
			return fmt.Errorf("usage assistance next operation is empty")
		}
	}
	if digestUsageAssistance(assistance) != assistance.AssistanceDigest {
		return fmt.Errorf("usage assistance digest does not match its evidence")
	}
	return nil
}

func RenderUsageAssistance(assistance UsageAssistance) (string, error) {
	if err := assistance.Validate(); err != nil {
		return "", err
	}
	var builder strings.Builder
	builder.WriteString("gooo capability discovery\n")
	builder.WriteString(assistance.Summary)
	builder.WriteString("\n\n")
	for _, item := range assistance.Items {
		fmt.Fprintf(
			&builder,
			"- %s [%s]: %s (next: %s; reason: %s)\n",
			item.ID,
			item.State,
			item.Description,
			item.Operation,
			item.Reason,
		)
	}
	fmt.Fprintf(
		&builder,
		"\nnext operations: %s\nsource digest: %s\ndiscovery digest: %s\nplan digest: %s\nassistance digest: %s\nnon-executing: true\nnon-authorizing: true\n",
		strings.Join(assistance.NextOperations, ", "),
		assistance.SourceDigest,
		assistance.DiscoveryDigest,
		assistance.PlanDigest,
		assistance.AssistanceDigest,
	)
	return builder.String(), nil
}

func digestUsageAssistance(assistance UsageAssistance) string {
	value := assistance.SourceDigest + "|" + assistance.DiscoveryDigest + "|" + assistance.PlanDigest + "|" + assistance.Summary
	for _, item := range assistance.Items {
		value += "|" + item.ID + "|" + string(item.State) + "|" + item.Description + "|" + item.Operation + "|" + item.Reason
	}
	for _, operation := range assistance.NextOperations {
		value += "|" + operation
	}
	return digestString(value)
}
