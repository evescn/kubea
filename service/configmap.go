package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var ConfigMap configmap

type configmap struct{}

// ConfigMapsResp 定义列表的返回类型
type ConfigMapsResp struct {
	Items []corev1.ConfigMap `json:"items"`
	Total int                `json:"total"`
}

// toCells 方法用于将 configmap 类型数组，转换成DataCell类型数组
func (c *configmap) toCells(std []corev1.ConfigMap) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = configmapCell(std[i])
	}
	return cells
}

// fromCells 方法用于将DataCell类型数组，转换成 configmap 类型数组
func (c *configmap) fromCells(cells []DataCell) []corev1.ConfigMap {
	configmaps := make([]corev1.ConfigMap, len(cells))
	for i := range cells {
		configmaps[i] = corev1.ConfigMap(cells[i].(configmapCell))
	}
	return configmaps
}

// GetConfigMaps 获取 ConfigMap 列表
func (c *configmap) GetConfigMaps(client *kubernetes.Clientset, filterName, namespace string, limit, page int) (configMapsResp *ConfigMapsResp, err error) {
	// context.TODO()用于声明一个空的context上下文，用于List方法内设置这个请求的超时（源码），这里的常用用法
	cmList, err := client.CoreV1().ConfigMaps(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		zap.L().
			Error(fmt.Sprintf("获取 ConfigMap 列表失败, %v\n", err))
		return nil, errors.New(fmt.Sprintf("获取 ConfigMap 列表失败, %v\n", err))
	}
	//实例化dataSelector对象，把 d 结构体中获取到的 StatefulSet 列表转化为 dataSelector 结构体，方便使用 dataSelector 结构体中 过滤，排序，分页功能
	selectableData := &dataSelector{
		GenericDataList: c.toCells(cmList.Items),
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
	//将[]DataCell类型的 StatefulSet 列表转为 v1.StatefulSet 列表
	configMaps := c.fromCells(data.GenericDataList)
	return &ConfigMapsResp{
		Items: configMaps,
		Total: total,
	}, nil
}

// GetConfigMapDetail 获取 ConfigMap 详情
func (c *configmap) GetConfigMapDetail(client *kubernetes.Clientset, cmName, namespace string) (cm *corev1.ConfigMap, err error) {
	cmDetail, err := client.CoreV1().ConfigMaps(namespace).Get(context.TODO(), cmName, metav1.GetOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("获取 ConfigMap 详情失败, %v\n", err))
		return nil, errors.New(fmt.Sprintf("获取 ConfigMap 详情失败, %v\n", err))
	}
	return cmDetail, nil
}

// DeleteConfigMap 删除 ConfigMap
func (c *configmap) DeleteConfigMap(client *kubernetes.Clientset, cmName, namespace string) (err error) {
	err = client.CoreV1().ConfigMaps(namespace).Delete(context.TODO(), cmName, metav1.DeleteOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("删除 ConfigMap 失败, %v\n", err))
		return errors.New(fmt.Sprintf("删除 ConfigMap 失败, %v\n", err))
	}
	return nil
}

// UpdateConfigMap 更新 ConfigMap
// content就是StatefulSet的整个json体
func (c *configmap) UpdateConfigMap(client *kubernetes.Clientset, content, namespace string) (err error) {
	//content转成pod结构体
	var cm = &corev1.ConfigMap{}
	//反序列化成pod对象
	err = json.Unmarshal([]byte(content), &cm)
	if err != nil {
		zap.L().Error(fmt.Sprintf("反序列化失败, %v\n", err))
		return errors.New(fmt.Sprintf("反序列化失败, %v\n", err))
	}
	//更新pod
	_, err = client.CoreV1().ConfigMaps(namespace).Update(context.TODO(), cm, metav1.UpdateOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("更新 ConfigMap 失败, %v\n", err))
		return errors.New(fmt.Sprintf("更新 ConfigMap 失败, %v\n", err))
	}
	return nil
}
