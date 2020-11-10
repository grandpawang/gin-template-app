package template

import (
	"gin-template-app/pkg/net/mqtt"
)

func getTopic(mac string) string {
	return "/dev/mw/" + mac
}

// MQTTOrder mqtt publish 命令
type MQTTOrder struct {
	mqtt *mqtt.Client
}

// NewMQTTOrder 新建mqtt命令对象
func NewMQTTOrder(cfg *mqtt.Config) *MQTTOrder {
	return &MQTTOrder{mqtt: mqtt.New(cfg)}
}
