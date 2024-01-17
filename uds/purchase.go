package uds

// PurchaseDetail
// Информация об операции.
type PurchaseDetail struct {
	MaxPoints          float64 `json:"maxPoints"`          // Максимальное количество бонусных баллов, доступное для списания.
	Total              float64 `json:"total"`              //  Общая сумма счета (в денежных единицах).
	SkipLoyaltyTotal   float64 `json:"skipLoyaltyTotal"`   // Часть суммы счета, на которую не начисляется кешбэк и на которую не распространяется скидка (в денежных единицах).
	UnredeemableTotal  float64 `json:"unredeemableTotal"`  // Часть суммы счета, которую нельзя погасить баллами.
	DiscountAmount     float64 `json:"discountAmount"`     // Размер скидки (в денежных единицах).
	DiscountPercent    float64 `json:"discountPercent"`    // Предоставленная скидка (в процентах).
	Points             float64 `json:"points"`             // Бонусных баллов к оплате.
	PointsPercent      float64 `json:"pointsPercent"`      // Размер скидки за счет бонусных баллов (в процентах).
	NetDiscount        float64 `json:"netDiscount"`        // Общий размер скидки (в денежных единицах).
	NetDiscountPercent float64 `json:"netDiscountPercent"` // Общий размер скидки (в процентах от общей суммы счета).
	CertificatePoints  float64 `json:"certificatePoints"`  // Количество списываемых бонусных баллов сертификата (в денежных единицах).
	Cash               float64 `json:"cash"`               // Сумма к оплате (в денежных единицах).
	CashTotal          float64 `json:"cashTotal"`          // Итоговая сумма к оплате с учетом доставки.
	CashBack           float64 `json:"cashBack"`           // Вознаграждение (кешбэк), которое получит клиент после проведения операции (в бонусных баллах).
	Extras             struct {
		Delivery float64 `json:"delivery"` // Стоимость доставки.
	} `json:"extras"` // Дополнительный платежи, на которые не распространяется программа лояльности.
	MaxScoresDiscount float64 `json:"maxScoresDiscount"` // Процент счета, который можно оплатить бонусными баллами.
}
