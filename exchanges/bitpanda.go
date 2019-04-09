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

type Bitpanda struct {
	TransactionID		string
	Type						string				
	Category				string				
	AmountFiat			float64				
	Fee							float64				
	FiatCurrency		string 				
	AmountCrypto		float64				
	CryptoCoin			string 				
	Status					string 				
	Timestamp 			string 				
}

func (txn *Bitpanda) ProcessData() (out model.Transaction){
	//t, err := time.Parse("2006-01-02 15:04:05",txn.Timestamp)
	t, err := time.Parse(time.RFC3339,txn.Timestamp)
	if err != nil {
		log.Fatal(err)
	}
	out.Date = t
	out.ExchangeID = txn.TransactionID
	out.Exchange = "Bitpanda"
	if (strings.ToUpper(txn.Type)=="BUY") {
		out.BaseCurrency = txn.FiatCurrency
		out.BaseSpent = txn.AmountFiat
		out.QuoteCurrency = txn.CryptoCoin
		out.QuoteReceived = txn.AmountCrypto
	} else if (strings.ToUpper(txn.Type)=="DEPOSIT") {
		out.BaseCurrency = txn.CryptoCoin
		out.BaseSpent = txn.Fee
		out.QuoteCurrency = "BTC"
		out.QuoteReceived = *new(float64)
	} else {
		out.BaseCurrency = txn.CryptoCoin
		out.BaseSpent = txn.AmountCrypto
		out.QuoteCurrency = txn.FiatCurrency
		out.QuoteReceived = txn.AmountFiat
	}


	return	
}

func BitpandaFile(filename string) ([]model.Transaction) {
		log.Info("Opening File")
		csvFile,err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer csvFile.Close()

		log.Info("Reading File")
		r := csv.NewReader(csvFile)
		_, _ = r.Read()
		_, _ = r.Read()
		_, _ = r.Read()
		_, _ = r.Read()
		_, _ = r.Read()
		lines, err := r.ReadAll()
		if err != nil {
			log.Fatal("failure to read the CSV File: ",err)
		}

		var out []model.Transaction

		log.Info("Looping over the lines")
		log.Info("Total Lines: ", len(lines[4:]))
		for idx, line := range lines[1:] {
			log.Info("Processing line: ", idx)
			f3, _ := strconv.ParseFloat(line[3],64)
			f4, _ := strconv.ParseFloat(line[4],64)
			f6, _ := strconv.ParseFloat(line[6],64)
			log.Info("Processed the floats")
			data := Bitpanda{
				TransactionID:						line[0],
				Type:											line[1],
				Category:									line[2],
				AmountFiat:								f3,
				Fee:											f4,
				FiatCurrency:							line[5],
				AmountCrypto:							f6,
				CryptoCoin:								line[7],
				Status:										line[8],
				Timestamp:								line[9],
			}

			log.Info(data)
			log.Info("Processed the data")
			if data.Type == "buy" || data.Type == "sell" || data.Type == "deposit" {
				out = append(out, data.ProcessData())
			}
		}

		return out

}
