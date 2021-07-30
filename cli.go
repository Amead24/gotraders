package main

import (
	"fmt"
	"log"
	"os"

	"github.com/amead24/gotraders/pkgs/account"
	"github.com/amead24/gotraders/pkgs/loans"
	"github.com/amead24/gotraders/pkgs/ships"
	"github.com/amead24/gotraders/pkgs/utils"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "health",
				Usage: "check status of server",
				Action: func(c *cli.Context) error {
					resp, err := utils.GetServerHealth()

					if err != nil {
						return err
					}

					fmt.Printf("Server Responded with %s\n", resp)
					return nil
				},
			},
			{
				Name:  "init",
				Usage: "get and set a new token",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "username",
						Usage:   "Username to create",
						Aliases: []string{"u"},
					},
				},
				Action: func(c *cli.Context) error {
					// it would be a nice todo to get the starter loan on init
					username := c.String("username")
					_, err := account.SetUsernameAndToken(username)

					if err != nil {
						return err
					}

					fmt.Printf("Username & Token written to ~/.spacetravels/credential\n")
					return nil
				},
			},
			{
				Name:  "status",
				Usage: "get status of account",
				Action: func(c *cli.Context) error {
					acct, err := account.GetAccount("", "")
					if err != nil {
						return err
					}

					fmt.Printf("Account Information:\n%+v\n", acct)
					return nil
				},
			},
			{
				Name:  "loans",
				Usage: "do loans-stuff",
				Subcommands: []*cli.Command{
					{
						Name:  "list",
						Usage: "List out all available loans",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "filter",
								Usage:   "a key & value pair to map on",
								Aliases: []string{"f"},
							},
						},
						Action: func(c *cli.Context) error {
							filter := c.String("filter")

							loans, err := loans.ListLoans(filter)
							if err != nil {
								return err
							}

							fmt.Printf("%+v", loans)
							return nil
						},
					},
					{
						Name:  "buy",
						Usage: "take out a loan",
						Action: func(c *cli.Context) error {
							loanType := c.String("type")
							loan, err := loans.BuyLoan(loanType)
							if err != nil {
								return err
							}

							fmt.Printf("Recieved loan:\n%+v", loan)
							return nil
						},
					},
					{
						Name:  "owned",
						Usage: "list out all owned loans",
						Action: func(c *cli.Context) error {
							loans, err := loans.ListOwnedLoans()
							if err != nil {
								return err
							}

							fmt.Printf("These are the loans you owe:\n%+v\n", loans)
							return nil
						},
					},
				},
			},
			{
				Name:  "ships",
				Usage: "do shippy-stuff",
				Subcommands: []*cli.Command{
					{
						Name:  "list",
						Usage: "List out all available ships",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "filter",
								Usage:   "a key & value pair to map on",
								Aliases: []string{"f"},
							},
						},
						Action: func(c *cli.Context) error {
							filter := c.String("filter")

							ships, err := ships.ListOtherShips(filter)
							if err != nil {
								return err
							}

							fmt.Println(ships)
							return nil
						},
					},
					{
						Name:  "buy",
						Usage: "buy a ship based on a filter and location",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "type",
								Usage:    "which type to buy",
								Aliases:  []string{"t"},
								Required: true,
							},
							&cli.StringFlag{
								Name:     "location",
								Usage:    "location of the ship",
								Aliases:  []string{"l"},
								Required: true,
							},
						},
						Action: func(c *cli.Context) error {
							shipType := c.String("type")
							shipLocation := c.String("location")
							boughtShip, err := ships.BuyShip(shipType, shipLocation)
							if err != nil {
								return err
							}

							fmt.Printf("Bought ship: %+v\n", boughtShip)
							return nil
						},
					},
					{
						Name:  "owned",
						Usage: "display information on ships owned",
						Action: func(c *cli.Context) error {
							shipList, err := ships.ListMyShips()
							if err != nil {
								return err
							}

							fmt.Printf("Your Ships:\n%+v\n", shipList)
							return nil
						},
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
