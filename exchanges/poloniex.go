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

type Poloniex struct {
	Timestamp 			string 				
	Market 					string 				
	Category				string				
	Type						string				
	Price						float64				
	Amount					float64				
	Total						float64				
	Fee							string
	OrderNumber			string
	BaseTotalLessFee	float64				
	QuoteTotalLessFee	float64				
}

func (txn *Poloniex) ProcessData() (out model.Transaction){
	t, err := time.Parse("2006-01-02 15:04:05",txn.Timestamp)
	if err != nil {
		log.Fatal(err)
	}
	out.Date = t
	out.BaseCurrency, out.QuoteCurrency = model.MarketConverter(txn.Market, txn.Type)
	out.ExchangeID = txn.Timestamp + txn.OrderNumber
	if (strings.ToUpper(txn.Type)=="BUY") {
		out.QuoteReceived = txn.QuoteTotalLessFee
		out.BaseSpent = txn.Total
	} else {
		out.BaseSpent = txn.Amount
		out.QuoteReceived = txn.BaseTotalLessFee
	}

	out.Exchange = "Poloniex"

	return	
}

func PoloniexFile(filename string) ([]model.Transaction) {
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
			f1, _ := strconv.ParseFloat(line[4],64)
			f2, _ := strconv.ParseFloat(line[5],64)
			f3, _ := strconv.ParseFloat(line[6],64)
			f4, _ := strconv.ParseFloat(line[9],64)
			f5, _ := strconv.ParseFloat(line[10],64)
			log.Info("Processed the floats")
			data := Poloniex{
					Timestamp: 					line[0],
					Market: 						line[1],
					Category:						line[2],
					Type: 							line[3],
					Price:							f1,
					Amount:							f2,
					Total:							f3,
					Fee:								line[7],
					OrderNumber:				line[8],
					BaseTotalLessFee: 	f4,
					QuoteTotalLessFee: 	f5,
			}

			log.Info(data)
			log.Info("Processed the data")
			out = append(out, data.ProcessData())
		}

		return out

}
