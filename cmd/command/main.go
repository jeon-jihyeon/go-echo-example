package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/jed-jeon/go-echo-example/components/api"
	"github.com/jed-jeon/go-echo-example/components/shop"
	"github.com/jed-jeon/go-echo-example/internal/config"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "go-echo-example",
		Short: "go-echo-example CLI",
	}

	rootCmd.AddCommand(seedCmd())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func seedCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "seed",
		Short: "Seed sample data into the database",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.Load()

			db, err := gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{})
			if err != nil {
				return fmt.Errorf("failed to connect database: %w", err)
			}

			// 필요한 컴포넌트만 선택적 Init
			shopComponent := shop.New(db)
			if err = shopComponent.Migrate(); err != nil {
				return fmt.Errorf("failed to migrate: %w", err)
			}

			creator := shopComponent.ProductCreator()
			products := []api.CreateProductCommand{
				{Name: "Laptop", Price: 1500000, Stock: 10},
				{Name: "Keyboard", Price: 80000, Stock: 50},
				{Name: "Mouse", Price: 45000, Stock: 100},
				{Name: "Monitor", Price: 350000, Stock: 20},
				{Name: "Headset", Price: 120000, Stock: 30},
			}

			for _, p := range products {
				if err = creator.CreateProduct(context.Background(), p); err != nil {
					return fmt.Errorf("failed to seed: %w", err)
				}
			}

			log.Println("seed completed successfully")
			return nil
		},
	}
}
