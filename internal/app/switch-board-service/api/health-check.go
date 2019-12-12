package api

import (
	"bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/responses"
	"bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/util"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// Generalized Health Check endpoint
// @Summary Validates that service is up and running
// @Description Health Check for service
// @Accept  json
// @Produce  json
// @Tags Exchange
// @Tags  Health Check
// @Success 200 {object} responses.HealthCheckResponse ""
// @Failure 200 {object} responses.HealthCheckResponse ""
// @Router /v1/health [get]
func HealthCheckAction(c echo.Context) error {
	cc := c.(*util.EndpointContext)

	err := cc.AppContext.Database.DB().Ping()

	if err != nil {
		cc.AppContext.Errors.DBErrors = cc.AppContext.Errors.DBErrors + 1

		return c.JSON(http.StatusOK, &responses.HealthCheckResponse{CurrentlyHealthy: false,
			AuthenticationErrors: cc.AppContext.Errors.AuthErrors, DatabaseConnectionErrors: cc.AppContext.Errors.DBErrors})

	}

	return c.JSON(http.StatusOK, &responses.HealthCheckResponse{CurrentlyHealthy: true,
		AuthenticationErrors: cc.AppContext.Errors.AuthErrors, DatabaseConnectionErrors: cc.AppContext.Errors.DBErrors})
}

func ContinuousHealthCheck(dc *util.DataContext) {

	logrus.Info("Starting Health-Check routine.")
	for {
		time.Sleep(time.Duration(dc.Config.HealthCheckFrequency) * time.Second)

		err := dc.Database.DB().Ping()

		if err != nil {
			dc.Errors.DBErrors = dc.Errors.DBErrors + 1

			sentry.CaptureException(err)
			logrus.Error("Error when connecting to database")

		}
	}
}
