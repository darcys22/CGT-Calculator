package exchanges

import (
	"encoding/csv"
	"os"
	"strconv"
	"time"

	"cgtcalc/model"

	log "github.com/Sirupsen/logrus"
)

type Hitbtc struct {
	Date       string
	Instrument string
	TradeID    string
	OrderID    string
	Side       string
	Quantity   float64
	Price      float64
	Volume     float64
	Fee        float64
	Rebate     float64
	Total      float64
}

func (txn *Hitbtc) ProcessData() (out model.Transaction) {
	t, err := time.Parse("02/01/2006 15:04", txn.Date)
	if err != nil {
		log.Fatal(err)
	}
	out.Date = t
	out.BaseCurrency, out.QuoteCurrency = model.MarketConverter(txn.Instrument, txn.Side)
	out.ExchangeID = txn.TradeID
	out.QuoteReceived = txn.Quantity
	out.BaseSpent = txn.Volume

	out.Exchange = "Hitbtc"

	return
}

func HitbtcFile(filename string) []model.Transaction {
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
		f10, _ := strconv.ParseFloat(line[10], 64)
		log.Info("Processed the floats")
		data := Hitbtc{
			Date:       line[0],
			Instrument: line[1],
			TradeID:    line[2],
			OrderID:    line[3],
			Side:       line[4],
			Quantity:   f5,
			Price:      f6,
			Volume:     f7,
			Fee:        f8,
			Rebate:     f9,
			Total:      f10,
		}

		log.Info(data)
		log.Info("Processed the data")
		out = append(out, data.ProcessData())
	}

	return out

}
