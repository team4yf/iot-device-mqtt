package hooks

import (
	"fmt"
	"strings"
	"time"
	"github.com/team4yf/iot-device-mqtt/pkg/utils"
	"github.com/team4yf/yf-fpm-server-go/fpm"
)

const (
	cmdFfmpeg = "ffmpeg -i rtsp://%s:%s@%s:554/h264/1/sub/av_stream -an -f mpegts -codec:v mpeg1video -s 640x480 -b:v 100k -bf 0 -muxdelay 0.001 http://open.yunplus.io:18081/fpmpassword/%s"
	cmdOnvif  = "onvif-ptz %s --baseUrl=http://%s:80 -u=%s -p=%s -x=%f -y=%f -z=%f"
	cmdNmap   = `nmap -p %s localhost | grep -E "^[1-9]" | awk '{print $1","$2 }' | sed "s/\/tcp//" | sed "s/\/udp//"`
	cmdPs     = "ps -ef | grep %s"
	cmdLsof   = "lsof -i:%d"
	cmdIP     = "ip addr | grep 'inet' | grep -v '127.0.0.1' | grep -v 'inet6' | cut -d: -f2 | awk '{print $2}' | head -1 | awk -F / '{print $1}'"
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

type beatBody struct {
	LocalIP   string          `json:"localIP"`
	GatewayID string          `json:"gatewayID"`
	TimeStamp int64           `json:"timestamp"`
	Cameras   map[string]bool `json:"cameras"`
}

func interval(cameras []string) (beatMessage *beatBody, err error) {
	beatMessage = &beatBody{
		LocalIP: utils.GetLocalIP(),
		// TODO: read from the env/arg
		GatewayID: "demo",
		TimeStamp: time.Now().Unix(),
		Cameras:   make(map[string]bool, len(cameras)),
	}
	for _, cameraIP := range cameras {
		beatMessage.Cameras[cameraIP] = false
		out, err := utils.RunCmd(fmt.Sprintf(`nmap -p 554 %s | grep -E "^[1-9]" | awk '{print $1","$2 }' | sed "s/\/tcp//" | sed "s/\/udp//" | grep open| wc -l`, cameraIP))
		if err != nil {
			return nil, err
		}
		count := strings.Trim((string)(out), " \n")
		beatMessage.Cameras[cameraIP] = (count != "0")
	}
	return
}

func runCommand(app *fpm.Fpm, execute *executeBody) {
	finalCommand := ""
	switch execute.Command {
	case "ffmpeg":
		//check exists
		out, err := utils.RunCmd("ps -ef | grep ffmpeg | grep " + execute.Argument[len(execute.Argument)-1].(string) + ` | grep -v "grep" |wc -l`)
		if err == nil {
			//count the ffmpeg process instance
			count := strings.Trim((string)(out), " \n")
			fmt.Printf("count=|%s|\n", count)
			if count != "0" {
				// exists
				return
			}
		}
		finalCommand = fmt.Sprintf(cmdFfmpeg, execute.Argument...)
	case "onvif":
		finalCommand = fmt.Sprintf(cmdOnvif, execute.Argument...)
	case "nmap":
		finalCommand = fmt.Sprintf(cmdNmap, execute.Argument...)
	case "ps":
		finalCommand = fmt.Sprintf(cmdPs, execute.Argument...)
	case "lsof":
		finalCommand = fmt.Sprintf(cmdLsof, execute.Argument...)
	case "ip":
		finalCommand = cmdIP
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
	go app.Execute("mqttclient.publish", &fpm.BizParam{
		"topic": "$d2s/yunplus/ipc/feedback",
		"payload": ([]byte)(feedbackStr),
	})
}

//ConsumerHook the hook of the consumer
//it will run after init,
//make mqtt connection
func ConsumerHook(app *fpm.Fpm) {

	app.Execute("mqttclient.subscribe", &fpm.BizParam{
		"topics": "$s2d/yunplus/ipc/demo/config",
	})
	// the demo is the device id
	//{ "commad": "ps","argument":["vscode"],"messageID":"123", "feedback": 1}
	app.Subscribe("$s2d/+/ipc/demo/execute", func(topic string, payload interface{}) {
		fmt.Println(topic, (string)(payload.([]byte)))
		//TODO: here
		execute := executeBody{}
		if err := utils.DataToStruct(payload.([]byte), &execute); err != nil {
			fmt.Println("convert the execute message error:", err)
			return
		}
		go runCommand(app, &execute)
	})

	cameras := []string{"192.168.0.108"}

	app.Subscribe("$s2d/yunplus/ipc/demo/config", func(topic string, payload interface{}) {
		fmt.Println(topic, (string)(payload.([]byte)))
		
		conf := make(map[string]interface{})
		if err := utils.DataToStruct(payload.([]byte), &conf); err != nil {
			fmt.Println(err)
			return
		}
		if cameraList, ok := conf["cameras"]; ok {
			cameras = make([]string, 0)
			for _, c := range cameraList.([]interface{}) {
				cameras = append(cameras, c.(string))
			}

		}
		
	})
	//auto push beat info
	t := time.NewTicker(10 * time.Second)
	defer t.Stop()
	for {

		select {

		case <-t.C:
			if beatMessage, err := interval(cameras); err != nil {
				fmt.Println(err)
			} else {
				app.Execute("mqttclient.publish", &fpm.BizParam{
					"topic": "$d2s/yunplus/ipc/beat",
					"payload": utils.Struct2Bytes(beatMessage),
				})
			}

		}

	}
}
