package ships

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/amead24/gotraders/pkgs/account"
)

// thinking about these structs - It only makes sense to return these objects
// if there would be a use for them.  If the goal is only to display them to the terminal
// then we might not need to make them public.
type PurchaseLocation struct {
	Location string `json:"location,omitempty"`
	Price    int    `json:"price,omitempty"`
	System   string `json:"system,omitempty"`
}

type Ship struct {
	Class             string             `json:"class,omitempty"`
	LoadingSpeed      int                `json:"loadingSpeed,omitempty"`
	Manufacturer      string             `json:"manufacturer,omitempty"`
	MaxCargo          int                `json:"maxCargo,omitempty"`
	Plating           int                `json:"plating,omitempty"`
	PurchaseLocations []PurchaseLocation `json:"purchaseLocations,omitempty"`
	RestrictedGoods   []string           `json:"restrictedGoods,omitempty"`
	Speed             int                `json:"speed,omitempty"`
	Type              string             `json:"type,omitempty"`
	Weapons           int                `json:"weapons,omitempty"`
}

type ShipListing struct {
	Ships []Ship `json:"shipListings,omitempty"`
}

// i've given up trying to create names -
// there's really got to be a better way
type Cargo struct {
	Good        string `json:"good,omitempty"`
	Quantity    int    `json:"quantity,omitempty"`
	TotalVolume int    `json:"totalVolume,omitempty"`
}

type BuyShipShipStruct struct {
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

type BuyShipStruct struct {
	Credits int               `json:"credits,omitempty"`
	Ship    BuyShipShipStruct `json:"ship,omitempty"`
}

func ListOtherShips(filter string) (ShipListing, error) {
	creds, err := account.GetUsernameAndToken()
	if err != nil {
		log.Fatalln(err)
	}

	url := fmt.Sprintf("https://api.spacetraders.io/systems/OE/ship-listings?token=%s", url.QueryEscape(creds.Token))
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	var shipListing ShipListing
	err = json.NewDecoder(resp.Body).Decode(&shipListing)
	if err != nil {
		log.Fatalln(err)
	}

	if filter == "" {
		return shipListing, nil

	} else {
		splitFilter := strings.Split(filter, "=")
		if len(splitFilter) != 2 {
			log.Fatalln("Filter must be formatted as key=value")
		}

		// key, value := splitFilter[0], splitFilter[1]
		_, value := "class", splitFilter[1]

		// Copy pasta - looks like a slice copies over the original list,
		// curious why this doesn't cause a buffer overflow of sorts, ex: [new, new, old]
		// https: //github.com/golang/go/wiki/SliceTricks#filtering-without-allocating

		// TODO: better understand more performative ways
		// not sure if there's a large downside to init a large arrray
		// https://www.ardanlabs.com/blog/2013/08/collections-of-unknown-length-in-go.html
		filteredShipListing := make([]Ship, 0, len(shipListing.Ships))

		for _, ship := range shipListing.Ships {
			// reflectedShip := reflect.ValueOf(&ship).Elem()
			// if reflectedShip.FieldByName(key) == value {
			// 	filteredShipListing = append(filteredShipListing, ship)
			// }

			if ship.Class == value {
				filteredShipListing = append(filteredShipListing, ship)
			}
		}

		return ShipListing{Ships: filteredShipListing}, nil
	}
}

func ListMyShips() ([]BuyShipShipStruct, error) {
	type MyShips struct {
		Ships []BuyShipShipStruct `json:"ships,omitempty"`
	}

	creds, err := account.GetUsernameAndToken()
	if err != nil {
		log.Fatalln(err)
	}

	url := fmt.Sprintf("https://api.spacetraders.io/my/ships?token=%s", url.QueryEscape(creds.Token))
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	var myShips MyShips
	err = json.NewDecoder(resp.Body).Decode(&myShips)
	if err != nil {
		log.Fatalln(err)
	}

	return myShips.Ships, nil
}

func BuyShip(shipType string, location string) (BuyShipStruct, error) {
	creds, err := account.GetUsernameAndToken()
	if err != nil {
		log.Fatalln(err)
	}

	url := fmt.Sprintf("https://api.spacetraders.io/my/ships?token=%s", url.QueryEscape(creds.Token))
	postBody, _ := json.Marshal(map[string]string{
		"type":     shipType,
		"location": location,
	})
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(url, "application/json", responseBody)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	// given the response, format it into the new struct
	var buyShipResponse BuyShipStruct
	err = json.NewDecoder(resp.Body).Decode(&buyShipResponse)
	if err != nil {
		log.Fatalln(err)
	}

	return buyShipResponse, nil
}
