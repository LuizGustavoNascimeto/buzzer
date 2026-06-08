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

	"backend-go/internal/server"
	"backend-go/internal/shared/infra"
	db "backend-go/pkg/gormutil/db"
	"backend-go/pkg/observability"

	"github.com/rollbar/rollbar-go"
)

func gracefulShutdown(apiServer *http.Server, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop() // Allow Ctrl+C to force shutdown

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

func main() {
	// Rollbar setup
	rollbar.SetToken(os.Getenv("ROLLBAR_TOKEN"))
	rollbar.SetEnvironment(os.Getenv("ROLLBAR_ENVIRONMENT"))
	rollbar.SetServerRoot("github.com/seu-org/seu-repo")

	// OTel setup
	otelShutdown, err := observability.Setup(context.Background())
	if err != nil {
		rollbar.Critical(err)
		log.Printf("OpenTelemetry setup failed: %v", err)
	}

	//DB setup
	conn, err := db.NewDBConn(db.DefaultConfig(), os.Getenv("POSTGRES_URL"))
	if err != nil {
		rollbar.Critical(err)
		log.Printf("Gorm Postgres setup failed: %v", err)
	}

	infra.RunMigrations(conn.Gorm)

	rollbar.WrapAndWait(func() {
		// OTel shutdown
		if otelShutdown != nil {
			defer func() {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				if err := otelShutdown(ctx); err != nil {
					rollbar.Warning(err) // reporta falha no shutdown também
					log.Printf("OpenTelemetry shutdown failed: %v", err)
				}
			}()
		}

		server := server.NewServer()
		done := make(chan bool, 1)
		go gracefulShutdown(server, done)
		fmt.Println(os.Getenv("ROLLBAR_ACCESS_TOKEN"))
		err = server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic(fmt.Sprintf("http server error: %s", err))
		}

		<-done
		log.Println("Graceful shutdown complete.")
	})

	rollbar.Close()
}
