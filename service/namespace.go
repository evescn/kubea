package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var Namespace namespace

type namespace struct{}

// NamespacesResp 定义列表的返回类型
type NamespacesResp struct {
	Items []corev1.Namespace `json:"items"`
	Total int                `json:"total"`
}

// toCells 方法用于将Node类型数组，转换成DataCell类型数组
func (n *namespace) toCells(std []corev1.Namespace) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = namespaceCell(std[i])
	}
	return cells
}

// fromCells 方法用于将DataCell类型数组，转换成Node类型数组
func (n *namespace) fromCells(cells []DataCell) []corev1.Namespace {
	nodes := make([]corev1.Namespace, len(cells))
	for i := range cells {
		nodes[i] = corev1.Namespace(cells[i].(namespaceCell))
	}
	return nodes
}

// GetNamespaces 获取 Namespace 列表
func (n *namespace) GetNamespaces(client *kubernetes.Clientset, filterName string, limit, page int) (namespacesResp *NamespacesResp, err error) {
	namespaceList, err := client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(fmt.Sprintf("获取 Namespace 列表失败, %v\n", err))
		return nil, errors.New(fmt.Sprintf("获取 Namespace 列表失败, %v\n", err))
	}
	selectableData := &dataSelector{
		GenericDataList: n.toCells(namespaceList.Items),
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
	namespaces := n.fromCells(data.GenericDataList)
	return &NamespacesResp{
		Items: namespaces,
		Total: total,
	}, nil
}

// GetNamespaceDetail 获取 Namespace 详情
func (n *namespace) GetNamespaceDetail(client *kubernetes.Clientset, namespaceName string) (namespace *corev1.Namespace, err error) {
	namespaces, err := client.CoreV1().Namespaces().Get(context.TODO(), namespaceName, metav1.GetOptions{})
	if err != nil {
		logger.Error(fmt.Sprintf("获取 Namespace 详情失败, %v\n", err))
		return nil, errors.New(fmt.Sprintf("获取 Namespace 详情失败, %v\n", err))
	}
	return namespaces, nil
}

// DeleteNamespace 删除 Namespace
func (n *namespace) DeleteNamespace(client *kubernetes.Clientset, namespaceName string) (err error) {
	err = client.CoreV1().Namespaces().Delete(context.TODO(), namespaceName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error(fmt.Sprintf("删除 Namespace 失败, %v\n", err))
		return errors.New(fmt.Sprintf("删除 Namespace 失败, %v\n", err))
	}
	return nil
}
