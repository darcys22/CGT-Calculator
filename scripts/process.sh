#!/bin/bash

#rm logfile.log
#rm outputcalculation.pdf
#rm result.csv
#rm transactions_list.txt
#rm transactions_trace.txt

rm -rf output
rm output.zip
mkdir output

cgtcalc del

cgtcalc load BinanceTradeHistory.xlsx -exchange binance
cgtcalc load BTCMarketsALL.csv -exchange btcmarkets
cgtcalc load Coinbase.csv -exchange custom
cgtcalc load Coinexchange.csv -exchange custom
cgtcalc load Cryptopia.csv -exchange cryptopia
cgtcalc load KuCoin.csv -exchange kucoin
cgtcalc load "IR Coins.csv" -exchange ir

cgtcalc process
cgtcalc dump > output/transactions_list.txt

cp -rf ~/.config/cgtcalc/gainsdb ./output/
cp -rf ~/.config/cgtcalc/pricedb ./output/
cp ~/.config/cgtcalc/cgt.conf ./output/

mv outputcalculation.pdf ./output/
mv result.csv ./output/
mv logfile.log ./output/
mv transactions_trace.txt ./output/

zip -rq output output/*
#rm -rf output
