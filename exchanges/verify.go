package exchanges

import (

	"errors"
	"strings"

	"cgtcalc/model"

	//log "github.com/Sirupsen/logrus"
)

type Processor func(string) []model.Transaction

func ExchangeFuncSearch(exchange string) (Processor, error){
	var f Processor

	switch strings.ToUpper(exchange) {
		case "BINANCE":
			f = BinanceFile
		case "BITFINEX":
			f = BitfinexFile
		case "BITMEX":
			f = BitmexFile
		case "BITPANDA":
			f = BitpandaFile
		case "BITTREX":
			f = BittrexFile
		case "BTCMARKETS":
			f = BtcMarketsFile
		case "COINSPOT":
			f = CoinspotFile
		case "COINTRACKING":
			f = CointrackingFile
		case "COSS":
			f = CossFile
		case "CRYPTOPIA":
			f = CryptopiaFile
		case "ETHERDELTA":
			f = EtherDeltaFile
		case "ETHFINEX":
			f = BitfinexFile
		case "KRAKEN":
			f = KrakenFile
		case "KUCOIN":
			f = KuCoinFile
		case "IR":
			f = IndependantReserveFile
		case "POLONIEX":
			f = PoloniexFile
		case "CUSTOM":
			f = CustomFile
		default:
			return nil, errors.New("Exchange File Not Implemented")
		}
	return f,nil
}
