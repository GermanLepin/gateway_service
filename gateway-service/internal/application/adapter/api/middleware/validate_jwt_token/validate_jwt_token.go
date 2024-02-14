package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"gateway-service/internal/application/dto"
	"gateway-service/internal/application/helper/jsonwrapper"
	"gateway-service/internal/application/helper/logging"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		logger := logging.LoggerFromContext(ctx)
		ctx = logging.ContextWithLogger(ctx, logger)

		fmt.Println(r.Cookie("Authorization"))
		tokenString, err := r.Cookie("Authorization")
		if err != nil {
			jsonwrapper.ErrorJSON(rw, err, http.StatusUnauthorized)
			logger.Error(
				"decoding of login request is failed",
				zap.Error(err),
			)
			return
		}

		token, err := jwt.Parse(tokenString.String(), func(token *jwt.Token) (interface{}, error) {

			return []byte(os.Getenv("SECRET")), nil
		})
		if err != nil {
			return
		}

		var loginRequest dto.PaymentRequest
		if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
			jsonwrapper.ErrorJSON(rw, err, http.StatusInternalServerError)
			logger.Error(
				"decoding of login request is failed",
				zap.Error(err),
			)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				jsonwrapper.ErrorJSON(rw, err, http.StatusUnauthorized)
			}

			// claims["user id"]
			// jwtClaims := jwt.MapClaims{
			// 	"user id": user.ID,
			// 	"email":   user.Email,
			// 	"name":    user.Name,
			// 	"exp":     expirationTime,
			// }

			fmt.Println(claims["email"])
		} else {
			jsonwrapper.ErrorJSON(rw, err, http.StatusUnauthorized)
		}
	})
}
