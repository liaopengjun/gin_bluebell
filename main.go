package main

import (
	"context"
	"fmt"
	"gin_bluebell/controllers"
	"gin_bluebell/dao/mysql"
	"gin_bluebell/dao/redis"
	"gin_bluebell/logger"
	"gin_bluebell/pkg/snowflake"
	"gin_bluebell/routes"
	"gin_bluebell/settings"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	//1.加载配置
	if err := settings.Init(); err != nil {
		fmt.Printf("init settings failed err : %v\n", err)
		return
	}

	//2.初始化日志
	if err := logger.Init(settings.Conf.LogConfig,settings.Conf.Mode); err != nil {
		zap.L().Debug("init logger failed err", zap.Error(err))
		return
	}

	defer func(l *zap.Logger) {
		err := l.Sync()
		if err != nil {
			zap.L().Debug("init mysql failed err", zap.Error(err))
		}
	}(zap.L()) //追加缓存区log

	//3.初始化mysql redis连接
	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		zap.L().Debug("init mysql failed err", zap.Error(err))
		return
	}
	defer mysql.Close() //释放mysql资源
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		zap.L().Debug("init redis failed err", zap.Error(err))
		return
	}
	defer redis.Close() //释放redis资源
	//注册雪花算法类
	if err := snowflake.Init(1);err != nil{
		zap.L().Debug("init snowflake failed err", zap.Error(err))
		return
	}

	//4.注册路由
	r := routes.SetupRouter(settings.Conf.Mode)

	//gin框架内置校验翻译器
	if err := controllers.InitTrans("zh");err !=nil{
		zap.L().Debug("init trans failed err", zap.Error(err))
		return
	}
	//5.启动服务
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", settings.Conf.Port),
		Handler: r,
	}
	// 开启一个goroutine启动服务
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Error("listen: %s\n", zap.Error(err))
		}
	}()
	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	zap.L().Info("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Error("Server Shutdown: ", zap.Error(err))
	}
	zap.L().Info("Server exiting")

}
