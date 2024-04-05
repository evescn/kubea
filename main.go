package main

import (
	"context"
	"fmt"
	"kubea/db"
	"kubea/middle/snowflake"
	"kubea/routers"
	"kubea/service"
	"net/http"
	"os"
	"os/signal"
	"time"

	//"go.uber.org/zap"
	"go.uber.org/zap"
	"kubea/logger"

	"kubea/settings"
)

func main() {

	// 1. 加载配置
	if err := settings.Init(); err != nil {
		fmt.Printf("init settings failed, err:%v\n", err)
		return
	}

	// 2. 初始化日志
	if err := logger.Init(settings.Conf.LogConfig); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}
	defer zap.L().Sync()
	zap.L().Debug("logger init success...")

	// 3. 初始化MySQL连接
	if err := db.Init(settings.Conf.MySQLConfig); err != nil {
		zap.L().Error("init mysql failed, err:%v\n", zap.Error(err))
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			zap.L().Fatal("数据库关闭异常:", zap.Error(err))
		}
	}()

	// 4. 初始化k8s client
	service.K8s.Init(settings.Conf.KubeConfigs)

	// 5. 注册雪花算法 ID 生成器
	if err := snowflake.Init(settings.Conf.StartTime, settings.Conf.MachineID); err != nil {
		fmt.Printf("init snowflake failed, err:%v\n", err)
		return
	}

	// 6. 注册路由
	r := routers.Setup()

	// 7. 启动task
	//go func() {
	//	service.Event.WatchEventTask("DEV")
	//}()
	go func() {
		service.Event.WatchEventTask("TST")
	}()

	// 8. websocket 启动
	wsHandler := http.NewServeMux()
	wsHandler.HandleFunc("/ws", service.Terminal.WsHandler)
	ws := &http.Server{
		Addr:    fmt.Sprintf(":%d", settings.Conf.WsPort),
		Handler: wsHandler,
	}
	go func() {
		if err := ws.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listen: %s", zap.Error(err))
		}
	}()

	// 9. gin server 启动
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", settings.Conf.Port),
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listen: %s", zap.Error(err))
		}
	}()

	// 10. 优雅关闭server
	// 声明一个系统信号的channel，并监听他，如果没有信号，就一直阻塞，如果有，就继续执行
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// 11 设置ctx超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//cancel用于释放ctx
	defer cancel()

	// 12 关闭 websocket
	if err := ws.Shutdown(ctx); err != nil {
		zap.L().Fatal("Websocket关闭异常:", zap.Error(err))
	}
	zap.L().Info("Websocket退出成功")

	// 13 关闭 gin server
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Gin Server 关闭异常：", zap.Error(err))
	}
	zap.L().Info("Gin Server 退出成功")
}
