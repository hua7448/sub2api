package claude

import "testing"

func TestDefaultModels_ContainsCurrentModelsWithoutDuplicates(t *testing.T) {
	t.Parallel()

	seen := make(map[string]struct{}, len(DefaultModels))
	for _, model := range DefaultModels {
		if _, ok := seen[model.ID]; ok {
			t.Fatalf("duplicate model %q in DefaultModels", model.ID)
		}
		seen[model.ID] = struct{}{}
	}

	for _, id := range []string{"claude-fable-5", "claude-opus-4-8", "claude-opus-5", "claude-sonnet-5"} {
		if _, ok := seen[id]; !ok {
			t.Fatalf("expected model %q to be exposed in DefaultModels", id)
		}
	}
}
