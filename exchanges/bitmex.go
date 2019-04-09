package exchanges

import (
	//"strings"
	"strconv"
	"os"
	"encoding/csv"
	//"time"


	"cgtcalc/model"

	log "github.com/Sirupsen/logrus"
)

type Bitmex struct {
	Timestamp 			string 				
	Symbol					string 				
	ExecType 				string				
	Side						string				
	LastQty					float64				
	LastPx					float64				
	ExecCost				float64				
	Commission			float64				
	ExecComm				float64	
	OrdType					string
	OrdQty					float64				
	LeavesQty				float64				
	Price						float64				
	Text						string
	OrderID					string
}

func (txn *Bitmex) ProcessData() (out model.Transaction){
	//TODO: The CSV didnt actually have any data in it, Unsure of the rest only know the headers
	//t, err := time.Parse("2006-01-02 15:04:05",txn.Timestamp)
	//if err != nil {
		//log.Fatal(err)
	//}
	//out.Date = t
	//out.BaseCurrency, out.QuoteCurrency = model.MarketConverter(txn.Market, txn.Type)
	//out.ExchangeID = txn.Timestamp + txn.OrderNumber
	//if (strings.ToUpper(txn.Type)=="BUY") {
		//out.QuoteReceived = txn.QuoteTotalLessFee
		//out.BaseSpent = txn.Total
	//} else {
		//out.BaseSpent = txn.Amount
		//out.QuoteReceived = txn.BaseTotalLessFee
	//}

	out.Exchange = "Bitmex"

	return	
}

func BitmexFile(filename string) ([]model.Transaction) {

		log.Info("Opening File")
		log.Warn("BITMEX IS NOT FULLY IMPLEMENTED THIS NEEDS TO BE UPDATED")
		log.Warn("DO NOT USE")
		csvFile,err := os.Open(filename)
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
			f4, _ := strconv.ParseFloat(line[4],64)
			f5, _ := strconv.ParseFloat(line[5],64)
			f6, _ := strconv.ParseFloat(line[6],64)
			f7, _ := strconv.ParseFloat(line[7],64)
			f8, _ := strconv.ParseFloat(line[8],64)
			f10, _ := strconv.ParseFloat(line[10],64)
			f11, _ := strconv.ParseFloat(line[11],64)
			f12, _ := strconv.ParseFloat(line[12],64)
			log.Info("Processed the floats")
			data := Bitmex{
				Timestamp:						line[0],
				Symbol:								line[1],
				ExecType:							line[2],
				Side:									line[3],
				LastQty:							f4,
				LastPx:								f5,
				ExecCost:							f6,
				Commission:						f7,
				ExecComm:							f8,
				OrdType:							line[9],
				OrdQty:								f10,
				LeavesQty:						f11,
				Price:								f12,
				Text:									line[13],
				OrderID:							line[14],
			}

			log.Info(data)
			log.Info("Processed the data")
			out = append(out, data.ProcessData())
		}

		return out

}
