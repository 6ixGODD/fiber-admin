package app

import (
	"context"
	"fmt"
	"time"
)

// ShutdownHandler hooks to close the application gracefully
func ShutdownHandler(ctx context.Context, app *App) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Stop Cron
	app.Tasks.Stop()
	// Sync Zap
	if err := app.Zap.Sync(); err != nil {
		fmt.Println("Failed to sync zap logger")
	}
	// Close Mongo
	if err := app.Mongo.Close(ctx); err != nil {
		fmt.Println("Failed to close mongo")
	}
	// Close Redis
	if err := app.Redis.Close(); err != nil {
		fmt.Println("Failed to close redis")
	}
	// Close Fiber
	if err := app.App.Shutdown(); err != nil {
		fmt.Println("Failed to close fiber")
		return err
	}

	return nil
}
