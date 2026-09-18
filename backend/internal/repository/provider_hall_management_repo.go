package repository

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/apikey"
	"github.com/Wei-Shaw/sub2api/ent/group"
	"github.com/Wei-Shaw/sub2api/ent/providerhallconfig"
	"github.com/Wei-Shaw/sub2api/ent/providerhallprobekey"
	"github.com/Wei-Shaw/sub2api/ent/providerhalltarget"
	"github.com/Wei-Shaw/sub2api/ent/usagelog"
	"github.com/Wei-Shaw/sub2api/internal/pkg/openai"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/lib/pq"
)

func (r *providerHallRepository) ListAdminGroups(ctx context.Context, f service.ProviderHallGroupFilter) (*service.ProviderHallGroupPage, error) {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.PageSize < 1 || f.PageSize > 100 {
		f.PageSize = 20
	}
	q := r.client.Group.Query().Where(group.DeletedAtIsNil())
	if f.Search != "" {
		q.Where(group.NameContainsFold(f.Search))
	}
	if f.Platform != "" {
		q.Where(group.PlatformEQ(f.Platform))
	}
	listings, err := r.client.ProviderHallGroup.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	lm := map[int64]*ent.ProviderHallGroup{}
	listedIDs := []int64{}
	for _, l := range listings {
		lm[l.ID] = l
		if l.Listed {
			listedIDs = append(listedIDs, l.ID)
		}
	}
	if f.Listed == "true" {
		q.Where(group.IDIn(listedIDs...))
	}
	if f.Listed == "false" {
		q.Where(group.IDNotIn(listedIDs...))
	}
	total, err := q.Count(ctx)
	if err != nil {
		return nil, err
	}
	if f.Sort == "name_desc" {
		q.Order(ent.Desc(group.FieldName), ent.Asc(group.FieldID))
	} else {
		q.Order(ent.Asc(group.FieldName, group.FieldID))
	}
	gs, err := q.Limit(f.PageSize).Offset((f.Page - 1) * f.PageSize).All(ctx)
	if err != nil {
		return nil, err
	}
	ids := []int64{}
	for _, g := range gs {
		ids = append(ids, g.ID)
	}
	ts, err := r.client.ProviderHallTarget.Query().Where(providerhalltarget.GroupIDIn(ids...)).Order(ent.Asc(providerhalltarget.FieldProfileID)).All(ctx)
	if err != nil {
		return nil, err
	}
	ps, err := r.client.ProviderHallProfile.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	pm := map[int64]*ent.ProviderHallProfile{}
	for _, p := range ps {
		pm[p.ID] = p
	}
	cfg, err := r.client.ProviderHallConfig.Get(ctx, 1)
	if err != nil {
		return nil, err
	}
	out := &service.ProviderHallGroupPage{Items: []service.ProviderHallGroupSummary{}, Total: total, Page: f.Page, PageSize: f.PageSize}
	keyIDs := []int64{}
	for _, t := range ts {
		if t.ProbeKeyID != nil {
			keyIDs = append(keyIDs, *t.ProbeKeyID)
		}
	}
	keyRows, err := r.client.APIKey.Query().Where(apikey.IDIn(keyIDs...), apikey.DeletedAtIsNil()).All(ctx)
	if err != nil {
		return nil, err
	}
	keyStatuses := map[int64]string{}
	for _, k := range keyRows {
		keyStatuses[k.ID] = providerHallKeyOption(k, true).Status
		if cfg.OperatorUserID == nil || k.UserID != *cfg.OperatorUserID {
			keyStatuses[k.ID] = "inactive"
		}
	}
	recent, err := r.client.QueryContext(ctx, `SELECT DISTINCT ON (j.group_id,j.kind) j.group_id,j.kind,j.id,j.status,COALESCE(v.verdict,''),COALESCE(j.finished_at,j.created_at) FROM provider_hall_jobs j LEFT JOIN provider_hall_verifications v ON v.job_id=j.id WHERE j.group_id=ANY($1::bigint[]) ORDER BY j.group_id,j.kind,j.created_at DESC,j.id DESC`, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	latest := map[int64]map[string]*service.ProviderHallAdminResult{}
	for recent.Next() {
		var gid int64
		var kind string
		item := &service.ProviderHallAdminResult{}
		if err := recent.Scan(&gid, &kind, &item.JobID, &item.Status, &item.Verdict, &item.At); err != nil {
			_ = recent.Close()
			return nil, err
		}
		if latest[gid] == nil {
			latest[gid] = map[string]*service.ProviderHallAdminResult{}
		}
		latest[gid][kind] = item
	}
	if err := recent.Err(); err != nil {
		_ = recent.Close()
		return nil, err
	}
	_ = recent.Close()
	for _, g := range gs {
		v := service.ProviderHallGroupSummary{ProviderHallGroup: service.ProviderHallGroup{GroupID: g.ID}, Name: g.Name, Platform: g.Platform, Status: g.Status, Supported: g.Platform == service.PlatformOpenAI || g.Platform == service.PlatformComposite, Targets: []service.ProviderHallTarget{}, Models: []string{}}
		v.KeyStatuses = map[string]int{}
		v.LatestProbe = latest[g.ID]["probe"]
		v.LatestVerification = latest[g.ID]["verification"]
		if !v.Supported {
			v.Reason = "unsupported_platform"
		}
		if l := lm[g.ID]; l != nil {
			v.ProviderHallGroup = *providerHallGroupFromEnt(l)
		}
		for _, t := range ts {
			if t.GroupID != g.ID {
				continue
			}
			p := pm[t.ProfileID]
			status := "missing"
			if t.ProbeKeyID != nil {
				if s, ok := keyStatuses[*t.ProbeKeyID]; ok {
					status = s
				}
			}
			v.KeyStatuses[status]++
			v.Targets = append(v.Targets, service.ProviderHallTarget{ProviderHallVersion: service.ProviderHallVersion{Version: t.Version, UpdatedAt: t.UpdatedAt, UpdatedBy: t.UpdatedBy}, ID: t.ID, GroupID: g.ID, ProfileID: t.ProfileID, ProbeKeyID: t.ProbeKeyID, Enabled: t.Enabled, AutoScheduleEnabled: &t.AutoScheduleEnabled, ProbeIntervalSeconds: t.ProbeIntervalSeconds, VerificationIntervalSeconds: t.VerificationIntervalSeconds})
			if p != nil {
				v.Models = append(v.Models, p.Model)
				if t.Enabled && (v.EffectiveProfileID == nil || p.Model == cfg.DefaultModel && string(p.Protocol) == string(cfg.DefaultProtocol)) {
					id := p.ID
					v.EffectiveProfileID = &id
				}
			}
		}
		out.Items = append(out.Items, v)
	}
	return out, nil
}

func providerHallKeyOption(k *ent.APIKey, registered bool) service.ProviderHallProbeKeyOption {
	s := "unused"
	switch {
	case k.Status != service.StatusActive:
		s = "inactive"
	case k.ExpiresAt != nil && !k.ExpiresAt.After(time.Now()):
		s = "expired"
	case k.Quota > 0 && k.QuotaUsed >= k.Quota:
		s = "quota_exhausted"
	case registered:
		s = "registered"
	case k.LastUsedAt != nil || k.QuotaUsed != 0 || k.Usage5h != 0 || k.Usage1d != 0 || k.Usage7d != 0:
		s = "used"
	}
	gid := int64(0)
	if k.GroupID != nil {
		gid = *k.GroupID
	}
	return service.ProviderHallProbeKeyOption{ID: k.ID, Name: k.Name, GroupID: gid, OperatorUserID: k.UserID, Status: s, Registered: registered}
}

func (r *providerHallRepository) ListProbeKeys(ctx context.Context, groupID int64) ([]service.ProviderHallProbeKeyOption, error) {
	cfg, err := r.client.ProviderHallConfig.Get(ctx, 1)
	if err != nil {
		return nil, err
	}
	out := []service.ProviderHallProbeKeyOption{}
	if cfg.OperatorUserID == nil {
		return out, nil
	}
	ks, err := r.client.APIKey.Query().Where(apikey.UserIDEQ(*cfg.OperatorUserID), apikey.GroupIDEQ(groupID), apikey.DeletedAtIsNil()).Order(ent.Asc(apikey.FieldID)).All(ctx)
	if err != nil {
		return nil, err
	}
	regs, err := r.client.ProviderHallProbeKey.Query().Where(providerhallprobekey.OperatorUserIDEQ(*cfg.OperatorUserID), providerhallprobekey.GroupIDEQ(groupID)).All(ctx)
	if err != nil {
		return nil, err
	}
	rm := map[int64]bool{}
	for _, k := range regs {
		rm[k.ID] = true
	}
	for _, k := range ks {
		opt := providerHallKeyOption(k, rm[k.ID])
		if opt.Status == "unused" {
			used, e := r.client.UsageLog.Query().Where(usagelog.APIKeyIDEQ(k.ID)).Exist(ctx)
			if e != nil {
				return nil, e
			}
			if used {
				opt.Status = "used"
			}
		}
		out = append(out, opt)
	}
	return out, nil
}

// The config row serializes ensure operations with operator changes and retries.
func (r *providerHallRepository) EnsureProbeKey(ctx context.Context, groupID, profileID, actorID int64) (*service.ProviderHallProbeKeyOption, error) {
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()
	cfg, err := tx.ProviderHallConfig.Query().Where(providerhallconfig.IDEQ(1)).ForUpdate().Only(ctx)
	if err != nil {
		return nil, err
	}
	if cfg.OperatorUserID == nil {
		return nil, service.ErrProviderHallOperator
	}
	g, err := tx.Group.Query().Where(group.IDEQ(groupID), group.DeletedAtIsNil()).ForUpdate().Only(ctx)
	if err != nil {
		return nil, err
	}
	p, err := tx.ProviderHallProfile.Get(ctx, profileID)
	if err != nil {
		return nil, err
	}
	regs, err := tx.ProviderHallProbeKey.Query().Where(providerhallprobekey.OperatorUserIDEQ(*cfg.OperatorUserID), providerhallprobekey.GroupIDEQ(groupID)).All(ctx)
	if err != nil {
		return nil, err
	}
	ids := []int64{}
	for _, k := range regs {
		ids = append(ids, k.ID)
	}
	keys, err := tx.APIKey.Query().Where(apikey.IDIn(ids...), apikey.UserIDEQ(*cfg.OperatorUserID), apikey.GroupIDEQ(groupID), apikey.DeletedAtIsNil()).Order(ent.Asc(apikey.FieldID)).All(ctx)
	if err != nil {
		return nil, err
	}
	var selected *ent.APIKey
	for _, k := range keys {
		if providerHallKeyOption(k, true).Status == "registered" {
			selected = k
			break
		}
	}
	if selected == nil {
		b := make([]byte, 32)
		if _, err := rand.Read(b); err != nil {
			return nil, err
		}
		selected, err = tx.APIKey.Create().SetUserID(*cfg.OperatorUserID).SetGroupID(groupID).SetName("Provider Hall").SetKey("sk-" + hex.EncodeToString(b)).Save(ctx)
		if err != nil {
			return nil, err
		}
	}
	if err := providerHallValidateAndRegisterKey(ctx, tx.Client(), cfg, g, p, &selected.ID, actorID); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	opt := providerHallKeyOption(selected, true)
	return &opt, nil
}

func (r *providerHallRepository) ListModelCandidates(ctx context.Context, groupID int64) ([]service.ProviderHallModelCandidate, error) {
	g, err := r.client.Group.Query().Where(group.IDEQ(groupID), group.DeletedAtIsNil()).WithAccounts().Only(ctx)
	if err != nil {
		return nil, err
	}
	out := []service.ProviderHallModelCandidate{}
	seen := map[string]bool{}
	routes := NewCompositeModelRouteRepository(r.client)
	resolver := service.NewCompositeRouteResolver(routes)
	add := func(model, upstream, source string, accountID int64, active bool) error {
		if model == "" || strings.ContainsAny(model, "*?\r\n") {
			return nil
		}
		for _, protocol := range []string{"responses", "chat_completions", "messages"} {
			key := model + ":" + protocol + ":" + source + ":" + strconv.FormatInt(accountID, 10)
			if seen[key] {
				continue
			}
			seen[key] = true
			available := active && g.Status == service.StatusActive && (g.Platform == service.PlatformOpenAI || g.Platform == service.PlatformComposite) && (protocol != "messages" || g.AllowMessagesDispatch)
			routedModel := model
			if g.Platform == service.PlatformComposite {
				route, e := resolver.Resolve(ctx, groupID, model, protocol)
				if e != nil {
					return e
				}
				available = available && route.Matched && route.TargetPlatform == service.PlatformOpenAI
				routedModel = route.UpstreamModel
			}
			hasAccount := false
			for _, a := range g.Edges.Accounts {
				if a.DeletedAt == nil && a.Status == service.StatusActive && a.Platform == service.PlatformOpenAI && (accountID == 0 || accountID == a.ID) && accountEntityToService(a).IsModelSupported(routedModel) {
					hasAccount = true
					break
				}
			}
			available = available && hasAccount
			reason := ""
			if !available {
				reason = "route_or_account_unavailable"
			}
			out = append(out, service.ProviderHallModelCandidate{Model: model, Protocol: protocol, Source: source, AccountID: accountID, UpstreamModel: upstream, Available: available, Reason: reason})
		}
		return nil
	}
	for _, a := range g.Edges.Accounts {
		if a.DeletedAt != nil {
			continue
		}
		account := accountEntityToService(a)
		if a.Platform == service.PlatformOpenAI {
			for _, model := range openai.DefaultModels {
				if err := add(model.ID, account.GetMappedModel(model.ID), "platform", a.ID, a.Status == service.StatusActive && account.IsModelSupported(model.ID)); err != nil {
					return nil, err
				}
			}
		}
		for model, upstream := range account.GetModelMapping() {
			if err := add(model, upstream, "account_mapping", a.ID, a.Status == service.StatusActive && a.Platform == service.PlatformOpenAI); err != nil {
				return nil, err
			}
		}
	}
	if g.Platform == service.PlatformComposite {
		rs, e := routes.ListByGroup(ctx, groupID, false)
		if e != nil {
			return nil, e
		}
		for _, route := range rs {
			if route.MatchType == service.CompositeRouteMatchExact {
				if e := add(route.PublicModel, route.UpstreamModel, "composite_route", 0, true); e != nil {
					return nil, e
				}
			}
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Model != out[j].Model {
			return out[i].Model < out[j].Model
		}
		return out[i].Protocol < out[j].Protocol
	})
	return out, nil
}
