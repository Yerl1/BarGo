package middleware

import (
	"context"
	"net/http"
	"strings"

	"backend/internal/consumer-service/adapters/driver/myhttp/handlers"
	"backend/internal/mylogger"

	"github.com/golang-jwt/jwt"
)

type Middleware struct {
	ctx       context.Context
	jwtSecret string
	mylogger.Logger
}

func New(ctx context.Context, jwtSecret string, mylog mylogger.Logger) *Middleware {
	return &Middleware{
		ctx:       ctx,
		jwtSecret: jwtSecret,
		Logger:    mylog,
	}
}

func (m *Middleware) Wrap(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := m.With().Str("middleware", "consumer").Logger()
		// Middleware logic here (e.g., logging, authentication, etc.)
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			log.Warn().Msg("Missing Authorization header")
			handlers.JsonResponse(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.jwtSecret), nil
		})

		if err != nil || !token.Valid {
			log.Warn().Err(err).Msg("Invalid token")
			handlers.JsonResponse(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Warn().Msg("Invalid token claims")
			handlers.JsonResponse(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			log.Warn().Msg("Invalid user ID in token")
			handlers.JsonResponse(w, "Invalid user ID in token", http.StatusUnauthorized)
			return
		}

		role, ok := claims["role"].(string)
		if !ok {
			log.Warn().Msg("Invalid role in token")
			handlers.JsonResponse(w, "Invalid role in token", http.StatusUnauthorized)
			return
		}

		// ✅ Add user info to request context
		ctx := context.WithValue(r.Context(), "user_id", userID)
		ctx = context.WithValue(ctx, "role", role)
		r = r.WithContext(ctx)

		log.Info().Str("user_id", userID).Str("role", role).Msg("Middleware passed")
		next.ServeHTTP(w, r)
	}
}

func (m *Middleware) CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
