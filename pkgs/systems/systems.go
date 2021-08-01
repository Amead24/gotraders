package systems

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/amead24/gotraders/pkgs/account"
)

type Planet struct {
	Name   string `json:"name,omitempty"`
	Symbol string `json:"symbol,omitempty"`
	Type   string `json:"type,omitempty"`
	X      int    `json:"x,omitempty"`
	Y      int    `json:"y,omitempty"`
}

// lesson 4 - boiler plate & memory management
//
// feels like unlimited potential when you have access to this level memory
// it's not clear if I should be rturning anything when modifying an object?
func Get(url string, params map[string]string, obj interface{}) error {
	creds, err := account.GetUsernameAndToken()
	if err != nil {
		log.Fatalln(err)
	}

	var queryParams []string
	for key, value := range params {
		queryParams = append(queryParams, fmt.Sprintf("?%s=%s", key, value))
	}

	queryParamString := strings.Join(queryParams, "")
	urlWithParams := fmt.Sprintf("%s?token=%s&%s", url, creds.Token, queryParamString)

	resp, err := http.Get(urlWithParams)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	bytesBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	err = json.Unmarshal(bytesBody, obj)
	return err
}

func List() ([]Planet, error) {
	type Locations struct {
		Planets []Planet `json:"locations,omitempty"`
	}

	params := map[string]string{
		"type": "PLANET",
	}

	var location Locations
	Get("https://api.spacetraders.io/systems/OE/locations", params, &location)

	return location.Planets, nil
}
