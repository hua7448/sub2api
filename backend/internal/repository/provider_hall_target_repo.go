package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/apikey"
	"github.com/Wei-Shaw/sub2api/ent/group"
	"github.com/Wei-Shaw/sub2api/ent/providerhallconfig"
	"github.com/Wei-Shaw/sub2api/ent/providerhallgroup"
	"github.com/Wei-Shaw/sub2api/ent/providerhallprobekey"
	"github.com/Wei-Shaw/sub2api/ent/providerhallprofile"
	"github.com/Wei-Shaw/sub2api/ent/providerhalltarget"
	"github.com/Wei-Shaw/sub2api/ent/usagelog"
	"github.com/Wei-Shaw/sub2api/ent/user"
	"github.com/Wei-Shaw/sub2api/ent/usersubscription"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

func (r *providerHallRepository) GetTargets(ctx context.Context, groupID int64) (*service.ProviderHallTargetSet, error) {
	tx, err := r.client.BeginTx(ctx, &sql.TxOptions{ReadOnly: true, Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()
	g, err := tx.Group.Query().Where(group.IDEQ(groupID), group.DeletedAtIsNil()).Only(ctx)
	if ent.IsNotFound(err) {
		return nil, service.ErrProviderHallNotFound
	}
	if err != nil {
		return nil, err
	}
	if g.Platform != service.PlatformOpenAI && g.Platform != service.PlatformComposite {
		return nil, service.ErrProviderHallTargetRoute
	}
	result, err := providerHallReadTargets(ctx, tx.Client(), groupID)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return result, nil
}

func providerHallReadTargets(ctx context.Context, client *ent.Client, groupID int64) (*service.ProviderHallTargetSet, error) {
	result := &service.ProviderHallTargetSet{GroupID: groupID, Items: []service.ProviderHallTarget{}}
	listing, err := client.ProviderHallGroup.Get(ctx, groupID)
	if ent.IsNotFound(err) {
		return result, nil
	}
	if err != nil {
		return nil, err
	}
	result.ProviderHallVersion = providerHallGroupFromEnt(listing).ProviderHallVersion
	rows, err := client.ProviderHallTarget.Query().Where(providerhalltarget.GroupIDEQ(groupID)).Order(ent.Asc(providerhalltarget.FieldProfileID)).All(ctx)
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		result.Items = append(result.Items, service.ProviderHallTarget{
			ProviderHallVersion: service.ProviderHallVersion{Version: row.Version, UpdatedAt: row.UpdatedAt.UTC(), UpdatedBy: row.UpdatedBy},
			ID:                  row.ID, GroupID: row.GroupID, ProfileID: row.ProfileID, ProbeKeyID: row.ProbeKeyID, Enabled: row.Enabled,
			ProbeIntervalSeconds: row.ProbeIntervalSeconds, VerificationIntervalSeconds: row.VerificationIntervalSeconds,
		})
	}
	return result, nil
}

func (r *providerHallRepository) SaveTargets(ctx context.Context, input service.ProviderHallTargetSet) (*service.ProviderHallTargetSet, error) {
	if input.UpdatedBy == nil || *input.UpdatedBy < 1 || len(input.Items) > service.ProviderHallMaxTargetsPerGroup {
		return nil, service.ErrProviderHallConflict
	}
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()
	// Shared lock order with profile/config writes: config -> group -> listing.
	cfg, err := tx.ProviderHallConfig.Query().Where(providerhallconfig.IDEQ(1)).ForUpdate().Only(ctx)
	if err != nil {
		return nil, err
	}
	g, err := tx.Group.Query().Where(group.IDEQ(input.GroupID), group.DeletedAtIsNil()).ForUpdate().Only(ctx)
	if ent.IsNotFound(err) {
		return nil, service.ErrProviderHallNotFound
	}
	if err != nil {
		return nil, err
	}
	if g.Platform != service.PlatformOpenAI && g.Platform != service.PlatformComposite {
		return nil, service.ErrProviderHallTargetRoute
	}
	listing, err := tx.ProviderHallGroup.Query().Where(providerhallgroup.IDEQ(input.GroupID)).ForUpdate().Only(ctx)
	if ent.IsNotFound(err) {
		if input.Version != 0 {
			return nil, service.ErrProviderHallConflict
		}
		listing, err = tx.ProviderHallGroup.Create().SetID(input.GroupID).SetNillableUpdatedBy(input.UpdatedBy).Save(ctx)
	} else if err == nil {
		if listing.Version != input.Version {
			return nil, service.ErrProviderHallConflict
		}
		listing, err = tx.ProviderHallGroup.UpdateOneID(input.GroupID).AddVersion(1).SetNillableUpdatedBy(input.UpdatedBy).Save(ctx)
	}
	if err != nil {
		return nil, err
	}
	oldRows, err := tx.ProviderHallTarget.Query().Where(providerhalltarget.GroupIDEQ(input.GroupID)).All(ctx)
	if err != nil {
		return nil, err
	}
	old := make(map[int64]*ent.ProviderHallTarget, len(oldRows))
	for _, row := range oldRows {
		old[row.ProfileID] = row
	}
	profiles := make([]int64, 0, len(input.Items))
	seen := map[int64]bool{}
	union := len(old)
	for _, item := range input.Items {
		if item.ProfileID < 1 || seen[item.ProfileID] {
			return nil, service.ErrProviderHallConflict
		}
		seen[item.ProfileID] = true
		profiles = append(profiles, item.ProfileID)
		if old[item.ProfileID] == nil {
			union++
		}
	}
	if union > service.ProviderHallMaxTargetsPerGroup {
		return nil, service.ErrProviderHallConflict
	}
	profileRows, err := tx.ProviderHallProfile.Query().Where(providerhallprofile.IDIn(profiles...)).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(profileRows) != len(profiles) {
		return nil, service.ErrProviderHallNotFound
	}
	profileMap := make(map[int64]*ent.ProviderHallProfile, len(profileRows))
	for _, row := range profileRows {
		profileMap[row.ID] = row
	}

	for _, item := range input.Items {
		previous := old[item.ProfileID]
		newBinding := item.ProbeKeyID != nil && (previous == nil || previous.ProbeKeyID == nil || *previous.ProbeKeyID != *item.ProbeKeyID)
		if item.Enabled || newBinding {
			if err := providerHallValidateAndRegisterKey(ctx, tx.Client(), cfg, g, profileMap[item.ProfileID], item.ProbeKeyID, *input.UpdatedBy); err != nil {
				return nil, err
			}
		}
		if previous == nil {
			_, err = tx.ProviderHallTarget.Create().SetGroupID(input.GroupID).SetProfileID(item.ProfileID).
				SetNillableProbeKeyID(item.ProbeKeyID).SetEnabled(item.Enabled).
				SetProbeIntervalSeconds(item.ProbeIntervalSeconds).SetVerificationIntervalSeconds(item.VerificationIntervalSeconds).
				SetNillableUpdatedBy(input.UpdatedBy).Save(ctx)
		} else {
			u := tx.ProviderHallTarget.UpdateOneID(previous.ID).AddVersion(1).ClearProbeKeyID().SetEnabled(item.Enabled).
				SetProbeIntervalSeconds(item.ProbeIntervalSeconds).SetVerificationIntervalSeconds(item.VerificationIntervalSeconds).
				SetNillableUpdatedBy(input.UpdatedBy)
			if item.ProbeKeyID != nil {
				u.SetProbeKeyID(*item.ProbeKeyID)
			}
			_, err = u.Save(ctx)
		}
		if err != nil {
			return nil, err
		}
		delete(old, item.ProfileID)
	}
	for _, omitted := range old {
		if omitted.Enabled {
			if _, err := tx.ProviderHallTarget.UpdateOneID(omitted.ID).SetEnabled(false).AddVersion(1).SetNillableUpdatedBy(input.UpdatedBy).Save(ctx); err != nil {
				return nil, err
			}
		}
	}
	result, err := providerHallReadTargets(ctx, tx.Client(), input.GroupID)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return result, nil
}

func providerHallValidateAndRegisterKey(ctx context.Context, client *ent.Client, cfg *ent.ProviderHallConfig, g *ent.Group, p *ent.ProviderHallProfile, keyID *int64, actorID int64) error {
	if cfg.OperatorUserID == nil {
		return service.ErrProviderHallOperator
	}
	if keyID == nil {
		return service.ErrProviderHallProbeKey
	}
	u, err := client.User.Query().Where(user.IDEQ(*cfg.OperatorUserID), user.DeletedAtIsNil()).WithAllowedGroups().Only(ctx)
	if ent.IsNotFound(err) {
		return service.ErrProviderHallOperator
	}
	if err != nil {
		return err
	}
	operator := userEntityToService(u)
	for _, allowed := range u.Edges.AllowedGroups {
		operator.AllowedGroups = append(operator.AllowedGroups, allowed.ID)
	}
	now := time.Now().UTC()
	hasSubscription := false
	if g.SubscriptionType == service.SubscriptionTypeSubscription {
		hasSubscription, err = client.UserSubscription.Query().Where(
			usersubscription.UserIDEQ(u.ID), usersubscription.GroupIDEQ(g.ID), usersubscription.DeletedAtIsNil(),
			usersubscription.StatusEQ(service.SubscriptionStatusActive), usersubscription.ExpiresAtGT(now)).Exist(ctx)
		if err != nil {
			return err
		}
	}
	k, err := client.APIKey.Query().Where(apikey.IDEQ(*keyID), apikey.DeletedAtIsNil()).ForUpdate().Only(ctx)
	if ent.IsNotFound(err) {
		return service.ErrProviderHallProbeKey
	}
	if err != nil {
		return err
	}
	var route service.CompositeRouteDecision
	if g.Platform == service.PlatformComposite {
		route, err = service.NewCompositeRouteResolver(NewCompositeModelRouteRepository(client)).Resolve(ctx, g.ID, p.Model, string(p.Protocol))
		if err != nil {
			return err
		}
	}
	if err := service.ValidateProviderHallTargetBinding(groupEntityToService(g), providerHallProfileFromEnt(p), apiKeyEntityToService(k), operator, hasSubscription, route, now); err != nil {
		return err
	}
	registered, err := client.ProviderHallProbeKey.Get(ctx, *keyID)
	if err == nil {
		if registered.OperatorUserID != u.ID || registered.GroupID != g.ID {
			return service.ErrProviderHallProbeKey
		}
		return nil
	}
	if !ent.IsNotFound(err) {
		return err
	}
	// First registration accepts only unused keys. Retired probe keys may be
	// re-enabled, but a normal user's historical traffic cannot be reclassified.
	if k.LastUsedAt != nil || k.QuotaUsed != 0 || k.Usage5h != 0 || k.Usage1d != 0 || k.Usage7d != 0 {
		return service.ErrProviderHallProbeKey
	}
	used, err := client.UsageLog.Query().Where(usagelog.APIKeyIDEQ(*keyID)).Exist(ctx)
	if err != nil {
		return err
	}
	if used {
		return service.ErrProviderHallProbeKey
	}
	_, err = client.ProviderHallProbeKey.Create().SetID(*keyID).SetOperatorUserID(u.ID).SetGroupID(g.ID).
		SetRegisteredBy(actorID).SetRegisteredAt(now).Save(ctx)
	return err
}

func (r *providerHallRepository) IsProbeKey(ctx context.Context, keyID int64, startedAt time.Time) (bool, error) {
	if keyID < 1 || startedAt.IsZero() {
		return false, nil
	}
	return r.client.ProviderHallProbeKey.Query().Where(providerhallprobekey.IDEQ(keyID), providerhallprobekey.RegisteredAtLTE(startedAt.UTC())).Exist(ctx)
}
