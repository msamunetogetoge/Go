package main

import (
	"goweather/app/controllers"
	"goweather/config"
	"goweather/utils"
)

func main() {
	utils.LoggingSettings(config.Config.LogFile)
	controllers.StartWebServer()
}
