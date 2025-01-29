package order_test

import (
	"encoding/json"
	"math/rand/v2"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/hossein1376/querier/pkg/model"
	"github.com/hossein1376/querier/pkg/order"
)

func TestEncode_Decode(t *testing.T) {
	a := assert.New(t)
	o := model.Order{
		OrderID:        rand.Int64(),
		Priority:       model.Normal,
		ProcessingTime: time.Second,
	}
	data := order.Encode(o)

	dec, err := order.Decode(data)
	a.NoError(err)
	a.Equal(o, *dec)
}

func BenchmarkEncode_Custom(b *testing.B) {
	o := model.Order{
		OrderID:        rand.Int64(),
		Priority:       model.Normal,
		ProcessingTime: time.Second,
	}

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		order.Encode(o)
	}
}

func BenchmarkEncode_JSON(b *testing.B) {
	o := &model.Order{
		OrderID:        rand.Int64(),
		Priority:       model.Normal,
		ProcessingTime: time.Second,
	}

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		_, _ = json.Marshal(o)
	}
}

func BenchmarkDecode_Custom(b *testing.B) {
	o := model.Order{
		OrderID:        rand.Int64(),
		Priority:       model.Normal,
		ProcessingTime: time.Second,
	}
	data := order.Encode(o)

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		_, _ = order.Decode(data)
	}
}

func BenchmarkDecode_JSON(b *testing.B) {
	o := &model.Order{
		OrderID:        rand.Int64(),
		Priority:       model.Normal,
		ProcessingTime: time.Second,
	}
	data, _ := json.Marshal(o)

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		var dec model.Order
		_ = json.Unmarshal(data, &dec)
	}
}
