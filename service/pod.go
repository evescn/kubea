package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	"io"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"kubea/config"
)

var Pod pod

type pod struct{}

// 定义列表的返回类型
type PodsResp struct {
	Items []corev1.Pod `json:"items"`
	Total int          `json:"total"`
}

// toCells 方法用于将pod类型数组，转换成DataCell类型数组
func (p *pod) toCells(std []corev1.Pod) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = podCell(std[i])
	}
	return cells
}

// fromCells 方法用于将DataCell类型数组，转换成pod类型数组
func (p *pod) fromCells(cells []DataCell) []corev1.Pod {
	pods := make([]corev1.Pod, len(cells))
	for i := range cells {
		pods[i] = corev1.Pod(cells[i].(podCell))
	}
	return pods
}

// GetPods 获取pod列表，支持过滤，排序，分页，
// client用于选择哪个集群
func (p *pod) GetPods(client *kubernetes.Clientset, filterName, namespace string, limit, page int) (podsResp *PodsResp, err error) {
	//context.TODO()用于声明一个空的context上下文，用于List方法内设置这个请求的超时（源码），这里的常
	//用用法
	//metav1.ListOptions{}用于过滤List数据，如使用label，field等
	podList, err := client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(fmt.Sprintf("获取 Pod 列表失败, %v\n", err))
		return nil, errors.New(fmt.Sprintf("获取 Pod 列表失败, %v\n", err))
	}
	//实例化dataSelector对象，把 p 结构体中获取到的 Pod 列表转化为 dataSelector 结构体，方便使用 dataSelector 结构体中 过滤，排序，分页功能
	selectableData := &dataSelector{
		GenericDataList: p.toCells(podList.Items),
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
	//将[]DataCell类型的pod列表转为v1.pod列表
	pods := p.fromCells(data.GenericDataList)
	return &PodsResp{
		Items: pods,
		Total: total,
	}, nil
}

// GetPodDetail 获取 Pod 详情
func (p *pod) GetPodDetail(client *kubernetes.Clientset, podName, namespace string) (pod *corev1.Pod, err error) {
	pod, err = client.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		logger.Error(fmt.Sprintf("获取 Pod 详情失败, %v\n", err))
		return nil, errors.New(fmt.Sprintf("获取 Pod 详情失败, %v\n", err))
	}
	return pod, nil
}

// DeletePod 删除 Pod
func (p *pod) DeletePod(client *kubernetes.Clientset, podName, namespace string) (err error) {
	err = client.CoreV1().Pods(namespace).Delete(context.TODO(), podName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error(fmt.Sprintf("删除 Pod 失败, %v\n", err))
		return errors.New(fmt.Sprintf("删除 Pod 失败, %v\n", err))
	}
	return nil
}

// UpdatePod 更新 Pod
// content就是pod的整个json体
func (p *pod) UpdatePod(client *kubernetes.Clientset, content, namespace string) (err error) {
	//content转成pod结构体
	var pods = &corev1.Pod{}
	//反序列化成pod对象
	err = json.Unmarshal([]byte(content), &pods)
	if err != nil {
		logger.Error(fmt.Sprintf("反序列化失败, %v\n", err))
		return errors.New(fmt.Sprintf("反序列化失败, %v\n", err))
	}
	//更新pod
	_, err = client.CoreV1().Pods(namespace).Update(context.TODO(), pods, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(fmt.Sprintf("更新 Pod 失败, %v\n", err))
		return errors.New(fmt.Sprintf("更新 Pod 失败, %v\n", err))
	}
	return nil
}

// GetPodContainer 获取 Pod 中的容器名
func (p *pod) GetPodContainer(client *kubernetes.Clientset, podName, namespace string) (containers []string, err error) {
	//获取pod详情
	pod, err := p.GetPodDetail(client, podName, namespace)
	if err != nil {
		return nil, err
	}
	//从pod对象中拿到容器名
	for _, container := range pod.Spec.Containers {
		containers = append(containers, container.Name)
	}
	return containers, nil
}

// GetPodLog 获取 Pod 中容器日志
func (p *pod) GetPodLog(client *kubernetes.Clientset, containerName, podName, namespace string) (log string, err error) {
	//设置日志的配置，容器名以及tail的行数
	lineLimit := int64(config.PodLogTailLine)
	option := &corev1.PodLogOptions{
		Container: containerName,
		TailLines: &lineLimit,
	}
	//获取request实例
	req := client.CoreV1().Pods(namespace).GetLogs(podName, option)
	//发起request请求，返回一个ioReadCloser类型（等同于response.body）
	podLogs, err := req.Stream(context.TODO())
	if err != nil {
		logger.Error(fmt.Sprintf("获取PodLog失败, %v\n", err))
		return "", errors.New(fmt.Sprintf("获取PodLog失败, %v\n", err))
	}
	defer podLogs.Close()
	//将response body写入缓冲区，目的是为了转成string返回
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		logger.Error(fmt.Sprintf("复制PodLog失败, %v\n", err))
		return "", errors.New(fmt.Sprintf("复制PodLog失败, %v\n", err))
	}
	return buf.String(), nil
}
