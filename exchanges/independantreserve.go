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

type IndependantReserve struct {
	SettledDateTimestamp 			string 				
	Timestamp 								string 				
	OrderGUID									string 				
	OrderType									string 				
	Status										string 				
	OpenClosed								string 				
	PrimaryCurrency						string 				
	Volume										float64				
	Outstanding								float64				
	SecondaryCurrency					string 				
	Price											float64				
	AveragePrice							float64				
	Value											float64				
	Brokerage									float64
	BrokerageFee							string
	GST												float64				
	TotalCost			 						float64				
}

func (txn *IndependantReserve) ProcessData() (out model.Transaction){
	t, err := time.Parse("2 Jan 2006 15:04",txn.SettledDateTimestamp)
	if err != nil {
		log.Fatal(err)
	}
	out.Date = t
	out.ExchangeID = txn.OrderGUID
	if (strings.Contains(strings.ToUpper(txn.OrderType),"BUY")) {
		out.BaseCurrency = txn.SecondaryCurrency
		out.BaseSpent = txn.TotalCost
		out.QuoteCurrency = txn.PrimaryCurrency
		out.QuoteReceived = txn.Volume
	} else {
		out.BaseCurrency = txn.PrimaryCurrency
		out.BaseSpent = txn.Volume
		out.QuoteCurrency = txn.SecondaryCurrency
		out.QuoteReceived = txn.TotalCost
	}

	out.Exchange = "Independant Reserve"

	return	
}

func IndependantReserveFile(filename string) ([]model.Transaction) {
		csvFile,err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer csvFile.Close()
		r := csv.NewReader(csvFile)

		//_, _ = r.Read()
		lines, err := r.ReadAll()

		if err != nil {
			log.Warn("Failing to read the CSV")
			log.Fatal(err)
		}

		var out []model.Transaction

		for _, line := range lines[1:] {
				f1, _ := strconv.ParseFloat(line[7],64)
				f2, _ := strconv.ParseFloat(line[8],64)
				f3, _ := strconv.ParseFloat(line[10],64)
				//f3 := 0.0
				f4, _ := strconv.ParseFloat(line[11],64)
				//f4 := 0.0
				f5, _ := strconv.ParseFloat(line[12],64)
				f6, _ := strconv.ParseFloat(line[13],64)
				f7, _ := strconv.ParseFloat(line[15],64)
				f8, _ := strconv.ParseFloat(line[16],64)
				data := IndependantReserve{
							SettledDateTimestamp: 			line[0],
							Timestamp: 									line[1],
							OrderGUID:									line[2],
							OrderType:									line[3],
							Status:											line[4],
							OpenClosed:									line[5],
							PrimaryCurrency:						line[6],
							Volume:											f1,	
							Outstanding:								f2,
							SecondaryCurrency:					line[9],
							Price:											f3,
							AveragePrice:								f4,
							Value:											f5,
							Brokerage:									f6,
							BrokerageFee:								line[14],
							GST:												f7,
							TotalCost:			 						f8,
				}

				log.Info(data)
				out = append(out, data.ProcessData())
		}

		return out

}
