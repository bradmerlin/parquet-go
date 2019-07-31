package go_parquet

import (
	"io"

	"github.com/fraugster/parquet-go/parquet"

	"github.com/pkg/errors"
)

type Int96 [12]byte

type int96PlainDecoder struct {
	r io.Reader
}

func (i *int96PlainDecoder) init(r io.Reader) error {
	i.r = r

	return nil
}

func (i *int96PlainDecoder) decodeValues(dst []interface{}) (int, error) {
	idx := 0
	for range dst {
		var data Int96
		// this one is a little tricky do not use ReadFull here
		n, err := i.r.Read(data[:])
		// make sure we handle the read data first then handle the error
		if n == 12 {
			dst[idx] = data
			idx++
		}

		if err != nil && (n == 0 || n == 12) {
			return idx, err
		}

		if err != nil {
			return idx, errors.Wrap(err, "not enough byte to read the Int96")
		}
	}
	return len(dst), nil
}

type int96PlainEncoder struct {
	w io.Writer
}

func (i *int96PlainEncoder) Close() error {
	return nil
}

func (i *int96PlainEncoder) init(w io.Writer) error {
	i.w = w

	return nil
}

func (i *int96PlainEncoder) encodeValues(values []interface{}) error {
	data := make([]byte, len(values)*12)
	for j := range values {
		i96 := values[j].(Int96)
		copy(data[j*12:], i96[:])
	}

	return writeFull(i.w, data)
}

type int96Store struct {
	byteArrayStore
}

func (*int96Store) sizeOf(v interface{}) int {
	return 12
}

func (is *int96Store) parquetType() parquet.Type {
	return parquet.Type_INT96
}

func (is *int96Store) typeLen() *int32 {
	return nil
}

func (is *int96Store) repetitionType() parquet.FieldRepetitionType {
	return is.repTyp
}

func (is *int96Store) convertedType() *parquet.ConvertedType {
	return nil
}

func (is *int96Store) scale() *int32 {
	return nil
}

func (is *int96Store) precision() *int32 {
	return nil
}

func (is *int96Store) logicalType() *parquet.LogicalType {
	return nil
}

func (is *int96Store) getValues(v interface{}) ([]interface{}, error) {
	var vals []interface{}
	switch typed := v.(type) {
	case Int96:
		is.setMinMax(typed[:])
		vals = []interface{}{typed}
	case []Int96:
		if is.repTyp != parquet.FieldRepetitionType_REPEATED {
			return nil, errors.Errorf("the value is not repeated but it is an array")
		}
		vals = make([]interface{}, len(typed))
		for j := range typed {
			is.setMinMax(typed[j][:])
			vals[j] = typed[j]
		}
	default:
		return nil, errors.Errorf("unsupported type for storing in Int96 column %T => %+v", v, v)
	}

	return vals, nil
}
