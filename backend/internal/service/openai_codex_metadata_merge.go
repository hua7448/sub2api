package service

import (
	"encoding/json"
	"net/http"
	"strings"
)

// Keep client-provided metadata from either transport carrier. CPA preserves
// opaque turn metadata while rewriting identity; do the same without replacing
// arbitrary occurrences of a session string in unrelated fields.
func seedCodexFingerprintHeaderMetadata(ids *codexFingerprintIDs, headers http.Header) {
	if ids != nil {
		ids.headerTurnMetadata = headers.Get("x-codex-turn-metadata")
	}
}

func decodeCodexTurnMetadata(raw string) map[string]any {
	var metadata map[string]any
	if json.Valid([]byte(raw)) {
		decoder := json.NewDecoder(strings.NewReader(raw))
		decoder.UseNumber()
		if decoder.Decode(&metadata) == nil && metadata != nil {
			return metadata
		}
	}
	return make(map[string]any)
}

// A cache key matching a declared session is an identity alias. Leave explicit
// independent cache keys untouched and never introduce a missing cache key.
func rewriteCodexMetadataDefaultCacheKey(metadata, fields map[string]any) {
	session, ok := fields["session_id"].(string)
	if !ok || session == "" {
		return
	}
	cache, ok := metadata["prompt_cache_key"].(string)
	if !ok || cache == "" {
		return
	}
	for _, key := range []string{"session_id", "session-id", "x-codex-session-id", "conversation_id", "conversation-id"} {
		if original, ok := metadata[key].(string); ok && cache == original {
			metadata["prompt_cache_key"] = session
			return
		}
	}
}

// Normalize each carrier before merging. A fallback cache alias refers to its
// own original session, which may differ from the primary carrier's session.
func mergeCodexFingerprintTurnMetadata(primary, fallback string, fields map[string]any) string {
	merged := decodeCodexTurnMetadata(primary)
	defaults := decodeCodexTurnMetadata(fallback)
	rewriteCodexMetadataDefaultCacheKey(merged, fields)
	rewriteCodexMetadataDefaultCacheKey(defaults, fields)
	for key, value := range defaults {
		if _, exists := merged[key]; !exists {
			merged[key] = value
		}
	}
	for key, value := range fields {
		merged[key] = value
	}
	alignCodexFingerprintMetadataAliases(merged, fields)
	encoded, err := marshalCodexASCIIJSON(merged)
	if err != nil {
		return primary
	}
	return string(encoded)
}
