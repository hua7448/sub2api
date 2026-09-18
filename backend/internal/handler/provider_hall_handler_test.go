//go:build unit

package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/handler/dto"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type hallCfgRepoStub struct {
	service.ProviderHallRepository
	cfg service.ProviderHallConfig
}

func (r *hallCfgRepoStub) GetConfig(context.Context) (*service.ProviderHallConfig, error) {
	c := r.cfg
	return &c, nil
}
func (r *hallCfgRepoStub) GetGroup(_ context.Context, id int64) (*service.ProviderHallGroup, error) {
	return &service.ProviderHallGroup{GroupID: id, Listed: true, DisplayName: "G"}, nil
}
func (r *hallCfgRepoStub) ListProfiles(context.Context) ([]service.ProviderHallProfile, error) {
	return nil, nil
}

type hallReadStub struct{ window time.Time }

func (r *hallReadStub) ListListedGroups(context.Context) ([]service.ProviderHallListedGroup, error) {
	return []service.ProviderHallListedGroup{{ProviderHallGroup: service.ProviderHallGroup{GroupID: 1, Listed: true, DisplayName: "G"}, GroupName: "g"}}, nil
}
func (r *hallReadStub) ListEnabledTargets(context.Context) ([]service.ProviderHallReadTarget, error) {
	return []service.ProviderHallReadTarget{{TargetID: 1, GroupID: 1, ProfileID: 10, Model: "gpt-5", Protocol: "responses", ProbeIntervalSeconds: 300}}, nil
}
func (r *hallReadStub) LatestWindow(context.Context) (*time.Time, int, error) {
	t := r.window
	return &t, 1, nil
}
func (r *hallReadStub) LoadSnapshots(context.Context, time.Time, int) ([]service.ProviderHallSnapshotRow, error) {
	return nil, nil
}
func (r *hallReadStub) LoadSeries(context.Context, []time.Time, int) ([]service.ProviderHallSeriesPoint, error) {
	return nil, nil
}
func (r *hallReadStub) LatestVerifications(context.Context) ([]service.ProviderHallVerification, error) {
	return nil, nil
}
func (r *hallReadStub) ProbeSeriesAll(context.Context, time.Time, time.Time, int) ([]service.ProviderHallGroupProbePoint, error) {
	return nil, nil
}

type hallAccessStub struct{}

func (hallAccessStub) GetAvailableGroups(context.Context, int64) ([]service.Group, error) {
	return []service.Group{{ID: 1, Name: "g", Platform: service.PlatformOpenAI, Status: service.StatusActive, RateMultiplier: 1}}, nil
}

func newHallHandler(t *testing.T, display bool) *ProviderHallUserHandler {
	t.Helper()
	cfg := &hallCfgRepoStub{cfg: service.ProviderHallConfig{CollectionEnabled: true, DisplayEnabled: display, DefaultModel: "gpt-5", DefaultProtocol: "responses", DefaultRange: "24h"}}
	hall := service.NewProviderHallService(cfg, nil, nil, nil)
	hall.SetGroupAccessForTest(hallAccessStub{})
	svc := service.NewProviderHallQueryService(&hallReadStub{window: time.Now().UTC().Truncate(time.Minute)}, cfg, hall, hallAccessStub{}, nil, nil, nil)
	return NewProviderHallUserHandler(svc)
}

func hallRequest(t *testing.T, h *ProviderHallUserHandler, target string, fn func(*gin.Context), params gin.Params) (*httptest.ResponseRecorder, map[string]any) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, target, nil)
	c.Params = params
	c.Set(string(middleware.ContextKeyUser), middleware.AuthSubject{UserID: 7})
	fn(c)
	var body map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &body)
	return rec, body
}

func TestProviderHallUserHandlerValidation(t *testing.T) {
	h := newHallHandler(t, true)
	bad := []string{
		"/provider-hall?sort=name",
		"/provider-hall?sort=rate:up",
		"/provider-hall?sort=rate,rate",
		"/provider-hall?sort=rate,ttft_fast95,cache_rate,success_rate",
		"/provider-hall?page=0",
		"/provider-hall?page=x",
		"/provider-hall?page_size=101",
		"/provider-hall?range=1h",
		"/provider-hall?protocol=grpc",
		"/provider-hall?snapshot_id=nope",
		"/provider-hall?search=" + strings.Repeat("a", 101),
	}
	for _, target := range bad {
		rec, body := hallRequest(t, h, target, h.List, nil)
		require.Equal(t, http.StatusBadRequest, rec.Code, target)
		require.Equal(t, "PROVIDER_HALL_INVALID_QUERY", body["reason"], target)
	}
	rec, body := hallRequest(t, h, "/provider-hall?sort=rate:desc,ttft_fast95&page=1&page_size=100&search=g", h.List, nil)
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	data := body["data"].(map[string]any)
	require.Equal(t, "24h", data["range"], "range defaults to the hall configuration")
	require.EqualValues(t, 100, data["pagination"].(map[string]any)["page_size"])
	require.Len(t, data["items"], 1)

	for _, target := range []string{"/provider-hall/groups/1?range=2h", "/provider-hall/groups/1/verifications?page_size=51", "/provider-hall/groups/1/verifications?profile_id=0"} {
		fn := h.GetGroup
		if strings.Contains(target, "verifications") {
			fn = h.ListVerifications
		}
		rec, body := hallRequest(t, h, target, fn, gin.Params{{Key: "id", Value: "1"}})
		require.Equal(t, http.StatusBadRequest, rec.Code, target)
		require.Equal(t, "PROVIDER_HALL_INVALID_QUERY", body["reason"], target)
	}
	for _, id := range []string{"0", "abc", "999"} {
		rec, _ := hallRequest(t, h, "/provider-hall/groups/"+id, h.GetGroup, gin.Params{{Key: "id", Value: id}})
		if id == "999" {
			// Listed and permitted per the stubs? No: access stub only grants group 1.
			require.Equal(t, http.StatusNotFound, rec.Code, id)
			continue
		}
		require.Equal(t, http.StatusNotFound, rec.Code, id)
	}
	rec, body = hallRequest(t, h, "/provider-hall/groups/1?range=6h", h.GetGroup, gin.Params{{Key: "id", Value: "1"}})
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	detail := body["data"].(map[string]any)
	require.Equal(t, "6h", detail["trend"].(map[string]any)["range"])
	require.Contains(t, detail, "detail")
	require.Contains(t, detail, "profiles")
	rec, body = hallRequest(t, h, "/provider-hall/groups/1/verifications", h.ListVerifications, gin.Params{{Key: "id", Value: "1"}})
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	require.EqualValues(t, 0, body["data"].(map[string]any)["pagination"].(map[string]any)["total"])
}

func TestProviderHallUserHandlerDisplayGuardAndAuth(t *testing.T) {
	off := newHallHandler(t, false)
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/provider-hall", nil)
	off.DisplayGuard()(c)
	require.True(t, c.IsAborted())
	require.Equal(t, http.StatusNotFound, rec.Code)

	on := newHallHandler(t, true)
	rec = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/provider-hall", nil)
	on.DisplayGuard()(c)
	require.False(t, c.IsAborted())

	// No auth subject → 401, even when display is on.
	rec = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/provider-hall", nil)
	on.List(c)
	require.Equal(t, http.StatusUnauthorized, rec.Code)

	var nilHandler *ProviderHallUserHandler
	rec = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/provider-hall", nil)
	nilHandler.DisplayGuard()(c)
	require.Equal(t, http.StatusNotFound, rec.Code)
}

// The user contract must never carry account, credential, upstream or absolute
// traffic-count fields. The walk covers every nested type of the three
// responses; names are checked on the JSON tag, which is what ships.
func TestProviderHallUserDTOHasNoSensitiveKeys(t *testing.T) {
	forbidden := []string{"account", "token", "key", "upstream", "count"}
	allowed := map[string]bool{"probe_tokens": true, "bucket_seconds": true}
	seen := map[reflect.Type]bool{}
	var walk func(t *testing.T, typ reflect.Type, path string)
	walk = func(t *testing.T, typ reflect.Type, path string) {
		for typ.Kind() == reflect.Ptr || typ.Kind() == reflect.Slice || typ.Kind() == reflect.Array || typ.Kind() == reflect.Map {
			typ = typ.Elem()
		}
		if typ.Kind() != reflect.Struct || seen[typ] || typ.PkgPath() == "time" {
			return
		}
		seen[typ] = true
		for i := 0; i < typ.NumField(); i++ {
			f := typ.Field(i)
			tag := strings.Split(f.Tag.Get("json"), ",")[0]
			name := tag
			if name == "" {
				name = strings.ToLower(f.Name)
			}
			if !allowed[name] && name != "-" {
				for _, bad := range forbidden {
					require.NotContains(t, name, bad, "%s.%s", path, name)
				}
			}
			walk(t, f.Type, path+"."+name)
		}
	}
	for _, typ := range []reflect.Type{
		reflect.TypeOf(dto.ProviderHallListResponse{}),
		reflect.TypeOf(dto.ProviderHallDetailResponse{}),
		reflect.TypeOf(dto.ProviderHallVerificationListResponse{}),
	} {
		walk(t, typ, typ.Name())
	}
	require.Greater(t, len(seen), 15, "the walk must reach the nested types")
}
