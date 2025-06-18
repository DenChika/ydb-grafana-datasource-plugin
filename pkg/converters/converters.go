package converters

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"time"

	"github.com/google/uuid"

	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/grafana-plugin-sdk-go/data/sqlutil"
	"github.com/shopspring/decimal"
)

type Converter struct {
	scanType   reflect.Type
	fieldType  data.FieldType
	matchRegex *regexp.Regexp
	convert    func(in interface{}) (interface{}, error)
}

var matchRegexes = map[string]*regexp.Regexp{
	"Decimal":           regexp.MustCompile(`^Decimal\(`),
	"Json":              regexp.MustCompile(`^Json\(`),
	"Uuid":              regexp.MustCompile(`^Uuid\(`),
	"List":              regexp.MustCompile(`^List<.*>`),
	"Dict":              regexp.MustCompile(`^Dict<.*>`),
	"Set":               regexp.MustCompile(`^Set<.*>`),
	"Tuple":             regexp.MustCompile(`^Tuple<.*>`),
	"Struct":            regexp.MustCompile(`^Struct<.*>`),
	"Variant":           regexp.MustCompile(`^Variant<.*>`),
	"Enum":              regexp.MustCompile(`^Enum<.*>`),
	"Optional<Decimal>": regexp.MustCompile(`^Optional<Decimal\(`),
	"Optional<Json>":    regexp.MustCompile(`^Optional<Json\(`),
	"Optional<Uuid>":    regexp.MustCompile(`^Optional<Uuid\(`),
	"Optional<List>":    regexp.MustCompile(`^Optional<List\(`),
	"Optional<Dict>":    regexp.MustCompile(`^Optional<Dict\(`),
	"Optional<Set>":     regexp.MustCompile(`^Optional<Set\(`),
	"Optional<Tuple>":   regexp.MustCompile(`^Optional<Tuple\(`),
	"Optional<Struct>":  regexp.MustCompile(`^Optional<Struct\(`),
	"Optional<Variant>": regexp.MustCompile(`^Optional<Variant\(`),
	"Optional<Enum>":    regexp.MustCompile(`^Optional<Enum\(`),
}

var Converters = map[string]Converter{
	"Bool": {
		scanType:  reflect.PointerTo(reflect.TypeOf(true)),
		fieldType: data.FieldTypeBool,
	},
	"Double": {
		scanType:  reflect.PointerTo(reflect.TypeOf(float64(0))),
		fieldType: data.FieldTypeFloat64,
	},
	"Float": {
		scanType:  reflect.PointerTo(reflect.TypeOf(float32(0))),
		fieldType: data.FieldTypeFloat32,
	},
	"Int64": {
		scanType:  reflect.PointerTo(reflect.TypeOf(int64(0))),
		fieldType: data.FieldTypeInt64,
	},
	"Int32": {
		scanType:  reflect.PointerTo(reflect.TypeOf(int32(0))),
		fieldType: data.FieldTypeInt32,
	},
	"Int16": {
		scanType:  reflect.PointerTo(reflect.TypeOf(int16(0))),
		fieldType: data.FieldTypeInt16,
	},
	"Int8": {
		scanType:  reflect.PointerTo(reflect.TypeOf(int8(0))),
		fieldType: data.FieldTypeInt8,
	},
	"Uint64": {
		scanType:  reflect.PointerTo(reflect.TypeOf(uint64(0))),
		fieldType: data.FieldTypeUint64,
	},
	"Uint32": {
		scanType:  reflect.PointerTo(reflect.TypeOf(uint32(0))),
		fieldType: data.FieldTypeUint32,
	},
	"Uint16": {
		scanType:  reflect.PointerTo(reflect.TypeOf(uint16(0))),
		fieldType: data.FieldTypeUint16,
	},
	"Uint8": {
		scanType:  reflect.PointerTo(reflect.TypeOf(uint8(0))),
		fieldType: data.FieldTypeUint8,
	},
	"Decimal": {
		convert:    decimalConvert,
		fieldType:  data.FieldTypeFloat64,
		matchRegex: matchRegexes["Decimal"],
		scanType:   reflect.PointerTo(reflect.TypeOf(decimal.Decimal{})),
	},
	"String": {
		fieldType: data.FieldTypeString,
		scanType:  reflect.PointerTo(reflect.TypeOf("")),
	},
	"Utf8": {
		fieldType: data.FieldTypeString,
		scanType:  reflect.PointerTo(reflect.TypeOf("")),
	},
	"Json": {
		convert:    jsonConvert,
		fieldType:  data.FieldTypeJSON,
		matchRegex: matchRegexes["Json"],
		scanType:   reflect.TypeOf((*interface{})(nil)).Elem(),
	},
	"Uuid": {
		convert:    uuidConvert,
		fieldType:  data.FieldTypeString,
		matchRegex: matchRegexes["Uuid"],
		scanType:   reflect.PointerTo(reflect.TypeOf("")),
	},
	"Date": {
		fieldType: data.FieldTypeTime,
		scanType:  reflect.PointerTo(reflect.TypeOf(time.Time{})),
	},
	"Datetime": {
		fieldType: data.FieldTypeTime,
		scanType:  reflect.PointerTo(reflect.TypeOf(time.Time{})),
	},
	"Timestamp": {
		fieldType: data.FieldTypeTime,
		scanType:  reflect.PointerTo(reflect.TypeOf(time.Time{})),
	},
	"Interval": {
		convert:   intervalConvert,
		scanType:  reflect.PointerTo(reflect.TypeOf(int64(0))),
		fieldType: data.FieldTypeInt64,
	},
	"SmallSerial": {
		scanType:  reflect.PointerTo(reflect.TypeOf(int16(0))),
		fieldType: data.FieldTypeInt16,
	},
	"Serial2": {
		scanType:  reflect.PointerTo(reflect.TypeOf(int16(0))),
		fieldType: data.FieldTypeInt16,
	},
	"Serial": {
		scanType:  reflect.PointerTo(reflect.TypeOf(int32(0))),
		fieldType: data.FieldTypeInt32,
	},
	"Serial4": {
		scanType:  reflect.PointerTo(reflect.TypeOf(int32(0))),
		fieldType: data.FieldTypeInt32,
	},
	"Serial8": {
		scanType:  reflect.PointerTo(reflect.TypeOf(int64(0))),
		fieldType: data.FieldTypeInt64,
	},
	"BigSerial": {
		scanType:  reflect.PointerTo(reflect.TypeOf(int64(0))),
		fieldType: data.FieldTypeInt64,
	},
	"List<>": {
		convert:    jsonConvert,
		fieldType:  data.FieldTypeJSON,
		matchRegex: matchRegexes["List"],
		scanType:   reflect.TypeOf((*interface{})(nil)).Elem(),
	},
	"Dict<>": {
		convert:    jsonConvert,
		fieldType:  data.FieldTypeJSON,
		matchRegex: matchRegexes["Dict"],
		scanType:   reflect.TypeOf((*interface{})(nil)).Elem(),
	},
	"Set<>": {
		convert:    jsonConvert,
		fieldType:  data.FieldTypeJSON,
		matchRegex: matchRegexes["Set"],
		scanType:   reflect.TypeOf((*interface{})(nil)).Elem(),
	},
	"Tuple<>": {
		convert:    jsonConvert,
		fieldType:  data.FieldTypeJSON,
		matchRegex: matchRegexes["Tuple"],
		scanType:   reflect.TypeOf((*interface{})(nil)).Elem(),
	},
	"Struct<>": {
		convert:    jsonConvert,
		fieldType:  data.FieldTypeJSON,
		matchRegex: matchRegexes["Struct"],
		scanType:   reflect.TypeOf((*interface{})(nil)).Elem(),
	},
	"Variant<>": {
		convert:    jsonConvert,
		fieldType:  data.FieldTypeJSON,
		matchRegex: matchRegexes["Variant"],
		scanType:   reflect.TypeOf((*interface{})(nil)).Elem(),
	},
	"Enum<>": {
		convert:    jsonConvert,
		fieldType:  data.FieldTypeJSON,
		matchRegex: matchRegexes["Enum"],
		scanType:   reflect.TypeOf((*interface{})(nil)).Elem(),
	},
	"Optional<Bool>": {
		scanType:  reflect.PointerTo(reflect.PointerTo(reflect.TypeOf(true))),
		fieldType: data.FieldTypeNullableBool,
	},
	"Optional<Double>": {
		scanType:  reflect.PointerTo(reflect.PointerTo(reflect.TypeOf(float64(0)))),
		fieldType: data.FieldTypeNullableFloat64,
	},
	"Optional<Float>": {
		scanType:  reflect.PointerTo(reflect.PointerTo(reflect.TypeOf(float32(0)))),
		fieldType: data.FieldTypeNullableFloat32,
	},
	"Optional<Int64>": {
		scanType:  reflect.PointerTo(reflect.PointerTo(reflect.TypeOf(int64(0)))),
		fieldType: data.FieldTypeNullableInt64,
	},
	"Optional<Int32>": {
		scanType:  reflect.PointerTo(reflect.PointerTo(reflect.TypeOf(int32(0)))),
		fieldType: data.FieldTypeNullableInt32,
	},
	"Optional<Int16>": {
		scanType:  reflect.PointerTo(reflect.PointerTo(reflect.TypeOf(int16(0)))),
		fieldType: data.FieldTypeNullableInt16,
	},
	"Optional<Int8>": {
		scanType:  reflect.PointerTo(reflect.PointerTo(reflect.TypeOf(int8(0)))),
		fieldType: data.FieldTypeNullableInt8,
	},
	"Optional<Uint64>": {
		scanType:  reflect.PointerTo(reflect.PointerTo(reflect.TypeOf(uint64(0)))),
		fieldType: data.FieldTypeNullableUint64,
	},
	"Optional<Uint32>": {
		scanType:  reflect.PointerTo(reflect.PointerTo(reflect.TypeOf(uint32(0)))),
		fieldType: data.FieldTypeNullableUint32,
	},
	"Optional<Uint16>": {
		scanType:  reflect.PointerTo(reflect.PointerTo(reflect.TypeOf(uint16(0)))),
		fieldType: data.FieldTypeNullableUint16,
	},
	"Optional<Uint8>": {
		scanType:  reflect.PointerTo(reflect.PointerTo(reflect.TypeOf(uint8(0)))),
		fieldType: data.FieldTypeNullableUint8,
	},
	"Optional<Decimal>": {
		convert:    decimalNullConvert,
		fieldType:  data.FieldTypeNullableFloat64,
		matchRegex: matchRegexes["Optional<Decimal>"],
		scanType:   reflect.PointerTo(reflect.PointerTo(reflect.TypeOf(decimal.Decimal{}))),
	},
	"Optional<String>": {
		fieldType: data.FieldTypeNullableString,
		scanType:  reflect.PointerTo(reflect.PointerTo(reflect.TypeOf(""))),
	},
	"Optional<Utf8>": {
		fieldType: data.FieldTypeNullableString,
		scanType:  reflect.PointerTo(reflect.PointerTo(reflect.TypeOf(""))),
	},
	"Optional<Json>": {
		convert:    jsonConvert,
		fieldType:  data.FieldTypeJSON,
		matchRegex: matchRegexes["Optional<Json>"],
		scanType:   reflect.PointerTo(reflect.TypeOf((*interface{})(nil)).Elem()),
	},
	"Optional<Uuid>": {
		convert:    uuidNullConvert,
		fieldType:  data.FieldTypeNullableString,
		matchRegex: matchRegexes["Optional<Uuid>"],
		scanType:   reflect.PointerTo(reflect.PointerTo(reflect.TypeOf(""))),
	},
	"Optional<Date>": {
		fieldType: data.FieldTypeNullableTime,
		scanType:  reflect.PointerTo(reflect.PointerTo(reflect.TypeOf(time.Time{}))),
	},
	"Optional<Datetime>": {
		fieldType: data.FieldTypeNullableTime,
		scanType:  reflect.PointerTo(reflect.PointerTo(reflect.TypeOf(time.Time{}))),
	},
	"Optional<Timestamp>": {
		fieldType: data.FieldTypeNullableTime,
		scanType:  reflect.PointerTo(reflect.PointerTo(reflect.TypeOf(time.Time{}))),
	},
	"Optional<Interval>": {
		convert:   intervalNullConvert,
		scanType:  reflect.PointerTo(reflect.PointerTo(reflect.TypeOf(int64(0)))),
		fieldType: data.FieldTypeNullableInt64,
	},
	"Optional<SmallSerial>": {
		scanType:  reflect.PointerTo(reflect.PointerTo(reflect.TypeOf(int16(0)))),
		fieldType: data.FieldTypeNullableInt16,
	},
	"Optional<Serial2>": {
		scanType:  reflect.PointerTo(reflect.PointerTo(reflect.TypeOf(int16(0)))),
		fieldType: data.FieldTypeNullableInt16,
	},
	"Optional<Serial>": {
		scanType:  reflect.PointerTo(reflect.PointerTo(reflect.TypeOf(int32(0)))),
		fieldType: data.FieldTypeNullableInt32,
	},
	"Optional<Serial4>": {
		scanType:  reflect.PointerTo(reflect.PointerTo(reflect.TypeOf(int32(0)))),
		fieldType: data.FieldTypeNullableInt32,
	},
	"Optional<Serial8>": {
		scanType:  reflect.PointerTo(reflect.PointerTo(reflect.TypeOf(int64(0)))),
		fieldType: data.FieldTypeNullableInt64,
	},
	"Optional<BigSerial>": {
		scanType:  reflect.PointerTo(reflect.PointerTo(reflect.TypeOf(int64(0)))),
		fieldType: data.FieldTypeNullableInt64,
	},
	"Optional<List>": {
		convert:    jsonConvert,
		fieldType:  data.FieldTypeJSON,
		matchRegex: matchRegexes["Optional<List>"],
		scanType:   reflect.PointerTo(reflect.TypeOf((*interface{})(nil)).Elem()),
	},
	"Optional<Dict>": {
		convert:    jsonConvert,
		fieldType:  data.FieldTypeJSON,
		matchRegex: matchRegexes["Optional<Dict>"],
		scanType:   reflect.PointerTo(reflect.TypeOf((*interface{})(nil)).Elem()),
	},
	"Optional<Set>": {
		convert:    jsonConvert,
		fieldType:  data.FieldTypeJSON,
		matchRegex: matchRegexes["Optional<Set>"],
		scanType:   reflect.PointerTo(reflect.TypeOf((*interface{})(nil)).Elem()),
	},
	"Optional<Tuple>": {
		convert:    jsonConvert,
		fieldType:  data.FieldTypeJSON,
		matchRegex: matchRegexes["Optional<Tuple>"],
		scanType:   reflect.PointerTo(reflect.TypeOf((*interface{})(nil)).Elem()),
	},
	"Optional<Struct>": {
		convert:    jsonConvert,
		fieldType:  data.FieldTypeJSON,
		matchRegex: matchRegexes["Optional<Struct>"],
		scanType:   reflect.PointerTo(reflect.TypeOf((*interface{})(nil)).Elem()),
	},
	"Optional<Variant>": {
		convert:    jsonConvert,
		fieldType:  data.FieldTypeJSON,
		matchRegex: matchRegexes["Optional<Variant>"],
		scanType:   reflect.PointerTo(reflect.TypeOf((*interface{})(nil)).Elem()),
	},
	"Optional<Enum>": {
		convert:    jsonConvert,
		fieldType:  data.FieldTypeJSON,
		matchRegex: matchRegexes["Optional<Enum>"],
		scanType:   reflect.PointerTo(reflect.TypeOf((*interface{})(nil)).Elem()),
	},
}

var YdbConverters = YDBConverters()

func YDBConverters() []sqlutil.Converter {
	var list []sqlutil.Converter
	for name, converter := range Converters {
		list = append(list, createConverter(name, converter))
	}
	return list
}

func GetConverter(columnType string) sqlutil.Converter {
	converter, ok := Converters[columnType]
	if ok {
		return createConverter(columnType, converter)
	}
	for name, converter := range Converters {
		if name == columnType {
			return createConverter(name, converter)
		}
		if converter.matchRegex != nil && converter.matchRegex.MatchString(columnType) {
			return createConverter(name, converter)
		}
	}
	return sqlutil.Converter{}
}

func createConverter(name string, converter Converter) sqlutil.Converter {
	convert := defaultConvert
	if converter.convert != nil {
		convert = converter.convert
	}
	return sqlutil.Converter{
		Name:           name,
		InputScanType:  converter.scanType,
		InputTypeRegex: converter.matchRegex,
		InputTypeName:  name,
		FrameConverter: sqlutil.FrameConverter{
			FieldType:     converter.fieldType,
			ConverterFunc: convert,
		},
	}
}

func jsonConvert(in interface{}) (interface{}, error) {
	if in == nil {
		return (*string)(nil), nil
	}
	jBytes, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}

	rawJSON := json.RawMessage(jBytes)
	return &rawJSON, nil
}

func uuidConvert(in interface{}) (interface{}, error) {
	if in == nil {
		return (*string)(nil), nil
	}
	v, ok := in.(*uuid.UUID)
	if !ok {
		return nil, fmt.Errorf("invalid uuid - %v", in)
	}
	f := (*v).String()
	return f, nil
}

func uuidNullConvert(in interface{}) (interface{}, error) {
	if in == nil {
		return (*string)(nil), nil
	}
	v, ok := in.(**uuid.UUID)
	if !ok {
		return nil, fmt.Errorf("invalid uuid - %v", in)
	}
	if *v == nil {
		return (*string)(nil), nil
	}
	f := (**v).String()
	return f, nil
}

func intervalConvert(in interface{}) (interface{}, error) {
	if in == nil {
		return int64(0), nil
	}
	v, ok := in.(*int64)
	if !ok {
		return nil, fmt.Errorf("invalid interval - %v", in)
	}
	return *v / 1000, nil
}

func intervalNullConvert(in interface{}) (interface{}, error) {
	if in == nil {
		return int64(0), nil
	}
	v, ok := in.(**int64)
	if !ok {
		return nil, fmt.Errorf("invalid interval - %v", in)
	}
	if *v == nil {
		return (*int64)(nil), nil
	}
	t := **v / 1000
	return &t, nil
}

func defaultConvert(in interface{}) (interface{}, error) {
	if in == nil {
		return reflect.Zero(reflect.TypeOf(in)).Interface(), nil
	}
	return reflect.ValueOf(in).Elem().Interface(), nil
}

func decimalConvert(in interface{}) (interface{}, error) {
	if in == nil {
		return float64(0), nil
	}
	v, ok := in.(*decimal.Decimal)
	if !ok {
		return nil, fmt.Errorf("invalid decimal - %v", in)
	}
	f, _ := (*v).Float64()
	return f, nil
}

func decimalNullConvert(in interface{}) (interface{}, error) {
	if in == nil {
		return float64(0), nil
	}
	v, ok := in.(**decimal.Decimal)
	if !ok {
		return nil, fmt.Errorf("invalid decimal - %v", in)
	}
	if *v == nil {
		return (*float64)(nil), nil
	}
	f, _ := (*v).Float64()
	return &f, nil
}
