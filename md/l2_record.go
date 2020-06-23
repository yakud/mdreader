package md

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"
)

type RecordType uint8

const (
	RecordTypeReset  RecordType = 0
	RecordTypeUpdate RecordType = 1
	RecordTypeTrade  RecordType = 2
)

type SideType byte

const (
	SideTypeBid         SideType = 'B'
	SideTypeAsk         SideType = 'A'
	SideTypeUnknown     SideType = 'U'
	SideTypeUnknownZero SideType = 0
)

const (
	tsExchangeUTCOffset = 0
	tsExchangeUTCLength = 8

	tsReceiveUTCOffset = tsExchangeUTCLength
	tsReceiveUTCLength = 8

	typeOffset = tsReceiveUTCOffset + tsReceiveUTCLength
	typeLength = 1

	sideOffset = typeOffset + typeLength
	sideLength = 1

	sizeOffset = sideOffset + sideLength
	sizeLength = 8

	priceOffset = sizeOffset + sizeLength
	priceLength = 8

	L2HeaderSize = 16
	L2RecordSize = tsExchangeUTCLength + tsReceiveUTCLength + typeLength + sideLength + sizeLength + priceLength

	timeOffsetYears = 30
)

var emptyRecord = L2Record{}

type L2Record struct {
	TsExchangeUTC time.Time
	TsReceiveUTC  time.Time
	Type          uint8
	Side          uint8
	Size          uint64
	Price         float64
}

func (r *L2Record) ToScvHeaderLine() []string {
	return []string{
		"TsExchangeUTC",
		"TsReceiveUTC",
		"Type",
		"Side",
		"Size",
		"Price",
	}
}
func (r *L2Record) ToScvLine() []string {
	csv := make([]string, 6)
	csv[0] = r.TsExchangeUTC.Format(time.RFC3339Nano)
	csv[1] = r.TsReceiveUTC.Format(time.RFC3339Nano)

	switch RecordType(r.Type) {
	case RecordTypeReset:
		csv[2] = "reset"
	case RecordTypeUpdate:
		csv[2] = "update"
	case RecordTypeTrade:
		csv[2] = "trade"
	default:
		log.Fatalf("undefined type: %d", r.Type)
	}

	switch SideType(r.Side) {
	case SideTypeBid:
		csv[3] = "bid"
	case SideTypeAsk:
		csv[3] = "ask"
	case SideTypeUnknown, SideTypeUnknownZero:
		csv[3] = "unknown"
	default:
		log.Fatalf("undefined side: %d", r.Side)
	}

	csv[4] = strconv.FormatUint(r.Size, 10)
	csv[5] = strconv.FormatFloat(r.Price, 'g', -1, 64)
	return csv
}

func (r *L2Record) Clean() {
	*r = emptyRecord
}

func ConvertBytesToL2Record(buffer []byte, record *L2Record) error {
	if len(buffer) != L2RecordSize {
		return fmt.Errorf("expected buffer size: %d actual: %d", L2RecordSize, len(buffer))
	}

	// tsExchangeUTC
	tsExchangeUTCUInt := binary.LittleEndian.Uint64(buffer[tsExchangeUTCOffset : tsExchangeUTCOffset+tsExchangeUTCLength])
	record.TsExchangeUTC = time.Unix(0, int64(time.Duration(tsExchangeUTCUInt))).AddDate(timeOffsetYears, 0, 0)

	// tsReceiveUTC
	tsReceiveUTC := buffer[tsReceiveUTCOffset : tsReceiveUTCOffset+tsReceiveUTCLength]
	tsReceiveUTCUint := binary.LittleEndian.Uint64(tsReceiveUTC)
	record.TsReceiveUTC = time.Unix(0, int64(tsReceiveUTCUint)).AddDate(timeOffsetYears, 0, 0)

	// type
	record.Type = buffer[typeOffset]

	// side
	record.Side = buffer[sideOffset]

	// size
	record.Size = binary.LittleEndian.Uint64(buffer[sizeOffset : sizeOffset+sizeLength])

	//double   price; // 8
	record.Price = float64FromBytes(buffer[priceOffset : priceOffset+priceLength])

	return nil
}

func float64FromBytes(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}
