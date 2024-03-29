package config

import "time"

const (
	WsAddr     = "0.0.0.0:8082"
	ListenAddr = "0.0.0.0:9000"
	// 1个集群使用为多集群
	Kubeconfigs = `{
                "DEV": "./config/dev-config",
                "TST": "./config/test-config"
    }`

	// 查看容器日志时，显示的tail行数 tail -n 5000
	PodLogTailLine = 2000
	//数据库配置
	DbType = "mysql"
	DbHost = "mysql"
	DbPort = 3306
	//DbHost = "10.0.0.101"
	//DbPort = 24858
	DbName = "kubea_cicd"
	DbUser = "kubea_cicd"
	DbPwd  = "kubea_cicd12321.."
	//打印mysql debug sql日志
	LogMode = false
	//连接池配置
	MaxIdleConns = 10               //最大空闲连接
	MaxOpenConns = 100              //最大连接数
	MaxLifeTime  = 30 * time.Second //最大生存时间
	//helm上传路径
	UploadPath = "/Users/evescn/Documents/GitLab/kubea/demo/"

	//账号密码
	AdminUser = "admin"
	AdminPwd  = "admin12321"

	//cicd 编译记录使用
	JenkinsUrl   = "https://test-jenkins.dayuan1997.com/view/%s/job/%s/"
	UserPassword = "admin:11b5ec1e1c83e2647675012d26e1381530"
)
