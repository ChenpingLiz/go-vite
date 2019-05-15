package dex

import (
	"encoding/binary"
	"fmt"
	"github.com/vitelabs/go-vite/vm_db"
	"strconv"
	"strings"
)

const PriceBytesLength = 10

//MarketId[0..2]Side[3]Price[4..13]timestamp[14..18]serialNo[19..21] = 22
func ComposeOrderId(db vm_db.VmDb,  marketId int32, param *ParamDexFundNewOrder) (idBytes []byte) {
	idBytes = make([]byte, OrderIdBytesLength)
	copy(idBytes[:3], Uint32ToBytes(uint32(marketId))[1:])
	if param.Side {
		idBytes[3] = 1
	}
	priceBytes := PriceToBytes(param.Price)
	if !param.Side { // buy
		BitwiseNotBytes(priceBytes)
	}
	copy(idBytes[4:14], priceBytes)
	timestamp := GetTimestampInt64(db)
	copy(idBytes[14:19], Uint64ToBytes(uint64(timestamp))[3:])
	serialNo := NewAndSaveOrderSerialNo(db, timestamp)
	copy(idBytes[19:], Uint32ToBytes(uint32(serialNo))[1:])
	return
}

func DeComposeOrderId(idBytes []byte) (marketId int32, side bool, price[]byte, timestamp int64, err error) {
	if len(idBytes) != OrderIdBytesLength {
		err = DeComposeOrderIdFailErr
		return
	}
	marketIdBytes := make([]byte, 4)
	copy(marketIdBytes[1:4], idBytes[:3])
	marketId = int32(BytesToUint32(marketIdBytes))
	side = int8(idBytes[3]) == 1
	price = make([]byte, 10)
	copy(price[:], idBytes[4:14])
	if !side { // buy
		BitwiseNotBytes(price)
	}
	timestampBytes := make([]byte, 8)
	timestamp = int64(BytesToUint64(timestampBytes))
	return
}

func PriceToBytes(price string) []byte {
	parts := strings.Split(price, ".")
	var intPart, decimalPart string
	priceBytes := make([]byte, PriceBytesLength)
	if len(parts) == 2 {
		intPart = parts[0]
		decimalPart = parts[1]
	} else {
		intPart = parts[0]
	}
	if len(intPart) > 0 {
		intValue, _ := strconv.ParseUint(intPart, 10, 64)
		bs := make([]byte, 8)
		binary.BigEndian.PutUint64(bs, intValue)
		copy(priceBytes[:5], bs[3:8])
	}
	decimalLen := len(decimalPart)
	if decimalLen > 0 {
		if decimalLen < priceDecimalMaxLen {
			decimalPartArr := make([]byte, priceDecimalMaxLen)
			copy(decimalPartArr, decimalPart)
			//right pad 0
			for i := decimalLen; i < priceDecimalMaxLen; i++ {
				decimalPartArr[i] = '0'
			}
			decimalPart = string(decimalPartArr)
		}
		decimalValue, _ := strconv.ParseUint(decimalPart, 10, 64)
		bs := make([]byte, 8)
		binary.BigEndian.PutUint64(bs, decimalValue)
		copy(priceBytes[5:], bs[3:8])
	}
	return priceBytes
}

func BytesToPrice(priceBytes []byte) string {
	intBytes := make([]byte, 8)
	copy(intBytes[3:], priceBytes[:5])
	intValue := binary.BigEndian.Uint64(intBytes)
	decimalBytes := make([]byte, 8)
	copy(decimalBytes[3:], priceBytes[5:])
	decimalValue := binary.BigEndian.Uint64(decimalBytes)
	var intStr, decimalStr string
	if intValue == 0 {
		intStr = "0"
	} else {
		intStr = strconv.FormatUint(intValue, 10)
	}
	if decimalValue == 0 {
		return intStr
	} else {
		decimalStr = strconv.FormatUint(decimalValue, 10)
		decimalLen := len(decimalStr)
		decimalPartArr := make([]byte, priceDecimalMaxLen)
		if decimalLen == priceDecimalMaxLen {
			copy(decimalPartArr, decimalStr)
		} else {
			//left pad 0
			for i := 0; i < priceDecimalMaxLen-decimalLen; i++ {
				decimalPartArr[i] = '0'
			}
			copy(decimalPartArr[priceDecimalMaxLen-decimalLen:], decimalStr)
		}
		var rightTruncate = 0
		for i := priceDecimalMaxLen - 1; i >= 0; i-- {
			if decimalPartArr[i] == '0' {
				rightTruncate++
			} else {
				break
			}
		}
		return fmt.Sprintf("%s.%s", intStr, string(decimalPartArr[:priceDecimalMaxLen-rightTruncate]))
	}
}

func getValueFromDb(db vm_db.VmDb, key []byte) []byte {
	if data, err := db.GetValue(key); err != nil {
		panic(err)
	} else {
		return data
	}
}

func setValueToDb(db vm_db.VmDb, key, value []byte) {
	if err := db.SetValue(key, value); err != nil {
		panic(err)
	}
}

func Uint64ToBytes(value uint64) []byte {
	bs := make([]byte, 8)
	binary.BigEndian.PutUint64(bs, value)
	return bs
}

func BytesToUint64(bytes []byte) uint64 {
	return binary.BigEndian.Uint64(bytes)
}

func Uint32ToBytes(value uint32) []byte {
	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, value)
	return bs
}

func BytesToUint32(bytes []byte) uint32 {
	return binary.BigEndian.Uint32(bytes)
}

func BitwiseNotBytes(bytes []byte) {
	for i, b := range bytes {
		bytes[i] = ^b
	}
}