package uds

import (
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"strconv"
	"time"
)

// OperationList
// Список операций
type OperationList struct {
	Rows   []Operation `json:"rows"`   // Информация об операциях.
	Cursor string      `json:"cursor"` // Используйте это значение для следующего запроса
}

// ActionState Статус операции.
type ActionState string

const (
	ActionStateNormal   ActionState = "NORMAL"
	ActionStateCanceled ActionState = "CANCELED"
	ActionStateReversal ActionState = "REVERSAL"
)

// Operation
// Информация об операции.
type Operation struct {
	Id            int64             `json:"id"`            // D операции в базе UDS.
	DateCreated   time.Time         `json:"dateCreated"`   // Дата операции.
	Action        string            `json:"action"`        // Тип операции.
	State         ActionState       `json:"state"`         // Статус операции.
	Customer      CustomerShortInfo `json:"customer"`      // Информация о клиенте.
	Cashier       Cashier           `json:"cashier"`       // Информация о сотруднике.
	Branch        BranchInfo        `json:"branch"`        // Информация о филиале.
	Points        float64           `json:"points"`        // Количество бонусных баллов, которое будет списано с клиента после завершения операции. Отрицательное значение говорит о списании, а положительное - о начислении бонусных баллов.
	ReceiptNumber string            `json:"receiptNumber"` // Номер чека.
	Origin        Origin            `json:"origin"`        // Для сторнирующей операции - ссылка на оригинальную операцию.
	Total         float64           `json:"total"`         // Общая сумма чека до применения скидок в денежных единицах.
	Cash          float64           `json:"cash"`          // Оплачиваемая сумма в денежных единицах.
}

// Cashier
// Информация о сотруднике.
type Cashier struct {
	Id          int64  `json:"id"`
	DisplayName string `json:"displayName"`
}

// Branch
// Информация о филиале.
type BranchInfo struct {
	Id          int64  `json:"id"`
	DisplayName string `json:"displayName"`
}

// Origin Для сторнирующей операции - ссылка на оригинальную операцию.
type Origin struct {
	Id int64 `json:"id"` // Идентификатор исходной (оригинальной) операции.
}

type ParticipantShort struct {
	Uid   *string `json:"uid,omitempty"`   // Идентификатор клиента в UDS (UID).
	Phone *string `json:"phone,omitempty"` // Номер телефона.
}

func (p *ParticipantShort) SetUid(uid string) *ParticipantShort {
	p.Uid = &uid
	return p
}

func (p *ParticipantShort) SetPhone(phone string) *ParticipantShort {
	p.Phone = &phone
	return p
}

// CashierExternal
// Информация о сотруднике.
type CashierExternal struct {
	ExternalId string  `json:"externalId"`     // Внешний идентификатор сотрудника.
	Name       *string `json:"name,omitempty"` // Имя сотрудника.
}

// OperationGetList
// Получить Список операций
// https://docs.uds.app/#tag/Operations/paths/~1operations/get
func (u *Client) OperationGetList(maxValue int, cursor string) (*OperationList, *resty.Response, error) {
	operationList := new(OperationList)
	apiErr := new(ApiError)

	req := u.client.R()

	if maxValue > 0 {
		maxValue = max(1, min(50, maxValue)) // от 1 до 50
		maxString := strconv.Itoa(maxValue)
		req.SetQueryParam("max", maxString)
	}

	if cursor != "" {
		req.SetQueryParam("cursor", cursor)
	}

	resp, err := req.
		SetResult(operationList).
		SetError(apiErr).
		Get("operations")

	if err != nil {
		return nil, resp, err
	}

	if resp.Error() != nil {
		return nil, resp, apiErr
	}

	return operationList, resp, nil
}

// CreateOperationRequest
// Объект запроса на создание операции
type CreateOperationRequest struct {
	Code        *string           `json:"code,omitempty"`        // Код на оплату.
	Participant *ParticipantShort `json:"participant,omitempty"` // Информация о клиенте.
	Nonce       string            `json:"nonce,omitempty"`       // Уникальный идентификатор операции (UUID).
	Cashier     *CashierExternal  `json:"cashier,omitempty"`     // Информация о сотруднике.
	Receipt     Receipt           `json:"receipt"`               // Информация о чеке.
	Tags        []int64           `json:"tags,omitempty"`        // Список id тегов компании, назначаемых клиенту при проведении операции. Передача null означает, что существующее значение не будет изменено
}

func (o *CreateOperationRequest) SetCode(code string) *CreateOperationRequest {
	o.Code = &code
	return o
}

// Receipt
// Информация о чеке.
type Receipt struct {
	Total             float64  `json:"total"`                       // Сумма счета в денежных единицах.
	Cash              float64  `json:"cash"`                        // Оплачиваемая сумма в денежных единицах.
	Points            float64  `json:"points"`                      // Оплачиваемая сумма в бонусных баллах.
	Number            *string  `json:"number,omitempty"`            // Номер чека.
	SkipLoyaltyTotal  *float64 `json:"skipLoyaltyTotal,omitempty"`  // Часть суммы счета, на которую не начисляется кешбэк и на которую не распространяется скидка (в денежных единицах).
	UnredeemableTotal *float64 `json:"unredeemableTotal,omitempty"` // Часть суммы счета, которую нельзя погасить баллами.
}

// CreateOperationResponse
// Объект ответа на создание операции
type CreateOperationResponse struct {
	Id            int64             `json:"id"`            // ID операции в базе UDS.
	DateCreated   time.Time         `json:"dateCreated"`   // Дата операции.
	Action        string            `json:"action"`        // Тип операции.
	State         ActionState       `json:"state"`         // Статус операции.
	Customer      CustomerShortInfo `json:"customer"`      // Информация о клиенте.
	Cashier       Cashier           `json:"cashier"`       // Информация о сотруднике.
	Branch        BranchInfo        `json:"branch"`        // Информация о филиале.
	Points        float64           `json:"points"`        // Количество бонусных баллов, которое будет списано с клиента после завершения операции. Отрицательное значение говорит о списании, а положительное - о начислении бонусных баллов.
	ReceiptNumber string            `json:"receiptNumber"` // Номер чека.
	Origin        Origin            `json:"origin"`        // Для сторнирующей операции - ссылка на оригинальную операцию.
	Total         float64           `json:"total"`         // Общая сумма чека до применения скидок в денежных единицах.
	Cash          float64           `json:"cash"`          // Оплачиваемая сумма в денежных единицах.
}

// OperationCreate
// Проведение операции
// https://docs.uds.app/#tag/Operations/paths/~1operations/post
func (u *Client) OperationCreate(operation *CreateOperationRequest) (*CreateOperationResponse, *resty.Response, error) {
	createResp := new(CreateOperationResponse)

	if operation.Nonce == "" {
		operation.Nonce = uuid.New().String()
	}

	resp, err := u.client.R().
		SetBody(operation).
		SetResult(createResp).
		Post("operations")

	if err != nil {
		return nil, resp, err
	}

	return createResp, resp, nil
}

// OperationGetByID
// Получение информации об операции
// https://docs.uds.app/#tag/Operations/paths/~1operations~1{id}/get
func (u *Client) OperationGetByID(id int64) (*Operation, *resty.Response, error) {
	operation := new(Operation)

	idString := strconv.FormatInt(id, 10)

	resp, err := u.client.R().SetPathParam("id", idString).
		SetResult(operation).
		Get("operations/{id}")

	if err != nil {
		return nil, resp, err
	}

	return operation, resp, nil
}

// RefundOperationRequest
// Объект запроса на возврат по операции
type RefundOperationRequest struct {
	PartialAmount float64 `json:"partialAmount"` // Сумма возврата.
}

// OperationRefund
// Операция возврата
// https://docs.uds.app/#tag/Operations/paths/~1operations~1{id}~1refund/post
func (u *Client) OperationRefund(id int64, partialAmount float64) (*Operation, *resty.Response, error) {
	operation := new(Operation)

	idString := strconv.FormatInt(id, 10)
	refund := RefundOperationRequest{partialAmount}

	resp, err := u.client.R().SetPathParam("id", idString).
		SetBody(refund).
		SetResult(operation).
		Post("operations/{id}/refund")

	if err != nil {
		return nil, resp, err
	}

	return operation, resp, nil
}

// CalcOperationRequest
// Объект запроса на расчёт информации по операции
type CalcOperationRequest struct {
	Code        *string              `json:"code,omitempty"`        // Код на оплату.
	Participant *ParticipantShort    `json:"participant,omitempty"` // Информация о клиенте.
	Receipt     CalcOperationReceipt `json:"receipt"`
}

func (o *CalcOperationRequest) SetCode(code string) *CalcOperationRequest {
	o.Code = &code
	return o
}

type CalcOperationReceipt struct {
	Total             float64  `json:"total"`                       // Общая сумма чека до применения скидок в денежных единицах.
	SkipLoyaltyTotal  *float64 `json:"skipLoyaltyTotal,omitempty"`  // Часть суммы счета, на которую не начисляется кешбэк и на которую не распространяется скидка (в денежных единицах).
	UnredeemableTotal *float64 `json:"unredeemableTotal,omitempty"` // Часть суммы счета, которую нельзя погасить баллами.
	Points            *float64 `json:"points,omitempty"`            // Количество бонусных баллов, которое клиент хочет списать. По умолчанию максимально доступное число баллов.
}

func (r *CalcOperationReceipt) SetSkipLoyaltyTotal(value float64) *CalcOperationReceipt {
	r.SkipLoyaltyTotal = &value
	return r
}

func (r *CalcOperationReceipt) SetUnredeemableTotal(value float64) *CalcOperationReceipt {
	r.UnredeemableTotal = &value
	return r
}

func (r *CalcOperationReceipt) SetPoints(value float64) *CalcOperationReceipt {
	r.Points = &value
	return r
}

// CalcOperationResponse
// Объект ответа на запрос по расчёту информации по операции
type CalcOperationResponse struct {
	User     CustomerShortInfo `json:"user"`     // Информация о клиенте.
	Purchase PurchaseDetail    `json:"purchase"` // Информация об операции.
}

// OperationCalc
// Рассчитать информацию по операции
// https://docs.uds.app/#tag/Operations/paths/~1operations~1calc/post
func (u *Client) OperationCalc(operation *CalcOperationRequest) (*CalcOperationResponse, *resty.Response, error) {
	calcResp := new(CalcOperationResponse)

	resp, err := u.client.R().
		SetBody(operation).
		SetResult(calcResp).
		Post("operations/calc")

	if err != nil {
		return nil, resp, err
	}

	return calcResp, resp, nil
}

// RewardOperationRequest
// Объект запроса на начисление бонусов клиенту
type RewardOperationRequest struct {
	Points       float64 `json:"points"`            // Количество бонусных баллов. (Можем иметь отрицательное значение - списание)
	Comment      string  `json:"comment,omitempty"` // Текст комментария, который увидит пользователь.
	Participants []int64 `json:"participants"`      // Список ID клиентов в компании.
	Silent       bool    `json:"silent"`            // Не отправлять пуш-уведомление клиенту (default false).
}

// RewardOperationResponse
// Объект ответа на запрос на начисление бонусов клиенту
type RewardOperationResponse struct {
	Accepted int `json:"accepted"` // Количество пользователей, которым будут начислены бонусные баллы.
}

// OperationReward
// Начисление бонусов клиенту (подарок)
// https://docs.uds.app/#tag/Operations/paths/~1operations~1reward/post
func (u *Client) OperationReward(operation RewardOperationRequest) (*RewardOperationResponse, *resty.Response, error) {
	rewardResp := new(RewardOperationResponse)
	apiErr := new(ApiError)

	resp, err := u.client.R().
		SetBody(operation).
		SetResult(rewardResp).
		SetError(apiErr).
		Post("operations/reward")

	if err != nil {
		return nil, resp, err
	}

	if resp.Error() != nil {
		return nil, resp, apiErr
	}

	return rewardResp, resp, nil
}
