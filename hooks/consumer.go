package hooks

import (
	"fmt"
	"strings"
	"time"

	"github.com/team4yf/iot-device-mqtt/pkg/utils"
	"github.com/team4yf/iot-device-mqtt/pkg/utils/kvstore"
	"github.com/team4yf/yf-fpm-server-go/fpm"
	"github.com/team4yf/yf-fpm-server-go/pkg/log"
)

const (
	cmdFfmpeg     = "ffmpeg -i rtsp://%s:%s@%s:554/h264/1/sub/av_stream -an -f mpegts -codec:v mpeg1video -s 1280x960 -b:v 1000k -bf 0 -muxdelay 0.001 http://open.yunplus.io:18081/fpmpassword/%s"
	cmdKillFfmpeg = `ps -ef | grep ffmpeg | grep %s | grep -v "grep" | awk '{print $2}' | xargs kill -9`
	cmdOnvif      = "onvif-ptz %s --baseUrl=http://%s:80 -u=%s -p=%s -x=%f -y=%f -z=%f"
	cmdNmap       = `nmap -p %s %s | grep -E "^[1-9]" | awk '{print $1","$2 }' | sed "s/\/tcp//" | sed "s/\/udp//"`
	cmdPs         = "ps -ef | grep %s"
	cmdLsof       = "lsof -i:%d"
	cmdIP         = "ip addr | grep 'inet' | grep -v '127.0.0.1' | grep -v 'inet6' | cut -d: -f2 | awk '{print $2}' | head -1 | awk -F / '{print $1}'"
)

var (
	appID    string
	deviceID string
)
var tickerHandler *time.Ticker

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

type deviceConfigBody struct {
	BeatInterval int      `json:"beatInterval"`
	Cameras      []string `json:"cameras"`
}

func interval(cameras []string) (beatMessage *beatBody, err error) {
	beatMessage = &beatBody{
		LocalIP: utils.GetLocalIP(),
		// TODO: read from the env/arg
		GatewayID: deviceID,
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
			// fmt.Printf("count=|%s|\n", count)
			if count != "0" {
				// exists
				return
			}
		}
		finalCommand = fmt.Sprintf(cmdFfmpeg, execute.Argument...)
	case "kill-ffmpeg":
		//check exists
		out, err := utils.RunCmd("ps -ef | grep ffmpeg | grep " + execute.Argument[len(execute.Argument)-1].(string) + ` | grep -v "grep" |wc -l`)
		if err == nil {
			count := strings.Trim((string)(out), " \n")
			if count == "0" {
				// not exists, ignore
				return
			}
		}
		finalCommand = fmt.Sprintf(cmdKillFfmpeg, execute.Argument...)
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
			log.Infof("run command error: %s, error:\n %v\n", finalCommand, err)
			return
		}
		log.Infof("run command success: %s, out:\n %s\n", finalCommand, (string)(out))
		return
	}

	feedback := feedbackBody{
		MessageID: execute.MessageID,
		DeviceID:  deviceID,
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
		"topic":   fmt.Sprintf("$d2s/%s/ipc/feedback", appID),
		"payload": ([]byte)(feedbackStr),
	})
}

//ConsumerHook the hook of the consumer
//it will run after init,
//make mqtt connection
func ConsumerHook(app *fpm.Fpm) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	kvstore.Init("device.db")

	appID = app.GetConfig("uuid").(string)
	deviceID = app.GetConfig("deviceID").(string)
	app.Logger.Debugf("inited appID: %s, deviceID: %s", appID, deviceID)

	configTopic := fmt.Sprintf(`$s2d/%s/ipc/%s/config`, appID, deviceID)
	executeTopic := fmt.Sprintf(`$s2d/%s/ipc/%s/execute`, appID, deviceID)
	beatTopic := fmt.Sprintf(`$d2s/%s/ipc/beat`, appID)
	// the demo is the device id
	//{ "commad": "ps","argument":["vscode"],"messageID":"123", "feedback": 1}
	deviceConfig := deviceConfigBody{
		BeatInterval: 10,
		Cameras:      []string{"192.168.0.64"},
	}

	//load from the leveldb
	if err := kvstore.GetObject("device-config", &deviceConfig); err != nil {
		log.Errorf("load device config error: %v", err)
	}

	beatHandler := func() {
		if beatMessage, err := interval(deviceConfig.Cameras); err != nil {
			log.Error(err)
		} else {
			app.Execute("mqttclient.publish", &fpm.BizParam{
				"topic":   beatTopic,
				"payload": utils.Struct2Bytes(beatMessage),
			})
		}
	}
	go startTicker(time.Second*time.Duration(deviceConfig.BeatInterval), beatHandler)

	app.Execute("mqttclient.subscribe", &fpm.BizParam{
		"topics": []string{configTopic, executeTopic},
	})

	app.Subscribe("#mqtt/receive", func(_ string, payload interface{}) {

		body := payload.(map[string]interface{})
		log.Debugf("receive data: %v", body)
		topic := body["topic"].(string)
		switch topic {
		case configTopic:
			oldInterval := deviceConfig.BeatInterval
			if err := utils.DataToStruct(body["payload"].([]byte), &deviceConfig); err != nil {
				log.Error(err)
				return
			}
			if err := kvstore.PutObject("device-config", &deviceConfig); err != nil {
				log.Error(err)
				return
			}
			app.Logger.Debugf("flush new config: %v", deviceConfig)
			if oldInterval != deviceConfig.BeatInterval {
				if tickerHandler != nil {
					tickerHandler.Stop()
					go startTicker(time.Second*time.Duration(deviceConfig.BeatInterval), beatHandler)
				}
			}

		case executeTopic:
			//TODO: here
			execute := executeBody{}
			if err := utils.DataToStruct(body["payload"].([]byte), &execute); err != nil {
				log.Infof("convert the execute message error:", err)
				return
			}
			go runCommand(app, &execute)
		}

	})

}

func startTicker(interval time.Duration, handler func()) {
	//auto push beat info
	tickerHandler = time.NewTicker(interval)
	for {

		select {

		case <-tickerHandler.C:
			handler()
		}

	}
}
