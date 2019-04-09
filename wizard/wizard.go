// Copyright 2017 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

package wizard

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"cgtcalc/config"
	"cgtcalc/exchanges"

	"github.com/urfave/cli"
	log "github.com/Sirupsen/logrus"
)

func WizardCommand(ctx *cli.Context) error {
	cfg, err := config.LoadConfig(ctx.String("configFileFlag"))
	if err != nil {
		cfg = config.NewConfig()
	}
	makeWizard(cfg).run()
	return nil
}

type wizard struct {

	conf    *config.Config // Configurations from previous runs
	in   *bufio.Reader // Wrapper around stdin to allow reading user input
}

// read reads a single line from stdin, trimming if from spaces.
func (w *wizard) read() string {
	fmt.Printf("> ")
	text, err := w.in.ReadString('\n')
	if err != nil {
		fmt.Println("Failed to read user input", "err", err)
	}
	return strings.TrimSpace(text)
}

// readString reads a single line from stdin, trimming if from spaces, enforcing
// non-emptyness.
func (w *wizard) readString() string {
	for {
		fmt.Printf("> ")
		text, err := w.in.ReadString('\n')
		if err != nil {
			fmt.Println("Failed to read user input", "err", err)
		}
		if text = strings.TrimSpace(text); text != "" {
			return text
		}
	}
}

// makeWizard creates and returns a new wizard.
func makeWizard(c *config.Config) *wizard {

	return &wizard{
		conf:			c,
		in:       bufio.NewReader(os.Stdin),
	}
}

// run displays some useful infos to the user, starting on the journey of
// setting up a new or managing an existing config file.
func (w *wizard) run() {
	fmt.Println("+-----------------------------------------------------------+")
	fmt.Println("| Welcome to CGT Calculator                                 |")
	fmt.Println("|                                                           |")
	fmt.Println("| This tool analyses your exchange trade history and        |")
	fmt.Println("| will output the total Capital Gains and/or Trading Gains. |")
	fmt.Println("|                                                           |")
	fmt.Println("| You will need to first create a new file and specify which|")
	fmt.Println("| files are to be analysed and which exchange they came     |")
	fmt.Println("| from. Following this the program does preprocessing to    |")
	fmt.Println("| identify duplicates, shortages and possible errors.       |")
	fmt.Println("|                                                           |")
	fmt.Println("| Finally the program will process the files and output     |")
	fmt.Println("| various pdf, csv and log files that can be added to the   |")
	fmt.Println("| client workpapers.                                        |")
	fmt.Println("+-----------------------------------------------------------+")
	fmt.Println()


	fmt.Println(w.conf.GainsDatabase)
	if w.conf.Name == "" {
		fmt.Println("Please specify the name of the entity for the Gains Calculation (ie Company Pty Ltd)")
		w.conf.Name = w.readString()
		fmt.Println("Using: " + w.conf.Name)
	}


	// Basics done, loop ad infinitum about what to do
Outerloop:
	for {
		fmt.Println()
		fmt.Println("What would you like to do? (default = 1)")
		fmt.Println(" 1. Set Exchange History Files to be analysed")
		fmt.Println(" 2. Preprocess the files ")
		fmt.Println(" 3. Run CGT Calculator")
		fmt.Println(" 4. Exit")

		choice := w.read()
		switch {
		case choice == "" || choice == "1":
			w.setFiles()

		case choice == "2":
			w.preprocess()
		case choice == "3":
			w.process()
		case choice == "4":
			break Outerloop

		default:
			fmt.Println("That's not something I can do")
		}
	}

	w.conf.Save("")
}

type File struct {
	Filename    string
	Exchange		string
}

func (w *wizard) setFiles() {
	var list []File

	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		var found File
		found.Filename = f.Name()
		if val, ok := w.conf.Files[f.Name()]; ok {
			found.Exchange = val
		} else {
			found.Exchange = ""
		}
		list = append(list, found)
	}
Outerloop:
	for {
		fmt.Println()
		fmt.Println("Select a file to describe the exchange (Select 0 to exit)")
		fmt.Println()
		for i, n := range list {
			fmt.Printf("%d) %s: %s\n", i+1, n.Filename, n.Exchange)
		}

		choice := w.read()
		i, _ := strconv.Atoi(choice)
		if choice == "" || choice == "0" {
			break Outerloop
		} else if i > len(list) {
			fmt.Println("That's not something I can do")
			continue
		}

		fmt.Println("type the name of the exchange")
		exchange := w.read()

		//validate and verify
		_, err := exchanges.ExchangeFuncSearch(exchange)
		if err != nil {
			fmt.Printf("That exchange is not implemented")
		} else {
			list[i-1].Exchange = strings.ToUpper(exchange)
		}

	}

	m := make(map[string]string)
	for _, n := range list {
		m[n.Filename] = n.Exchange
	}
	w.conf.Files = m

	w.conf.Save("")

}
func (w *wizard) preprocess() {
	// 1) Run the files that have been named and not yet in system
	// 2) identify duplicates
	// 3) deal with duplicates

Outerloop:
	for k, v := range w.conf.Files {
		if v == "" {
			continue Outerloop
		}
		cmd := exec.Command("cgtcalc", "load", k, "-exchange", v)
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			log.Warnf("cmd.Run() failed with %s\n", err)
		}
	}
}

func (w *wizard) preclean() {
}
func (w *wizard) process() {
	cmd := exec.Command("cgtcalc", "process")
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
}
func (w *wizard) postclean() {
}
