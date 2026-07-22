package proverkacheka

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/mizerov/pricepulse/receipts/pkg/domains"
	"github.com/samber/do/v2"
)

type ProverkaChekaClient struct {
	httpClient *http.Client
	token      string
}

type ProverkaCheckaConfig struct {
	HttpClient *http.Client
	Token      string
}

func NewProverkaChekaClient(i do.Injector) (*ProverkaChekaClient, error) {
	config, err := do.Invoke[ProverkaCheckaConfig](i)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke config: %w", err)
	}

	return &ProverkaChekaClient{httpClient: config.HttpClient, token: config.Token}, nil
}

func (c *ProverkaChekaClient) FetchByQR(ctx context.Context, qrString domains.QRString) (*domains.Receipt, error) {
	resp, err := c.GetReceiptData(ctx, qrString)
	if err != nil {
		return nil, fmt.Errorf("failed to get receipt data: %w", err)
	}

	if resp.Code != 1 {
		var errMsg string
		if len(resp.Data) > 0 && resp.Data[0] == '"' {
			_ = json.Unmarshal(resp.Data, &errMsg)
		}
		if errMsg == "" {
			errMsg = string(resp.Data)
		}
		return nil, fmt.Errorf("API returned error code: %d, message: %s", resp.Code, errMsg)
	}

	var data Data
	if err := json.Unmarshal(resp.Data, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal successful receipt data: %w", err)
	}

	dateTime, err := time.Parse("20060102T1504", qrString.T)
	if err != nil {
		return nil, fmt.Errorf("failed to parse receipt date time: %w", err)
	}

	domainData := &domains.Receipt{
		StoreName: data.JSON.User,
		StoreINN:  data.JSON.UserInn,
		TotalSum:  data.JSON.TotalSum,
		DateTime:  dateTime,
		FN:        data.JSON.FiscalDriveNo,
		FD:        fmt.Sprintf("%d", data.JSON.FiscalDocNo),
		FP:        fmt.Sprintf("%d", data.JSON.FiscalSign),
		Items:     make([]domains.ReceiptItem, 0, len(data.JSON.Items)),
	}

	for _, item := range data.JSON.Items {
		domainData.Items = append(domainData.Items, domains.ReceiptItem{
			Name:     item.Name,
			Price:    item.Price,
			Quantity: item.Quantity,
			Amount:   item.Sum,
		})
	}

	return domainData, nil
}

func (c *ProverkaChekaClient) GetReceiptData(ctx context.Context, qrString domains.QRString) (*Response, error) {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	_ = writer.WriteField("qrraw", qrString.String())
	_ = writer.WriteField("token", c.token)

	err := writer.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://proverkacheka.com/api/v1/check/get", payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result Response
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}
