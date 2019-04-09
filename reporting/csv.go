package reporting

import (
	"fmt"
	"os"
	"encoding/csv"

	"cgtcalc/model"
	//"cgtcalc/version"

	//"github.com/leekchan/accounting"
	log "github.com/Sirupsen/logrus"
)

func ExportCSV(m *model.Model) {

	file, err := os.Create("result.csv")
	if err != nil {
		log.Fatal("Couldnt create the CSV File to export", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	heading := []string{fmt.Sprintf("Gains between %s - %s", m.StartDate.Format("2 January 2006 (15:04:05)"),m.EndDate.Format("2 January 2006 (15:04:05)") )}
	writer.Write(heading)

	header := []string{"ID","Date", "Exchange", "Base", "Spent", "Quote", " Received", "Proceeds A$", "Cost Base A$", "Total Gains A$", "Discount Available Gain A$"}
	writer.Write(header)
	for _, txn := range m.Txns {
		row := []string{
			fmt.Sprintf("%d",txn.Nonce),
			txn.Date.Format("2006-01-02"),
			txn.Exchange,
			txn.BaseCurrency,
			fmt.Sprintf("%f",txn.BaseSpent),
			txn.QuoteCurrency,
			fmt.Sprintf("%f",txn.QuoteReceived),
			fmt.Sprintf("%f",txn.Proceeds),
			fmt.Sprintf("%f",txn.CostBase),
			fmt.Sprintf("%f",txn.Gain),
			fmt.Sprintf("%f",txn.DiscountAvailableCGTGain),
		}
		writer.Write(row)
	}
}
