package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/yakud/marketdata/md"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	in := flag.String("in", "", "./BTCUSD.dat.uncompressed")
	//out := flag.String("out", "", "./BTCUSD.csv")
	xmlPath := flag.String("xml", "", "./Bitfinex.xml")
	flag.Parse()

	xmlFile, err := os.Open(*xmlPath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully Opened file")
	defer xmlFile.Close()

	byteValue, _ := ioutil.ReadAll(xmlFile)

	var session md.Session
	if err := xml.Unmarshal(byteValue, &session); err != nil {
		log.Println(err)
	} else {
		log.Println("Unmarshaled")
	}

	exponents := md.NewExponents()
	for _, instrument := range session.DumpInstruments.Instruments {
		isInSplit := strings.Split(instrument.IsIn, "||")
		if len(isInSplit) != 2 {
			continue
		}
		prov := strings.TrimSuffix(isInSplit[1], ".MD")
		pr := instrument.Details.BaseCurrency + instrument.Details.QuoteCurrency
		exp := instrument.Details.ExchangeLotSize

		if err := exponents.SetString(prov, pr, exp); err != nil {
			log.Println("An error occured when trying to add value to exponents map. Probably problem with parsing exponent from scientific notation.")
			log.Println("error: ", err)
			continue
		}
	}

	providerCitySplit:= strings.Split(filepath.Base(filepath.Dir(*in)), ".")
	if len(providerCitySplit) != 2 {
		log.Fatal("in file's parents dir should be if provider.city format! Example: Bitfinex.london")
	}

	provider := providerCitySplit[0]
	city := providerCitySplit[1]

	var pairRegexp = regexp.MustCompile(`(?m)(_.*)?\..+`)
	pair := pairRegexp.ReplaceAllString(filepath.Base(*in), "")

	//dataFile, err := os.Open(*in)
	//if err != nil {
	//	log.Fatal(err)
	//}

	log.Println(exponents.Get(provider, pair))
	log.Println(city)
}
