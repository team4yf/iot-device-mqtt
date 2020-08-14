package main

import (
	"fmt"
	"testing"

	"github.com/team4yf/yf-fpm-server-go/fpm"

	_ "github.com/team4yf/yf-fpm-server-go/plugin"
)

func TestConfig(t *testing.T) {

	app := fpm.New()
	app.Init()
	setting := app.GetConfig("mqtt")
	fmt.Printf("setting: %+v", setting)
	app.Run(":9999")

}
