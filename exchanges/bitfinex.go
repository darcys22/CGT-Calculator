package exchanges

import (
	"strconv"
	"os"
	"encoding/csv"
	"time"


	"cgtcalc/model"

	log "github.com/Sirupsen/logrus"
)

type Bitfinex struct {
	TransactionID		string
	Pair						string 				
	Amount					float64				
	Price						float64				
	Fee							float64
	FeeCurrency			string				
	Timestamp 			string 				
	OrderID					string
}

func (txn *Bitfinex) ProcessData() (out model.Transaction){
	t, err := time.Parse("02-01-06 15:04:05",txn.Timestamp)
	if err != nil {
		log.Fatal(err)
	}
	out.Date = t
	buy := txn.Amount >= 0
	typ := "BUY"

	if buy {
		out.QuoteReceived = txn.Amount
		out.BaseSpent = txn.Amount * txn.Price - txn.Fee

	} else {
		typ = "SELL"
		out.BaseSpent = -txn.Amount - txn.Fee
		out.QuoteReceived = -txn.Amount * txn.Price
	}

	out.ExchangeID = txn.TransactionID
	out.BaseCurrency, out.QuoteCurrency = model.MarketConverter(txn.Pair, typ)

	out.Exchange = "Bitfinex"

	return	
}

func BitfinexFile(filename string) ([]model.Transaction) {
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
			f2, _ := strconv.ParseFloat(line[2],64)
			f3, _ := strconv.ParseFloat(line[3],64)
			f4, _ := strconv.ParseFloat(line[4],64)
			log.Info("Processed the floats")
			data := Bitfinex{
					TransactionID: 			line[0],
					Pair:								line[1],
					Amount: 						f2,
					Price:							f3,
					Fee:								f4,
					FeeCurrency:				line[5],
					Timestamp:					line[6],
					OrderID:						line[7],
			}

			log.Info(data)
			log.Info("Processed the data")
			out = append(out, data.ProcessData())
		}

		return out

}
