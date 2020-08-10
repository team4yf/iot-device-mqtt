package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

//ResponseWrapper the wrapper of the http response
type ResponseWrapper struct {
	StatusCode int
	Success    bool
	Err        error
	Body       []byte
	Header     http.Header
}

//ConvertBody Convert the body to another struct
func (rsp *ResponseWrapper) ConvertBody(data interface{}) (err error) {
	if json.Unmarshal(rsp.Body, data); err != nil {
		return
	}
	return
}

//GetStringBody get the string of the body
func (rsp *ResponseWrapper) GetStringBody() string {

	return (string)(rsp.Body)
}

//Get send a get request with timeout
func Get(url string, timeout int) ResponseWrapper {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return createRequestError(err)
	}

	return request(req, timeout)
}

//GetWithHeader send a get request with header and timeout
func GetWithHeader(url string, header map[string]string, timeout int) ResponseWrapper {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return createRequestError(err)
	}
	for k, v := range header {
		req.Header.Set(k, v)
	}

	return request(req, timeout)
}

//PostParams post a form request with timeout
func PostParams(url string, params string, timeout int) ResponseWrapper {
	buf := bytes.NewBufferString(params)
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return createRequestError(err)
	}
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")

	return request(req, timeout)
}

//PostJSON post a json data request with timeout
func PostJSON(url string, body []byte, timeout int) ResponseWrapper {
	buf := bytes.NewBuffer(body)
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return createRequestError(err)
	}
	req.Header.Set("Content-type", "application/json")

	return request(req, timeout)
}

//PostJSONWithHeader post json & header with timeout
func PostJSONWithHeader(url string, header map[string]string, body []byte, timeout int) ResponseWrapper {
	buf := bytes.NewBuffer(body)
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return createRequestError(err)
	}
	req.Header.Set("Content-type", "application/json")
	for k, v := range header {
		req.Header.Set(k, v)
	}

	return request(req, timeout)
}

func request(req *http.Request, timeout int) ResponseWrapper {
	wrapper := ResponseWrapper{StatusCode: 0, Success: false}
	client := &http.Client{}
	if timeout > 0 {
		client.Timeout = time.Duration(timeout) * time.Second
	}
	setRequestHeader(req)
	resp, err := client.Do(req)
	if err != nil {
		wrapper.Err = errors.Wrap(err, "执行HTTP请求错误")
		return wrapper
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		wrapper.Err = errors.Wrap(err, "读取HTTP请求返回值失败")
		return wrapper
	}
	wrapper.StatusCode = resp.StatusCode
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		wrapper.Success = true
	}
	wrapper.Body = body
	wrapper.Header = resp.Header

	return wrapper
}

func setRequestHeader(req *http.Request) {
	req.Header.Set("User-Agent", "fpm-iot-go-middleware")
}

func createRequestError(err error) ResponseWrapper {
	errorMessage := errors.Wrap(err, "创建HTTP请求错误")
	return ResponseWrapper{StatusCode: 0, Success: false, Err: errorMessage}
}
