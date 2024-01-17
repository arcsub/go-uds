package uds

import (
	"fmt"
)

type ErrorCode string

const (
	// ErrNotFound
	// Пользователь с данным кодом на оплату или идентификатором не найден.
	// В основном, возникает при неверно введенном коде клиента или время жизни кода истекло.
	// Также возникает при запросе find по номеру телефона,
	// если такой номер в системе не зарегистрирован и по нему в компании ранее не было операций.
	ErrNotFound ErrorCode = "notFound"

	// ErrBadRequest
	// Возникли ошибки при проведении валидации.
	// Для получения подробной информации об ошибке обратитесь к полю errors.
	// Возникает при передаче некорректного запроса.
	ErrBadRequest ErrorCode = "badRequest"

	// ErrForbidden
	// Доступ запрещен, данный токен аутентификации не установлен или не имеет соответствующего разрешения.
	ErrForbidden ErrorCode = "forbidden"

	// ErrInvalidChecksum
	// Указанные значения в полях cash, points или total не соответствуют настройке.
	// Способ предоставления скидки в настройках программы лояльности в UDS Бизнес.
	// Должно быть в случае начисления бонусных баллов total= cash + points.
	// Например, total =100, cash = 50, points = 50.
	// Для предоставления скидки total = (cash + point) / (1 - процент скидки).
	// Например, total =100, cash = 50, points = 40 (для скидки в 10%).
	// Ошибки часто возникают при применении сторонних скидок, неверном округлении.
	// В случае печати ваучера ошибка может возникать из-за способа предоставления скидки не CHARGE_SCORES;
	// размер начисляемых бонусных баллов равен 0.0; поле total меньше или равноskipLoyaltyTotal.
	ErrInvalidChecksum ErrorCode = "invalidChecksum"

	// ErrInsufficientFunds
	// Значение в поле points превышает доступное количество бонусных баллов на счете клиента.
	// Ошибка в основном возникает при неверном округлении, округлять баллы можно только в меньшую сторону.
	ErrInsufficientFunds ErrorCode = "insufficientFunds"

	// ErrDiscountLimitExceed
	// Соотношение points / total больше, чем указано в настройках UDS.
	// Компания может установить максимальный процент от чека, на который допускается списать бонусы.
	// Посмотреть допустимый процент списания баллов можно в настройках UDS Бизнес в Программа лояльности
	// в графе Какой процент счета можно оплатить баллами.
	ErrDiscountLimitExceed ErrorCode = "discountLimitExceed"

	// ErrUnauthorized
	// Неверно указан ID компании или API Key.
	// Необходимо проверить актуальность API Key и ID компании
	// на странице Интеграция и корректность аутентификации Basic.
	ErrUnauthorized ErrorCode = "unauthorized"

	// ErrWithdrawNotPermitted
	// В запросе было указано значение в поле participant -> uid, при этом значение в поле points не равно 0.0.
	// Ошибка может возникать, если передавать вместе с параметром code в запросе на проведение
	// операции со списанием баллов параметров uid или phone.
	// Необходимо удалить параметры uid и phone в запросе и оставить только code.
	// Ошибка может также возникать при попытке списать бонусы по uid или номеру телефона
	// (если в компании такая возможность отсутствует).
	// Необходимо внести корректировку в сумму баллов и указать 0.0.
	ErrWithdrawNotPermitted ErrorCode = "withdrawNotPermitted"

	// ErrPurchaseByPhoneDisabled
	// Проведение операции по номеру телефона не разрешено настройками UDS.
	// Необходимо включить настройку оплаты по номеру телефона в настройках UDS Бизнес.
	ErrPurchaseByPhoneDisabled ErrorCode = "purchaseByPhoneDisabled"

	// ErrGoodsNodeIndexInvalid
	// Для создания категории указание идентификатора nodeId недопустимо.
	// Необходимо проверить актуальный ID категории, в которую создается товар или категория.
	ErrGoodsNodeIndexInvalid ErrorCode = "goods.nodeIndex.invalid"

	// ErrGoodsLimitIsReached
	// Превышен лимит количества товаров.
	// Лимит товаров можно проверить в разделе Товары и услуги в UDS Бизнес.
	ErrGoodsLimitIsReached ErrorCode = "goods.limitIsReached"

	// ErrParticipantIsBlocked
	// Клиент в данной компании заблокирован.
	ErrParticipantIsBlocked ErrorCode = "participantIsBlocked"
)

type ApiError struct {
	ErrorCode ErrorCode         `json:"errorCode"` // Код ошибки.
	Message   string            `json:"message"`   // Описание ошибки.
	Errors    []BadRequestError `json:"errors"`    // Присутствует, если ErrorCode = ErrBadRequest
}

type BadRequestError struct {
	ErrorCode string `json:"errorCode"` // Код ошибки.
	Message   string `json:"message"`   // Описание ошибки.
	Field     string `json:"field"`     // Поле объекта, на которое указывает ошибка
	Value     any    `json:"value"`     // Значение поля
}

func (e ApiError) Error() string {
	var message string

	if e.ErrorCode == ErrBadRequest {
		message = fmt.Sprintf("[%s]: %s", e.ErrorCode, e.Message)
		for _, requestError := range e.Errors {
			message += fmt.Sprintf("\n[%s]: %s", requestError.ErrorCode, requestError.Message)
			if requestError.Field != "" {
				message += fmt.Sprintf("; field '%s', value: %v", requestError.Field, requestError.Value)
			}
		}
	} else {
		message = fmt.Sprintf("[%s]: %s", e.ErrorCode, e.Message)
	}

	return message
}
