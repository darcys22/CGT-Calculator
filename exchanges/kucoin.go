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

type KuCoin struct {
	Timestamp       string
	Market          string
	Type            string
	FilledPrice     float64
	FilledPriceCoin string
	Amount          float64
	AmountPriceCoin string
	Volume          float64
	VolumeCoin      string
	Fee             float64
	FeeCoin         string
}

func (txn *KuCoin) ProcessData() (out model.Transaction) {
	t, err := time.Parse("2006-01-02 15:04:05", txn.Timestamp)
	if err != nil {
		log.Fatal(err)
	}
	out.Date = t
	if strings.ToUpper(txn.Type) == "BUY" {
		out.QuoteCurrency = txn.AmountPriceCoin
		if txn.FeeCoin != out.QuoteCurrency {
			log.Fatal("Unexpected Fee in Wrong Currency for Transaction")
		}
		out.QuoteReceived = txn.Amount - txn.Fee
		out.BaseCurrency = txn.VolumeCoin
		out.BaseSpent = txn.Volume
	} else {
		out.BaseCurrency = txn.AmountPriceCoin
		out.QuoteCurrency = txn.VolumeCoin
		if txn.FeeCoin != out.QuoteCurrency {
			log.Fatal("Unexpected Fee in Wrong Currency for Transaction")
		}
		out.QuoteReceived = txn.Volume - txn.Fee
		out.BaseSpent = txn.Amount
	}

	out.Exchange = "KuCoin"

	return
}

func KuCoinFile(filename string) []model.Transaction {
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
		filledprice, _ := strconv.ParseFloat(line[3], 64)
		amt, _ := strconv.ParseFloat(line[5], 64)
		volume, _ := strconv.ParseFloat(line[7], 64)
		fee, _ := strconv.ParseFloat(line[9], 64)
		data := KuCoin{
			Timestamp:       line[0],
			Market:          line[1],
			Type:            line[2],
			FilledPrice:     filledprice,
			FilledPriceCoin: line[4],
			Amount:          amt,
			AmountPriceCoin: line[6],
			Volume:          volume,
			VolumeCoin:      line[8],
			Fee:             fee,
			FeeCoin:         line[10],
		}

		log.Info(data)
		out = append(out, data.ProcessData())
	}

	return out

}
