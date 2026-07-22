package http

import (
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strconv"

	"github.com/mizerov/pricepulse/receipts/internal/application/commands"
	"github.com/mizerov/pricepulse/receipts/pkg/domains"
	"github.com/mizerov/router"
)

type UploadReceiptByQRStringRequest struct {
	QRString string `json:"qr_string"`
}

func (self UploadReceiptByQRStringRequest) ToQRString() (domains.QRString, error) {
	parsed, err := url.ParseQuery(self.QRString)
	if err != nil {
		return domains.QRString{}, fmt.Errorf("invalid QR string format")
	}

	fn := parsed.Get("fn")
	fp := parsed.Get("fp")
	fd := parsed.Get("i")
	t := parsed.Get("t")
	s := parsed.Get("s")
	n := parsed.Get("n")

	if slices.Contains([]string{fn, fp, fd, t, s, n}, "") {
		return domains.QRString{}, fmt.Errorf("QR string contains empty values")
	}

	fSum, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return domains.QRString{}, fmt.Errorf("sum fields should be a number")
	}

	parsedN, err := strconv.ParseInt(n, 10, 64)
	if err != nil {
		return domains.QRString{}, fmt.Errorf("n should be a number")
	}

	return domains.QRString{
		T:  t,
		S:  fSum,
		FN: fn,
		FD: fd,
		FP: fp,
		N:  parsedN,
	}, nil
}

func (self *Handler) UploadReceiptByQRString(ctx router.ContextWithBody[UploadReceiptByQRStringRequest]) error {
	reqBody, err := ctx.Body()
	if err != nil {
		return ctx.Error(http.StatusBadRequest, ErrorResponse{
			Title:  "invalid input body",
			Detail: err.Error(),
		}, err)
	}

	qrString, err := reqBody.ToQRString()
	if err != nil {
		return ctx.Error(http.StatusBadRequest, ErrorResponse{
			Title:  "invalid input body",
			Detail: err.Error(),
		}, err)
	}

	accountID, ok := GetAccountIDFromContext(ctx)
	if !ok {
		return ctx.Error(http.StatusUnauthorized, ErrorResponse{
			Title: "unauthorized user",
		}, err)
	}

	err = self.app.Commands.UploadReceipt.Execute(ctx, commands.UploadReceiptData{
		AccountID: accountID,
		QRString:  qrString,
	})
	if err != nil {
		return ctx.Error(http.StatusInternalServerError, ErrorResponse{
			"",
			err.Error(),
		}, err)
	}

	return ctx.Status(http.StatusNoContent)
}
