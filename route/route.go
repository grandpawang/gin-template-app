package route

import (
	"flag"
	"gin-template-app/conf"
	"gin-template-app/pkg/net/http"
	"gin-template-app/pkg/net/mqtt"
	templateRoute "gin-template-app/route/template"
	templateService "gin-template-app/service/template"

	// _ "net/http/pprof" // pprof

	"github.com/gin-gonic/gin"
)

var iconPath string

func init() {
	flag.StringVar(&iconPath, "icon", "./cmd/favicon.ico", "default icon path")
}

// Init http server
func Init(temp *templateService.Service, c conf.Config) {
	mqttCli := mqtt.New(c.MQTT)
	engine := http.New()
	templateRoute.Init(temp, mqttCli)
	route(engine)
	listen()
	http.Init(c.HTTP, engine)
}

// Setup routers
func route(engine *gin.Engine) {
	engine.StaticFile("/favicon.ico", iconPath)
	templateRoute.Route(engine)
}

// setup mqtt route
func listen() {
	templateRoute.MqttRoute()
}
