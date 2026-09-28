package gooo

import (
	"fmt"
	"sort"
	"strings"
)

type CapabilityQueryActionKind string

const (
	CapabilityQueryActionInspectNextOperation    CapabilityQueryActionKind = "INSPECT_NEXT_OPERATION"
	CapabilityQueryActionAskClarifyingQuestion   CapabilityQueryActionKind = "ASK_CLARIFYING_QUESTION"
	CapabilityQueryActionRequireExternalBoundary CapabilityQueryActionKind = "REQUIRE_EXTERNAL_BOUNDARY"
)

// CapabilityQueryAction is a non-executing LSP action derived from a
// capability-query hover. It tells an editor what to inspect next without
// invoking a provider, mutating a declaration, or granting authority.
type CapabilityQueryAction struct {
	Title          string                   `json:"title"`
	Kind           CapabilityQueryActionKind `json:"kind"`
	Status         CapabilityQueryState     `json:"status"`
	CapabilityID   string                   `json:"capability_id,omitempty"`
	NextOperation  string                   `json:"next_operation"`
	Question       string                   `json:"question"`
	EvidenceDigest string                   `json:"evidence_digest"`
	NonExecuting   bool                     `json:"non_executing"`
	NonAuthorizing bool                     `json:"non_authorizing"`
}

func SuggestCapabilityQueryActions(hover CapabilityQueryHover) ([]CapabilityQueryAction, error) {
	if err := hover.Validate(); err != nil {
		return nil, fmt.Errorf("capability query hover: %w", err)
	}
	actions := make([]CapabilityQueryAction, 0, len(hover.Capabilities))
	switch hover.Status {
	case CapabilityQueryUnknown:
		actions = append(actions, CapabilityQueryAction{
			Title:          "Ask a narrower gooo capability question",
			Kind:           CapabilityQueryActionAskClarifyingQuestion,
			Status:         CapabilityQueryUnknown,
			NextOperation:  "ask_clarifying_question",
			Question:       firstCapabilityQueryActionQuestion(hover),
			NonExecuting:   true,
			NonAuthorizing: true,
		})
	case CapabilityQueryDeferred:
		question := firstCapabilityQueryActionQuestion(hover)
		if question == "" {
			question = capabilityQueryGuideExternalBoundaryQuestion
		}
		for _, capability := range hover.Capabilities {
			actions = append(actions, CapabilityQueryAction{
				Title:          fmt.Sprintf("Inspect the external boundary for %s", capability.ID),
				Kind:           CapabilityQueryActionRequireExternalBoundary,
				Status:         CapabilityQueryDeferred,
				CapabilityID:   capability.ID,
				NextOperation:  "provide_explicit_external_boundary",
				Question:       question,
				NonExecuting:   true,
				NonAuthorizing: true,
			})
		}
	case CapabilityQueryAvailable:
		for _, capability := range hover.Capabilities {
			actions = append(actions, CapabilityQueryAction{
				Title:          fmt.Sprintf("Inspect the next operation for %s", capability.ID),
				Kind:           CapabilityQueryActionInspectNextOperation,
				Status:         CapabilityQueryAvailable,
				CapabilityID:   capability.ID,
				NextOperation:  capability.NextOperation,
				Question:       capability.ExampleQuery,
				NonExecuting:   true,
				NonAuthorizing: true,
			})
		}
	default:
		return nil, fmt.Errorf("capability query action status %q is invalid", hover.Status)
	}
	sort.Slice(actions, func(i, j int) bool {
		if actions[i].CapabilityID != actions[j].CapabilityID {
			return actions[i].CapabilityID < actions[j].CapabilityID
		}
		return actions[i].Kind < actions[j].Kind
	})
	for index := range actions {
		actions[index].EvidenceDigest = digestCapabilityQueryAction(hover.EvidenceDigest, actions[index])
	}
	return actions, nil
}

func firstCapabilityQueryActionQuestion(hover CapabilityQueryHover) string {
	for _, question := range append(append([]string(nil), hover.NextQuestions...), hover.SuggestedQueries...) {
		if question = strings.TrimSpace(question); question != "" {
			return question
		}
	}
	return ""
}

func (action CapabilityQueryAction) Validate() error {
	if strings.TrimSpace(action.Title) == "" || strings.TrimSpace(action.NextOperation) == "" || strings.TrimSpace(action.Question) == "" {
		return fmt.Errorf("capability query action is incomplete")
	}
	if !validDigest(action.EvidenceDigest) {
		return fmt.Errorf("capability query action evidence digest is invalid")
	}
	if !action.NonExecuting || !action.NonAuthorizing {
		return fmt.Errorf("capability query action crossed an execution or authorization boundary")
	}
	switch action.Status {
	case CapabilityQueryUnknown:
		if action.Kind != CapabilityQueryActionAskClarifyingQuestion || action.CapabilityID != "" || action.NextOperation != "ask_clarifying_question" {
			return fmt.Errorf("unknown capability query action is incomplete")
		}
	case CapabilityQueryDeferred:
		if action.Kind != CapabilityQueryActionRequireExternalBoundary || strings.TrimSpace(action.CapabilityID) == "" || action.NextOperation != "provide_explicit_external_boundary" {
			return fmt.Errorf("deferred capability query action is incomplete")
		}
	case CapabilityQueryAvailable:
		if action.Kind != CapabilityQueryActionInspectNextOperation || strings.TrimSpace(action.CapabilityID) == "" {
			return fmt.Errorf("available capability query action is incomplete")
		}
	default:
		return fmt.Errorf("capability query action status %q is invalid", action.Status)
	}
	return nil
}

func ValidateCapabilityQueryActions(hover CapabilityQueryHover, actions []CapabilityQueryAction) error {
	if err := hover.Validate(); err != nil {
		return fmt.Errorf("capability query hover: %w", err)
	}
	if len(actions) == 0 {
		return fmt.Errorf("capability query actions are empty")
	}
	for _, action := range actions {
		if err := action.Validate(); err != nil {
			return err
		}
		if digestCapabilityQueryAction(hover.EvidenceDigest, action) != action.EvidenceDigest {
			return fmt.Errorf("capability query action evidence digest does not match")
		}
		if action.Status != hover.Status {
			return fmt.Errorf("capability query action status does not match hover")
		}
	}
	return nil
}

func digestCapabilityQueryAction(hoverDigest string, action CapabilityQueryAction) string {
	return digestString(strings.Join([]string{
		"gooo-capability-query-action",
		hoverDigest,
		action.Title,
		string(action.Kind),
		string(action.Status),
		action.CapabilityID,
		action.NextOperation,
		action.Question,
		fmt.Sprintf("%t|%t", action.NonExecuting, action.NonAuthorizing),
	}, "|"))
}
