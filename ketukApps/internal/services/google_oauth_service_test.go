package services

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"ketukApps/config"

	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func TestNewGoogleOAuthService(t *testing.T) {
	cfg := &config.Config{}
	cfg.Google.ClientID = "test-client-id"
	cfg.Google.ClientSecret = "test-client-secret"
	cfg.Google.RedirectURI = "http://localhost:8080/callback"

	service := NewGoogleOAuthService(cfg)

	assert.NotNil(t, service)
	assert.Equal(t, "test-client-id", service.config.ClientID)
	assert.Equal(t, "test-client-secret", service.config.ClientSecret)
	assert.Equal(t, "http://localhost:8080/callback", service.config.RedirectURL)
}

func TestGetAuthURL(t *testing.T) {
	cfg := &config.Config{}
	service := NewGoogleOAuthService(cfg)
	state := "random-state"

	url := service.GetAuthURL(state)

	assert.Contains(t, url, "response_type=code")
	assert.Contains(t, url, "client_id=")
	assert.Contains(t, url, "state="+state)
	assert.Contains(t, url, "scope=openid+email+profile")
}

func TestExchangeCode(t *testing.T) {
	// Mock Google Token Endpoint
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/token" {
			t.Errorf("Expected path /token, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"access_token": "mock-access-token", "token_type": "Bearer", "expires_in": 3600}`))
	}))
	defer server.Close()

	cfg := &config.Config{}
	service := NewGoogleOAuthService(cfg)

	// Override Endpoint to point to mock server
	service.config.Endpoint = oauth2.Endpoint{
		TokenURL: server.URL + "/token",
	}

	ctx := context.Background()
	// Need to provide a client in context if we were testing generic oauth, but here we redirect the endpoint.
	// oauth2 Exchange uses DefaultClient if context doesn't have one. DefaultClient works with localhost.

	token, err := service.ExchangeCode(ctx, "valid-code")

	assert.NoError(t, err)
	assert.NotNil(t, token)
	assert.Equal(t, "mock-access-token", token.AccessToken)
}

// customTransport allows us to intercept HTTP requests
type customTransport struct {
	roundTrip func(*http.Request) (*http.Response, error)
}

func (t *customTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.roundTrip(req)
}

func TestGetUserInfo(t *testing.T) {
	cfg := &config.Config{}
	service := NewGoogleOAuthService(cfg)

	// Mock token
	token := &oauth2.Token{
		AccessToken: "mock-access-token",
		TokenType:   "Bearer",
		Expiry:      time.Now().Add(time.Hour),
	}

	// Create a context with a custom HTTP client
	mockTransport := &customTransport{
		roundTrip: func(req *http.Request) (*http.Response, error) {
			if req.URL.String() == "https://www.googleapis.com/oauth2/v2/userinfo" {
				resp := GoogleUserInfo{
					ID:            "123456789",
					Email:         "test@example.com",
					VerifiedEmail: true,
					Name:          "Test User",
				}
				bodyBytes, _ := json.Marshal(resp)

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(bodyBytes)),
					Header:     make(http.Header),
				}, nil
			}
			return nil, http.ErrNotSupported
		},
	}

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{Transport: mockTransport})

	userInfo, err := service.GetUserInfo(ctx, token)

	assert.NoError(t, err)
	assert.NotNil(t, userInfo)
	assert.Equal(t, "test@example.com", userInfo.Email)
	assert.Equal(t, "Test User", userInfo.Name)
	assert.True(t, userInfo.VerifiedEmail)
}

func TestGetUserInfo_Error(t *testing.T) {
	cfg := &config.Config{}
	service := NewGoogleOAuthService(cfg)

	token := &oauth2.Token{AccessToken: "invalid-token"}

	mockTransport := &customTransport{
		roundTrip: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusUnauthorized,
				Body:       io.NopCloser(strings.NewReader(`{"error": "unauthorized"}`)),
				Header:     make(http.Header),
			}, nil
		},
	}

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{Transport: mockTransport})

	userInfo, err := service.GetUserInfo(ctx, token)

	assert.Error(t, err)
	assert.Nil(t, userInfo)
	assert.Contains(t, err.Error(), "status 401")
}
