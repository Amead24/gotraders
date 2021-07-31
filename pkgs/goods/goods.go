package goods

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/amead24/gotraders/pkgs/account"
)

type Good struct {
	PricePerUnit         int    `json:"pricePerUnit,omitempty"`
	PurchasePricePerUnit int    `json:"purchasePricePerUnit,omitempty"`
	QuantityAvailable    int    `json:"quantityAvailable,omitempty"`
	SellPricePerUnit     int    `json:"sellPricePerUnit,omitempty"`
	Spread               int    `json:"spread,omitempty"`
	Symbol               string `json:"symbol,omitempty"`
	VolumePerUnit        int    `json:"volumePerUnit,omitempty"`
}

func List(symbol string) ([]Good, error) {
	type Marketplace struct {
		GoodsList []Good `json:"marketplace"`
	}

	creds, err := account.GetUsernameAndToken()
	if err != nil {
		return make([]Good, 0), err
	}

	url := fmt.Sprintf("https://api.spacetraders.io/locations/OE-PM-TR/marketplace?token=%s", creds.Token)
	resp, err := http.Get(url)
	if err != nil {
		return make([]Good, 0), err
	}

	var mp Marketplace
	err = json.NewDecoder(resp.Body).Decode(&mp)
	if err != nil {
		return make([]Good, 0), err
	}

	if symbol != "" {
		filteredGoodsList := make([]Good, 0, len(mp.GoodsList))
		for _, good := range mp.GoodsList {
			if good.Symbol == symbol {
				filteredGoodsList = append(filteredGoodsList, good)
			}
		}

		return filteredGoodsList, nil
	}

	return mp.GoodsList, nil
}

func Buy(shipId string, good string, quantity int) (string, error) {
	// What does the user rreally neeed to know afteer buying something?
	type Cargo struct {
		Good        string `json:"good,omitempty"`
		Quantity    int    `json:"quantity,omitempty"`
		TotalVolume int    `json:"totalVolume,omitempty"`
	}

	type Ship struct {
		Id             string  `json:"id,omitempty"`
		Location       string  `json:"location,omitempty"`
		X              int     `json:"x,omitempty"`
		Y              int     `json:"y,omitempty"`
		Cargo          []Cargo `json:"cargo,omitempty"`
		SpaceAvailable int     `json:"spaceAvailable,omitempty"`
		Type           string  `json:"type,omitempty"`
		Class          string  `json:"class,omitempty"`
		MaxCargo       int     `json:"maxCargo,omitempty"`
		LoadingSpeed   int     `json:"loadingSpeed,omitempty"`
		Speed          int     `json:"speed,omitempty"`
		Manufacturer   string  `json:"manufacturer,omitempty"`
		Plating        int     `json:"plating,omitempty"`
		Weapons        int     `json:"weapons,omitempty"`
	}

	type Order struct {
		Good         string `json:"good,omitempty"`
		PricePerUnit int    `json:"pricePerUnit,omitempty"`
		Quantity     int    `json:"quantity,omitempty"`
		Total        int    `json:"total,omitempty"`
	}

	type GoodsBuy struct {
		Credits int   `json:"credits,omitempty"`
		Order   Order `json:"order,omitempty"`
		Ship    Ship  `json:"ship,omitempty"`
	}

	creds, err := account.GetUsernameAndToken()
	if err != nil {
		return "", nil
	}

	url := fmt.Sprintf("https://api.spacetraders.io/my/purchase-orders?token=%s", creds.Token)
	postBody, _ := json.Marshal(map[string]string{
		"shipId":   shipId,
		"good":     good,
		"quantity": fmt.Sprint(quantity),
	})
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(url, "application/json", responseBody)
	if err != nil {
		log.Println(err)
		return "", nil
	}
	defer resp.Body.Close()

	// Problem:
	// this causes  the next parsing of resp.Body to be empty
	// is this triggerinng the defer?
	// https://stackoverflow.com/a/43021236/5660197
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		log.Fatal(err)
	}
	fmt.Println(string(bodyBytes))

	var gb GoodsBuy
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&gb); err != nil {
		if terr, ok := err.(*json.UnmarshalTypeError); ok {
			fmt.Printf("Failed to unmarshal field %s \n", terr.Field)
		} else {
			fmt.Println("global fail")
			fmt.Printf("error == %s\n", err)
		}
	} else {
		fmt.Println(gb)
	}

	fmt.Printf("gb == %+v", gb)

	return fmt.Sprintf("New total: %d", gb.Order.Total), nil
}
