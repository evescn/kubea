package service

import (
	"bytes"
	"fmt"
	"go.uber.org/zap"
	"io"
	"kubea/model"
	"kubea/settings"
	"net/http"
	"regexp"

	"encoding/json"
)

var GitLab gitlab

type gitlab struct{}

func (*gitlab) GetGroupNameID(newProjectInfo *model.GitLab) error {
	client := &http.Client{}

	// 通过API查询group列表，获取已存在的 group_name 和对应的 group_id，并根据输入的 group_name 返回对应的 group_id，如果不存在输入的 group_name，返回-1
	req, err := http.NewRequest("GET", settings.Conf.GitLab.GitLabUrl+"/api/v4/groups?search="+newProjectInfo.GroupName, nil)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}

	req.Header.Add("PRIVATE-TOKEN", settings.Conf.GitLab.GitLabToken)
	res, err := client.Do(req)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	groupinfo := model.TempInfo{}
	json.Unmarshal(body, &groupinfo)

	if len(groupinfo) == 0 {
		zap.L().Error(fmt.Sprintf("ERROR: GitLab 不存在 %s 的 GroupName\n", newProjectInfo.GroupName))
		return fmt.Errorf("ERROR: GitLab 不存在 %s 的 GroupName\n", newProjectInfo.GroupName)
	} else {
		newProjectInfo.GroupID = groupinfo[0].ID
	}

	return nil
}

func (*gitlab) CheckVisibility(visibility string) error {
	if visibility != "private" && visibility != "internal" && visibility != "public" {
		zap.L().Error("ERROR: 可见度级别关键字错误，只接收 private、internal、public 级别 ")
		return fmt.Errorf("ERROR: 可见度级别关键字错误，只接收 private、internal、public 级别 ")
	}
	return nil
}

func (*gitlab) CheckProject(projectname string) error {
	// 检查 ProjectName 是否符合规则
	if matched, _ := regexp.MatchString("^[a-z0-9-]+$", projectname); !matched {
		zap.L().Error("Error: ProjectName 只能使用小写字母、数字或-\n")
		return fmt.Errorf("Error: ProjectName 只能使用小写字母、数字或-\n")
	}

	// 检查 ProjectName 是否存在
	client := &http.Client{}
	req, err := http.NewRequest("GET", settings.Conf.GitLab.GitLabUrl+"/api/v4/projects?search="+projectname, nil)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}

	req.Header.Add("PRIVATE-TOKEN", settings.Conf.GitLab.GitLabToken)
	res, err := client.Do(req)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	projectinfo := model.TempInfo{}
	json.Unmarshal(body, &projectinfo)

	if len(projectinfo) != 0 {
		zap.L().Error("ERROR: 已存在此项目！")
		return fmt.Errorf("ERROR: 已存在此项目！")
	}

	return nil
}

func (*gitlab) CreateProject(args *model.GitLab) error {
	// 创建项目
	client := &http.Client{}
	newProjectInfo := &model.ProjectInfo{
		Name:        args.ProjectName,
		Description: args.Description,
		Path:        args.ProjectName,
		NameSpaceID: args.GroupID,
		Visibility:  args.Visibility,
		ImportUrl:   settings.Conf.GitLab.GitLabUrl + "/init/bare.git",
	}

	dataByte, err := json.Marshal(newProjectInfo)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	zap.L().Info(string(dataByte))
	req, err := http.NewRequest("POST", settings.Conf.GitLab.GitLabUrl+"/api/v4/projects", bytes.NewReader(dataByte))
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	req.Header.Set("PRIVATE-TOKEN", settings.Conf.GitLab.GitLabToken)
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	resp, err := client.Do(req)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	zap.L().Info(string(body))
	return nil
}
