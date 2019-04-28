package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"cgtcalc/account"
	"cgtcalc/cgtdb"
	"cgtcalc/config"
	"cgtcalc/exchanges"
	"cgtcalc/model"
	"cgtcalc/prices"
	"cgtcalc/reporting"
	"cgtcalc/version"
	"cgtcalc/wizard"

	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

var year = time.Now().Year()

var app *cli.App

var loadCommand = cli.Command{
	Name:      "load",
	Aliases:   []string{"l"},
	Usage:     "Adds a CSV/XLXS content to the database",
	ArgsUsage: "first argument is the file (CSV) then -exchange <<name of exchange>>",
	Flags: []cli.Flag{
		config.ConfigFileFlag,
		cli.StringFlag{
			Name:  "exchange",
			Value: "btcmarkets",
			Usage: "Exchange that CSV was exported from",
		},
	},
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			log.Fatal("This command requires an argument.")
		}

		cfg, err := config.LoadConfig(ctx.String("configFileFlag"))
		if err != nil {
			return err
		}
		fp := ctx.Args().First()

		log.Info("The Database being added to is: ", cfg.GainsDatabase)

		db, err := cgtdb.NewLDBDatabase(cfg.GainsDatabase)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		m := *model.NewModel(db, nil, cfg)
		log.Info("Opening Number of Transactions before load: ", len(m.Txns))

		txns := []model.Transaction{}

		log.Info("The Exchange Loaded is: ", strings.ToLower(ctx.String("exchange")))
		if m.Checksum(fp) {
			processor, err := exchanges.ExchangeFuncSearch(ctx.String("exchange"))
			if err != nil {
				log.Fatal(err)
			}
			txns = processor(fp)
		} else {
			log.Fatal("Could not process file as it has already been committed to the database")
		}

		m.AddTxns(txns)

		return nil
	},
}

var processCommand = cli.Command{
	Name:    "process",
	Aliases: []string{"p"},
	Usage:   "Iterates through the transactions to establish the CGT",
	Flags: []cli.Flag{
		config.ConfigFileFlag,
	},
	Action: func(ctx *cli.Context) error {

		cfg, err := config.LoadConfig(ctx.String("configFileFlag"))
		if err != nil {
			log.Fatal(err)
		}

		db, err := cgtdb.NewLDBDatabase(cfg.GainsDatabase)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		m := *model.NewModel(db, nil, cfg)

		log.Info("==============================================")
		log.Info("Processing the Transactions presented")
		m.ProcessModel()
		log.Info(fmt.Sprintf("length: %d \n", len(m.Accounts)))

		reporting.Printpdf(&m)
		reporting.ExportCSV(&m)
		reporting.ExportTraces(&m)

		return nil
	},
}

func init() {
	app = cli.NewApp()
	app.Action = wizard.WizardCommand
	app.Name = "CGT Calculator"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Sean Darcy",
			Email: "sean@darcyfinancial.com",
		},
	}
	app.Flags = []cli.Flag{}
	app.Version = version.Version
	app.Usage = "Calculates Capital Gains on cryptocurrency transactions"
	app.Commands = []cli.Command{
		config.NewCommand,
		loadCommand,
		prices.PricesCommand,
		debugCommand,
		config.DumpConfigCommand,
		account.DelCommand,
		account.DumpCommand,
		processCommand,
	}

	Formatter := new(log.TextFormatter)
	Formatter.TimestampFormat = "2006-01-02 15:04:05"
	Formatter.FullTimestamp = true
	log.SetFormatter(Formatter)
	app.Before = func(c *cli.Context) error {
		var filename string = "logfile.log"
		f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			fmt.Println(err)
		} else {
			log.SetOutput(f)
		}
		log.Info("========================")
		log.Info("========================")
		log.Info("========================")
		log.Info("Running CGT CALCULATOR")
		log.Info("Time: ", time.Now())
		log.Info("Command: ", os.Args)
		fmt.Println("Command: ", os.Args)
		return nil
	}

}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
