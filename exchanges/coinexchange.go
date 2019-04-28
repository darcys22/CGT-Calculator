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

type Coinexchange struct {
	TradeID       string
	Type          string
	Related_order string
	Market        string
	Time          string
	Price         float64
	Amount        float64
	Total         float64
	Fee           float64
	Net_total     float64
}

func (txn *Coinexchange) ProcessData() (out model.Transaction) {
	t, err := time.Parse("2/01/2006 15:04", txn.Time)
	if err != nil {
		log.Fatal(err)
	}
	out.Date = t
	out.BaseCurrency, out.QuoteCurrency = model.MarketConverter(txn.Market, txn.Type)
	out.ExchangeID = txn.Related_order
	if strings.ToUpper(txn.Type) == "BUY" {
		out.QuoteReceived = txn.Amount
		out.BaseSpent = txn.Net_total
	} else {
		out.BaseSpent = txn.Amount
		out.QuoteReceived = txn.Net_total
	}

	out.Exchange = "Coinexchange"

	return
}

func CoinexchangeFile(filename string) []model.Transaction {
	log.Info("Opening File")
	csvFile, err := os.Open(filename)
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
		f5, _ := strconv.ParseFloat(line[5], 64)
		f6, _ := strconv.ParseFloat(line[6], 64)
		f7, _ := strconv.ParseFloat(line[7], 64)
		f8, _ := strconv.ParseFloat(line[8], 64)
		f9, _ := strconv.ParseFloat(line[9], 64)
		log.Info("Processed the floats")
		data := Coinexchange{
			TradeID:       line[0],
			Type:          line[1],
			Related_order: line[2],
			Market:        line[3],
			Time:          line[4],
			Price:         f5,
			Amount:        f6,
			Total:         f7,
			Fee:           f8,
			Net_total:     f9,
		}

		log.Info(data)
		log.Info("Processed the data")
		out = append(out, data.ProcessData())
	}

	return out

}
