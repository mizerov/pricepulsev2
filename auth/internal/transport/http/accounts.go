package http

import (
	"errors"
	"net/http"

	"github.com/mizerov/pricepulse/auth/internal/application/commands"
	"github.com/mizerov/pricepulse/auth/pkg/domains"

	"github.com/mizerov/router"
)

type CreateAccountRequest struct {
	Email    string `json:"email" valid:"email,required"`
	Password string `json:"password" valid:"required"`
}

func (self *Handler) CreateAccount(ctx router.ContextWithBody[CreateAccountRequest]) error {
	requestBody, err := ctx.Body()
	if err != nil {
		return ctx.Error(http.StatusBadRequest, ErrorResponse{
			Title:  "invalid input body",
			Detail: err.Error(),
		}, err)
	}

	err = self.app.Commands.CreateAccount.Execute(ctx, commands.CreateAccountData{
		Email:    requestBody.Email,
		Password: requestBody.Password,
	})
	switch {
	case err == nil:
	case errors.Is(err, domains.ErrEmailAlreadyExists):
		return ctx.Error(http.StatusConflict, ErrorResponse{
			Title:  "email address already exists",
			Detail: err.Error(),
		}, err)
	default:
		return ctx.Error(http.StatusInternalServerError, ErrorResponse{
			Title:  "internal server error",
			Detail: err.Error(),
		}, err)
	}

	return nil
}
