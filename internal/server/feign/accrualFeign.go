package feign

import (
	"encoding/json"
	"errors"
	"github.com/Lovchik/gophermart/internal/server/retry"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type AccuralFeing struct {
	AccuralSystemAddress string
}

func NewAccuralFeign(adress string) (*AccuralFeing, error) {
	return &AccuralFeing{AccuralSystemAddress: adress}, nil
}

type BonusInfo struct {
	Order   *string  `json:"order"`
	Status  *string  `json:"status"`
	Accrual *float64 `json:"accrual"`
}

func (ac AccuralFeing) GetBonusInfo(orderNumber string) (BonusInfo, error) {
	var bonusInfo BonusInfo
	client := &http.Client{}
	req, err := http.NewRequest("GET", ac.AccuralSystemAddress+"/api/orders/"+orderNumber, nil)
	if err != nil {
		return BonusInfo{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	response, err := retry.Retry(3, 1, func() (*http.Response, error) {
		return client.Do(req)
	})
	log.Info(response.StatusCode)
	if err != nil {
		return BonusInfo{}, err
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusNoContent {
		return BonusInfo{}, nil
	}
	if response.StatusCode != http.StatusOK {
		return BonusInfo{}, errors.New(response.Status)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return BonusInfo{}, err
	}
	err = json.Unmarshal(body, &bonusInfo)
	if err != nil {
		return BonusInfo{}, err
	}
	return bonusInfo, nil
}
