package template

import (
	"gin-template-app/pkg/net/mqtt"
	"gin-template-app/service/template"

	"github.com/gin-gonic/gin"
	jsoniterator "github.com/json-iterator/go"
)

var json = jsoniterator.ConfigCompatibleWithStandardLibrary

var (
	templateSvr *template.Service
	mqttCli     *mqtt.Client
)

// Init template route
func Init(hws *template.Service, mqttcli *mqtt.Client) {
	templateSvr = hws
	mqttCli = mqttcli
}

// Route template route
func Route(engine *gin.Engine) {
	// register router
}

// MqttRoute .
func MqttRoute() {

}
