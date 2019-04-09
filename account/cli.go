package account

import (
	"fmt"
	"os"

	"cgtcalc/config"
	"cgtcalc/model"
	"cgtcalc/cgtdb"

	"github.com/leekchan/accounting"
	"github.com/urfave/cli"
	log "github.com/Sirupsen/logrus"
)

var DelCommand = cli.Command{
	Name:      "delete",
	Aliases:   []string{"del"},
	Usage:     "Delete the gains Database",
	ArgsUsage: "blah",
	//Description: `
//`,
	Flags: []cli.Flag{
		config.ConfigFileFlag,
	},
	Action: func(ctx *cli.Context) error {
		cfg, err := config.LoadConfig(ctx.String("configFileFlag"))
		if err != nil {
			return err
		}
		log.Info("Deleting the Gains Database %s", cfg.GainsDatabase)

		err = os.RemoveAll(cfg.GainsDatabase)
		if err != nil {
			return err
			//log.Fatal(err)
		}
		return nil
	},
}

var DumpCommand = cli.Command{
	Name:      "dump",
	Usage:     "Prints all transactions in GainsDB",
	Flags: []cli.Flag{
	},
	Action: func(ctx *cli.Context) error {
		cry := accounting.Accounting{Precision: 4}
		cfg, err := config.LoadConfig(ctx.String("configFileFlag"))
		if err != nil {
			return err
		}

		db, err := cgtdb.NewLDBDatabase(cfg.GainsDatabase)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		m := *model.NewModel(db, nil,cfg)
		
		for _, i := range m.Txns {
			fmt.Printf("[%s - %s] Traded %s x %s to receive %s x %s\n",i.Date.Format("02 Jan 2006"), i.Exchange, cry.FormatMoney(i.BaseSpent), i.BaseCurrency, cry.FormatMoney(i.QuoteReceived), i.QuoteCurrency)
		}

		return nil
	},
}

