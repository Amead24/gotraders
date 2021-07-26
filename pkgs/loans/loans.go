package loans

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

func (l Loan) String() string {
	return fmt.Sprintf("\tAmount: %d\n\tColllateral: %t\n\tRate: %d\n\tTerm (days): %d\n\tType: %s", l.Amount, l.CollateralRequired, l.Rate, l.TermInDays, l.Type)
}

type Loans struct {
	Loans []Loan `json:"loans"`
}

func (ls Loans) String() string {
	s := fmt.Sprintln("Loans:")

	for _, l := range ls.Loans {
		s += fmt.Sprintf("%+v", l)
		s += fmt.Sprintln("\n\t----------------")
	}

	return s
}

type DebtLoan struct {
	Due             string `json:"due"`
	Id              string `json:"id"`
	RepaymentAmount int    `json:"repaymentAmount`
	Status          string `json:"status"`
	Type            string `json:"type"`
}

func (dl DebtLoan) String() string {
	return fmt.Sprintf("\tID: %d\n\tDue Date: %s\n\tBalance: %d\t\nStatus: %s\t\nType: %s\n", dl.Id, dl.Due, dl.RepaymentAmount, dl.Status, dl.Type)
}

type Debt struct {
	Credits int `json:"credits"`
	Loan    DebtLoan
}

func (d Debt) String() string {
	return fmt.Sprintf("Credits Recieved: %d\n", d.Credits)
}

func ListLoans(filter string) (Loans, error) {
	// TODO: i'm thinking something like boto: list loans -f type=STARTUP

	creds, err := account.GetUsernameAndToken()
	if err != nil {
		return Loans{}, err
	}

	url := fmt.Sprintf("https://api.spacetraders.io/types/loans?token=%s", url.QueryEscape(creds.Token))
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("body == %s", string(body))

	// given the response, format it into the new struct
	var loans Loans
	err = json.NewDecoder(resp.Body).Decode(&loans)
	if err != nil {
		return Loans{}, err
	}

	return loans, nil
}

func TakeoutLoan(loanType string) (Debt, error) {
	creds, err := account.GetUsernameAndToken()
	if err != nil {
		log.Fatalln(err)
	}

	url := fmt.Sprintf("https://api.spacetraders.io/my/loans?token=%s", url.QueryEscape(creds.Token))
	postBody, _ := json.Marshal(map[string]string{
		"type": loanType,
	})
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(url, "application/json", responseBody)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	// given the response, format it into the new struct
	var debt Debt
	err = json.NewDecoder(resp.Body).Decode(&debt)
	if err != nil {
		log.Fatalln(err)
	}

	return debt, nil
}
