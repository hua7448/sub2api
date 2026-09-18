package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// openAIStreamFailedEventClientPayload preserves the local client-facing
// response.failed envelope while relying on the current upstream error parser.
func openAIStreamFailedEventClientPayload(payload []byte, failedMessage string) []byte {
	body := openAIStreamFailedEventPassthroughBody(payload, failedMessage)
	errorPayload := gjson.GetBytes(body, "error")
	if !errorPayload.Exists() || !errorPayload.IsObject() {
		return payload
	}
	event, err := sjson.SetRawBytes([]byte(`{"type":"response.failed"}`), "error", []byte(errorPayload.Raw))
	if err != nil {
		return payload
	}
	return event
}

func normalizeOpenAIResponsesUpstreamBody(body []byte, modelCandidates ...string) ([]byte, bool, error) {
	if len(bytes.TrimSpace(body)) == 0 {
		return body, false, nil
	}
	modified := false
	var err error
	model := strings.TrimSpace(gjson.GetBytes(body, "model").String())
	preserveMaxEffort := shouldPreserveOpenAIMaxReasoningEffort(append(modelCandidates, model)...)
	if !preserveMaxEffort && strings.EqualFold(strings.TrimSpace(gjson.GetBytes(body, "reasoning.effort").String()), "max") {
		body, err = sjson.SetBytes(body, "reasoning.effort", "xhigh")
		if err != nil {
			return nil, false, err
		}
		modified = true
	}
	if !preserveMaxEffort && strings.EqualFold(strings.TrimSpace(gjson.GetBytes(body, "reasoning_effort").String()), "max") {
		body, err = sjson.SetBytes(body, "reasoning_effort", "xhigh")
		if err != nil {
			return nil, false, err
		}
		modified = true
	}
	if bytes.Contains(body, []byte(`"namespace"`)) && gjson.GetBytes(body, "input").IsArray() {
		var removed bool
		body, removed, err = removeOpenAIResponsesInputNamespaces(body)
		if err != nil {
			return nil, false, err
		}
		modified = modified || removed
	}
	return body, modified, nil
}

func shouldPreserveOpenAIMaxReasoningEffort(modelCandidates ...string) bool {
	for _, model := range modelCandidates {
		if isOpenAIGPT56Model(model) {
			return true
		}
	}
	return false
}

func removeOpenAIResponsesInputNamespaces(body []byte) ([]byte, bool, error) {
	var requestBody map[string]any
	if err := json.Unmarshal(body, &requestBody); err != nil {
		return nil, false, fmt.Errorf("parse request: %w", err)
	}
	input, _ := requestBody["input"].([]any)
	removed := false
	for _, item := range input {
		obj, _ := item.(map[string]any)
		if obj == nil {
			continue
		}
		if _, ok := obj["namespace"]; ok {
			delete(obj, "namespace")
			removed = true
		}
	}
	if !removed {
		return body, false, nil
	}
	out, err := json.Marshal(requestBody)
	if err != nil {
		return nil, false, fmt.Errorf("encode request: %w", err)
	}
	return out, true, nil
}
