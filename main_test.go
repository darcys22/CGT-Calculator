package main

import (
	"os"
	"io/ioutil"
	"path/filepath"
	"encoding/json"
	"time"

	"cgtcalc/cgtdb"
	"cgtcalc/config"
	"cgtcalc/model"

	"github.com/leekchan/accounting"

	"testing"
)

func TestMemorySingle(t *testing.T) {
	aud := accounting.Accounting{Precision: 2}

	cfg := defaultTestConfig()

	file := "./testdata/simplebuysell.golden"

	t.Log("------------------------------------------------------------")
	t.Log("Testing: " + file)
	g, err := ioutil.ReadFile(file)
	if err != nil {
		t.Fatalf("failed reading .golden: %s", err)
	}

	m := *model.NewModel(cgtdb.NewMemDatabase(), cgtdb.NewMemDatabase(),cfg)

	m.Pricer.AddSinglePrice("2018-06-17","BTC","6112.40301","AUD")
	m.Pricer.AddSinglePrice("2018-06-30","BTC","6391.5000","USD")
	m.Pricer.AddSinglePrice("2018-06-30","USD","1.3503","AUD")

	err = json.Unmarshal(g, &m)
	if err != nil {
		t.Fatalf("failed reading json: %s", err)
	}

	ExpectedGains := m.CapitalGains
	ExpectedTradingGains := m.TradingGains
	t.Logf("Number of Transactions Analysed: %d", len(m.Txns))
	t.Logf("Expected Capital Gains from Transactions: %s", aud.FormatMoney(ExpectedGains))
	m.ProcessModel()
	t.Logf("Actual Capital Gains Results: %s", aud.FormatMoney(m.CapitalGains))
	t.Logf("Expected Trading Gains from Transactions: %s", aud.FormatMoney(ExpectedTradingGains))
	t.Logf("Actual Trading Gains Results: %s", aud.FormatMoney(m.TradingGains))

	if ExpectedGains != m.CapitalGains {
		 t.Errorf("Expected Capital Gains for %s was incorrect, got: %f, want: %f.", file, m.CapitalGains, ExpectedGains)
	}
	if ExpectedTradingGains != m.TradingGains {
		 t.Errorf("Expected Trading Gains for %s was incorrect, got: %f, want: %f.", file, m.TradingGains, ExpectedTradingGains)
		 t.Logf("The length of the accounts was: %d", len(m.Accounts))
		 for k := range m.Accounts {
			 t.Logf("Account %s", k)
			 for idx, l := range m.Accounts[k].Transactions {
				 t.Logf("		[%d] Amount: %f Costbase: %f", idx,l.Amount, l.CostBase)
			 }
		 }
	}

	m.Close()
}

func TestDateMemoryStart(t *testing.T) {

	cfg := defaultTestConfig()

	m := *model.NewModel(cgtdb.NewMemDatabase(), cgtdb.NewMemDatabase(),cfg)

	m.Pricer.AddSinglePrice("2017-06-30","BTC","100","AUD")
	m.Pricer.AddSinglePrice("2017-07-01","BTC","100","AUD")
	m.Pricer.AddSinglePrice("2018-06-30","BTC","100","AUD")

	m.AddSingleTxn(*model.NewTransaction("2017-06-30T00:00:01.000000000","BTC","2","AUD","200"))
	m.AddSingleTxn(*model.NewTransaction("2017-07-01T00:00:00.000000001","AUD","150","BTC","1"))

	ExpectedGains := 50.0
	ExpectedTradingGains := 50.0

	m.ProcessModel()

	if ExpectedGains != m.CapitalGains {
		 t.Errorf("Expected Capital Gains was incorrect, got: %f, want: %f.", m.CapitalGains, ExpectedGains)
	}
	if ExpectedTradingGains != m.TradingGains {
		 t.Errorf("Expected Trading Gains was incorrect, got: %f, want: %f.", m.TradingGains, ExpectedTradingGains)
	}

	m.Close()
}
func TestDateMemoryBeforeEnd(t *testing.T) {
	cfg := defaultTestConfig()
	m := *model.NewModel(cgtdb.NewMemDatabase(), cgtdb.NewMemDatabase(),cfg)

	m.Pricer.AddSinglePrice("2017-06-29","BTC","100","AUD")
	m.Pricer.AddSinglePrice("2017-06-30","BTC","100","AUD")
	m.Pricer.AddSinglePrice("2017-07-01","BTC","100","AUD")
	m.Pricer.AddSinglePrice("2018-06-30","BTC","100","AUD")
	m.Pricer.AddSinglePrice("2018-07-01","BTC","100","AUD")

	m.AddSingleTxn(*model.NewTransaction("2017-06-29T00:00:01.000000000","BTC","2","AUD","200"))
	m.AddSingleTxn(*model.NewTransaction("2017-06-30T23:59:59.999999999","AUD","150","BTC","1"))

	ExpectedGains := 0.0
	ExpectedTradingGains := 0.0

	m.ProcessModel()

	if ExpectedGains != m.CapitalGains {
		 t.Errorf("Expected Capital Gains was incorrect, got: %f, want: %f.", m.CapitalGains, ExpectedGains)
	}
	if ExpectedTradingGains != m.TradingGains {
		 t.Errorf("Expected Trading Gains was incorrect, got: %f, want: %f.", m.TradingGains, ExpectedTradingGains)
	}

	m.Close()
}

func TestDateMemoryEnd(t *testing.T) {
	cfg := defaultTestConfig()
	m := *model.NewModel(cgtdb.NewMemDatabase(), cgtdb.NewMemDatabase(),cfg)

	m.Pricer.AddSinglePrice("2017-06-30","BTC","100","AUD")
	m.Pricer.AddSinglePrice("2017-07-01","BTC","100","AUD")
	m.Pricer.AddSinglePrice("2018-06-30","BTC","100","AUD")

	m.AddSingleTxn(*model.NewTransaction("2017-06-30T00:00:01.000000000","BTC","2","AUD","200"))
	m.AddSingleTxn(*model.NewTransaction("2018-06-30T23:59:59.999999999","AUD","150","BTC","1"))

	ExpectedGains := 50.0
	ExpectedTradingGains := 50.0

	m.ProcessModel()

	if ExpectedGains != m.CapitalGains {
		 t.Errorf("Expected Capital Gains was incorrect, got: %f, want: %f.", m.CapitalGains, ExpectedGains)
	}
	if ExpectedTradingGains != m.TradingGains {
		 t.Errorf("Expected Trading Gains was incorrect, got: %f, want: %f.", m.TradingGains, ExpectedTradingGains)
	}

	m.Close()
}
func TestDateMemoryAfterEnd(t *testing.T) {
	cfg := defaultTestConfig()
	m := *model.NewModel(cgtdb.NewMemDatabase(), cgtdb.NewMemDatabase(),cfg)

	m.Pricer.AddSinglePrice("2017-06-30","BTC","100","AUD")
	m.Pricer.AddSinglePrice("2017-07-01","BTC","100","AUD")
	m.Pricer.AddSinglePrice("2018-06-30","BTC","100","AUD")
	m.Pricer.AddSinglePrice("2018-07-01","BTC","100","AUD")

	m.AddSingleTxn(*model.NewTransaction("2017-06-30T00:00:01.000000000","BTC","2","AUD","200"))
	m.AddSingleTxn(*model.NewTransaction("2018-07-01T00:00:00.000000001","AUD","150","BTC","1"))

	ExpectedGains := 0.0
	ExpectedTradingGains := 0.0

	m.ProcessModel()

	if ExpectedGains != m.CapitalGains {
		 t.Errorf("Expected Capital Gains was incorrect, got: %f, want: %f.", m.CapitalGains, ExpectedGains)
	}
	if ExpectedTradingGains != m.TradingGains {
		 t.Errorf("Expected Trading Gains was incorrect, got: %f, want: %f.", m.TradingGains, ExpectedTradingGains)
	}

	m.Close()
}

func TestTradingMemoryBeforeEnd(t *testing.T) {
	cfg := defaultTestConfig()
	m := *model.NewModel(cgtdb.NewMemDatabase(), cgtdb.NewMemDatabase(),cfg)

	m.Pricer.AddSinglePrice("2017-06-29","BTC","100","AUD")
	m.Pricer.AddSinglePrice("2017-06-30","BTC","150","AUD")
	m.Pricer.AddSinglePrice("2017-07-01","BTC","150","AUD")
	m.Pricer.AddSinglePrice("2018-06-30","BTC","150","AUD")
	m.Pricer.AddSinglePrice("2018-07-01","BTC","150","AUD")

	m.AddSingleTxn(*model.NewTransaction("2017-06-29T00:00:01.000000000","BTC","2","AUD","200"))
	m.AddSingleTxn(*model.NewTransaction("2018-06-30T23:59:59.999999999","AUD","150","BTC","1"))

	ExpectedGains := 50.0
	ExpectedTradingGains := 0.0

	m.ProcessModel()

	if ExpectedGains != m.CapitalGains {
		 t.Errorf("Expected Capital Gains was incorrect, got: %f, want: %f.", m.CapitalGains, ExpectedGains)
	}
	if ExpectedTradingGains != m.TradingGains {
		 t.Errorf("Expected Trading Gains was incorrect, got: %f, want: %f.", m.TradingGains, ExpectedTradingGains)
	}

	m.Close()
}

func TestTradingMemoryNoTrades(t *testing.T) {
	cfg := defaultTestConfig()
	m := *model.NewModel(cgtdb.NewMemDatabase(), cgtdb.NewMemDatabase(),cfg)

	m.Pricer.AddSinglePrice("2017-06-29","BTC","100","AUD")
	m.Pricer.AddSinglePrice("2017-06-30","BTC","100","AUD")
	m.Pricer.AddSinglePrice("2017-07-01","BTC","150","AUD")
	m.Pricer.AddSinglePrice("2018-06-30","BTC","150","AUD")
	m.Pricer.AddSinglePrice("2018-07-01","BTC","150","AUD")

	m.AddSingleTxn(*model.NewTransaction("2017-06-29T00:00:01.000000000","BTC","1","AUD","100"))

	ExpectedGains := 0.0
	ExpectedTradingGains := 50.0

	m.ProcessModel()

	if ExpectedGains != m.CapitalGains {
		 t.Errorf("Expected Capital Gains was incorrect, got: %f, want: %f.", m.CapitalGains, ExpectedGains)
	}
	if ExpectedTradingGains != m.TradingGains {
		 t.Errorf("Expected Trading Gains was incorrect, got: %f, want: %f.", m.TradingGains, ExpectedTradingGains)
	}

	m.Close()
}

func TestAllDir(t *testing.T) {
	aud := accounting.Accounting{Precision: 2}

	cfg := defaultTestConfig()
	searchDir := "./testdata"

	fileList := []string{}
	err := filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {

			if f.IsDir() {
				return nil
			}

			fileList = append(fileList, path)
			return nil
	})
	if err != nil {
		t.Fatalf("failed reading directory: %s", err)
	}

	for _, file := range fileList {
			t.Log("------------------------------------------------------------")
			t.Log("Testing: " + file)
			g, err := ioutil.ReadFile(file)
			if err != nil {
				t.Fatalf("failed reading .golden: %s", err)
			}
			m := *model.NewModel(cgtdb.NewMemDatabase(), cgtdb.NewMemDatabase(),cfg)

			m.Pricer.AddSinglePrice("2018-06-17","BTC","6435.7","USD")
			m.Pricer.AddSinglePrice("2018-06-17","ETH","496.43","USD")
			m.Pricer.AddSinglePrice("2018-06-17","BCH","847.51","USD")
			m.Pricer.AddSinglePrice("2018-06-17","USD","1.3437","AUD")
			m.Pricer.AddSinglePrice("2018-06-20","BTC","6760.40","USD")
			m.Pricer.AddSinglePrice("2018-06-20","ETH","536.31","USD")
			m.Pricer.AddSinglePrice("2018-06-20","BCH","888.91","USD")
			m.Pricer.AddSinglePrice("2018-06-20","USD","1.3563","AUD")
			m.Pricer.AddSinglePrice("2018-06-29","BTC","6208.1","USD")
			m.Pricer.AddSinglePrice("2018-06-29","ETH","435.37","USD")
			m.Pricer.AddSinglePrice("2018-06-29","BCH","717.29","USD")
			m.Pricer.AddSinglePrice("2018-06-29","USD","1.3503","AUD")
			m.Pricer.AddSinglePrice("2018-06-30","BTC","6391.5","USD")
			m.Pricer.AddSinglePrice("2018-06-30","ETH","452.67","USD")
			m.Pricer.AddSinglePrice("2018-06-30","BCH","748.95","USD")
			m.Pricer.AddSinglePrice("2018-06-30","USD","1.3503","AUD")

			err = json.Unmarshal(g, &m)
			if err != nil {
				t.Fatalf("failed reading json: %s", err)
			}
			t.Log(m.EndDate)

			ExpectedGains := m.CapitalGains
			ExpectedTradingGains := m.TradingGains
			t.Logf("Number of Transactions Analysed: %d", len(m.Txns))
			t.Logf("Expected Capital Gains from Transactions: %s", aud.FormatMoney(ExpectedGains))
			m.ProcessModel()
			t.Logf("Actual Capital Gains Results: %s", aud.FormatMoney(m.CapitalGains))
			t.Logf("Expected Trading Gains from Transactions: %s", aud.FormatMoney(ExpectedTradingGains))
			t.Logf("Actual Trading Gains Results: %s", aud.FormatMoney(m.TradingGains))

			if ExpectedGains != m.CapitalGains {
				 t.Errorf("Expected Capital Gains for %s was incorrect, got: %f, want: %f.", file, m.CapitalGains, ExpectedGains)
			}
			if ExpectedTradingGains != m.TradingGains {
				 t.Errorf("Expected Trading Gains for %s was incorrect, got: %.15f, want: %.15f.", file, m.TradingGains, ExpectedTradingGains)
			}

			m.Close()
	}
}

func defaultTestConfig() *config.Config {
	cfg := config.DefaultConfig

	cfg.GainsDatabase = "/home/sean/.config/cgtcalc/testdb"
	cfg.EndFY = time.Date(2018, 6, 30, 23, 59, 59, 999999999, time.UTC)
	cfg.StartFY = time.Date(2017, 7, 1, 0, 0, 0, 0, time.UTC)
	return &cfg
}
