package uds

import (
	"github.com/go-resty/resty/v2"
	"strconv"
	"time"
)

type Gender string

const (
	GenderMale         Gender = "MALE"          // Мужчина
	GenderFemale       Gender = "FEMALE"        // Женщина
	GenderNotSpecified Gender = "NOT_SPECIFIED" // Пол не установлен
)

// FindCustomerResponse
// Объект ответа на запрос поиска клиента
type FindCustomerResponse struct {
	User     CustomerDetail `json:"user"`     // Информация о клиенте.
	Code     string         `json:"code"`     // Новый долгоживущий код на оплату, если был запрошен параметр exchangeCode.
	Purchase PurchaseDetail `json:"purchase"` // Информация об операции.
}

// Customer
// Объект с данными клиента
type Customer struct {
	Uid         string      `json:"uid"`         // Идентификатор клиента в UDS (UID).
	Avatar      string      `json:"avatar"`      // URL изображения клиента.
	DisplayName string      `json:"displayName"` // Имя клиента.
	Gender      Gender      `json:"gender"`      // Пол.
	Phone       string      `json:"phone"`       // Номер телефона клиента.
	BirthDate   string      `json:"birthDate"`   // Дата рождения клиента.
	Participant Participant `json:"participant"` // Информация о клиенте.
	ChannelName string      `json:"channelName"` // Источник трафика.
	Email       string      `json:"email"`       // Email клиента.
}

// CustomerDetail
// Детальный объект с данными клиента
type CustomerDetail struct {
	Customer
	Tags []TagModel `json:"tags"` // Список тегов клиента.
}

// Participant
// Информация о клиенте.
type Participant struct {
	Id                  int64          `json:"id"`                  // ID клиента в компании.
	InviterId           int            `json:"inviterId"`           // ID клиента в компании, пригласившего данного клиента.
	Points              float64        `json:"points"`              // Баланс бонусных баллов клиента.
	DiscountRate        float64        `json:"discountRate"`        // Размер скидки (в процентах).
	CashbackRate        float64        `json:"cashbackRate"`        // Размер кешбэка (в процентах) для данного клиента UDS.
	MembershipTier      MembershipTier `json:"membershipTier"`      // Настройки статусов клиентов.
	DateCreated         time.Time      `json:"dateCreated"`         // Дата, когда клиент присоединился к компании
	LastTransactionTime time.Time      `json:"lastTransactionTime"` // Дата и время, когда клиент совершил последнюю транзакцию.
}

// CustomerShortInfo
// Минимальная информация о клиенте.
type CustomerShortInfo struct {
	Id             int64          `json:"id"`             // ID клиента в компании
	DisplayName    string         `json:"displayName"`    // Имя и фамилия клиента.
	Uid            string         `json:"uid"`            // Идентификатор клиента в UDS (UID).
	MembershipTier MembershipTier `json:"membershipTier"` // Настройки статусов клиентов.
}

// CustomerGetList
// Получить список клиентов
// https://docs.uds.app/#tag/Customers/paths/~1customers/get
func (u *Client) CustomerGetList(maxValue int, offset int) (*List[Customer], *resty.Response, error) {
	customers := new(List[Customer])

	req := u.client.R()

	if maxValue > 0 {
		maxValue = max(1, min(50, maxValue)) // от 1 до 50
		maxString := strconv.Itoa(maxValue)
		req.SetQueryParam("max", maxString)
	}

	if offset > 0 {
		offset = max(1, min(10000, offset)) // от 1 до 10000
		offsetString := strconv.Itoa(offset)
		req.SetQueryParam("offset", offsetString)
	}

	resp, err := u.client.R().
		SetResult(customers).
		Get("customers")

	if err != nil {
		return nil, resp, err
	}

	return customers, resp, nil
}

type FindCustomerParams struct {
	ExchangeCode      bool    // Если указан, то в ответе будет отправлен новый долгоживущий код на оплату
	Total             float64 // Общая сумма счета в денежных единицах
	SkipLoyaltyTotal  float64 // Часть суммы счета, на которую не начисляется кешбэк и на которую не распространяется скидка (в денежных единицах).
	UnredeemableTotal float64 // Часть суммы счета, которую нельзя погасить баллами.
}

func float64ToString(value float64) string {
	return strconv.FormatFloat(value, 'g', -1, 64)
}

// findCustomerProcess осуществляет запрос на поиск клиента по нужному ключу
func findCustomerProcess(req *resty.Request, params *FindCustomerParams) (*FindCustomerResponse, *resty.Response, error) {
	customer := new(FindCustomerResponse)
	apiErr := new(ApiError)

	if params != nil {
		if params.ExchangeCode {
			req.SetQueryParam("exchangeCode", "true")
		}

		if params.Total > 0 {
			strVal := float64ToString(params.Total)
			req.SetQueryParam("total", strVal)
		}

		if params.SkipLoyaltyTotal > 0 {
			strVal := float64ToString(params.SkipLoyaltyTotal)
			req.SetQueryParam("skipLoyaltyTotal", strVal)
		}

		if params.SkipLoyaltyTotal > 0 {
			strVal := float64ToString(params.SkipLoyaltyTotal)
			req.SetQueryParam("skipLoyaltyTotal", strVal)
		}
	}

	resp, err := req.
		SetResult(customer).
		SetError(apiErr).
		Get("customers/find")

	if err != nil {
		return nil, resp, err
	}

	if resp.Error() != nil {
		return nil, resp, apiErr
	}

	return customer, resp, nil
}

// CustomerFindByCode
// Поиск клиента по коду из приложения
// params может быть nil
// https://docs.uds.app/#tag/Customers/paths/~1customers~1find/get
func (u *Client) CustomerFindByCode(code string, params *FindCustomerParams) (*FindCustomerResponse, *resty.Response, error) {
	req := u.client.R().SetQueryParam("code", code)
	return findCustomerProcess(req, params)
}

// CustomerFindByPhone
// Поиск клиента по номеру телефона в формате +79998887766
// params может быть nil
// https://docs.uds.app/#tag/Customers/paths/~1customers~1find/get
func (u *Client) CustomerFindByPhone(phone string, params *FindCustomerParams) (*FindCustomerResponse, *resty.Response, error) {
	req := u.client.R().SetQueryParam("phone", phone)
	return findCustomerProcess(req, params)
}

// CustomerFindByUID
// Поиск клиента по uid
// params может быть nil
// https://docs.uds.app/#tag/Customers/paths/~1customers~1find/get
func (u *Client) CustomerFindByUID(uid string, params *FindCustomerParams) (*FindCustomerResponse, *resty.Response, error) {
	req := u.client.R().SetQueryParam("uid", uid)
	return findCustomerProcess(req, params)
}

// CustomerGetByID
// Получение информации о клиенте по ID
// https://docs.uds.app/#tag/Customers/paths/~1customers~1{id}/get
func (u *Client) CustomerGetByID(id int64) (*CustomerDetail, *resty.Response, error) {
	customer := new(CustomerDetail)
	apiErr := new(ApiError)

	idString := strconv.FormatInt(id, 10)

	resp, err := u.client.R().
		SetPathParam("id", idString).
		SetResult(customer).
		SetError(apiErr).
		Get("customers/{id}")

	if err != nil {
		return nil, resp, err
	}

	if resp.Error() != nil {
		return nil, resp, apiErr
	}

	return customer, resp, nil
}

// CustomerTagList
// Список тегов клиента
type CustomerTagList struct {
	Rows  []TagModel `json:"rows"`  // Список тегов.
	Total int        `json:"total"` // Количество тегов.
}

// TagModel
// Тег клиента
type TagModel struct {
	Id   int64  `json:"id"`   // Идентификатор тега.
	Name string `json:"name"` // Наименование тега.
}

// CustomerGetTags
// Получение списка тегов клиента
// https://docs.uds.app/#tag/Customers/paths/~1customers~1{id}~1tags/get
func (u *Client) CustomerGetTags(id int64) (*CustomerTagList, *resty.Response, error) {
	tags := new(CustomerTagList)
	apiErr := new(ApiError)

	idString := strconv.FormatInt(id, 10)

	resp, err := u.client.R().
		SetPathParam("id", idString).
		SetResult(tags).
		SetError(apiErr).
		Get("customers/{id}/tags")

	if err != nil {
		return nil, resp, err
	}

	if resp.Error() != nil {
		return nil, resp, apiErr
	}

	return tags, resp, nil
}

// SetCustomerTagsRequest
// Объект запроса на установку тегов клиенту
type SetCustomerTagsRequest struct {
	IDs []int64 `json:"ids"`
}

// CustomerSetTags
// Установка тегов клиенту
// https://docs.uds.app/#tag/Customers/paths/~1customers~1{id}~1tags/get
func (u *Client) CustomerSetTags(id int64, tagsReq SetCustomerTagsRequest) (*List[TagModel], *resty.Response, error) {
	tags := new(List[TagModel])
	apiErr := new(ApiError)

	idString := strconv.FormatInt(id, 10)

	resp, err := u.client.R().
		SetPathParam("id", idString).
		SetBody(tagsReq).
		SetResult(tags).
		SetError(apiErr).
		Post("customers/{id}/tags")

	if err != nil {
		return nil, resp, err
	}

	if resp.Error() != nil {
		return nil, resp, apiErr
	}

	return tags, resp, nil
}
