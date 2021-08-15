package ships

import (
	"log"
	"strings"

	"github.com/amead24/gotraders/pkgs/utils"
)

type PurchaseLocation struct {
	Location string `json:"location,omitempty"`
	Price    int    `json:"price,omitempty"`
	System   string `json:"system,omitempty"`
}

type Cargo struct {
	Good        string `json:"good,omitempty"`
	Quantity    int    `json:"quantity,omitempty"`
	TotalVolume int    `json:"totalVolume,omitempty"`
}

type Ship struct {
	Id                string             `json:"id,omitempty"`
	Location          string             `json:"location,omitempty"`
	X                 int                `json:"x,omitempty"`
	Y                 int                `json:"y,omitempty"`
	Cargo             []Cargo            `json:"cargo,omitempty"`
	SpaceAvailable    int                `json:"spaceAvailable,omitempty"`
	Type              string             `json:"type,omitempty"`
	Class             string             `json:"class,omitempty"`
	MaxCargo          int                `json:"maxCargo,omitempty"`
	LoadingSpeed      int                `json:"loadingSpeed,omitempty"`
	Speed             int                `json:"speed,omitempty"`
	Manufacturer      string             `json:"manufacturer,omitempty"`
	Plating           int                `json:"plating,omitempty"`
	Weapons           int                `json:"weapons,omitempty"`
	PurchaseLocations []PurchaseLocation `json:"purchaseLocations,omitempty"`
	RestrictedGoods   []string           `json:"restrictedGoods,omitempty"`
}

type User struct {
	Credits int `json:"credits,omitempty"`
}

type ShipReceipt struct {
	User User `json:"user,omitempty"`
	Ship Ship `json:"ship,omitempty"`
}

type FlightPlan struct {
	ArrivesAt              string `json:"arrivesAt,omitempty"`
	CreatedAt              string `json:"createdAt,omitempty"`
	Departure              string `json:"departure,omitempty"`
	Destination            string `json:"destination,omitempty"`
	Distance               int    `json:"distance,omitempty"`
	FuelConsumed           int    `json:"fuelConsumed,omitempty"`
	FuelRemaining          int    `json:"fuelRemaining,omitempty"`
	Id                     string `json:"id,omitempty"`
	ShipId                 string `json:"shipId,omitempty"`
	TerminatedAt           string `json:"terminatedAt,omitempty"`
	TimeRemainingInSeconds int    `json:"timeRemainingInSeconds,omitempty"`
}

func List(filter string) ([]Ship, error) {
	type ShipListing struct {
		Ships []Ship `json:"shipListings,omitempty"`
	}

	var shipListing ShipListing
	url := "https://api.spacetraders.io/systems/OE/ship-listings"
	params := map[string]string{}

	err := utils.Get(url, params, &shipListing)
	if err != nil {
		return make([]Ship, 0), err
	}

	if filter == "" {
		return shipListing.Ships, nil

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

		return filteredShipListing, nil
	}
}

func Buy(shipType string, location string) (ShipReceipt, error) {
	var shipReceipt ShipReceipt
	url := "https://api.spacetraders.io/my/ships"
	params := map[string]string{
		"type":     shipType,
		"location": location,
	}

	ok, err := utils.Post(url, params, &shipReceipt)
	if !ok {
		return ShipReceipt{}, err
	}

	return shipReceipt, nil
}

func Owned(shipId string) ([]Ship, error) {
	type MyShips struct {
		Ships []Ship `json:"ships,omitempty"`
	}

	var myShips MyShips
	params := map[string]string{}
	url := "https://api.spacetraders.io/my/ships"

	err := utils.Get(url, params, &myShips)
	if err != nil {
		return nil, err
	}

	if shipId == "" {
		return myShips.Ships, nil
	}

	filteredShipList := make([]Ship, 0, len(myShips.Ships))
	for _, ship := range myShips.Ships {
		if ship.Id == shipId {
			filteredShipList = append(filteredShipList, ship)
		}
	}

	return filteredShipList, nil
}

func CreateFlightPlan(shipId string, destination string) (FlightPlan, error) {
	type FPResponse struct {
		FP FlightPlan `json:"flightPlan,omitempty"`
	}

	var fpr FPResponse
	url := "https://api.spacetraders.io/my/flight-plans"
	params := map[string]string{
		"shipId":      shipId,
		"destination": destination,
	}

	ok, err := utils.Post(url, params, &fpr)
	if !ok {
		return FlightPlan{}, err
	}

	return fpr.FP, nil
}
