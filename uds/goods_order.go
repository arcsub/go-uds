package uds

import (
	"github.com/go-resty/resty/v2"
	"strconv"
	"time"
)

// GoodsOrderState
// Статус заказа.
type GoodsOrderState string

const (
	GoodsOrderStateNew            GoodsOrderState = "NEW"
	GoodsOrderStateCompleted      GoodsOrderState = "COMPLETED"
	GoodsOrderStateDeleted        GoodsOrderState = "DELETED"
	GoodsOrderStateWaitingPayment GoodsOrderState = "WAITING_PAYMENT"
)

// GoodsOrderDetailed
// Структура заказа
type GoodsOrderDetailed struct {
	Id                int               `json:"id"`                // ID заказа.
	DateCreated       time.Time         `json:"dateCreated"`       // Дата заказа.
	Comment           string            `json:"comment"`           // Комментарий к заказу.
	State             GoodsOrderState   `json:"state"`             // Статус заказа.
	Cash              float64           `json:"cash"`              // Сумма, оплачиваемая деньгами.
	Points            float64           `json:"points"`            // Количество списываемых баллов.
	Total             float64           `json:"total"`             // Сумма заказа.
	CertificatePoints float64           `json:"certificatePoints"` // Количество списываемых баллов сертификата.
	Customer          CustomerShortInfo `json:"customer"`          // Информация о клиенте.
	Delivery          Delivery          `json:"delivery"`          // Способ получения заказа.
	OnlinePayment     OnlinePayment     `json:"onlinePayment"`     // Информация об онлайн-оплате.
	PaymentMethod     PaymentMethod     `json:"paymentMethod"`     // Информация об оплате.
	Items             []GoodOrderItem   `json:"items"`             // Информация о товарах.
	Purchase          PurchaseDetail    `json:"purchase"`          // Информация об операции.
}

// DeliveryType
// Способ получения заказа.
type DeliveryType string

const (
	DeliveryTypePickUp   DeliveryType = "PICKUP"
	DeliveryTypeDelivery DeliveryType = "DELIVERY"
)

// Delivery
// Способ получения заказа.
type Delivery struct {
	ReceiverName  string       `json:"receiverName"`  // Имя клиента, который заберет заказа.
	ReceiverPhone string       `json:"receiverPhone"` // Номер телефона клиента, который заберет заказ.
	UserComment   string       `json:"userComment"`   // Комментарий клиента к заказу
	Branch        BranchInfo   `json:"branch"`        // Информация о филиале.
	Type          DeliveryType `json:"type"`          // Способ получения заказа.
}

// PaymentProviderType
// Тип платежной системы
type PaymentProviderType string

const (
	PaymentProviderB2P           PaymentProviderType = "B2P"
	PaymentProviderCloudPayments PaymentProviderType = "CLOUD_PAYMENTS"
	PaymentProviderCustom        PaymentProviderType = "CUSTOM"
)

// OnlinePayment
// Информация об онлайн-оплате.
type OnlinePayment struct {
	PaymentProvider PaymentProviderType `json:"paymentProvider"` // Тип платежной системы
	Id              string              `json:"id"`              // Идентификатор платежа во внешней платежной системе
	Completed       bool                `json:"completed"`       // Статус оплаты.
}

// PaymentMethodType
// Тип оплаты.
type PaymentMethodType string

const (
	PaymentMethodBestToPay     PaymentMethodType = "BEST_TO_PAY"
	PaymentMethodCloudPayments PaymentMethodType = "CLOUD_PAYMENTS"
	PaymentMethodCash          PaymentMethodType = "CASH"
	PaymentMethodManual        PaymentMethodType = "MANUAL"
	PaymentMethodCustom        PaymentMethodType = "CUSTOM"
)

// PaymentMethod
// Информация об оплате.
type PaymentMethod struct {
	Type PaymentMethodType `json:"type"` // Тип оплаты.
	Name string            `json:"name"` // Название пользовательского метода оплаты с типом MANUAL
}

// GoodsItemType
// Тип товара.
type GoodsItemType string

const (
	GoodsItemTypeItem        GoodsItemType = "ITEM"
	GoodsItemTypeVaryingItem GoodsItemType = "VARYING_ITEM"
)

// GoodsMeasurement
// Единицы измерения товаров.
type GoodsMeasurement string

const (
	GoodsMeasurementPiece      GoodsMeasurement = "PIECE"
	GoodsMeasurementCentimeter GoodsMeasurement = "CENTIMETRE"
	GoodsMeasurementMetre      GoodsMeasurement = "METRE"
	GoodsMeasurementMillilitre GoodsMeasurement = "MILLILITRE"
	GoodsMeasurementLitre      GoodsMeasurement = "LITRE"
	GoodsMeasurementGram       GoodsMeasurement = "GRAM"
	GoodsMeasurementKiloGram   GoodsMeasurement = "KILOGRAM"
)

// GoodOrderItem
// Информация о товаре в заказе.
type GoodOrderItem struct {
	Id          int              `json:"id"`          // ID товара в UDS.
	ExternalId  string           `json:"externalId"`  // Внешний идентификатор товара.
	Name        string           `json:"name"`        // Название товара
	VariantName string           `json:"variantName"` // Имя варианта товара, если тип этого товара VARYING_ITEM
	Sku         string           `json:"sku"`         // Артикул товара.
	Type        GoodsItemType    `json:"type"`        // Тип товара.
	Qty         float64          `json:"qty"`         // Количество.
	Price       float64          `json:"price"`       // Цена товара.
	Measurement GoodsMeasurement `json:"measurement"` // Единицы измерения товаров.
}

// GoodsOrderGetByID
// Подробная информация о заказе
// https://docs.uds.app/#tag/Goods-Order/paths/~1goods-orders~1{id}/get
func (u *Client) GoodsOrderGetByID(id int64) (*GoodsOrderDetailed, *resty.Response, error) {
	goodsOrder := new(GoodsOrderDetailed)

	idString := strconv.FormatInt(id, 10)

	resp, err := u.client.R().SetPathParam("id", idString).
		SetResult(goodsOrder).
		Get("goods-orders/{id}")

	if err != nil {
		return nil, resp, err
	}

	return goodsOrder, resp, nil
}

// DeliveryCase
// Информация о доставке.
type DeliveryCase struct {
	Name  string  `json:"name"`  // Название доставки.
	Value float64 `json:"value"` // Стоимость доставки.
}

// GoodsOrderItemUpdate
// Объект обновления товара
type GoodsOrderItemUpdate struct {
	Id          int    `json:"id"`          // ID товара в UDS.
	VariantName string `json:"variantName"` // Имя варианта товара, если тип этого товара VARYING_ITEM
	Qty         int    `json:"qty"`         // Количество.
}

// GoodsOrderItemNew
// Объект добавления нового товара в заказ
type GoodsOrderItemNew struct {
	ExternalId  string  `json:"externalId"`  // Внешний идентификатор товара.
	Name        string  `json:"name"`        // Название товара
	VariantName string  `json:"variantName"` // Имя варианта товара, если тип этого товара VARYING_ITEM
	QTY         int     `json:"qty"`         // Количество.
	Price       float64 `json:"price"`       // Цена товара.
	SkipLoyalty bool    `json:"skipLoyalty"` // Не применять бонусную программу к товару.
}

// UpdateGoodsOrderRequest
// Объект запроса на изменение заказа
type UpdateGoodsOrderRequest[T GoodsOrderItemUpdate | GoodsOrderItemNew] struct {
	DeliveryCase DeliveryCase `json:"deliveryCase"` // Информация о доставке.
	Items        []T          `json:"items"`        // Информация о товарах.
}

func updateGoodsOrderItemsProcess(req *resty.Request, id int64) (*GoodsOrderDetailed, *resty.Response, error) {
	goodsOrder := new(GoodsOrderDetailed)
	idString := strconv.FormatInt(id, 10)
	resp, err := req.SetPathParam("id", idString).SetResult(goodsOrder).Put("goods-orders/{id}")

	if err != nil {
		return nil, resp, err
	}

	return goodsOrder, resp, nil
}

// GoodsOrderUpdateItems
// Изменить товары заказа
// https://docs.uds.app/#tag/Goods-Order/paths/~1goods-orders~1{id}/put
func (u *Client) GoodsOrderUpdateItems(id int64, updatedOrder *UpdateGoodsOrderRequest[GoodsOrderItemUpdate]) (*GoodsOrderDetailed, *resty.Response, error) {
	req := u.client.R().SetBody(updatedOrder)
	return updateGoodsOrderItemsProcess(req, id)
}

// GoodsOrderAddItems
// Добавить товары в заказ
// https://docs.uds.app/#tag/Goods-Order/paths/~1goods-orders~1{id}/put
func (u *Client) GoodsOrderAddItems(id int64, updatedOrder *UpdateGoodsOrderRequest[GoodsOrderItemNew]) (*GoodsOrderDetailed, *resty.Response, error) {
	req := u.client.R().SetBody(updatedOrder)
	return updateGoodsOrderItemsProcess(req, id)
}

// CompleteGoodsOrder
// Объект ответа на запрос о завершении заказа
type CompleteGoodsOrder struct {
	Transaction struct {
		Id int64 `json:"id"` // ID операции.
	} `json:"transaction"`
	Order GoodsOrderDetailed `json:"order"` // Информация о заказе.
}

// GoodsOrderComplete
// Завершает заказ товара с идентификатором и создает транзакцию
// https://docs.uds.app/#tag/Goods-Order/paths/~1goods-orders~1{id}~1complete/post
func (u *Client) GoodsOrderComplete(id int64) (*CompleteGoodsOrder, *resty.Response, error) {
	completeGoodsOrder := new(CompleteGoodsOrder)

	idString := strconv.FormatInt(id, 10)

	resp, err := u.client.R().SetPathParam("id", idString).
		SetResult(completeGoodsOrder).
		Post("goods-orders/{id}/complete")

	if err != nil {
		return nil, resp, err
	}

	return completeGoodsOrder, resp, nil
}

// GoodsOrderGenerateCode
// Сгенерировать код для завершения заказа товара с идентификатором
// https://docs.uds.app/#tag/Goods-Order/paths/~1goods-orders~1{id}~1code/post
func (u *Client) GoodsOrderGenerateCode(id int64) (string, *resty.Response, error) {
	type s struct {
		Code string `json:"code"`
	}
	code := new(s)

	idString := strconv.FormatInt(id, 10)

	resp, err := u.client.R().SetPathParam("id", idString).
		SetResult(code).
		Post("goods-orders/{id}/code")

	if err != nil {
		return "", resp, err
	}

	return code.Code, resp, nil
}
