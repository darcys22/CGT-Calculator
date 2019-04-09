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

type Coinspot struct {
	Timestamp 			string 				
	Type						string				
	Market 					string 				
	Amount					float64				
	RateIncFee			float64				
	RateExclFee			float64				
	Fee							string
	FeeAUDInclGST		float64				
	GSTAud					float64				
	TotalAud				float64				
	TotalAudInclGST	string
}

func (txn *Coinspot) ProcessData() (out model.Transaction){
	t, err := time.Parse("2006-01-02T15:04:05.999Z",txn.Timestamp)
	if err != nil {
		log.Fatal(err)
	}
	out.Date = t
	out.BaseCurrency, out.QuoteCurrency = model.MarketConverter(txn.Market, txn.Type)
	if (strings.ToUpper(txn.Type)=="BUY") {
		//TODO: I have not actually seen a coinspot "BUY" this is a guess
		log.Warn("BUY ORDERS ARE UNCERTAIN, Review before this gets finalised")
		out.QuoteReceived = txn.Amount
		out.BaseSpent = txn.Amount / txn.RateIncFee
	} else {
		out.BaseSpent = txn.Amount
		out.QuoteReceived = txn.Amount * txn.RateIncFee 
	}

	out.Exchange = "Coinspot"

	return	
}

func CoinspotFile(filename string) ([]model.Transaction) {
		log.Warn("Only Sell orders have been made for Coinspot, need to analyse BUY")
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
			f1, _ := strconv.ParseFloat(line[3],64)
			f2, _ := strconv.ParseFloat(line[4],64)
			f3, _ := strconv.ParseFloat(line[5],64)
			f4, _ := strconv.ParseFloat(line[7],64)
			f5, _ := strconv.ParseFloat(line[8],64)
			f6, _ := strconv.ParseFloat(line[9],64)
			log.Info("Processed the floats")
			data := Coinspot{
					Timestamp: 		line[0],
					Type: 				line[1],
					Market: 			line[2],
					Amount:				f1,
					RateIncFee:		f2,
					RateExclFee:	f3,
					Fee:					line[6],
					FeeAUDInclGST: f4,
					GSTAud:				f5,
					TotalAud:			f6,
					TotalAudInclGST: line[10],
			}

			log.Info(data)
			log.Info("Processed the data")
			out = append(out, data.ProcessData())
		}

		return out

}
