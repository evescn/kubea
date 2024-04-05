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

var Pvc pvc

type pvc struct{}

// PvcsResp 定义列表的返回类型
type PvcsResp struct {
	Items []corev1.PersistentVolumeClaim `json:"items"`
	Total int                            `json:"total"`
}

// toCells 方法用于将 pvc 类型数组，转换成DataCell类型数组
func (c *pvc) toCells(std []corev1.PersistentVolumeClaim) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = pvcCell(std[i])
	}
	return cells
}

// fromCells 方法用于将DataCell类型数组，转换成 pvc 类型数组
func (c *pvc) fromCells(cells []DataCell) []corev1.PersistentVolumeClaim {
	configmaps := make([]corev1.PersistentVolumeClaim, len(cells))
	for i := range cells {
		configmaps[i] = corev1.PersistentVolumeClaim(cells[i].(pvcCell))
	}
	return configmaps
}

// GetPvcs 获取 PVC 列表
func (c *pvc) GetPvcs(client *kubernetes.Clientset, filterName, namespace string, limit, page int) (pvcsResp *PvcsResp, err error) {
	// context.TODO()用于声明一个空的context上下文，用于List方法内设置这个请求的超时（源码），这里的常用用法
	pvcList, err := client.CoreV1().PersistentVolumeClaims(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		zap.L().
			Error(fmt.Sprintf("获取 PVC 列表失败, %v\n", err))
		return nil, errors.New(fmt.Sprintf("获取 PVC 列表失败, %v\n", err))
	}
	//实例化dataSelector对象，把 d 结构体中获取到的 StatefulSet 列表转化为 dataSelector 结构体，方便使用 dataSelector 结构体中 过滤，排序，分页功能
	selectableData := &dataSelector{
		GenericDataList: c.toCells(pvcList.Items),
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
	pvcs := c.fromCells(data.GenericDataList)
	return &PvcsResp{
		Items: pvcs,
		Total: total,
	}, nil
}

// GetPvcDetail 获取 PVC 详情
func (c *pvc) GetPvcDetail(client *kubernetes.Clientset, pvcName, namespace string) (pvc *corev1.PersistentVolumeClaim, err error) {
	pvcDetail, err := client.CoreV1().PersistentVolumeClaims(namespace).Get(context.TODO(), pvcName, metav1.GetOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("获取 PVC 详情失败, %v\n", err))
		return nil, errors.New(fmt.Sprintf("获取 PVC 详情失败, %v\n", err))
	}
	return pvcDetail, nil
}

// DeletePvc 删除 PVC
func (c *pvc) DeletePvc(client *kubernetes.Clientset, pvcName, namespace string) (err error) {
	err = client.CoreV1().PersistentVolumeClaims(namespace).Delete(context.TODO(), pvcName, metav1.DeleteOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("删除 PVC 失败, %v\n", err))
		return errors.New(fmt.Sprintf("删除 PVC 失败, %v\n", err))
	}
	return nil
}

// UpdatePvc 更新 PVC
func (c *pvc) UpdatePvc(client *kubernetes.Clientset, content, namespace string) (err error) {
	//content转成pod结构体
	var pvcs = &corev1.PersistentVolumeClaim{}
	//反序列化成pod对象
	err = json.Unmarshal([]byte(content), &pvcs)
	if err != nil {
		zap.L().Error(fmt.Sprintf("反序列化失败, %v\n", err))
		return errors.New(fmt.Sprintf("反序列化失败, %v\n", err))
	}
	//更新pod
	_, err = client.CoreV1().PersistentVolumeClaims(namespace).Update(context.TODO(), pvcs, metav1.UpdateOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("更新 PVC 失败, %v\n", err))
		return errors.New(fmt.Sprintf("更新 PVC 失败, %v\n", err))
	}
	return nil
}
