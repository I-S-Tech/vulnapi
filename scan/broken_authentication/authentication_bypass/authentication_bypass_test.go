package authenticationbypass_test

import (
	"net/http"
	"testing"

	"github.com/cerberauth/vulnapi/internal/auth"
	"github.com/cerberauth/vulnapi/internal/operation"
	"github.com/cerberauth/vulnapi/internal/request"
	authenticationbypass "github.com/cerberauth/vulnapi/scan/broken_authentication/authentication_bypass"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthenticationByPassScanHandler_Skipped_WhenNoAuthSecurityScheme(t *testing.T) {
	securityScheme := auth.MustNewNoAuthSecurityScheme()
	operation := operation.MustNewOperation(http.MethodGet, "http://localhost:8080/", nil, nil)

	report, err := authenticationbypass.ScanHandler(operation, securityScheme)

	require.NoError(t, err)
	assert.True(t, report.Issues[0].HasBeenSkipped())
}

func TestAuthenticationByPassScanHandler_Failed_WhenAuthIsByPassed(t *testing.T) {
	client := request.GetDefaultClient()
	httpmock.ActivateNonDefault(client.Client)
	defer httpmock.DeactivateAndReset()

	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	securityScheme := auth.MustNewAuthorizationBearerSecurityScheme("default", &token)
	operation := operation.MustNewOperation(http.MethodGet, "http://localhost:8080/", nil, client)
	httpmock.RegisterResponder(operation.Method, operation.URL.String(), httpmock.NewBytesResponder(http.StatusNoContent, nil))

	report, err := authenticationbypass.ScanHandler(operation, securityScheme)

	require.NoError(t, err)
	assert.True(t, report.Issues[0].HasFailed())
}

func TestAuthenticationByPassScanHandler_Passed_WhenAuthIsNotByPassed(t *testing.T) {
	client := request.GetDefaultClient()
	httpmock.ActivateNonDefault(client.Client)
	defer httpmock.DeactivateAndReset()

	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	securityScheme := auth.MustNewAuthorizationBearerSecurityScheme("default", &token)
	operation := operation.MustNewOperation(http.MethodGet, "http://localhost:8080/", nil, client)
	httpmock.RegisterResponder(operation.Method, operation.URL.String(), httpmock.NewBytesResponder(http.StatusUnauthorized, nil))

	report, err := authenticationbypass.ScanHandler(operation, securityScheme)

	require.NoError(t, err)
	assert.True(t, report.Issues[0].HasPassed())
}
