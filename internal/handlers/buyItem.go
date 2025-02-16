package handlers

import (
	"go/avito-test/config"
	"go/avito-test/internal/services"
	"go/avito-test/pkg/midleware"
	"go/avito-test/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BuyItemDeps struct {
	OrderService services.IOrderService
	Config       *config.Config
}

type BuyItemHandler struct {
	OrderService services.IOrderService
}

func NewBuyItemHandler(router *http.ServeMux, deps BuyItemDeps) *BuyItemHandler {
	handler := &BuyItemHandler{
		OrderService: deps.OrderService,
	}
	router.Handle("GET /api/buy/{item}", midleware.CheckAuthed(handler.BuyItem(), deps.Config))
	return handler
}

func (h *BuyItemHandler) BuyItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, ok := r.Context().Value(midleware.ContextUserNameKey).(string)
		if ok {
			println("Get username in BuyItem handler", username)
		}

		itemName := r.PathValue("item")

		// Business logic
		if err := h.OrderService.PurchaseItem(username, itemName); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response.JsonResponse(w, gin.H{"message": "Item purchased successfully"}, http.StatusOK)
	}
}
