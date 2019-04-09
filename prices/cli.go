package prices

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strings"

	"cgtcalc/config"
	"cgtcalc/cgtdb"

	"github.com/leekchan/accounting"
	"github.com/urfave/cli"

	log "github.com/Sirupsen/logrus"
)

var PricesCommand = cli.Command{
      Name:        "prices",
      Usage:       "options for adding to prices database",
      Subcommands: []cli.Command{
				loadPriceCommand,
				loadPriceDirCommand,
				loadSinglePriceCommand,
				searchPriceCommand,
				printPriceCommand,
      },
}
var loadPriceCommand = cli.Command{
	Name:      "loadPriceFile",
	Usage:     "Adds a CSV of prices to the database",
	ArgsUsage: "First parameter is the CSV containting the prices, base will default to USD unless set with -base",
	Flags: []cli.Flag{
		config.ConfigFileFlag,
		cli.StringFlag{
			Name: "base",
			Value: "USD",
			Usage: "currency the prices are shown",
		},
		cli.StringFlag{
			Name: "currency",
			Value: "",
			Usage: "currency the prices are shown",
		},
		cli.StringFlag{
			Name: "dateFormat",
			Value: "2-Jan-06",
			Usage: "the format that the date is in, using Go's stupid format eg 2-Jan-06, 2006-01-02",
		},
		cli.IntFlag{
			Name: "dateColumn",
			Value: 1,
			Usage: "The Column that the date is in the CSV",
		},
		cli.IntFlag{
			Name: "priceColumn",
			Value: 5,
			Usage: "The Column that the price is in the CSV",
		},
	},
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			log.Fatal("This command requires an argument.")
		}
		if ctx.String("currency") == "" {
			log.Fatal("No Currency was specified for file")
		}
		fp := ctx.Args().First()
		cfg, err := config.LoadConfig(ctx.String("configFileFlag"))
		if err != nil {
			return err
		}

		db, err := cgtdb.NewLDBDatabase(cfg.PriceDatabase)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		p := *NewPricer(db)


		p.AddPriceFile(fp,strings.ToUpper(ctx.String("base")),strings.ToUpper(ctx.String("currency")),ctx.String("dateFormat"),ctx.Int("dateColumn"),ctx.Int("priceColumn"))

		return nil
	},
}
var loadPriceDirCommand = cli.Command{
	Name:      "loadPriceDir",
	Usage:     "Adds a directory of price CSVs to the database, Name reflects currency",
	ArgsUsage: "First parameter is the Directory containting the prices, base will default to USD unless set with -base",
	Flags: []cli.Flag{
		config.ConfigFileFlag,
		cli.StringFlag{
			Name: "base",
			Value: "USD",
			Usage: "currency the prices are shown",
		},
		cli.StringFlag{
			Name: "dateFormat",
			Value: "2-Jan-06",
			Usage: "the format that the date is in, using Go's stupid format eg 2-Jan-06, 2006-01-02",
		},
		cli.IntFlag{
			Name: "dateColumn",
			Value: 1,
			Usage: "The Column that the date is in the CSV",
		},
		cli.IntFlag{
			Name: "priceColumn",
			Value: 5,
			Usage: "The Column that the price is in the CSV",
		},
	},
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			log.Fatal("This command requires an argument.")
		}
		fp := ctx.Args().First()
		cfg, err := config.LoadConfig(ctx.String("configFileFlag"))
		if err != nil {
			return err
		}

		db, err := cgtdb.NewLDBDatabase(cfg.PriceDatabase)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		p := *NewPricer(db)

		p.AddPriceDir(fp,strings.ToUpper(ctx.String("base")),ctx.String("dateFormat"),ctx.Int("dateColumn"),ctx.Int("priceColumn"))

		return nil
	},
}

var searchPriceCommand = cli.Command{
	Name:      "search",
	Usage:     "searches the price database for the price at that date",
	ArgsUsage: "priceformat = yyyy-mm-dd",
	Flags: []cli.Flag{
		config.ConfigFileFlag,
		cli.StringFlag{
			Name: "currency",
			Value: "ETH",
			Usage: "currency being searched",
		},
		//cli.StringFlag{
			//Name: "base",
			//Value: "AUD",
			//Usage: "currency the prices are shown",
		//},
		//cli.StringFlag{
			//Name: "dateFormat",
			//Value: "2006-01-02",
			//Usage: "the format that the date is in, using Go's stupid format eg 2-Jan-06, 2006-01-02",
		//},
	},
	Action: func(ctx *cli.Context) error {

		if len(ctx.Args()) < 1 {
			log.Fatal("This command requires an argument.")
		}
		date := ctx.Args().First()

		cry := accounting.Accounting{Precision: 4}
		cfg, err := config.LoadConfig(ctx.String("configFileFlag"))
		if err != nil {
			return err
		}

		db, err := cgtdb.NewLDBDatabase(cfg.PriceDatabase)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		p := *NewPricer(db)

		price := p.GetPrice(date,strings.ToUpper(ctx.String("currency")))

		fmt.Printf("[%s] Currency: %s, Price: %s \n", date, strings.ToUpper(ctx.String("currency")), cry.FormatMoney(price))

		return nil
	},
}

var printPriceCommand = cli.Command{
	Name:      "print",
	Usage:     "Prints a list of all prices in the priceDB",
	ArgsUsage: "blah",
	Flags: []cli.Flag{
		config.ConfigFileFlag,
	},
	Action: func(ctx *cli.Context) error {

		cry := accounting.Accounting{Precision: 4}
		cfg, err := config.LoadConfig(ctx.String("configFileFlag"))
		if err != nil {
			return err
		}

		db, err := cgtdb.NewLDBDatabase(cfg.PriceDatabase)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		gob.Register(PriceItem{})

		iter := db.NewIterator()
		
		for iter.Next() {

			dbuf := bytes.NewBuffer(iter.Value())
			if err != nil {
				log.Fatal(err)
			} 
			var p PriceItem
			d := gob.NewDecoder(dbuf)
			if err := d.Decode(&p); err != nil {
				log.Fatal(err)
			}

			key := strings.Split(string(iter.Key()),"-")
			fmt.Printf("[%s] Currency: %s, Price: %s %s \n",p.Date, key[3], cry.FormatMoney(p.Amount), p.Base)
		}

		return nil
	},
}

var loadSinglePriceCommand = cli.Command{
	Name:      "addSinglePrice",
	Usage:     "Adds a single to the database",
	ArgsUsage: "First parameter is the amount prices, base will default to USD unless set with -base",
	Flags: []cli.Flag{
		config.ConfigFileFlag,
		cli.StringFlag{
			Name: "base",
			Value: "USD",
			Usage: "currency the prices are quoted in shown",
		},
		cli.StringFlag{
			Name: "currency",
			Value: "",
			Usage: "which currency to add",
		},
		cli.StringFlag{
			Name: "date",
			Value: "",
			Usage: "date (yyyy-mm-dd)",
		},
	},
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			log.Fatal("This command requires an argument.")
		}
		if ctx.String("currency") == "" {
			log.Fatal("No Currency was specified for file")
		}
		if ctx.String("date") == "" {
			log.Fatal("No Date was specified for file")
		}
		amount := ctx.Args().First()
		cfg, err := config.LoadConfig(ctx.String("configFileFlag"))
		if err != nil {
			return err
		}

		db, err := cgtdb.NewLDBDatabase(cfg.PriceDatabase)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		p := *NewPricer(db)


		p.AddSinglePrice(ctx.String("date"), strings.ToUpper(ctx.String("currency")),amount,strings.ToUpper(ctx.String("base")))

		return nil
	},
}
