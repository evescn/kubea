package service

import (
	"fmt"
	"kubea/model"
	"kubea/settings"
	"kubea/utils"
	"strings"
)

var Jenkins jenkins

type jenkins struct{}

func (*jenkins) CreatePipeline(data *model.Jenkins) error {
	var url string

	fmt.Println(data)

	// 判断 cocos 安卓 ios 项目和前后端不在一个 jenkins 上
	if strings.Contains(data.Name, "cocos") || strings.Contains(data.Name, "android") || strings.Contains(data.Name, "ios") {
		// 拼接字符串
		fmt.Println("111111111111111111111111111111111111111111111111111111111111111111111111111")
		url = fmt.Sprintf("%s/createItem?name=%s&mode=copy&from=%s", settings.Conf.CiCd.CocosJenkinsUrl, fmt.Sprintf("test-%s-%s", data.GroupName, data.Name), data.CopyJobName)

		_, err := utils.CiCd(url, settings.Conf.CiCd.CocosUserPassword)
		if err != nil {
			return err
		}

		//zap.L().Info(string(body))

	} else {
		// 拼接字符串
		url = fmt.Sprintf("%s/createItem?name=%s&mode=copy&from=%s", settings.Conf.CiCd.JenkinsUrl, fmt.Sprintf("%s-%s", data.GroupName, data.Name), data.CopyJobName)

		_, err := utils.CiCd(url, settings.Conf.CiCd.UserPassword)
		if err != nil {
			return err
		}

		//zap.L().Info(string(body))
	}

	return nil
}
