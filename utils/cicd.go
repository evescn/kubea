package utils

import (
	"encoding/base64"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
)

func CiCd(tmpUrl, auth string) (body []byte, err error) {
	url := tmpUrl
	method := "POST"

	payload := strings.NewReader(``)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		zap.L().Error("New HTTP 报错: ", zap.Error(err))
		return nil, errors.New(fmt.Sprintf("New HTTP 请求报错: ", err.Error()))
	}

	// 添加 Basic Authentication 头
	//auth := config.UserPassword
	//auth := settings.Conf.CiCd.UserPassword
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", authHeader)

	// 发起请求
	res, err := client.Do(req)
	if err != nil {
		zap.L().Error("HTTP 请求报错: ", zap.Error(err))
		return nil, errors.New(fmt.Sprintf("HTTP 请求报错: ", err.Error()))
	}
	defer res.Body.Close()

	body, err = io.ReadAll(res.Body)
	if err != nil {
		zap.L().Error("IO 数据解析报错: ", zap.Error(err))
		return nil, errors.New(fmt.Sprintf("IO 数据解析报错: ", err.Error()))
	}

	return body, nil
}
