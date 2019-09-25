package api

import (
	"bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/responses"
	"bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/util"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
	"strconv"
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
func HealthCheckAction(c echo.Context) error {
	cc := c.(*util.EndpointContext)

	err := cc.AppContext.Database.DB().Ping()

	if err != nil {
		cc.AppContext.Errors.DBErrors++

		return c.JSON(http.StatusOK, &responses.HealthCheckResponse{CurrentlyHealthy: false,
			AuthenticationErrors: cc.AppContext.Errors.AuthErrors, DatabaseConnectionErrors: cc.AppContext.Errors.DBErrors})

	}

	return c.JSON(http.StatusOK, &responses.HealthCheckResponse{CurrentlyHealthy: true,
		AuthenticationErrors: cc.AppContext.Errors.AuthErrors, DatabaseConnectionErrors: cc.AppContext.Errors.DBErrors})
}

func ContinuousHealthCheck(dc *util.DataContext) {
	dbHostname := "http://"+dc.Config.ListenIp+":"+strconv.Itoa(dc.Config.ListenPort)+"/v1/health"
	logrus.Info("starting heath checks at: "+dbHostname)
	for {
		time.Sleep(time.Duration(dc.Config.HealthCheckFrequency) * time.Second)

		healthCheckResponse, healthCheckResponseError := http.Get(dbHostname)

		if healthCheckResponseError != nil {
			dc.Errors.DBErrors++
		} else {

			var healthCheckResponseModel responses.HealthCheckResponse

			decodeError := json.NewDecoder(healthCheckResponse.Body).Decode(&healthCheckResponseModel)

			if decodeError != nil {
				logrus.Error(util.StrHealthCheckFailed)
				dc.Errors.DBErrors++
			} else {

				if healthCheckResponseModel.CurrentlyHealthy != true {
					dc.Errors.DBErrors++
				}
			}
		}
	}
}
