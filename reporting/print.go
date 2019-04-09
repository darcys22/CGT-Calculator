package reporting

import (
	"fmt"

	"cgtcalc/model"
	"cgtcalc/version"

	"github.com/jung-kurt/gofpdf"
	"github.com/leekchan/accounting"
	log "github.com/Sirupsen/logrus"
)

func Printpdf(m *model.Model) {
	aud := accounting.Accounting{Precision: 2}
	cry := accounting.Accounting{Precision: 4}
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Courier", "B", 12)

	pdf.SetFooterFunc(func() {
			pdf.SetY(-15)
			pdf.SetFont("Arial", "I", 8)
			pdf.CellFormat(0, 10, fmt.Sprintf("Page %d", pdf.PageNo()),
					"", 0, "R", false, 0, "")
	})

	heading := fmt.Sprintf("Between %s - %s", m.StartDate.Format("2 January 2006 (15:04:05)"),m.EndDate.Format("2 January 2006 (15:04:05)") )

	pdf.Cell(40, 10, fmt.Sprintf("Crypto Capital Gains Calculation %s",version.VersionWithMeta))
	pdf.Ln(-1)
	pdf.Cell(40, 10, heading)
	pdf.Ln(-1)
	pdf.Ln(-1)
	pdf.Bookmark("Summary", 0, 0)
	pdf.Cell(40, 10, fmt.Sprintf("Summary"))
	pdf.Ln(5)
	pdf.Cell(40, 10, fmt.Sprintf("---------------"))
	pdf.Ln(-1)
	pdf.Cell(40, 10, fmt.Sprintf("Transactions Analysed: %d", (len(m.Txns))))
	pdf.Ln(-1)
	pdf.Cell(40, 10, fmt.Sprintf("Capital Gains During Period: A$ %s", aud.FormatMoney(m.CapitalGains)))
	pdf.Ln(-1)
	pdf.Cell(40, 10, fmt.Sprintf("Capital Gains that can utilise CGT Discount During Period: A$ %s", aud.FormatMoney(m.DiscountAvailableCGTGain)))
	pdf.Ln(-1)
	pdf.Cell(40, 10, fmt.Sprintf("Trading Gains During Period: A$ %s [Includes Unrealised Gains (Loss) of A$ %s]", aud.FormatMoney(m.TradingGains), aud.FormatMoney(m.UnrealisedGains)))
	pdf.Ln(-1)
	pdf.Ln(-1)
	pdf.Ln(-1)

	pdf.Bookmark("Transactions - Capital Gains", 0, 0)
	pdf.Cell(40, 10, fmt.Sprintf("Transactions - Capital Gains"))
	pdf.Ln(5)
	pdf.Cell(40, 10, fmt.Sprintf("----------------------------"))
	pdf.Ln(-1)



	w := []float64{28.0, 50.0, 18.0, 33.0, 18.0, 33.0, 33.0, 33.0, 33.0}
	wSum := 0.0
	for _, v := range w {
			wSum += v
	}

	// 	Header
	header := []string{"Date", "Exchange", "Base", "Spent", "Quote", " Received", "Proceeds A$", "Cost Base $A", "Gains $A"}
	for j, str := range header {
			pdf.CellFormat(w[j], 7, str, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	for _, txn := range m.Txns {
			if pdf.GetY() > 180 {
				pdf.CellFormat(wSum, 0, "", "T", 0, "", false, 0, "")
				pdf.AddPage()
				for j, str := range header {
						pdf.CellFormat(w[j], 7, str, "1", 0, "C", false, 0, "")
				}
				pdf.Ln(-1)
			}
			pdf.CellFormat(w[0], 6, txn.Date.Format("2006-01-02"), "LR", 0, "", false, 0, "")
			pdf.CellFormat(w[1], 6, txn.Exchange, "LR", 0, "C", false, 0, "")
			pdf.CellFormat(w[2], 6, txn.BaseCurrency, "LR", 0, "C", false, 0, "")
			pdf.CellFormat(w[3], 6, aud.FormatMoney(txn.BaseSpent), "LR", 0, "R", false, 0, "")
			pdf.CellFormat(w[4], 6, txn.QuoteCurrency, "LR", 0, "C", false, 0, "")
			pdf.CellFormat(w[5], 6, aud.FormatMoney(txn.QuoteReceived), "LR", 0, "R", false, 0, "")
			pdf.CellFormat(w[6], 6, aud.FormatMoney(txn.Proceeds), "LR", 0, "R", false, 0, "")
			pdf.CellFormat(w[7], 6, aud.FormatMoney(txn.CostBase), "LR", 0, "R", false, 0, "")
			pdf.CellFormat(w[8], 6, aud.FormatMoney(txn.Gain), "LR", 0, "R", false, 0, "")
			pdf.Ln(-1)
	}
	pdf.CellFormat(wSum, 0, "", "T", 0, "", false, 0, "")


	pdf.AddPage()
	pdf.Bookmark("Transactions - Trading Gains", 0, 0)
	pdf.Cell(40, 10, fmt.Sprintf("Transactions - Trading Gains"))
	pdf.Ln(5)
	pdf.Cell(40, 10, fmt.Sprintf("----------------------------"))
	pdf.Ln(-1)
	for j, str := range header {
			pdf.CellFormat(w[j], 7, str, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	for _, txn := range m.Txns {
			if pdf.GetY() > 180 {
				pdf.CellFormat(wSum, 0, "", "T", 0, "", false, 0, "")
				pdf.AddPage()
				for j, str := range header {
						pdf.CellFormat(w[j], 7, str, "1", 0, "C", false, 0, "")
				}
				pdf.Ln(-1)
			}
			pdf.CellFormat(w[0], 6, txn.Date.Format("2006-01-02"), "LR", 0, "", false, 0, "")
			pdf.CellFormat(w[1], 6, txn.Exchange, "LR", 0, "C", false, 0, "")
			pdf.CellFormat(w[2], 6, txn.BaseCurrency, "LR", 0, "C", false, 0, "")
			pdf.CellFormat(w[3], 6, aud.FormatMoney(txn.BaseSpent), "LR", 0, "R", false, 0, "")
			pdf.CellFormat(w[4], 6, txn.QuoteCurrency, "LR", 0, "C", false, 0, "")
			pdf.CellFormat(w[5], 6, aud.FormatMoney(txn.QuoteReceived), "LR", 0, "R", false, 0, "")
			pdf.CellFormat(w[6], 6, aud.FormatMoney(txn.Proceeds), "LR", 0, "R", false, 0, "")
			pdf.CellFormat(w[7], 6, aud.FormatMoney(txn.TradingCostBase), "LR", 0, "R", false, 0, "")
			pdf.CellFormat(w[8], 6, aud.FormatMoney(txn.TradingGain), "LR", 0, "R", false, 0, "")
			pdf.Ln(-1)
	}
	pdf.CellFormat(wSum, 0, "", "T", 0, "", false, 0, "")
	pdf.AddPageFormat("P",pdf.GetPageSizeStr("a4"))
	pdf.SetFont("Courier", "B", 10)
	pdf.Bookmark("Account Balances - Capital Gains", 0, 0)
	pdf.Cell(40, 10, fmt.Sprintf("Account Balances"))
	pdf.Ln(5)
	pdf.Cell(40, 10, fmt.Sprintf("---------------"))
	for _, account := range m.Accounts {
		if (len(account.Transactions) > 0) && (account.Name != "AUD") {
			if pdf.GetY() > 180 {
				pdf.AddPageFormat("P",pdf.GetPageSizeStr("a4"))
				pdf.SetFont("Courier", "B", 10)
			}
			pdf.Ln(-1)
			pdf.Cell(40, 10, fmt.Sprintf("Account: %s", account.Name))
			pdf.Ln(-1)
			pdf.Cell(40, 10, fmt.Sprintf("    Balance: %s %s", cry.FormatMoney(account.Balance), account.Name))
			pdf.Ln(-1)
			pdf.Cell(40, 10, fmt.Sprintf("    Cost Base: A$ %s", aud.FormatMoney(account.CostBase)))
			pdf.Ln(-1)
			for idx, txn := range account.Transactions {
				pdf.Cell(40, 10, fmt.Sprintf("    [%d] %s %s, purchased on %s with a cost base of A$ %s",idx, cry.FormatMoney(txn.Amount),account.Name, txn.Date.Format("2 January 2006"), aud.FormatMoney(txn.CostBase)))
				pdf.Ln(-1)
			}
		}
	}



	//err := pdf.OutputFileAndClose(heading + ".pdf")
	err := pdf.OutputFileAndClose("outputcalculation.pdf")
	if err != nil {
		log.Fatal(err)
	} 
}
