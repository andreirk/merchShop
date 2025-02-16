package main

import (
	"context"
	"fmt"
	"go/avito-test/config"
	"go/avito-test/internal/auth"
	"go/avito-test/internal/handlers"
	"go/avito-test/internal/health"
	"go/avito-test/internal/repositories"
	"go/avito-test/internal/services"
	"go/avito-test/pkg/db"
	"go/avito-test/pkg/midleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	SHUTDOWN_TIMEOUT = 5 * time.Second
)

type App struct {
	Handler  http.Handler
	Shutdown func()
}

func NewApp(conf *config.Config) *App {
	router := http.NewServeMux()
	dbConnection, dbClose := db.NewDb(conf)

	// Repositories
	userRepo := repositories.NewUserRepository(dbConnection)
	coinTransRepo := repositories.NewCoinTransactionRepository(dbConnection)
	orderRepo := repositories.NewOrderRepository(dbConnection)
	itemRepo := repositories.NewItemRepository(dbConnection)

	//Services
	authService := auth.NewAuthService(userRepo)
	userService := services.NewUserService(userRepo)
	coinService := services.NewCoinService(userRepo, coinTransRepo)
	orderService := services.NewOrderService(userRepo, itemRepo, orderRepo)

	//Handlers
	auth.NewHandler(router, auth.HandlerDeps{
		Config:      conf,
		AuthService: authService,
		UserService: userService,
	})
	handlers.NewInfoHandler(router, handlers.InfoDeps{
		UserService: userService,
		Config:      conf,
	})

	health.NewHealthHandler(router)
	handlers.NewSendCoinHandler(router, handlers.SendCoinDeps{
		Config:      conf,
		CoinService: coinService,
	})
	handlers.NewBuyItemHandler(router, handlers.BuyItemDeps{
		Config:       conf,
		OrderService: orderService,
	})

	//Midlewares
	midlewareStack := midleware.Chain(
		midleware.CORS,
		midleware.Logging,
		midleware.HttpError,
	)

	app := &App{
		Handler: midlewareStack(router),
		Shutdown: func() {
			dbClose()
		},
	}

	return app
}

func main() {
	conf := config.LoadConfig("local")
	app := NewApp(conf)
	server := http.Server{
		Addr:    "localhost:" + conf.Port,
		Handler: app.Handler,
	}

	fmt.Println("Server is listening on:", server.Addr)
	go func() {
		// service connections
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\\n", err)
		}
	}()

	// Handle OS signals to gracefully shut down
	stopChan := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	// Wait for a termination signal
	log.Println("Service is running. Press Ctrl+C to exit.")
	<-stopChan
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), SHUTDOWN_TIMEOUT)
	defer cancel()
	cleanupResources(app)

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	select {
	case <-ctx.Done():
		log.Println("timeout of " + SHUTDOWN_TIMEOUT.String() + " microseconds elapsed.")
	}
	log.Println("Server exiting")

}

// handles resource cleanup during shutdown
func cleanupResources(app *App) {
	log.Println("Cleaning up resources...")

	app.Shutdown()

	// ... cleanup other resources here as services, connections or workers

	log.Println("All resources cleaned up.")
}
