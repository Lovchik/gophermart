package feign

import (
	"encoding/json"
	"errors"
	"gofermart/internal/server/config"
	"io"
	"net/http"
)

type BonusInfo struct {
	Order   *string  `json:"order"`
	Status  *string  `json:"status"`
	Accrual *float64 `json:"accrual"`
}

func GetBonusInfo(orderNumber string) (BonusInfo, error) {
	var bonusInfo BonusInfo
	client := &http.Client{}
	req, err := http.NewRequest("PUT", config.GetConfig().AccuralSystemAddress+"/api/orders/"+orderNumber, nil) // Example PUT request
	if err != nil {
		return BonusInfo{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	response, err := client.Do(req)
	if err != nil {
		return BonusInfo{}, err
	}
	defer response.Body.Close()
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
