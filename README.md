# iot-device-mqtt
iot-device-mqtt


## command

```sh

# all device to server topics $d2s/yunplus/ipc/+

# topic:  $s2d/yunplus/ipc/demo/execute
# ffmpeg
{ "command": "ip" }

# ffmpeg
{ "command": "ffmpeg","argument":["admin","Mima123456","172.16.11.64", "abc"],"feedback":0,"messageID":"123" }

# onvif
{ "command": "onvif","argument":["move", "172.16.11.64", "admin","Mima123456",0.5, 0,0],"feedback":0,"messageID":"123" }

# nmap
{ "command": "nmap", "argument": ["80,123,254"], "feedback": 1, "messageID":"234"}
```