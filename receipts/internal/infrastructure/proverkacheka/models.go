package proverkacheka

import (
	"encoding/json"
	"time"
)

// Главный контейнер ответа
type Response struct {
	Code  int             `json:"code"`
	First int             `json:"first"`
	Data  json.RawMessage `json:"data"` // Используем RawMessage, так как при ошибке здесь приходит строка, а при успехе - объект
}

type Data struct {
	JSON ReceiptJSON `json:"json"`
}

// Содержимое чека
type ReceiptJSON struct {
	User          string   `json:"user"`                 // Название организации (АО ТАНДЕР)
	UserInn       string   `json:"userInn"`              // ИНН
	DateTime      string   `json:"dateTime"`             // Дата и время (в строке)
	TotalSum      int64    `json:"totalSum"`             // Итоговая сумма в копейках
	RetailPlace   string   `json:"retailPlace"`          // Место расчета
	KktRegId      string   `json:"kktRegId"`             // Регистрационный номер ККТ
	FiscalDriveNo string   `json:"fiscalDriveNumber"`    // ФН
	FiscalDocNo   int      `json:"fiscalDocumentNumber"` // ФД
	FiscalSign    int64    `json:"fiscalSign"`           // ФП
	Items         []Item   `json:"items"`                // Список товаров
	Metadata      Metadata `json:"metadata"`
}

// Товарная позиция
type Item struct {
	Name            string       `json:"name"`
	Price           int64        `json:"price"`
	Quantity        float64      `json:"quantity"`
	Sum             int64        `json:"sum"`
	Nds             int          `json:"nds"`
	ProductType     int          `json:"productType"`
	ProductCodeNew  *ProductCode `json:"productCodeNew,omitempty"`
	QuantityMeasure int          `json:"itemsQuantityMeasure"` // Единицы измерения
}

// Коды маркировки (EAN-13, Честный Знак)
type ProductCode struct {
	Ean13 *Ean13 `json:"ean13,omitempty"`
	Gs1m  *Gs1m  `json:"gs1m,omitempty"`
}

type Ean13 struct {
	Gtin string `json:"gtin"`
}

type Gs1m struct {
	Gtin   string `json:"gtin"`
	Sernum string `json:"sernum"`
}

type Metadata struct {
	Address     string    `json:"address"`
	OfdId       string    `json:"ofdId"`
	ReceiveDate time.Time `json:"receiveDate"`
}
