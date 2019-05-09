package model

import (
	"fmt"
	"io"
	"sort"
	"time"

	log "github.com/Sirupsen/logrus"
)

type LineItem struct {
	Date      time.Time
	Amount    float64
	UnitPrice float64
	CostBase  float64
}

type Account struct {
	Name         string
	Balance      float64
	CostBase     float64
	Transactions []LineItem
}

func NewAccount(name string) *Account {
	return &Account{name, 0, 0, []LineItem{}}
}

//Amount of commodity to add to account and the total AUD cost of the
func (a *Account) Debit(date time.Time, amount, unitprice, costbase float64) {
	log.Info("===========================================")
	log.Info("Debiting Account: ", a.Name)
	log.Info("  Date: ", date)
	log.Info("  Opening Balance: ", a.Balance)
	log.Info("  Amount: ", amount)
	a.Transactions = append(a.Transactions, LineItem{date, amount, unitprice, costbase})
	a.Recalculate()
	log.Info("  Closing Balance: ", a.Balance)
}

func (a *Account) SortTxns() {
	sort.Slice(a.Transactions, func(i, j int) bool {
		return a.Transactions[i].Date.Before(a.Transactions[j].Date)
	})
}

//Takes a date and an amount to credit, returns the total cost base, the units sold that were over a year old and the cost base amount of those units
func (a *Account) removeFirst(date time.Time, amount float64, w io.Writer) (float64, float64, float64) {
	a.SortTxns()
	numberSourceTransactions := 1
	var count, costBase, discountedCostBase, discountedUnits float64
	if a.Balance < amount {
		//Balance is Not Sufficient
		log.Warn("===================================================")
		log.Warn("Not enough coins are recorded for that transactions")
		log.Warn("Account: ", a.Name)
		log.Warn("Balance: ", a.Balance)
		log.Warn("Amount Requested: ", amount)
		fmt.Fprintln(w, "")
		fmt.Fprintln(w, "Balance in account before Calculation")
		fmt.Fprintln(w, "Account:          ", a.Name)
		fmt.Fprintln(w, "Balance:          ", a.Balance)
		fmt.Fprintln(w, "Amount Requested: ", amount)
		fmt.Fprintln(w, "Insufficient balance in Account of the coin spent")
		fmt.Fprintln(w, "Costbase will only be the current balance:", a.Balance)

		amount = a.Balance
	} else {
		//Balance is Sufficient
		log.Info("===================================================")
		log.Info("Sufficient Balance in Account for Credit")
		fmt.Fprintln(w, "")
		fmt.Fprintln(w, "Balance in account before Calculation")
		fmt.Fprintln(w, "Account:          ", a.Name)
		fmt.Fprintln(w, "Balance:          ", a.Balance)
		fmt.Fprintln(w, "Amount Requested: ", amount)
		fmt.Fprintln(w, "Sufficient Balance in Account for Credit")
	}
	for count < amount && len(a.Transactions) > 0 {
		log.Info("Number of Transactions in Account: ", len(a.Transactions))
		current := a.Transactions[0]
		a.Transactions = append(a.Transactions[:0], a.Transactions[1:]...)
		if (count + current.Amount) < amount {
			log.Info("---------------------------------------------------")
			log.Info("Transaction in account not sufficient to cover total Credit")
			log.Info("Amount: ", amount)
			log.Info("Carried Forward from Previous Entries: ", count)
			log.Info("Transaction Date: ", current.Date)
			log.Info("Transaction Amount: ", current.Amount)
			log.Info("Transaction Costbase: ", current.CostBase)
			fmt.Fprintln(w, "")
			fmt.Fprintln(w, "Costbase Source Transaction #", numberSourceTransactions)
			numberSourceTransactions += 1
			fmt.Fprintln(w, "Amount carried Forward from Previous Transactions: ", count)
			fmt.Fprintln(w, "Transaction Date:                                  ", current.Date)
			fmt.Fprintln(w, "Transaction Amount:                                ", current.Amount)
			fmt.Fprintln(w, "Transaction Unitprice:                             ", current.UnitPrice)
			fmt.Fprintln(w, "Transaction Costbase:                              ", current.CostBase)
			costBase += current.CostBase
			timediff := date.Sub(current.Date)
			log.Info("Time Held: ", timediff)
			fmt.Fprintln(w, "Time Held: ", timediff)
			if timediff.Hours() > (24 * 365) {
				log.Info("Asset Held for over a year")
				fmt.Fprintln(w, "Asset Held for over a year")
				discountedCostBase += current.CostBase
				discountedUnits += current.Amount
				log.Info("Transaction Discounted Amount: ", current.CostBase)
				log.Info("Transaction Discounted Costbase: ", current.CostBase)
				fmt.Fprintln(w, "Transaction Discounted Amount:                   ", current.CostBase)
				fmt.Fprintln(w, "Transaction Discounted Costbase:                 ", current.CostBase)
			}
			count += current.Amount
			a.Recalculate()
		} else {
			log.Info("---------------------------------------------------")
			log.Info("Transaction in account is sufficient to cover total Credit")
			log.Info("Amount: ", amount)
			log.Info("Carried Forward from Previous Transactions: ", count)
			log.Info("Transaction Date: ", current.Date)
			log.Info("Transaction Amount: ", current.Amount)
			log.Info("Transaction Costbase: ", current.CostBase)
			fmt.Fprintln(w, "")
			fmt.Fprintln(w, "Costbase Source Transaction #", numberSourceTransactions)
			numberSourceTransactions += 1
			fmt.Fprintln(w, "Amount carried Forward from Previous Transactions: ", count)
			fmt.Fprintln(w, "Transaction Date:                                  ", current.Date)
			fmt.Fprintln(w, "Transaction Amount:                                ", current.Amount)
			fmt.Fprintln(w, "Transaction Costbase:                              ", current.CostBase)
			cb := current.CostBase / current.Amount
			remAmt := (amount - count)        //Remaining amount needed for Source Costbase
			repAmt := current.Amount - remAmt //New Transaction (Replacement) for the transaction partially used
			replacement := LineItem{
				Date:     current.Date,
				Amount:   repAmt,
				CostBase: repAmt * cb,
			}
			costBase += remAmt * cb
			count += remAmt
			log.Info("Prorata amount from Transaction: ", remAmt)
			log.Info("Prorata amount from Transaction Costbase: ", remAmt*cb)
			fmt.Fprintln(w, "Prorata amount from Transaction:                   ", remAmt)
			fmt.Fprintln(w, "Prorata amount from Transaction Costbase:          ", remAmt*cb)
			timediff := date.Sub(current.Date)
			log.Info("Time Held: ", timediff)
			fmt.Fprintln(w, "Time Held: ", timediff)
			//If the time held is greater than a year
			if timediff.Hours() > (24 * 365) {
				fmt.Fprintln(w, "")
				log.Info("Asset Held for over a year")
				fmt.Fprintln(w, "Asset Held for over a year")
				discountedCostBase += remAmt * cb
				discountedUnits += remAmt
				log.Info("Transaction Discounted Amount: ", remAmt)
				log.Info("Transaction Discounted Costbase: ", remAmt*cb)
				fmt.Fprintln(w, "Transaction Discounted Amount:                   ", remAmt)
				fmt.Fprintln(w, "Transaction Discounted Costbase:                 ", remAmt*cb)
			}
			a.Transactions = append([]LineItem{replacement}, a.Transactions...)
			a.Recalculate()
		}
	}
	log.Info("---------------------------------------------------")
	log.Info("Balance remaining in account after Calculation")
	log.Info("Account: ", a.Name)
	log.Info("Balance: ", a.Balance)
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Balance remaining in account after Calculation")
	fmt.Fprintln(w, "Account: ", a.Name)
	fmt.Fprintln(w, "Balance: ", a.Balance)

	return costBase, discountedCostBase, discountedUnits
}

//Takes a date and an amount to credit, returns the total cost base, the units sold that were over a year old and the cost base amount of those units
func (a *Account) Credit(date time.Time, amount float64, w io.Writer) (float64, float64, float64) {
	_ = date
	log.Info("===========================================")
	log.Info("Crediting Account: ", a.Name)
	log.Info("  Date: ", date)
	log.Info("  Opening Balance: ", a.Balance)
	log.Info("  Amount: ", amount)
	costBase, discountedCostBase, discountedAmount := a.removeFirst(date, amount, w)
	log.Info("  Closing Balance: ", a.Balance)

	return costBase, discountedCostBase, discountedAmount
}

func (a *Account) Revalue(price float64) float64 {
	a.SortTxns()
	var costBase float64
	for index, _ := range a.Transactions {
		newBase := a.Transactions[index].Amount * price
		costBase += newBase - a.Transactions[index].CostBase
		a.Transactions[index].CostBase = newBase
	}
	return costBase
}

func (a *Account) Recalculate() {
	amount := 0.0
	costBase := 0.0
	for _, txn := range a.Transactions {
		amount += txn.Amount
		costBase += txn.CostBase
	}
	a.Balance = amount
	a.CostBase = costBase
}
