package service

import (
	"fmt"
	"github.com/wonderivan/logger"
	"kubea-demo/dao"
	"kubea-demo/model"
)

var Deploy deploy

type deploy struct{}

// List 列表
func (*deploy) List(en, appName string, page, limit int) (*dao.Deploys, error) {
	return dao.Deploy.List(en, appName, page, limit)
}

// Add 新增
func (*deploy) Add(d *model.Deploy) error {
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
