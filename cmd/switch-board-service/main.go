package main

import (
	"bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/api"
	"bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/models"
	"bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/util"
	"encoding/gob"
	"flag"
	"github.com/getsentry/sentry-go"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"html/template"
	"os"
	"strconv"
)

func connectToDatabase(config *util.Configuration) (*gorm.DB, error) {

	database, databaseError := gorm.Open("mysql", config.DatabaseUsername+":"+config.DatabasePassword+"@/"+config.DatabaseName+"?charset=utf8&parseTime=True&loc=Local")

	if databaseError != nil {
		return nil, databaseError
	}

	database.AutoMigrate(&models.UserModel{})
	database.AutoMigrate(&models.ExchangeModel{})
	database.AutoMigrate(&models.TradePairModel{})
	database.AutoMigrate(&models.ServiceKeyModel{})

	// Check to see if the default user exit
	util.CheckForDefaultUser(database, config)

	return database, nil

}

func init() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}

// @title Alkemi - Switch Board Service
// @version 1.0
// @description This service manages the authorization for exchanges and trade pairs

// @contact.name Raynaldo Rivera
// @contact.url https://alkemi.ai
// @contact.email royal@alkemi.ai

// @license.name Proprietary
// @license.url

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @host 127.0.0.1:7020
// @BasePath /

func main() {
	util.ShowBanner()

	prod := flag.Bool("prod", false, "Use '$go run ./cmd/... -prod' to set to true, changing default configuration file from 'dev.json' to 'config.json'")
	flag.Parse()

	var configFilename string
	if *prod {
		configFilename = "config.json"
	} else {
		configFilename = "dev.json"
	}

	config, configLoadError := util.LoadConfigurationFile(configFilename)
	if configLoadError != nil {
		logrus.Fatalf("failed to load app configuration; %v", configLoadError)
		os.Exit(1)
	}
	if config.LogLevel == 0 {
		logrus.SetLevel(logrus.InfoLevel)
		logrus.Info("Log Level: Info")
	} else {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Info("Log Level: Debug")
	}

	sentrySetDSNError := sentry.Init(sentry.ClientOptions{
		Dsn: config.SentryUrl,
	})

	if sentrySetDSNError != nil {
		logrus.WithFields(logrus.Fields{}).Error(sentrySetDSNError)
	}

	databaseConnection, databaseError := connectToDatabase(config)
	if databaseError != nil {
		logrus.Errorf("failed to connect to database; %v", databaseError, nil)
		sentry.CaptureException(databaseError)
		os.Exit(1)
	} else {
		logrus.Info(util.StrDatabaseConnect)
		sentry.CaptureMessage(util.StrDatabaseConnect)
	}

	// Create the data context
	dataContext := util.DataContext{}

	dataContext.Config = config
	dataContext.Database = databaseConnection
	dataContext.Errors = &util.ErrorsStruct{}


	e := echo.New()
	e.Debug = true

	// Register our custom app context
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {

		return func(c echo.Context) error {
			cc := &util.EndpointContext{Context: c, AppContext: &dataContext}
			return h(cc)
		}
	})

	e.Use(session.Middleware(sessions.NewCookieStore([]byte(uuid.NewV4().String()))))

	e.Use(middleware.Logger())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"authorization", "Content-Type"},
		AllowCredentials: true,
		AllowMethods:     []string{echo.OPTIONS, echo.GET, echo.POST},
	}))

	templates := make(map[string]*template.Template)

	e.Renderer = &util.TemplateRegistry{
		Templates: templates,
	}

	gob.Register(models.UserModel{})

	e.HideBanner = true

	// Public
	e.POST("/v1/login", api.LoginAction)
	e.POST("/v1/generate_new_temp_token", api.GenerateNewTempTokenAction)
	e.POST("/v1/revoke_temp_access_token", api.RevokeTempAccessTokenAction)
	e.POST("/v1/retrieve_exchange_by_temp_token", api.RetrieveExchangeInfoByExchangeTempTokenAction)

	// Services
	e.POST("/v1/validate_access", api.ValidateAccessAction)
	e.POST("/v1/service/service_retrieve_all_exchanges", api.ServiceRetrieveAllExchangesAction)
	e.POST("/v1/service/service_retrieve_all_trade_pairs", api.ServiceRetrieveAllTradePairsAction)

	//Health
	go api.ContinuousHealthCheck(&dataContext)
	e.GET("/v1/health", api.HealthCheckAction)

	r := e.Group("/v1/action")

	r.Use(middleware.JWT([]byte(config.JWTSecret)))

	// Exchanges
	r.POST("/create_exchange", api.CreateExchangeAction)
	r.GET("/remove_exchange", api.RemoveExchangeAction)
	r.GET("/authorize_exchange", api.AuthorizeExchangeAction)
	r.GET("/retrieve_all_exchanges", api.RetrieveAllExchangesAction)

	// Trade pair
	r.POST("/create_trade_pair", api.CreateTradePairAction)
	r.GET("/remove_trade_pair", api.RemoveTradePairAction)
	r.GET("/retrieve_all_trade_pairs", api.RetrieveAllTradePairsAction)
	r.POST("/create_exchange_trade_pair_association", api.CreateExchangeTradePairsAssociationAction)
	r.POST("/remove_exchange_trade_pair_association", api.RemoveExchangeTradePairsAssociationAction)

	// Service Keys
	r.POST("/create_service_key", api.CreateServiceKeyAction)
	r.GET("/remove_service_key", api.RemoveServiceKeyAction)
	r.GET("/retrieve_all_service_keys", api.RetrieveAllServiceKeysAction)

	e.Logger.Fatal(e.Start(config.ListenIp + ":" + strconv.Itoa(config.ListenPort)))

	databaseCloseError := dataContext.Database.Close()

	if databaseCloseError != nil {
		sentry.CaptureException(databaseCloseError)
	}
}
