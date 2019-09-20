package api

import (
	"net/http"
	"github.com/labstack/echo/v4"
)


// Generalized Health Check endpoint
// @Summary Validates that service is up and running
// @Description Health Check for service
// @Accept ---TBD---
// @Produce http.response code
// @Tags  Health Check
// @Param ---TBD---
// @Success 200
// @Failure 503
// @Router /v1/health [get]

func HealthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "Health Check")
}
