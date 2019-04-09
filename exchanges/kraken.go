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

type Kraken struct {
	Txid						string
	OrderTxid				string
	Pair						string 				
	Timestamp 			string 				
	Type						string				
	Ordertype				string 				
	Price						float64				
	Cost						float64				
	Fee							float64
	Volume					float64				
}

func krakenConverter(market, typ string) (string, string){
	var x,y string

	if (strings.ToUpper(market)[len(market)-1] == 'T') {
		x = strings.ToUpper(market)[:len(market)-4]
		y = strings.ToUpper(market)[len(market)-4:]
		if (strings.ToUpper(y)[0] == strings.ToUpper(y)[1]) {
			y = strings.ToUpper(y)[1:]
		}
	} else {
		x = strings.ToUpper(market)[:len(market)-3]
		y = strings.ToUpper(market)[len(market)-3:]
	}

	if (strings.ToUpper(x)[len(x)-1] == 'X') {
		x = strings.ToUpper(x)[:len(x)-1]
	}
	if len(x) == 4 || (strings.ToUpper(x)[0] == 'X') {
		x = strings.ToUpper(x)[1:]
	}
	if (strings.ToUpper(typ)=="BUY") {
		return y,x
	} else {
		return x,y
	}
}

func (txn *Kraken) ProcessData() (out model.Transaction){
	t, err := time.Parse("2006-01-02 15:04:05.9999",txn.Timestamp)
	if err != nil {
		log.Fatal(err)
	}
	out.Date = t
	out.BaseCurrency, out.QuoteCurrency = krakenConverter(txn.Pair, txn.Type)
	out.ExchangeID = txn.Txid
	if (strings.ToUpper(txn.Type)=="BUY") {
		out.QuoteReceived = txn.Volume
		out.BaseSpent = txn.Cost + txn.Fee
	} else {
		out.BaseSpent = txn.Volume
		out.QuoteReceived = txn.Cost - txn.Fee
	}

	out.Exchange = "Kraken"

	return	
}

func KrakenFile(filename string) ([]model.Transaction) {
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
			f1, _ := strconv.ParseFloat(line[6],64)
			f2, _ := strconv.ParseFloat(line[7],64)
			f3, _ := strconv.ParseFloat(line[8],64)
			f4, _ := strconv.ParseFloat(line[9],64)
			log.Info("Processed the floats")
			data := Kraken{
					Txid: 							line[0],
					OrderTxid: 					line[1],
					Pair: 							line[2],
					Timestamp: 					line[3],
					Type:								line[4],
					Ordertype:					line[5],
					Price:							f1,
					Cost:								f2,
					Fee:								f3,	
					Volume:							f4,
			}

			log.Info(data)
			log.Info("Processed the data")
			out = append(out, data.ProcessData())
		}

		return out

}
