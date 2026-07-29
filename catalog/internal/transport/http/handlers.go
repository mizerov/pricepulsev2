package http

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/mizerov/pricepulse/catalog/internal/application/commands"
	"github.com/mizerov/router"
)

type proposeRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
type approveParams struct {
	ID string `json:"id"` // It might be mapped from query or path? Wait, standard mux path params aren't in URL.Query().
}
type searchParams struct {
	Q string `json:"q" validate:"required"`
}

func (h *Handler) Search(ctx router.ContextWithParams[searchParams]) error {
	params, err := ctx.Params()
	if err != nil {
		return ctx.Error(http.StatusBadRequest, ErrorResponse{
			Title:  "invalid input params",
			Detail: err.Error(),
		}, err)
	}

	products, err := h.app.Queries.SearchProducts.Execute(ctx, params.Q)
	if err != nil {
		return ctx.Error(http.StatusInternalServerError, err.Error(), err)
	}

	return ctx.JSON(http.StatusOK, products)
}

func (h *Handler) Propose(ctx router.ContextWithBody[proposeRequest]) error {
	reqBody, err := ctx.Body()
	if err != nil {
		return ctx.Error(http.StatusBadRequest, ErrorResponse{
			Title:  "invalid input body",
			Detail: err.Error(),
		}, err)
	}

	accountIDStr := req.Header.Get("X-User-Id")
	if accountIDStr == "" {
		return ctx.Error(http.StatusUnauthorized, "unauthorized", nil)
	}

	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		return ctx.Error(http.StatusBadRequest, "invalid user id", err)
	}

	product, err := h.app.Commands.ProposeProduct.Execute(ctx, commands.ProposeProductData{
		AccountID:   accountID,
		Name:        reqBody.Name,
		Description: reqBody.Description,
	})
	if err != nil {
		return ctx.Error(http.StatusInternalServerError, err.Error(), err)
	}

	return ctx.JSON(http.StatusCreated, product)
}

func (h *Handler) Approve(ctx router.ContextWithParams[approveParams]) error {
	req := ctx.Request()
	role := req.Header.Get("X-User-Role")
	if role != "admin" {
		return ctx.Error(http.StatusForbidden, "forbidden", nil)
	}

	// Go 1.22 standard mux PathValue
	idParam := req.PathValue("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Error(http.StatusBadRequest, "invalid product id", err)
	}

	if err := h.app.Commands.ApproveProduct.Execute(ctx, id); err != nil {
		return ctx.Error(http.StatusInternalServerError, err.Error(), err)
	}

	return ctx.JSON(http.StatusOK, map[string]string{"status": "ok"})
}
