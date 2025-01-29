package order

import (
	"encoding/binary"
	"errors"
	"time"

	"github.com/hossein1376/querier/pkg/model"
)

var (
	ErrInvalidSize = errors.New("unexpected data size")
)

// Encode encodes [model.Order] into a byte slice with the length of 18,
// ignoring the Status field.
func Encode(o model.Order) []byte {
	// int64 takes 8 bytes and int16 takes 2 bytes, and [time.Duration]
	// is just int64 under the hood, hence:
	// 8 + 8 + 2 = 18
	buf := make([]byte, 18)

	// int64 and uint64 have the same number of bits, with difference
	// being how sign bit is being treated. The same goes for int16 and
	// uint16. Therefore, it's safe to cast from int to uint.
	// [Decode] function will cast these bits into uint and our data
	// remains valid.
	binary.BigEndian.PutUint64(buf[0:], uint64(o.OrderID))
	binary.BigEndian.PutUint64(buf[8:], uint64(o.ProcessingTime))
	binary.BigEndian.PutUint16(buf[16:], uint16(o.Priority))

	return buf
}

// Decode extracts the [model.Order] instance out of the given data. The
// byte slice must be of size 18.
func Decode(data []byte) (*model.Order, error) {
	if len(data) != 18 {
		return nil, ErrInvalidSize
	}

	orderID := int64(binary.BigEndian.Uint64(data[0:]))
	processingTime := time.Duration(binary.BigEndian.Uint64(data[8:]))
	priority := model.Priority(binary.BigEndian.Uint16(data[16:]))

	return &model.Order{
		OrderID:        orderID,
		ProcessingTime: processingTime,
		Priority:       priority,
	}, nil
}
