package loans

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/amead24/gotraders/pkgs/account"
)

type Loan struct {
	Amount             int    `json:"amount,omitempty"`
	CollateralRequired bool   `json:"collateralRequired,omitempty"`
	Rate               int    `json:"rate,omitempty"`
	TermInDays         int    `json:"termInDays,omitempty"`
	Type               string `json:"type,omitempty"`
}

type Loans struct {
	Loans []Loan `json:"loans"`
}

func ListLoans() (string, error) {
	creds, err := account.GetUsernameAndToken()
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://api.spacetraders.io/types/loans?token=%s", url.QueryEscape(creds.Token))
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	// given the response, format it into the new struct
	var responseLoans Loans
	err = json.NewDecoder(resp.Body).Decode(&responseLoans)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%+v", responseLoans), nil
}
