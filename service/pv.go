package service

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var Pv pv

type pv struct{}

// PvsResp 定义列表的返回类型
type PvsResp struct {
	Items []corev1.PersistentVolume `json:"items"`
	Total int                       `json:"total"`
}

// toCells 方法用于将Node类型数组，转换成DataCell类型数组
func (p *pv) toCells(std []corev1.PersistentVolume) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = pvCell(std[i])
	}
	return cells
}

// fromCells 方法用于将DataCell类型数组，转换成Node类型数组
func (p *pv) fromCells(cells []DataCell) []corev1.PersistentVolume {
	nodes := make([]corev1.PersistentVolume, len(cells))
	for i := range cells {
		nodes[i] = corev1.PersistentVolume(cells[i].(pvCell))
	}
	return nodes
}

// GetPvs 获取 Pv 列表
func (p *pv) GetPvs(client *kubernetes.Clientset, filterName string, limit, page int) (pvsResp *PvsResp, err error) {
	pvList, err := client.CoreV1().PersistentVolumes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("获取 PV 列表失败, %v\n", err))
		return nil, errors.New(fmt.Sprintf("获取 PV 列表失败, %v\n", err))
	}
	selectableData := &dataSelector{
		GenericDataList: p.toCells(pvList.Items),
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
	//将[]DataCell类型的Service列表转为v1.Service列表
	pvs := p.fromCells(data.GenericDataList)
	return &PvsResp{
		Items: pvs,
		Total: total,
	}, nil
}

// GetPvDetail 获取 Pv 详情
func (p *pv) GetPvDetail(client *kubernetes.Clientset, pvName string) (pvs *corev1.PersistentVolume, err error) {
	pvDetail, err := client.CoreV1().PersistentVolumes().Get(context.TODO(), pvName, metav1.GetOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("获取 PV 详情失败, %v\n", err))
		return nil, errors.New(fmt.Sprintf("获取 PV 详情失败, %v\n", err))
	}
	return pvDetail, nil
}

// DeletePv 删除 Pv
func (p *pv) DeletePv(client *kubernetes.Clientset, pvName string) (err error) {
	err = client.CoreV1().PersistentVolumes().Delete(context.TODO(), pvName, metav1.DeleteOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("删除 PV 失败, %v\n", err))
		return errors.New(fmt.Sprintf("删除 PV 失败, %v\n", err))
	}
	return nil
}
