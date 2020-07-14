package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/yakud/marketdata/md"
)

func main() {
	in := flag.String("in", "", "./BTCUSD.dat.uncompressed")
	out := flag.String("out", "", "./BTCUSD.csv")
	flag.Parse()

	buffer, err := ioutil.ReadFile(*in)
	if err != nil {
		log.Fatal(err)
	}

	csvResult, err := os.Create(*out)
	if err != nil {
		log.Fatal(err)
	}
	defer csvResult.Close()

	writer := csv.NewWriter(csvResult)
	defer writer.Flush()

	var offset = 16 // skip header
	var record = &md.L2Record{}

	fmt.Println("L2RecordSize:", md.L2RecordSize)

	if err := writer.Write(record.ToScvHeaderLine()); err != nil {
		log.Fatal(err)
	}

	fmt.Println("start reading:", *in)

	lines := 0
	records := make([][]string, 0)
	for {
		if offset+md.L2RecordSize > len(buffer) {
			fmt.Println("end read")
			break
		}
		if err := md.ConvertBytesToL2Record(buffer[offset:offset+md.L2RecordSize], record); err != nil {
			log.Fatal(err)
		}
		offset += md.L2RecordSize

		records = append(records, record.ToScvLine())

		record.Clean()
		lines++
	}

	if err := writer.WriteAll(records); err != nil {
		log.Fatal(err)
	}

	fmt.Println("total lines:", lines)
	fmt.Println("saved to:", *out)
}
