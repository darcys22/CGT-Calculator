package exchanges

import (
	"errors"
	"strings"

	"cgtcalc/model"
	//log "github.com/Sirupsen/logrus"
)

type Processor func(string) []model.Transaction

func ExchangeFuncSearch(exchange string) (Processor, error) {
	var f Processor

	switch strings.ToUpper(exchange) {
	case "BIBOX":
		f = BiboxFile
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
	case "COINEXCHANGE":
		f = CoinexchangeFile
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
	case "HITBTC":
		f = HitbtcFile
	case "KRAKEN":
		f = KrakenFile
	case "KUCOIN":
		f = KuCoinFile
	case "IR":
		f = IndependantReserveFile
	case "MERCATOX":
		f = MercatoxFile
	case "POLONIEX":
		f = PoloniexFile
	case "CUSTOM":
		f = CustomFile
	default:
		return nil, errors.New("Exchange File Not Implemented")
	}
	return f, nil
}
