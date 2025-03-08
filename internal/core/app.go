package core

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

// NewApp creates the Goblin application
func NewApp() *fx.App {
	return fx.New(
		AppModule,
		fx.Invoke(startServer),
	)
}

// Start/Stop lifecycle hooks
func startServer(lc fx.Lifecycle, engine *gin.Engine) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Println("Starting Goblin App...")
			go func() {
				engine.Run(":8080")
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Println("Stopping Goblin App...")
			return nil
		},
	})
}
