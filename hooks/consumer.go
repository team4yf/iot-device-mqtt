package hooks

import (
	"fmt"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/team4yf/iot-device-mqtt/pkg/pubsub"
	"github.com/team4yf/yf-fpm-server-go/fpm"
	"github.com/team4yf/yf-fpm-server-go/pkg/utils"
)

//ConsumerHook the hook of the consumer
//it will run after init,
//make mqtt connection
func ConsumerHook(fpm *fpm.Fpm) {
	//fpm.GetConfig("mqtt") , it's not workiong.
	setting := &pubsub.MqttSetting{
		Options:  &MQTT.ClientOptions{},
		Retained: false,
		Qos:      (byte)(0),
	}

	setting.Options.AddBroker(fmt.Sprintf("tcp://%s:%d",
		"mqtt.yunplus.io",
		1883))
	setting.Options.SetClientID("iot-device-" + utils.GenUUID())
	setting.Options.SetUsername("fpmuser")
	setting.Options.SetPassword("fpmpassword")

	mq := pubsub.NewMQTTPubSub(setting)
	fmt.Println("mqtt client inited!")

	mq.Subscribe("$s2d/+/device/execute", func(topic, payload interface{}) {
		fmt.Println(topic, (string)(payload.([]byte)))
		//TODO: here

	})

	//auto push beat info
	t := time.NewTicker(10 * time.Second)
	defer t.Stop()
	for {

		select {

		case <-t.C:

			// fmt.Println(time.Now())
			mq.Publish("$d2s/aa/bb/beat", ([]byte)(`{"Status":"UP"}`))

		}

	}
}
