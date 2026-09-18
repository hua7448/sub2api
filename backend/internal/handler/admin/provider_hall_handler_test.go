//go:build unit

package admin

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type providerHallHandlerRepo struct {
	service.ProviderHallRepository
	conflict bool
}

func (r *providerHallHandlerRepo) UpdateConfig(_ context.Context, cfg service.ProviderHallConfig) (*service.ProviderHallConfig, error) {
	if r.conflict {
		return nil, service.ErrProviderHallConflict
	}
	cfg.Version++
	return &cfg, nil
}

func TestProviderHallAdminAuthentication(t *testing.T) {
	h := NewProviderHallHandler(nil, nil, nil, nil, nil, nil)
	for _, method := range []gin.HandlerFunc{h.GetConfig, h.UpdateConfig, h.GetGroup, h.UpdateGroup, h.ListProfiles, h.CreateProfile, h.UpdateProfile, h.GetTargets, h.UpdateTargets} {
		for _, role := range []string{"", service.RoleUser} {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
			if role != "" {
				c.Set(string(middleware.ContextKeyUser), middleware.AuthSubject{UserID: 1})
				c.Set(string(middleware.ContextKeyUserRole), role)
			}
			method(c)
			if role == "" {
				require.Equal(t, 401, w.Code)
			} else {
				require.Equal(t, 403, w.Code)
			}
		}
	}
}

func TestProviderHallConfigHTTPContract(t *testing.T) {
	for _, tc := range []struct {
		name, body string
		conflict   bool
		status     int
	}{
		{"saved", `{"version":1,"default_protocol":"responses","default_range":"6h","daily_budget":"0"}`, false, 200},
		{"conflict", `{"version":1,"default_protocol":"responses","default_range":"6h","daily_budget":"0"}`, true, 409},
		{"missing_version", `{}`, false, 400},
		{"numeric_money", `{"version":1,"daily_budget":1.2}`, false, 400},
		{"forged_actor", `{"version":1,"updated_by":999}`, false, 400},
		{"multiple_objects", `{"version":1} {}`, false, 400},
		{"null", `null`, false, 400},
		{"enable_not_ready", `{"version":1,"default_protocol":"responses","default_range":"6h","daily_budget":"0","tasks_enabled":true}`, false, 409},
	} {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodPut, "/config", strings.NewReader(tc.body))
			c.Set(string(middleware.ContextKeyUser), middleware.AuthSubject{UserID: 7})
			c.Set(string(middleware.ContextKeyUserRole), service.RoleAdmin)
			h := NewProviderHallHandler(service.NewProviderHallService(&providerHallHandlerRepo{conflict: tc.conflict}, nil, nil, nil), nil, nil, nil, nil, nil)
			h.UpdateConfig(c)
			require.Equal(t, tc.status, w.Code, w.Body.String())
			if tc.status == 200 {
				require.Contains(t, w.Body.String(), `"updated_by":7`)
				require.Contains(t, w.Body.String(), `"daily_budget":"0"`)
			}
			if tc.name == "conflict" {
				require.Contains(t, w.Body.String(), "PROVIDER_HALL_VERSION_CONFLICT")
			}
		})
	}
}
