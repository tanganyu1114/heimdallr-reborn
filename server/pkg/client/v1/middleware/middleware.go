package middleware

import (
	txpclientv1 "github.com/tanganyu1114/heimdallr-reborn/pkg/client/v1/transport"
)

// Middleware defines the middleware function type.
type Middleware func(factory txpclientv1.Factory) txpclientv1.Factory

var middlewares = make(map[string]Middleware)

// Register registers a middleware with the given name.
func Register(name string, m Middleware) {
	middlewares[name] = m
}

// GetMiddlewares returns all registered middlewares.
func GetMiddlewares() map[string]Middleware {
	return middlewares
}
