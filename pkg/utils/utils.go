package utils

import (
	"encoding/json"
	"sync"

	"github.com/google/uuid"
	"github.com/teris-io/shortid"
	tnet "github.com/toolkits/net"
)

//RespJSON the common json
type RespJSON struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var (
	once     sync.Once
	clientIP = "127.0.0.1"
)

//CheckErr panic if err is not nil
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

// GetLocalIP 获取本地内网IP
func GetLocalIP() string {
	once.Do(func() {
		ips, _ := tnet.IntranetIP()
		if len(ips) > 0 {
			clientIP = ips[0]
		} else {
			clientIP = "127.0.0.1"
		}
	})
	return clientIP
}

//CutBigIntSlice cut long int slice to [][] slice
func CutBigIntSlice(origin []int, size int) (desc [][]int) {
	total := len(origin)
	rows := total / size

	if total%size > 0 {
		rows++
	}
	for i := 0; i < rows; i++ {
		desc = append(desc, []int{})
	}
	for i, d := range origin {
		row := (i + 1) / size
		if (i+1)%size == 0 {
			row--
		}
		desc[row] = append(desc[row], d)
	}
	return
}

//JSON2String convert the json object to string
func JSON2String(j interface{}) (str string) {
	bytes, err := json.Marshal(j)
	if err != nil {
		return "{}"
	}
	str = (string)(bytes)
	return
}

//StringToStruct convert the string to struct
func StringToStruct(data string, desc interface{}) (err error) {
	if err = json.Unmarshal(([]byte)(data), desc); err != nil {
		return
	}
	return
}

// GenShortID 生成一个id
func GenShortID() string {
	sid, _ := shortid.Generate()
	return sid
}

// GenUUID 生成随机字符串，eg: 76d27e8c-a80e-48c8-ad20-e5562e0f67e4
func GenUUID() string {
	u, _ := uuid.NewRandom()
	return u.String()
}

// GetReqID 获取请求中的request_id
// func GetReqID(c *gin.Context) string {
// 	v, ok := c.Get("X-Request-ID")
// 	if !ok {
// 		return ""
// 	}
// 	if requestID, ok := v.(string); ok {
// 		return requestID
// 	}
// 	return ""
// }

// SendResponse 返回json
// func SendResponse(c *gin.Context, err error, data interface{}) {
// 	code, message := errno.DecodeErr(err)

// 	// always return http.StatusOK
// 	c.JSON(http.StatusOK, RespJSON{
// 		Code:    code,
// 		Message: message,
// 		Data:    data,
// 	})
// }
