package exchanges

import (
	"encoding/csv"
	"os"
	"strconv"
	"time"

	"cgtcalc/model"

	log "github.com/Sirupsen/logrus"
)

type TempKuCoin struct {
	Market    string
	DealPrice float64
	DealValue float64
	Amount    float64
	Fee       float64
	Type      string
	Timestamp string
}

func (txn *TempKuCoin) ProcessData() (out model.Transaction) {
	t, err := time.Parse("2/01/2006", txn.Timestamp)
	if err != nil {
		log.Fatal(err)
	}
	out.Date = t
	out.BaseCurrency, out.QuoteCurrency = model.InvMarketConverter(txn.Market, txn.Type)
	out.ExchangeID = txn.Timestamp + txn.Market
	out.QuoteReceived = txn.Amount
	out.BaseSpent = txn.DealValue
	//if strings.ToUpper(txn.Type) == "BUY" {
	//} else {
	//out.QuoteReceived = txn.Amount
	//out.BaseSpent = txn.DealValue
	//}

	out.Exchange = "KuCoin"

	return
}

func TempKuCoinFile(filename string) []model.Transaction {
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
		f1, _ := strconv.ParseFloat(line[1], 64)
		f2, _ := strconv.ParseFloat(line[2], 64)
		f3, _ := strconv.ParseFloat(line[3], 64)
		f4, _ := strconv.ParseFloat(line[4], 64)
		data := TempKuCoin{
			Market:    line[0],
			DealPrice: f1,
			DealValue: f2,
			Amount:    f3,
			Fee:       f4,
			Type:      line[5],
			Timestamp: line[6],
		}

		log.Info(data)
		out = append(out, data.ProcessData())
	}

	return out

}
