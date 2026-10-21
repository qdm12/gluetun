package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_VPNHandler_PatchSettings_rejects_unknown_field_when_payload_uses_legacy_names(t *testing.T) {
	t.Parallel()

	const body = `{"Provider":{"ServerSelection":{"Countries":["Netherlands"]}}}`
	request := httptest.NewRequest(http.MethodPut, "/v1/vpn/settings", strings.NewReader(body))
	response := httptest.NewRecorder()
	handler := &vpnHandler{looper: stubVPNLooper{}}

	handler.patchSettings(response, request)

	require.Equal(t, http.StatusBadRequest, response.Code)
	assert.Contains(t, response.Body.String(), `json: unknown field "ServerSelection"`)
	assert.NotContains(t, response.Body.String(), "settings left unchanged")
}

type stubVPNLooper struct{}

func (stubVPNLooper) GetStatus() models.LoopStatus {
	return ""
}

func (stubVPNLooper) ApplyStatus(context.Context, models.LoopStatus) (string, error) {
	return "", nil
}

func (stubVPNLooper) GetSettings() settings.VPN {
	return settings.VPN{}
}

func (stubVPNLooper) SetSettings(context.Context, settings.VPN) string {
	return ""
}
