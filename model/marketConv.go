package model

import (
	"strings"
)

func MarketConverter(market, typ string) (string, string) {
	var x, y string
	if strings.Index(strings.ToUpper(market), "/") >= 0 {
		arr := strings.Split(strings.ToUpper(market), "/")
		x = arr[0]
		y = arr[1]
	} else if strings.Index(strings.ToUpper(market), "-") >= 0 {
		arr := strings.Split(strings.ToUpper(market), "-")
		x = arr[1]
		y = arr[0]
	} else {
		if strings.ToUpper(market)[len(market)-1] == 'T' {
			x = strings.ToUpper(market)[:len(market)-4]
			y = strings.ToUpper(market)[len(market)-4:]
		} else {
			x = strings.ToUpper(market)[:len(market)-3]
			y = strings.ToUpper(market)[len(market)-3:]
		}
	}
	if strings.ToUpper(typ) == "BUY" {
		return y, x
	} else {
		return x, y
	}
}

func InvMarketConverter(market, typ string) (string, string) {
	var x, y string
	if strings.Index(strings.ToUpper(market), "-") >= 0 {
		arr := strings.Split(strings.ToUpper(market), "-")
		x = arr[0]
		y = arr[1]
	} else {
		if strings.ToUpper(market)[len(market)-1] == 'T' {
			x = strings.ToUpper(market)[:len(market)-4]
			y = strings.ToUpper(market)[len(market)-4:]
		} else {
			x = strings.ToUpper(market)[:len(market)-3]
			y = strings.ToUpper(market)[len(market)-3:]
		}
	}
	if strings.ToUpper(typ) == "BUY" {
		return y, x
	} else {
		return x, y
	}
}

//First Returned is the Base Currency, Second is Quote Currency
//out.BaseCurrency, out.QuoteCurrency = model.MarketConverter(txn.Market, txn.Type)

//var name = []string{
//"ETHBTC",
//"ETHBTC",
//"LTCBTC",
//"BNBBTC",
//"NEOBTC",
//"QTUMETH",
//"EOSETH",
//"SNTETH",
//"BNTETH",
//"BCHBTC",
//"GASBTC",
//"BNBETH",
//"BTCUSDT",
//"ETHUSDT",
//"OAXETH",
//"DNTETH",
//"MCOETH",
//"MCOBTC",
//"WTCBTC",
//"WTCETH",
//"LRCBTC",
//"LRCETH",
//"QTUMBTC",
//"YOYOWBTC",
//"OMGBTC",
//"OMGETH",
//"ZRXBTC",
//"ZRXETH",
//"STRATBTC",
//"STRATETH",
//"SNGLSBTC",
//"SNGLSETH",
//"ETHOSBTC",
//"ETHOSETH",
//"KNCBTC",
//"KNCETH",
//"FUNBTC",
//"FUNETH",
//"SNMBTC",
//"SNMETH",
//"NEOETH",
//"MIOTABTC",
//"MIOTAETH",
//"LINKBTC",
//"LINKETH",
//"XVGBTC",
//"XVGETH",
//"SALTBTC",
//"SALTETH",
//"MDABTC",
//"MDAETH",
//"MTLBTC",
//"MTLETH",
//"SUBBTC",
//"SUBETH",
//"EOSBTC",
//"SNTBTC",
//"ETCETH",
//"ETCBTC",
//"MTHBTC",
//"MTHETH",
//"ENGBTC",
//"ENGETH",
//"DNTBTC",
//"ZECBTC",
//"ZECETH",
//"BNTBTC",
//"ASTBTC",
//"ASTETH",
//"DASHBTC",
//"DASHETH",
//"OAXBTC",
//"BTGBTC",
//"BTGETH",
//"EVXBTC",
//"EVXETH",
//"REQBTC",
//"REQETH",
//"VIBBTC",
//"VIBETH",
//"TRXBTC",
//"TRXETH",
//"POWRBTC",
//"POWRETH",
//"ARKBTC",
//"ARKETH",
//"YOYOWETH",
//"XRPBTC",
//"XRPETH",
//"MODBTC",
//"MODETH",
//"ENJBTC",
//"ENJETH",
//"STORJBTC",
//"STORJETH",
//"BNBUSDT",
//"YOYOWBNB",
//"POWRBNB",
//"KMDBTC",
//"KMDETH",
//"NULSBNB",
//"RCNBTC",
//"RCNETH",
//"RCNBNB",
//"NULSBTC",
//"NULSETH",
//"RDNBTC",
//"RDNETH",
//"RDNBNB",
//"XMRBTC",
//"XMRETH",
//"DLTBNB",
//"WTCBNB",
//"DLTBTC",
//"DLTETH",
//"AMBBTC",
//"AMBETH",
//"AMBBNB",
//"BCHETH",
//"BCHUSDT",
//"BCHBNB",
//"BATBTC",
//"BATETH",
//"BATBNB",
//"BCPTBTC",
//"BCPTETH",
//"BCPTBNB",
//"ARNBTC",
//"ARNETH",
//"GVTBTC",
//"GVTETH",
//"CDTBTC",
//"CDTETH",
//"GXSBTC",
//"GXSETH",
//"NEOUSDT",
//"NEOBNB",
//"POEBTC",
//"POEETH",
//"QSPBTC",
//"QSPETH",
//"QSPBNB",
//"BTSBTC",
//"BTSETH",
//"BTSBNB",
//"XZCBTC",
//"XZCETH",
//"XZCBNB",
//"LSKBTC",
//"LSKETH",
//"LSKBNB",
//"TNTBTC",
//"TNTETH",
//"FUELBTC",
//"FUELETH",
//"MANABTC",
//"MANAETH",
//"BCDBTC",
//"BCDETH",
//"DGDBTC",
//"DGDETH",
//"MIOTABNB",
//"ADXBTC",
//"ADXETH",
//"ADXBNB",
//"ADABTC",
//"ADAETH",
//"PPTBTC",
//"PPTETH",
//"CMTBTC",
//"CMTETH",
//"CMTBNB",
//"XLMBTC",
//"XLMETH",
//"XLMBNB",
//"CNDBTC",
//"CNDETH",
//"CNDBNB",
//"LENDBTC",
//"LENDETH",
//"WABIBTC",
//"WABIETH",
//"WABIBNB",
//"LTCETH",
//"LTCUSDT",
//"LTCBNB",
//"TNBBTC",
//"TNBETH",
//"WAVESBTC",
//"WAVESETH",
//"WAVESBNB",
//"GTOBTC",
//"GTOETH",
//"GTOBNB",
//"ICXBTC",
//"ICXETH",
//"ICXBNB",
//"OSTBTC",
//"OSTETH",
//"OSTBNB",
//"ELFBTC",
//"ELFETH",
//"AIONBTC",
//"AIONETH",
//"AIONBNB",
//"NEBLBTC",
//"NEBLETH",
//"NEBLBNB",
//"BRDBTC",
//"BRDETH",
//"BRDBNB",
//"MCOBNB",
//"EDOBTC",
//"EDOETH",
//"WINGSBTC",
//"WINGSETH",
//"NAVBTC",
//"NAVETH",
//"NAVBNB",
//"LUNBTC",
//"LUNETH",
//"APPCBTC",
//"APPCETH",
//"APPCBNB",
//"VIBEBTC",
//"VIBEETH",
//"RLCBTC",
//"RLCETH",
//"RLCBNB",
//"INSBTC",
//"INSETH",
//"PIVXBTC",
//"PIVXETH",
//"PIVXBNB",
//"IOSTBTC",
//"IOSTETH",
//"STEEMBTC",
//"STEEMETH",
//"STEEMBNB",
//"NANOBTC",
//"NANOETH",
//"NANOBNB",
//"VIABTC",
//"VIAETH",
//"VIABNB",
//"BLZBTC",
//"BLZETH",
//"BLZBNB",
//"AEBTC",
//"AEETH",
//"AEBNB",
//"NCASHBTC",
//"NCASHETH",
//"NCASHBNB",
//"POABTC",
//"POAETH",
//"POABNB",
//"ZILBTC",
//"ZILETH",
//"ZILBNB",
//"ONTBTC",
//"ONTETH",
//"ONTBNB",
//"STORMBTC",
//"STORMETH",
//"STORMBNB",
//"QTUMBNB",
//"QTUMUSDT",
//"XEMBTC",
//"XEMETH",
//"XEMBNB",
//"WANBTC",
//"WANETH",
//"WANBNB",
//"WPRBTC",
//"WPRETH",
//"QLCBTC",
//"QLCETH",
//"SYSBTC",
//"SYSETH",
//"SYSBNB",
//"QLCBNB",
//"GRSBTC",
//"GRSETH",
//"ADAUSDT",
//"ADABNB",
//"CLOAKBTC",
//"CLOAKETH",
//"GNTBTC",
//"GNTETH",
//"GNTBNB",
//"LOOMBTC",
//"LOOMETH",
//"LOOMBNB",
//"XRPUSDT",
//"REPBTC",
//"REPETH",
//"REPBNB",
//"TUSDBTC",
//"TUSDETH",
//"TUSDBNB",
//"ZENBTC",
//"ZENETH",
//"ZENBNB",
//"SKYBTC",
//"SKYETH",
//"SKYBNB",
//"EOSUSDT",
//"EOSBNB",
//"CVCBTC",
//"CVCETH",
//"CVCBNB",
//"THETABTC",
//"THETAETH",
//"THETABNB",
//"XRPBNB",
//"TUSDUSDT",
//"MIOTAUSDT",
//"XLMUSDT",
//"IOTXBTC",
//"IOTXETH",
//"QKCBTC",
//"QKCETH",
//"AGIBTC",
//"AGIETH",
//"AGIBNB",
//"NXSBTC",
//"NXSETH",
//"NXSBNB",
//"ENJBNB",
//"DATABTC",
//"DATAETH",
//"ONTUSDT",
//"TRXBNB",
//"TRXUSDT",
//"ETCUSDT",
//"ETCBNB",
//"ICXUSDT",
//"SCBTC",
//"SCETH",
//"SCBNB",
//"NPXSBTC",
//"NPXSETH",
//"KEYBTC",
//"KEYETH",
//"NASBTC",
//"NASETH",
//"NASBNB",
//"MFTBTC",
//"MFTETH",
//"MFTBNB",
//"DENTBTC",
//"DENTETH",
//"ARDRBTC",
//"ARDRETH",
//"ARDRBNB",
//"NULSUSDT",
//"HOTBTC",
//"HOTETH",
//"VETBTC",
//"VETETH",
//"VETUSDT",
//"VETBNB",
//"DOCKBTC",
//"DOCKETH",
//"POLYBTC",
//"POLYBNB",
//"PHXBTC",
//"PHXETH",
//"PHXBNB",
//"HCBTC",
//"HCETH",
//"GOBTC",
//"GOBNB",
//"PAXBTC",
//"PAXBNB",
//"PAXUSDT",
//"PAXETH",
//"RVNBTC",
//"RVNBNB",
//}
