package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/jeon-jihyeon/go-echo-example/actors"
	"github.com/jeon-jihyeon/go-echo-example/handlers"
)

func main() {
	system := actor.NewActorSystem()
	registry := actors.NewRegistry(system)

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	root := system.Root

	shop := e.Group("/shop")
	handlers.NewProductHandler(root, registry.ProductPID).RegisterRoutes(shop)
	handlers.NewOrderHandler(root, registry.OrderPID).RegisterRoutes(shop)
	handlers.NewPaymentHandler(root, registry.PaymentPID).RegisterRoutes(e)

	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Printf("echo shutdown error: %v", err)
	}
	system.Shutdown()
	log.Println("server exited")
}
