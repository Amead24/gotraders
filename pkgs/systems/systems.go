package systems

import "github.com/amead24/gotraders/pkgs/utils"

type Planet struct {
	Name   string `json:"name,omitempty"`
	Symbol string `json:"symbol,omitempty"`
	Type   string `json:"type,omitempty"`
	X      int    `json:"x,omitempty"`
	Y      int    `json:"y,omitempty"`
}

func List() ([]Planet, error) {
	type Locations struct {
		Planets []Planet `json:"locations,omitempty"`
	}

	var location Locations
	url := "https://api.spacetraders.io/systems/OE/locations"
	params := map[string]string{
		"type": "PLANET",
	}

	err := utils.Get(url, params, &location)
	if err != nil {
		return nil, err
	}

	return location.Planets, nil
}
