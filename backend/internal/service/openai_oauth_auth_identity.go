package service

import "context"

type codexAccountAuthIdentityContextKey struct{}

// WithCodexAccountAuthIdentity carries only the administrator-configured client
// metadata for an account refresh. Account IDs and credentials are never part
// of this identity; endpoint selection is unaffected.
func WithCodexAccountAuthIdentity(ctx context.Context, userAgent string) context.Context {
	return context.WithValue(ctx, codexAccountAuthIdentityContextKey{}, userAgent)
}

// CodexAuthIdentityFromContext uses the same UA/originator pairing and current
// version as inference. Requests without an account retain canonical defaults.
// As on the existing authentication path, no inference-only version header is
// introduced.
func CodexAuthIdentityFromContext(ctx context.Context) (userAgent, originator string) {
	candidate, _ := ctx.Value(codexAccountAuthIdentityContextKey{}).(string)
	identity := resolveCodexOutboundIdentity(candidate)
	return identity.userAgent, identity.originator
}
