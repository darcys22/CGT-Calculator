package exchanges

import (
	"strings"

	"cgtcalc/model"

	"github.com/tealeg/xlsx"
	"github.com/araddon/dateparse"
	log "github.com/Sirupsen/logrus"
)

type Binance struct {
	Date 						string 				`xlsx:"0"`
	Market 					string 				`xlsx:"1"`
	BuySell 				string				`xlsx:"2"`
	Price						float64				`xlsx:"3"`
	Amount					float64				`xlsx:"4"`
	Total						float64				`xlsx:"5"`
	Fee							float64				`xlsx:"6"`
	Feecur 					string 				`xlsx:"7"`
}

func (txn *Binance) ProcessData() (out model.Transaction){
	t, err := dateparse.ParseAny(txn.Date)
	if err != nil {
		log.Fatal(err)
	}
	out.Date = t
	out.BaseCurrency, out.QuoteCurrency = model.MarketConverter(txn.Market, txn.BuySell)
	out.Exchange = "Binance"
	if (strings.ToUpper(txn.BuySell)=="BUY") {
		out.BaseSpent = txn.Total
		out.QuoteReceived = txn.Amount - txn.Fee
	} else {
		out.BaseSpent = txn.Amount
		out.QuoteReceived = txn.Total - txn.Fee
	}

	return	
}

func BinanceFile(filename string) ([]model.Transaction) {
		xlFile,err := xlsx.OpenFile(filename)
		if err != nil {
			log.Fatal(err)
		}

		var x []Binance

		for _, sheet := range xlFile.Sheets {
			for _, row := range sheet.Rows[1:] {
				readStruct := &Binance{}
				err := row.ReadStruct(readStruct)
				if err != nil {
					panic(err)
				}
				x = append(x,*readStruct)
				log.Info(readStruct)
			}
		}

		var out []model.Transaction

		for _, btxn := range x {
			out = append(out, btxn.ProcessData())
		}

		return out
		
}
