package app

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Gursevak56/food-delivery-platform/services/rider-service/config"
	"github.com/Gursevak56/food-delivery-platform/services/rider-service/internal/handler"
	"github.com/Gursevak56/food-delivery-platform/services/rider-service/internal/repository"
	"github.com/Gursevak56/food-delivery-platform/services/rider-service/internal/service"
	"github.com/Gursevak56/food-delivery-platform/services/rider-service/pkg/auth"
)

type Container struct {
	Config              config.Config
	Logger              *slog.Logger
	Postgres            *sql.DB
	Mongo               *mongo.Client
	Repo                *repository.Repository
	Services            *service.Services
	AuthManager         *auth.Manager
	AuthHandler         *handler.AuthHandler
	RiderHandler        *handler.RiderHandler
	ShiftHandler        *handler.ShiftHandler
	DispatchHandler     *handler.DispatchHandler
	LocationHandler     *handler.LocationHandler
	FinanceHandler      *handler.FinanceHandler
	FeedbackHandler     *handler.FeedbackHandler
	NotificationHandler *handler.NotificationHandler
	SupportHandler      *handler.SupportHandler
	AdminHandler        *handler.AdminHandler
	HealthHandler       *handler.HealthHandler
	DocsHandler         *handler.DocsHandler
}

func NewContainer(cfg config.Config, logger *slog.Logger) (*Container, error) {
	postgres, mongoClient, err := initDatastores(cfg, logger)
	if err != nil {
		return nil, err
	}

	repo := repository.NewMemoryRepository(cfg.Auth.PasswordHashCost)
	manager := auth.NewManager(cfg.Auth.JWTSecret, cfg.Auth.JWTIssuer)
	services := service.NewServices(repo, cfg, manager)
	base := handler.NewBaseHandler()

	return &Container{
		Config:              cfg,
		Logger:              logger,
		Postgres:            postgres,
		Mongo:               mongoClient,
		Repo:                repo,
		Services:            services,
		AuthManager:         manager,
		AuthHandler:         &handler.AuthHandler{Base: base, Service: services.Auth},
		RiderHandler:        &handler.RiderHandler{Base: base, Service: services.Rider},
		ShiftHandler:        &handler.ShiftHandler{Base: base, Service: services.Shift},
		DispatchHandler:     &handler.DispatchHandler{Base: base, Service: services.Dispatch},
		LocationHandler:     &handler.LocationHandler{Base: base, Service: services.Dispatch},
		FinanceHandler:      &handler.FinanceHandler{Base: base, Service: services.Finance},
		FeedbackHandler:     &handler.FeedbackHandler{Base: base, Service: services.Feedback},
		NotificationHandler: &handler.NotificationHandler{Base: base, Service: services.Notification},
		SupportHandler:      &handler.SupportHandler{Base: base, Service: services.Support},
		AdminHandler:        &handler.AdminHandler{Base: base, Service: services.Admin},
		HealthHandler:       &handler.HealthHandler{Postgres: postgres, Mongo: mongoClient},
		DocsHandler:         &handler.DocsHandler{},
	}, nil
}

func (c *Container) Close() {
	if c.Postgres != nil {
		_ = c.Postgres.Close()
	}
	if c.Mongo != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = c.Mongo.Disconnect(ctx)
	}
}

func initDatastores(cfg config.Config, logger *slog.Logger) (*sql.DB, *mongo.Client, error) {
	var postgres *sql.DB
	var mongoClient *mongo.Client

	if cfg.Postgres.URL != "" {
		db, err := sql.Open("postgres", cfg.Postgres.URL)
		if err != nil {
			return nil, nil, fmt.Errorf("open postgres: %w", err)
		}
		db.SetMaxOpenConns(cfg.Postgres.MaxOpenConns)
		db.SetMaxIdleConns(cfg.Postgres.MaxIdleConns)
		db.SetConnMaxLifetime(cfg.Postgres.ConnMaxLifetime)
		if err := db.Ping(); err != nil {
			logger.Warn("postgres ping failed, keeping service in memory mode", "error", err)
		} else {
			postgres = db
		}
	}

	if cfg.Mongo.URI != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Mongo.URI))
		if err != nil {
			logger.Warn("mongo connect failed, keeping telemetry in memory mode", "error", err)
		} else {
			mongoClient = client
		}
	}

	return postgres, mongoClient, nil
}
