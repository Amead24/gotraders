package main

import (
	"fmt"
	"log"
	"os"

	"github.com/amead24/gotraders/pkgs/account"
	"github.com/amead24/gotraders/pkgs/loans"
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

					fmt.Printf("Server Responded with %s", resp)
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

					fmt.Printf("Username & Token written to ~/.spacetravels/credentials")
					return nil
				},
			},
			{
				Name:  "status",
				Usage: "get status of account",
				Action: func(c *cli.Context) error {
					acct, err := account.ListAccount()
					if err != nil {
						return err
					}

					fmt.Printf("Account info == %+v", acct)
					return nil
				},
			},
			{
				Name:  "list",
				Usage: "list out loans, to be paramed for: https://api.spacetraders.io/#api-types",
				Action: func(c *cli.Context) error {
					loans, err := loans.ListLoans()
					if err != nil {
						return err
					}

					fmt.Printf("Loans:\n%+v", loans)
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
