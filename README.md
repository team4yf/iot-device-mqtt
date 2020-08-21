# iot-device-mqtt
iot-device-mqtt


## command

```sh

# all device to server topics $d2s/ceaa191a/ipc/+
# $d2s/ceaa191a/ipc/beat

# topic:  $s2d/ceaa191a/ipc/foodevice/execute
# ffmpeg
{ "command": "ip" }

# ffmpeg
{ "command": "ffmpeg","argument":["admin","Mima123456","172.16.11.64", "abc"],"feedback":0,"messageID":"123" }

# kill-ffmpeg
{ "command": "kill-ffmpeg","argument":["172.16.11.64"],"feedback":0,"messageID":"123" }

# onvif
{ "command": "onvif","argument":["move", "172.16.11.64", "admin","Mima123456",0.5, 0,0],"feedback":0,"messageID":"123" }

# nmap
{ "command": "nmap", "argument": ["80,123,254", "192.168.0.108"], "feedback": 1, "messageID":"234"}

# config
{ "beatInterval": 20, "cameras": ["172.16.11.64"]}
```

## ref

leveldb: https://github.com/syndtr/goleveldb


## config

config.json -> uuid: [app_id]

env -> FPM_DEVICE: [device_id]

## mosquitto command shell

```sh
# install mosquitto
sudo apt  install mosquitto-clients

# sub
mosquitto_sub -h open.yunplus.io -t '$d2s/ceaa191a/ipc/beat' -u "fpmuser" -P

# pub

## config
mosquitto_pub -h open.yunplus.io -t '$s2d/ceaa191a/ipc/iot-device-ipc1/config' -m '{ "beatInterval": 20, "cameras": ["172.16.11.64"]}' -u "fpmuser" -P

## onvif
mosquitto_pub -h open.yunplus.io -t '$s2d/ceaa191a/ipc/iot-device-ipc1/execute' -m '{ "command": "onvif","argument":["move", "172.16.11.64", "admin","Mima123456",0.5, 0,0],"feedback":0,"messageID":"123" }' -u "fpmuser" -P


## ffmpeg
mosquitto_pub -h open.yunplus.io -t '$s2d/ceaa191a/ipc/iot-device-ipc1/execute' -m '{ "command": "ffmpeg","argument":["admin","Mima123456","172.16.11.64", "abc"],"feedback":0,"messageID":"123" }
' -u "fpmuser" -P

http://open.yunplus.io:18081/static/demo.html

```