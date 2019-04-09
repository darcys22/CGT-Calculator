package prices

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"strconv"
	"time"

	"cgtcalc/cgtdb"

	log "github.com/Sirupsen/logrus"
)

type PriceItem struct {
	Currency		string
	Date				string
	Base				string
	Amount			float64
}

type Pricer struct {
	DB					cgtdb.Database
}

func NewPricer(db cgtdb.Database) *Pricer {
	return &Pricer{
		DB:				db,
	}
}

func (p *Pricer) Close() {
	log.Info("Closing PricesDB")
	p.DB.Close()
}

func (p *Pricer) AddPriceFile(filename, base, currency, dateFormat string, dateColumn, priceColumn int) {

		log.Info("===================================")
		log.Info("Adding a Price File to the Database")
		log.Info("File: %s",filename)
		log.Info("Currency: %s",currency)
		log.Info("Base: %s",base)
		log.Info("dateFormat: %s",dateFormat)
		log.Info("dateColumn: %d",dateColumn)
		log.Info("priceColumn: %d",priceColumn)

		f := strings.TrimSuffix(filename, filepath.Ext(filename))
		_, name := path.Split(f)
		name = strings.TrimSpace(currency)

		gob.Register(PriceItem{})

		csvFile , _ := os.Open(filename)
		reader := csv.NewReader(bufio.NewReader(csvFile))
		reader.Read()

		if (dateFormat=="") {
			dateFormat = "2-Jan-06"
		}
		if (dateColumn==0) {
			dateColumn = 1
		}
		if (priceColumn==0) {
			priceColumn = 5
		}

		for {
			line, err := reader.Read()

			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}
			s := strings.Replace(line[priceColumn-1], "$", "", -1)
			fl, err := strconv.ParseFloat(strings.TrimSpace(strings.Replace(s, ",", "", -1)), 64)
			if err != nil {
				panic(err)
			}

			t, err := time.Parse(dateFormat, line[dateColumn-1])
			if err != nil {
				fmt.Println(err)
			}

			price := PriceItem {
				Currency: strings.ToUpper(name),
				Date: t.Format("2006-01-02"),
				Base: strings.ToUpper(base),
				Amount: fl,
			}

			buf2 := bytes.NewBuffer([]byte{})
			enc := gob.NewEncoder(buf2)
			err = enc.Encode(price)
			if err != nil {
				log.Fatal(err)
			}

			err = p.DB.Put([]byte(price.Date+"-"+price.Currency),buf2.Bytes())
			if err != nil {
				log.Fatal(err)
			}

			log.Info("-----------------------------------------")
			log.Info("Adding to the prices Database")
			log.Info("Name: ", price.Currency)
			log.Info("Date: ", price.Date)
			log.Info("Base: ", price.Base)
			log.Info("Price: ", price.Amount)
			log.Info("Key: ", price.Date+"-"+price.Currency)
		}
}

func (p *Pricer) AddPriceDir(dirname, base, dateFormat string, dateColumn, priceColumn int) {
		log.Info("Walking the Directory for Pricefiles: ",dirname)
		err := filepath.Walk("./"+dirname, func(filename string, info os.FileInfo, err error) error {
				if !info.IsDir() {
					log.Info("Adding to the Price Database: ",filename)
					f := strings.TrimSuffix(filename, filepath.Ext(filename))
					_, name := path.Split(f)
					currency := strings.ToUpper(name)
					p.AddPriceFile(filename, base, currency, dateFormat, dateColumn, priceColumn)
				}
				return nil
		})
		if err != nil {
				panic(err)
		}
}

func (pr *Pricer) GetPrice(date, currency string) (float64){
		log.Info("=====================================")
		log.Info("Searching for the price of a currency (prices.GetPrice)")
		log.Info(fmt.Sprintf("%s - %s\n",date,currency))
		log.Info(fmt.Sprintf("%s\n",hex.Dump([]byte(date+"-"+currency))))

		gob.Register(PriceItem{})

		if strings.ToUpper(currency) == "AUD" {
			return 1.0
		}
		if strings.ToUpper(currency) == "XBT" {
			currency = "BTC"
		}

		var priceItems []PriceItem

		BaseAUD := false
		for !BaseAUD {
			var p PriceItem
			data, err := pr.DB.Get([]byte(date+"-"+currency))
			if err != nil {
				resp, err := PriceSearch(currency, date) 
				log.Info("Currency Was Not Found in the Database searching on CryptoCompare")
				if err != nil {
					log.Warn("Currency Was Not Found in the Database or in CryptoCompare")
					log.Warn("Currency: ", currency)
					log.Warn("Date: ", date)
					log.Fatal(err)
				}
				p = *resp
				pr.AddSinglePrice(p.Date, p.Currency,strconv.FormatFloat(p.Amount, 'f', -1, 64), p.Base)
			} else {
				buf := bytes.NewBuffer(data)
				d := gob.NewDecoder(buf)
				if err := d.Decode(&p); err != nil {
					log.Warn("Could not convert database entry to model")
					log.Fatal(err)
				}
			}

			priceItems = append(priceItems, p)
			currency = p.Base
			BaseAUD = (strings.ToUpper(p.Base) == "AUD")

			log.Info(fmt.Sprintf("%s - %.2f\n",p.Base, p.Amount))
		}

		ret := priceItems[0].Amount
		for _, item := range priceItems[1:] {
			ret *= item.Amount
		}

		

	return ret 
}

func (p *Pricer) AddSinglePrice(date, currency, amount, base string) {


		s := strings.Replace(amount, "$", "", -1)
		fl, err := strconv.ParseFloat(strings.TrimSpace(strings.Replace(s, ",", "", -1)), 64)
		if err != nil {
			panic(err)
		}

		layout := "2006-01-02"
		t, err := time.Parse(layout, date)
		if err != nil {
			fmt.Println(err)
		}

		price := PriceItem {
			Currency: strings.ToUpper(currency),
			Date: t.Format("2006-01-02"),
			Base: strings.ToUpper(base),
			Amount: fl,
		}

		buf2 := bytes.NewBuffer([]byte{})
		enc := gob.NewEncoder(buf2)
		err = enc.Encode(price)
		if err != nil {
			log.Fatal(err)
		}

		err = p.DB.Put([]byte(price.Date+"-"+price.Currency),buf2.Bytes())
		if err != nil {
			log.Fatal(err)
		}

		log.Info("-----------------------------------------")
		log.Info("Adding to the prices Database")
		log.Info("Name: ", price.Currency)
		log.Info("Date: ", price.Date)
		log.Info("Base: ", price.Base)
		log.Info("Price: ", price.Amount)
		log.Info("Key: ", price.Date+"-"+price.Currency)
}
