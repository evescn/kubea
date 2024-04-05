package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var DaemonSet daemonSet

type daemonSet struct{}

// DaemonSetsResp 定义列表的返回类型
type DaemonSetsResp struct {
	Items []appsv1.DaemonSet `json:"items"`
	Total int                `json:"total"`
}

// toCells 方法用于将 ds 类型数组，转换成 DataCell 类型数组
func (d *daemonSet) toCells(std []appsv1.DaemonSet) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = daemonSetCell(std[i])
	}
	return cells
}

// fromCells 方法用于将 DataCell 类型数组，转换成 ds 类型数组
func (d *daemonSet) fromCells(cells []DataCell) []appsv1.DaemonSet {
	daemonsets := make([]appsv1.DaemonSet, len(cells))
	for i := range cells {
		daemonsets[i] = appsv1.DaemonSet(cells[i].(daemonSetCell))
	}
	return daemonsets
}

// GetDaemonSets 获取daemonset列表，支持过滤，排序，分页，
// client用于选择哪个集群
func (d *daemonSet) GetDaemonSets(client *kubernetes.Clientset, filterName, namespace string, limit, page int) (dssResp *DaemonSetsResp, err error) {
	// context.TODO()用于声明一个空的context上下文，用于List方法内设置这个请求的超时（源码），这里的常用用法
	//metav1.ListOptions{}用于过滤List数据，如使用label，field等
	dsList, err := client.AppsV1().DaemonSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("获取 DaemonSet 列表失败, %v\n", err))
		return nil, errors.New(fmt.Sprintf("获取 DaemonSet 列表失败, %v\n", err))
	}
	//实例化dataSelector对象，把 d 结构体中获取到的 DaemonSet 列表转化为 dataSelector 结构体，方便使用 dataSelector 结构体中 过滤，排序，分页功能
	selectableData := &dataSelector{
		GenericDataList: d.toCells(dsList.Items),
		dataSelectQuery: &DataSelectQuery{
			FilterQuery: &FilterQuery{Name: filterName},
			PaginateQuery: &PaginateQuery{
				Limit: limit,
				Page:  page,
			},
		},
	}
	//先过滤，filtered中的数据才是总数据，data中的数据是排序分页后的数据，可能每次只有10行
	filtered := selectableData.Filter()
	total := len(filtered.GenericDataList)

	//在排序和分页
	data := filtered.Sort().Paginate()
	//将[]DataCell类型的 DaemonSet 列表转为 v1.DaemonSet 列表
	dss := d.fromCells(data.GenericDataList)
	return &DaemonSetsResp{
		Items: dss,
		Total: total,
	}, nil
}

// GetDaemonSetDetail 获取 DaemonSet 详情
func (d *daemonSet) GetDaemonSetDetail(client *kubernetes.Clientset, dsName, namespace string) (ds *appsv1.DaemonSet, err error) {
	dsDetail, err := client.AppsV1().DaemonSets(namespace).Get(context.TODO(), dsName, metav1.GetOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("获取 DaemonSet 详情失败, %v\n", err))
		return nil, errors.New(fmt.Sprintf("获取 DaemonSet 详情失败, %v\n", err))
	}
	return dsDetail, nil
}

// DeleteDaemonSet 删除 DaemonSet
func (d *daemonSet) DeleteDaemonSet(client *kubernetes.Clientset, dsName, namespace string) (err error) {
	err = client.AppsV1().DaemonSets(namespace).Delete(context.TODO(), dsName, metav1.DeleteOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("删除 DaemonSet 失败, %v\n", err))
		return errors.New(fmt.Sprintf("删除 DaemonSet 失败, %v\n", err))
	}
	return nil
}

// UpdateDaemonSet 更新 DaemonSet
// content就是DaemonSet的整个json体
func (d *daemonSet) UpdateDaemonSet(client *kubernetes.Clientset, content, namespace string) (err error) {
	//content转成pod结构体
	var dss = &appsv1.DaemonSet{}
	//反序列化成pod对象
	err = json.Unmarshal([]byte(content), &dss)
	if err != nil {
		zap.L().Error(fmt.Sprintf("反序列化失败, %v\n", err))
		return errors.New(fmt.Sprintf("反序列化失败, %v\n", err))
	}
	//更新pod
	_, err = client.AppsV1().DaemonSets(namespace).Update(context.TODO(), dss, metav1.UpdateOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("更新 DaemonSet 失败, %v\n", err))
		return errors.New(fmt.Sprintf("更新 DaemonSet 失败, %v\n", err))
	}
	return nil
}
