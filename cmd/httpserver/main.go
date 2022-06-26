package main

import (
	"context"
	"github.com/ShareChat/service-template/config"
	"github.com/ShareChat/service-template/pkg/application/httpserver"
	"github.com/ShareChat/service-template/pkg/domain/services"
	"github.com/ShareChat/service-template/pkg/infrastructure/dbmysql"
	"github.com/ShareChat/service-template/pkg/infrastructure/memory"
	"github.com/ShareChat/service-template/pkg/infrastructure/transport/rest"
	"github.com/ShareChat/service-template/third_party/assetmnger"
	"github.com/ShareChat/service-template/third_party/platlogger"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	serviceName = "service-template"
)


func main() {
	assetMng := assetmnger.Initialize()
	var config *config.Store = config.NewConfig(assetMng)
	logger, _ := platlogger.NewLogger(serviceName, config, platlogger.ConsoleOutput(true), platlogger.StackDriverOutput(true))
	var memStore = memory.NewMemoryStore()
	var userSvc = rest.NewUserSvc(config)
	var reviewStore = dbmysql.NewReviewStore(config)

	var appLogic services.AppInterface = services.NewAppLogic(memStore, userSvc, reviewStore, logger)

	idleConnsClosed := make(chan struct{})
	go func() {
		done := make(chan os.Signal, 1)
		signal.Notify(done, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
		<-done

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		httpserver.Shutdown(ctx)
		reviewStore.Close()
		defer cancel()
		close(idleConnsClosed)
	}()

	httpserver.NewServer(appLogic, logger)
	<-idleConnsClosed
}