package api

import (
	"net/http"
	"github.com/labstack/echo/v4"
	"fmt"
	"bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/util"
	"time"
	"github.com/sirupsen/logrus"
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
	cc := c.(*util.EndpointContext)
	ccc := *cc.AppContext
	status := fmt.Sprintf("DB Connection Healthy, Errors: %d", ccc.DBErrors)
	if err := ccc.Database.DB().Ping(); err != nil {
		ccc.DBErrors++
		status = fmt.Sprintf("DB Connection Unhealthy, Errors: %d", ccc.DBErrors)
		return c.String(503, status)
	}
	return c.String(http.StatusOK, status)
}

func ContinuousHealthCheck(seconds int, dc *util.DataContext) {
	for {
		time.Sleep(time.Duration(seconds) * time.Second)
		_, err := http.Get("http://127.0.0.1:7020/v1/health")
		if err != nil {
			dc.DBErrors++
			logrus.Info("err\n\n\n")
			logrus.Info(err)
		}
	}
}
