package service

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"kubea/settings"
)

var K8s k8s

type k8s struct {
	ClientMap   map[string]*kubernetes.Clientset
	KubeConfMap map[string]string
}

func (k *k8s) GetClient(cluster string) (*kubernetes.Clientset, error) {
	client, ok := k.ClientMap[cluster]
	if !ok {
		zap.L().Error(fmt.Sprintf("集群：%s 不存在，无法获取 client\n", client))
		return nil, errors.New(fmt.Sprintf("集群：%s 不存在，无法获取 client\n", client))
	}
	return client, nil
}

func (k *k8s) Init(cfg *settings.KubeConfigs) {
	//mp := map[string]string{}
	//
	//if err := json.Unmarshal([]byte(config.Kubeconfigs), &mp); err != nil {
	//	panic(fmt.Sprintf("Kubeconfigs 反序列化失败 %v\n", err))
	//}

	mp := map[string]string{
		//"DEV": cfg.DEV,
		"TST": cfg.TST,
	}

	k.ClientMap = map[string]*kubernetes.Clientset{}

	k.KubeConfMap = mp
	for key, value := range mp {
		conf, err := clientcmd.BuildConfigFromFlags("", value)
		if err != nil {
			zap.L().Panic(fmt.Sprintf("集群 %s: 创建 K8S 配置失败 %v\n", key, err))
			//panic(fmt.Sprintf("集群 %s: 创建 K8S 配置失败 %v\n", key, err))
		}
		clientSet, err := kubernetes.NewForConfig(conf)
		if err != nil {
			zap.L().Panic(fmt.Sprintf("集群 %s: 创建 K8sClient 失败 %v\n", key, err))
			//panic(fmt.Sprintf("集群 %s: 创建 K8sClient 失败 %v\n", key, err))
		}

		k.ClientMap[key] = clientSet
		zap.L().Info(fmt.Sprintf("集群 %s: 创建 K8sClient 成功", key))
	}
}
