package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	nwv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var Ingress ingress

type ingress struct{}

// IngressesResp 定义列表的返回类型
type IngressesResp struct {
	Items []nwv1.Ingress `json:"items"`
	Total int            `json:"total"`
}

// IngressCreate 定义 IngressCreate 的结构体
type IngressCreate struct {
	Name      string                 `json:"name"`
	Namespace string                 `json:"namespace"`
	Lable     map[string]string      `json:"lable"`
	Hosts     map[string][]*HttpPath `json:"hosts"`
	Cluster   string                 `json:"cluster"`
}

// HttpPath 定义 ingress 的 path 结构体
type HttpPath struct {
	Path        string        `json:"path"`
	PathType    nwv1.PathType `json:"path_type"`
	ServiceName string        `json:"service_name"`
	ServicePort int32         `json:"service_port"`
}

// toCells 方法用于将 Ingress 类型数组，转换成 DataCell 类型数组
func (i *ingress) toCells(std []nwv1.Ingress) []DataCell {
	cells := make([]DataCell, len(std))
	for j := range std {
		cells[j] = ingressCell(std[j])
	}
	return cells
}

// fromCells 方法用于将 DataCell 类型数组，转换成 Ingress 类型数组
func (i *ingress) fromCells(cells []DataCell) []nwv1.Ingress {
	ingresses := make([]nwv1.Ingress, len(cells))
	for j := range cells {
		ingresses[j] = nwv1.Ingress(cells[j].(ingressCell))
	}
	return ingresses
}

// GetIngresses 获取 Ingress 列表，支持过滤，排序，分页，
// client用于选择哪个集群
func (i *ingress) GetIngresses(client *kubernetes.Clientset, filterName, namespace string, limit, page int) (ingressesResp *IngressesResp, err error) {
	//context.TODO()用于声明一个空的context上下文，用于List方法内设置这个请求的超时（源码），这里的常
	//用用法
	//metav1.ListOptions{}用于过滤List数据，如使用label，field等
	ingressList, err := client.NetworkingV1().Ingresses(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("获取 Ingress 列表失败, %v\n", err))
		return nil, errors.New(fmt.Sprintf("获取 Ingress 列表失败, %v\n", err))
	}
	//实例化dataSelector对象，把 p 结构体中获取到的 Ingress 列表转化为 dataSelector 结构体，方便使用 dataSelector 结构体中 过滤，排序，分页功能
	selectableData := &dataSelector{
		GenericDataList: i.toCells(ingressList.Items),
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
	//将[]DataCell类型的Ingress列表转为v1.Ingress列表
	ingresses := i.fromCells(data.GenericDataList)
	return &IngressesResp{
		Items: ingresses,
		Total: total,
	}, nil
}

// GetIngressDetail 获取 Ingress 详情
func (i *ingress) GetIngressDetail(client *kubernetes.Clientset, ingressName, namespace string) (Ingress *nwv1.Ingress, err error) {
	ingresses, err := client.NetworkingV1().Ingresses(namespace).Get(context.TODO(), ingressName, metav1.GetOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("获取 Ingress 详情失败, %v\n", err))
		return nil, errors.New(fmt.Sprintf("获取 Ingress 详情失败, %v\n", err))
	}
	return ingresses, nil
}

// DeleteIngress 删除 Ingress
func (i *ingress) DeleteIngress(client *kubernetes.Clientset, ingressName, namespace string) (err error) {
	err = client.NetworkingV1().Ingresses(namespace).Delete(context.TODO(), ingressName, metav1.DeleteOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("删除 Ingress 失败, %v\n", err))
		return errors.New(fmt.Sprintf("删除 Ingress 失败, %v\n", err))
	}
	return nil
}

// UpdateIngress 更新 Ingress
// content就是Ingress的整个json体
func (i *ingress) UpdateIngress(client *kubernetes.Clientset, content, namespace string) (err error) {
	//content转成Ingress结构体
	var ingresses = &nwv1.Ingress{}
	//反序列化成Ingress对象
	err = json.Unmarshal([]byte(content), &ingresses)
	if err != nil {
		zap.L().Error(fmt.Sprintf("反序列化失败, %v\n", err))
		return errors.New(fmt.Sprintf("反序列化失败, %v\n", err))
	}
	//更新Ingress
	_, err = client.NetworkingV1().Ingresses(namespace).Update(context.TODO(), ingresses, metav1.UpdateOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("更新 Ingress 失败, %v\n", err))
		return errors.New(fmt.Sprintf("更新 Ingress 失败, %v\n", err))
	}
	return nil
}

// CreateIngress 创建 Ingress
func (i *ingress) CreateIngress(client *kubernetes.Clientset, data *IngressCreate) (err error) {
	//将data中的数据组装成nwv1.Ingress对象
	ingresses := &nwv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      data.Name,
			Namespace: data.Namespace,
			Labels:    data.Lable,
		},
		Status: nwv1.IngressStatus{},
	}

	//声明nwv1.IngressRule和nwv1.HTTPIngressPath变量，后面用于数据组装
	//ingressRule代表的是Hosts
	var ingressRules = make([]nwv1.IngressRule, 0)

	//httpIngressPaths代表的是Paths
	var httpIngressPaths = make([]nwv1.HTTPIngressPath, 0)

	//第一层for循环是将host组装成nwv1.IngressRule类型的对象
	//一个host对应一个ingressrule，每隔ingressrule中包含一个host和多个path
	for key, value := range data.Hosts {
		// 先把 host 放进去
		ir := nwv1.IngressRule{
			Host: key,
			IngressRuleValue: nwv1.IngressRuleValue{
				HTTP: &nwv1.HTTPIngressRuleValue{Paths: nil}},
		}

		//第二层for循环是将path组装成nwv1.HTTPIngressPath类型的对象
		for _, httpPath := range value {
			hip := nwv1.HTTPIngressPath{
				Path:     httpPath.Path,
				PathType: &httpPath.PathType,
				Backend: nwv1.IngressBackend{
					Service: &nwv1.IngressServiceBackend{
						Name: httpPath.ServiceName,
						Port: nwv1.ServiceBackendPort{
							Number: httpPath.ServicePort,
						},
					},
				},
			}
			//将每个 hip 对象组装成数组
			httpIngressPaths = append(httpIngressPaths, hip)
		}
		// 给 Paths 赋值，前面置空了
		ir.IngressRuleValue.HTTP.Paths = httpIngressPaths
		// 将每个 ir 组装为数组
		ingressRules = append(ingressRules, ir)
	}
	//将ingressRules放到ingress中
	ingresses.Spec.Rules = ingressRules

	// 创建 Ingress
	_, err = client.NetworkingV1().Ingresses(data.Namespace).Create(context.TODO(), ingresses, metav1.CreateOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("创建 Ingress 失败, %v\n", err))
		return errors.New(fmt.Sprintf("创建 Ingress 失败, %v\n", err))
	}
	return nil
}
