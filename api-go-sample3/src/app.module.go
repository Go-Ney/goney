package main

import (
	"api-go-sample3/config"
	"api-go-sample3/pkg/core"
)

type AppModule struct {
	Config *config.Config
	Core   *core.Application
}

func NewAppModule() *AppModule {
	cfg := config.Load()

	coreConfig := &core.Config{
		Port: cfg.Port,
		Database: core.DatabaseConfig{
			Host:     cfg.Database.Host,
			Port:     cfg.Database.Port,
			User:     cfg.Database.User,
			Password: cfg.Database.Password,
			Name:     cfg.Database.Name,
		},
		Grpc: core.GrpcConfig{
			Port: cfg.Grpc.Port,
		},
		Nats: core.NatsConfig{
			URL: cfg.Nats.URL,
		},
	}

	app := core.NewApplication(coreConfig)

	return &AppModule{
		Config: cfg,
		Core:   app,
	}
}

func (app *AppModule) Bootstrap() error {
	// Registrar módulos aquí
	// app.RegisterModule(&users.UsersModule{})
	// app.RegisterModule(&products.ProductsModule{})

	return nil
}
