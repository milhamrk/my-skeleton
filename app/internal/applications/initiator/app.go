package initiator

import (
	"context"
	"errors"
	"fazz/app/internal/config"
	"fazz/app/internal/handlers/users"
	"fazz/app/internal/repository/postgres"
	usersService "fazz/app/internal/service/users"
	"fazz/app/pkg/auth"
	"fazz/app/pkg/database"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/uptrace/bun"
	"go.uber.org/zap"
)

type App struct {
	cfg        *config.Config
	logger     *zap.Logger
	router     *mux.Router
	database   *bun.DB
	httpServer *http.Server
}

func NewApp(cfg *config.Config, logger *zap.Logger) (*App, error) {
	logger.Info("router init")
	router := mux.NewRouter()

	db := postgres.NewPostgresDatabase(database.DatabaseParametersToDSN("postgres", cfg.PostgreSQL.Host,
		cfg.PostgreSQL.Database, cfg.PostgreSQL.Username, cfg.PostgreSQL.Password, false))

	tokenManager, err := auth.NewManager(cfg.SigningKey)
	if err != nil {
		logger.Error("error init token-manager", zap.Error(err))
	}
	appUsersRepository := postgres.NewUsersRepository(db)
	appUsersService := usersService.NewService(appUsersRepository, logger, tokenManager)
	appHandlerUsers := users.NewUsersHandler(appUsersService)
	appHandlerUsers.Register(router)

	return &App{cfg: cfg, logger: logger, router: router}, nil
}

func (a *App) startHTTP() {
	a.logger.Info("start HTTP")

	var listener net.Listener

	a.logger.Info(fmt.Sprintf("bind application to host: %s and port: %s", a.cfg.Listen.BindIP, a.cfg.Listen.Port))
	var err error
	listener, err = net.Listen("tcp", fmt.Sprintf("%s:%s", a.cfg.Listen.BindIP, a.cfg.Listen.Port))
	if err != nil {
		a.logger.Fatal("", zap.Error(err))
	}

	c := cors.New(cors.Options{
		AllowedMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodOptions, http.MethodDelete},
		AllowedOrigins:     []string{"http://localhost:3000", "http://localhost:8080"},
		AllowCredentials:   true,
		AllowedHeaders:     []string{"Location", "Charset", "Access-Control-Allow-Origin", "Content-Type", "content-type", "Origin", "Accept", "Content-Length", "Accept-Encoding", "X-CSRF-Token"},
		OptionsPassthrough: true,
		ExposedHeaders:     []string{"Location", "Authorization", "Content-Disposition"},
		Debug:              false,
	})

	handler := c.Handler(a.router)

	a.httpServer = &http.Server{
		Handler:      handler,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	a.logger.Info("application completely initialized and started")

	if err := a.httpServer.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			a.logger.Warn("server shutdown")
		default:
			a.logger.Fatal("error run listener", zap.Error(err))
		}
	}
	err = a.httpServer.Shutdown(context.Background())
	if err != nil {
		a.logger.Fatal("", zap.Error(err))
	}
}

func (a *App) Run() {
	a.startHTTP()
}
