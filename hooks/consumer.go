package hooks

import (
	"fmt"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/team4yf/iot-device-mqtt/pkg/pubsub"
	"github.com/team4yf/iot-device-mqtt/pkg/utils"
	"github.com/team4yf/yf-fpm-server-go/fpm"
)

const (
	cmdFfmpeg = "ffmpeg "
	cmdOnvif  = "onvif"
	cmdNmap   = "nmap"
	cmdPs     = "ps -ef | grep %s"
	cmdLsof   = "lsof -i:%d"
)

type executeBody struct {
	Command   string        `json:"command"`
	MessageID string        `json:"messageID"`
	Argument  []interface{} `json:"argument"`
	Feedback  int           `json:"feedback"`
}

type feedbackBody struct {
	MessageID string `json:"messageID"`
	DeviceID  string `json:"deviceID"`
	Code      int    `json:"code"`
	Error     string `json:"error,omitempty"`
	Data      string `json:"data,omitempty"`
}

func runCommand(fpm *fpm.Fpm, mq *pubsub.PubSub, execute *executeBody) {
	finalCommand := ""
	switch execute.Command {
	case "ffmpeg":
		finalCommand = fmt.Sprintf(cmdFfmpeg, execute.Argument...)
	case "onvif":
		finalCommand = fmt.Sprintf(cmdOnvif, execute.Argument...)
	case "nmap":
		finalCommand = fmt.Sprintf(cmdNmap, execute.Argument...)
	case "ps":
		finalCommand = fmt.Sprintf(cmdPs, execute.Argument...)
	case "lsof":
		finalCommand = fmt.Sprintf(cmdLsof, execute.Argument...)
	}

	out, err := utils.RunCmd(finalCommand)
	if execute.Feedback == 0 {
		if err != nil {
			fmt.Printf("run command error: %s, error:\n %v\n", finalCommand, err)
			return
		}
		fmt.Printf("run command success: %s, out:\n %s\n", finalCommand, (string)(out))
		return
	}

	feedback := feedbackBody{
		MessageID: execute.MessageID,
		DeviceID:  "demo",
		Code:      -1,
	}
	if err != nil {
		feedback.Error = err.Error()
		feedback.Code = -9
	} else {
		feedback.Code = 0
		feedback.Data = (string)(out)
	}
	feedbackStr := utils.JSON2String(feedback)
	go (*mq).Publish("$d2s/aa/ipc/feedback", ([]byte)(feedbackStr))
}

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

	// the demo is the device id
	//{ "commad": "ps","argument":["vscode"],"messageID":"123", "feedback": 1}
	mq.Subscribe("$s2d/+/ipc/demo/execute", func(topic, payload interface{}) {
		fmt.Println(topic, (string)(payload.([]byte)))
		//TODO: here
		execute := executeBody{}
		if err := utils.DataToStruct(payload.([]byte), &execute); err != nil {
			fmt.Println("convert the execute message error:", err)
			return
		}
		go runCommand(fpm, &mq, &execute)
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
