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
				Name:  "list",
				Usage: "list out smaller types: https://api.spacetraders.io/#api-types",
				Subcommands: []*cli.Command{
					{
						Name:  "loans",
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
						Name:  "ships",
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

							ships, err := ships.ListShips(filter)
							if err != nil {
								return err
							}

							fmt.Println(ships)
							return nil
						},
					},
				},
			},
			{
				// Has me thinking if you should do both with this command?
				// gotraders loan list vs. gotraders list loans
				Name:  "loan",
				Usage: "take out a loan",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "type",
						Usage:   "Type of loan to takeout, use `gotraders list loans` to see them all",
						Aliases: []string{"t"},
					},
				},
				Action: func(c *cli.Context) error {
					loanType := c.String("type")
					loan, err := loans.TakeoutLoan(loanType)
					if err != nil {
						return err
					}

					fmt.Printf("Recieved loan:\n%+v", loan)
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
