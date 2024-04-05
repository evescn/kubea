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

var Secret secret

type secret struct{}

// SecretsResp 定义列表的返回类型
type SecretsResp struct {
	Items []corev1.Secret `json:"items"`
	Total int             `json:"total"`
}

// toCells 方法用于将 configmap 类型数组，转换成DataCell类型数组
func (c *secret) toCells(std []corev1.Secret) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = secretCell(std[i])
	}
	return cells
}

// fromCells 方法用于将DataCell类型数组，转换成 configmap 类型数组
func (c *secret) fromCells(cells []DataCell) []corev1.Secret {
	configmaps := make([]corev1.Secret, len(cells))
	for i := range cells {
		configmaps[i] = corev1.Secret(cells[i].(secretCell))
	}
	return configmaps
}

// GetSecrets 获取 Secret 列表
func (c *secret) GetSecrets(client *kubernetes.Clientset, filterName, namespace string, limit, page int) (secretsResp *SecretsResp, err error) {
	// context.TODO()用于声明一个空的context上下文，用于List方法内设置这个请求的超时（源码），这里的常用用法
	secretList, err := client.CoreV1().Secrets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("获取 Secret 列表失败, %v\n", err))
		return nil, errors.New(fmt.Sprintf("获取 Secret 列表失败, %v\n", err))
	}
	//实例化dataSelector对象，把 d 结构体中获取到的 StatefulSet 列表转化为 dataSelector 结构体，方便使用 dataSelector 结构体中 过滤，排序，分页功能
	selectableData := &dataSelector{
		GenericDataList: c.toCells(secretList.Items),
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
	secrets := c.fromCells(data.GenericDataList)
	return &SecretsResp{
		Items: secrets,
		Total: total,
	}, nil
}

// GetSecretDetail 获取 Secret 详情
func (c *secret) GetSecretDetail(client *kubernetes.Clientset, secretName, namespace string) (secret *corev1.Secret, err error) {
	secretDetail, err := client.CoreV1().Secrets(namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("获取 Secret 详情失败, %v\n", err))
		return nil, errors.New(fmt.Sprintf("获取 Secret 详情失败, %v\n", err))
	}
	return secretDetail, nil
}

// DeleteSecret 删除 Secret
func (c *secret) DeleteSecret(client *kubernetes.Clientset, secretName, namespace string) (err error) {
	err = client.CoreV1().Secrets(namespace).Delete(context.TODO(), secretName, metav1.DeleteOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("删除 Secret 失败, %v\n", err))
		return errors.New(fmt.Sprintf("删除 Secret 失败, %v\n", err))
	}
	return nil
}

// UpdateSecret 更新 Secret
// content就是StatefulSet的整个json体
func (c *secret) UpdateSecret(client *kubernetes.Clientset, content, namespace string) (err error) {
	//content转成pod结构体
	var secrets = &corev1.Secret{}
	//反序列化成pod对象
	err = json.Unmarshal([]byte(content), &secrets)
	if err != nil {
		zap.L().Error(fmt.Sprintf("反序列化失败, %v\n", err))
		return errors.New(fmt.Sprintf("反序列化失败, %v\n", err))
	}
	//更新pod
	_, err = client.CoreV1().Secrets(namespace).Update(context.TODO(), secrets, metav1.UpdateOptions{})
	if err != nil {
		zap.L().Error(fmt.Sprintf("更新 Secret 失败, %v\n", err))
		return errors.New(fmt.Sprintf("更新 Secret 失败, %v\n", err))
	}
	return nil
}
