package middleware

import (
	"crypto-temka/internal/delivery/middleware/auth"
	"crypto-temka/pkg/log"
)

type Middleware struct {
	logger  *log.Logs
	jwtUtil auth.JWTUtil
}

func InitMiddleware(logger *log.Logs) Middleware {
	jwtUtil := auth.InitJWTUtil()

	return Middleware{
		logger:  logger,
		jwtUtil: jwtUtil,
	}
}
