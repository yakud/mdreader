package md

import (
	"encoding/xml"
)

type Session struct {
	XMLName         xml.Name        `xml:"Session"`
	ServerId        string          `xml:"ServerId"`
	Date            string          `xml:"Date"`
	Revision        int             `xml:"Revision"`
	DumpInstruments DumpInstruments `xml:"DumpInstruments"`
}

type DumpInstruments struct {
	XMLName     xml.Name      `xml:"DumpInstruments"`
	Instruments []*Instrument `xml:"Instrument"`
}

type Instrument struct {
	XMLName  xml.Name          `xml:"Instrument"`
	IsIn     string            `xml:"Isin,attr"`
	Id       string            `xml:"Id,attr"`
	DumpType string            `xml:"DumpType,attr"`
	Details  InstrumentDetails `xml:"Details"`
}

type InstrumentDetails struct {
	XMLName         xml.Name  `xml:"Details"`
	Type            string    `xml:"Type,attr"`
	QuoteCurrency   string    `xml:"QuoteCurrency,attr"`
	BaseCurrency    string    `xml:"BaseCurrency,attr"`
	AuxCurrency     string    `xml:"AuxCurrency,attr"`
	PlatformLotSize string    `xml:"PlatformLotSize,attr"`
	ExchangeLotSize string    `xml:"ExchangeLotSize,attr"`
	OrderSizeStep   string    `xml:"OrderSizeStep,attr"`
	PriceStep       PriceStep `xml:"PriceStep"`
}

type PriceStep struct {
	XMLName xml.Name `xml:"PriceStep"`
	Type    string   `xml:"Type,attr"`
	Data    string   `xml:"Data,attr"`
}
