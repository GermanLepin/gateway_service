package validate_jwt_token

import (
	"context"
	"fmt"

	"gateway-service/internal/application/helper/jsonwrapper"
	"gateway-service/internal/application/helper/logging"

	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		logger := logging.LoggerFromContext(ctx)

		tokenString := r.Header.Get("Authorization")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				err := fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				logger.Error("token method is incorrect", zap.Error(err))
				return nil, err
			}

			return []byte(os.Getenv("SECRET")), nil
		})
		if err != nil {
			logger.Error("token parsing is failed", zap.Error(err))
			jsonwrapper.ErrorJSON(rw, err, http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				logger.Error("token is expired", zap.Error(err))
				jsonwrapper.ErrorJSON(rw, err, http.StatusUnauthorized)
			}

			next.ServeHTTP(rw, r)
		} else {
			logger.Error("token claims error", zap.Error(err))
			jsonwrapper.ErrorJSON(rw, err, http.StatusUnauthorized)
			return
		}
	})
}
