//go:build unit

package admin

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type keepStatusActiveAdminServiceStub struct {
	*stubAdminService
	account *service.Account
	id      int64
	enabled bool
}

func (s *keepStatusActiveAdminServiceStub) SetAccountKeepStatusActive(_ context.Context, id int64, enabled bool) (*service.Account, error) {
	s.id = id
	s.enabled = enabled
	if s.account == nil {
		s.account = &service.Account{ID: id, Name: "token", Platform: service.PlatformOpenAI, Type: service.AccountTypeAPIKey, Status: service.StatusActive, Schedulable: true}
	}
	if s.account.Extra == nil {
		s.account.Extra = map[string]any{}
	}
	s.account.Extra[service.KeepStatusActiveExtraKey] = enabled
	return s.account, nil
}

func TestAccountHandlerSetKeepStatusActive(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &keepStatusActiveAdminServiceStub{stubAdminService: newStubAdminService()}
	h := NewAccountHandler(svc, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	router := gin.New()
	router.POST("/api/v1/admin/accounts/:id/keep-status-active", h.SetKeepStatusActive)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/accounts/7579/keep-status-active", strings.NewReader(`{"enabled":true}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, int64(7579), svc.id)
	require.True(t, svc.enabled)
	var payload struct {
		Data struct {
			KeepStatusActive bool `json:"keep_status_active"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &payload))
	require.True(t, payload.Data.KeepStatusActive)
}
