package exchanges

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"cgtcalc/model"

	log "github.com/Sirupsen/logrus"
)

type BiboxLine struct {
	Time          string
	User_id       string
	Coin_symbol   string
	Bill_type     string
	Change_amount float64
	Fee           float64
	Fee_symbol    string
	Result_amount float64
	Relay_id      string
	Comment       string
}

type Bibox struct {
	RawTransactions []BiboxLine
	Txns            map[string]model.Transaction
}

func NewBibox(in []BiboxLine) (*Bibox, error) {
	marketstruct := new(Bibox)
	marketstruct.RawTransactions = in
	marketstruct.Txns = make(map[string]model.Transaction)

	return marketstruct, nil
}

func (mdl *Bibox) ProcessData() []model.Transaction {
	for _, btxn := range mdl.RawTransactions {
		if btxn.Bill_type == "3" || btxn.Bill_type == "5" {
			if value, exist := mdl.Txns[btxn.Relay_id]; exist {
				log.Info(fmt.Sprintf("Key already found in map, adding to id: %s", btxn.Relay_id))
				if btxn.Change_amount < 0 {
					if mdl.Txns[btxn.Relay_id].BaseSpent != 0 && mdl.Txns[btxn.Relay_id].BaseCurrency != btxn.Coin_symbol {
						log.Warn(fmt.Sprintf("Have a negative amount in %s, however the amount is different then the base currency", btxn.Relay_id))
						log.Warn(fmt.Sprintf("CSV %s, %.2f", btxn.Coin_symbol, btxn.Change_amount))
						log.Warn(fmt.Sprintf("Current %s, %.2f", mdl.Txns[btxn.Relay_id].BaseCurrency, mdl.Txns[btxn.Relay_id].BaseSpent))
					} else {
						value.BaseSpent += btxn.Change_amount * -1
						value.BaseCurrency = btxn.Coin_symbol
					}
				} else {
					if mdl.Txns[btxn.Relay_id].QuoteReceived != 0 && mdl.Txns[btxn.Relay_id].QuoteCurrency != btxn.Coin_symbol {
						if mdl.Txns[btxn.Relay_id].QuoteCurrency == "" {
							value.QuoteCurrency = btxn.Coin_symbol
						} else {
							log.Warn(fmt.Sprintf("Have a positive amount in %s, however the amount is different then the quote currency", btxn.Relay_id))

							log.Warn(fmt.Sprintf("CSV %s, %.2f", btxn.Coin_symbol, btxn.Change_amount))
							log.Warn(fmt.Sprintf("Current %s, %.2f", mdl.Txns[btxn.Relay_id].QuoteCurrency, mdl.Txns[btxn.Relay_id].QuoteReceived))
						}

					}
					value.QuoteReceived += btxn.Change_amount
					value.QuoteCurrency = btxn.Coin_symbol
				}
				mdl.Txns[btxn.Relay_id] = value
			} else {
				log.Info(fmt.Sprintf("Key not found in map, creating new with id: %s", btxn.Relay_id))
				var out model.Transaction
				t, err := time.Parse("2/01/2006 15:04", btxn.Time)
				if err != nil {
					log.Fatal(err)
				}
				out.Date = t
				out.BaseSpent = 0.0
				out.QuoteReceived = 0.0
				out.ExchangeID = btxn.Relay_id
				if btxn.Change_amount < 0 {
					out.BaseCurrency = btxn.Coin_symbol
					out.BaseSpent += btxn.Change_amount * -1
				} else {
					out.QuoteCurrency = btxn.Coin_symbol
					out.QuoteReceived += btxn.Change_amount
				}
				out.Exchange = "Bibox"
				mdl.Txns[btxn.Relay_id] = out
			}
		}
	}

	v := make([]model.Transaction, 0, len(mdl.Txns))

	for _, value := range mdl.Txns {
		v = append(v, value)
	}

	return v
}

func BiboxFile(filename string) []model.Transaction {
	csvFile, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	lines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var x []BiboxLine

	for _, line := range lines[1:] {
		f4, _ := strconv.ParseFloat(line[4], 64)
		f5, _ := strconv.ParseFloat(line[5], 64)
		f7, _ := strconv.ParseFloat(line[7], 64)
		data := BiboxLine{
			Time:          line[0],
			User_id:       line[1],
			Coin_symbol:   line[2],
			Bill_type:     line[3],
			Change_amount: f4,
			Fee:           f5,
			Fee_symbol:    line[6],
			Result_amount: f7,
			Relay_id:      line[8],
			Comment:       line[9],
		}
		x = append(x, data)
		log.Info(data)
	}

	mdl, _ := NewBibox(x)
	return mdl.ProcessData()

}
