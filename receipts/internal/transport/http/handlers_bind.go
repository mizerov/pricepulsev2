package http

import (
	"errors"
	"net/http"

	"github.com/mizerov/pricepulse/receipts/internal/application/commands"
	"github.com/mizerov/router"
)

type bindRequest struct {
	ProductID string `json:"product_id" valid:"required"`
}

func (self *Handler) BindProduct(ctx router.ContextWithBody[bindRequest]) error {
	reqBody, err := ctx.Body()
	if err != nil {
		return ctx.Error(http.StatusBadRequest, err.Error(), err)
	}

	accountID, ok := GetAccountIDFromContext(ctx)
	if !ok {
		return ctx.Error(http.StatusUnauthorized, ErrorResponse{
			Title: "unauthorized",
		}, errors.New("unauthorized access"))
	}

	receiptID := ctx.Request().PathValue("id")
	itemID := ctx.Request().PathValue("item_id")

	if err := self.app.Commands.BindProduct.Execute(ctx, commands.BindProductData{
		AccountID: accountID,
		ReceiptID: receiptID,
		ItemID:    itemID,
		ProductID: reqBody.ProductID,
	}); err != nil {
		if err.Error() == "product does not exist or access denied" {
			return ctx.Error(http.StatusBadRequest, err.Error(), nil)
		}
		return ctx.Error(http.StatusInternalServerError, err.Error(), err)
	}

	return ctx.Status(http.StatusOK)
}
