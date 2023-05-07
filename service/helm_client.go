package service

import (
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	"helm.sh/helm/v3/pkg/action"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"log"
	"os"
)

var HelmConfig helmConfig

type helmConfig struct{}

// GetAction 获取 helm action 配置
func (*helmConfig) GetAction(cluster, namespace string) (*action.Configuration, error) {
	// 获取 kubeconfig
	kubeconfig, ok := K8s.KubeConfMap[cluster]
	if !ok {
		logger.Error("actionConfig初始化失败,cluster不存在")
		return nil, errors.New("actionConfig初始化失败,cluster不存在")
	}

	// new 一个 actionConfig 对象
	actionConfig := new(action.Configuration)
	cf := &genericclioptions.ConfigFlags{
		KubeConfig:  &kubeconfig,
		ClusterName: &namespace,
	}
	if err := actionConfig.Init(cf, namespace, os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		logger.Error(fmt.Sprintf("actionConfig初始化失败, %v\n", err))
		return nil, errors.New(fmt.Sprintf("actionConfig初始化失败, %v\n", err))
	}
	return actionConfig, nil
}
