package converters_test

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/stretchr/testify/assert"

	"github.com/ydb/grafana-ydb-datasource/pkg/converters"
)

func TestBool(t *testing.T) {
	value := true
	sut := converters.GetConverter("Bool")
	v, err := sut.FrameConverter.ConverterFunc(&value)
	assert.Nil(t, err)
	actual := v.(bool)
	assert.True(t, actual)
}

func TestDouble(t *testing.T) {
	value := 1.1
	sut := converters.GetConverter("Double")
	v, err := sut.FrameConverter.ConverterFunc(&value)
	assert.Nil(t, err)
	actual := v.(float64)
	assert.Equal(t, value, actual)
}

func TestFloat(t *testing.T) {
	value := float32(1.1)
	sut := converters.GetConverter("Float")
	v, err := sut.FrameConverter.ConverterFunc(&value)
	assert.Nil(t, err)
	actual := v.(float32)
	assert.Equal(t, value, actual)
}

func TestInt64(t *testing.T) {
	value := int64(1)
	sut := converters.GetConverter("Int64")
	v, err := sut.FrameConverter.ConverterFunc(&value)
	assert.Nil(t, err)
	actual := v.(int64)
	assert.Equal(t, value, actual)
}

func TestInt32(t *testing.T) {
	value := int32(1)
	sut := converters.GetConverter("Int32")
	v, err := sut.FrameConverter.ConverterFunc(&value)
	assert.Nil(t, err)
	actual := v.(int32)
	assert.Equal(t, value, actual)
}

func TestInt16(t *testing.T) {
	value := int16(1)
	sut := converters.GetConverter("Int16")
	v, err := sut.FrameConverter.ConverterFunc(&value)
	assert.Nil(t, err)
	actual := v.(int16)
	assert.Equal(t, value, actual)
}

func TestInt8(t *testing.T) {
	value := int8(1)
	sut := converters.GetConverter("Int8")
	v, err := sut.FrameConverter.ConverterFunc(&value)
	assert.Nil(t, err)
	actual := v.(int8)
	assert.Equal(t, value, actual)
}

func TestUint64(t *testing.T) {
	value := uint64(1)
	sut := converters.GetConverter("Uint64")
	v, err := sut.FrameConverter.ConverterFunc(&value)
	assert.Nil(t, err)
	actual := v.(uint64)
	assert.Equal(t, value, actual)
}

func TestUint32(t *testing.T) {
	value := uint32(1)
	sut := converters.GetConverter("Uint32")
	v, err := sut.FrameConverter.ConverterFunc(&value)
	assert.Nil(t, err)
	actual := v.(uint32)
	assert.Equal(t, value, actual)
}

func TestUint16(t *testing.T) {
	value := uint16(1)
	sut := converters.GetConverter("Uint16")
	v, err := sut.FrameConverter.ConverterFunc(&value)
	assert.Nil(t, err)
	actual := v.(uint16)
	assert.Equal(t, value, actual)
}

func TestUint8(t *testing.T) {
	value := uint8(1)
	sut := converters.GetConverter("Uint8")
	v, err := sut.FrameConverter.ConverterFunc(&value)
	assert.Nil(t, err)
	actual := v.(uint8)
	assert.Equal(t, value, actual)
}

func TestDecimal(t *testing.T) {
	value := decimal.New(22, 9)
	sut := converters.GetConverter("Decimal(22,9)")
	v, err := sut.FrameConverter.ConverterFunc(&value)
	assert.Nil(t, err)
	actual := v.(float64)
	f, _ := value.Float64()
	assert.Equal(t, f, actual)
}

func TestString(t *testing.T) {
	value := "test"
	sut := converters.GetConverter("String")
	v, err := sut.FrameConverter.ConverterFunc(&value)
	assert.Nil(t, err)
	actual := v.(string)
	assert.Equal(t, value, actual)
}

func TestUtf8(t *testing.T) {
	value := "test"
	sut := converters.GetConverter("Utf8")
	v, err := sut.FrameConverter.ConverterFunc(&value)
	assert.Nil(t, err)
	actual := v.(string)
	assert.Equal(t, value, actual)
}

func toJson(obj interface{}) (json.RawMessage, error) {
	bytes, err := json.Marshal(obj)
	if err != nil {
		return nil, errors.New("unable to marshal")
	}
	var rawJSON json.RawMessage
	err = json.Unmarshal(bytes, &rawJSON)
	if err != nil {
		return nil, errors.New("unable to unmarshal")
	}
	return rawJSON, nil
}

func TestJson(t *testing.T) {
	value := map[string]interface{}{
		"1": map[string]interface{}{
			"test": 1,
		},
		"2": uint16(2),
	}
	sut := converters.GetConverter("Json")
	v, err := sut.FrameConverter.ConverterFunc(&value)
	assert.Nil(t, err)
	msg, err := toJson(value)
	assert.Nil(t, err)
	assert.Equal(t, msg, *v.(*json.RawMessage))
}

func TestUuid(t *testing.T) {
	value, _ := uuid.Parse("5ddda1c0-bc6d-49ba-8062-38769ce088ce")
	sut := converters.GetConverter("Uuid")
	v, err := sut.FrameConverter.ConverterFunc(&value)
	assert.Nil(t, err)
	actual := v.(string)
	assert.Equal(t, value.String(), actual)
}

func TestDate(t *testing.T) {
	layout := "2006-01-02"
	str := "2006-01-02"
	d, _ := time.Parse(layout, str)
	sut := converters.GetConverter("Date")
	v, err := sut.FrameConverter.ConverterFunc(&d)
	assert.Nil(t, err)
	actual := v.(time.Time)
	assert.Equal(t, d, actual)
}

func TestDatetime(t *testing.T) {
	layout := "2006-01-02T15:04:05Z"
	str := "2006-01-02T15:04:05Z"
	d, _ := time.Parse(layout, str)
	sut := converters.GetConverter("Datetime")
	v, err := sut.FrameConverter.ConverterFunc(&d)
	assert.Nil(t, err)
	actual := v.(time.Time)
	assert.Equal(t, d, actual)
}

func TestTimestamp(t *testing.T) {
	layout := "2006-01-02T15:04:05.123456Z"
	str := "2006-01-02T15:04:05.123456Z"
	d, _ := time.Parse(layout, str)
	sut := converters.GetConverter("Timestamp")
	v, err := sut.FrameConverter.ConverterFunc(&d)
	assert.Nil(t, err)
	actual := v.(time.Time)
	assert.Equal(t, d, actual)
}

func TestInterval(t *testing.T) {
	layout := "2006-01-02T15:04:05.123456Z"
	str1 := "2006-01-02T15:04:05.123456Z"
	str2 := "2007-02-03T16:05:06.234567Z"
	d1, _ := time.Parse(layout, str1)
	d2, _ := time.Parse(layout, str2)
	i := d2.Sub(d1).Microseconds()
	sut := converters.GetConverter("Interval")
	v, err := sut.FrameConverter.ConverterFunc(&i)
	assert.Nil(t, err)
	actual := v.(int64)
	assert.Equal(t, i, actual)
}

func TestSmallSerial(t *testing.T) {
	value := int16(1)
	sut := converters.GetConverter("SmallSerial")
	v, err := sut.FrameConverter.ConverterFunc(&value)
	assert.Nil(t, err)
	actual := v.(int16)
	assert.Equal(t, value, actual)
}

func TestSerial2(t *testing.T) {
	value := int16(1)
	sut := converters.GetConverter("Serial2")
	v, err := sut.FrameConverter.ConverterFunc(&value)
	assert.Nil(t, err)
	actual := v.(int16)
	assert.Equal(t, value, actual)
}

func TestSerial(t *testing.T) {
	value := int32(1)
	sut := converters.GetConverter("Serial")
	v, err := sut.FrameConverter.ConverterFunc(&value)
	assert.Nil(t, err)
	actual := v.(int32)
	assert.Equal(t, value, actual)
}

func TestSerial4(t *testing.T) {
	value := int32(1)
	sut := converters.GetConverter("Serial4")
	v, err := sut.FrameConverter.ConverterFunc(&value)
	assert.Nil(t, err)
	actual := v.(int32)
	assert.Equal(t, value, actual)
}

func TestSerial8(t *testing.T) {
	value := int64(1)
	sut := converters.GetConverter("Serial8")
	v, err := sut.FrameConverter.ConverterFunc(&value)
	assert.Nil(t, err)
	actual := v.(int64)
	assert.Equal(t, value, actual)
}

func TestBigSerial(t *testing.T) {
	value := int64(1)
	sut := converters.GetConverter("BigSerial")
	v, err := sut.FrameConverter.ConverterFunc(&value)
	assert.Nil(t, err)
	actual := v.(int64)
	assert.Equal(t, value, actual)
}

func TestList(t *testing.T) {
	value := []string{"1", "2", "3"}
	sut := converters.GetConverter("List<String>")
	v, err := sut.FrameConverter.ConverterFunc(&value)
	assert.Nil(t, err)
	msg, err := toJson(value)
	assert.Nil(t, err)
	assert.Equal(t, msg, *v.(*json.RawMessage))
}

func TestDict(t *testing.T) {
	value := map[string]interface{}{
		"1": uint16(1),
		"2": uint16(2),
		"3": uint16(3),
		"4": uint16(4),
	}
	sut := converters.GetConverter("Dict<String, Uint16>")
	v, err := sut.FrameConverter.ConverterFunc(&value)
	assert.Nil(t, err)
	msg, err := toJson(value)
	assert.Nil(t, err)
	assert.Equal(t, msg, *v.(*json.RawMessage))
}

func TestSet(t *testing.T) {
	value := []string{"1", "2", "3"}
	ipConverter := converters.GetConverter("Set<String>")
	v, err := ipConverter.FrameConverter.ConverterFunc(&value)
	assert.Nil(t, err)
	msg, err := toJson(value)
	assert.Nil(t, err)
	assert.Equal(t, msg, *v.(*json.RawMessage))
}

func TestTuple(t *testing.T) {
	value := map[string]interface{}{
		"Id":   int64(1),
		"Name": "Test",
	}
	sut := converters.GetConverter("Tuple<Int64, String>")
	v, err := sut.FrameConverter.ConverterFunc(&value)
	assert.Nil(t, err)
	msg, err := toJson(value)
	assert.Nil(t, err)
	assert.Equal(t, msg, *v.(*json.RawMessage))
}

func TestStruct(t *testing.T) {
	value := map[string]interface{}{
		"Id":   int64(1),
		"Name": "Test",
	}
	sut := converters.GetConverter("Struct<Id:Int64, Name:String>")
	v, err := sut.FrameConverter.ConverterFunc(&value)
	assert.Nil(t, err)
	msg, err := toJson(value)
	assert.Nil(t, err)
	assert.Equal(t, msg, *v.(*json.RawMessage))
}

func TestVariant(t *testing.T) {
	value := map[string]interface{}{
		"Id":   int64(1),
		"Name": "Test",
	}
	sut := converters.GetConverter("Variant<Id:Int64, Name:String>")
	v, err := sut.FrameConverter.ConverterFunc(&value)
	assert.Nil(t, err)
	msg, err := toJson(value)
	assert.Nil(t, err)
	assert.Equal(t, msg, *v.(*json.RawMessage))
}

func TestEnum(t *testing.T) {
	value := []string{"Test1", "Test2", "Test3"}
	sut := converters.GetConverter("Enum<Test1, Test2, Test3>")
	v, err := sut.FrameConverter.ConverterFunc(&value)
	assert.Nil(t, err)
	msg, err := toJson(value)
	assert.Nil(t, err)
	assert.Equal(t, msg, *v.(*json.RawMessage))
}

func TestOptionalBool(t *testing.T) {
	value := true
	val := &value
	sut := converters.GetConverter("Optional<Bool>")
	v, err := sut.FrameConverter.ConverterFunc(&val)
	assert.Nil(t, err)
	actual := v.(*bool)
	assert.True(t, value, *actual)
}

func TestOptionalBoolShouldBeNil(t *testing.T) {
	var value *bool
	sut := converters.GetConverter("Optional<Bool>")
	v, err := sut.FrameConverter.ConverterFunc(&value)
	assert.Nil(t, err)
	actual := v.(*bool)
	assert.Equal(t, value, actual)
}

func TestOptionalDouble(t *testing.T) {
	value := 1.1
	val := &value
	sut := converters.GetConverter("Optional<Double>")
	v, err := sut.FrameConverter.ConverterFunc(&val)
	assert.Nil(t, err)
	actual := v.(*float64)
	assert.Equal(t, value, *actual)
}

func TestOptionalDoubleShouldBeNil(t *testing.T) {
	var value *float64
	sut := converters.GetConverter("Optional<Double>")
	v, err := sut.FrameConverter.ConverterFunc(&value)
	assert.Nil(t, err)
	actual := v.(*float64)
	assert.Equal(t, value, actual)
}

func TestOptionalFloat(t *testing.T) {
	value := float32(1.1)
	val := &value
	sut := converters.GetConverter("Optional<Float>")
	v, err := sut.FrameConverter.ConverterFunc(&val)
	assert.Nil(t, err)
	actual := v.(*float32)
	assert.Equal(t, value, *actual)
}

func TestOptionalFloatShouldBeNil(t *testing.T) {
	var value *float32
	sut := converters.GetConverter("Optional<Float>")
	v, err := sut.FrameConverter.ConverterFunc(&value)
	assert.Nil(t, err)
	actual := v.(*float32)
	assert.Equal(t, value, actual)
}

func TestOptionalInt64(t *testing.T) {
	value := int64(1)
	val := &value
	sut := converters.GetConverter("Optional<Int64>")
	v, err := sut.FrameConverter.ConverterFunc(&val)
	assert.Nil(t, err)
	actual := v.(*int64)
	assert.Equal(t, value, *actual)
}

func TestOptionalInt64ShouldBeNil(t *testing.T) {
	var value *int64
	sut := converters.GetConverter("Optional<Int64>")
	v, err := sut.FrameConverter.ConverterFunc(&value)
	assert.Nil(t, err)
	actual := v.(*int64)
	assert.Equal(t, value, actual)
}

func TestOptionalInt32(t *testing.T) {
	value := int32(1)
	val := &value
	sut := converters.GetConverter("Optional<Int32>")
	v, err := sut.FrameConverter.ConverterFunc(&val)
	assert.Nil(t, err)
	actual := v.(*int32)
	assert.Equal(t, value, *actual)
}

func TestOptionalInt32ShouldBeNil(t *testing.T) {
	var value *int32
	sut := converters.GetConverter("Optional<Int32>")
	v, err := sut.FrameConverter.ConverterFunc(&value)
	assert.Nil(t, err)
	actual := v.(*int32)
	assert.Equal(t, value, actual)
}

func TestOptionalInt16(t *testing.T) {
	value := int16(1)
	val := &value
	sut := converters.GetConverter("Optional<Int16>")
	v, err := sut.FrameConverter.ConverterFunc(&val)
	assert.Nil(t, err)
	actual := v.(*int16)
	assert.Equal(t, value, *actual)
}

func TestOptionalInt16ShouldBeNil(t *testing.T) {
	var value *int16
	sut := converters.GetConverter("Optional<Int16>")
	v, err := sut.FrameConverter.ConverterFunc(&value)
	assert.Nil(t, err)
	actual := v.(*int16)
	assert.Equal(t, value, actual)
}

func TestOptionalInt8(t *testing.T) {
	value := int8(1)
	val := &value
	sut := converters.GetConverter("Optional<Int8>")
	v, err := sut.FrameConverter.ConverterFunc(&val)
	assert.Nil(t, err)
	actual := v.(*int8)
	assert.Equal(t, value, *actual)
}

func TestOptionalInt8ShouldBeNil(t *testing.T) {
	var value *int8
	sut := converters.GetConverter("Optional<Int8>")
	v, err := sut.FrameConverter.ConverterFunc(&value)
	assert.Nil(t, err)
	actual := v.(*int8)
	assert.Equal(t, value, actual)
}

func TestOptionalUint64(t *testing.T) {
	value := uint64(1)
	val := &value
	sut := converters.GetConverter("Optional<Uint64>")
	v, err := sut.FrameConverter.ConverterFunc(&val)
	assert.Nil(t, err)
	actual := v.(*uint64)
	assert.Equal(t, value, *actual)
}

func TestOptionalUint64ShouldBeNil(t *testing.T) {
	var value *uint64
	sut := converters.GetConverter("Optional<Uint64>")
	v, err := sut.FrameConverter.ConverterFunc(&value)
	assert.Nil(t, err)
	actual := v.(*uint64)
	assert.Equal(t, value, actual)
}

func TestOptionalUint32(t *testing.T) {
	value := uint32(1)
	val := &value
	sut := converters.GetConverter("Optional<Uint32>")
	v, err := sut.FrameConverter.ConverterFunc(&val)
	assert.Nil(t, err)
	actual := v.(*uint32)
	assert.Equal(t, value, *actual)
}

func TestOptionalUint32ShouldBeNil(t *testing.T) {
	var value *uint32
	sut := converters.GetConverter("Optional<Uint32>")
	v, err := sut.FrameConverter.ConverterFunc(&value)
	assert.Nil(t, err)
	actual := v.(*uint32)
	assert.Equal(t, value, actual)
}

func TestOptionalUint16(t *testing.T) {
	value := uint16(1)
	val := &value
	sut := converters.GetConverter("Optional<Uint16>")
	v, err := sut.FrameConverter.ConverterFunc(&val)
	assert.Nil(t, err)
	actual := v.(*uint16)
	assert.Equal(t, value, *actual)
}

func TestOptionalUint16ShouldBeNil(t *testing.T) {
	var value *uint16
	sut := converters.GetConverter("Optional<Uint16>")
	v, err := sut.FrameConverter.ConverterFunc(&value)
	assert.Nil(t, err)
	actual := v.(*uint16)
	assert.Equal(t, value, actual)
}

func TestOptionalUint8(t *testing.T) {
	value := uint8(1)
	val := &value
	sut := converters.GetConverter("Optional<Uint8>")
	v, err := sut.FrameConverter.ConverterFunc(&val)
	assert.Nil(t, err)
	actual := v.(*uint8)
	assert.Equal(t, value, *actual)
}

func TestOptionalUint8ShouldBeNil(t *testing.T) {
	var value *uint8
	sut := converters.GetConverter("Optional<Uint8>")
	v, err := sut.FrameConverter.ConverterFunc(&value)
	assert.Nil(t, err)
	actual := v.(*uint8)
	assert.Equal(t, value, actual)
}

func TestOptionalDecimal(t *testing.T) {
	value := decimal.New(22, 9)
	val := &value
	sut := converters.GetConverter("Optional<Decimal(22,9)>")
	v, err := sut.FrameConverter.ConverterFunc(&val)
	assert.Nil(t, err)
	actual := v.(*float64)
	f, _ := value.Float64()
	assert.Equal(t, f, *actual)
}

func TestOptionalDecimalShouldBeNil(t *testing.T) {
	var value *decimal.Decimal
	sut := converters.GetConverter("Optional<Decimal(22,9)>")
	v, err := sut.FrameConverter.ConverterFunc(&value)
	assert.Nil(t, err)
	actual := v.(*float64)
	assert.Equal(t, (*float64)(nil), actual)
}

func TestOptionalString(t *testing.T) {
	value := "test"
	val := &value
	sut := converters.GetConverter("Optional<String>")
	v, err := sut.FrameConverter.ConverterFunc(&val)
	assert.Nil(t, err)
	actual := v.(*string)
	assert.Equal(t, value, *actual)
}

func TestOptionalStringShouldBeNil(t *testing.T) {
	var value *string
	sut := converters.GetConverter("Optional<String>")
	v, err := sut.FrameConverter.ConverterFunc(&value)
	assert.Nil(t, err)
	actual := v.(*string)
	assert.Equal(t, value, actual)
}

func TestOptionalUtf8(t *testing.T) {
	value := "test"
	val := &value
	sut := converters.GetConverter("Optional<Utf8>")
	v, err := sut.FrameConverter.ConverterFunc(&val)
	assert.Nil(t, err)
	actual := v.(*string)
	assert.Equal(t, value, *actual)
}

func TestOptionalUtf8ShouldBeNil(t *testing.T) {
	var value *string
	sut := converters.GetConverter("Optional<Utf8>")
	v, err := sut.FrameConverter.ConverterFunc(&value)
	assert.Nil(t, err)
	actual := v.(*string)
	assert.Equal(t, value, actual)
}

func TestOptionalJson(t *testing.T) {
	value := map[string]interface{}{
		"1": map[string]interface{}{
			"test": 1,
		},
		"2": uint16(2),
	}
	var data interface{} = value
	val := &data
	sut := converters.GetConverter("Optional<Json>")
	v, err := sut.FrameConverter.ConverterFunc(&val)
	assert.Nil(t, err)
	msg, err := toJson(value)
	assert.Nil(t, err)
	assert.Equal(t, msg, **v.(**json.RawMessage))
}

func TestOptionalJsonShouldBeNil(t *testing.T) {
	var ptr *interface{} = nil
	sut := converters.GetConverter("Optional<Json>")
	v, err := sut.FrameConverter.ConverterFunc(&ptr)
	assert.Nil(t, err)
	assert.Nil(t, v)
}

func TestOptionalUuid(t *testing.T) {
	value, _ := uuid.Parse("5ddda1c0-bc6d-49ba-8062-38769ce088ce")
	val := &value
	sut := converters.GetConverter("Optional<Uuid>")
	v, err := sut.FrameConverter.ConverterFunc(&val)
	assert.Nil(t, err)
	actual := v.(*string)
	assert.Equal(t, value.String(), *actual)
}

func TestOptionalUuidShouldBeNil(t *testing.T) {
	var value *uuid.UUID
	sut := converters.GetConverter("Optional<Uuid>")
	v, err := sut.FrameConverter.ConverterFunc(&value)
	assert.Nil(t, err)
	actual := v.(*string)
	assert.Equal(t, (*string)(nil), actual)
}

func TestOptionalDate(t *testing.T) {
	layout := "2006-01-02"
	str := "2006-01-02"
	d, _ := time.Parse(layout, str)
	val := &d
	sut := converters.GetConverter("Optional<Date>")
	v, err := sut.FrameConverter.ConverterFunc(&val)
	assert.Nil(t, err)
	actual := v.(*time.Time)
	assert.Equal(t, d, *actual)
}

func TestOptionalDateShouldBeNil(t *testing.T) {
	var d *time.Time
	sut := converters.GetConverter("Optional<Date>")
	v, err := sut.FrameConverter.ConverterFunc(&d)
	assert.Nil(t, err)
	actual := v.(*time.Time)
	assert.Equal(t, d, actual)
}

func TestOptionalDatetime(t *testing.T) {
	layout := "2006-01-02T15:04:05Z"
	str := "2006-01-02T15:04:05Z"
	d, _ := time.Parse(layout, str)
	val := &d
	sut := converters.GetConverter("Optional<Datetime>")
	v, err := sut.FrameConverter.ConverterFunc(&val)
	assert.Nil(t, err)
	actual := v.(*time.Time)
	assert.Equal(t, d, *actual)
}

func TestOptionalDatetimeShouldBeNil(t *testing.T) {
	var d *time.Time
	sut := converters.GetConverter("Optional<Datetime>")
	v, err := sut.FrameConverter.ConverterFunc(&d)
	assert.Nil(t, err)
	actual := v.(*time.Time)
	assert.Equal(t, d, actual)
}

func TestOptionalTimestamp(t *testing.T) {
	layout := "2006-01-02T15:04:05.123456Z"
	str := "2006-01-02T15:04:05.123456Z"
	d, _ := time.Parse(layout, str)
	val := &d
	sut := converters.GetConverter("Optional<Timestamp>")
	v, err := sut.FrameConverter.ConverterFunc(&val)
	assert.Nil(t, err)
	actual := v.(*time.Time)
	assert.Equal(t, d, *actual)
}

func TestOptionalTimestampShouldBeNil(t *testing.T) {
	var d *time.Time
	sut := converters.GetConverter("Optional<Timestamp>")
	v, err := sut.FrameConverter.ConverterFunc(&d)
	assert.Nil(t, err)
	actual := v.(*time.Time)
	assert.Equal(t, d, actual)
}

func TestOptionalInterval(t *testing.T) {
	layout := "2006-01-02T15:04:05.123456Z"
	str1 := "2006-01-02T15:04:05.123456Z"
	str2 := "2007-02-03T16:05:06.234567Z"
	d1, _ := time.Parse(layout, str1)
	d2, _ := time.Parse(layout, str2)
	i := d2.Sub(d1).Microseconds()
	val := &i
	sut := converters.GetConverter("Optional<Interval>")
	v, err := sut.FrameConverter.ConverterFunc(&val)
	assert.Nil(t, err)
	actual := v.(*int64)
	assert.Equal(t, i, *actual)
}

func TestOptionalIntervalShouldBeNil(t *testing.T) {
	var value *int64
	sut := converters.GetConverter("Optional<Interval>")
	v, err := sut.FrameConverter.ConverterFunc(&value)
	assert.Nil(t, err)
	actual := v.(*int64)
	assert.Equal(t, value, actual)
}

func TestOptionalSmallSerial(t *testing.T) {
	value := int16(1)
	val := &value
	sut := converters.GetConverter("Optional<SmallSerial>")
	v, err := sut.FrameConverter.ConverterFunc(&val)
	assert.Nil(t, err)
	actual := v.(*int16)
	assert.Equal(t, value, *actual)
}

func TestOptionalSmallSerialShouldBeNil(t *testing.T) {
	var value *int16
	sut := converters.GetConverter("Optional<SmallSerial>")
	v, err := sut.FrameConverter.ConverterFunc(&value)
	assert.Nil(t, err)
	actual := v.(*int16)
	assert.Equal(t, value, actual)
}

func TestOptionalSerial2(t *testing.T) {
	value := int16(1)
	val := &value
	sut := converters.GetConverter("Optional<Serial2>")
	v, err := sut.FrameConverter.ConverterFunc(&val)
	assert.Nil(t, err)
	actual := v.(*int16)
	assert.Equal(t, value, *actual)
}

func TestOptionalSerial2ShouldBeNil(t *testing.T) {
	var value *int16
	sut := converters.GetConverter("Optional<Serial2>")
	v, err := sut.FrameConverter.ConverterFunc(&value)
	assert.Nil(t, err)
	actual := v.(*int16)
	assert.Equal(t, value, actual)
}

func TestOptionalSerial(t *testing.T) {
	value := int32(1)
	val := &value
	sut := converters.GetConverter("Optional<Serial>")
	v, err := sut.FrameConverter.ConverterFunc(&val)
	assert.Nil(t, err)
	actual := v.(*int32)
	assert.Equal(t, value, *actual)
}

func TestOptionalSerialShouldBeNil(t *testing.T) {
	var value *int32
	sut := converters.GetConverter("Optional<Serial>")
	v, err := sut.FrameConverter.ConverterFunc(&value)
	assert.Nil(t, err)
	actual := v.(*int32)
	assert.Equal(t, value, actual)
}

func TestOptionalSerial4(t *testing.T) {
	value := int32(1)
	val := &value
	sut := converters.GetConverter("Optional<Serial4>")
	v, err := sut.FrameConverter.ConverterFunc(&val)
	assert.Nil(t, err)
	actual := v.(*int32)
	assert.Equal(t, value, *actual)
}

func TestOptionalSerial4ShouldBeNil(t *testing.T) {
	var value *int32
	sut := converters.GetConverter("Optional<Serial4>")
	v, err := sut.FrameConverter.ConverterFunc(&value)
	assert.Nil(t, err)
	actual := v.(*int32)
	assert.Equal(t, value, actual)
}

func TestOptionalSerial8(t *testing.T) {
	value := int64(1)
	val := &value
	sut := converters.GetConverter("Optional<Serial8>")
	v, err := sut.FrameConverter.ConverterFunc(&val)
	assert.Nil(t, err)
	actual := v.(*int64)
	assert.Equal(t, value, *actual)
}

func TestOptionalSerial8ShouldBeNil(t *testing.T) {
	var value *int64
	sut := converters.GetConverter("Optional<Serial8>")
	v, err := sut.FrameConverter.ConverterFunc(&value)
	assert.Nil(t, err)
	actual := v.(*int64)
	assert.Equal(t, value, actual)
}

func TestOptionalBigSerial(t *testing.T) {
	value := int64(1)
	val := &value
	sut := converters.GetConverter("Optional<BigSerial>")
	v, err := sut.FrameConverter.ConverterFunc(&val)
	assert.Nil(t, err)
	actual := v.(*int64)
	assert.Equal(t, value, *actual)
}

func TestOptionalBigSerialShouldBeNil(t *testing.T) {
	var value *int64
	sut := converters.GetConverter("Optional<BigSerial>")
	v, err := sut.FrameConverter.ConverterFunc(&value)
	assert.Nil(t, err)
	actual := v.(*int64)
	assert.Equal(t, value, actual)
}

func TestOptionalList(t *testing.T) {
	value := []string{"1", "2", "3"}
	var data interface{} = value
	val := &data
	sut := converters.GetConverter("Optional<List<String>>")
	v, err := sut.FrameConverter.ConverterFunc(&val)
	assert.Nil(t, err)
	msg, err := toJson(value)
	assert.Nil(t, err)
	assert.Equal(t, msg, **v.(**json.RawMessage))
}

func TestOptionalListShouldBeNil(t *testing.T) {
	var ptr *interface{} = nil
	sut := converters.GetConverter("Optional<List<String>>")
	v, err := sut.FrameConverter.ConverterFunc(&ptr)
	assert.Nil(t, err)
	assert.Nil(t, v)
}

func TestOptionalDict(t *testing.T) {
	value := map[string]interface{}{
		"1": uint16(1),
		"2": uint16(2),
		"3": uint16(3),
		"4": uint16(4),
	}
	var data interface{} = value
	val := &data
	sut := converters.GetConverter("Optional<Dict<String, Uint16>>")
	v, err := sut.FrameConverter.ConverterFunc(&val)
	assert.Nil(t, err)
	msg, err := toJson(value)
	assert.Nil(t, err)
	assert.Equal(t, msg, **v.(**json.RawMessage))
}

func TestOptionalDictShouldBeNil(t *testing.T) {
	var ptr *interface{} = nil
	sut := converters.GetConverter("Optional<Dict<String, Uint16>>")
	v, err := sut.FrameConverter.ConverterFunc(&ptr)
	assert.Nil(t, err)
	assert.Nil(t, v)
}

func TestOptionalSet(t *testing.T) {
	value := []string{"1", "2", "3"}
	var data interface{} = value
	val := &data
	sut := converters.GetConverter("Optional<Set<String>>")
	v, err := sut.FrameConverter.ConverterFunc(&val)
	assert.Nil(t, err)
	msg, err := toJson(value)
	assert.Nil(t, err)
	assert.Equal(t, msg, **v.(**json.RawMessage))
}

func TestOptionalSetShouldBeNil(t *testing.T) {
	var ptr *interface{} = nil
	sut := converters.GetConverter("Optional<Dict<String, Uint16>>")
	v, err := sut.FrameConverter.ConverterFunc(&ptr)
	assert.Nil(t, err)
	assert.Nil(t, v)
}

func TestOptionalTuple(t *testing.T) {
	value := map[string]interface{}{
		"Id":   int64(1),
		"Name": "Test",
	}
	var data interface{} = value
	val := &data
	sut := converters.GetConverter("Optional<Tuple<Int64, String>>")
	v, err := sut.FrameConverter.ConverterFunc(&val)
	assert.Nil(t, err)
	msg, err := toJson(value)
	assert.Nil(t, err)
	assert.Equal(t, msg, **v.(**json.RawMessage))
}

func TestOptionalTupleShouldBeNil(t *testing.T) {
	var ptr *interface{} = nil
	sut := converters.GetConverter("Optional<Dict<String, Uint16>>")
	v, err := sut.FrameConverter.ConverterFunc(&ptr)
	assert.Nil(t, err)
	assert.Nil(t, v)
}

func TestOptionalStruct(t *testing.T) {
	value := map[string]interface{}{
		"Id":   int64(1),
		"Name": "Test",
	}
	var data interface{} = value
	val := &data
	sut := converters.GetConverter("Optional<Struct<Id:Int64, Name:String>>")
	v, err := sut.FrameConverter.ConverterFunc(&val)
	assert.Nil(t, err)
	msg, err := toJson(value)
	assert.Nil(t, err)
	assert.Equal(t, msg, **v.(**json.RawMessage))
}

func TestOptionalStructShouldBeNil(t *testing.T) {
	var ptr *interface{} = nil
	sut := converters.GetConverter("Optional<Struct<Id:Int64, Name:String>>")
	v, err := sut.FrameConverter.ConverterFunc(&ptr)
	assert.Nil(t, err)
	assert.Nil(t, v)
}

func TestOptionalVariant(t *testing.T) {
	value := map[string]interface{}{
		"Id":   int64(1),
		"Name": "Test",
	}
	var data interface{} = value
	val := &data
	sut := converters.GetConverter("Optional<Variant<Id:Int64, Name:String>>")
	v, err := sut.FrameConverter.ConverterFunc(&val)
	assert.Nil(t, err)
	msg, err := toJson(value)
	assert.Nil(t, err)
	assert.Equal(t, msg, **v.(**json.RawMessage))
}

func TestOptionalVariantShouldBeNil(t *testing.T) {
	var ptr *interface{} = nil
	sut := converters.GetConverter("Optional<Variant<Id:Int64, Name:String>>")
	v, err := sut.FrameConverter.ConverterFunc(&ptr)
	assert.Nil(t, err)
	assert.Nil(t, v)
}

func TestOptionalEnum(t *testing.T) {
	value := []string{"Test1", "Test2", "Test3"}
	var data interface{} = value
	val := &data
	sut := converters.GetConverter("Optional<Enum<Test1, Test2, Test3>>")
	v, err := sut.FrameConverter.ConverterFunc(&val)
	assert.Nil(t, err)
	msg, err := toJson(value)
	assert.Nil(t, err)
	assert.Equal(t, msg, **v.(**json.RawMessage))
}

func TestOptionalEnumShouldBeNil(t *testing.T) {
	var ptr *interface{} = nil
	sut := converters.GetConverter("Optional<Enum<Test1, Test2, Test3>>")
	v, err := sut.FrameConverter.ConverterFunc(&ptr)
	assert.Nil(t, err)
	assert.Nil(t, v)
}
