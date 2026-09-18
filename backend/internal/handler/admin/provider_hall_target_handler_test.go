//go:build unit

package admin

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type providerHallTargetHandlerRepo struct {
	service.ProviderHallRepository
	saved *service.ProviderHallTargetSet
	err   error
}

func (r *providerHallTargetHandlerRepo) SaveTargets(_ context.Context, input service.ProviderHallTargetSet) (*service.ProviderHallTargetSet, error) {
	r.saved = &input
	return &input, r.err
}

func TestProviderHallTargetsHTTPContract(t *testing.T) {
	for _, tc := range []struct {
		name, body string
		err        error
		status     int
	}{
		{"missing_items", `{"version":1}`, nil, 400},
		{"null_items", `{"version":1,"items":null}`, nil, 400},
		{"missing_version", `{"items":[]}`, nil, 400},
		{"forged_key", `{"version":1,"items":[{"profile_id":2,"api_key":"secret"}]}`, nil, 400},
		{"forged_id", `{"version":1,"items":[{"profile_id":2,"id":99}]}`, nil, 400},
		{"disable_all", `{"version":1,"items":[]}`, nil, 200},
		{"conflict", `{"version":1,"items":[]}`, service.ErrProviderHallConflict, 409},
		{"operator", `{"version":1,"items":[]}`, service.ErrProviderHallOperator, 400},
	} {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("PUT", "/groups/4/targets", strings.NewReader(tc.body))
			c.Params = gin.Params{{Key: "id", Value: "4"}}
			c.Set(string(middleware.ContextKeyUser), middleware.AuthSubject{UserID: 7})
			c.Set(string(middleware.ContextKeyUserRole), service.RoleAdmin)
			repo := &providerHallTargetHandlerRepo{err: tc.err}
			h := NewProviderHallHandler(service.NewProviderHallService(repo, nil, nil, nil), nil, nil, nil, nil, nil)
			h.UpdateTargets(c)
			require.Equal(t, tc.status, w.Code, w.Body.String())
			if tc.status == 200 {
				require.Equal(t, int64(4), repo.saved.GroupID)
				require.Equal(t, int64(7), *repo.saved.UpdatedBy)
			}
		})
	}
}
