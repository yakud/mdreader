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
	Size          string
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

	csv[4] = r.Size
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
	record.Size = lotToDecimalQuantity(binary.LittleEndian.Uint64(buffer[sizeOffset : sizeOffset+sizeLength]), -8)

	//double   price; // 8
	record.Price = float64FromBytes(buffer[priceOffset : priceOffset+priceLength])

	return nil
}

func float64FromBytes(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}

const smallsString = "00010203040506070809" +
	"10111213141516171819" +
	"20212223242526272829" +
	"30313233343536373839" +
	"40414243444546474849" +
	"50515253545556575859" +
	"60616263646566676869" +
	"70717273747576777879" +
	"80818283848586878889" +
	"90919293949596979899"

func formatBits(u uint64, exp int) []byte {
	var a [40]byte
	i := len(a)

	if exp > 0 {
		i -= exp
		for j := 0; j < exp; j++ {
			a[i+j] = '0'
		}
	}

	p := 0
	for u >= 100 {
		is := u % 100 * 2
		u /= 100
		i -= 2

		if p == exp {
			a[i+1] = '.'
			i--
		}
		p--
		a[i+1] = smallsString[is+1]
		if p == exp {
			a[i] = '.'
			i--
		}
		p--
		a[i+0] = smallsString[is+0]
	}

	// us < 100
	is := u * 2
	i--
	if p == exp {
		a[i] = '.'
		i--
	}
	p--
	a[i] = smallsString[is+1]
	if u >= 10 {
		i--
		if p == exp {
			a[i] = '.'
			i--
		}
		p--
		a[i] = smallsString[is]
	}
	for p >= exp {
		i--
		if p == exp {
			a[i] = '.'
			i--
		}
		p--
		a[i] = '0'
	}
	return a[i:]
}

func lotToDecimalQuantity(lot uint64, exp int) string {
	if lot == 0 {
		return "0"
	}
	return string(formatBits(lot, exp))
}

func iLotToQuantity(lot int64, exp int) string {
	if lot == 0 {
		return "0"
	}
	if lot > 0 {
		return string(formatBits(uint64(lot), exp))
	} else {
		return "-" + string(formatBits(uint64(-lot), exp))
	}
}