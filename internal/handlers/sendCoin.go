package handlers

import (
	"fmt"
	"go/avito-test/config"
	"go/avito-test/internal/services"
	"go/avito-test/pkg/midleware"
	"go/avito-test/pkg/request"
	"go/avito-test/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SendCoinRequestDTO struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

type SendCoinDeps struct {
	CoinService services.ICoinService
	Config      *config.Config
}

type SendCoinHandler struct {
	CoinService services.ICoinService
	Config      *config.Config
}

func NewSendCoinHandler(router *http.ServeMux, deps SendCoinDeps) *SendCoinHandler {
	handler := &SendCoinHandler{
		CoinService: deps.CoinService,
	}
	router.Handle("POST /api/sendCoin", midleware.CheckAuthed(handler.SendCoin(), deps.Config))
	return handler
}

func (h *SendCoinHandler) SendCoin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentUserName, ok := r.Context().Value(midleware.ContextUserNameKey).(string)
		if ok {
			println(currentUserName)
		}
		body, err := request.HandleBody[SendCoinRequestDTO](&w, r)
		if err != nil {
			return
		}
		fmt.Println(" ---- Body here", body)
		// Business logic
		if err := h.CoinService.SendCoins(currentUserName, body.ToUser, body.Amount); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response.JsonResponse(w, gin.H{"message": "Coins sent successfully"}, http.StatusOK)
	}
}
