package main

import (
	//"fmt"

	//"cgtcalc/exchanges"
	//"cgtcalc/model"
	//"cgtcalc/config"
	//"cgtcalc/utils"

	"github.com/urfave/cli"

	log "github.com/Sirupsen/logrus"

)

const (
)

var debugCommand = cli.Command{
	Name:      "debug",
	Usage:     "for testing",
	Category:	"MISCELLANEOUS COMMANDS",
	Flags: []cli.Flag{
	},
	Action: func(ctx *cli.Context) error {
		//////
		log.Warn("TESTING DEBUG SHIT\n")

		//config.TestWiz()

		//////
		return nil
	},
}

