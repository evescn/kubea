package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

var Servicev1 servicev1

type servicev1 struct{}

// ServicesResp 定义列表的返回类型
type ServicesResp struct {
	Items []corev1.Service `json:"items"`
	Total int              `json:"total"`
}

type ServiceCreate struct {
	Name          string            `json:"name"`
	Namespace     string            `json:"namespace"`
	Type          string            `json:"type"`
	ContainerPort int32             `json:"container_port"`
	Port          int32             `json:"port"`
	NodePort      int32             `json:"node_port"`
	Label         map[string]string `json:"label"`
	Cluster       string            `json:"cluster"`
}

// toCells 方法用于将Service类型数组，转换成DataCell类型数组
func (s *servicev1) toCells(std []corev1.Service) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = serviceCell(std[i])
	}
	return cells
}

// fromCells 方法用于将DataCell类型数组，转换成Service类型数组
func (s *servicev1) fromCells(cells []DataCell) []corev1.Service {
	svcs := make([]corev1.Service, len(cells))
	for i := range cells {
		svcs[i] = corev1.Service(cells[i].(serviceCell))
	}
	return svcs
}

// GetServices 获取Service列表，支持过滤，排序，分页，
// client用于选择哪个集群
func (s *servicev1) GetServices(client *kubernetes.Clientset, filterName, namespace string, limit, page int) (servicesResp *ServicesResp, err error) {
	//context.TODO()用于声明一个空的context上下文，用于List方法内设置这个请求的超时（源码），这里的常
	//用用法
	//metav1.ListOptions{}用于过滤List数据，如使用label，field等
	serviceList, err := client.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("获取 Service 列表失败, %v\n", err))
		return nil, errors.New(fmt.Sprintf("获取 Service 列表失败, %v\n", err))
	}
	//实例化dataSelector对象，把 p 结构体中获取到的 Service 列表转化为 dataSelector 结构体，方便使用 dataSelector 结构体中 过滤，排序，分页功能
	selectableData := &dataSelector{
		GenericDataList: s.toCells(serviceList.Items),
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
	svcs := s.fromCells(data.GenericDataList)
	return &ServicesResp{
		Items: svcs,
		Total: total,
	}, nil
}

// GetServiceDetail 获取 Service 详情
func (s *servicev1) GetServiceDetail(client *kubernetes.Clientset, serviceName, namespace string) (service *corev1.Service, err error) {
	service, err = client.CoreV1().Services(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("获取 Service 详情失败, %v\n", err))
		return nil, errors.New(fmt.Sprintf("获取 Service 详情失败, %v\n", err))
	}
	return service, nil
}

// DeleteService 删除 Service
func (s *servicev1) DeleteService(client *kubernetes.Clientset, serviceName, namespace string) (err error) {
	err = client.CoreV1().Services(namespace).Delete(context.TODO(), serviceName, metav1.DeleteOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("删除 Service 失败, %v\n", err))
		return errors.New(fmt.Sprintf("删除 Service 失败, %v\n", err))
	}
	return nil
}

// UpdateService 更新 Service
// content就是Service的整个json体
func (s *servicev1) UpdateService(client *kubernetes.Clientset, content, namespace string) (err error) {
	//content转成Service结构体
	var services = &corev1.Service{}
	//反序列化成Service对象
	err = json.Unmarshal([]byte(content), &services)
	if err != nil {
		zap.L().Error(fmt.Sprintf("反序列化失败, %v\n", err))
		return errors.New(fmt.Sprintf("反序列化失败, %v\n", err))
	}
	//更新Service
	_, err = client.CoreV1().Services(namespace).Update(context.TODO(), services, metav1.UpdateOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("更新 Service 失败, %v\n", err))
		return errors.New(fmt.Sprintf("更新 Service 失败, %v\n", err))
	}
	return nil
}

// CreateService 创建 Service
func (s *servicev1) CreateService(client *kubernetes.Clientset, data *ServiceCreate) (err error) {
	//将data中的数据组装成corev1.Service对象
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      data.Name,
			Namespace: data.Namespace,
			Labels:    data.Label,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:     "http",
					Port:     data.Port,
					Protocol: "TCP",
					TargetPort: intstr.IntOrString{
						Type:   0,
						IntVal: data.ContainerPort,
					},
				},
			},
			Type:     corev1.ServiceType(data.Type),
			Selector: data.Label,
		},
		Status: corev1.ServiceStatus{},
	}

	// 根据 Service 类型来判断不同的配置
	if data.NodePort != 0 && data.Type == "NodePort" {
		service.Spec.Ports[0].NodePort = data.NodePort
	}

	// 创建 Service
	_, err = client.CoreV1().Services(data.Namespace).Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("创建 Service 失败, %v\n", err))
		return errors.New(fmt.Sprintf("创建 Service 失败, %v\n", err))
	}
	return nil
}
