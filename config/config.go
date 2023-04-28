package config

const (
	ListenAddr = "0.0.0.0:9000"
	// 1个集群使用为多集群
	Kubeconfigs = `{
		"TST-1": "/Users/evescn/Documents/GitHub/kubea-demo/config/config.config",
		"TST-2": "/Users/evescn/Documents/GitHub/kubea-demo/config/config.config"
	}`
	PodLogTailLine = 2000
)
