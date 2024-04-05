package settings

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
)

// Conf 全局变量，用来保存程序的所有配置信息
var Conf = new(AppConfig)

type AppConfig struct {
	Name           string `mapstructure:"name"`
	Mode           string `mapstructure:"mode"`
	Version        string `mapstructure:"version"`
	StartTime      string `mapstructure:"start_time"`
	MachineID      int64  `mapstructure:"machine_id"`
	Port           int    `mapstructure:"port"`
	WsPort         int    `mapstructure:"ws_port"`
	PodLogTailLine int    `mapstructure:"pod_log_tail_line"`
	UploadPath     string `mapstructure:"upload_path"`
	*Admin         `mapstructure:"admin"`
	*LogConfig     `mapstructure:"log"`
	*KubeConfigs   `mapstructure:"kube_configs"`

	*MySQLConfig `mapstructure:"mysql"`
	//*RedisConfig `mapstructure:"redis"`
	*CiCd `mapstructure:"ci_cd"`

	*GitLab `mapstructure:"gitlab"`
}

type Admin struct {
	UserName string `mapstructure:"username"`
	PassWord string `mapstructure:"password"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type KubeConfigs struct {
	DEV string `mapstructure:"dev"`
	TST string `mapstructure:"tst"`
}

type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DbName       string `mapstructure:"dbname"`
	Port         int    `mapstructure:"port"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxLifeTime  int    `mapstructure:"max_life_time"`
	DbType       string `mapstructure:"db_type"`
	LogMode      bool   `mapstructure:"log_mode"`
}

//type RedisConfig struct {
//	Host     string `mapstructure:"host"`
//	Password string `mapstructure:"password"`
//	Port     int    `mapstructure:"port"`
//	DB       int    `mapstructure:"db"`
//	PoolSize int    `mapstructure:"pool_size"`
//}

type CiCd struct {
	CopyJobName       string `mapstructure:"copy_job_name"`
	JenkinsUrl        string `mapstructure:"jenkins_url"`
	UserPassword      string `mapstructure:"user_password"`
	CocosJenkinsUrl   string `mapstructure:"cocos_jenkins_url"`
	CocosUserPassword string `mapstructure:"cocos_user_password"`
}

type GitLab struct {
	GitLabUrl   string `mapstructure:"gitlab_url"`
	GitLabToken string `mapstructure:"gitlab_token"`
}

func Init() (err error) {

	// 从环境变量中读取 ENV 的值
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev" // 默认为 dev
	}

	fmt.Println("env: ", env)

	// 设置配置文件的名称和路径
	viper.SetConfigName("config." + env) // 使用环境变量的值拼接配置文件名
	viper.SetConfigType("yaml")          // 指定配置文件类型(专用于从远程获取配置信息时指定配置文件类型的)
	viper.AddConfigPath("./config/")     // 指定查找配置文件的路径（这里使用相对路径）

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file: %s\n", err)
		return err
	}

	// 把读取到的配置信息反序列化到 Conf 变量中
	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改了...")
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
		}
	})
	return
}
