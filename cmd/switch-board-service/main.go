package main

import (
    "bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/util"
    "github.com/Sirupsen/logrus"
    "os"
)

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


    }
}