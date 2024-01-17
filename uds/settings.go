package uds

import "github.com/go-resty/resty/v2"

type DiscountPolicy string

const (
	DiscountPolicyApplyDiscount DiscountPolicy = "APPLY_DISCOUNT" // понижать сумму счета (скидка)
	DiscountPolicyChargeScores  DiscountPolicy = "CHARGE_SCORES"  // начислять бонусные баллы (кешбэк)
)

// Settings
// Настройки компании.
type Settings struct {
	Id                     int                    `json:"id,omitempty"`                     // Идентификатор компании в UDS.
	Name                   string                 `json:"name,omitempty"`                   // Название компании.
	PromoCode              string                 `json:"promoCode,omitempty"`              // Промокод компании для вступления.
	Currency               string                 `json:"currency,omitempty"`               // Валюта компании (стандарт ISO-4217).
	BaseDiscountPolicy     DiscountPolicy         `json:"baseDiscountPolicy,omitempty"`     // Определяет тип программы лояльности
	LoyaltyProgramSettings LoyaltyProgramSettings `json:"loyaltyProgramSettings,omitempty"` // Настройки бонусной программы компании.
	PurchaseByPhone        bool                   `json:"purchaseByPhone,omitempty"`        // Возможность проведения операции, используя номер телефона клиента.
	WriteInvoice           bool                   `json:"writeInvoice,omitempty"`           // Необходимо ли указывать номер счета при проведении оплаты через UDS Кассир.
	Slug                   string                 `json:"slug,omitempty"`                   // Доменное имя, которое отображается в ссылке на веб-страницу вашей компании.
}

// LoyaltyProgramSettings
// Настройки бонусной программы компании.
type LoyaltyProgramSettings struct {
	BaseMembershipTier    MembershipTier   `json:"baseMembershipTier,omitempty"`    // Настройки статусов клиентов.
	MembershipTiers       []MembershipTier `json:"membershipTiers,omitempty"`       // Настройки статусов.
	ReferralCashbackRates [3]float64       `json:"referralCashbackRates,omitempty"` // Коэффициенты начисления кешбэка для рефералов (3 уровня в %).
	CashierAward          float64          `json:"cashierAward,omitempty"`          // Процент вознаграждения кассиру за проведенную операцию.
	ReferralReward        float64          `json:"referralReward,omitempty"`        // Вознаграждение клиенту за эффективную рекомендацию.
	ReceiptLimit          float64          `json:"receiptLimit,omitempty"`          // Максимальная сумма операции, которую можно провести через UDS Кассир.
	DeferPointsForDays    int              `json:"deferPointsForDays,omitempty"`    // Количество дней, после которых будут начислены отложенные бонусные баллы.
	FirstPurchasePoints   float64          `json:"firstPurchasePoints,omitempty"`   // Number of points for the first purchase
}

// MembershipTier
// Статус клиента
type MembershipTier struct {
	Uid               string  `json:"uid,omitempty"`               // Идентификатор статуса.
	Name              string  `json:"name,omitempty"`              // Название статуса.
	Rate              float64 `json:"rate,omitempty"`              // Коэффициент статуса.
	MaxScoresDiscount float64 `json:"maxScoresDiscount,omitempty"` // Процент счета, который можно оплатить бонусными баллами.
	// Условия для автоматического назначения статуса.
	Conditions struct {
		// Повысить статус, когда сумма покупок достигнет данного значения.
		TotalCashSpent struct {
			// Сумма покупок.
			Target float64 `json:"target,omitempty"`
		} `json:"totalCashSpent,omitempty"`
		// Повысить уровень, когда клиент достигнет значения effectiveInvitedCount.
		EffectiveInvitedCount struct {
			// Количество эффективных рекомендаций.
			Target int `json:"target,omitempty"`
		} `json:"effectiveInvitedCount,omitempty"`
	} `json:"conditions,omitempty"`
}

// GetSettings Получение настроек компании
// https://docs.uds.app/#tag/Settings/paths/~1settings/get
func (u *Client) GetSettings() (*Settings, *resty.Response, error) {
	settings := new(Settings)

	resp, err := u.client.R().SetResult(settings).Get("settings")
	if err != nil {
		return nil, resp, err
	}

	return settings, resp, nil
}
