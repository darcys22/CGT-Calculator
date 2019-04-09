package exchanges

import (
	"encoding/csv"
	"os"
	"strconv"
	"time"
	"strings"

	"cgtcalc/model"

	log "github.com/Sirupsen/logrus"
)

type Custom struct {
	Timestamp  string
	BasePrice  float64
	BaseCoin   string
	QuotePrice float64
	QuoteCoin  string
	Exchange   string
}

func (txn *Custom) ProcessData() (out model.Transaction) {
	t, err := time.Parse("2006-01-02", txn.Timestamp)
	//t, err := time.Parse("2006-01-2 15:04:05", txn.Timestamp)
	if err != nil {
		log.Fatal(err)
	}
	out.Date = t

	out.QuoteCurrency = strings.ToUpper(txn.QuoteCoin)
	out.QuoteReceived = txn.QuotePrice
	out.BaseCurrency = strings.ToUpper(txn.BaseCoin)
	out.BaseSpent = txn.BasePrice

	out.Exchange = txn.Exchange

	return
}

func CustomFile(filename string) []model.Transaction {
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
		baseprice, _ := strconv.ParseFloat(line[1], 64)
		quoteprice, _ := strconv.ParseFloat(line[3], 64)
		data := Custom{
			Timestamp:  line[0],
			BasePrice:  baseprice,
			BaseCoin:   line[2],
			QuotePrice: quoteprice,
			QuoteCoin:  line[4],
			Exchange:  	line[5],
		}

		log.Info(data)
		out = append(out, data.ProcessData())
	}

	return out

}
