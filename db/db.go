package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"go.uber.org/zap"
	"kubea/model"
	"kubea/settings"
	"time"
)

var (
	isInit bool
	GORM   *gorm.DB
	err    error
)

func Init(cfg *settings.MySQLConfig) (err error) {
	//判断是否已经初始化
	if isInit {
		return
	}

	//组装数据库连接的数据
	//parseTime是查询结果是否自动解析为时间
	//loc是Mysql的时区设置
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DbName)

	GORM, err = gorm.Open(cfg.DbType, dsn)
	if err != nil {
		zap.L().Error("connect DB failed", zap.Error(err))
		return err
	}
	//打印sql语句
	GORM.LogMode(cfg.LogMode)
	//开启连接池
	GORM.DB().SetMaxIdleConns(cfg.MaxIdleConns)
	GORM.DB().SetMaxOpenConns(cfg.MaxOpenConns)
	GORM.DB().SetConnMaxLifetime(time.Duration(cfg.MaxLifeTime) * time.Second)

	//isInit = true
	GORM.AutoMigrate(
		model.App{},
		model.Chart{},
		model.Deploy{},
		model.DeployLog{},
		model.Event{},
		model.User{},
		model.Env{},
		model.Password{},
		model.Service{},
		model.Menu{},
		model.SubMenu{},
		model.SubSubMenu{},
		model.Role{},
		model.RoleMenuRelation{},
	)
	zap.L().Info("数据库连接成功")
	return
}

func Close() error {
	zap.L().Info("关闭数据库连接", zap.Error(err))
	return GORM.Close()
}
