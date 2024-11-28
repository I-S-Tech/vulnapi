package weaksecret_test

import (
	"net/http"
	"testing"

	"github.com/cerberauth/vulnapi/internal/auth"
	"github.com/cerberauth/vulnapi/internal/operation"
	"github.com/cerberauth/vulnapi/internal/request"
	weaksecret "github.com/cerberauth/vulnapi/scan/broken_authentication/jwt/weak_secret"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWeakHMACSecretScanHandler_WithoutSecurityScheme(t *testing.T) {
	securityScheme := auth.MustNewNoAuthSecurityScheme()
	operation := operation.MustNewOperation(http.MethodGet, "http://localhost:8080/", nil, nil)

	report, err := weaksecret.ScanHandler(operation, securityScheme)

	require.NoError(t, err)
	assert.True(t, report.Issues[0].HasBeenSkipped())
}

func TestWeakHMACSecretScanHandler_WithJWTUsingOtherAlg(t *testing.T) {
	token := "eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJhYmMxMjMifQ.vLBmArLmAKEshqJa3px6qYfrkAfiwBrKPs5dCMxqj9bdiEKR5W4o0Srxt6VHZKzsxIGMTTsqpW21lKnYsLw5DA"
	securityScheme := auth.MustNewAuthorizationBearerSecurityScheme("token", &token)
	operation := operation.MustNewOperation(http.MethodGet, "http://localhost:8080/", nil, nil)

	report, err := weaksecret.ScanHandler(operation, securityScheme)

	require.NoError(t, err)
	assert.True(t, report.Issues[0].HasBeenSkipped())
}

func TestWeakHMACSecretScanHandler_WithoutJWT(t *testing.T) {
	securityScheme := auth.MustNewAuthorizationBearerSecurityScheme("token", nil)
	operation := operation.MustNewOperation(http.MethodGet, "http://localhost:8080/", nil, nil)

	report, err := weaksecret.ScanHandler(operation, securityScheme)

	require.NoError(t, err)
	assert.True(t, report.Issues[0].HasBeenSkipped())
}

func TestWeakHMACSecretScanHandler_Failed_WithWeakJWT(t *testing.T) {
	client := request.GetDefaultClient()
	httpmock.ActivateNonDefault(client.Client)
	defer httpmock.DeactivateAndReset()

	secret := "secret"
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.t-IDcSemACt8x4iTMCda8Yhe3iZaWbvV5XKSTbuAn0M"
	securityScheme := auth.MustNewAuthorizationBearerSecurityScheme("token", &token)
	operation := operation.MustNewOperation(http.MethodGet, "http://localhost:8080/", nil, client)
	httpmock.RegisterResponder(operation.Method, operation.URL.String(), httpmock.NewBytesResponder(http.StatusOK, nil))

	report, err := weaksecret.ScanHandler(operation, securityScheme)

	assert.NoError(t, err)
	assert.Equal(t, 1, httpmock.GetTotalCallCount())
	assert.True(t, report.Issues[0].HasFailed())
	assert.NotNil(t, report.Data)
	assert.Equal(t, &secret, report.Data.(*weaksecret.WeakSecretData).Secret)
}

func TestWeakHMACSecretScanHandler_Failed_WithExpiredJWTSignedWithWeakSecret(t *testing.T) {
	client := request.GetDefaultClient()
	httpmock.ActivateNonDefault(client.Client)
	defer httpmock.DeactivateAndReset()

	secret := "secret"
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTYyMzkwMjJ9.7BbIenT4-HobiMHaMUQdNcJ6lD_QQkKnImP9IprJFvU"
	securityScheme := auth.MustNewAuthorizationBearerSecurityScheme("token", &token)
	operation := operation.MustNewOperation(http.MethodGet, "http://localhost:8080/", nil, client)
	httpmock.RegisterResponder(operation.Method, operation.URL.String(), httpmock.NewBytesResponder(http.StatusOK, nil))

	report, err := weaksecret.ScanHandler(operation, securityScheme)

	assert.NoError(t, err)
	assert.Equal(t, 1, httpmock.GetTotalCallCount())
	assert.True(t, report.Issues[0].HasFailed())
	assert.NotNil(t, report.Data)
	assert.Equal(t, &secret, report.Data.(*weaksecret.WeakSecretData).Secret)
}

func TestWeakHMACSecretScanHandler_Passed_WithStrongerJWT(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.MWUarT7Q4e5DqnZbdr7VKw3rx9VW-CrvoVkfpllS4CY"
	securityScheme := auth.MustNewAuthorizationBearerSecurityScheme("token", &token)
	operation := operation.MustNewOperation(http.MethodGet, "http://localhost:8080/", nil, nil)
	httpmock.RegisterResponder(operation.Method, operation.URL.String(), httpmock.NewBytesResponder(http.StatusUnauthorized, nil))

	report, err := weaksecret.ScanHandler(operation, securityScheme)

	assert.NoError(t, err)
	assert.Equal(t, 0, httpmock.GetTotalCallCount())
	assert.True(t, report.Issues[0].HasPassed())
	assert.Nil(t, report.Data)
}

func TestWeakHMACSecretScanHandler_Failed_WithUnorderedClaims(t *testing.T) {
	client := request.GetDefaultClient()
	httpmock.ActivateNonDefault(client.Client)
	defer httpmock.DeactivateAndReset()

	secret := "secret"
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJuYmYiOjIwMTYyMzkwMjJ9.ymnE0GznV0dMkjANTQl8IqBSlTi9RFWfBeT42jBNrU4"
	securityScheme := auth.MustNewAuthorizationBearerSecurityScheme("token", &token)
	operation := operation.MustNewOperation(http.MethodGet, "http://localhost:8080/", nil, client)
	httpmock.RegisterResponder(operation.Method, operation.URL.String(), httpmock.NewBytesResponder(http.StatusOK, nil))

	report, err := weaksecret.ScanHandler(operation, securityScheme)

	assert.NoError(t, err)
	assert.Equal(t, 1, httpmock.GetTotalCallCount())
	assert.True(t, report.Issues[0].HasFailed())
	assert.NotNil(t, report.Data)
	assert.Equal(t, &secret, report.Data.(*weaksecret.WeakSecretData).Secret)
}
