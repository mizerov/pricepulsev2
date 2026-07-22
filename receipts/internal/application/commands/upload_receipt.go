package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mizerov/pricepulse/receipts/pkg/domains"
	"github.com/samber/do/v2"
)

type ReceiptProvider interface {
	FetchByQR(ctx context.Context, qrString domains.QRString) (*domains.Receipt, error)
}

type UploadReceiptCommand struct {
	receiptProvider   ReceiptProvider
	receiptRepository domains.ReceiptRepository
}

type UploadReceiptData struct {
	AccountID uuid.UUID
	QRString  domains.QRString
}

func NewUploadReceiptCommand(i do.Injector) (*UploadReceiptCommand, error) {
	receiptProvider, err := do.InvokeAs[ReceiptProvider](i)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke receipt provider: %w", err)
	}
	receiptRepository, err := do.InvokeAs[domains.ReceiptRepository](i)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke receipt repository: %w", err)
	}

	return &UploadReceiptCommand{
		receiptProvider:   receiptProvider,
		receiptRepository: receiptRepository,
	}, nil
}

func (self *UploadReceiptCommand) Execute(ctx context.Context, data UploadReceiptData) error {
	exists, err := self.receiptRepository.Exists(ctx, data.QRString.FN, data.QRString.FD, data.QRString.FP)
	if err != nil {
		return fmt.Errorf("failed to check exists: %w", err)
	}

	if exists {
		return nil
	}

	rawReceipt, err := self.receiptProvider.FetchByQR(ctx, data.QRString)
	if err != nil {
		return fmt.Errorf("failed to fetch raw receipt: %w", err)
	}

	rawReceipt.ID = uuid.New()
	rawReceipt.CreatedBy = data.AccountID
	rawReceipt.CreatedAt = time.Now()

	for i, item := range rawReceipt.Items {
		item.ID = uuid.New()
		item.ReceiptID = uuid.New()
		rawReceipt.Items[i] = item
	}

	err = self.receiptRepository.Save(ctx, rawReceipt)
	if err != nil {
		return fmt.Errorf("failed to save receipt: %w", err)
	}

	return nil
}
