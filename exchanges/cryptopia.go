package exchanges

import (
	"encoding/csv"
	"os"
	"strconv"
	"strings"
	"time"

	"cgtcalc/model"

	log "github.com/Sirupsen/logrus"
)

type Cryptopia struct {
	Reference string
	Market    string
	Type      string
	Rate      float64
	Amount    float64
	Total     float64
	Fee       float64
	Timestamp string
}

func (txn *Cryptopia) ProcessData() (out model.Transaction) {
	//t, err := time.Parse("2/01/2006 15:04:05 PM",txn.Timestamp)
	t, err := time.Parse("2/01/2006 15:04", txn.Timestamp)
	if err != nil {
		log.Fatal(err)
	}
	out.Date = t
	out.ExchangeID = txn.Reference
	out.BaseCurrency, out.QuoteCurrency = model.MarketConverter(txn.Market, txn.Type)
	if strings.ToUpper(txn.Type) == "BUY" {
		out.BaseSpent = txn.Total
		out.QuoteReceived = txn.Amount
	} else {
		out.BaseSpent = txn.Amount
		out.QuoteReceived = txn.Total
	}

	out.Exchange = "Cryptopia"

	return
}

func CryptopiaFile(filename string) []model.Transaction {
	csvFile, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	lines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var out []model.Transaction

	for _, line := range lines[1:] {
		price, _ := strconv.ParseFloat(line[3], 64)
		amt, _ := strconv.ParseFloat(line[4], 64)
		total, _ := strconv.ParseFloat(line[5], 64)
		fee, _ := strconv.ParseFloat(line[6], 64)
		data := Cryptopia{
			Reference: line[0],
			Market:    line[1],
			Type:      line[2],
			Rate:      price,
			Amount:    amt,
			Total:     total,
			Fee:       fee,
			Timestamp: line[7],
		}

		log.Info(data)
		out = append(out, data.ProcessData())
	}

	return out

}
