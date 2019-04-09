package exchanges

import (
	"strings"
	"strconv"
	"os"
	"encoding/csv"
	"time"


	"cgtcalc/model"

	"golang.org/x/text/encoding/unicode"

	log "github.com/Sirupsen/logrus"
)

type Bittrex struct {
	OrderUuid				string
	Exchange				string 				
	Type						string				
	Quantity				float64				
	Limit						float64
	CommissionPaid	float64				
	Price						float64				
	Opened					string 				
	Timestamp 			string 				
}

func (txn *Bittrex) ProcessData() (out model.Transaction){
	t, err := time.Parse("1/2/2006 3:04:05 PM",txn.Timestamp)
	if err != nil {
		log.Fatal(err)
	}
	out.Date = t
	typ := strings.Split(txn.Type,"_")
	out.BaseCurrency, out.QuoteCurrency = model.MarketConverter(txn.Exchange, typ[len(typ)-1])
	out.ExchangeID = txn.OrderUuid
	if (typ[len(typ)-1])=="BUY" {
		out.QuoteReceived = txn.Quantity
		out.BaseSpent = txn.Price + txn.CommissionPaid
	} else {
		out.BaseSpent = txn.Quantity
		out.QuoteReceived = txn.Price - txn.CommissionPaid
	}

	out.Exchange = "Bittrex"

	return	
}

func BittrexFile(filename string) ([]model.Transaction) {
		log.Info("Opening File")

		csvFile,err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer csvFile.Close()

		dec := unicode.UTF16(unicode.LittleEndian,0).NewDecoder()
		reader := dec.Reader(csvFile)

		log.Info("Reading File")
		lines, err := csv.NewReader(reader).ReadAll()
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
			data := Bittrex{
					OrderUuid: 					line[0],
					Exchange: 					line[1],
					Type: 							line[2],
					Quantity: 					f3,
					Limit:							f4,
					CommissionPaid:			f5,
					Price:							f6,
					Opened:							line[7],
					Timestamp:					line[8],
			}

			log.Info(data)
			log.Info("Processed the data")
			out = append(out, data.ProcessData())
		}

		return out

}
