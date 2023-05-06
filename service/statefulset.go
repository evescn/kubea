package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var StatefulSet statefulSet

type statefulSet struct{}

// StatefulSetResp 定义列表的返回类型
type StatefulSetResp struct {
	Items []appsv1.StatefulSet `json:"items"`
	Total int                  `json:"total"`
}

// toCells 方法用于将 ds 类型数组，转换成 DataCell 类型数组
func (s *statefulSet) toCells(std []appsv1.StatefulSet) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = statefulSetCell(std[i])
	}
	return cells
}

// fromCells 方法用于将 DataCell 类型数组，转换成 ds 类型数组
func (s *statefulSet) fromCells(cells []DataCell) []appsv1.StatefulSet {
	statefulsets := make([]appsv1.StatefulSet, len(cells))
	for i := range cells {
		statefulsets[i] = appsv1.StatefulSet(cells[i].(statefulSetCell))
	}
	return statefulsets
}

// GetStatefulSets 获取StatefulSet列表，支持过滤，排序，分页，
// client用于选择哪个集群
func (s *statefulSet) GetStatefulSets(client *kubernetes.Clientset, filterName, namespace string, limit, page int) (statefulSetResp *StatefulSetResp, err error) {
	// context.TODO()用于声明一个空的context上下文，用于List方法内设置这个请求的超时（源码），这里的常用用法
	//metav1.ListOptions{}用于过滤List数据，如使用label，field等
	stsList, err := client.AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(fmt.Sprintf("获取 StatefulSet 列表失败, %v\n", err))
		return nil, errors.New(fmt.Sprintf("获取 StatefulSet 列表失败, %v\n", err))
	}
	//实例化dataSelector对象，把 d 结构体中获取到的 StatefulSet 列表转化为 dataSelector 结构体，方便使用 dataSelector 结构体中 过滤，排序，分页功能
	selectableData := &dataSelector{
		GenericDataList: s.toCells(stsList.Items),
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
	statefulSets := s.fromCells(data.GenericDataList)
	return &StatefulSetResp{
		Items: statefulSets,
		Total: total,
	}, nil
}

// GetStatefulSetDetail 获取 StatefulSet 详情
func (s *statefulSet) GetStatefulSetDetail(client *kubernetes.Clientset, stsName, namespace string) (sts *appsv1.StatefulSet, err error) {
	stsDetail, err := client.AppsV1().StatefulSets(namespace).Get(context.TODO(), stsName, metav1.GetOptions{})
	if err != nil {
		logger.Error(fmt.Sprintf("获取 StatefulSet 详情失败, %v\n", err))
		return nil, errors.New(fmt.Sprintf("获取 StatefulSet 详情失败, %v\n", err))
	}
	return stsDetail, nil
}

// DeleteStatefulSet 删除 StatefulSet
func (s *statefulSet) DeleteStatefulSet(client *kubernetes.Clientset, stsName, namespace string) (err error) {
	err = client.AppsV1().StatefulSets(namespace).Delete(context.TODO(), stsName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error(fmt.Sprintf("删除 StatefulSet 失败, %v\n", err))
		return errors.New(fmt.Sprintf("删除 StatefulSet 失败, %v\n", err))
	}
	return nil
}

// UpdateStatefulSet 更新 StatefulSet
// content就是StatefulSet的整个json体
func (s *statefulSet) UpdateStatefulSet(client *kubernetes.Clientset, content, namespace string) (err error) {
	//content转成pod结构体
	var sts = &appsv1.StatefulSet{}
	//反序列化成pod对象
	err = json.Unmarshal([]byte(content), &sts)
	if err != nil {
		logger.Error(fmt.Sprintf("反序列化失败, %v\n", err))
		return errors.New(fmt.Sprintf("反序列化失败, %v\n", err))
	}
	//更新pod
	_, err = client.AppsV1().StatefulSets(namespace).Update(context.TODO(), sts, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(fmt.Sprintf("更新 StatefulSet 失败, %v\n", err))
		return errors.New(fmt.Sprintf("更新 StatefulSet 失败, %v\n", err))
	}
	return nil
}
