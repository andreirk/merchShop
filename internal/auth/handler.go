package auth

import (
	"fmt"
	"go/avito-test/config"
	"go/avito-test/internal/services"
	"go/avito-test/pkg/jwt"
	"go/avito-test/pkg/request"
	"go/avito-test/pkg/response"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerDeps struct {
	*config.Config
	*AuthService
	*services.UserService
}

type Handler struct {
	Config      *config.Config
	AuthService *AuthService
	UserService services.IUserService
}

func NewHandler(router *http.ServeMux, deps HandlerDeps) {
	handler := &Handler{
		Config:      deps.Config,
		AuthService: deps.AuthService,
		UserService: deps.UserService,
	}
	router.HandleFunc("POST /api/auth", handler.Auth())
	router.HandleFunc("POST /auth/register", handler.Register())
	router.HandleFunc("POST /auth/login", handler.Login())
}

func (handler *Handler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := request.HandleBody[RegisterRequest](&w, r)
		if err != nil {
			return
		}
		log.Println("Register with payload", payload)
		_, err = handler.AuthService.Register(payload.Username, payload.Password)
		if err != nil {
			response.JsonResponse(w, nil, http.StatusConflict)
			return
		}
		res := RegisterResponse{
			RegisterSuccess: true,
			Message:         "Please login",
		}
		response.JsonResponse(w, res, http.StatusCreated)
		fmt.Println("Register")
	}
}

func (handler *Handler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userDTO, err := request.HandleBody[LoginRequestDTO](&w, r)
		if err != nil {
			return
		}
		username, err := handler.AuthService.Login(userDTO.Username, userDTO.Password)

		if err != nil {
			response.JsonResponse(w, nil, http.StatusUnauthorized)
			return
		}
		jwtToken, err := jwt.NewJwt(handler.Config.Auth.Secret).Sign(jwt.JwtData{
			Username: username,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		resp := LoginResponse{
			AccessToken: jwtToken,
		}
		response.JsonResponse(w, resp, http.StatusOK)
	}
}

func (handler *Handler) Auth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// parse and validate
		authDTO, err := request.HandleBody[AuthRequestDTO](&w, r)
		fmt.Println("hit auth, here is authDTO", authDTO, err)
		if err != nil {
			response.JsonResponse(w, map[string]string{"errors": "bad credentials"}, http.StatusBadRequest)
			return
		}

		// Create or find user
		user, err := handler.UserService.FindOrCreateUser(authDTO.Username)
		if err != nil {
			response.JsonResponse(w, gin.H{"errors": "Failed to create user"}, http.StatusUnauthorized)
			return
		}

		// Generate JWT token
		jwtToken, err := jwt.NewJwt(handler.Config.Auth.Secret).Sign(jwt.JwtData{
			Username: user.Username,
		})
		if err != nil {
			response.JsonResponse(w, gin.H{"errors": "Failed to create user"}, http.StatusInternalServerError)
			return
		}

		resp := AuthResponseDTO{
			AccessToken: jwtToken,
		}
		response.JsonResponse(w, resp, http.StatusOK)
	}
}
