package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/jed-jeon/go-echo-example/components"
	"github.com/jed-jeon/go-echo-example/components/payment"
	"github.com/jed-jeon/go-echo-example/components/shop"
	"github.com/jed-jeon/go-echo-example/internal/config"
)

func main() {
	cfg := config.Load()

	// DB 초기화
	db, err := gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// componentRegistry: 컴포넌트 생성 및 DI
	shopComponent := shop.New(db)
	paymentComponent := payment.New(db)

	// 컴포넌트 간 의존성 주입
	paymentComponent.SetProductFinder(shopComponent.ProductFinder())

	// 컴포넌트 lifecycle
	comps := []components.Component{shopComponent, paymentComponent}

	for _, c := range comps {
		if err = c.Migrate(); err != nil {
			log.Fatalf("failed to migrate: %v", err)
		}
	}
	for _, c := range comps {
		if err = c.Init(); err != nil {
			log.Fatalf("failed to init component: %v", err)
		}
	}

	// Echo 설정
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// RESTServer 라우트 등록
	for _, c := range comps {
		srv := c.NewServer()
		group := e.Group(srv.RoutePrefix())
		srv.RegisterRoutes(group)
	}

	// net/http.Server로 래핑
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.ServerPort),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Graceful shutdown
	go func() {
		log.Printf("starting server on :%s", cfg.ServerPort)
		if srvErr := e.StartServer(httpServer); srvErr != nil && srvErr != http.ErrServerClosed {
			log.Fatalf("server error: %v", srvErr)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = e.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	// 컴포넌트 종료
	for _, c := range comps {
		if err = c.Shutdown(); err != nil {
			log.Printf("failed to shutdown component: %v", err)
		}
	}

	log.Println("server exited")
}
