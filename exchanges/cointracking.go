package exchanges

import (
	"strings"
	"strconv"
	"os"
	"encoding/csv"
	"time"


	"cgtcalc/model"

	"github.com/leekchan/accounting"

	log "github.com/Sirupsen/logrus"
)

type Cointracking struct {
	Type						string				
	BuyAmount				float64				
	BuyCurrency			string 				
	SellAmount			float64				
	SellCurrency		string				
	Fee							float64
	FeeCurrency			string
	Exchange				string
	Group						string
	Comment					string
	Timestamp 			string 				
}

func (txn *Cointracking) ProcessData() (out model.Transaction){
	cry := accounting.Accounting{Precision: 4}
	t, err := time.Parse("2006-01-02 15:04:05",txn.Timestamp)
	if err != nil {
		log.Fatal("Could not Parse the Time ", err)
	}
	out.Date = t
	out.Exchange = txn.Exchange
	out.ExchangeID = txn.Timestamp + "b" + cry.FormatMoney(txn.BuyAmount) + txn.BuyCurrency + "s" + cry.FormatMoney(txn.SellAmount) + txn.SellCurrency
	if (strings.ToUpper(txn.Type)=="TRADE") {
		out.BaseCurrency = txn.SellCurrency
		out.BaseSpent = txn.SellAmount
		out.QuoteCurrency = txn.BuyCurrency
		out.QuoteReceived = txn.BuyAmount
		if txn.FeeCurrency == txn.SellCurrency {
			out.BaseSpent += txn.Fee
		} else if txn.FeeCurrency == txn.BuyCurrency {
			out.QuoteReceived -= txn.Fee
		}
	} else if (strings.ToUpper(txn.Type)=="LOST") {
		out.BaseCurrency = txn.SellCurrency
		out.BaseSpent = txn.SellAmount
		out.QuoteCurrency = "AUD"
		out.QuoteReceived = *new(float64)

	} else if (strings.ToUpper(txn.Type)=="GIFT/TIP") {
		out.BaseCurrency = "BTC"
		out.BaseSpent = *new(float64)
		out.QuoteCurrency = txn.BuyCurrency
		out.QuoteReceived = txn.BuyAmount
	} else {
		out = model.Transaction{}
	}

	return	
}

func CointrackingFile(filename string) ([]model.Transaction) {
		log.Info("Opening File")
		csvFile,err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer csvFile.Close()

		log.Info("Reading File")
		lines, err := csv.NewReader(csvFile).ReadAll()
		if err != nil {
			log.Fatal(err)
		}

		var out []model.Transaction

		log.Info("Looping over the lines")
		log.Info("Total Lines: ", len(lines[1:]))
		for idx, line := range lines[1:] {
			log.Info("Processing line: ", idx)
			f1, _ := strconv.ParseFloat(line[1],64)
			f3, _ := strconv.ParseFloat(line[3],64)
			f5, _ := strconv.ParseFloat(line[5],64)
			log.Info("Processed the floats")
			data := Cointracking{
					Type: 						line[0],
					BuyAmount:				f1,
					BuyCurrency: 			line[2],
					SellAmount:				f3,
					SellCurrency:			line[4],
					Fee:							f5,
					FeeCurrency:			line[6],
					Exchange:					line[7],
					Group: 						line[8],
					Comment:					line[9],
					Timestamp: 				line[10],
			}

			log.Info(data)
			log.Info("Processed the data")
			if data.Type == "Trade" || data.Type == "Lost" || data.Type == "Gift/Tip" {
				out = append(out, data.ProcessData())
			}
		}

		return out

}
