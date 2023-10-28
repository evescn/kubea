package service

import (
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	"kubea-demo/config"
	"kubea-demo/dao"
	"kubea-demo/model"
	"kubea-demo/utils"
	"time"
)

var Deploy deploy

type deploy struct{}

// List 列表
func (*deploy) List(en, appName, repoName string, page, limit int) (*dao.Deploys, error) {
	return dao.Deploy.List(en, appName, repoName, page, limit)
}

// Add 新增
func (*deploy) Add(d *model.Deploy) error {
	data, has, err := dao.App.Get(d.AppId)
	if err != nil {
		return err
	}

	if !has {
		return errors.New("查询无此应用")
	}

	d.BuildUrl = fmt.Sprintf(config.JenkinsUrl, data.RepoName, fmt.Sprintf("%s-%s", data.RepoName, data.AppName))
	return dao.Deploy.Add(d)
}

// Update 更新
func (*deploy) Update(d *model.Deploy) error {
	return dao.Deploy.Update(d)
}

// Delete 删除
func (*deploy) Delete(id uint) error {
	return dao.Deploy.Delete(id)
}

// CiCd 开始部署
func (*deploy) CiCd(d *model.Deploy) error {
	data, has, err := dao.Deploy.Get(d.ID)
	if err != nil {
		return err
	}

	if !has {
		return errors.New("查询无此部署任务")
	}

	// 更新数据
	d.Status = 1
	//d.
	err = dao.Deploy.Update(d)
	if err != nil {
		return err
	}

	// 请求 jenkins 服务
	var tag bool
	if data.Tag == 1 {
		tag = true
	} else if data.Tag == 2 {
		tag = false
	}

	url := fmt.Sprintf("%sbuildWithParameters?ENV=%s&BRANCH=%s&IS_CREATE_TAG=%v", data.BuildUrl, data.En, data.Branch, tag)
	_, err = utils.CiCd(url)

	if err != nil {
		return err
	}
	// 写入日志
	Deploy.LogInfo(data.ID, "服务开始部署！！！")

	return nil
}

// JenkinsCiCd 开始部署
func (*deploy) JenkinsCiCd(appId uint, en, startTime, builder string) error {
	data, has, err := dao.Deploy.Has(en, appId)
	if err != nil {
		return err
	}

	if !has {
		return errors.New("查询无此部署任务")
	}

	// 更新数据
	data.Status = 1
	data.StartTime = startTime
	data.Builder = builder

	err = dao.Deploy.Update(data)
	if err != nil {
		return err
	}

	// 写入日志
	Deploy.LogInfo(data.ID, "服务开始部署！！！")

	return nil
}

// UpdateCiCd 更新CiCd流程
func (*deploy) UpdateCiCd(en, appName, repoName, branch string, codeCheck, buildStatus, deployStatus int) error {
	var deplayID uint

	// 使用 en 和 appName 获取数据
	deploysData, err := dao.Deploy.List(en, appName, repoName, 1, 10)
	if err != nil {
		return err
	}

	// 返回第一条记录即可
	for _, item := range deploysData.Items {
		deplayID = item.ID
		break
	}

	// 获取数据
	deployData, _, err := dao.Deploy.Get(deplayID)
	if err != nil {
		return err
	}

	if deplayID == 0 {
		return errors.New("查询无此应用")
	}

	// 判断 codeCheck 是否 != 0
	if codeCheck != 0 {
		deployData.CodeCheck = codeCheck
		if codeCheck == 2 {
			deployData.Status = 4
			// 部署耗时
			deployData.Duration = Deploy.getDuration(deployData.StartTime)
			// 更新日志
			Deploy.LogInfo(deployData.ID, "代码检查失败！代码 checkout 未通过，当前分支存在未合并代码")
		}
	}

	// 判断 buildStatus 是否 == 1
	if buildStatus == 1 {
		deployData.BuildStatus = buildStatus
		deployData.Status = 2
		// 更新日志
		Deploy.LogInfo(deployData.ID, "服务开始编译！！！")
	} else if buildStatus == 2 {
		deployData.BuildStatus = buildStatus
		deployData.Status = 4
		// 部署耗时
		deployData.Duration = Deploy.getDuration(deployData.StartTime)
		// 更新日志
		Deploy.LogInfo(deployData.ID, "服务编译失败！！！")
	}

	// 判断 deployStatus 是否 == 2
	if deployStatus == 1 && deployData.Status == 2 {
		deployData.DeployStatus = deployStatus
		deployData.Status = 3
		// 程序启动日志查看
		// 部署耗时
		deployData.Duration = Deploy.getDuration(deployData.StartTime)
		// 更新日志
		Deploy.LogInfo(deployData.ID, "服务开始部署！！！")
	} else if deployStatus == 2 {
		deployData.DeployStatus = deployStatus
		deployData.Status = 4
		// 部署耗时
		deployData.Duration = Deploy.getDuration(deployData.StartTime)
		// 更新日志
		Deploy.LogInfo(deployData.ID, "服务部署失败！！！")
	}

	// 更新 Tag 信息
	if len(branch) > 0 {
		deployData.Branch = branch
	}

	// 更新 Tag 信息
	if len(branch) > 0 {
		deployData.Branch = branch
	}

	// 更新服务
	return dao.Deploy.Update(deployData)
}

// getDuration 获取时间，计算耗时
func (*deploy) getDuration(startTimeStr string) string {
	// 解析开始时间
	loc, _ := time.LoadLocation("Asia/Shanghai")

	startTime, err := time.ParseInLocation("2006-01-02 15:04:05", startTimeStr, loc)
	if err != nil {
		logger.Error("解析开始时间时出错:", err)
		return ""
	}

	// 获取当前时间作为结束时间
	endTime := time.Now()

	// 计算时间差
	duration := endTime.Sub(startTime)
	// 使用 String 方法将时间差格式化为分秒形式
	durationString := duration.Truncate(time.Second).String()

	// 打印时间差
	return durationString
}

// GetLog 查询日志
func (*deploy) GetLog(deployId uint) (*model.DeployLog, bool, error) {
	return dao.Deploy.GetLog(deployId)
}

// LogInfo 写日志
func (*deploy) LogInfo(deployId uint, content string) {
	// 查日志
	data, has, err := dao.Deploy.GetLog(deployId)
	if err != nil {
		logger.Error(err)
	}
	// 若存在就修改
	if has {
		logContent := fmt.Sprintf("%s%s\n", data.Log, content)
		data.Log = logContent
		err = dao.Deploy.UpdateLog(data)
		if err != nil {
			logger.Error(err)
		}
	} else {
		// 若不存在就新增
		logContent := fmt.Sprintf("%s\n", content)

		err = dao.Deploy.AddLog(&model.DeployLog{
			DeployId: deployId,
			Log:      logContent,
		})
		if err != nil {
			logger.Error(err)
		}
	}

}
