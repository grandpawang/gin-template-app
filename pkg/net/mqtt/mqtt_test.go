package mqtt

import (
	"context"
	"fmt"
	"testing"
	"time"

	pahoMqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/smartystreets/goconvey/convey"
)

func TestShareSubscribe(t *testing.T) {
	cfg := Config{
		URL: "ws://coint.pgy:8083/mqtt",
	}
	mqttCli1 := New(&cfg)
	mqttCli2 := New(&cfg)
	mqttCli3 := New(&cfg)
	mqttCli4 := New(&cfg)
	convey.Convey("TestShareSubscribe", t, func(ctx convey.C) {

		ctx.Convey("$share/<group>/`topic`", func(ctx convey.C) {
			var (
				count = 0
			)
			mqttCli1.Listen("$share/test/hello", 1, func(ctx context.Context, cli pahoMqtt.Client, msg pahoMqtt.Message) (err error) {
				fmt.Println("mqtt cli1 get message")
				count++
				return
			})

			mqttCli2.Listen("$share/test/hello", 1, func(ctx context.Context, cli pahoMqtt.Client, msg pahoMqtt.Message) (err error) {
				fmt.Println("mqtt cli2 get message")
				count++
				return
			})

			mqttCli3.Listen("$share/test/hello", 1, func(ctx context.Context, cli pahoMqtt.Client, msg pahoMqtt.Message) (err error) {
				fmt.Println("mqtt cli3 get message")
				count++
				return
			})
			ctx.Println(`publish test_queue message: "{ "msg": "Hello, World!" }`)
			mqttCli4.Publish("hello", 1, false, `{ "msg": "Hello, World!" }`)
			mqttCli4.Publish("hello", 1, false, `{ "msg": "Hello, World!" }`)
			mqttCli4.Publish("hello", 1, false, `{ "msg": "Hello, World!" }`)
			mqttCli4.Publish("hello", 1, false, `{ "msg": "Hello, World!" }`)
			time.Sleep(1 * time.Second)
			convey.So(count, convey.ShouldEqual, 4)
		})

		ctx.Convey("$queue/`topic`", func(ctx convey.C) {
			var (
				count = 0
			)
			mqttCli1.Listen("$queue/topic", 1, func(ctx context.Context, cli pahoMqtt.Client, msg pahoMqtt.Message) (err error) {
				fmt.Println("mqtt cli1 get message")
				count++
				return
			})

			mqttCli2.Listen("$queue/topic", 1, func(ctx context.Context, cli pahoMqtt.Client, msg pahoMqtt.Message) (err error) {
				fmt.Println("mqtt cli2 get message")
				count++
				return
			})

			mqttCli3.Listen("$queue/topic", 1, func(ctx context.Context, cli pahoMqtt.Client, msg pahoMqtt.Message) (err error) {
				fmt.Println("mqtt cli3 get message")
				count++
				return
			})
			ctx.Println(`publish test_queue message: "{ "msg": "Hello, World!" }`)
			mqttCli4.Publish("topic", 1, false, `{ "msg": "Hello, World!" }`)
			mqttCli4.Publish("topic", 1, false, `{ "msg": "Hello, World!" }`)
			mqttCli4.Publish("topic", 1, false, `{ "msg": "Hello, World!" }`)
			mqttCli4.Publish("topic", 1, false, `{ "msg": "Hello, World!" }`)
			time.Sleep(1 * time.Second)
			convey.So(count, convey.ShouldEqual, 4)
		})
	})
}
