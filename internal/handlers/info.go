package handlers

import (
	"fmt"
	"go/avito-test/config"
	"go/avito-test/internal/services"
	"go/avito-test/pkg/midleware"
	"go/avito-test/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type InfoDeps struct {
	UserService services.IUserService
	Config      *config.Config
}

type InfoHandler struct {
	UserService services.IUserService
}

func NewInfoHandler(router *http.ServeMux, deps InfoDeps) *InfoHandler {
	handler := &InfoHandler{
		UserService: deps.UserService,
	}
	router.Handle("GET /api/info", midleware.CheckAuthed(handler.Info(), deps.Config))
	return handler
}

func (h *InfoHandler) Info() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Info handler before getting username")
		username, ok := r.Context().Value(midleware.ContextUserNameKey).(string)
		if ok {
			println(username)
		}

		// Business logic
		user, err := h.UserService.GetUserInfo(username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response.JsonResponse(w, gin.H{
			"coins":     user.CoinBalance,
			"inventory": user.Orders,
			"coinHistory": gin.H{
				"received": user.ReceivedCoinTransactions,
				"sent":     user.SentCoinTransactions,
			},
		}, http.StatusOK)
	}
}
