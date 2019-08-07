package main

import (
	"bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/api"
	"bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/models"
	"bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/util"
	"bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/views"
	"encoding/gob"
	"github.com/Sirupsen/logrus"
	"github.com/getsentry/raven-go"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/middleware"
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

	// Check to see if the default user exit
	util.CheckForDefaultUser(database, config)

	return database, nil

}

func init() {

	logrus.SetOutput(os.Stdout)

	logrus.SetLevel(logrus.InfoLevel)
}

func main() {

	util.ShowBanner()

	config, configLoadError := util.LoadConfigurationFile()

	_ = config

	if configLoadError == nil {

		// Set the logrus config level
		if config.LogLevel == 0 {
			logrus.SetLevel(logrus.InfoLevel)
			logrus.Info("Log Level: Info")
		} else {
			logrus.SetLevel(logrus.DebugLevel)
			logrus.Info("Log Level: Debug")

		}

		databaseConnection, databaseError := connectToDatabase(&config)

		if databaseError != nil {
			logrus.Error(databaseError, nil)
			return
		} else {
			logrus.Info("Successfully connected to database")
		}

		// Create the data context
		dataContext := util.DataContext{}

		dataContext.Database = databaseConnection

		sentrySetDSNError := raven.SetDSN(config.SentryUrl)

		if sentrySetDSNError != nil {
			logrus.WithFields(logrus.Fields{}).Error(sentrySetDSNError)
		}

		e := echo.New()
		e.Debug = true

		// Register our custom app context
		e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				cc := &util.EndpointContext{Context: c, AppContext: &dataContext}
				return h(cc)
			}
		})

		e.Static("/static", "assets")

		e.Use(session.Middleware(sessions.NewCookieStore([]byte("herpderp"))))
		e.Use(middleware.Logger())

		templates := make(map[string]*template.Template)

		templates["login.html"] = template.Must(template.ParseFiles("templates/login.html", "templates/base.html"))
		templates["dashboard.html"] = template.Must(template.ParseFiles("templates/dashboard.html", "templates/base.html"))

		e.Renderer = &util.TemplateRegistry{
			Templates: templates,
		}

		gob.Register(models.UserModel{})

		e.HideBanner = true

		e.GET("/", views.RedirectView)
		e.GET("/login", views.LoginView)
		e.GET("/dashboard", views.DashboardView)

		e.GET("/action/logout", api.LogoutAction)
		e.POST("/action/login", api.LoginAction)

		e.Logger.Fatal(e.Start(config.ListenIp + ":" + strconv.Itoa(config.ListenPort)))

		databaseCloseError := dataContext.Database.Close()

		if databaseCloseError != nil {
			raven.CaptureError(databaseCloseError, nil)
		}

	}
}
