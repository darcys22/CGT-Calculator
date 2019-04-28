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

type Mercatox struct {
	MXTransactionID string
	NTTransactionID string
	WithdrawAddr    string
	Type            string
	Currency        string
	Pair            string
	Fee             string
	Amount          float64
	Price           float64
	Total           float64
	Action          string
	From            string
	To              string
	Time            string
}

func (txn *Mercatox) ProcessData() (out model.Transaction) {
	t, err := time.Parse("Jan 2, 2006, 15:04:05 PM", txn.Time)
	if err != nil {
		log.Fatal(err)
	}
	out.Date = t
	out.BaseCurrency, out.QuoteCurrency = model.MarketConverter(txn.Pair, txn.Action)
	out.ExchangeID = txn.MXTransactionID
	if strings.ToUpper(txn.Action) == "BUY" {
		out.QuoteReceived = txn.Amount
		out.BaseSpent = txn.Total
	} else {
		log.Warn("Sell side not full implemented, finish this")
		out.BaseSpent = txn.Amount
		out.QuoteReceived = txn.Total
	}

	out.Exchange = "Mercatox"

	return
}

func MercatoxFile(filename string) []model.Transaction {
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
		f7, _ := strconv.ParseFloat(line[7], 64)
		f8, _ := strconv.ParseFloat(line[8], 64)
		f9, _ := strconv.ParseFloat(line[9], 64)
		log.Info("Processed the floats")
		data := Mercatox{
			MXTransactionID: line[0],
			NTTransactionID: line[1],
			WithdrawAddr:    line[2],
			Type:            line[3],
			Currency:        line[4],
			Pair:            line[5],
			Fee:             line[6],
			Amount:          f7,
			Price:           f8,
			Total:           f9,
			Action:          line[10],
			From:            line[11],
			To:              line[12],
			Time:            line[13],
		}

		log.Info(data)
		log.Info("Processed the data")
		if data.Type == "Deal" {
			out = append(out, data.ProcessData())
		}
	}

	return out

}
