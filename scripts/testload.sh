#!/bin/bash

rm logfile.log
rm transactions_list.txt

cgtcalc del

cgtcalc load history.csv -exchange cointracking
cgtcalc dump > transactions_list.txt
