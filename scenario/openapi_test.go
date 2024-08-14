package scenario_test

import (
	"context"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/cerberauth/vulnapi/internal/auth"
	"github.com/cerberauth/vulnapi/openapi"
	"github.com/cerberauth/vulnapi/scenario"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var server *http.Server

func TestMain(m *testing.M) {
	// Start the HTTP server
	server = &http.Server{
		Addr:    ":8080",
		Handler: http.DefaultServeMux,
	}
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Run the tests
	code := m.Run()

	// Shutdown the server
	err := server.Shutdown(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(code)
}

func TestNewOpenAPIScanWithHttpBearer(t *testing.T) {
	token := "token"
	doc, _ := openapi.LoadOpenAPI(context.Background(), "../test/stub/simple_http_bearer.openapi.json")
	securitySchemeValues := auth.NewSecuritySchemeValues(map[string]interface{}{
		"bearer_auth": &token,
	})

	s, err := scenario.NewOpenAPIScan(doc, securitySchemeValues, nil, nil)

	require.NoError(t, err)
	assert.Equal(t, 1, len(s.Operations))
	assert.Equal(t, "http://localhost:8080/", s.Operations[0].Request.URL.String())
	assert.Equal(t, http.MethodGet, s.Operations[0].Request.Method)
	assert.Equal(t, http.Header{}, s.Operations[0].Request.Header)
	assert.Equal(t, []auth.SecurityScheme{auth.NewAuthorizationBearerSecurityScheme("bearer_auth", &token)}, s.Operations[0].SecuritySchemes)
}

func TestNewOpenAPIScanWithJWTHttpBearer(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U"
	doc, _ := openapi.LoadOpenAPI(context.Background(), "../test/stub/simple_http_bearer_jwt.openapi.json")
	expectedSecurityScheme, _ := auth.NewAuthorizationJWTBearerSecurityScheme("bearer_auth", &token)
	securitySchemeValues := auth.NewSecuritySchemeValues(map[string]interface{}{
		"bearer_auth": &token,
	})

	s, err := scenario.NewOpenAPIScan(doc, securitySchemeValues, nil, nil)

	require.NoError(t, err)
	assert.Equal(t, 1, len(s.Operations))
	assert.Equal(t, "http://localhost:8080/", s.Operations[0].Request.URL.String())
	assert.Equal(t, http.MethodGet, s.Operations[0].Request.Method)
	assert.Equal(t, http.Header{}, s.Operations[0].Request.Header)
	assert.Equal(t, []auth.SecurityScheme{expectedSecurityScheme}, s.Operations[0].SecuritySchemes)
}

func TestNewOpenAPIScanWithMultipleOperations(t *testing.T) {
	gofakeit.Seed(1)

	token := "token"
	doc, _ := openapi.LoadOpenAPI(context.Background(), "../test/stub/basic_http_bearer.openapi.json")
	securitySchemes := []auth.SecurityScheme{auth.NewAuthorizationBearerSecurityScheme("bearer_auth", &token)}
	securitySchemeValues := auth.NewSecuritySchemeValues(map[string]interface{}{
		"bearer_auth": &token,
	})

	s, err := scenario.NewOpenAPIScan(doc, securitySchemeValues, nil, nil)

	require.NoError(t, err)
	assert.Equal(t, 2, len(s.Operations))
	for _, s := range s.Operations {
		assert.Equal(t, s.SecuritySchemes, securitySchemes)
	}
}

func TestNewOpenAPIScanWithoutParamsExample(t *testing.T) {
	gofakeit.Seed(1)

	token := "token"
	doc, _ := openapi.LoadOpenAPI(context.Background(), "../test/stub/basic_http_bearer.openapi.json")
	securitySchemes := []auth.SecurityScheme{auth.NewAuthorizationBearerSecurityScheme("bearer_auth", &token)}
	securitySchemeValues := auth.NewSecuritySchemeValues(map[string]interface{}{
		"bearer_auth": &token,
	})

	s, err := scenario.NewOpenAPIScan(doc, securitySchemeValues, nil, nil)

	require.NoError(t, err)
	assert.Equal(t, 2, len(s.Operations))
	for _, s := range s.Operations {
		assert.Equal(t, s.SecuritySchemes, securitySchemes)
	}
}
