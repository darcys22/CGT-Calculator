package reporting

import (
	"fmt"
	"os"
	"bufio"

	"cgtcalc/model"
	//"cgtcalc/version"

	//"github.com/leekchan/accounting"
	log "github.com/Sirupsen/logrus"
)

func ExportTraces(m *model.Model) {

	file, err := os.Create("transactions_trace.txt")
	if err != nil {
		log.Fatal("Couldnt create the txt File to export", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for idx, _ := range m.Txns {
		writer.WriteString(fmt.Sprintf("#### Transaction #%d\n", m.Txns[idx].Nonce))
		writer.Write(m.Txns[idx].Trace)
		writer.WriteString(fmt.Sprintf("\n"))
	}
}
