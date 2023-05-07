package service

import (
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/release"
	"io"
	"kubea-demo/config"
	"kubea-demo/dao"
	"kubea-demo/model"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
)

var HelmStore helmStore

type helmStore struct{}

// 定义列表返回的内容
type releaseElement struct {
	Name         string `json:"name"`
	Namespace    string `json:"namespace"`
	Revision     string `json:"revision"`
	Updated      string `json:"updated"`
	Status       string `json:"status"`
	Chart        string `json:"chart"`
	ChartVersion string `json:"chart_version"`
	AppVersion   string `json:"app_version"`
	Notes        string `json:"notes,omitempty"`
}

type releaseElements struct {
	Items []*releaseElement `json:"items"`
	Total int               `json:"total"`
}

// ListReleases release列表
// 这里没有使用page和limit,这里的分页是前端实现的,翻页不发起请求
// k8s资源使用了page和limit获取列表，每次翻页都发起请求
func (*helmStore) ListReleases(actionConfig *action.Configuration, filterName string) (*releaseElements, error) {
	// new一个列表的client
	client := action.NewList(actionConfig)
	client.Filter = filterName

	// 显示所有数据
	client.All = true
	//client.Limit = limit
	//client.Offset = offset
	client.TimeFormat = "2006-01-02 15:04:05"

	// 只获取已部署 helm
	client.Deployed = true
	results, err := client.Run()
	if err != nil {
		logger.Error(fmt.Sprintf("获取Release列表失败, %v\n", err))
		return nil, errors.New(fmt.Sprintf("获取Release列表失败, %v\n", err))
	}

	total := len(results)
	elements := make([]*releaseElement, 0)
	for _, r := range results {
		elements = append(elements, constructReleaseElement(r, false))
	}

	return &releaseElements{
		Items: elements,
		Total: total,
	}, nil
}

// DetailRelease 获取release详情
func (*helmStore) DetailRelease(actionConfig *action.Configuration, release string) (*release.Release, error) {
	client := action.NewGet(actionConfig)
	data, err := client.Run(release)
	if err != nil {
		logger.Error(fmt.Sprintf("获取Release详情失败, %v\n", err))
		return nil, errors.New(fmt.Sprintf("获取Release详情失败, %v\n", err))
	}
	return data, nil
}

// InstallRelease 安装Release
// release release的名字
// chart chart文件所在的路径
func (*helmStore) InstallRelease(actionConfig *action.Configuration, release, chart, namespace string) error {
	client := action.NewInstall(actionConfig)
	client.ReleaseName = release
	//这里的namespace没啥用，主要安装在哪个namespace还是要看actionConfig初始化的namespace
	client.Namespace = namespace
	splitChart := strings.Split(chart, ".")
	if splitChart[len(splitChart)-1] == "tgz" && !strings.Contains(chart, ":") {
		chart = config.UploadPath + chart
	}

	//加载chart文件，并给予文件内容生成k8s资源
	chartRequested, err := loader.Load(chart)
	if err != nil {
		logger.Error(fmt.Sprintf("加载Chart文件失败, %v\n", err))
		return errors.New(fmt.Sprintf("加载Chart文件失败, %v\n", err))
	}
	vals := make(map[string]interface{}, 0)
	_, err = client.Run(chartRequested, vals)
	if err != nil {
		logger.Error(fmt.Sprintf("安装Release失败, %v\n", err))
		return errors.New(fmt.Sprintf("安装Release失败, %v\n", err))
	}

	return nil
}

// UninstallRelease 卸载release
func (*helmStore) UninstallRelease(actionConfig *action.Configuration, release string) error {
	client := action.NewUninstall(actionConfig)
	_, err := client.Run(release)
	if err != nil {
		logger.Error(fmt.Sprintf("卸载Release失败, %v\n", err))
		return errors.New(fmt.Sprintf("卸载Release失败, %v\n", err))
	}
	return nil
}

// UploadChartFile chart文件上传
// *multipart.FileHeader 获取上传文件header中的相关信息
// multipart.File封装了读文件的操作
func (*helmStore) UploadChartFile(file multipart.File, header *multipart.FileHeader) error {
	filename := header.Filename
	t := strings.Split(filename, ".")
	if t[len(t)-1] != "tgz" {
		logger.Error("chart文件必须以.tgz结尾")
		return errors.New("chart文件必须以.tgz结尾")
	}
	filePath := config.UploadPath + filename
	_, err := os.Stat(filePath)
	if os.IsExist(err) {
		logger.Error("chart文件已存在")
		return errors.New("chart文件已存在")
	}
	out, err := os.Create(filePath)
	if err != nil {
		logger.Error(fmt.Sprintf("创建chart文件失败 %v\n", err))
		return errors.New(fmt.Sprintf("创建chart文件失败 %v\n", err))
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		logger.Error(fmt.Sprintf("创建chart文件失败2 %v\n", err))
		return errors.New(fmt.Sprintf("创建chart文件失败2 %v\n", err))
	}
	return nil
}

// ListCharts 获取 chart 列表
func (*helmStore) ListCharts(name string, page, limit int) (*dao.Charts, error) {
	return dao.Chart.GetList(name, page, limit)
}

// DeleteChartFile chart 文件删除
func (*helmStore) DeleteChartFile(chart string) error {
	//路径拼接这种写法只支持mac或者linux，如果是windows，则要改成 filePath := config.UploadPath + "\\" + chart
	filePath := config.UploadPath + chart
	_, err := os.Stat(filePath)
	if err != nil || os.IsNotExist(err) {
		logger.Error(fmt.Sprintf("chart文件不存在 %v\n", err))
		return errors.New(fmt.Sprintf("chart文件不存在 %v\n", err))
	}

	// 直接删除
	err = os.Remove(filePath)
	if err != nil {
		logger.Error(fmt.Sprintf("chart文件删除失败 %v\n", err))
		return errors.New(fmt.Sprintf("chart文件删除失败 %v\n", err))
	}
	return nil
}

// AddChart chart 新增
func (*helmStore) AddChart(chart *model.Chart) error {
	_, has, err := dao.Chart.Has(chart.Name)
	if err != nil {
		return err
	}
	if has {
		return errors.New("该数据已存在，请重新添加")
	}

	if err := dao.Chart.Add(chart); err != nil {
		return err
	}
	return nil
}

// UpdateChart chart 更新
func (h *helmStore) UpdateChart(chart *model.Chart) error {
	oldChart, _, err := dao.Chart.Has(chart.Name)
	if err != nil {
		return err
	}
	fmt.Println(chart.FileName, oldChart.FileName)
	//如果更新了新的上传文件，则老的文件要删除
	if chart.FileName != "" && chart.FileName != oldChart.FileName {
		err = h.DeleteChartFile(oldChart.FileName)
		if err != nil {
			return err
		}
	}
	return dao.Chart.Update(chart)
}

// DeleteChart chart 删除
func (h *helmStore) DeleteChart(chart *model.Chart) error {
	//删除文件
	err := h.DeleteChartFile(chart.FileName)
	if err != nil {
		return err
	}
	//删除数据
	return dao.Chart.Delete(chart.ID)
}

// constructReleaseElement release内容过滤
func constructReleaseElement(r *release.Release, showStatus bool) *releaseElement {
	element := &releaseElement{
		Name:         r.Name,
		Namespace:    r.Namespace,
		Revision:     strconv.Itoa(r.Version),
		Status:       r.Info.Status.String(),
		Chart:        r.Chart.Metadata.Name,
		ChartVersion: r.Chart.Metadata.Version,
		AppVersion:   r.Chart.Metadata.AppVersion,
	}

	if showStatus {
		element.Notes = r.Info.Notes
	}

	// 判断输出
	t := "-"
	if tspb := r.Info.LastDeployed; !tspb.IsZero() {
		t = tspb.String()
	}
	element.Updated = t

	return element
}
