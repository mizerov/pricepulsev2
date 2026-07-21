package http

import (
	"errors"
	"net/http"

	"github.com/mizerov/pricepulse/auth/internal/application/commands"
	"github.com/mizerov/pricepulse/auth/pkg/domains"

	"github.com/mizerov/router"
)

type IssueJWTRequest struct {
	Email    string `json:"email" valid:"email,required"`
	Password string `json:"password" valid:"required"`
}

type IssueJWTResponse struct {
	SignedToken  string `json:"signedToken"`
	RefreshToken string `json:"refreshToken"`
}

func (self *Handler) IssueJWT(ctx router.ContextWithBody[IssueJWTRequest]) error {
	requestBody, err := ctx.Body()
	if err != nil {
		return ctx.Error(http.StatusBadRequest, ErrorResponse{
			Title:  "invalid input body",
			Detail: err.Error(),
		}, err)
	}

	tokenPair, err := self.app.Commands.IssueJWT.Execute(ctx, commands.IssueJWTData{
		Email:    requestBody.Email,
		Password: requestBody.Password,
	})
	switch {
	case err == nil:
	case errors.Is(err, domains.ErrAccountNotFound), errors.Is(err, domains.ErrInvalidPassword):
		return ctx.Error(http.StatusUnauthorized, ErrorResponse{
			Title:  "invalid credentials",
			Detail: err.Error(),
		}, err)
	default:
		return ctx.Error(http.StatusInternalServerError, ErrorResponse{
			Title:  "internal server error",
			Detail: err.Error(),
		}, err)
	}

	return ctx.JSON(http.StatusOK, IssueJWTResponse{
		SignedToken:  tokenPair.SignedToken,
		RefreshToken: tokenPair.RefreshToken,
	})
}

type RefreshJWTRequest struct {
	RefreshToken string `json:"refreshToken" valid:"required"`
}

type RefreshJWTResponse struct {
	SignedToken  string `json:"signedToken"`
	RefreshToken string `json:"refreshToken"`
}

func (self *Handler) RefreshJWT(ctx router.ContextWithBody[RefreshJWTRequest]) error {
	requestBody, err := ctx.Body()
	if err != nil {
		return ctx.Error(http.StatusBadRequest, ErrorResponse{
			Title:  "invalid input body",
			Detail: err.Error(),
		}, err)
	}

	tokenPair, err := self.app.Commands.RefreshJWT.Execute(ctx, commands.RefreshJWTData{
		RefreshToken: requestBody.RefreshToken,
	})
	switch {
	case err == nil:
	case errors.Is(err, domains.ErrSessionNotFound), errors.Is(err, domains.ErrSessionIsBlocked):
		return ctx.Error(http.StatusUnauthorized, ErrorResponse{
			Title:  "refresh token not found or blocked",
			Detail: err.Error(),
		}, err)
	default:
		return ctx.Error(http.StatusInternalServerError, ErrorResponse{
			Title:  "internal server error",
			Detail: err.Error(),
		}, err)
	}

	return ctx.JSON(http.StatusOK, RefreshJWTResponse{
		SignedToken:  tokenPair.SignedToken,
		RefreshToken: tokenPair.RefreshToken,
	})
}
