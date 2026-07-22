package application

import (
	"fmt"

	"github.com/mizerov/pricepulse/receipts/internal/application/commands"

	"github.com/samber/do/v2"
)

type Application struct {
	Commands Commands
}

type Commands struct {
	UploadReceipt *commands.UploadReceiptCommand
}

func NewApplication(i do.Injector) (*Application, error) {
	uploadReceiptCommand, err := do.Invoke[*commands.UploadReceiptCommand](i)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke upload receipt command: %w", err)
	}

	return &Application{
		Commands: Commands{
			UploadReceipt: uploadReceiptCommand,
		},
	}, nil
}
