package loans

import (
	"strings"

	"github.com/amead24/gotraders/pkgs/utils"
)

type Loan struct {
	Amount             int    `json:"amount,omitempty"`
	CollateralRequired bool   `json:"collateralRequired,omitempty"`
	Rate               int    `json:"rate,omitempty"`
	TermInDays         int    `json:"termInDays,omitempty"`
	Type               string `json:"type,omitempty"`
}

type LoanReceipt struct {
	DueDate         string `json:"due,omitempty"`
	Id              string `json:"id,omitempty"`
	RepaymentAmount int    `json:"repaymentAmount,omitempty"`
	Status          string `json:"status,omitempty"`
	Type            string `json:"type,omitempty"`
}

func List(loanTypeFilter string) ([]Loan, error) {
	// TODO: i'm thinking something like boto: list loans -f type=STARTUP
	type Loans struct {
		Loans []Loan `json:"loans,omitempty"`
	}

	var loans Loans
	url := "https://api.spacetraders.io/types/loans"
	params := map[string]string{}

	err := utils.Get(url, params, &loans)
	if err != nil {
		return nil, err
	}

	if loanTypeFilter != "" {
		filteredLoans := make([]Loan, 0, len(loans.Loans))
		for _, loan := range loans.Loans {
			if loan.Type == strings.ToUpper(loanTypeFilter) {
				filteredLoans = append(filteredLoans, loan)
			}
		}

		return filteredLoans, nil
	}

	return loans.Loans, nil
}

func Buy(loanType string) (LoanReceipt, error) {
	type Debt struct {
		Credits int         `json:"credits,omitempty"`
		Receipt LoanReceipt `json:"loan,omitempty"`
	}

	var debt Debt
	url := "https://api.spacetraders.io/my/loans"
	params := map[string]string{
		"type": strings.ToUpper(loanType),
	}

	ok, err := utils.Post(url, params, &debt)
	if !ok {
		return LoanReceipt{}, err
	}

	return debt.Receipt, nil
}

func Owned() ([]LoanReceipt, error) {
	type ListLoans struct {
		Receipts []LoanReceipt `json:"loans,omitempty"`
	}

	var listLoans ListLoans
	url := "https://api.spacetraders.io/my/loans"
	params := map[string]string{}

	err := utils.Get(url, params, &listLoans)
	if err != nil {
		return nil, err
	}

	return listLoans.Receipts, nil
}
