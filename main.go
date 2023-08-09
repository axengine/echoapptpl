package main

import (
	"context"
	"echoapptpl/config"
	"echoapptpl/dbo"
	"echoapptpl/handler"
	"echoapptpl/http"
	"echoapptpl/svc"
	"echoapptpl/task"
	"echoapptpl/util"
	"fmt"
	"github.com/axengine/utils/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	db := dbo.New(viper.GetString("database.dsn"))
	if viper.GetBool("dev") {
		db.Sync()
	}

	httpServer := http.New(handler.New(svc.New(db)))

	ctx, cancel := context.WithCancel(context.Background())
	taskServer := task.New(db)
	taskServer.Start(ctx)

	httpServer.Start(ctx, viper.GetString("http.bind"))

	// handle exit signal
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	<-exit
	cancel()

	log.Logger.Info("shutdown....")

	// wait gracefully shutdown
	httpServer.Stop(ctx)
	taskServer.Wait()

	log.Logger.Info("shutdown ok")
}

func init() {
	pflag.Bool("dev", false, "Dev mode")
	pflag.String("http.bind", ":8080", "Port to run the http server")
	pflag.String("database.dsn", "", "MySQL connection:user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local")

	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	// 设置环境变量前缀
	viper.SetEnvPrefix("HSA")
	viper.AllowEmptyEnv(true)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	// 读取环境变量
	viper.AutomaticEnv()

	// 设置配置文件名和路径
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Logger.Warn("Error reading config file", zap.Error(err))
	}

	// 解析到Config
	if err := viper.Unmarshal(&config.Cfg); err != nil {
		log.Logger.Warn("Error Unmarshal config file", zap.Error(err))
	}
	fmt.Println(util.JsonPretty(&config.Cfg))
}
