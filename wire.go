//go:build wireinject
// +build wireinject

package main

import (
	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/infras"
	"github.com/evermos/boilerplate-go/internal/domain/auth"
	"github.com/evermos/boilerplate-go/internal/domain/users"
	"github.com/evermos/boilerplate-go/internal/handlers"
	"github.com/evermos/boilerplate-go/transport/http"
	"github.com/evermos/boilerplate-go/transport/http/middleware"
	"github.com/evermos/boilerplate-go/transport/http/router"
	"github.com/google/wire"
)

// Wiring for configurations.
var configurations = wire.NewSet(
	configs.Get,
)

// Wiring for persistences.
var persistences = wire.NewSet(
	infras.ProvideMySQLConn,
)

// Wiring for domain.
var domainAuth = wire.NewSet(
	// AuthService interface and implementation
	auth.ProvideAuthServiceImpl,
	wire.Bind(new(auth.AuthService), new(*auth.AuthServiceImpl)),
	// AuthRepository interface and implementation
	auth.ProvideAuthRepositoryMySQL,
	wire.Bind(new(auth.AuthRepository), new(*auth.AuthRepositoryMySQL)),
)

var domainUser = wire.NewSet(
	// UserService interface and implementation
	users.ProvideUserServiceImpl,
	wire.Bind(new(users.UserService), new(*users.UserServiceImpl)),
	// UserRepository interface and implementation
	users.ProvideUserRepositoryMySQL,
	wire.Bind(new(users.UserRepository), new(*users.UserRepositoryMySQL)),
)

// Wiring for all domains.
var domains = wire.NewSet(
	domainAuth,
	domainUser,
)

var authMiddleware = wire.NewSet(
	middleware.ProvideAuthentication,
)

// Wiring for HTTP routing.
var routing = wire.NewSet(
	wire.Struct(new(router.DomainHandlers), "AuthHandler", "UserHandler"),
	handlers.ProvideAuthHandler,
	handlers.ProvideUserHandler,
	router.ProvideRouter,
)

// Wiring for everything.
func InitializeService() *http.HTTP {
	wire.Build(
		// configurations
		configurations,
		// persistences
		persistences,
		// middleware
		authMiddleware,
		// domains
		domains,
		// routing
		routing,
		// selected transport layer
		http.ProvideHTTP)
	return &http.HTTP{}
}
