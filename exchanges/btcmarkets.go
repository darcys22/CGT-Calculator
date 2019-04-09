package exchanges

import (
	"fmt"
	"strconv"
	"os"
	"encoding/csv"


	"cgtcalc/model"

	"github.com/araddon/dateparse"

	log "github.com/Sirupsen/logrus"
)

type BtcmarketsLine struct {
	CreationTime 		string 				`xlsx:"0"`
	RecordType 			string 				`xlsx:"1"`
	Action 					string 				`xlsx:"2"`
	Currency 				string 				`xlsx:"3"`
	Amount 					float64				`xlsx:"4"`
	Description 		string 				`xlsx:"5"`
	ReferenceId 		string 				`xlsx:"6"`
}

type Btcmarket struct {
	RawTransactions	[]BtcmarketsLine
	Txns						map[string]model.Transaction
}

func NewBtcmarket(in []BtcmarketsLine) (*Btcmarket, error) {
	marketstruct := new(Btcmarket)
	marketstruct.RawTransactions = in
	marketstruct.Txns = make(map[string]model.Transaction)

	return marketstruct, nil
}

func (mdl *Btcmarket) ProcessData() ([]model.Transaction){
	for _, btxn := range mdl.RawTransactions {
		if btxn.RecordType == "Trade" {
			if value, exist := mdl.Txns[btxn.ReferenceId]; exist {
				log.Info(fmt.Sprintf("Key already found in map, adding to id: %s",btxn.ReferenceId))
				if btxn.Amount < 0 {
					if mdl.Txns[btxn.ReferenceId].BaseCurrency != btxn.Currency {
						if btxn.Action == "Trading Fee" {
							value.QuoteReceived += btxn.Amount
						} else {
							log.Warn(fmt.Sprintf("Have a negative amount in %s, however the amount is different then the base currency",btxn.ReferenceId))
							log.Warn(fmt.Sprintf("CSV %s, %s, %.2f",btxn.Currency, btxn.Action, btxn.Amount))
							log.Warn(fmt.Sprintf("Current %s, %.2f",mdl.Txns[btxn.ReferenceId].BaseCurrency, mdl.Txns[btxn.ReferenceId].BaseSpent))
						}
					} else {
						value.BaseSpent += btxn.Amount * -1
					}
				} else {
					if mdl.Txns[btxn.ReferenceId].QuoteCurrency != btxn.Currency {
						if mdl.Txns[btxn.ReferenceId].QuoteCurrency == "" {
							value.QuoteCurrency = btxn.Currency
						} else {
							log.Warn(fmt.Sprintf("Have a positive amount in %s, however the amount is different then the quote currency",btxn.ReferenceId))

							log.Warn(fmt.Sprintf("CSV %s, %s, %.2f",btxn.Currency, btxn.Action, btxn.Amount))
							log.Warn(fmt.Sprintf("Current %s, %.2f",mdl.Txns[btxn.ReferenceId].QuoteCurrency, mdl.Txns[btxn.ReferenceId].QuoteReceived))
						}

					}
					value.QuoteReceived += btxn.Amount
				}
				mdl.Txns[btxn.ReferenceId] = value
			} else {
				log.Info(fmt.Sprintf("Key not found in map, creating new with id: %s",btxn.ReferenceId))
				var out model.Transaction
				t, err := dateparse.ParseAny(btxn.CreationTime)
				if err != nil {
					log.Fatal(err)
				}
				out.Date = t
				out.BaseSpent = 0.0
				out.QuoteReceived = 0.0
				out.ExchangeID = btxn.ReferenceId
				if btxn.Amount < 0 {
					out.BaseCurrency = btxn.Currency
					out.BaseSpent += btxn.Amount * -1
				} else {
					out.QuoteCurrency = btxn.Currency
					out.QuoteReceived += btxn.Amount
				}
				out.Exchange = "BTC Markets"
				mdl.Txns[btxn.ReferenceId] = out
			}
		}
	}

	v := make([]model.Transaction, 0, len(mdl.Txns))

	for  _, value := range mdl.Txns {
		v = append(v, value)
	}

	return v
}

func BtcMarketsFile(filename string) ([]model.Transaction) {
		csvFile,err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer csvFile.Close()

		lines, err := csv.NewReader(csvFile).ReadAll()
		if err != nil {
			log.Fatal(err)
		}

		var x []BtcmarketsLine

		for _, line := range lines[1:] {
				amt, _ := strconv.ParseFloat(line[4],64)
        data := BtcmarketsLine{
							CreationTime: 		line[0],
							RecordType: 			line[1],
							Action: 					line[2],
							Currency: 				line[3],
							Amount: 					amt,
							Description: 			line[5],
							ReferenceId: 			line[6],
        }
				x = append(x,data)
				log.Info(data)
		}

		mdl,_ := NewBtcmarket(x)
		return mdl.ProcessData()



		
}
