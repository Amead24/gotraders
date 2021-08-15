package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"text/tabwriter"
	"time"

	"github.com/amead24/gotraders/pkgs/goods"
	"github.com/amead24/gotraders/pkgs/loans"
	"github.com/amead24/gotraders/pkgs/ships"
	"github.com/amead24/gotraders/pkgs/systems"
	"github.com/amead24/gotraders/pkgs/utils"
	"github.com/dustin/go-humanize"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{ // health
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
			{ // init
				Name:  "init",
				Usage: "get and set a new token",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "username",
						Usage:    "Username to create",
						Aliases:  []string{"u"},
						Required: true,
					},
					&cli.BoolFlag{
						Name:     "takeout-starter-loan",
						Aliases:  []string{"tsl"},
						Required: false,
					},
				},
				Action: func(c *cli.Context) error {
					username := c.String("username")
					tsl := c.Bool("takeout-starter-loan")
					ok, err := utils.SetUsernameAndToken(username)

					if !ok {
						fmt.Printf("Error: %s", err)
					} else {
						fmt.Println("Account creation successful!")
						if tsl {
							_, err := loans.Buy("STARTUP")
							if err != nil {
								return nil
							}
							fmt.Println("Startup loan recieved - Good luck.")
						}
					}

					return nil
				},
			},
			{ // status
				Name:  "status",
				Usage: "get status of account",
				Action: func(c *cli.Context) error {
					acct, err := utils.GetAccount()
					if err != nil {
						return err
					}

					w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
					fmt.Fprintf(w, "Username\t%s\n", acct.Username)
					fmt.Fprintf(w, "Ships\t%d\t\n", acct.ShipCount)
					fmt.Fprintf(w, "Structure\t%d\t\n", acct.StructureCount)
					fmt.Fprintf(w, "Credits\t%s\t\n", humanize.FormatInteger("", acct.Credits))
					w.Flush()

					return nil
				},
			},
			{ // loans
				Name:  "loans",
				Usage: "do loans-stuff",
				Subcommands: []*cli.Command{
					{ // loans list
						Name:  "list",
						Usage: "List out all available loans",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "filter",
								Usage:    "filter on loan type",
								Aliases:  []string{"f"},
								Required: false,
							},
						},
						Action: func(c *cli.Context) error {
							filter := c.String("filter")
							loans, err := loans.List(filter)
							if err != nil {
								return err
							}

							w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.AlignRight)
							fmt.Fprintln(w, " \tType\tTerm\tRate\tCollateral\tAmount\t")
							for i, loan := range loans {
								fmt.Fprintf(w, "%d\t%s\t%d\t%d\t%t\t%d\t\n", i, loan.Type, loan.TermInDays, loan.Rate, loan.CollateralRequired, loan.Amount)
							}
							w.Flush()

							return nil
						},
					},
					{ // loans buy
						Name:  "buy",
						Usage: "take out a loan",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "type",
								Usage:    "type of loan",
								Aliases:  []string{"t"},
								Required: true,
							},
						},
						Action: func(c *cli.Context) error {
							loanType := c.String("type")
							loan, err := loans.Buy(loanType)
							if err != nil {
								return err
							}

							fmt.Printf("Recieved loan:\n%+v", loan)
							return nil
						},
					},
					{ // loans owned
						Name:  "owned",
						Usage: "list out all owned loans",
						Action: func(c *cli.Context) error {
							loanReceipts, err := loans.Owned()
							if err != nil {
								return err
							}

							fmt.Println("These are the loans you owe:")
							w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.AlignRight)
							fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
							fmt.Fprintln(w, " \tID\tType\tStatus\tDue Date\tAmount\t")
							for i, receipt := range loanReceipts {
								fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\t\n", i, receipt.Id, receipt.Type, receipt.Status, receipt.DueDate, humanize.FormatInteger("", receipt.RepaymentAmount))
							}
							w.Flush()
							fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")

							return nil
						},
					},
				},
			},
			{ // ships
				Name:  "ships",
				Usage: "do shippy-stuff",
				Subcommands: []*cli.Command{
					{ // ships list
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

							ships, err := ships.List(filter)
							if err != nil {
								return err
							}

							fmt.Println("All ships for sale, to purchase run: `gotraders ships buy -t <type> -l <location>`")
							fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")

							w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.AlignRight)
							fmt.Fprintln(w, " \tManufacturer\tClass\tType\tSpeed\tPlating\tWeapons\tMax Cargo\tLoading Speed\tLocation\tPrice\t")
							tableIndex := 0
							for _, ship := range ships {
								for _, pl := range ship.PurchaseLocations {
									fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%d\t%d\t%d\t%d\t%d\t%s\t%s\t\n", tableIndex, ship.Manufacturer, ship.Class, ship.Type, ship.Speed, ship.Plating, ship.Weapons, ship.MaxCargo, ship.LoadingSpeed, pl.Location, humanize.FormatInteger("", pl.Price))
									tableIndex++
								}
							}
							w.Flush()

							fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
							return nil
						},
					},
					{ // ships buy
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

							shipReciept, err := ships.Buy(shipType, shipLocation)
							if err != nil {
								return err
							}

							fmt.Printf("Purchase of ship %s %s confirmed.\n", shipReciept.Ship.Manufacturer, shipReciept.Ship.Type)
							fmt.Printf("New account balance: %d\n", shipReciept.User.Credits)
							fmt.Printf("New ship id: %s\n", shipReciept.Ship.Id)
							return nil
						},
					},
					{ // ships owned
						Name:  "owned",
						Usage: "display information on ships owned",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "shipId",
								Aliases:  []string{"s"},
								Usage:    "Restrict readout to a specific ship",
								Required: false,
							},
						},
						Action: func(c *cli.Context) error {
							shipId := c.String("shipId")
							shipList, err := ships.Owned(shipId)
							if err != nil {
								return err
							}

							if len(shipList) == 0 {
								fmt.Println("Ship not found - Exiting")

							} else if len(shipList) == 1 {
								ship := shipList[0]
								left := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

								fmt.Fprintf(left, "Information on ship: %s\n", ship.Id)
								fmt.Fprintf(left, "Ship System\t%s\t\n", ship.Location)
								fmt.Fprintf(left, "Coordinates\t( %d, %d)\t\n", ship.X, ship.Y)

								left.Flush()

								right := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.AlignRight)
								fmt.Fprintln(right, "~~~~~~~~~~~~~~~ CARGO ~~~~~~~~~~~~~~~")
								fmt.Fprintln(right, " \tGood\tQuantity\tTotal Volume\t")
								for i, cargo := range ship.Cargo {
									fmt.Fprintf(right, "%d\t%s\t%d\t%d\t\n", i, cargo.Good, cargo.Quantity, cargo.TotalVolume)
								}
								fmt.Fprintln(right, "~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")

							} else {
								w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.AlignRight)
								fmt.Fprintln(w, "Your currrent ships; to get stats on just one use: `gotraders ships owned --id <id>`")
								fmt.Fprintln(w, "~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")

								// fmt.Fprintln(w, " \tID\tClass\tType\tLocation\tSystem\tSpeed\tPlating\tWeapons\tLoading Speed\tMax Cargo\tCargo Space Available\t")
								fmt.Fprintln(w, " \tID\tClass\tType\tLocation\tSystem\tMax Cargo\tCargo Space Available\t")
								for i, ship := range shipList {
									// fmt.Fprintf(w, "%d\t%s\t%s\t%s\t(%d, %d)\t%s\t%d\t%d\t%d\t%d\t%d\t%d\t\n", i, ship.Id, ship.Class, ship.Type, ship.X, ship.Y, ship.Location, ship.Speed, ship.Plating, ship.Weapons, ship.LoadingSpeed, ship.MaxCargo, ship.SpaceAvailable)
									fmt.Fprintf(w, "%d\t%s\t%s\t%s\t(%d, %d)\t%s\t%d\t%d\t\n", i, ship.Id, ship.Class, ship.Type, ship.X, ship.Y, ship.Location, ship.MaxCargo, ship.SpaceAvailable)
								}

								fmt.Fprintln(w, "~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
								w.Flush()
							}

							return nil
						},
					},
					{ // ships fp
						Name:  "fp",
						Usage: "command your ships flight paths",
						Subcommands: []*cli.Command{
							{ // ships fp start
								Name:  "start",
								Usage: "start a flight path for a given ship",
								Flags: []cli.Flag{
									&cli.StringFlag{
										Name:     "shipId",
										Usage:    "the ship id to interact with",
										Aliases:  []string{"s"},
										Required: true,
									},
									&cli.StringFlag{
										Name:     "destination",
										Usage:    "Destination for your ship",
										Aliases:  []string{"d"},
										Required: true,
									},
								},
								Action: func(c *cli.Context) error {
									shipId := c.String("shipId")
									destination := c.String("destination")
									fp, err := ships.CreateFlightPlan(shipId, destination)
									if err != nil {
										return err
									}

									fmt.Printf("Flight confirmed, track with flight ID: %+v\n", fp)

									return nil
								},
							},
							{
								Name:  "map",
								Usage: "list out ship to planet distance",
								Flags: []cli.Flag{
									&cli.StringFlag{
										Name:     "shipid",
										Usage:    "list relative to a specific ship",
										Aliases:  []string{"s"},
										Required: false,
									},
								},
								Action: func(c *cli.Context) error {
									shipid := c.String("shipid")
									ships, _ := ships.Owned(shipid)
									planets, _ := systems.List()

									fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
									w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.AlignRight)
									fmt.Fprintln(w, " \tShip\tPlanet\tDistance\t")
									idx := 0
									for _, ship := range ships {
										for _, planet := range planets {
											distance := math.Ceil(math.Sqrt(math.Pow(float64(ship.Y)-float64(planet.Y), 2) + math.Pow(float64(ship.X)-float64(planet.X), 2)))
											if math.IsNaN(distance) || math.IsInf(distance, 0) {
												distance = 0
											}

											fmt.Fprintf(w, "%d\t%s - %s\t%s - %s\t%d\t\n", idx, ship.Id, ship.Location, planet.Name, planet.Symbol, int64(distance))
											idx++
										}
									}
									w.Flush()
									fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")

									return nil
								},
							},
						},
					},
				},
			},
			{ // goods
				Name:  "goods",
				Usage: "do market stuff",
				Subcommands: []*cli.Command{
					{ // goods list
						Name:  "list",
						Usage: "list goods at market",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "good",
								Usage:    "specify a symbol as filter",
								Aliases:  []string{"g"},
								Required: false,
							},
							&cli.StringFlag{
								Name:     "shipId",
								Usage:    "restrict goods relative to ship location",
								Aliases:  []string{"s"},
								Required: false,
							},
							&cli.BoolFlag{
								Name:     "show-spread",
								Required: false,
							},
						},
						Action: func(c *cli.Context) error {
							good := c.String("good")
							shipId := c.String("shipId")
							showSpread := c.Bool("show-spread")

							ships, err := ships.Owned(shipId)
							if err != nil {
								return err
							}

							// Use map to easily de-dup a list
							systems := make(map[string]bool)
							for _, ship := range ships {
								if ship.Location != "" {
									systems[ship.Location] = true
								}
							}

							var goodsList []goods.Good
							for system := range systems {
								systemGoods, _ := goods.List(good, system)
								for i := 0; i < len(systemGoods); i++ {
									good := &systemGoods[i]
									good.System = system
								}

								goodsList = append(goodsList, systemGoods...)
							}

							if showSpread {
								goodsMap := make(map[string]map[string][]int)
								for _, good := range goodsList {
									if goodsMap[good.Symbol] == nil {
										goodsMap[good.Symbol] = make(map[string][]int)
										goodsMap[good.Symbol][good.System] = make([]int, 0)
									}

									goodsMap[good.Symbol][good.System] = append(goodsMap[good.Symbol][good.System], good.PricePerUnit)
									goodsMap[good.Symbol][good.System] = append(goodsMap[good.Symbol][good.System], good.SellPricePerUnit)

								}

								w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.AlignRight)
								fmt.Fprintln(w, " \tGood\tSystem\tBuy\tSell\t")
								idx := 0
								for symbol, systemMap := range goodsMap {
									if len(systemMap) > 1 {
										for system, spreadList := range systemMap {
											fmt.Fprintf(w, "%d\t%s\t%s\t%d\t%d\t\n", idx, symbol, system, spreadList[0], spreadList[1])
											idx++
										}
									}
								}
								w.Flush()

							} else {
								fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
								w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.AlignRight)
								fmt.Fprintln(w, " \tSymbol\tSystem\tQuantity\tVolume\tPrice\tPurchase Price\tSell Price\tSpread\t")
								for i, good := range goodsList {
									fmt.Fprintf(w, "%d\t%s\t%s\t%d\t%d\t%d\t%d\t%d\t%d\t\n", i, good.Symbol, good.System, good.QuantityAvailable, good.VolumePerUnit, good.PricePerUnit, good.PurchasePricePerUnit, good.SellPricePerUnit, good.Spread)
								}
								w.Flush()
								fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
							}
							return nil
						},
					},
					{
						Name:  "buy",
						Usage: "buy a good",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "shipId",
								Usage:    "-s shipId",
								Aliases:  []string{"s"},
								Required: true,
							},
							&cli.StringFlag{
								Name:     "good",
								Usage:    "nname of good",
								Aliases:  []string{"g"},
								Required: true,
							},
							&cli.IntFlag{
								Name:     "quantity",
								Usage:    "how much",
								Aliases:  []string{"q"},
								Required: true,
							},
						},
						Action: func(c *cli.Context) error {
							shipId := c.String("shipId")
							good := c.String("good")
							quantity := c.Int("quantity")
							_, err := goods.Buy(shipId, good, quantity)
							if err != nil {
								return nil
							}

							shipList, err := ships.Owned(shipId)
							if err != nil {
								return err
							}

							// probably could use better error checking
							// and a function to do this printing
							ship := shipList[0]

							fmt.Println("Order confirmed - Updating Ship Manifesto")
							w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.AlignRight)
							fmt.Fprintf(w, "Information on ship: %s\n", ship.Id)
							fmt.Fprintf(w, "Ship System\t%s\t\n", ship.Location)
							fmt.Fprintf(w, "Coordinates\t( %d, %d)\t\n", ship.X, ship.Y)

							w.Flush()

							fmt.Fprintln(w, "~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
							fmt.Fprintln(w, " \tGood\tQuantity\tTotal Volume\t")
							for i, cargo := range ship.Cargo {
								fmt.Fprintf(w, "%d\t%s\t%d\t%d\t\n", i, cargo.Good, cargo.Quantity, cargo.TotalVolume)
							}
							fmt.Fprintln(w, "~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
							return nil
						},
					},
					{
						Name:  "sell",
						Usage: "sell a good",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "shipId",
								Usage:    "-s shipId",
								Aliases:  []string{"s"},
								Required: true,
							},
							&cli.StringFlag{
								Name:     "good",
								Usage:    "name of good",
								Aliases:  []string{"g"},
								Required: true,
							},
							&cli.IntFlag{
								Name:     "quantity",
								Usage:    "how much",
								Aliases:  []string{"q"},
								Required: true,
							},
						},
						Action: func(c *cli.Context) error {
							shipId := c.String("shipId")
							good := c.String("good")
							quantity := c.Int("quantity")
							updatedShip, err := goods.Sell(shipId, good, quantity)
							if err != nil {
								return nil
							}

							fmt.Println("Order confirmed - Updating Ship Manifesto")
							w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.AlignRight)
							fmt.Fprintf(w, "Information on ship: %s\n", updatedShip.Ship.Id)
							fmt.Fprintf(w, "Ship System\t%s\t\n", updatedShip.Ship.Location)
							fmt.Fprintf(w, "Coordinates\t( %d, %d)\t\n", updatedShip.Ship.X, updatedShip.Ship.Y)

							w.Flush()

							fmt.Fprintln(w, "~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
							fmt.Fprintln(w, " \tGood\tQuantity\tTotal Volume\t")
							for i, cargo := range updatedShip.Ship.Cargo {
								fmt.Fprintf(w, "%d\t%s\t%d\t%d\t\n", i, cargo.Good, cargo.Quantity, cargo.TotalVolume)
							}
							fmt.Fprintln(w, "~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
							return nil
						},
					},
				},
			},
			{ // systems
				Name:  "systems",
				Usage: "list all systems",
				Action: func(c *cli.Context) error {
					planetList, err := systems.List()
					if err != nil {
						return err
					}

					fmt.Println("Listing out all nearby planets:")
					fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")

					w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.AlignRight)
					fmt.Fprintln(w, " \tName\tSymbol\tType\tCoords\t")

					for i, planet := range planetList {
						fmt.Fprintf(w, "%d\t%s\t%s\t%s\t( %d, %d)\t\n", i, planet.Name, planet.Symbol, planet.Type, planet.X, planet.Y)
					}

					w.Flush()
					fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")

					return nil
				},
			},
			{
				Name: "algo",
				Action: func(c *cli.Context) error {
					fmt.Printf("%s - Starting back and forth...\n", time.Now())
					shipId := "cks8soxdn3430611ds6nhd90kp4"
					acct, _ := utils.GetAccount()
					credits := acct.Credits
					// startingLocation := "OE-PM-TR"
					// nextLocation := "OE-PM"

					for {
						// Moon -> Prime
						// fmt.Printf("%s - Starting at system %s and flying to %s\n", time.Now(), "OE-PM-TR", "OM-PM")
						goods.Buy(shipId, "metals", 290)
						goods.Buy(shipId, "fuel", 5)

						fpToPrime, _ := ships.CreateFlightPlan(shipId, "OE-PM")
						fmt.Printf("%s - Leaving %s for %s, ETA: %d\n", time.Now(), fpToPrime.Departure, fpToPrime.Destination, fpToPrime.TimeRemainingInSeconds)
						time.Sleep(time.Duration(fpToPrime.TimeRemainingInSeconds) * time.Second)

						receipt, _ := goods.Sell(shipId, "metals", 290)
						fmt.Printf("%s - Metals sold, new account balance: %d\n", time.Now(), receipt.Credits)

						// return flight home
						fpToMoon, _ := ships.CreateFlightPlan(shipId, "OE-PM-TR")
						fmt.Printf("%s - Leaving %s for %s, ETA: %d\n", time.Now(), fpToMoon.Departure, fpToMoon.Destination, fpToMoon.TimeRemainingInSeconds)
						time.Sleep(time.Duration(fpToMoon.TimeRemainingInSeconds) * time.Second)

						if credits >= receipt.Credits {
							fmt.Printf("Credits decreasing (gain = %d) - Aborting.", receipt.Credits-credits)
							break
						}
						credits = receipt.Credits
					}

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
