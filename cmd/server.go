package cmd

import (
	"context"
	"eventdrivensystem/internal/domain"
	"eventdrivensystem/internal/handler/rest"
	"eventdrivensystem/internal/usecase"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/spf13/cobra"
)

var apiServerCmd = &cobra.Command{
	Use:   "api-server",
	Short: "Runs the API server",
	Run: func(cmd *cobra.Command, args []string) {
		NewServer()
	},
}

func NewServer() {
	ctx := context.Background()
	e := echo.New()
	e.Use(echoMiddleware.Recover())
	dp := GetAppDependency()

	dom := domain.NewDomain(dp.cfg, dp.db, dp.log)
	uc := usecase.NewUsecase(dp.cfg, dp.log, dom)

	rest.NewRouterHandler(e, dp.validator, uc).RegisterRoutes()

	go func() {
		address := fmt.Sprintf("%s:%d", dp.cfg.ApiServer.Host, dp.cfg.ApiServer.Port)
		if err := e.Start(address); err != nil {
			log.Info("shutting down the server -> ", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Info("shut down started")
	if err := e.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
		log.Errorf("error shutting down api server", err)
	}

	log.Info("shut down completed")

}
