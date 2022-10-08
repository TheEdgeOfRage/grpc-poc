package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"

	"grpc-test/handlers"
)

func ginFormatterWithUTCAndBodySize(param gin.LogFormatterParams) string {
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}
	return fmt.Sprintf("[DUNE] %v |%s %3d %s| %13v | %6v bytes | %15s |%s %-7s %s %#v\n%s",
		param.TimeStamp.UTC().Format("2006/01/02 15:04:05.000"),
		statusColor, param.StatusCode, resetColor,
		param.Latency,
		param.BodySize,
		param.ClientIP,
		methodColor, param.Method, resetColor,
		param.Path,
		param.ErrorMessage,
	)
}

func setupGinEngine() (*gin.Engine, error) {
	engine := gin.New()
	// Middlewares are executed top to bottom in a stack-like manner
	engine.Use(
		gin.LoggerWithFormatter(ginFormatterWithUTCAndBodySize),
		gin.Recovery(), // Recovery needs to go before other middlewares to catch panics
		gzip.Gzip(gzip.BestSpeed),
	)
	engine.GET("/results", handlers.GetResults)
	engine.GET("/status", handlers.GetStatus)

	return engine, nil
}

func StartServer(engine *gin.Engine) {
	server := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: engine,
	}

	quit := make(chan os.Signal, 1)
	// handle Interrupt (ctrl-c) Term, used by `kill` et al, HUP which is commonly used to reload configs
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		<-quit
		if err := server.Shutdown(context.Background()); err != nil {
			log.Fatalf("failed to shut down server gracefully: %v", err)
		}
	}()

	if err := server.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}
}

func main() {
	var err error

	gin.SetMode("debug")
	engine, err := setupGinEngine()
	if err != nil {
		log.Fatalf("failed to set up gin engine:%v\n", err)
		return
	}

	StartServer(engine)
}
