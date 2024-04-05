package service

import (
	"errors"
	"fmt"
	"kubea/dao"
	"kubea/model"
	"kubea/settings"
)

var App app

type app struct{}

// List 列表
func (*app) List(appName, repoName string, page, limit int) (*dao.Apps, error) {
	return dao.App.List(appName, repoName, page, limit)
}

// Get 查询单个应用信息
func (*app) Get(repoName, appName string) (*model.App, error) {
	data, has, err := dao.App.Has(repoName, appName)
	if err != nil {
		return nil, err
	}

	if !has {
		return nil, errors.New("查询无此应用")
	}
	return data, nil
}

// GetAll 查询所有应用
func (*app) GetAll() ([]*model.App, error) {
	return dao.App.GetAll()
}

// GetApp 根据仓库信息查询App信息
func (*app) GetApp(repoName string) ([]*model.App, error) {
	data, has, err := dao.App.GetApp(repoName)
	if err != nil {
		return nil, err
	}

	if !has {
		return nil, errors.New("repo 仓库下没有数据")
	}
	return data, nil
}

// Update 更新
func (*app) Update(app *model.App) error {
	return dao.App.Update(app)
}

// Add 新增
func (*app) Add(apps *model.Apps) error {
	//查看应用名是否存在
	_, has, err := dao.App.Has(apps.RepoName, apps.AppName)
	if err != nil {
		return err
	}
	if has {
		return errors.New("该数据已存在，请重新添加")
	}

	fmt.Println(apps)

	//不存在则创建，但是需要先创建 GitLab 和 Jenkins 流水线任务
	// 创建 GitLab
	if apps.GitLabJenkins.HasGitLab {
		// 存在则创建 GitLab
		newProjectInfo := &model.GitLab{
			GroupName:   apps.RepoName,
			ProjectName: apps.AppName,
			Visibility:  apps.GitLabJenkins.Visibility,
			Description: apps.Description,
		}

		// 检查 RepoName（组） 是否存在
		if err := GitLab.GetGroupNameID(newProjectInfo); err != nil {
			return err
		}

		// 检查 Visibility
		if err := GitLab.CheckVisibility(newProjectInfo.Visibility); err != nil {
			return err
		}

		// 检查 AppName（项目） 是否存在
		if err := GitLab.CheckProject(newProjectInfo.ProjectName); err != nil {
			return err
		}

		// 创建项目
		if err := GitLab.CreateProject(newProjectInfo); err != nil {
			return err
		}
	}

	// 创建 Jenkins 流水线任务
	if apps.GitLabJenkins.HasJenkins {
		// 存在则创建 Jenkins 流水线
		newPipelineInfo := &model.Jenkins{
			GroupName:   apps.RepoName,
			Name:        apps.AppName,
			CopyJobName: settings.Conf.CiCd.CopyJobName,
		}

		// GitLab 确保仓库唯一，此处不用检查
		// 创建
		if err := Jenkins.CreatePipeline(newPipelineInfo); err != nil {
			return err
		}
	}

	// 创建数据库数据，页面展示
	app := &model.App{
		AppName:     apps.AppName,
		RepoName:    apps.RepoName,
		Lang:        apps.Lang,
		Type:        apps.Type,
		Owner:       apps.Owner,
		Description: apps.Description,
	}

	if err := dao.App.Add(app); err != nil {
		return err
	}

	return nil
}

// Delete 删除
func (*app) Delete(id uint) error {
	return dao.App.Delete(id)
}
