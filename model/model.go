package model

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"cgtcalc/cgtdb"
	"cgtcalc/config"
	"cgtcalc/prices"

	log "github.com/Sirupsen/logrus"
)


type Transaction struct {
	Date           						time.Time
	Exchange									string
	ExchangeID								string
	BaseCurrency   						string
	BaseSpent     						float64
	QuoteCurrency   					string
	QuoteReceived   					float64
	Proceeds									float64
	CostBase									float64
	TradingCostBase						float64
	Gain											float64
	TradingGain								float64
	DiscountAvailableCGTGain	float64
	Trace											[]byte
	Nonce											int
}

func NewTransactionShortDate(date, buycurrency, buyamount, sellcurrency, sellamount string) *Transaction{
	return NewTransaction(date+"T00:00:01.000000000",buycurrency,buyamount,sellcurrency,sellamount)
}

func NewTransaction(date, buycurrency, buyamount, sellcurrency, sellamount string) *Transaction{

		b := strings.Replace(buyamount, "$", "", -1)
		buy, err := strconv.ParseFloat(strings.TrimSpace(strings.Replace(b, ",", "", -1)), 64)
		if err != nil {
			panic(err)
		}

		s := strings.Replace(sellamount, "$", "", -1)
		sell, err := strconv.ParseFloat(strings.TrimSpace(strings.Replace(s, ",", "", -1)), 64)
		if err != nil {
			panic(err)
		}

		layout := "2006-01-02T15:04:05.999999999"
		t, err := time.Parse(layout, date)
		if err != nil {
			fmt.Println(err)
		}

		return &Transaction{
			Date:											t,
			Exchange:									"Custom",
			BaseCurrency:							strings.ToUpper(sellcurrency),
			BaseSpent:								sell,
			QuoteCurrency:						strings.ToUpper(buycurrency),
			QuoteReceived:						buy,
			Nonce:										0,
		}
}

type Model struct {
	GainsDB										cgtdb.Database
	Pricer										prices.Pricer
	Config										*config.Config	

	Nonce											int
	Txns											[]Transaction
	Files											[]string

	Accounts									map[string]*Account


	CapitalGains							float64
	DiscountAvailableCGTGain	float64
	TradingGains							float64
	UnrealisedGains						float64
	StartDate									time.Time
	EndDate										time.Time
}

func NewModel(gainsDB, priceDB cgtdb.Database, cfg *config.Config) *Model{
	log.Info("New Model Called")

	if gainsDB == nil {
		log.Info("Nil parameter passed for new model gainDB, using Config to open")
		log.Info("Location: ", cfg.GainsDatabase)
		newgainsDB, err := cgtdb.NewLDBDatabase(cfg.GainsDatabase)
		if err != nil {
			log.Fatal("Could not Open the Gain DB ",err)
		}
		gainsDB = newgainsDB
	}
	if priceDB == nil {
		log.Info("Nil parameter passed for new model priceDB, using Config to open")
		newpriceDB, err := cgtdb.NewLDBDatabase(cfg.PriceDatabase)
		log.Info("Location: ", cfg.PriceDatabase)
		if err != nil {
			log.Fatal("Could not Open the Price DB ",err)
		}
		priceDB = newpriceDB
	}

	p := *prices.NewPricer(priceDB)

	gob.Register(Transaction{})
	var txns []Transaction
	nonce := 0
	data , err  := gainsDB.Get([]byte("txns"))
	if err == nil {
		log.Info("Transactions found to load")
		buf := bytes.NewBuffer(data)
		d := gob.NewDecoder(buf)
		if err := d.Decode(&txns); err != nil {
			log.Fatal("Decoding of the Transactions failed: ",err)
		}
		log.Info("txns loaded")
		nonce = len(txns)
	}

	var files []string
	data2 , err := gainsDB.Get([]byte("files"))
	if err == nil {
		log.Info("Previous file checksums found to load")
		buf2 := bytes.NewBuffer(data2)
		d2 := gob.NewDecoder(buf2)
		if err := d2.Decode(&files); err != nil {
			log.Warn("Decoding of the checksums failed: ",err)
		}
		log.Info("checksums loaded")
	}

	return &Model{
		GainsDB:									gainsDB,
		Pricer:										p,
		Config:										cfg,
		Nonce:										nonce,
		Txns:											txns,
		Files:										files,
		Accounts:									make(map[string]*Account),
		CapitalGains:							0.0,
		DiscountAvailableCGTGain: 0.0,
		TradingGains:							0.0,
		UnrealisedGains:					0.0,
		StartDate:								cfg.StartFY,
		EndDate:									cfg.EndFY,
	}
}
func (m *Model) addAccount(name string) {
	if _, ok := m.Accounts[name]; !ok {
		m.Accounts[name] = NewAccount(name)
	}
}


func (m *Model) ProcessModel() {
	log.Info("Process Model Called -- Number of Transactions: ",len(m.Txns))
	m.resetGains()
	m.processTradingGains()
	m.resetGains()
	m.processCapitalGains()
	log.Info("Finished Processing Capital Gains")
}
func (m *Model) resetGains() {
		m.Accounts = make(map[string]*Account)
}

func (m *Model) processCapitalGains() {
	m.CapitalGains = 0.0
	m.DiscountAvailableCGTGain = 0.0
	//loc, _ := time.LoadLocation("Australia/Sydney")
	log.Info("Processing the Capital Gains")
	log.Info("Number of Transactions: ",len(m.Txns))

	for idx, _ := range m.Txns {
		var buff bytes.Buffer
		fmt.Fprintln(&buff, "#### Transaction Information")
		fmt.Fprintln(&buff)
		fmt.Fprintln(&buff, "Financial Year Start Date:", m.StartDate)
		fmt.Fprintln(&buff, "Financial Year End Dated: ", m.EndDate)
		fmt.Fprintln(&buff, "CGT Event Date:           ", m.Txns[idx].Date)
		fmt.Fprintln(&buff, "Exchange:                 ", m.Txns[idx].Exchange)
		fmt.Fprintln(&buff, "Transaction ID:           ", m.Txns[idx].ExchangeID)
		fmt.Fprintln(&buff)
		fmt.Fprintln(&buff, "Base Currency (Spent):    ", m.Txns[idx].BaseCurrency)
		fmt.Fprintln(&buff, "Base Amount:              ", m.Txns[idx].BaseSpent)
		fmt.Fprintln(&buff)
		fmt.Fprintln(&buff, "Quote Currency (Received):", m.Txns[idx].QuoteCurrency)
		fmt.Fprintln(&buff, "Quote Amount:             ", m.Txns[idx].QuoteReceived)
		if m.Txns[idx].Date.After(m.EndDate) {
			fmt.Fprintln(&buff)
			fmt.Fprintln(&buff, "Did not continue as the date was after the financial year")
			m.Txns[idx].Trace = buff.Bytes()
			continue
		}

		//Handle XBT & BTC being different identifiers for the same currency
		if strings.ToUpper(m.Txns[idx].BaseCurrency) == "XBT" {
			m.Txns[idx].BaseCurrency = "BTC"
		}
		if strings.ToUpper(m.Txns[idx].QuoteCurrency) == "XBT" {
			m.Txns[idx].QuoteCurrency = "BTC"
		}

		//Create new accounts in the model for both the base and the quote, if they already exist addAccount will ignore
		m.addAccount(m.Txns[idx].BaseCurrency)
		m.addAccount(m.Txns[idx].QuoteCurrency)

		//Consults the price DB for the AUD price of the quote currency (What is being received)
		price := m.Pricer.GetPrice(m.Txns[idx].Date.Format("2006-01-02"), m.Txns[idx].QuoteCurrency)

		//Sets the proceeds for the transaction to be the price x how much was received
		m.Txns[idx].Proceeds = price * m.Txns[idx].QuoteReceived

		//If the Base is AUD (Spending AUD to buy crypto) then you dont have to calculate the gain but you need to set the costbase
		if m.Txns[idx].BaseCurrency == "AUD" {
			fmt.Fprintln(&buff)
			fmt.Fprintln(&buff, "Buying Crypto with AUD")
			fmt.Fprintln(&buff, "No CGT Event -> Capital Gains: A$0")
			m.Txns[idx].CostBase = m.Txns[idx].BaseSpent
			m.Txns[idx].Proceeds = m.Txns[idx].CostBase
			m.Txns[idx].Gain = 0
		} else {
		//Crediting an existing account within the model is how the costbase is established
			discountedCostBase := 0.0
			discountedUnits := 0.0
			m.Txns[idx].CostBase, discountedCostBase, discountedUnits = m.Accounts[m.Txns[idx].BaseCurrency].Credit(m.Txns[idx].Date,m.Txns[idx].BaseSpent,&buff)
			m.Txns[idx].Gain = m.Txns[idx].Proceeds - m.Txns[idx].CostBase
			m.Txns[idx].DiscountAvailableCGTGain = m.Txns[idx].Proceeds/m.Txns[idx].BaseSpent*discountedUnits - discountedCostBase

			fmt.Fprintln(&buff)
			fmt.Fprintln(&buff, "Total Costbase Information from Source Transactions:")
			fmt.Fprintln(&buff, "Total Units:        ", m.Txns[idx].BaseSpent)
			fmt.Fprintln(&buff, "Total Costbase:     ", m.Txns[idx].CostBase)
			fmt.Fprintln(&buff, "Discounted Units:   ", discountedUnits)
			fmt.Fprintln(&buff, "Discounted Costbase:", discountedCostBase)
			fmt.Fprintln(&buff)
			fmt.Fprintln(&buff, "Information regarding Proceeds:")
			fmt.Fprintln(&buff, "AUD unit price of Quote Currency:  ", price)
			fmt.Fprintln(&buff, "Total AUD Value of Quote received: ", m.Txns[idx].Proceeds)
			fmt.Fprintln(&buff)
			fmt.Fprintln(&buff, "Capital Gain Information:")
			fmt.Fprintln(&buff, "Total Gain:                        ", m.Txns[idx].Gain)
			fmt.Fprintln(&buff, "Gains that can utilise Discount:   ", m.Txns[idx].DiscountAvailableCGTGain)
			fmt.Fprintln(&buff, "Discount Available:                ", m.Txns[idx].DiscountAvailableCGTGain/2)
			fmt.Fprintln(&buff, "Gain After Discount:               ", m.Txns[idx].Gain - m.Txns[idx].DiscountAvailableCGTGain/2)

		}
		
		//Increase (Debit) the account for the currency received
		m.Accounts[m.Txns[idx].QuoteCurrency].Debit(m.Txns[idx].Date,m.Txns[idx].QuoteReceived,price,m.Txns[idx].Proceeds)

		//Checking that the date is correct before adding to total in model

		log.Info("==============================")
		log.Info("THE Date is")
		log.Info(m.Txns[idx].Date)
		log.Info("THE Start Date is")
		log.Info(m.StartDate)
		log.Info("THE End Date is")
		log.Info(m.EndDate)

		if (m.Txns[idx].Date.After(m.StartDate) && m.Txns[idx].Date.Before(m.EndDate)) || m.Txns[idx].Date == m.EndDate {
			m.CapitalGains += m.Txns[idx].Gain
			m.DiscountAvailableCGTGain += m.Txns[idx].DiscountAvailableCGTGain
			m.Txns[idx].Trace = buff.Bytes()
		}
		log.Info("Index: ")	
		log.Info(idx)	
		log.Info("Total: ")	
		log.Info(len(m.Txns))
	}
}

func (m *Model) processTradingGains() {
	m.TradingGains = 0.0

	//Build up the Accounts Before the start of the trading period
	for idx, _ := range m.Txns {
		if m.Txns[idx].Date.After(m.StartDate) {
			continue
		}
		var buff bytes.Buffer

		if strings.ToUpper(m.Txns[idx].BaseCurrency) == "XBT" {
			m.Txns[idx].BaseCurrency = "BTC"
		}
		if strings.ToUpper(m.Txns[idx].QuoteCurrency) == "XBT" {
			m.Txns[idx].QuoteCurrency = "BTC"
		}
		m.addAccount(m.Txns[idx].BaseCurrency)
		m.addAccount(m.Txns[idx].QuoteCurrency)

		price := m.Pricer.GetPrice(m.Txns[idx].Date.Format("2006-01-02"), m.Txns[idx].QuoteCurrency)
		m.Txns[idx].Proceeds = price * m.Txns[idx].QuoteReceived
		if m.Txns[idx].BaseCurrency == "AUD" {
			m.Txns[idx].TradingCostBase = m.Txns[idx].BaseSpent
			m.Txns[idx].Proceeds = m.Txns[idx].TradingCostBase
			m.Txns[idx].Gain = 0
		} else {
			m.Txns[idx].TradingCostBase, _, _ = m.Accounts[m.Txns[idx].BaseCurrency].Credit(m.Txns[idx].Date,m.Txns[idx].BaseSpent, &buff)
			m.Txns[idx].Gain = m.Txns[idx].Proceeds - m.Txns[idx].TradingCostBase
		}
		
		m.Accounts[m.Txns[idx].QuoteCurrency].Debit(m.Txns[idx].Date,m.Txns[idx].QuoteReceived,price,m.Txns[idx].Proceeds)
	}
	
	//Set inventory start value
	_ = m.Revalue(m.StartDate.Add(time.Hour * -24))

	//Process the transactions in the trading period

	for idx, _ := range m.Txns {
		if m.Txns[idx].Date.Before(m.StartDate) || m.Txns[idx].Date.After(m.EndDate) {
			continue
		}
		var buff bytes.Buffer
		m.addAccount(m.Txns[idx].BaseCurrency)
		m.addAccount(m.Txns[idx].QuoteCurrency)

		price := m.Pricer.GetPrice(m.Txns[idx].Date.Format("2006-01-02"), m.Txns[idx].QuoteCurrency)
		m.Txns[idx].Proceeds = price * m.Txns[idx].QuoteReceived
		if m.Txns[idx].BaseCurrency == "AUD" {
			m.Txns[idx].TradingCostBase = m.Txns[idx].BaseSpent
			m.Txns[idx].Proceeds = m.Txns[idx].TradingCostBase
			m.Txns[idx].TradingGain = 0
		} else {
			m.Txns[idx].TradingCostBase, _, _ = m.Accounts[m.Txns[idx].BaseCurrency].Credit(m.Txns[idx].Date,m.Txns[idx].BaseSpent,&buff)
			m.Txns[idx].TradingGain = m.Txns[idx].Proceeds - m.Txns[idx].TradingCostBase
		}
		
		m.Accounts[m.Txns[idx].QuoteCurrency].Debit(m.Txns[idx].Date,m.Txns[idx].QuoteReceived,price,m.Txns[idx].Proceeds)


		m.TradingGains += m.Txns[idx].TradingGain
	}

	//Set Unrealised Gains on closing balance
	m.UnrealisedGains = m.Revalue(m.EndDate)
	m.TradingGains += m.UnrealisedGains
	
}

func (m *Model) Revalue(date time.Time) float64 {
	log.Info("======================")
	log.Info("Revaluing the Accounts")
	log.Info("Date: ", date)
	unrealisedGain := 0.0
	for idx, _ := range m.Accounts {
		log.Info("----------------------")
		name := m.Accounts[idx].Name
		log.Info("Currency: ", name)
		price := m.Pricer.GetPrice(date.Format("2006-01-02"), name)
		unrealisedGain += m.Accounts[idx].Revalue(price)
	}
	return unrealisedGain
}
func (m *Model) Close() error {
		log.Info("Closing the Model Databases")
		log.Info("Closing GainsDB")
		m.GainsDB.Close()
		m.Pricer.Close()
		return nil
}

func (m *Model) Commit() error {

	log.Info("Commiting the Txns to the DB")
	log.Info("Length: ", len(m.Txns))
	buf := bytes.NewBuffer([]byte{})
	enc := gob.NewEncoder(buf)
	err := enc.Encode(m.Txns)
	if err != nil {
		return err
	}

	err = m.GainsDB.Put([]byte("txns"), buf.Bytes())
	if err != nil {
		return err
	}

	log.Info("Commiting the Checksums to the DB")
	log.Info("Length: ", len(m.Files))
	buf2 := bytes.NewBuffer([]byte{})
	enc2 := gob.NewEncoder(buf2)
	err = enc2.Encode(m.Files)
	if err != nil {
		return err
	}

	err = m.GainsDB.Put([]byte("files"), buf2.Bytes())
	if err != nil {
		return err
	}

	log.Info("Completed Commiting the Txns and Checksums to the Database")
	return nil
}

func (m *Model) AddSingleTxn(txn Transaction) error {
	log.Info("Adding Transaction to the Database")
	m.Txns = append(m.Txns, txn)
	m.SortTxns()
	return m.Commit()
}

func (m *Model) AddTxns(txns []Transaction) error {
	log.Info("Adding Transactions to the Database")
	m.Txns = append(m.Txns, txns...)
	m.SortTxns()
	m.UpdateNonces()
	return m.Commit()
}

func (m *Model) SortTxns() {
	log.Info("Sorting the Transactions withing the database")
	log.Info("Length: ", len(m.Txns))
	sort.Slice(m.Txns, func(i, j int) bool {
	  return m.Txns[i].Date.Before(m.Txns[j].Date)
	})
}

func (m *Model) GetNonce() int {
	m.Nonce += 1
	return m.Nonce
}

func (m *Model) UpdateNonces() {
	m.Nonce = 0
	for idx, _ := range m.Txns {
		m.Txns[idx].Nonce = m.GetNonce()
	}

}

func (m *Model) Checksum(filename string) (bool){

	log.Info("Reviewing Checksum of ", filename)
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		log.Fatal(err)
	}

	hashInBytes := hash.Sum(nil)[:16]

	MD5String := hex.EncodeToString(hashInBytes)
	log.Info("Checksum: ", MD5String)

	for _, n := range m.Files {
		if MD5String == n {
			log.Warn("Checksum of file has already been used before")
			return false
		}
	}
	m.Files = append(m.Files, MD5String)
	
	return true
}

func (m *Model) removeDuplicates() {   
	encountered := map[string]bool{}          
	result := []Transaction{}                  

	for v := range m.Txns {
		if encountered[m.Txns[v].Date.String()] == true {
		} else {
			encountered[m.Txns[v].Date.String()] = true
			result = append(result, m.Txns[v])
		}
	}
	m.Txns = result
}

