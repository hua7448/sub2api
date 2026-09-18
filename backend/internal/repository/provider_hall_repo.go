package repository

import (
	"context"
	"errors"
	"time"

	"github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/group"
	"github.com/Wei-Shaw/sub2api/ent/providerhallconfig"
	"github.com/Wei-Shaw/sub2api/ent/providerhallgroup"
	"github.com/Wei-Shaw/sub2api/ent/providerhallprofile"
	"github.com/Wei-Shaw/sub2api/ent/providerhalltarget"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/lib/pq"
	"github.com/shopspring/decimal"
)

type providerHallRepository struct{ client *ent.Client }

func NewProviderHallRepository(client *ent.Client) service.ProviderHallRepository {
	return &providerHallRepository{client: client}
}

func (r *providerHallRepository) GetConfig(ctx context.Context) (*service.ProviderHallConfig, error) {
	cfg, err := r.client.ProviderHallConfig.Get(ctx, 1)
	if err != nil {
		return nil, err
	}
	return providerHallConfigFromEnt(cfg), nil
}

func (r *providerHallRepository) UpdateConfig(ctx context.Context, cfg service.ProviderHallConfig) (*service.ProviderHallConfig, error) {
	budget, err := decimal.NewFromString(cfg.DailyBudget)
	if err != nil {
		return nil, err
	}
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()
	previous, err := tx.ProviderHallConfig.Query().Where(providerhallconfig.IDEQ(1)).ForUpdate().Only(ctx)
	if err != nil {
		return nil, err
	}
	if previous.Version != cfg.Version {
		return nil, service.ErrProviderHallConflict
	}
	operatorChanged := (previous.OperatorUserID == nil) != (cfg.OperatorUserID == nil) ||
		previous.OperatorUserID != nil && cfg.OperatorUserID != nil && *previous.OperatorUserID != *cfg.OperatorUserID
	if operatorChanged {
		var groupIDs []int64
		err := tx.ProviderHallTarget.Query().Where(providerhalltarget.EnabledEQ(true)).Select(providerhalltarget.FieldGroupID).Scan(ctx, &groupIDs)
		if err != nil {
			return nil, err
		}
		if len(groupIDs) > 0 {
			if _, err := tx.ProviderHallGroup.Update().Where(providerhallgroup.IDIn(groupIDs...)).AddVersion(1).SetNillableUpdatedBy(cfg.UpdatedBy).Save(ctx); err != nil {
				return nil, err
			}
			if _, err := tx.ProviderHallTarget.Update().Where(providerhalltarget.EnabledEQ(true)).SetEnabled(false).AddVersion(1).SetNillableUpdatedBy(cfg.UpdatedBy).Save(ctx); err != nil {
				return nil, err
			}
		}
	}
	u := tx.ProviderHallConfig.Update().Where(providerhallconfig.IDEQ(1), providerhallconfig.VersionEQ(cfg.Version)).
		AddVersion(1).SetUpdatedAt(time.Now().UTC()).SetNillableUpdatedBy(cfg.UpdatedBy).
		SetCollectionEnabled(cfg.CollectionEnabled).SetDisplayEnabled(cfg.DisplayEnabled).SetTasksEnabled(cfg.TasksEnabled).
		SetDefaultModel(cfg.DefaultModel).SetDefaultProtocol(providerhallconfig.DefaultProtocol(cfg.DefaultProtocol)).
		SetDefaultRange(providerhallconfig.DefaultRange(cfg.DefaultRange)).SetGatewayOrigin(cfg.GatewayOrigin).
		SetDailyBudget(budget).SetExpectedNodes(cfg.ExpectedNodes).ClearOperatorUserID()
	if cfg.AutoScheduleEnabled != nil {
		u.SetAutoScheduleEnabled(*cfg.AutoScheduleEnabled)
	}
	if cfg.OperatorUserID != nil {
		u.SetOperatorUserID(*cfg.OperatorUserID)
	}
	n, err := u.Save(ctx)
	if err != nil {
		return nil, err
	}
	if n != 1 {
		return nil, service.ErrProviderHallConflict
	}
	if cfg.DefaultModel != "" {
		exists, err := tx.ProviderHallProfile.Query().Where(providerhallprofile.ModelEQ(cfg.DefaultModel), providerhallprofile.ProtocolEQ(providerhallprofile.Protocol(cfg.DefaultProtocol))).Exist(ctx)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, service.ErrProviderHallConflict
		}
	}
	updated, err := tx.ProviderHallConfig.Get(ctx, 1)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return providerHallConfigFromEnt(updated), nil
}

func providerHallConfigFromEnt(cfg *ent.ProviderHallConfig) *service.ProviderHallConfig {
	nodes := append([]string{}, cfg.ExpectedNodes...)
	return &service.ProviderHallConfig{
		ProviderHallVersion: service.ProviderHallVersion{Version: cfg.Version, UpdatedAt: cfg.UpdatedAt.UTC(), UpdatedBy: cfg.UpdatedBy},
		CollectionEnabled:   cfg.CollectionEnabled, DisplayEnabled: cfg.DisplayEnabled, TasksEnabled: cfg.TasksEnabled, AutoScheduleEnabled: &cfg.AutoScheduleEnabled,
		DefaultModel: cfg.DefaultModel, DefaultProtocol: string(cfg.DefaultProtocol), DefaultRange: string(cfg.DefaultRange),
		GatewayOrigin: cfg.GatewayOrigin, OperatorUserID: cfg.OperatorUserID, DailyBudget: cfg.DailyBudget.StringFixed(8), ExpectedNodes: nodes,
	}
}

func (r *providerHallRepository) GetGroup(ctx context.Context, id int64) (*service.ProviderHallGroup, error) {
	// Soft deletion does not invoke the FK cascade. Check the original row on
	// every read so a deleted group cannot retain a visible hall listing.
	exists, err := r.client.Group.Query().Where(group.IDEQ(id), group.DeletedAtIsNil()).Exist(ctx)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, service.ErrProviderHallNotFound
	}
	g, err := r.client.ProviderHallGroup.Get(ctx, id)
	if ent.IsNotFound(err) {
		return &service.ProviderHallGroup{GroupID: id}, nil
	}
	if err != nil {
		return nil, err
	}
	return providerHallGroupFromEnt(g), nil
}

func (r *providerHallRepository) SaveGroup(ctx context.Context, cfg service.ProviderHallGroup) (*service.ProviderHallGroup, error) {
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()
	g, err := tx.Group.Query().Where(group.IDEQ(cfg.GroupID), group.DeletedAtIsNil()).ForUpdate().Only(ctx)
	if ent.IsNotFound(err) {
		return nil, service.ErrProviderHallNotFound
	}
	if err != nil {
		return nil, err
	}
	if (g.Platform != service.PlatformOpenAI && g.Platform != service.PlatformComposite) || (cfg.Listed && g.Status != service.StatusActive) {
		return nil, service.ErrProviderHallConflict
	}
	var updated *ent.ProviderHallGroup
	if cfg.Version == 0 {
		updated, err = tx.ProviderHallGroup.Create().SetID(cfg.GroupID).
			SetListed(cfg.Listed).SetDisplayName(cfg.DisplayName).SetDescription(cfg.Description).
			SetDisplayOrder(cfg.DisplayOrder).SetNillableUpdatedBy(cfg.UpdatedBy).Save(ctx)
		if providerHallUniqueViolation(err) {
			return nil, service.ErrProviderHallConflict
		}
	} else {
		var n int
		n, err = tx.ProviderHallGroup.Update().Where(providerhallgroup.IDEQ(cfg.GroupID), providerhallgroup.VersionEQ(cfg.Version)).
			AddVersion(1).SetUpdatedAt(time.Now().UTC()).SetListed(cfg.Listed).SetDisplayName(cfg.DisplayName).
			SetDescription(cfg.Description).SetDisplayOrder(cfg.DisplayOrder).SetNillableUpdatedBy(cfg.UpdatedBy).Save(ctx)
		if err == nil && n != 1 {
			return nil, service.ErrProviderHallConflict
		}
		if err == nil {
			updated, err = tx.ProviderHallGroup.Get(ctx, cfg.GroupID)
		}
	}
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return providerHallGroupFromEnt(updated), nil
}

func providerHallGroupFromEnt(g *ent.ProviderHallGroup) *service.ProviderHallGroup {
	return &service.ProviderHallGroup{
		ProviderHallVersion: service.ProviderHallVersion{Version: g.Version, UpdatedAt: g.UpdatedAt.UTC(), UpdatedBy: g.UpdatedBy},
		GroupID:             g.ID, Listed: g.Listed, DisplayName: g.DisplayName, Description: g.Description, DisplayOrder: g.DisplayOrder,
	}
}

func (r *providerHallRepository) ListProfiles(ctx context.Context) ([]service.ProviderHallProfile, error) {
	profiles, err := r.client.ProviderHallProfile.Query().Order(ent.Asc(providerhallprofile.FieldModel, providerhallprofile.FieldProtocol, providerhallprofile.FieldID)).All(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]service.ProviderHallProfile, 0, len(profiles))
	targets, err := r.client.ProviderHallTarget.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	groups, err := r.client.Group.Query().Where(group.DeletedAtIsNil()).All(ctx)
	if err != nil {
		return nil, err
	}
	names := map[int64]string{}
	for _, g := range groups {
		names[g.ID] = g.Name
	}
	cfg, err := r.client.ProviderHallConfig.Get(ctx, 1)
	if err != nil {
		return nil, err
	}
	for _, p := range profiles {
		profile := providerHallProfileFromEnt(p)
		profile.IsDefault = cfg.DefaultModel == p.Model && cfg.DefaultProtocol == providerhallconfig.DefaultProtocol(p.Protocol)
		for _, target := range targets {
			if target.ProfileID == p.ID {
				if name, ok := names[target.GroupID]; ok {
					profile.Groups = append(profile.Groups, service.ProviderHallProfileGroup{ID: target.GroupID, Name: name})
				}
			}
		}
		result = append(result, *profile)
	}
	return result, nil
}

func (r *providerHallRepository) SaveProfile(ctx context.Context, p service.ProviderHallProfile) (*service.ProviderHallProfile, error) {
	var inputPrice, cachePrice, cacheRate *decimal.Decimal
	for _, pair := range []struct {
		source *string
		target **decimal.Decimal
	}{
		{p.ReferenceInputPrice, &inputPrice}, {p.ReferenceCachePrice, &cachePrice}, {p.ReferenceCacheRate, &cacheRate},
	} {
		if pair.source != nil {
			value, err := decimal.NewFromString(*pair.source)
			if err != nil {
				return nil, err
			}
			*pair.target = &value
		}
	}
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()
	var updated *ent.ProviderHallProfile
	if p.ID == 0 {
		updated, err = tx.ProviderHallProfile.Create().SetModel(p.Model).SetProtocol(providerhallprofile.Protocol(p.Protocol)).
			SetSupportsTools(p.SupportsTools).SetOutputLimit(p.OutputLimit).SetModelAliases(p.ModelAliases).
			SetNillableReferenceInputPrice(inputPrice).SetNillableReferenceCachePrice(cachePrice).
			SetNillableReferenceCacheRate(cacheRate).SetNillableReferenceConfirmedAt(p.ReferenceConfirmedAt).
			SetNillableUpdatedBy(p.UpdatedBy).Save(ctx)
	} else {
		// Changing the default profile identity must not leave config pointing at
		// a model/protocol pair that no longer exists. Use config -> profile lock order.
		cfg, readErr := tx.ProviderHallConfig.Query().Where(providerhallconfig.IDEQ(1)).ForUpdate().Only(ctx)
		if readErr != nil {
			return nil, readErr
		}
		old, readErr := tx.ProviderHallProfile.Get(ctx, p.ID)
		if ent.IsNotFound(readErr) {
			return nil, service.ErrProviderHallNotFound
		}
		if readErr != nil {
			return nil, readErr
		}
		if cfg.DefaultModel == old.Model && string(cfg.DefaultProtocol) == string(old.Protocol) && (p.Model != old.Model || p.Protocol != string(old.Protocol)) {
			return nil, service.ErrProviderHallConflict.WithMetadata(map[string]string{"reason": "default_profile_identity", "field": "default_model"})
		}
		u := tx.ProviderHallProfile.Update().Where(providerhallprofile.IDEQ(p.ID), providerhallprofile.VersionEQ(p.Version)).
			AddVersion(1).SetUpdatedAt(time.Now().UTC()).SetModel(p.Model).SetProtocol(providerhallprofile.Protocol(p.Protocol)).
			SetSupportsTools(p.SupportsTools).SetOutputLimit(p.OutputLimit).SetModelAliases(p.ModelAliases).
			SetNillableUpdatedBy(p.UpdatedBy).ClearReferenceInputPrice().ClearReferenceCachePrice().ClearReferenceCacheRate().ClearReferenceConfirmedAt()
		if inputPrice != nil {
			u.SetReferenceInputPrice(*inputPrice)
		}
		if cachePrice != nil {
			u.SetReferenceCachePrice(*cachePrice)
		}
		if cacheRate != nil {
			u.SetReferenceCacheRate(*cacheRate)
		}
		if p.ReferenceConfirmedAt != nil {
			u.SetReferenceConfirmedAt(*p.ReferenceConfirmedAt)
		}
		var n int
		n, err = u.Save(ctx)
		if err == nil && n != 1 {
			return nil, service.ErrProviderHallConflict
		}
		if err == nil {
			updated, err = tx.ProviderHallProfile.Get(ctx, p.ID)
		}
	}
	if providerHallUniqueViolation(err) {
		return nil, service.ErrProviderHallProfileExists
	}
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return providerHallProfileFromEnt(updated), nil
}

func providerHallUniqueViolation(err error) bool {
	var pg *pq.Error
	return errors.As(err, &pg) && pg.Code == "23505"
}

func providerHallProfileFromEnt(p *ent.ProviderHallProfile) *service.ProviderHallProfile {
	result := &service.ProviderHallProfile{
		ProviderHallVersion: service.ProviderHallVersion{Version: p.Version, UpdatedAt: p.UpdatedAt.UTC(), UpdatedBy: p.UpdatedBy},
		ID:                  p.ID, Model: p.Model, Protocol: string(p.Protocol), SupportsTools: p.SupportsTools,
		OutputLimit: p.OutputLimit, ModelAliases: append([]string{}, p.ModelAliases...),
	}
	if p.ReferenceConfirmedAt != nil {
		utc := p.ReferenceConfirmedAt.UTC()
		result.ReferenceConfirmedAt = &utc
	}
	for _, pair := range []struct {
		source *decimal.Decimal
		target **string
	}{
		{p.ReferenceInputPrice, &result.ReferenceInputPrice}, {p.ReferenceCachePrice, &result.ReferenceCachePrice}, {p.ReferenceCacheRate, &result.ReferenceCacheRate},
	} {
		if pair.source != nil {
			value := pair.source.StringFixed(10)
			*pair.target = &value
		}
	}
	return result
}
