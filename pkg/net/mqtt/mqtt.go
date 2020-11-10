package mqtt

import (
	"context"
	"encoding/json"
	"gin-template-app/pkg/ecode"
	"gin-template-app/pkg/log"
	"time"

	pahoMqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	// ServerOnline server online channel
	ServerOnline = "OL"
)

// Config .
type Config struct {
	URL   string
	Debug bool
}

// Client mqtt client
type Client struct {
	Client pahoMqtt.Client
	cb     func(cli *Client, err error)
}

// Handler message handler
type Handler func(ctx context.Context, cli pahoMqtt.Client, msg pahoMqtt.Message) error

// New pahoMqtt
func New(cfg *Config) *Client {
	opts := pahoMqtt.NewClientOptions().AddBroker(cfg.URL)
	opts.SetAutoReconnect(true) //自动重连
	opts.SetWill(ServerOnline, "", 1, true)
	cli := &Client{
		cb:     func(cli *Client, err error) {},
		Client: pahoMqtt.NewClient(opts),
	}
	if token := cli.Client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	if cfg.Debug {
		cli.cb = debug
	}
	log.Infoln("🚀 MQTT Client Connected!")
	cli.Publish(ServerOnline, 1, true, "1")
	return cli
}

// Publish pahoMqtt pulish message payload can be string or []bytes
func (cli *Client) Publish(topic string, qos byte, retained bool, payload interface{}) error {
	return cli.Client.Publish(topic, qos, retained, payload).Error()
}

// Listen pahoMqtt Subscribe topic
func (cli *Client) Listen(topic string, qos byte, callback Handler) {
	log.Info("🚀 MQTT Client Listen: %v\n", topic)
	token := cli.Client.Subscribe(topic, qos, func(client pahoMqtt.Client, msg pahoMqtt.Message) {
		now := time.Now()
		ctx := context.Background()
		err := callback(ctx, client, msg)
		log.Info("[MQTT] %v|%v|%v|", now.Format("0000-00-00 00:00:00"), time.Since(now), topic)
		if err != nil {
			log.Error("handle error(%v)", err)
		}
		cli.cb(cli, err)
	})
	token.Wait()
}

// log nil
func debug(cli *Client, err error) {
	bcode := ecode.Cause(err)
	data := struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}{
		Code:    bcode.Code(),
		Message: bcode.Message(),
	}
	bs, _ := json.Marshal(data)
	cli.Publish("/server/response", 1, false, bs)
}
