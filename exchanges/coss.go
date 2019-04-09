package exchanges

import (
	"strings"
	"strconv"
	"os"
	"encoding/csv"
	"time"


	"cgtcalc/model"

	log "github.com/Sirupsen/logrus"
)

type Coss struct {
	Date 						string 				
	Market 					string 				
	BuySell 				string				
	Amount					float64				
	Price						float64				
	Total						float64				
	Fee							float64				
}

func (txn *Coss) ProcessData() (out model.Transaction){
	t, err := time.Parse("2006-01-02T15:04:05.999Z",txn.Date)
	if err != nil {
		log.Fatal(err)
	}
	out.Date = t
	out.ExchangeID = txn.Date + txn.Market + txn.BuySell
	out.BaseCurrency, out.QuoteCurrency = model.MarketConverter(txn.Market, txn.BuySell)
	if (strings.ToUpper(txn.BuySell)=="BUY") {
		out.BaseSpent = txn.Total + txn.Fee
		out.QuoteReceived = txn.Amount
	} else {
		out.BaseSpent = txn.Amount
		out.QuoteReceived = txn.Total - txn.Fee
	}
	out.Exchange = "Coss"

	return	
}

func CossFile(filename string) ([]model.Transaction) {
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
			f3, _ := strconv.ParseFloat(line[3],64)
			f4, _ := strconv.ParseFloat(line[4],64)
			f5, _ := strconv.ParseFloat(line[5],64)
			f6, _ := strconv.ParseFloat(line[6],64)
			log.Info("Processed the floats")
			data := Coss{
					Date: 								line[0],
					Market: 							line[1],
					BuySell: 							line[2],
					Amount:								f3,
					Price:								f4,
					Total:								f5,
					Fee:									f6,
			}

			log.Info(data)
			log.Info("Processed the data")
			out = append(out, data.ProcessData())
		}

		return out

}
