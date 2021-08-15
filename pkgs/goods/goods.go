package goods

import (
	"fmt"
	"strings"

	"github.com/amead24/gotraders/pkgs/ships"
	"github.com/amead24/gotraders/pkgs/utils"
)

type Good struct {
	PricePerUnit         int    `json:"pricePerUnit,omitempty"`
	PurchasePricePerUnit int    `json:"purchasePricePerUnit,omitempty"`
	QuantityAvailable    int    `json:"quantityAvailable,omitempty"`
	SellPricePerUnit     int    `json:"sellPricePerUnit,omitempty"`
	Spread               int    `json:"spread,omitempty"`
	Symbol               string `json:"symbol,omitempty"`
	VolumePerUnit        int    `json:"volumePerUnit,omitempty"`
	System               string
}

func List(goodFilter string, system string) ([]Good, error) {
	// TODO: These should look up where all your ships are
	// and list the union of all of them with that location
	type Marketplace struct {
		GoodsList []Good `json:"marketplace"`
	}

	url := fmt.Sprintf("https://api.spacetraders.io/locations/%s/marketplace", system)
	params := map[string]string{}
	var mp Marketplace

	err := utils.Get(url, params, &mp)
	if err != nil {
		return make([]Good, 0), err
	}

	if goodFilter != "" {
		filteredGoodsList := make([]Good, 0, len(mp.GoodsList))
		for _, good := range mp.GoodsList {
			if good.Symbol == goodFilter {
				filteredGoodsList = append(filteredGoodsList, good)
			}
		}

		return filteredGoodsList, nil
	}

	return mp.GoodsList, nil
}

func Buy(shipId string, good string, quantity int) (string, error) {
	type Order struct {
		Good         string `json:"good,omitempty"`
		PricePerUnit int    `json:"pricePerUnit,omitempty"`
		Quantity     int    `json:"quantity,omitempty"`
		Total        int    `json:"total,omitempty"`
	}

	type GoodsBuy struct {
		Credits int        `json:"credits,omitempty"`
		Order   Order      `json:"order,omitempty"`
		Ship    ships.Ship `json:"ship,omitempty"`
	}

	var gb GoodsBuy
	url := "https://api.spacetraders.io/my/purchase-orders"
	params := map[string]string{
		"shipId":   shipId,
		"good":     strings.ToUpper(good),
		"quantity": fmt.Sprint(quantity),
	}

	_, err := utils.Post(url, params, &gb)
	if err != nil {
		return "", err
	}

	return "", nil
}

type SellReceipt struct {
	Credits int        `json:"credits,omitempty"`
	Order   Good       `json:"order,omitempty"`
	Ship    ships.Ship `json:"ship,omitempty"`
}

func Sell(shipId string, good string, quantity int) (SellReceipt, error) {
	var sr SellReceipt
	url := "https://api.spacetraders.io/my/sell-orders"
	params := map[string]string{
		"shipId":   shipId,
		"good":     good,
		"quantity": fmt.Sprintf("%d", quantity),
	}

	ok, err := utils.Post(url, params, &sr)
	if !ok {
		fmt.Printf("Error: %s", err)
		return SellReceipt{}, err
	}

	return sr, nil
}
