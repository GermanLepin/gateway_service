package validate_token_service

// func RequireAuth(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
// 		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 		defer cancel()

// 		logger := logging.LoggerFromContext(ctx)

// 		tokenString := r.Header.Get("Authorization")
// 		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 				err := fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 				logger.Error("token method is incorrect", zap.Error(err))
// 				return nil, err
// 			}

// 			return []byte(os.Getenv("SECRET")), nil
// 		})
// 		if err != nil {
// 			logger.Error("token parsing is failed", zap.Error(err))
// 			jsonwrapper.ErrorJSON(rw, err, http.StatusUnauthorized)
// 			return
// 		}

// 		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 			if float64(time.Now().Unix()) > claims["exp"].(float64) {
// 				logger.Error("token is expired", zap.Error(err))
// 				jsonwrapper.ErrorJSON(rw, err, http.StatusUnauthorized)
// 			}

// 			next.ServeHTTP(rw, r)
// 		} else {
// 			logger.Error("token claims error", zap.Error(err))
// 			jsonwrapper.ErrorJSON(rw, err, http.StatusUnauthorized)
// 			return
// 		}
// 	})
// }

// func (s *service) ValidateTokenRequest(
// 	ctx context.Context,
// 	w http.ResponseWriter,
// 	loginRequest dto.LoginRequest,
// ) (session dto.Session, err error) {
// 	logger := logging.LoggerFromContext(ctx)

// 	jsonData, err := json.MarshalIndent(loginRequest, "", "\t")
// 	if err != nil {
// 		logger.Error("login request marshalling is failed", zap.Error(err))
// 		return session, err
// 	}

// 	authenticationServiceURL := "http://authentication-service/login"
// 	request, err := http.NewRequest(http.MethodPost, authenticationServiceURL, bytes.NewBuffer(jsonData))
// 	if err != nil {
// 		logger.Error("cannot reach out to the authentication-service", zap.Error(err))
// 		return session, err
// 	}
// 	request.Close = true

// 	client := &http.Client{}
// 	response, err := client.Do(request)
// 	if err != nil {
// 		logger.Error("sending an HTTP request is a failure", zap.Error(err))
// 		return session, err
// 	}
// 	defer response.Body.Close()

// 	// make sure we get back the correct status code
// 	if response.StatusCode == http.StatusUnauthorized {
// 		err := errors.New("invalid credentials")
// 		logger.Error("invalid credentials", zap.Error(err))
// 		return session, err
// 	} else if response.StatusCode != http.StatusAccepted {
// 		err := errors.New("error calling the authentication-service")
// 		logger.Error("error calling the authentication-service", zap.Error(err))
// 		return session, err
// 	}

// 	var loginResponse dto.LoginResponse
// 	if err = json.NewDecoder(response.Body).Decode(&loginResponse); err != nil {
// 		jsonwrapper.ErrorJSON(w, err, http.StatusInternalServerError)
// 		logger.Error("the decoding of the login request is failed", zap.Error(err))
// 		return
// 	}

// 	session = dto.Session{
// 		ID:                    loginResponse.SessionID,
// 		IsBlocked:             loginResponse.IsBlocked,
// 		AccessToken:           loginResponse.AccessToken,
// 		AccessTokenExpiresAt:  loginResponse.AccessTokenExpiresAt,
// 		RefreshToken:          loginResponse.RefreshToken,
// 		RefreshTokenExpiresAt: loginResponse.RefreshTokenExpiresAt,
// 		UserID:                loginResponse.UserID,
// 	}

// 	return session, nil
// }

// func New(userRepository UserRepository) *service {
// 	return &service{
// 		userRepository: userRepository,
// 	}
// }

// type service struct {
// 	userRepository UserRepository
// }
