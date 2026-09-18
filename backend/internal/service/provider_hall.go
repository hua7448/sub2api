package service

import (
	"context"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/shopspring/decimal"
)

var (
	ErrProviderHallNotFound      = infraerrors.NotFound("PROVIDER_HALL_NOT_FOUND", "provider hall resource not found")
	ErrProviderHallConflict      = infraerrors.Conflict("PROVIDER_HALL_VERSION_CONFLICT", "provider hall configuration changed; reload before saving")
	ErrProviderHallProfileExists = infraerrors.Conflict("PROVIDER_HALL_PROFILE_EXISTS", "model and protocol already configured")
	ErrProviderHallNotReady      = infraerrors.Conflict("PROVIDER_HALL_NOT_READY", "provider hall runtime is not available in this implementation batch")
)

type ProviderHallVersion struct {
	Version   int64     `json:"version"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy *int64    `json:"updated_by"`
}

type ProviderHallConfig struct {
	ProviderHallVersion
	CollectionEnabled   bool     `json:"collection_enabled"`
	DisplayEnabled      bool     `json:"display_enabled"`
	TasksEnabled        bool     `json:"tasks_enabled"`
	AutoScheduleEnabled *bool    `json:"auto_schedule_enabled"`
	DefaultModel        string   `json:"default_model"`
	DefaultProtocol     string   `json:"default_protocol"`
	DefaultRange        string   `json:"default_range"`
	GatewayOrigin       string   `json:"gateway_origin"`
	OperatorUserID      *int64   `json:"operator_user_id"`
	DailyBudget         string   `json:"daily_budget"`
	ExpectedNodes       []string `json:"expected_nodes"`
}

type ProviderHallGroup struct {
	ProviderHallVersion
	GroupID      int64  `json:"group_id"`
	Listed       bool   `json:"listed"`
	DisplayName  string `json:"display_name"`
	Description  string `json:"description"`
	DisplayOrder int    `json:"display_order"`
}

type ProviderHallProfile struct {
	Groups    []ProviderHallProfileGroup `json:"groups,omitempty"`
	IsDefault bool                       `json:"is_default"`
	ProviderHallVersion
	ID                   int64      `json:"id"`
	Model                string     `json:"model"`
	Protocol             string     `json:"protocol"`
	SupportsTools        bool       `json:"supports_tools"`
	OutputLimit          int        `json:"output_limit"`
	ModelAliases         []string   `json:"model_aliases"`
	ReferenceInputPrice  *string    `json:"reference_input_price"`
	ReferenceCachePrice  *string    `json:"reference_cache_price"`
	ReferenceCacheRate   *string    `json:"reference_cache_rate"`
	ReferenceConfirmedAt *time.Time `json:"reference_confirmed_at"`
}

type ProviderHallProfileGroup struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type ProviderHallRepository interface {
	GetConfig(context.Context) (*ProviderHallConfig, error)
	UpdateConfig(context.Context, ProviderHallConfig) (*ProviderHallConfig, error)
	GetGroup(context.Context, int64) (*ProviderHallGroup, error)
	SaveGroup(context.Context, ProviderHallGroup) (*ProviderHallGroup, error)
	ListProfiles(context.Context) ([]ProviderHallProfile, error)
	SaveProfile(context.Context, ProviderHallProfile) (*ProviderHallProfile, error)
	GetTargets(context.Context, int64) (*ProviderHallTargetSet, error)
	SaveTargets(context.Context, ProviderHallTargetSet) (*ProviderHallTargetSet, error)
	IsProbeKey(context.Context, int64, time.Time) (bool, error)
}

type ProviderHallGroupAccess interface {
	GetAvailableGroups(context.Context, int64) ([]Group, error)
}

// ProviderHallReadiness records which runtime pieces are wired into this
// build. A switch can only be turned on once its implementation exists.
type ProviderHallReadiness struct {
	Collection bool
	Tasks      bool
	Display    bool
}

type ProviderHallService struct {
	repo         ProviderHallRepository
	groups       GroupRepository
	users        UserRepository
	access       ProviderHallGroupAccess
	ready        ProviderHallReadiness
	jobs         ProviderHallJobControl
	accounts     AccountRepository
	modelFetcher *AccountTestService
	modelRoutes  CompositeModelRouteRepository
	// onConfigSaved runs after a successful config save so process caches
	// derived from it (public settings, injected index.html) refresh at once.
	onConfigSaved func()
}

// SetConfigSavedCallback registers the cache-invalidation hook.
func (s *ProviderHallService) SetConfigSavedCallback(fn func()) {
	if s != nil {
		s.onConfigSaved = fn
	}
}

// SetJobControl lets configuration writes cancel queued jobs and stale reports
// immediately. Without it the runner converges on its next tick.
func (s *ProviderHallService) SetJobControl(jobs ProviderHallJobControl) {
	if s != nil {
		s.jobs = jobs
	}
}

func NewProviderHallService(repo ProviderHallRepository, groups GroupRepository, users UserRepository, keys *APIKeyService) *ProviderHallService {
	return &ProviderHallService{repo: repo, groups: groups, users: users, access: keys}
}

func (s *ProviderHallService) SetReadiness(r ProviderHallReadiness) {
	if s != nil {
		s.ready = r
	}
}

// SetGroupAccessForTest swaps the permission reader; test-only seam for
// packages that cannot reach the unexported field.
func (s *ProviderHallService) SetGroupAccessForTest(access ProviderHallGroupAccess) {
	if s != nil {
		s.access = access
	}
}

// SetDisplayReady is flipped by the user query service provider once the
// user API and page are wired; only then may display_enabled be turned on.
func (s *ProviderHallService) SetDisplayReady(ready bool) {
	if s != nil {
		s.ready.Display = ready
	}
}

// DisplayEnabled reports the user hall switch. It backs the public setting
// provider_hall_enabled and the user route guard.
func (s *ProviderHallService) DisplayEnabled(ctx context.Context) bool {
	if s == nil || s.repo == nil {
		return false
	}
	cfg, err := s.repo.GetConfig(ctx)
	return err == nil && cfg != nil && cfg.DisplayEnabled
}

func (s *ProviderHallService) Readiness() ProviderHallReadiness {
	if s == nil {
		return ProviderHallReadiness{}
	}
	return s.ready
}

func providerHallNotReady(switchName string) error {
	return ErrProviderHallNotReady.WithMetadata(map[string]string{"switch": switchName})
}

func (s *ProviderHallService) GetConfig(ctx context.Context) (*ProviderHallConfig, error) {
	return s.repo.GetConfig(ctx)
}

func providerHallInvalid(field string) error {
	return infraerrors.BadRequest("PROVIDER_HALL_INVALID_CONFIG", "invalid provider hall configuration").WithMetadata(map[string]string{"field": field})
}

func providerHallProtocol(value string) bool {
	return value == "responses" || value == "chat_completions" || value == "messages"
}

var providerHallDecimalPattern = regexp.MustCompile(`^(0|[1-9][0-9]*)(\.[0-9]+)?$`)
var providerHallNodePattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_.:-]{0,127}$`)

// Validate precision before parsing: PostgreSQL numeric would otherwise silently
// round excess fractional digits. Amounts never pass through float64.
func providerHallDecimal(value string, precision, scale int) (decimal.Decimal, bool) {
	if len(value) > precision+2 || !providerHallDecimalPattern.MatchString(value) {
		return decimal.Zero, false
	}
	parts := strings.SplitN(value, ".", 2)
	if len(parts[0]) > precision-scale || (len(parts) == 2 && len(parts[1]) > scale) {
		return decimal.Zero, false
	}
	d, err := decimal.NewFromString(value)
	return d, err == nil
}

func (s *ProviderHallService) UpdateConfig(ctx context.Context, cfg ProviderHallConfig, actorID int64) (*ProviderHallConfig, error) {
	validated, err := s.ValidateConfig(ctx, cfg, actorID)
	if err != nil {
		return nil, err
	}
	saved, err := s.repo.UpdateConfig(ctx, *validated)
	if err == nil && s.onConfigSaved != nil {
		s.onConfigSaved()
	}
	return saved, err
}

func (s *ProviderHallService) ValidateConfig(ctx context.Context, cfg ProviderHallConfig, actorID int64) (*ProviderHallConfig, error) {
	if cfg.Version < 1 || actorID < 1 {
		return nil, providerHallInvalid("version")
	}
	if !providerHallProtocol(cfg.DefaultProtocol) {
		return nil, providerHallInvalid("default_protocol")
	}
	switch cfg.DefaultRange {
	case "6h", "24h", "7d", "30d":
	default:
		return nil, providerHallInvalid("default_range")
	}
	if cfg.DefaultModel != strings.TrimSpace(cfg.DefaultModel) || utf8.RuneCountInString(cfg.DefaultModel) > 200 {
		return nil, providerHallInvalid("default_model")
	}
	if _, ok := providerHallDecimal(cfg.DailyBudget, 20, 8); !ok {
		return nil, providerHallInvalid("daily_budget")
	}
	// A single trailing slash is what browsers and copy-paste produce for an
	// origin; accept it instead of treating it as a path.
	cfg.GatewayOrigin = strings.TrimSuffix(strings.TrimSpace(cfg.GatewayOrigin), "/")
	if cfg.GatewayOrigin != "" {
		u, err := url.Parse(cfg.GatewayOrigin)
		if err != nil || len(cfg.GatewayOrigin) > 2048 || u.Hostname() == "" || u.User != nil || u.Path != "" || u.RawQuery != "" || u.ForceQuery || u.Fragment != "" || u.Opaque != "" || strings.Contains(cfg.GatewayOrigin, "#") {
			return nil, providerHallInvalid("gateway_origin")
		}
		// Loopback origins are only for local end-to-end runs and must be
		// enabled explicitly; the runner enforces the same rule before sending.
		loopback := ProviderHallIsLoopbackOrigin(cfg.GatewayOrigin)
		if loopback && !ProviderHallAllowLoopback() || !loopback && u.Scheme != "https" || loopback && u.Scheme != "https" && u.Scheme != "http" {
			return nil, providerHallInvalid("gateway_origin")
		}
	}
	if len(cfg.ExpectedNodes) > 256 {
		return nil, providerHallInvalid("expected_nodes")
	}
	seen := make(map[string]bool)
	for _, node := range cfg.ExpectedNodes {
		if !providerHallNodePattern.MatchString(node) || seen[node] {
			return nil, providerHallInvalid("expected_nodes")
		}
		seen[node] = true
	}
	if cfg.ExpectedNodes == nil {
		cfg.ExpectedNodes = []string{}
	}
	// Each switch is gated by its own readiness flag, set by the wire provider
	// once that runtime piece exists. Saving settings cannot activate partial work.
	if cfg.CollectionEnabled && !s.ready.Collection {
		return nil, providerHallNotReady("collection_enabled")
	}
	if cfg.TasksEnabled {
		if !s.ready.Tasks {
			return nil, providerHallNotReady("tasks_enabled")
		}
		if !cfg.CollectionEnabled || cfg.GatewayOrigin == "" || cfg.OperatorUserID == nil {
			return nil, providerHallInvalid("tasks_enabled")
		}
	}
	if cfg.DisplayEnabled {
		if !s.ready.Display {
			return nil, providerHallNotReady("display_enabled")
		}
		if !cfg.CollectionEnabled {
			return nil, providerHallInvalid("display_enabled")
		}
	}
	if cfg.OperatorUserID != nil {
		if *cfg.OperatorUserID < 1 {
			return nil, providerHallInvalid("operator_user_id")
		}
		user, err := s.users.GetByID(ctx, *cfg.OperatorUserID)
		if err != nil {
			return nil, err
		}
		if !user.IsActive() {
			return nil, providerHallInvalid("operator_user_id")
		}
	}
	if cfg.DefaultModel != "" {
		profiles, err := s.repo.ListProfiles(ctx)
		if err != nil {
			return nil, err
		}
		found := false
		for _, p := range profiles {
			found = found || p.Model == cfg.DefaultModel && p.Protocol == cfg.DefaultProtocol
		}
		if !found {
			return nil, providerHallInvalid("default_model")
		}
	}
	cfg.UpdatedBy = &actorID
	return &cfg, nil
}

func (s *ProviderHallService) GetGroup(ctx context.Context, id int64) (*ProviderHallGroup, error) {
	g, err := s.groups.GetByIDLite(ctx, id)
	if err != nil {
		return nil, err
	}
	if g.Platform != PlatformOpenAI && g.Platform != PlatformComposite {
		return nil, providerHallInvalid("group_id")
	}
	return s.repo.GetGroup(ctx, id)
}

func (s *ProviderHallService) SaveGroup(ctx context.Context, cfg ProviderHallGroup, actorID int64) (*ProviderHallGroup, error) {
	if cfg.GroupID < 1 || cfg.Version < 0 || actorID < 1 {
		return nil, providerHallInvalid("version")
	}
	if utf8.RuneCountInString(cfg.DisplayName) > 100 || utf8.RuneCountInString(cfg.Description) > 2000 || cfg.DisplayOrder < -2147483648 || cfg.DisplayOrder > 2147483647 {
		return nil, providerHallInvalid("group")
	}
	g, err := s.groups.GetByIDLite(ctx, cfg.GroupID)
	if err != nil {
		return nil, err
	}
	if g.Platform != PlatformOpenAI && g.Platform != PlatformComposite || cfg.Listed && !g.IsActive() {
		return nil, providerHallInvalid("group_id")
	}
	cfg.DisplayName = strings.TrimSpace(cfg.DisplayName)
	cfg.UpdatedBy = &actorID
	return s.repo.SaveGroup(ctx, cfg)
}

func (s *ProviderHallService) ListProfiles(ctx context.Context) ([]ProviderHallProfile, error) {
	return s.repo.ListProfiles(ctx)
}

func (s *ProviderHallService) SaveProfile(ctx context.Context, p ProviderHallProfile, actorID int64) (*ProviderHallProfile, error) {
	if p.ID < 0 || actorID < 1 || (p.ID == 0 && p.Version != 0) || (p.ID > 0 && p.Version < 1) {
		return nil, providerHallInvalid("version")
	}
	if p.Model == "" || strings.TrimSpace(p.Model) != p.Model || utf8.RuneCountInString(p.Model) > 200 || !providerHallProtocol(p.Protocol) {
		return nil, providerHallInvalid("model_protocol")
	}
	if p.OutputLimit < 1 || p.OutputLimit > 1024 {
		return nil, providerHallInvalid("output_limit")
	}
	if len(p.ModelAliases) > 32 {
		return nil, providerHallInvalid("model_aliases")
	}
	seen := map[string]bool{p.Model: true}
	for _, alias := range p.ModelAliases {
		if alias == "" || strings.TrimSpace(alias) != alias || utf8.RuneCountInString(alias) > 200 || strings.ContainsAny(alias, "*?\r\n") || seen[alias] {
			return nil, providerHallInvalid("model_aliases")
		}
		seen[alias] = true
	}
	if p.ModelAliases == nil {
		p.ModelAliases = []string{}
	}
	if p.ReferenceInputPrice != nil || p.ReferenceCachePrice != nil || p.ReferenceCacheRate != nil || p.ReferenceConfirmedAt != nil {
		if p.ReferenceInputPrice == nil || p.ReferenceCachePrice == nil || p.ReferenceCacheRate == nil || p.ReferenceConfirmedAt == nil {
			return nil, providerHallInvalid("reference")
		}
		_, inputOK := providerHallDecimal(*p.ReferenceInputPrice, 24, 10)
		_, cacheOK := providerHallDecimal(*p.ReferenceCachePrice, 24, 10)
		rate, rateOK := providerHallDecimal(*p.ReferenceCacheRate, 11, 10)
		if !inputOK || !cacheOK || !rateOK || rate.GreaterThan(decimal.NewFromInt(1)) || p.ReferenceConfirmedAt.IsZero() || p.ReferenceConfirmedAt.After(time.Now()) {
			return nil, providerHallInvalid("reference")
		}
		confirmedAt := p.ReferenceConfirmedAt.UTC()
		p.ReferenceConfirmedAt = &confirmedAt
	}
	p.UpdatedBy = &actorID
	var previous *ProviderHallProfile
	if p.ID > 0 && s.jobs != nil {
		profiles, err := s.repo.ListProfiles(ctx)
		if err != nil {
			return nil, err
		}
		for i := range profiles {
			if profiles[i].ID == p.ID {
				previous = &profiles[i]
			}
		}
	}
	saved, err := s.repo.SaveProfile(ctx, p)
	if err != nil {
		return nil, err
	}
	// Reports depend on model/protocol/tools/aliases. The stale mark runs right
	// after the profile commit (the repository owns that transaction); the
	// runner also re-checks the fingerprint every tick, so a crash here only
	// delays the mark.
	if previous != nil && ProviderHallProfileFingerprint(previous.Model, previous.Protocol, previous.SupportsTools, previous.ModelAliases) !=
		ProviderHallProfileFingerprint(saved.Model, saved.Protocol, saved.SupportsTools, saved.ModelAliases) {
		if _, err := s.jobs.MarkVerificationsStale(ctx, saved.ID); err != nil {
			return saved, err
		}
	}
	return saved, nil
}

// AuthorizeUserGroup is shared by detail, trend and report reads. Never use the
// model-plaza permission API: it deliberately includes expired subscriptions.
func (s *ProviderHallService) AuthorizeUserGroup(ctx context.Context, userID, groupID int64) (*ProviderHallGroup, error) {
	cfg, err := s.repo.GetConfig(ctx)
	if err != nil {
		return nil, err
	}
	if !cfg.DisplayEnabled || userID < 1 || groupID < 1 {
		return nil, ErrProviderHallNotFound
	}
	groups, err := s.access.GetAvailableGroups(ctx, userID)
	if err != nil {
		return nil, err
	}
	for _, group := range groups {
		if group.ID != groupID || !group.IsActive() || (group.Platform != PlatformOpenAI && group.Platform != PlatformComposite) {
			continue
		}
		listing, err := s.repo.GetGroup(ctx, groupID)
		if err != nil {
			return nil, err
		}
		if listing.Listed {
			return listing, nil
		}
	}
	return nil, ErrProviderHallNotFound
}
