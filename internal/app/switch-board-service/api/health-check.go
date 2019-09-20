package api

import (
	"net/http"
	"github.com/labstack/echo/v4"
//	"context"
//	"database/sql"
//	"time"
	"github.com/jinzhu/gorm"
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

//func HealthCheck(c echo.Context, db *gorm.DB) echo.HandlerFunc {
//	status := "DB Connection Healthy"
	//if err := db.DB().Ping(); err != nil {
	//	status = "DB Connection Unhealthy"
	//	return echo.HandlerFunc(func(c echo.Context) {
	//		return	c.String(503, status)
	//	}())
	//}
//	return func(c *echo.Context) error {
//		return echo.NewHTTPError(http.StatusOK, status)
//	}(*c)
//	return func(c *echo.Context) error {
//		return echo.NewHTTPError(http.StatusUnauthorized)
//	}
//}
func HealthCheck(db *gorm.DB) echo.HandlerFunc {
    return func(c echo.Context) error {
	status := "DB Connection Healthy"
	if err := db.DB().Ping(); err != nil {
		status = "DB Connection Unhealthy"
		return c.String(503, status)
	}
        return c.String(http.StatusOK, status)
    }
}
