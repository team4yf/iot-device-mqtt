package main

import (
	"github.com/team4yf/iot-device-mqtt/hooks"

	"github.com/team4yf/yf-fpm-server-go/fpm"

	_ "github.com/team4yf/fpm-go-plugin-mqtt-client/plugin"
)

func main() {

	app := fpm.New()

	app.AddHook("AFTER_INIT", hooks.ConsumerHook, 10)
	app.Init()

	app.Run()

}
