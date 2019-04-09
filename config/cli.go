package config

import (
	"fmt"

	"github.com/naoina/toml"
	"github.com/urfave/cli"
)

var (
	DumpConfigCommand = cli.Command{
		Action: 	dumpConfig,
		Name: 		"dumpconfig",
		Flags: []cli.Flag{
			ConfigFileFlag,
		},
		Usage: 		"Show configuration values",
		Category:	"MISCELLANEOUS COMMANDS",
		Description: `The dumpconfig command shows configuration values.`,
	}

	ConfigFileFlag = cli.StringFlag{
		Name:		"config",
		Usage:	"TOML configuration file",
	}
)

var NewCommand = cli.Command{
	Name:      "new",
	Usage:     "new config file",
	Flags: []cli.Flag{
	},
	Action: func(ctx *cli.Context) error {
		entity := ctx.Args().First()

		cfg := NewConfig()
		cfg.Name = entity

		configfile := defaultConfiglocation()
		cfg.Save(configfile)

		return nil
	},
}



func dumpConfig(ctx *cli.Context) error {
	cfg, err := LoadConfig(ctx.String("configFileFlag"))
	if err != nil {
		return err
	}

	fmt.Println("Dumping Config File")
	fmt.Println("-------------------")
	output, _ := toml.Marshal(&cfg)
	fmt.Printf("\n%s", output)

	return nil
}
