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

type EtherDelta struct {
	Date 						string 				
	Action					string 				
	Source					string 				
	Volume					float64				
	Symbol					string 				
	Price						float64				
	Currency				string				
	Fee							float64				
	FeeCurrency			string				
	Memo						string				
}

func (txn *EtherDelta) ProcessData() (out model.Transaction){
	t, err := time.Parse("2006-01-02T15:04:05-07:00",txn.Date)
	if err != nil {
		log.Fatal(err)
	}
	out.Date = t
	out.ExchangeID = txn.Date + txn.Action + txn.Symbol + txn.Currency
	if (strings.ToUpper(txn.Action)=="BUY") {
		out.BaseCurrency = txn.Currency
		out.QuoteCurrency = txn.Symbol
		out.BaseSpent = txn.Price * txn.Volume + txn.Fee
		out.QuoteReceived = txn.Volume
	} else {
		out.BaseCurrency = txn.Symbol
		out.QuoteCurrency = txn.Currency
		out.BaseSpent = txn.Volume + txn.Fee
		out.QuoteReceived = txn.Volume * txn.Price
	}
	out.Exchange = "EtherDelta"

	return	
}

func EtherDeltaFile(filename string) ([]model.Transaction) {
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
			f5, _ := strconv.ParseFloat(line[5],64)
			f7, _ := strconv.ParseFloat(line[7],64)
			log.Info("Processed the floats")
			data := EtherDelta{
					Date:									line[0],
					Action:								line[1],
					Source:								line[2],
					Volume:								f3,
					Symbol:								line[4],
					Price:								f5,
					Currency:							line[6],
					Fee:									f7,
					FeeCurrency:					line[8],
					Memo:									line[9],
			}

			log.Info(data)
			log.Info("Processed the data")
			out = append(out, data.ProcessData())
		}

		return out

}
