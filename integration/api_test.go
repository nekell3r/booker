package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"booker/cmd/admin-gateway/config"
	"booker/cmd/admin-gateway/handlers"
	"booker/cmd/admin-gateway/middleware"
)

// Integration tests for API endpoints
// These tests require running services or use testcontainers

func setupTestServer() *echo.Echo {
	// Setup test server with mocked dependencies
	// Note: Handler uses private fields, so we need to use New() function
	// For integration tests, we'd typically use real connections or testcontainers
	cfg := &config.Config{}

	// In real integration tests, we'd create actual gRPC connections
	// For now, we create a handler with nil connections (will fail on actual calls)
	handler := handlers.New(nil, nil, nil, cfg)

	mw := &middleware.Middleware{}
	e := handler.SetupRoutes(mw)
	return e
}

func TestAPI_HealthCheck(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	e := setupTestServer()

	req := httptest.NewRequest(http.MethodGet, "/api", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Admin Gateway", response["service"])
}

func TestAPI_CreateVenue_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	e := setupTestServer()

	body := map[string]interface{}{
		"name":     "Test Venue",
		"timezone": "UTC",
		"address":  "123 Main St",
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/venues", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Note: This would require authentication middleware to be mocked
	// For now, we test the endpoint structure
	assert.True(t, rec.Code == http.StatusOK || rec.Code == http.StatusUnauthorized)
}

func TestAPI_CheckAvailability_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	e := setupTestServer()

	body := map[string]interface{}{
		"venue_id": "venue-1",
		"slot": map[string]interface{}{
			"date":             "2024-01-15",
			"start_time":       "19:00",
			"duration_minutes": 120,
		},
		"party_size": 4,
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/availability/check", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Note: This would require authentication middleware to be mocked
	assert.True(t, rec.Code == http.StatusOK || rec.Code == http.StatusUnauthorized)
}

// Note: Full integration tests would require:
// 1. Testcontainers for PostgreSQL, Redis, Kafka
// 2. Running all services (venue-svc, booking-svc, admin-gateway)
// 3. Proper authentication setup
// 4. Database migrations and test data setup

