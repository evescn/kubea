package config

import "time"

const (
	WsAddr     = "0.0.0.0:8082"
	ListenAddr = "0.0.0.0:9090"
	// 1个集群使用为多集群
	Kubeconfigs = `{
		"DEV": "/Users/evescn/Documents/GitHub/kubea-demo/config/config",
		"TST": "/Users/evescn/Documents/GitHub/kubea-demo/config/config"
	}`
	// 查看容器日志时，显示的tail行数 tail -n 5000
	PodLogTailLine = 2000
	//数据库配置
	DbType = "mysql"
	DbHost = "localhost"
	DbPort = 3306
	DbName = "kubea_cicd_demo"
	DbUser = "root"
	DbPwd  = "123456"
	//打印mysql debug sql日志
	LogMode = false
	//连接池配置
	MaxIdleConns = 10               //最大空闲连接
	MaxOpenConns = 100              //最大连接数
	MaxLifeTime  = 30 * time.Second //最大生存时间
	//helm上传路径
	UploadPath = "/Users/evescn/Documents/GitLab/kubea-demo/demo/"

	//账号密码
	AdminUser = "admin"
	AdminPwd  = "123456"

	//cicd
	//触发tekton ci
	TektonKubeConfig = "C:\\custom\\project\\config"
	//CI上传镜像，CD下载镜像
	RegistryHost = "harbor.dayuan1997.com/adoo_k8s"
	//部署到的namespace
	DeployNamespace = "default"
	//CI使用，git的下载地址
	GitUrl = "http://192.168.0.14:30180"
	//编译记录使用
	TektonUrl = "http://192.168.0.14:32000/#/namespaces/%s/pipelineruns/%s"
)
