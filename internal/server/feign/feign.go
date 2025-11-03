package feign

type Feign interface {
	GetBonusInfo(orderNumber string) (BonusInfo, error)
}
