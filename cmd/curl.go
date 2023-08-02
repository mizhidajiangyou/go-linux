package cmd

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"moul.io/http2curl"
	"net/http"
	"os"
	"strings"
	"time"
)

const CurlKey = "c_key"

type CurlCommand struct {
	// -F Specify multipart MIME data
	File string
	// -O Write to file instead of stdout
	Output string
	// -H Pass custom header(s) to server
	Header map[string][]string
	// --data-raw <data>  HTTP POST data, '@' allowed
	Data interface{}
	// -X Specify request command to use
	Method string
	// --retry <num>   Retry request if transient problems occur
	RetryTimes int
	// --max-time <fractional seconds>  Maximum time allowed for the transfer
	MaxTimes time.Duration
	// --retry-delay <seconds>  Wait time between retries
	RetryDelay time.Duration
}

type CurlData struct {
	UsedTime time.Duration
	CurlCmd  string
}

// Curl http请求
func Curl(ctx context.Context, url string, com CurlCommand) (statusCode int, reqData map[string]interface{}, err error) {
	var resp []byte
	jsonBody, err := json.Marshal(com.Data)
	if err != nil {
		return
	}
	var request *http.Request
	if com.File != "" {
		err = LsFile(com.File)
		if err != nil {
			return
		}
		payload := &bytes.Buffer{}
		writer := multipart.NewWriter(payload)
		reads, _ := ioutil.ReadFile(com.File)
		ss := strings.Split(com.File, "/")
		fileName := fmt.Sprintf("%s", ss[len(ss)-1])
		writes, _ := writer.CreateFormFile("file", fileName)
		_, err = writes.Write(reads)
		if err != nil {
			return
		}
		err = writer.Close()
		if err != nil {
			return
		}
		request, err = http.NewRequestWithContext(ctx, com.Method, url, payload)
		request.Header = com.Header
		realContentType := writer.FormDataContentType()
		request.Header.Set("Content-Type", realContentType)
	} else {
		request, err = http.NewRequestWithContext(ctx, com.Method, url, strings.NewReader(string(jsonBody)))
	}

	if err != nil {
		return
	}
	request.Header = com.Header
	client := http.DefaultClient
	client.Timeout = com.MaxTimes
	var used time.Duration
	for i := 0; i < com.RetryTimes; i++ {
		start := time.Now()
		statusCode, resp, err = doCurl(client, request)
		end := time.Now()
		used = end.Sub(start)
		if err != nil {
			if i == com.RetryTimes-1 {
				return
			}
			time.Sleep(com.RetryDelay)
			continue
		}
		break
	}
	if com.Output != "" {
		//fmt.Print(file)
		var fileObj *os.File
		fileObj, err = os.OpenFile(com.Output, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return
		}
		writeObj := bufio.NewWriterSize(fileObj, 4096)
		buf := resp
		if _, err = writeObj.Write(buf); err == nil {
			if err = writeObj.Flush(); err != nil {
				return
			}
		}
		defer func(fileObj *os.File) {
			err = fileObj.Close()
			if err != nil {
				return
			}
		}(fileObj)
		resp = nil
	}

	err = json.Unmarshal(resp, &reqData)
	command, _ := http2curl.GetCurlCommand(request)
	ctx = context.WithValue(ctx, CurlKey, CurlData{
		CurlCmd:  command.String(),
		UsedTime: used,
	})
	return

}

// doCurl 使用提供的客户端发送 HTTP 请求，并返回响应状态码、响应体和错误（如果有）。
func doCurl(client *http.Client, request *http.Request) (statusCode int, resp []byte, err error) {
	// 发送 HTTP 请求
	response, err := client.Do(request)
	if err != nil {
		return 0, nil, fmt.Errorf("发送 HTTP 请求失败：%w", err)
	}
	defer func() {
		cerr := response.Body.Close()
		if cerr != nil && err == nil {
			err = fmt.Errorf("关闭 HTTP 响应体失败：%w", cerr)
		}
	}()

	// 检查 HTTP 响应状态码
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return response.StatusCode, nil, fmt.Errorf("HTTP 错误：%d", response.StatusCode)
	}

	// 读取 HTTP 响应体
	resp, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return response.StatusCode, nil, fmt.Errorf("读取 HTTP 响应体失败：%w", err)
	}

	return response.StatusCode, resp, nil
}
